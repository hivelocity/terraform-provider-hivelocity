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

func resourceVpsNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Creates a private VPS network in a facility. Attach instances with " +
			"`hivelocity_vps_network_attachment`. Set `vlan_id` to bridge the network to a VLAN shared " +
			"with bare metal devices.",
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceVpsNetworkCreate,
		ReadContext:   resourceVpsNetworkRead,
		UpdateContext: resourceVpsNetworkUpdate,
		DeleteContext: resourceVpsNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"facility_code": {
				Description: "Facility of the network, for example `TPA1`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Name of the network.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"vlan_id": {
				Description: "ID of a `hivelocity_vlan` to bridge the network to.",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"network_id": computedString("Unique ID of the network."),
			"vlan_tag":   computedInt("VLAN tag of the bridged VLAN, if any."),
		},
	}
}

func resourceVpsNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := swagger.VpsNetworkCreate{
		FacilityCode: d.Get("facility_code").(string),
		Name:         d.Get("name").(string),
		VlanId:       int32(d.Get("vlan_id").(int)),
	}
	network, _, err := hv.client.VPSApi.PostVpsNetworkResource(hv.auth, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /vps/network")
	}
	if network.NetworkId == "" {
		return diag.Errorf("POST /vps/network did not return a networkId")
	}

	d.SetId(network.NetworkId)
	log.Printf("[INFO] Created VPS network %s", d.Id())

	return resourceVpsNetworkRead(ctx, d, m)
}

func resourceVpsNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	network, resp, err := hv.client.VPSApi.GetVpsNetworkIdResource(hv.auth, d.Id(), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] VPS network %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /vps/network/%s", d.Id())
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"facility_code": network.FacilityCode,
		"name":          network.Name,
		"vlan_id":       int(network.VlanId),
		"network_id":    network.NetworkId,
		"vlan_tag":      int(network.VlanTag),
	}))
}

func resourceVpsNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	if d.HasChange("name") {
		payload := swagger.VpsNetworkUpdate{Name: d.Get("name").(string)}
		task, _, err := hv.client.VPSApi.PutVpsNetworkIdResource(hv.auth, d.Id(), payload, nil)
		if err != nil {
			return apiErrorDiag(err, "PUT /vps/network/%s", d.Id())
		}
		if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVpsNetworkRead(ctx, d, m)
}

func resourceVpsNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Deleting VPS network %s", d.Id())
	task, resp, err := hv.client.VPSApi.DeleteVpsNetworkIdResource(hv.auth, d.Id(), nil)
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /vps/network/%s", d.Id())
	}
	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func dataSourceVpsNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up a VPS network, for example by `name` or `facility_code`.",
		ReadContext: dataSourceVpsNetworkRead,
		Schema: map[string]*schema.Schema{
			"filter":        dataSourceFiltersSchema(),
			"first":         dataSourceFilterFirstSchema(),
			"network_id":    computedString("Unique ID of the network."),
			"name":          computedString("Name."),
			"facility_code": computedString("Facility."),
			"vlan_id":       computedInt("Bridged VLAN ID, if any."),
			"vlan_tag":      computedInt("Bridged VLAN tag, if any."),
		},
	}
}

func dataSourceVpsNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	networks, _, err := hv.client.VPSApi.GetVpsNetworkResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /vps/network")
	}

	return readFilteredList(d, networks, dataSourceVpsNetwork().Schema, "network_id", nil)
}

// --- network attachment ----------------------------------------------------

func resourceVpsNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Attaches a VPS instance to a `hivelocity_vps_network`. The resource ID is " +
			"`<device_id>:<network_id>`.",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceVpsNetworkAttachmentCreate,
		ReadContext:   resourceVpsNetworkAttachmentRead,
		DeleteContext: resourceVpsNetworkAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVpsNetworkAttachmentImport,
		},
		Schema: map[string]*schema.Schema{
			"device_id": {
				Description: "ID of the VPS instance.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"network_id": {
				Description: "ID of the VPS network.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"nic_id":     computedString("ID of the NIC created for the attachment."),
			"ip_address": computedString("IPv4 address of the NIC, if any."),
		},
	}
}

func parseVpsNetworkAttachmentId(id string) (int, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return 0, "", fmt.Errorf("unexpected ID format %q, expected <device_id>:<network_id>", id)
	}
	deviceId, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("unexpected ID format %q, expected <device_id>:<network_id>: %w", id, err)
	}
	return deviceId, parts[1], nil
}

func resourceVpsNetworkAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	deviceId := int32(d.Get("device_id").(int))
	networkId := d.Get("network_id").(string)

	log.Printf("[INFO] Attaching VPS %d to network %s", deviceId, networkId)
	payload := swagger.VpsInstanceNetworkAdd{NetworkId: networkId}
	task, _, err := hv.client.VPSApi.PostVpsIdNetworkAddResource(hv.auth, deviceId, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /vps/%d/network", deviceId)
	}
	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", deviceId, networkId))
	return resourceVpsNetworkAttachmentRead(ctx, d, m)
}

func resourceVpsNetworkAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, networkId, err := parseVpsNetworkAttachmentId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vps, resp, err := hv.client.VPSApi.GetVpsIdResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /vps/%d", deviceId)
	}

	nics, err := flattenVpsNics(vps.Nics)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, n := range nics {
		nic := n.(map[string]interface{})
		if nic["network_id"] == networkId {
			return diag.FromErr(setResourceData(d, map[string]interface{}{
				"device_id":  deviceId,
				"network_id": networkId,
				"nic_id":     nic["nic_id"],
				"ip_address": nic["ip_address"],
			}))
		}
	}

	log.Printf("[WARN] VPS %d is no longer attached to network %s, removing from state", deviceId, networkId)
	d.SetId("")
	return nil
}

func resourceVpsNetworkAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, networkId, err := parseVpsNetworkAttachmentId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Detaching VPS %d from network %s", deviceId, networkId)
	task, resp, err := hv.client.VPSApi.DeleteVpsIdNetworkRemoveResource(hv.auth, int32(deviceId), networkId, nil)
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /vps/%d/network/%s", deviceId, networkId)
	}
	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceVpsNetworkAttachmentImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	deviceId, networkId, err := parseVpsNetworkAttachmentId(d.Id())
	if err != nil {
		return nil, err
	}
	if err := d.Set("device_id", deviceId); err != nil {
		return nil, err
	}
	if err := d.Set("network_id", networkId); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
