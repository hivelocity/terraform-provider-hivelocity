package hivelocity

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceVpsSnapshot() *schema.Resource {
	return &schema.Resource{
		Description: "Takes a snapshot of a VPS volume. Snapshots are immutable; restoring one is an " +
			"operational action done outside Terraform.",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		CreateContext: resourceVpsSnapshotCreate,
		ReadContext:   resourceVpsSnapshotRead,
		DeleteContext: resourceVpsSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Description: "ID of the volume to snapshot.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Name of the snapshot.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"facility_code": {
				Description: "Facility of the volume, for example `TPA1`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"snapshot_id":   computedString("Unique ID of the snapshot."),
			"size":          computedInt("Size of the snapshot."),
			"state":         computedString("State of the snapshot."),
			"snapshot_type": computedString("Type of the snapshot."),
			"created_at":    computedString("When the snapshot was created."),
		},
	}
}

func resourceVpsSnapshotCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	volumeId := d.Get("volume_id").(string)
	name := d.Get("name").(string)

	before, _, err := hv.client.VPSApi.GetVpsSnapshotResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /vps/snapshot")
	}
	existing := map[string]bool{}
	for _, s := range before {
		existing[s.SnapshotId] = true
	}

	log.Printf("[INFO] Creating snapshot %q of volume %s", name, volumeId)
	payload := swagger.VpsSnapshotCreate{
		VolumeId:     volumeId,
		Name:         name,
		FacilityCode: d.Get("facility_code").(string),
	}
	task, _, err := hv.client.VPSApi.PostVpsSnapshotResource(hv.auth, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /vps/snapshot")
	}
	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	after, _, err := hv.client.VPSApi.GetVpsSnapshotResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /vps/snapshot")
	}
	var found *swagger.VpsSnapshot
	for i := range after {
		s := &after[i]
		if existing[s.SnapshotId] || s.VolumeId != volumeId || s.Name != name {
			continue
		}
		if found == nil || s.CreatedAt.After(found.CreatedAt) {
			found = s
		}
	}
	if found == nil {
		return diag.Errorf("snapshot %q of volume %s was created (task %s) but could not be identified in the "+
			"snapshot list; import it once you know its ID", name, volumeId, task.TaskId)
	}

	d.SetId(found.SnapshotId)
	return resourceVpsSnapshotRead(ctx, d, m)
}

func resourceVpsSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	snapshot, resp, err := hv.client.VPSApi.GetVpsSnapshotIdResource(hv.auth, d.Id(), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] VPS snapshot %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /vps/snapshot/%s", d.Id())
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"volume_id":     snapshot.VolumeId,
		"name":          snapshot.Name,
		"snapshot_id":   snapshot.SnapshotId,
		"size":          int(snapshot.Size),
		"state":         snapshot.State,
		"snapshot_type": snapshot.SnapshotType,
		"created_at":    timeString(snapshot.CreatedAt),
	}))
}

func resourceVpsSnapshotDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Deleting VPS snapshot %s", d.Id())
	task, resp, err := hv.client.VPSApi.DeleteVpsSnapshotIdResource(hv.auth, d.Id(), nil)
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /vps/snapshot/%s", d.Id())
	}
	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func dataSourceVpsSnapshot() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up a VPS volume snapshot, for example by `name` or `volume_id`.",
		ReadContext: dataSourceVpsSnapshotRead,
		Schema: map[string]*schema.Schema{
			"filter":        dataSourceFiltersSchema(),
			"first":         dataSourceFilterFirstSchema(),
			"snapshot_id":   computedString("Unique ID of the snapshot."),
			"volume_id":     computedString("Volume the snapshot was taken of."),
			"name":          computedString("Name of the snapshot."),
			"size":          computedInt("Size of the snapshot."),
			"state":         computedString("State of the snapshot."),
			"snapshot_type": computedString("Type of the snapshot."),
			"created_at":    computedString("When the snapshot was created."),
		},
	}
}

func dataSourceVpsSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	snapshots, _, err := hv.client.VPSApi.GetVpsSnapshotResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /vps/snapshot")
	}

	return readFilteredList(d, snapshots, dataSourceVpsSnapshot().Schema, "snapshot_id", nil)
}

func resourceVpsSnapshotSchedule() *schema.Resource {
	return &schema.Resource{
		Description: "Schedules automatic snapshots of a VPS volume. Schedules cannot be modified through " +
			"the API, so any change replaces the schedule.",
		CreateContext: resourceVpsSnapshotScheduleCreate,
		ReadContext:   resourceVpsSnapshotScheduleRead,
		DeleteContext: resourceVpsSnapshotScheduleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Description: "ID of the volume to snapshot.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"facility_code": {
				Description: "Facility of the volume, for example `TPA1`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"interval_type": {
				Description:  "`DAILY`, `WEEKLY` or `MONTHLY`.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DAILY", "WEEKLY", "MONTHLY"}, false),
			},
			"hour": {
				Description:  "Hour of the day to run (0-23).",
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 23),
			},
			"minute": {
				Description:  "Minute of the hour to run (0-59).",
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 59),
			},
			"max_snapshots": {
				Description:  "Maximum number of snapshots to retain for this schedule.",
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"timezone": {
				Description: "IANA timezone of the schedule, for example `America/New_York`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"weekday": {
				Description:  "Day of the week to run (1-7, 1 being Monday). Required for `WEEKLY`.",
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 7),
			},
			"day": {
				Description:  "Day of the month to run (1-28). Required for `MONTHLY`.",
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 28),
			},
			"snapshot_schedule_id": computedString("Unique ID of the schedule."),
		},
	}
}

func resourceVpsSnapshotScheduleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := swagger.VpsSnapshotScheduleCreate{
		VolumeId:     d.Get("volume_id").(string),
		FacilityCode: d.Get("facility_code").(string),
		IntervalType: d.Get("interval_type").(string),
		Hour:         int32(d.Get("hour").(int)),
		Minute:       int32(d.Get("minute").(int)),
		MaxSnapshots: int32(d.Get("max_snapshots").(int)),
		Timezone:     d.Get("timezone").(string),
		Weekday:      int32(d.Get("weekday").(int)),
		Day:          int32(d.Get("day").(int)),
	}
	schedule, _, err := hv.client.VPSApi.PostVpsSnapshotScheduleResource(hv.auth, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /vps/snapshotSchedule")
	}
	if schedule.SnapshotScheduleId == "" {
		return diag.Errorf("POST /vps/snapshotSchedule did not return a snapshotScheduleId")
	}

	d.SetId(schedule.SnapshotScheduleId)
	return resourceVpsSnapshotScheduleRead(ctx, d, m)
}

func resourceVpsSnapshotScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	schedule, resp, err := hv.client.VPSApi.GetVpsSnapshotScheduleIdResource(hv.auth, d.Id(), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] VPS snapshot schedule %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /vps/snapshotSchedule/%s", d.Id())
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"volume_id":            schedule.VolumeId,
		"interval_type":        schedule.IntervalType,
		"hour":                 int(schedule.Hour),
		"minute":               int(schedule.Minute),
		"max_snapshots":        int(schedule.MaxSnapshots),
		"timezone":             schedule.Timezone,
		"weekday":              int(schedule.Weekday),
		"day":                  int(schedule.Day),
		"snapshot_schedule_id": schedule.SnapshotScheduleId,
	}))
}

func resourceVpsSnapshotScheduleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Deleting VPS snapshot schedule %s", d.Id())
	if resp, err := hv.client.VPSApi.DeleteVpsSnapshotScheduleIdResource(hv.auth, d.Id(), nil); err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /vps/snapshotSchedule/%s", d.Id())
	}

	d.SetId("")
	return nil
}
