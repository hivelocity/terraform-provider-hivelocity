package hivelocity

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceDeviceTags() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the full list of tags on a device. Useful for tagging devices that were not " +
			"provisioned by Terraform; devices created with `hivelocity_bare_metal_device` can use its " +
			"`tags` argument instead. Destroying the resource removes all tags. The resource ID is the " +
			"device ID.",
		CreateContext: resourceDeviceTagsCreate,
		ReadContext:   resourceDeviceTagsRead,
		UpdateContext: resourceDeviceTagsUpdate,
		DeleteContext: resourceDeviceTagsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"device_id": {
				Description: "ID of the device.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"tags": {
				Description: "Tags of the device.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func (hv *Client) setDeviceTags(ctx context.Context, deviceId int, tags []string) error {
	// Sent by hand so an empty list is serialized as [] instead of being omitted.
	payload := map[string]interface{}{"tags": tags}
	if _, err := hv.rawJSON(ctx, "PUT", fmt.Sprintf("/device/%d/tags", deviceId), payload, nil); err != nil {
		return err
	}
	return nil
}

func resourceDeviceTagsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	deviceId := d.Get("device_id").(int)

	if err := hv.setDeviceTags(ctx, deviceId, expandStringList(d.Get("tags"))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(deviceId))
	return resourceDeviceTagsRead(ctx, d, m)
}

func resourceDeviceTagsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tags, resp, err := hv.client.DeviceApi.GetDeviceTagIdResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] Device %d not found, removing tags from state", deviceId)
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /device/%d/tags", deviceId)
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"device_id": deviceId,
		"tags":      flattenStringList(tags.Tags),
	}))
}

func resourceDeviceTagsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	if d.HasChange("tags") {
		if err := hv.setDeviceTags(ctx, d.Get("device_id").(int), expandStringList(d.Get("tags"))); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDeviceTagsRead(ctx, d, m)
}

func resourceDeviceTagsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := hv.setDeviceTags(ctx, deviceId, []string{}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func dataSourceDeviceTags() *schema.Resource {
	return &schema.Resource{
		Description: "Lists every tag in use on the account, or the tags of one device when `device_id` is set.",
		ReadContext: dataSourceDeviceTagsRead,
		Schema: map[string]*schema.Schema{
			"device_id": {
				Description: "Limit the result to the tags of this device.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"tags": computedStringList("Tags."),
		},
	}
}

func dataSourceDeviceTagsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	var tags swagger.DeviceTag
	var err error
	id := "all"

	if v, ok := d.GetOk("device_id"); ok {
		id = fmt.Sprint(v.(int))
		tags, _, err = hv.client.DeviceApi.GetDeviceTagIdResource(hv.auth, int32(v.(int)), nil)
		if err != nil {
			return apiErrorDiag(err, "GET /device/%d/tags", v.(int))
		}
	} else {
		tags, _, err = hv.client.DeviceApi.GetClientDeviceTagResource(hv.auth, nil)
		if err != nil {
			return apiErrorDiag(err, "GET /device/tags")
		}
	}

	if err := d.Set("tags", flattenStringList(tags.Tags)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}
