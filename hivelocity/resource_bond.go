package hivelocity

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
	"log"
	"strconv"
)

func resourceBond() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBondCreate,
		DeleteContext: resourceBondDelete,

		Schema: map[string]*schema.Schema{
			"device_id": {
				Description: "Device ID",
				Type:        schema.TypeInt,
				Required:    true,
			},
		},
	}
}

func resourceBondCreate(d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	hv, _ := m.(*Client)
	deviceId, err := strconv.Atoi(d.Id())

	log.Printf("[INFO] Creating Bond for device %d", deviceId)

	if err != nil {
		return diag.FromErr(err)
	}

	// Creating a bond
	task, _, err := hv.client.DeviceApi.PostDeviceBondResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		return diag.FromErr(formatSwaggerError(err, "POST /device/{deviceId}/ports/bond"))
	}

	log.Printf("[INFO] Waiting for task %d to complete", task.TaskId)
	// Wait for task to complete
	if _, err := waitForNetworkTaskByClient(hv.auth, d.Timeout(schema.TimeoutCreate), hv, task.TaskId); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Bond created for device %d", deviceId)
	return diags
}

func resourceBondDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting bond for device %d", deviceId)
	task, resp, err := hv.client.DeviceApi.DeleteDeviceBondResource(hv.auth, int32(deviceId), nil)

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /device/{deviceId}/ports/bond returned %d, expected 202", fmt.Sprint(deviceId), err, myErr.Body())
	}

	// wait for task to complete
	if _, err := waitForNetworkTaskByClient(hv.auth, d.Timeout(schema.TimeoutDelete), hv, task.TaskId); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Bond deleted for device %d", deviceId)
	d.SetId("")
	return diags
}
