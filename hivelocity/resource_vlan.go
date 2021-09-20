package hivelocity

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
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
		Schema: map[string]*schema.Schema{
			"device_ids": &schema.Schema{
				Description: "IDs of devices to include in this VLAN",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Required: true,
			},
		},
	}
}

func resourceVlanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := makeVlanCreatePayload(d)

	vlan, _, err := hv.client.VLANApi.PostVlanResource(hv.auth, payload, nil)
	if err != nil {
		d.SetId("")
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("POST /vlan failed! (%s)\n\n %s", err, myErr.Body())
	}

	_, err = waitForVlan(ctx, d, hv, vlan.VlanId, payload.DeviceIds)
	if err != nil {
		d.SetId("")
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("Error creating VLAN ID (%d). The Hivelocity team will investigate: (%s)\n\n %s", vlan.VlanId, err, myErr.Body())
	}

	log.Printf("[INFO] Created VLAN ID: %d", vlan.VlanId)

	d.SetId(fmt.Sprint(vlan.VlanId))

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
		if response.StatusCode == 404 {
			log.Printf("[WARN] Vlan ID not found: (%s)", d.Id())
			d.SetId("")
			return diags
		}
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /vlan failed! (%s)\n\n %s", err, myErr.Body())
	}

	deviceIds := getVlanDeviceIds(&vlan)
	d.Set("device_ids", deviceIds)

	return diags
}

func resourceVlanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Updating VLAN ID: %s", d.Id())

	vlanId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload := makeVlanCreatePayload(d)

	_, _, err = hv.client.VLANApi.PutVlanIdResource(hv.auth, int32(vlanId), payload, nil)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /vlan/%d failed! (%s)\n\n %s", vlanId, err, myErr.Body())
	}

	_, err = waitForVlan(ctx, d, hv, int32(vlanId), payload.DeviceIds)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("Error updating VLAN ID (%d). The Hivelocity team will investigate: (%s)\n\n %s", vlanId, err, myErr.Body())
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
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
		if response.StatusCode == 404 {
			log.Printf("[WARN] Vlan ID not found: (%s)", d.Id())
			d.SetId("")
			return diags
		}
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /vlan/%d failed! (%s)\n\n %s", vlanId, err, myErr.Body())
	}

	_, err = waitForVlanDeletion(ctx, d, hv, vlanId)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("Error deleting VLAN ID (%d). The Hivelocity team will investigate: (%s)\n\n %s", vlanId, err, myErr.Body())
	}

	// Delete resource from state
	d.SetId("")

	return diags
}

func makeVlanCreatePayload(d *schema.ResourceData) swagger.VlanCreate {
	deviceIds := make([]int32, d.Get("device_ids.#").(int))
	for i, id := range d.Get("device_ids").([]interface{}) {
		deviceIds[i] = int32(id.(int))
	}

	payload := swagger.VlanCreate{
		DeviceIds: deviceIds,
	}

	return payload
}

func getVlanDeviceIds(vlan *swagger.Vlan) []int {
	deviceIds := make([]int, len(vlan.DeviceIds))
	for i, id := range vlan.DeviceIds {
		deviceIds[i] = int(id)
	}

	return deviceIds
}

func arraysEqual(a []int32, b []int32) bool {
	sort.SliceStable(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	sort.SliceStable(b, func(i, j int) bool {
		return b[i] < b[j]
	})

	return reflect.DeepEqual(a, b)
}

// Wait for VLAN to be created or updated by polling the pendingDevices list until empty
func waitForVlan(ctx context.Context, d *schema.ResourceData, hv *Client, vlanId int32, deviceIds []int32) (interface{}, error) {
	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{
			"pending",
		},
		Target: []string{
			"complete",
		},
		Refresh: func() (interface{}, string, error) {
			vlan, _, err := hv.client.VLANApi.GetVlanIdResource(hv.auth, vlanId, nil)
			if err != nil {
				return 0, "", err
			}
			if len(vlan.PendingDevices) == 0 && arraysEqual(vlan.DeviceIds, deviceIds) {
				return vlan, "complete", nil
			}
			return vlan, "pending", nil
		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     30 * time.Second,
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 1,
	}
	return stateChangeConf.WaitForStateContext(ctx)
}

// Wait for VLAN to be deleted by polling its status code until a 404
func waitForVlanDeletion(ctx context.Context, d *schema.ResourceData, hv *Client, vlanId int) (interface{}, error) {
	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{
			"pending",
		},
		Target: []string{
			"complete",
		},
		Refresh: func() (interface{}, string, error) {
			vlan, response, err := hv.client.VLANApi.GetVlanIdResource(hv.auth, int32(vlanId), nil)
			if err == nil {
				return vlan, "pending", nil
			} else if response.StatusCode == 404 {
				return vlan, "complete", nil
			}

			return 0, "", err
		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     30 * time.Second,
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 1,
	}
	return stateChangeConf.WaitForStateContext(ctx)
}
