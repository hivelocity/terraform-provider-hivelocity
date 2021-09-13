package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func buildDeviceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"first":  dataSourceFilterFirstSchema(),
		"device_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"device_type": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"hostname": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_addresses": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"ipmi_address": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"ipmi_enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Computed: true,
		},
		"location": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"metadata": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"power_status": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"service_plan": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"status": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func dataSourceDevice() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceRead,
		Schema:      buildDeviceSchema(),
	}
}

func dataSourceDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hv, _ := m.(*Client)

	hivelocityDevices, _, err := hv.client.DeviceApi.GetDeviceResource(hv.auth, nil)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /device failed! (%s)\n\n %s", err, myErr.Body())
	}

	jsonDevices, err := json.Marshal(hivelocityDevices)
	if err != nil {
		return diag.FromErr(err)
	}

	devices := make([]map[string]interface{}, 0)
	err = json.Unmarshal(jsonDevices, &devices)
	if err != nil {
		return diag.FromErr(err)
	}

	devices = convertKeysOfList(devices)
	devices = forceValuesToStringOfList(devices, "location")
	devices = forceValuesToStringOfList(devices, "metadata")
	devices = filterNonSchemaKeysForList(devices, buildDeviceSchema())

	device, err := doFiltering(d, devices)
	if err != nil {
		return diag.FromErr(err)
	}

	if device == nil {
		return nil
	}

	for k, v := range device {
		d.Set(k, v)
	}

	d.SetId(fmt.Sprint(device["device_id"]))

	return diags
}
