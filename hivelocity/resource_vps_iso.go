package hivelocity

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceVpsIso() *schema.Resource {
	return &schema.Resource{
		Description: "Imports an ISO image from a URL into a facility so it can be attached to VPS " +
			"instances (`hivelocity_vps.iso_id`). The provider waits until the download finishes. ISOs " +
			"cannot be modified; any change replaces the image.",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		CreateContext: resourceVpsIsoCreate,
		ReadContext:   resourceVpsIsoRead,
		DeleteContext: resourceVpsIsoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the ISO.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "Description of the ISO.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"url": {
				Description: "URL the image is downloaded from.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"facility_code": {
				Description: "Facility to store the image in, for example `TPA1`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"bootable": {
				Description: "Whether the image is bootable.",
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
			},
			"iso_id":     computedString("Unique ID of the ISO."),
			"ready":      computedBool("True once the image is downloaded and usable."),
			"status":     computedString("Status of the image."),
			"size":       computedInt("Size of the image."),
			"device_ids": computedIntList("VPS instances the ISO is attached to."),
		},
	}
}

func (hv *Client) waitForVpsIso(ctx context.Context, isoId string, timeout time.Duration) error {
	conf := &retry.StateChangeConf{
		Pending: []string{"waiting"},
		Target:  []string{"ready"},
		Refresh: func() (interface{}, string, error) {
			iso, _, err := hv.client.VPSApi.GetVpsisoIdResource(hv.auth, isoId, nil)
			if err != nil {
				return nil, "", formatSwaggerError(err, "GET /vps/iso/%s", isoId)
			}
			status := strings.ToLower(iso.Status)
			if strings.Contains(status, "fail") || strings.Contains(status, "error") {
				return &iso, status, fmt.Errorf("ISO %s entered status %q", isoId, iso.Status)
			}
			if iso.Ready {
				return &iso, "ready", nil
			}
			return &iso, "waiting", nil
		},
		Timeout:                   timeout,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 1,
	}
	_, err := conf.WaitForStateContext(ctx)
	return err
}

func resourceVpsIsoCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := swagger.VpsisoCreate{
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Url:          d.Get("url").(string),
		FacilityCode: d.Get("facility_code").(string),
		Bootable:     d.Get("bootable").(bool),
	}
	log.Printf("[INFO] Importing ISO %q into %s", payload.Name, payload.FacilityCode)
	iso, _, err := hv.client.VPSApi.PostVpsisoResource(hv.auth, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /vps/iso")
	}
	if iso.IsoId == "" {
		return diag.Errorf("POST /vps/iso did not return an isoId")
	}

	d.SetId(iso.IsoId)
	if err := hv.waitForVpsIso(ctx, iso.IsoId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	return resourceVpsIsoRead(ctx, d, m)
}

func resourceVpsIsoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	iso, resp, err := hv.client.VPSApi.GetVpsisoIdResource(hv.auth, d.Id(), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] VPS ISO %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /vps/iso/%s", d.Id())
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"name":          iso.Name,
		"description":   iso.Description,
		"url":           iso.Url,
		"facility_code": iso.FacilityCode,
		"bootable":      iso.Bootable,
		"iso_id":        iso.IsoId,
		"ready":         iso.Ready,
		"status":        iso.Status,
		"size":          int(iso.Size),
		"device_ids":    flattenInt32List(iso.DeviceIds),
	}))
}

func resourceVpsIsoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Deleting VPS ISO %s", d.Id())
	task, resp, err := hv.client.VPSApi.DeleteVpsisoIdResource(hv.auth, d.Id(), nil)
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /vps/iso/%s", d.Id())
	}
	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func dataSourceVpsIso() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up a VPS ISO image, for example by `name`.",
		ReadContext: dataSourceVpsIsoRead,
		Schema: map[string]*schema.Schema{
			"filter":        dataSourceFiltersSchema(),
			"first":         dataSourceFilterFirstSchema(),
			"iso_id":        computedString("Unique ID of the ISO."),
			"name":          computedString("Name."),
			"description":   computedString("Description."),
			"url":           computedString("Source URL."),
			"facility_code": computedString("Facility."),
			"bootable":      computedBool("Whether the image is bootable."),
			"ready":         computedBool("True once downloaded."),
			"status":        computedString("Status."),
			"size":          computedInt("Size."),
			"device_ids":    computedIntList("VPS instances the ISO is attached to."),
		},
	}
}

func dataSourceVpsIsoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	isos, _, err := hv.client.VPSApi.GetVpsisoResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /vps/iso")
	}

	return readFilteredList(d, isos, dataSourceVpsIso().Schema, "iso_id", nil)
}
