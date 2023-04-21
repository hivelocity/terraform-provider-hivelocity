package hivelocity

import (
	"context"
	"fmt"
	"github.com/hivelocity/terraform-provider-hivelocity/hivelocity/pkg/mod/github.com/hashicorp/terraform-plugin-sdk/v2@v2.17.0/helper/resource"
	"log"
	"strconv"
	"time"

	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
	"github.com/hivelocity/terraform-provider-hivelocity/hivelocity/pkg/mod/github.com/hashicorp/terraform-plugin-sdk/v2@v2.17.0/diag"
	"github.com/hivelocity/terraform-provider-hivelocity/hivelocity/pkg/mod/github.com/hashicorp/terraform-plugin-sdk/v2@v2.17.0/helper/schema"
)

func resourceVLAN() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceVLANCreate,
		ReadContext:   resourceVLANRead,
		UpdateContext: resourceVLANUpdate,
		DeleteContext: resourceVLANDelete,
		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Description: "Type of VLAN to be created, can be either `private` or `public`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"facility_code": &schema.Schema{
				Description: "Location where to create this VLAN",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ip_ids": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
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

// resourceVLANRead reads a VLAN
func resourceVLANRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	hv, _ := m.(*Client)

	log.Printf("[INFO] Reading VLAN %s", d.Id())

	vlanId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vlan, resp, err := hv.client.VLANApi.GetVlanIdResource(hv.auth, int32(vlanId), nil)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] VLAN %s not found", d.Id())
			d.SetId("")
			return diags
		}
		return diag.FromErr(formatSwaggerError(err, "GET /vlan/%d", vlanId))
	}

	d.Set("q_in_q", vlan.QInQ)
	d.Set("port_ids", SetFromInt32List(vlan.PortIds))
	d.Set("ip_ids", SetFromInt32List(vlan.IpIds))
	d.Set("automated", vlan.Automated)
	d.Set("vlan_tag", vlan.VlanTag)
	d.Set("facility_code", vlan.FacilityCode)
	d.Set("type", vlan.Type_)
	d.Set("vlan_id", vlan.VlanId)

	return diags
}

// resourceVLANCreate creates a VLAN
func resourceVLANCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	hv, _ := m.(*Client)
	log.Printf("[INFO] Creating")
	payload := swagger.VlanCreate{
		FacilityCode: d.Get("facility_code").(string),
		Type_:        d.Get("type").(string),
	}

	vlan, _, err := hv.client.VLANApi.PostVlanResource(hv.auth, payload, nil)
	if err != nil {
		d.SetId("")
		return diag.FromErr(formatSwaggerError(err, "POST /vlan"))
	}

	if diags.HasError() {
		d.SetId(fmt.Sprint(vlan.VlanId))
		diags = append(diags, resourceVLANRead(ctx, d, m)...)
		d.SetId("")
		return diags
	}
	log.Printf("[INFO] Created VLAN ID: %d", vlan.VlanId)
	d.SetId(fmt.Sprint(vlan.VlanId))
	return resourceVLANRead(ctx, d, m)
}

// resourceVLANUpdate updates a VLAN
func resourceVLANUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var payload swagger.VlanUpdate
	hv, _ := m.(*Client)
	_vlanId, err := strconv.Atoi(d.Id())

	log.Printf("[INFO] Updating VLAN ID: %s", d.Id())

	if err != nil {
		return diag.FromErr(err)
	}
	vlanId := int32(_vlanId)

	ipIds := make([]int32, 0)
	portIds := make([]int32, 0)

	if portIdSet, ok := d.GetOk("port_ids"); ok {
		portIds = Int32ListFromSet(portIdSet.(*schema.Set))
		payload.IpIds = portIds
	}

	if ipIdSet, ok := d.GetOk("ip_ids"); ok {
		ipIds = Int32ListFromSet(ipIdSet.(*schema.Set))
		payload.IpIds = ipIds
	}

	// Update VLAN
	task, _, err := hv.client.VLANApi.PutVlanIdResource(hv.auth, vlanId, payload, nil)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /vlan/%d: %s\n\n %s", fmt.Sprint(vlanId), err, myErr.Body())
	}

	//wait for task to complete
	if _, err := waitForNetworkTaskByClient(hv.auth, timeout, hv, task.TaskId); err != nil {
		return diag.FromErr(err)
	}

	return resourceVLANRead(ctx, d, m)
}

// resourceVLANDelete deletes a VLAN
func resourceVLANDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	hv, _ := m.(*Client)

	log.Printf("[INFO] Deleting VLAN ID: %s", d.Id())

	var vlanId int32
	if vlanId_, err := strconv.Atoi(d.Id()); err != nil {
		return diag.FromErr(err)
	} else {
		vlanId = int32(vlanId_)
	}

	resp, err := hv.client.VLANApi.DeleteVlanIdResource(hv.auth, vlanId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
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
			task, _, err := hv.client.NetworkApi.GetNetworkTaskIdResource(ctx, taskId, nil)
			if err != nil {
				return nil, "", formatSwaggerError(err, "network/status/%s", taskId)
			}

			return &task, task.Result, nil
		},
		Timeout:                   timeout,
		Delay:                     10 * time.Second,
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 1,
	}

	r, err := stateChangeConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.(*swagger.NetworkTaskDump), err
}
