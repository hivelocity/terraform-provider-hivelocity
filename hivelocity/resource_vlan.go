package hivelocity

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"vlan_tag": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "The VLAN Tag id from the switch. Use this value when configuring" +
					" your OS interfaces to use the VLAN.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "If `public`, this VLAN can have IPs assigned to become reachable " +
					" from the internet. If `private`, this VLAN can not have IPs assigned and " +
					" will never be reachabled from the internet. All VLANs are subject to traffic " +
					" billing and overages, with the exception of private VLAN traffic on unbonded " +
					" Devices.",
			},
			"ip_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Unique IDs of IP Assignments.",
			},
			"automated": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, VLAN can be automated via API. If false, contact support to enable automation.",
			},
			"facility_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "For example: `NYC1`.",
			},
			"q_in_q": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "If true, VLAN is configured in Q-in-Q Mode. Automation is not" +
					" currently supported on Q-in-Q VLANs.",
			},
			"port_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Unique IDs of ports or bonds.",
			},
		},
	}
}

func resourceVlanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	hv, _ := m.(*Client)

	log.Printf("[INFO] Creating VLAN")
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
		diags = append(diags, resourceVlanDelete(ctx, d, m)...)
		d.SetId("")
		return diags
	}

	log.Printf("[INFO] Created VLAN ID: %d", vlan.VlanId)
	d.SetId(fmt.Sprint(vlan.VlanId))
	return resourceVlanRead(ctx, d, m)

}

func resourceVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceVlanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var payload swagger.VlanUpdate
	hv, _ := m.(*Client)
	_vlanId, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}
	vlanId := int32(_vlanId)

	log.Printf("[INFO] Updating VLAN ID: %s", d.Id())

	// Construct payload for VLAN update
	if portIdSet, ok := d.GetOk("port_ids"); ok {
		portIds := Int32ListFromSet(portIdSet.(*schema.Set))
		payload.PortIds = portIds
	}

	if ipIdSet, ok := d.GetOk("ip_ids"); ok {
		ipIds := Int32ListFromSet(ipIdSet.(*schema.Set))
		payload.IpIds = ipIds
	}

	// Update VLAN
	task, _, err := hv.client.VLANApi.PutVlanIdResource(hv.auth, vlanId, payload, nil)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /vlan/%d: %s\n\n %s", vlanId, err, myErr.Body())
	}

	// Wait for task to complete
	if _, err := waitForNetworkTaskByClient(hv.auth, d.Timeout(schema.TimeoutUpdate), hv, task.TaskId); err != nil {
		return diag.FromErr(err)
	}

	return resourceVlanRead(ctx, d, m)
}

func resourceVlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
