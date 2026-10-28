package hivelocity

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

// ipAssignmentAttributes are the read-only attributes of an IP assignment.
func ipAssignmentAttributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"assignment_id":   computedInt("Unique ID of the assignment."),
		"client_id":       computedInt("Owning client ID."),
		"version":         computedInt("IP version (4 or 6)."),
		"subnet":          computedString("CIDR representation of the assignment."),
		"netmask":         computedString("Netmask of the assignment."),
		"broadcast_ip":    computedString("Broadcast address of the subnet."),
		"gateway_ip":      computedString("Gateway address of the subnet."),
		"first_usable_ip": computedString("First address usable by devices."),
		"last_usable_ip":  computedString("Last address usable by devices."),
		"usable_ips":      computedStringList("Usable addresses (IPv4 only)."),
		"dns_resolvers":   computedStringList("DNS resolvers for the subnet (VPS only)."),
		"facility_code":   computedString("Facility of the assignment."),
		"description":     computedString("Description of the assignment."),
		"device_id":       computedInt("Device receiving the assignment's traffic, if any."),
		"port_id":         computedInt("Port receiving the assignment's traffic, if any."),
		"vlan_id":         computedInt("VLAN receiving the assignment's traffic, if any."),
		"next_hop_ip":     computedString("Next-hop IP when the assignment is statically routed."),
	}
}

func flattenIpAssignment(ip swagger.IpAssignment) map[string]interface{} {
	return map[string]interface{}{
		"assignment_id":   int(ip.AssignmentId),
		"client_id":       int(ip.ClientId),
		"version":         int(ip.Version),
		"subnet":          ip.Subnet,
		"netmask":         ip.Netmask,
		"broadcast_ip":    ip.BroadcastIp,
		"gateway_ip":      ip.GatewayIp,
		"first_usable_ip": ip.FirstUsableIp,
		"last_usable_ip":  ip.LastUsableIp,
		"usable_ips":      flattenStringList(ip.UsableIps),
		"dns_resolvers":   flattenStringList(ip.DnsResolvers),
		"facility_code":   ip.FacilityCode,
		"description":     ip.Description,
		"device_id":       int(ip.DeviceId),
		"port_id":         int(ip.PortId),
		"vlan_id":         int(ip.VlanId),
		"next_hop_ip":     ip.NextHopIp,
	}
}

func dataSourceIpAssignment() *schema.Resource {
	s := ipAssignmentAttributes()
	s["assignment_id"] = &schema.Schema{
		Description:   "Look up the assignment by ID.",
		Type:          schema.TypeInt,
		Optional:      true,
		Computed:      true,
		ConflictsWith: []string{"ip_address", "filter"},
	}
	s["ip_address"] = &schema.Schema{
		Description:   "Look up the assignment containing this IP address.",
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"assignment_id", "filter"},
	}
	s["filter"] = dataSourceFiltersSchema()
	s["first"] = dataSourceFilterFirstSchema()

	return &schema.Resource{
		Description: "Looks up one of your IP assignments (subnets) by ID, by an address it contains, or " +
			"by filtering the full list (for example on `facility_code` or `subnet`).",
		ReadContext: dataSourceIpAssignmentRead,
		Schema:      s,
	}
}

func dataSourceIpAssignmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	if id, ok := d.GetOk("assignment_id"); ok {
		ip, _, err := hv.client.IPAssignmentApi.GetIpAssignmentIdResource(hv.auth, int32(id.(int)), nil)
		if err != nil {
			return apiErrorDiag(err, "GET /ip/%d", id.(int))
		}
		return setIpAssignment(d, ip)
	}

	if addr, ok := d.GetOk("ip_address"); ok {
		ip, _, err := hv.client.IPAssignmentApi.GetIpAssignmentSubnetResource(hv.auth, addr.(string), nil)
		if err != nil {
			return apiErrorDiag(err, "GET /ip/address/%s", addr.(string))
		}
		return setIpAssignment(d, ip)
	}

	opts := &swagger.IPAssignmentApiGetIpAssignmentResourceOpts{
		DisplayIPv6:            optional.NewBool(true),
		DisplayPrivateNetworks: optional.NewBool(true),
	}
	ips, _, err := hv.client.IPAssignmentApi.GetIpAssignmentResource(hv.auth, opts)
	if err != nil {
		return apiErrorDiag(err, "GET /ip")
	}

	return readFilteredList(d, ips, ipAssignmentAttributes(), "assignment_id", nil)
}

func setIpAssignment(d *schema.ResourceData, ip swagger.IpAssignment) diag.Diagnostics {
	if err := setResourceData(d, flattenIpAssignment(ip)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprint(ip.AssignmentId))
	return nil
}

func dataSourceDeviceIps() *schema.Resource {
	return &schema.Resource{
		Description: "Lists the IP assignments routed to a device.",
		ReadContext: dataSourceDeviceIpsRead,
		Schema: map[string]*schema.Schema{
			"device_id": {
				Description: "ID of the device.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"ip_assignments": {
				Description: "IP assignments of the device.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource{Schema: ipAssignmentAttributes()},
			},
		},
	}
}

func dataSourceDeviceIpsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	deviceId := d.Get("device_id").(int)

	ips, _, err := hv.client.DeviceApi.GetDeviceIpAssignmentsResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		return apiErrorDiag(err, "GET /device/%d/ips", deviceId)
	}

	list := make([]interface{}, 0, len(ips))
	for _, ip := range ips {
		list = append(list, flattenIpAssignment(ip))
	}

	if err := d.Set("ip_assignments", list); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprint(deviceId))
	return nil
}
