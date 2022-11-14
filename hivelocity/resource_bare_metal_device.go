package hivelocity

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
			"bonded": {
				Description: "When set, prefer only bonded devices",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			// "private_network": {
			// 	Description: "Private network and IP of device in CIDR notation",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
		},
	}
}

func resourceBareMetalDeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	var tags []string
	for _, tag := range d.Get("tags").([]interface{}) {
		tags = append(tags, tag.(string))
	}

	// Load existing deviceId if any
	var deviceId, orderId int32
	if d.Id() != "" && d.Get("order_id") != "" {
		deviceId_, deviceErr := strconv.Atoi(d.Id())
		orderId_, orderErr := strconv.Atoi(d.Get("order_id").(string))

		if deviceErr == nil && orderErr == nil {
			deviceId = int32(deviceId_)
			orderId = int32(orderId_)
			tflog.Warn(ctx, "Device was tainted and may have failed after deployment")
		} else {
			tflog.Error(ctx, "Invalid device ID or order ID from tainted device")
			d.SetId("")
		}
	} else {
		tflog.Info(ctx, "Creating new device device resource")
	}

	if deviceId == 0 {
		// Create if no existing ID
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
			BondingSupport: d.Get("bonded").(bool),
			// PrivateNetwork: d.Get("private_network").(string),
			Tags: tags,
		}

		bareMetalDevice, _, err := hv.client.BareMetalDevicesApi.PostBareMetalDeviceResource(hv.auth, payload, nil)
		if err != nil {
			d.SetId("")
			err = formatSwaggerError(err, "POST /bare-metal-devices")
			return diag.FromErr(err)
		}
		deviceId = bareMetalDevice.DeviceId
		orderId = bareMetalDevice.OrderId

		// Set device ID and order ID after creation so any errors before the function completes will mark it as
		// "tainted" so we can wait on the device again. Though this does not work at the moment because of device
		// statuses being inaccessible for new devices
		d.SetId(fmt.Sprintf("%d", deviceId))
		d.Set("order_id", orderId)
	}

	// Wait for order if not complete
	timeout := d.Timeout(schema.TimeoutCreate)
	if status, err := waitForOrder(timeout, hv, orderId); err != nil {
		if status == "cancelled" {
			return diag.Errorf("Your deployment (order %d) has been 'cancelled'. Please contact Hivelocity"+
				"support if you believe this is a mistake.\n\n%s",
				orderId, err)
		}
		return diag.FromErr(err)
	}

	// Wait for device, if not provisioned
	if _, err := hv.waitForDevice(deviceId, timeout); err != nil {
		return diag.FromErr(err)
	}

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
	// d.Set("private_network", deviceResponse.PrivateNetwork)

	return nil
}

func resourceBareMetalDeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	var deviceId int32
	if deviceId_, err := strconv.Atoi(d.Id()); err == nil {
		deviceId = int32(deviceId_)
	} else {
		return diag.FromErr(err)
	}

	updatePayload := getBareMetalUpdatePayload(d)
	if err := hv._updateDevice(deviceId, updatePayload, d); err != nil {
		return diag.FromErr(err)
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
