package hivelocity

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceDevicePortBond() *schema.Resource {
	return &schema.Resource{
		Description: "Bonds the network ports of a device into a single LACP bond. The device must have " +
			"been provisioned with bonding support (`bonded = true` on `hivelocity_bare_metal_device`). " +
			"Destroying the resource unbonds the ports. The resource ID is the device ID.",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceDevicePortBondCreate,
		ReadContext:   resourceDevicePortBondRead,
		DeleteContext: resourceDevicePortBondDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"device_id": {
				Description: "ID of the device whose ports are bonded.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"bond_port_id": {
				Description: "Port ID of the resulting bond interface. Use it in `hivelocity_vlan.port_ids`.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"bond_name": {
				Description: "Interface name of the bond, for example `bond0`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// findBondPort returns the bond port of a device, or nil when the device is not bonded.
func (hv *Client) findBondPort(deviceId int32) (*swagger.DevicePort, error) {
	ports, _, err := hv.client.DeviceApi.GetDevicePortResource(hv.auth, deviceId, nil)
	if err != nil {
		return nil, formatSwaggerError(err, "GET /device/%d/ports", deviceId)
	}
	for i := range ports {
		if strings.HasPrefix(strings.ToLower(ports[i].Name), "bond") || strings.EqualFold(ports[i].Type_, "bond") {
			return &ports[i], nil
		}
	}
	return nil, nil
}

func resourceDevicePortBondCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	deviceId := int32(d.Get("device_id").(int))

	log.Printf("[INFO] Bonding ports of device %d", deviceId)
	task, _, err := hv.client.DeviceApi.PostDeviceBondResource(hv.auth, deviceId, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /device/%d/ports/bond", deviceId)
	}

	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(deviceId))
	return resourceDevicePortBondRead(ctx, d, m)
}

func resourceDevicePortBondRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	bond, err := hv.findBondPort(int32(deviceId))
	if err != nil {
		return diag.FromErr(err)
	}
	if bond == nil {
		log.Printf("[WARN] Device %d has no bond, removing from state", deviceId)
		d.SetId("")
		return nil
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"device_id":    deviceId,
		"bond_port_id": int(bond.PortId),
		"bond_name":    bond.Name,
	}))
}

func resourceDevicePortBondDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Unbonding ports of device %d", deviceId)
	task, resp, err := hv.client.DeviceApi.DeleteDeviceBondResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /device/%d/ports/bond", deviceId)
	}

	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
