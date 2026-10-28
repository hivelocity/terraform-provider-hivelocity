package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

// VpsTimeout is the default timeout for provisioning and reloading VPS instances.
const VpsTimeout = 60 * time.Minute

func vpsNicSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"nic_id":       computedString("Unique ID of the NIC."),
			"nic_name":     computedString("Name of the NIC."),
			"network_id":   computedString("ID of the VPS network the NIC is attached to."),
			"network_name": computedString("Name of the VPS network."),
			"mac_address":  computedString("MAC address."),
			"ip_address":   computedString("IPv4 address."),
			"netmask":      computedString("IPv4 netmask."),
			"gateway":      computedString("IPv4 gateway."),
			"ipv6_address": computedString("IPv6 address."),
			"ipv6_cidr":    computedString("IPv6 CIDR."),
			"ipv6_gateway": computedString("IPv6 gateway."),
			"vlan_id":      computedInt("VLAN ID, for VLAN-enabled networks."),
			"vlan_tag":     computedInt("VLAN tag, for VLAN-enabled networks."),
		},
	}
}

func resourceVps() *schema.Resource {
	return &schema.Resource{
		Description: "Provisions a VPS instance. Instances are ordered through the compute endpoint with a " +
			"VPS product ID (see `hivelocity_product` with `is_vps = true`) and then managed through the " +
			"VPS endpoints: `product_id` changes scale the instance, `ssh_key_ids` and `iso_id` changes " +
			"are applied in place, and `os_name`/`hostname`/`script`/`tags` changes go through a reload.",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(VpsTimeout),
			Update: schema.DefaultTimeout(VpsTimeout),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceVpsCreate,
		ReadContext:   resourceVpsRead,
		UpdateContext: resourceVpsUpdate,
		DeleteContext: resourceVpsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"product_id": {
				Description: "ID of the VPS product. Changing it scales the instance to the new product.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"os_name": {
				Description: "Operating system to install (see `hivelocity_product_operating_systems`).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"location_name": {
				Description: "Facility code to deploy in, for example `TPA1`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"hostname": {
				Description: "Hostname of the instance.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"script": {
				Description: "Cloud-init or post-install script.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tags": {
				Description: "Tags to apply to the instance.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"period": {
				Description:  "Billing period: hourly|monthly|quarterly|semi-annually|annually|biennial|triennial.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(billingPeriods, false),
			},
			"public_ssh_key_id": {
				Description: "ID of the SSH key to inject at provisioning time.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"ssh_key_ids": {
				Description: "IDs of the SSH keys authorized on the instance. Applied in place after provisioning.",
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"iso_id": {
				Description: "ID of a `hivelocity_vps_iso` to attach to the instance. Unset to detach.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"force_device_id": {
				Description: "Force deployment of this device ID (internal use only).",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"device_id":     computedInt("Device ID of the instance."),
			"order_id":      computedInt("Order the instance was created from."),
			"service_id":    computedInt("Service ID of the instance."),
			"product_name":  computedString("Product name."),
			"primary_ip":    computedString("Primary public IP."),
			"power_status":  computedString("Power status (ON|OFF)."),
			"status":        computedString("Instance status, for example `Running`."),
			"vcpu":          computedInt("Number of vCPUs."),
			"memory":        computedInt("Memory in MB."),
			"facility_code": computedString("Facility of the instance."),
			"template_id":   computedString("Template used for the install."),
			"cloud_init":    computedString("Cloud-init script in effect."),
			"password": {
				Description: "Root password as reported by the API.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"nics": {
				Description: "Network interfaces of the instance.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        vpsNicSchema(),
			},
		},
	}
}

// waitForVps polls the VPS endpoint until the instance reports a running state.
func (hv *Client) waitForVps(ctx context.Context, deviceId int32, timeout time.Duration) error {
	conf := &retry.StateChangeConf{
		Pending: []string{"waiting"},
		Target:  []string{"ready"},
		Refresh: func() (interface{}, string, error) {
			vps, resp, err := hv.client.VPSApi.GetVpsIdResource(hv.auth, deviceId, nil)
			if err != nil {
				if isNotFound(resp) {
					// Not visible yet while the order is being processed.
					return nil, "", nil
				}
				return nil, "", formatSwaggerError(err, "GET /vps/%d", deviceId)
			}
			status := strings.ToLower(vps.Status)
			if strings.Contains(status, "fail") || strings.Contains(status, "error") {
				return &vps, status, fmt.Errorf("VPS %d entered status %q", deviceId, vps.Status)
			}
			if status == "running" {
				return &vps, "ready", nil
			}
			return &vps, "waiting", nil
		},
		Timeout:                   timeout,
		Delay:                     30 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 1,
		NotFoundChecks:            360,
	}
	_, err := conf.WaitForStateContext(ctx)
	return err
}

func (hv *Client) updateVpsSettings(ctx context.Context, deviceId int32, payload swagger.VpsUpdate, timeout time.Duration) error {
	task, _, err := hv.client.VPSApi.PutVpsIdResource(hv.auth, deviceId, payload, nil)
	if err != nil {
		return formatSwaggerError(err, "PUT /vps/%d", deviceId)
	}
	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, timeout); err != nil {
		return fmt.Errorf("waiting for VPS %d update: %w", deviceId, err)
	}
	return nil
}

func (hv *Client) setVpsIso(ctx context.Context, deviceId int32, oldIso, newIso string, timeout time.Duration) error {
	if oldIso != "" {
		task, _, err := hv.client.VPSApi.PostVpsDetachIsoResource(hv.auth, deviceId, nil)
		if err != nil {
			return formatSwaggerError(err, "POST /vps/%d/detachiso", deviceId)
		}
		if _, err := hv.waitForNetworkTask(ctx, task.TaskId, timeout); err != nil {
			return err
		}
	}
	if newIso != "" {
		task, _, err := hv.client.VPSApi.PostVpsAttachIsoResource(hv.auth, deviceId, swagger.VpsisoAttach{IsoId: newIso}, nil)
		if err != nil {
			return formatSwaggerError(err, "POST /vps/%d/attachiso", deviceId)
		}
		if _, err := hv.waitForNetworkTask(ctx, task.TaskId, timeout); err != nil {
			return err
		}
	}
	return nil
}

func resourceVpsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := swagger.BareMetalDeviceCreate{
		ProductId:      int32(d.Get("product_id").(int)),
		Hostname:       d.Get("hostname").(string),
		OsName:         d.Get("os_name").(string),
		LocationName:   d.Get("location_name").(string),
		Script:         d.Get("script").(string),
		Period:         d.Get("period").(string),
		PublicSshKeyId: int32(d.Get("public_ssh_key_id").(int)),
		ForceDeviceId:  int32(d.Get("force_device_id").(int)),
		Tags:           expandStringList(d.Get("tags")),
	}

	log.Printf("[INFO] Provisioning VPS (product %d in %s)", payload.ProductId, payload.LocationName)
	device, _, err := hv.client.ComputeDevicesApi.PostComputeResource(hv.auth, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /compute")
	}

	d.SetId(fmt.Sprint(device.DeviceId))
	if err := d.Set("order_id", int(device.OrderId)); err != nil {
		return diag.FromErr(err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	if status, err := waitForOrder(ctx, timeout, hv, device.OrderId); err != nil {
		if status == "cancelled" {
			return diag.Errorf("Your deployment (order %d) has been 'cancelled'. Please contact Hivelocity "+
				"support if you believe this is a mistake.\n\n%s", device.OrderId, err)
		}
		return diag.FromErr(err)
	}

	if err := hv.waitForVps(ctx, device.DeviceId, timeout); err != nil {
		return diag.FromErr(err)
	}

	if configHas(d, "ssh_key_ids") {
		update := swagger.VpsUpdate{SshKeyIds: expandInt32List(d.Get("ssh_key_ids"))}
		if err := hv.updateVpsSettings(ctx, device.DeviceId, update, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	if iso := d.Get("iso_id").(string); iso != "" {
		if err := hv.setVpsIso(ctx, device.DeviceId, "", iso, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVpsRead(ctx, d, m)
}

func flattenVpsNics(raw []interface{}) ([]interface{}, error) {
	data, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	var nics []swagger.VpsNetworkInterface
	if err := json.Unmarshal(data, &nics); err != nil {
		return nil, fmt.Errorf("decoding VPS nics: %w", err)
	}

	list := make([]interface{}, 0, len(nics))
	for _, n := range nics {
		list = append(list, map[string]interface{}{
			"nic_id":       n.NicId,
			"nic_name":     n.NicName,
			"network_id":   n.NetworkId,
			"network_name": n.NetworkName,
			"mac_address":  n.MacAddress,
			"ip_address":   n.IpAddress,
			"netmask":      n.Netmask,
			"gateway":      n.Gateway,
			"ipv6_address": n.Ipv6Address,
			"ipv6_cidr":    n.Ipv6Cidr,
			"ipv6_gateway": n.Ipv6Gateway,
			"vlan_id":      int(n.VlanId),
			"vlan_tag":     int(n.VlanTag),
		})
	}
	return list, nil
}

func resourceVpsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	device, resp, err := hv.client.ComputeDevicesApi.GetComputeIdResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] VPS %d not found, removing from state", deviceId)
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /compute/%d", deviceId)
	}

	values := map[string]interface{}{
		"device_id":         int(device.DeviceId),
		"order_id":          int(device.OrderId),
		"service_id":        int(device.ServiceId),
		"product_id":        int(device.ProductId),
		"product_name":      device.ProductName,
		"hostname":          device.Hostname,
		"location_name":     device.LocationName,
		"os_name":           device.OsName,
		"script":            device.Script,
		"tags":              flattenStringList(device.Tags),
		"period":            device.Period,
		"public_ssh_key_id": int(device.PublicSshKeyId),
		"primary_ip":        device.PrimaryIp,
		"power_status":      device.PowerStatus,
	}

	vps, resp, err := hv.client.VPSApi.GetVpsIdResource(hv.auth, int32(deviceId), nil)
	if err != nil && !isNotFound(resp) {
		return apiErrorDiag(err, "GET /vps/%d", deviceId)
	}
	if err == nil {
		nics, err := flattenVpsNics(vps.Nics)
		if err != nil {
			return diag.FromErr(err)
		}
		values["status"] = vps.Status
		values["vcpu"] = int(vps.Vcpu)
		values["memory"] = int(vps.Memory)
		values["facility_code"] = vps.FacilityCode
		values["template_id"] = vps.TemplateId
		values["cloud_init"] = vps.CloudInit
		values["password"] = vps.Password
		values["iso_id"] = vps.IsoId
		values["ssh_key_ids"] = SetFromInt32List(vps.SshKeyIds)
		values["nics"] = nics
	}

	return diag.FromErr(setResourceData(d, values))
}

func resourceVpsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	deviceId := int32(d.Get("device_id").(int))
	timeout := d.Timeout(schema.TimeoutUpdate)

	// Reload-style changes go through the compute endpoint.
	if d.HasChanges("hostname", "os_name", "script", "tags", "public_ssh_key_id") {
		payload := swagger.BareMetalDeviceUpdate{
			Hostname:       d.Get("hostname").(string),
			OsName:         d.Get("os_name").(string),
			Script:         d.Get("script").(string),
			Tags:           expandStringList(d.Get("tags")),
			PublicSshKeyId: int32(d.Get("public_ssh_key_id").(int)),
		}
		log.Printf("[INFO] Updating VPS %d through /compute", deviceId)
		_, _, err := hv.client.ComputeDevicesApi.PutComputeIdResource(hv.auth, deviceId, payload, nil)
		if err != nil && isPowerOffError(formatSwaggerError(err, "PUT /compute/%d", deviceId)) {
			opts := &swagger.ComputeDevicesApiPutComputeIdResourceOpts{SkipPowerCheck: optional.NewBool(true)}
			_, _, err = hv.client.ComputeDevicesApi.PutComputeIdResource(hv.auth, deviceId, payload, opts)
		}
		if err != nil {
			return apiErrorDiag(err, "PUT /compute/%d", deviceId)
		}
		if d.HasChange("os_name") {
			if err := hv.waitForVps(ctx, deviceId, timeout); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChanges("product_id", "ssh_key_ids") {
		payload := swagger.VpsUpdate{}
		if d.HasChange("product_id") {
			payload.ProductId = int32(d.Get("product_id").(int))
		}
		if d.HasChange("ssh_key_ids") {
			payload.SshKeyIds = expandInt32List(d.Get("ssh_key_ids"))
		}
		log.Printf("[INFO] Updating VPS %d settings", deviceId)
		if err := hv.updateVpsSettings(ctx, deviceId, payload, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("iso_id") {
		oldIso, newIso := d.GetChange("iso_id")
		if err := hv.setVpsIso(ctx, deviceId, oldIso.(string), newIso.(string), timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVpsRead(ctx, d, m)
}

func resourceVpsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting VPS %d", deviceId)
	if resp, err := hv.client.ComputeDevicesApi.DeleteComputeIdResource(hv.auth, int32(deviceId)); err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /compute/%d", deviceId)
	}

	d.SetId("")
	return nil
}
