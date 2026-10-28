package hivelocity

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceDevicePort() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the network configuration of a single device port (or bond): whether it is " +
			"enabled, its native VLAN and the IP assignments routed to it. Only attributes set in the " +
			"configuration are changed; unset attributes keep their current value. Removing the resource " +
			"from Terraform leaves the port configured as-is (it does not clear the port, which could " +
			"disconnect the device).",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceDevicePortCreate,
		ReadContext:   resourceDevicePortRead,
		UpdateContext: resourceDevicePortUpdate,
		DeleteContext: resourceDevicePortDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDevicePortImport,
		},
		Schema: map[string]*schema.Schema{
			"device_id": {
				Description: "ID of the device the port belongs to.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"port_id": {
				Description: "ID of the port (see the `hivelocity_device_port` data source).",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"enabled": {
				Description: "Whether the port is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"native_vlan_id": {
				Description: "ID of the VLAN used as the port's native (untagged) VLAN. Set to `0` to remove it.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"ip_assignment_ids": {
				Description: "IDs of the IP assignments routed to this port.",
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Computed:    true,
			},
			"name":            computedString("Interface name, for example `eth1` or `bond0`."),
			"type":            computedString("Port type."),
			"status":          computedString("Port status as reported by the API."),
			"private":         computedBool("True for private (back-end) ports."),
			"mtu":             computedInt("MTU of the port."),
			"native_vlan_tag": computedInt("VLAN tag of the native VLAN, if any."),
		},
	}
}

func devicePortId(deviceId, portId int) string {
	return fmt.Sprintf("%d:%d", deviceId, portId)
}

// findDevicePort returns a single port of a device, or nil when it does not exist.
func (hv *Client) findDevicePort(deviceId, portId int32) (*swagger.DevicePort, error) {
	ports, _, err := hv.client.DeviceApi.GetDevicePortResource(hv.auth, deviceId, nil)
	if err != nil {
		return nil, formatSwaggerError(err, "GET /device/%d/ports", deviceId)
	}
	for i := range ports {
		if ports[i].PortId == portId {
			return &ports[i], nil
		}
	}
	return nil, nil
}

func portIpAssignmentIds(port *swagger.DevicePort) []int32 {
	ids := make([]int32, 0, len(port.Ips))
	for _, ip := range port.Ips {
		ids = append(ids, ip.AssignmentId)
	}
	return ids
}

// applyDevicePort sends the port configuration. Attributes not present in the
// configuration are filled from the port's current state so nothing is reset
// unintentionally.
func (hv *Client) applyDevicePort(ctx context.Context, d *schema.ResourceData, timeout time.Duration) error {
	deviceId := int32(d.Get("device_id").(int))
	portId := int32(d.Get("port_id").(int))

	current, err := hv.findDevicePort(deviceId, portId)
	if err != nil {
		return err
	}
	if current == nil {
		return fmt.Errorf("port %d not found on device %d", portId, deviceId)
	}

	// Build the payload by hand so that an empty ipAssignments list is sent as
	// `[]` (the generated model would drop it) and unset fields are omitted.
	port := map[string]interface{}{
		"portId":        portId,
		"ipAssignments": portIpAssignmentIds(current),
	}
	if configHas(d, "ip_assignment_ids") {
		port["ipAssignments"] = expandInt32List(d.Get("ip_assignment_ids"))
	}
	if configHas(d, "enabled") {
		port["enabled"] = d.Get("enabled").(bool)
	}
	if configHas(d, "native_vlan_id") {
		port["nativeVlanId"] = d.Get("native_vlan_id").(int)
	}

	payload := map[string]interface{}{"ports": []interface{}{port}}

	var tasks []swagger.NetworkTaskDump
	if _, err := hv.rawJSON(ctx, "PUT", fmt.Sprintf("/device/%d/ports", deviceId), payload, &tasks); err != nil {
		return err
	}

	return hv.waitForNetworkTasks(ctx, tasks, timeout)
}

func resourceDevicePortCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Configuring port %d of device %d", d.Get("port_id").(int), d.Get("device_id").(int))
	if err := hv.applyDevicePort(ctx, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(devicePortId(d.Get("device_id").(int), d.Get("port_id").(int)))
	return resourceDevicePortRead(ctx, d, m)
}

func resourceDevicePortRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, portId, err := parseTwoPartId(d.Id(), "<device_id>:<port_id>")
	if err != nil {
		return diag.FromErr(err)
	}

	port, err := hv.findDevicePort(int32(deviceId), int32(portId))
	if err != nil {
		return diag.FromErr(err)
	}
	if port == nil {
		log.Printf("[WARN] Port %d of device %d not found, removing from state", portId, deviceId)
		d.SetId("")
		return nil
	}

	values := map[string]interface{}{
		"device_id":         deviceId,
		"port_id":           portId,
		"native_vlan_id":    int(port.NativeVlanId),
		"native_vlan_tag":   int(port.NativeVlanTag),
		"ip_assignment_ids": SetFromInt32List(portIpAssignmentIds(port)),
		"name":              port.Name,
		"type":              port.Type_,
		"status":            port.Status,
		"private":           port.Private,
		"mtu":               int(port.Mtu),
	}

	// The API does not expose `enabled` directly; derive it from the status
	// when the status is recognisable and otherwise keep the configured value.
	switch strings.ToLower(port.Status) {
	case "enabled", "up", "active", "connected":
		values["enabled"] = true
	case "disabled", "down", "inactive", "disconnected":
		values["enabled"] = false
	}

	return diag.FromErr(setResourceData(d, values))
}

func resourceDevicePortUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	if d.HasChanges("enabled", "native_vlan_id", "ip_assignment_ids") {
		log.Printf("[INFO] Updating port %d of device %d", d.Get("port_id").(int), d.Get("device_id").(int))
		if err := hv.applyDevicePort(ctx, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDevicePortRead(ctx, d, m)
}

func resourceDevicePortDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Intentionally a no-op: clearing the port would also remove its primary IP
	// and disconnect the device. Users who want a clean port can set the
	// attributes explicitly before removing the resource.
	log.Printf("[INFO] Removing port %s from state without changing its configuration", d.Id())
	d.SetId("")
	return nil
}

func resourceDevicePortImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	deviceId, portId, err := parseTwoPartId(d.Id(), "<device_id>:<port_id>")
	if err != nil {
		return nil, err
	}
	if err := d.Set("device_id", deviceId); err != nil {
		return nil, err
	}
	if err := d.Set("port_id", portId); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
