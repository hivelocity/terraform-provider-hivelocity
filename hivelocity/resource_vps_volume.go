package hivelocity

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

const bytesPerGiB = 1024 * 1024 * 1024

func resourceVpsVolume() *schema.Resource {
	return &schema.Resource{
		Description: "Creates an additional block storage volume attached to a VPS instance. Volumes are " +
			"billed monthly. Changing `size` grows the volume in place; changing `device_id` moves it " +
			"to another instance.",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		CreateContext: resourceVpsVolumeCreate,
		ReadContext:   resourceVpsVolumeRead,
		UpdateContext: resourceVpsVolumeUpdate,
		DeleteContext: resourceVpsVolumeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"device_id": {
				Description: "ID of the VPS instance the volume is attached to.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"size": {
				Description: "Size of the volume in GiB.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"volume_id":     computedString("Unique ID of the volume."),
			"size_bytes":    computedInt("Size of the volume in bytes as reported by the API."),
			"type":          computedString("`ROOT` or `DATADISK`."),
			"facility_code": computedString("Facility of the volume."),
			"created_at":    computedString("When the volume was created."),
		},
	}
}

func (hv *Client) listVpsVolumeIds() (map[string]bool, error) {
	volumes, _, err := hv.client.VPSApi.GetVpsVolumeResource(hv.auth, nil)
	if err != nil {
		return nil, formatSwaggerError(err, "GET /vps/volume")
	}
	ids := make(map[string]bool, len(volumes))
	for _, v := range volumes {
		ids[v.VolumeId] = true
	}
	return ids, nil
}

func resourceVpsVolumeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	deviceId := int32(d.Get("device_id").(int))

	// The create call returns a task rather than the volume, so remember the
	// volumes that exist beforehand to identify the new one afterwards.
	before, err := hv.listVpsVolumeIds()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating %d GiB volume for VPS %d", d.Get("size").(int), deviceId)
	payload := swagger.VpsVolumeCreate{DeviceId: deviceId, Size: int32(d.Get("size").(int))}
	task, _, err := hv.client.VPSApi.PostVpsVolumeResource(hv.auth, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /vps/volume")
	}
	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	volumes, _, err := hv.client.VPSApi.GetVpsVolumeResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /vps/volume")
	}
	for _, v := range volumes {
		if !before[v.VolumeId] && v.DeviceId == deviceId {
			d.SetId(v.VolumeId)
			return resourceVpsVolumeRead(ctx, d, m)
		}
	}

	return diag.Errorf("volume for VPS %d was created (task %s) but could not be identified in the volume list; "+
		"import it with `terraform import` once you know its ID", deviceId, task.TaskId)
}

func resourceVpsVolumeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	volume, resp, err := hv.client.VPSApi.GetVpsVolumeIdResource(hv.auth, d.Id(), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] VPS volume %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /vps/volume/%s", d.Id())
	}

	values := map[string]interface{}{
		"device_id":     int(volume.DeviceId),
		"volume_id":     volume.VolumeId,
		"size_bytes":    int(volume.Size),
		"type":          volume.Type_,
		"facility_code": volume.FacilityCode,
		"created_at":    timeString(volume.CreatedAt),
	}
	// The API reports bytes; only derive the GiB size when it divides evenly
	// (or on import, when nothing is known yet) to avoid spurious diffs.
	if d.Get("size").(int) == 0 || int64(volume.Size)%bytesPerGiB == 0 {
		if gib := int64(volume.Size) / bytesPerGiB; gib > 0 {
			values["size"] = int(gib)
		}
	}

	return diag.FromErr(setResourceData(d, values))
}

func resourceVpsVolumeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	timeout := d.Timeout(schema.TimeoutUpdate)

	if d.HasChange("device_id") {
		oldDev, newDev := d.GetChange("device_id")
		if oldDev.(int) != 0 {
			log.Printf("[INFO] Detaching volume %s from VPS %d", d.Id(), oldDev.(int))
			task, _, err := hv.client.VPSApi.PostVpsDetachVolumeResource(hv.auth, d.Id(), int32(oldDev.(int)), nil)
			if err != nil {
				return apiErrorDiag(err, "POST /vps/%d/volume/%s/detach", oldDev.(int), d.Id())
			}
			if _, err := hv.waitForNetworkTask(ctx, task.TaskId, timeout); err != nil {
				return diag.FromErr(err)
			}
		}
		log.Printf("[INFO] Attaching volume %s to VPS %d", d.Id(), newDev.(int))
		task, _, err := hv.client.VPSApi.PostVpsAttachVolumeResource(hv.auth, d.Id(), int32(newDev.(int)), nil)
		if err != nil {
			return apiErrorDiag(err, "POST /vps/%d/volume/%s/attach", newDev.(int), d.Id())
		}
		if _, err := hv.waitForNetworkTask(ctx, task.TaskId, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("size") {
		log.Printf("[INFO] Resizing volume %s to %d GiB", d.Id(), d.Get("size").(int))
		payload := swagger.VpsVolumeUpdate{
			DeviceId: int32(d.Get("device_id").(int)),
			Size:     int32(d.Get("size").(int)),
		}
		task, _, err := hv.client.VPSApi.PutVpsVolumeIdResource(hv.auth, d.Id(), payload, nil)
		if err != nil {
			return apiErrorDiag(err, "PUT /vps/volume/%s", d.Id())
		}
		if _, err := hv.waitForNetworkTask(ctx, task.TaskId, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVpsVolumeRead(ctx, d, m)
}

func resourceVpsVolumeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Deleting VPS volume %s", d.Id())
	payload := swagger.VpsVolumeDelete{DeviceId: int32(d.Get("device_id").(int))}
	task, resp, err := hv.client.VPSApi.DeleteVpsVolumeIdResource(hv.auth, d.Id(), payload, nil)
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /vps/volume/%s", d.Id())
	}
	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func dataSourceVpsVolume() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up a VPS volume, for example by `device_id` and `type`.",
		ReadContext: dataSourceVpsVolumeRead,
		Schema: map[string]*schema.Schema{
			"filter":        dataSourceFiltersSchema(),
			"first":         dataSourceFilterFirstSchema(),
			"volume_id":     computedString("Unique ID of the volume."),
			"device_id":     computedInt("VPS instance the volume is attached to."),
			"size":          computedInt("Size in bytes."),
			"type":          computedString("`ROOT` or `DATADISK`."),
			"facility_code": computedString("Facility of the volume."),
			"created_at":    computedString("When the volume was created."),
		},
	}
}

func dataSourceVpsVolumeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	volumes, _, err := hv.client.VPSApi.GetVpsVolumeResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /vps/volume")
	}

	return readFilteredList(d, volumes, dataSourceVpsVolume().Schema, "volume_id", nil)
}
