package hivelocity

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceBond() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceBondCreate,
		UpdateContext: resourceBondUpdate,
		ReadContext:   resourceBondRead,
		DeleteContext: resourceBondDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the device to create the bond on. ",
			},
		},
	}
}

func resourceBondCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	hv, _ := m.(*Client)
	deviceId, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating Bond for device %d", deviceId)
	task, _, err := hv.client.DeviceApi.PostDeviceBondResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		return diag.FromErr(formatSwaggerError(err, "POST /device/{deviceId}/ports/bond"))
	}

	log.Printf("[INFO] Waiting for task %s to complete", task.TaskId)
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
			log.Printf("[INFO] Bond for device %d not found", deviceId)
			return diags
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /device/%d/ports/bond failed! (%s)\n\n %s", deviceId, err, myErr.Body())
	}

	// wait for task to complete
	if _, err := waitForNetworkTaskByClient(hv.auth, d.Timeout(schema.TimeoutDelete), hv, task.TaskId); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Bond deleted for device %d", deviceId)
	d.SetId("")
	return diags
}

func resourceBondUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceBondRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
