package hivelocity

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

// BareMetalDeviceTimeout is the timeout for creating/updating devices
const BareMetalDeviceTimeout = 60 * time.Minute

func resourceBareMetalDevice(forceNew bool) *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(BareMetalDeviceTimeout),
		},
		CreateContext: resourceBareMetalDeviceCreate,
		ReadContext:   resourceBareMetalDeviceRead,
		UpdateContext: resourceBareMetalDeviceUpdate,
		DeleteContext: resourceBareMetalDeviceDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Description: "Last time this device was updated",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"device_id": {
				Description: "Device ID",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"order_id": {
				Description: "Order ID",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"service_id": {
				Description: "Service ID",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"product_id": {
				Description: "Product ID to pick from the stock",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    forceNew,
			},
			"product_name": {
				Description: "Product Name",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"os_name": {
				Description: "Operating system to install on device",
				Type:        schema.TypeString,
				Required:    true,
			},
			"location_name": {
				Description: "Deploy device in this location",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    forceNew,
			},
			"hostname": {
				Description: "Hostname for this device",
				Type:        schema.TypeString,
				Required:    true,
			},
			"power_status": {
				Description: "Power status",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"vlan_id": {
				Description: "VLAN ID",
				Deprecated:  "This field is deprecated. Please use a hivelocity_vlan resource instead.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
			},
			"primary_ip": {
				Description: "Primary IP of device",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"tags": {
				Description: "Tags to apply for device",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"script": {
				Description: "Post-install script for device",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
			},
			"period": {
				Description: "Billing period for device",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"public_ssh_key_id": {
				Description: "ID of a SSH Key to apply for device",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
			},
			"force_device_id": {
				Description: "Force deployment of this Device ID (internal use only)",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
			},
			"ignition_id": {
				Description: "IgnitionConfig ID",
				Type:        schema.TypeInt,
				Optional:    true,
			},
		},
	}
}

func resourceBareMetalDeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	var tags []string
	for _, tag := range d.Get("tags").([]interface{}) {
		tags = append(tags, tag.(string))
	}

	payload := swagger.BareMetalDeviceCreate{
		ProductId:      int32(d.Get("product_id").(int)),
		Hostname:       d.Get("hostname").(string),
		OsName:         d.Get("os_name").(string),
		LocationName:   d.Get("location_name").(string),
		Script:         d.Get("script").(string),
		Period:         d.Get("period").(string),
		PublicSshKeyId: int32(d.Get("public_ssh_key_id").(int)),
		ForceDeviceId:  int32(d.Get("force_device_id").(int)),
		IgnitionId:     int32(d.Get("ignition_id").(int)),
		Tags:           tags,
	}

	bareMetalDevice, _, err := hv.client.BareMetalDevicesApi.PostBareMetalDeviceResource(hv.auth, payload, nil)
	if err != nil {
		d.SetId("")
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("POST /bare-metal-devices failed! (%s)\n\n %s", err, myErr.Body())
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	_, err = waitForOrder(timeout, hv, bareMetalDevice.OrderId)
	if err != nil {
		d.SetId("")
		myErr, _ := err.(swagger.GenericSwaggerError)
		if strings.Contains(fmt.Sprint(err), "'cancelled'") {
			return diag.Errorf("Your deployment (order %d) has been 'cancelled'. Please contact Hivelocity support if you believe this is a mistake.\n\n %s",
				bareMetalDevice.OrderId, myErr.Body())
		}
		return diag.Errorf("Error provisioning order %d. The Hivelocity team will investigate:\n\n%s\n\n %s",
			bareMetalDevice.OrderId, err, myErr.Body())
	}

	newDevice := []swagger.BareMetalDevice{bareMetalDevice}
	devices, err := waitForDevices(timeout, hv, bareMetalDevice.OrderId, newDevice)
	if err != nil {
		d.SetId("")
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("Error finding devices for order %d. The Hivelocity team will investigate:\n\n%s\n\n %s",
			bareMetalDevice.OrderId, err, myErr.Body())
	}

	newDeviceId := devices.([]swagger.BareMetalDevice)[0].DeviceId
	d.SetId(fmt.Sprint(newDeviceId))
	d.Set("device_id", newDeviceId)

	return resourceBareMetalDeviceRead(ctx, d, m)
}

func resourceBareMetalDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deviceResponse, httpResponse, err := hv.client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /bare-metal-devices/%d failed! (%s)\n\n %s", deviceId, err, myErr.Body())
	}

	d.Set("device_id", deviceResponse.DeviceId)
	d.Set("hostname", deviceResponse.Hostname)
	d.Set("location_name", deviceResponse.LocationName)
	d.Set("order_id", deviceResponse.OrderId)
	d.Set("os_name", deviceResponse.OsName)
	d.Set("power_status", deviceResponse.PowerStatus)
	d.Set("primary_ip", deviceResponse.PrimaryIp)
	d.Set("product_id", deviceResponse.ProductId)
	d.Set("product_name", deviceResponse.ProductName)
	d.Set("service_id", deviceResponse.ServiceId)
	d.Set("tags", deviceResponse.Tags)
	d.Set("public_ssh_key_id", deviceResponse.PublicSshKeyId)
	d.Set("script", deviceResponse.Script)

	return nil
}

func resourceBareMetalDeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload := swagger.BareMetalDeviceUpdate{}
	reload_required := false

	payload.Tags = getTags(d, "")

	ignitionId := d.Get("ignition_id").(int32)
	payload.IgnitionId = ignitionId
	if d.HasChange("ignition_id") {
		reload_required = true
	}

	hostname := d.Get("hostname").(string)
	payload.Hostname = hostname
	if d.HasChange("hostname") {
		reload_required = true
	}

	osName := d.Get("os_name").(string)
	payload.OsName = osName
	if d.HasChange("os_name") {
		reload_required = true
	}

	publicSshKeyId := d.Get("public_ssh_key_id").(int)
	payload.PublicSshKeyId = int32(publicSshKeyId)
	if d.HasChange("public_ssh_key_id") {
		reload_required = true
	}

	script := d.Get("script").(string)
	payload.Script = script
	if d.HasChange("script") {
		reload_required = true
	}

	// If a reload is required, it's necessary to turn the device off first
	if reload_required {
		devicePower, _, err := hv.client.DeviceApi.GetPowerResource(hv.auth, int32(deviceId), nil)
		if err != nil {
			myErr, _ := err.(swagger.GenericSwaggerError)
			return diag.Errorf("GET /device/%s/power failed! (%s)\n\n %s", fmt.Sprint(deviceId), err, myErr.Body())
		}

		if devicePower.PowerStatus == "ON" {
			_, _, err = hv.client.DeviceApi.PostPowerResource(hv.auth, int32(deviceId), "shutdown", nil)
			if err != nil {
				myErr, _ := err.(swagger.GenericSwaggerError)
				return diag.Errorf("POST /device/%s/power failed! (%s)\n\n %s", fmt.Sprint(deviceId), err, myErr.Body())
			}

			// Power status will transition to PENDING, then OFF
			_, err := waitForDevicePowerOff(d, hv, int32(deviceId))
			if err != nil {
				return diag.Errorf("error powering off device %s. The Hivelocity team will investigate:\n\n%s", fmt.Sprint(deviceId), err)
			}
		}
	}

	_, _, err = hv.client.BareMetalDevicesApi.PutBareMetalDeviceIdResource(hv.auth, int32(deviceId), payload, nil)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /bare-metal-devices/%d failed! (%s)\n\n %s", deviceId, err, myErr.Body())
	}

	if reload_required {
		_, err := waitForDeviceReload(d, hv, int32(deviceId))
		if err != nil {
			return diag.Errorf("error reloading device %s. The Hivelocity team will investigate:\n\n%s", fmt.Sprint(deviceId), err)
		}
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
	return resourceBareMetalDeviceRead(ctx, d, m)
}

func resourceBareMetalDeviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := hv.client.BareMetalDevicesApi.DeleteBareMetalDeviceIdResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /bare-metal-devices/%d failed! (%s)\n\n %s", deviceId, err, myErr.Body())
	}

	d.SetId("")

	return diags
}
