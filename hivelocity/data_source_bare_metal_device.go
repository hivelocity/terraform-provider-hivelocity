package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBareMetalDevice() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBareMetalDeviceRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"first":  dataSourceFilterFirstSchema(),
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"device_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"order_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"service_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"product_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"product_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"os_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"location_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"power_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"vlan_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			"primary_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ssh_key_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			"script": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceBareMetalDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	bareMetalDeviceInfo, _, err := hv.client.BareMetalDevicesApi.GetBareMetalDeviceResource(hv.auth, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonProductInfo, err := json.Marshal(bareMetalDeviceInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	bareMetalDevices := make([]map[string]interface{}, 0)
	err = json.Unmarshal(jsonProductInfo, &bareMetalDevices)
	if err != nil {
		return diag.FromErr(err)
	}

	bareMetalDevices = convertKeysOfList(bareMetalDevices)

	bareMetalDevice, err := doFiltering(d, bareMetalDevices)
	if err != nil {
		return diag.FromErr(err)
	}

	if bareMetalDevice == nil {
		return nil
	}

	for k, v := range bareMetalDevice {
		d.Set(k, v)
	}

	d.SetId(fmt.Sprint(bareMetalDevice["device_id"]))

	return nil
}
