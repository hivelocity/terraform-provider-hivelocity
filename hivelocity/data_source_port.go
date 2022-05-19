package hivelocity

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
	"sort"
	"strings"
)

// dataSourceDevicePort is a resource for retrieving a device's port.
// Because there is more than one port per device, the filter used must match down to one port, or error otherwise. If
// the filter is omitted entirely, the ports are automatically filtered to the private interface which is typically
// used to apply private VLANs.
func dataSourceDevicePort() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDevicePortRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"first":  dataSourceFilterFirstSchema(),
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
			"is_bond": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceDevicePortRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	deviceId := d.Get("device_id").(int)

	ports, _, err := hv.client.DeviceApi.GetDevicePortResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /device/%v/ports failed! (%s)\n\n %s", deviceId, err, myErr.Body())
	}

	if len(ports) == 0 {
		return diag.Errorf("Device %v has no ports to query", deviceId)
	}

	filterablePorts := make([]map[string]interface{}, len(ports))
	hasBonds := false

	for i, p := range ports {
		isBond := strings.HasPrefix(p.Name, "bond")
		filterablePorts[i] = map[string]interface{}{
			"private":   p.Private,
			"name":      p.Name,
			"port_id":   p.PortId,
			"device_id": p.DeviceId,
			"is_bond":   isBond,
		}
		if isBond {
			hasBonds = true
		}
	}

	// Sort ports by name, in case user uses a filter and a device with more than 2 ports
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

	// Select default filter based on
	var defaultFilter []filter
	if hasBonds {
		defaultFilter = []filter{{name: "is_bond", values: []string{"true"}}}
	} else {
		defaultFilter = []filter{
			{name: "private", values: []string{"true"}},
			{name: "name", values: []string{"eth1"}},
		}
	}

	matchedPort, err := doFiltering(d, filterablePorts, defaultFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	if matchedPort == nil {
		return diag.Errorf("No ports found with given filters")
	}

	for k, v := range matchedPort {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprint(matchedPort["port_id"]))

	return nil
}
