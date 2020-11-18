package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func buildDeviceInitialCredsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"device_id": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"user": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"locker_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"port": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"password": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceDeviceInitialCreds() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceInitialCredsRead,
		Schema:      buildDeviceInitialCredsSchema(),
	}
}

func dataSourceDeviceInitialCredsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hv, _ := m.(*Client)

	deviceID := int32(d.Get("device_id").(int))
	hivelocityInitialCreds, _, err := hv.client.DeviceApi.GetInitialCredsIdResource(hv.auth, deviceID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonInitialCreds, err := json.Marshal(hivelocityInitialCreds)
	if err != nil {
		return diag.FromErr(err)
	}

	initialCreds := make(map[string]interface{}, 0)
	err = json.Unmarshal(jsonInitialCreds, &initialCreds)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range convertKeys(initialCreds) {
		d.Set(k, v)
	}

	d.SetId(fmt.Sprint(deviceID))

	return diags
}
