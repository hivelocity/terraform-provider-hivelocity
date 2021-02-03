package hivelocity

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
	"strconv"
	"strings"
	"time"
)

func resourceBareMetalDevice() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},
		CreateContext: resourceBareMetalDeviceCreate,
		ReadContext:   resourceBareMetalDeviceRead,
		UpdateContext: resourceBareMetalDeviceUpdate,
		DeleteContext: resourceBareMetalDeviceDelete,
		Schema: map[string]*schema.Schema{
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
				Required: true,
				ForceNew: true,
			},
			"product_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"os_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
			"script": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

// TODO: Test what happens when you change hostname, tags, etc anything that is required.

func resourceBareMetalDeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	tags, err := getTags(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	payload := swagger.BareMetalDeviceCreate{
		ProductId:    int32(d.Get("product_id").(int)),
		Hostname:     d.Get("hostname").(string),
		OsName:       d.Get("os_name").(string),
		VlanId:       int32(d.Get("vlan_id").(int)),
		LocationName: d.Get("location_name").(string),
		Script:       d.Get("script").(string),
	}

	bareMetalDevice, _, err := hv.client.BareMetalDevicesApi.PostBareMetalDeviceResource(hv.auth, payload, nil)
	if err != nil {
		d.SetId("")
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("POST /bare-metal-device failed! (%s)\n\n %s", err, myErr.Body())
	}

	_, err = waitForOrder(d, hv, bareMetalDevice.OrderId)
	if err != nil {
		d.SetId("")
		if strings.Contains(fmt.Sprint(err), "'cancelled'") {
			return diag.Errorf("Your deployment (order %s) has been 'cancelled'. Please contact Hivelocity support if you believe this is a mistake.", fmt.Sprint(bareMetalDevice.OrderId))
		}
		return diag.Errorf("error provisioning order %s. The Hivelocity team will investigate:\n\n%s", fmt.Sprint(bareMetalDevice.OrderId), err)
	}

	device, err := waitForDevice(d, hv, bareMetalDevice.OrderId)
	if err != nil {
		d.SetId("")
		return diag.Errorf("error finding devices for order %s. The Hivelocity team will investigate:\n\n%s", fmt.Sprint(bareMetalDevice.OrderId), err)
	}

	newDeviceId := device.(swagger.BareMetalDevice).DeviceId
	_, err = updateTagsForCreate(hv, newDeviceId, tags)
	if err != nil {
		// TODO: The deployment was successful, so we should throw a warning here that tags failed for some reason.
	}
	d.SetId(fmt.Sprint(newDeviceId))

	return resourceBareMetalDeviceRead(ctx, d, m)
}

func resourceBareMetalDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deviceResponse, _, err := hv.client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		return diag.FromErr(err)
	}

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
	d.Set("vlanId", deviceResponse.VlanId)

	return diags
}

func resourceBareMetalDeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload := swagger.BareMetalDeviceUpdate{}

	if d.HasChange("hostname") {
		hostname := d.Get("hostname").(string)
		payload.Hostname = hostname
	}

	if d.HasChange("tags") {
		tags, err := getTags(d)
		if err != nil {
			return diag.FromErr(err)
		}
		payload.Tags = tags
	}
	if d.HasChange("vlan_id") {
		// TODO: Currently no-op until VLAN IDS deployed
	}

	_, _, err = hv.client.BareMetalDevicesApi.PutBareMetalDeviceIdResource(hv.auth, int32(deviceId), payload, nil)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /bare-metal-device/%s failed! (%s)\n\n %s", fmt.Sprint(deviceId), err, myErr.Body())
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

	// Check device exists still, if not mark as already destroyed.
	_, _, err = hv.client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		d.SetId("")
		return diags
	}

	_, err = hv.client.BareMetalDevicesApi.DeleteBareMetalDeviceIdResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /bare-metal-device/%s failed! (%s)\n\n %s", fmt.Sprint(deviceId), err, myErr.Body())
	}

	d.SetId("")

	return diags
}
