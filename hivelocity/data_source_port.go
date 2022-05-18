package hivelocity

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
	"sort"
	"strconv"
	"strings"
)

// dataSourceDevicePort is a resource for retrieving a device's port.
// Because there is more than one port per device, the filter used must match down to one port, or error otherwise. If
// the filter is omitted entirely, the ports are automatically filtered to the primary interface.
// If this resource is used to setup VLANs, you mostly likely will just want to omit the filter attribute to default
// to the primary port.
func dataSourceDevicePort() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDevicePortRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			// Primary is computed by sorting port names ascending ("bond*" first if any) and picking the first one
			"is_primary": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"private": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"device_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func dataSourceDevicePortRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ports, _, err := hv.client.DeviceApi.GetDevicePortResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /device/%v/ports failed! (%s)\n\n %s", deviceId, err, myErr.Body())
	}

	if len(ports) == 0 {
		return diag.Errorf("Device %v has no ports to query", deviceId)
	}

	filterablePorts := make([]map[string]interface{}, len(ports))

	for i, p := range ports {
		filterablePorts[i] = map[string]interface{}{
			"is_primary": false,
			"private":    p.Private,
			"name":       p.Name,
			"port_id":    p.PortId,
			"device_id":  p.DeviceId,
		}
	}

	// Sort ports by name, preferring bonds first
	sort.SliceStable(filterablePorts, func(i, j int) bool {
		a, b := filterablePorts[i]["name"].(string), filterablePorts[j]["name"].(string)
		if strings.HasPrefix(a, "bond") {
			if strings.HasPrefix(b, "bond") {
				return a < b
			} else {
				return true
			}
		}
		return a < b
	})

	// Set "is_primary" on first
	filterablePorts[0]["is_primary"] = "true"

	matchedPort, err := doFiltering(d, filterablePorts, []filter{
		{"is_primary", []string{"true"}},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range matchedPort {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("%v", matchedPort["port_id"]))

	return nil
}
