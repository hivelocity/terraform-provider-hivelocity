package hivelocity

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
	"log"
	"strconv"
	"time"
)

func resourceVlan() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceVlanCreate,
		ReadContext:   resourceVlanRead,
		UpdateContext: resourceVlanUpdate,
		DeleteContext: resourceVlanDelete,
		Schema: map[string]*schema.Schema{
			"facility_code": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_networking_only": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"port_ids": &schema.Schema{
				Description: "IDs of ports to include in this VLAN",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Required: true,
			},
		},
	}
}

func resourceVlanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	hv, _ := m.(*Client)

	payload := makeVlanCreatePayload(d)

	vlan, _, err := hv.client.VLANApi.PostVlanResource(hv.auth, payload, nil)
	if err != nil {
		d.SetId("")
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("POST /vlan failed! (%s)\n\n %s", err, myErr.Body())
	}

	// Update ports
	if len(vlan.PortIds) > 0 {
		diags = append(diags, _updateVlanPorts(ctx, hv, d, vlan.Id)...)
	}

	// If any errors happened, delete VLAN
	if diags.HasError() {
		// Set ID for delete to run
		d.SetId(fmt.Sprint(vlan.Id))
		for _, d := range resourceVlanDelete(ctx, d, m) {
			diags = append(diags, d)
		}
		d.SetId("")

		return diags
	}

	log.Printf("[INFO] Created VLAN ID: %d", vlan.Id)
	d.SetId(fmt.Sprint(vlan.Id))

	return resourceVlanRead(ctx, d, m)
}

func resourceVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	hv, _ := m.(*Client)

	vlanId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vlan, response, err := hv.client.VLANApi.GetVlanIdResource(hv.auth, int32(vlanId), nil)
	if err != nil {
		// If resource was deleted outside terraform, remove it from state and exit gracefully
		if response != nil && response.StatusCode == 404 {
			log.Printf("[WARN] Vlan ID not found: (%s)", d.Id())
			d.SetId("")
			return diags
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /vlan failed! (%s)\n\n %s", err, myErr.Body())
	}

	valuesToSet := map[string]interface{}{
		"port_ids":                SetFromInt32List(vlan.PortIds),
		"facility_code":           vlan.FacilityCode,
		"private_networking_only": vlan.PrivateNetworkingOnly,
	}

	for k, v := range valuesToSet {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func _updateVlanPorts(
	ctx context.Context,
	hv *Client,
	d *schema.ResourceData,
	vlanId int32,
) diag.Diagnostics {
	// Update ports
	updatePayload := makeUpdatePayload(d)
	task, _, err := hv.client.VLANApi.PutVlanIdResource(hv.auth, vlanId, updatePayload, nil)

	// Append errors if PUT or network task fails
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /vlan failed! (%s)\n\n %s", err, myErr.Body())
	}

	// Wait for task to finish
	_, err = waitForNetworkTaskByClient(hv.auth, d.Timeout(schema.TimeoutCreate), hv, task.TaskId)
	if err != nil {
		return diag.Errorf("Waiting for vlan update failed! (%s)", err)
	}

	return diag.Diagnostics{}
}

func resourceVlanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Updating VLAN ID: %s", d.Id())

	_vlanId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	vlanId := int32(_vlanId)

	diags := _updateVlanPorts(ctx, hv, d, vlanId)
	if diags.HasError() {
		return diags
	}

	return resourceVlanRead(ctx, d, m)
}

func resourceVlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	hv, _ := m.(*Client)

	log.Printf("[INFO] Deleting VLAN ID: %s", d.Id())

	vlanId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, response, err := hv.client.VLANApi.DeleteVlanIdResource(hv.auth, int32(vlanId), nil)
	if err != nil {
		// If resource was deleted outside terraform, remove it from state and exit gracefully
		if response != nil && response.StatusCode == 404 {
			log.Printf("[WARN] Vlan ID not found: (%s)", d.Id())
			d.SetId("")
			return diags
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /vlan/%d failed! (%s)\n\n %s", vlanId, err, myErr.Body())
	}

	// Delete resource from state
	d.SetId("")

	return diags
}

func makeVlanCreatePayload(d *schema.ResourceData) swagger.VlanCreate {
	payload := swagger.VlanCreate{
		FacilityCode:          d.Get("facility_code").(string),
		PrivateNetworkingOnly: d.Get("private_networking_only").(bool),
	}

	return payload
}

func Int32ListFromSet(s *schema.Set) []int32 {
	int32List := make([]int32, s.Len())
	for i, n := range s.List() {
		int32List[i] = int32(n.(int))
	}
	return int32List
}

func SetFromInt32List(items []int32) *schema.Set {
	intList := make([]interface{}, len(items))
	for i, n := range items {
		intList[i] = int(n)
	}
	return schema.NewSet(schema.HashInt, intList)
}

func makeUpdatePayload(d *schema.ResourceData) swagger.VlanUpdate {
	portIds := make([]int32, 0)

	if portIdSet, ok := d.GetOk("port_ids"); ok {
		portIds = Int32ListFromSet(portIdSet.(*schema.Set))
	}

	return swagger.VlanUpdate{
		PortIds: portIds,
	}
}

func waitForNetworkTaskByClient(
	ctx context.Context,
	timeout time.Duration,
	hv *Client,
	taskId string,
) (*swagger.NetworkTaskDump, error) {
	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{"Pending", ""},
		Target:  []string{"Success"},
		Refresh: func() (interface{}, string, error) {
			var matchedTask *swagger.NetworkTaskDump
			tasks, _, err := hv.client.NetworkApi.GetNetworkTaskClientResource(ctx, nil)
			if err != nil {
				return nil, "", err
			}

			for _, task := range tasks {
				if task.TaskId == taskId {
					matchedTask = &task
				}
			}

			if matchedTask == nil {
				return nil, "Failed", errors.New("could not find network task")
			}

			return matchedTask, matchedTask.Result, nil
		},
		Timeout:                   timeout,
		Delay:                     30 * time.Second,
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 1,
	}

	r, err := stateChangeConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.(*swagger.NetworkTaskDump), err
}
