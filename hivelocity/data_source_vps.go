package hivelocity

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func dataSourceVps() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up a VPS instance by `device_id`, or by filtering the instance list (for " +
			"example on `facility_code` or `status`).",
		ReadContext: dataSourceVpsRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"first":  dataSourceFilterFirstSchema(),
			"device_id": {
				Description:   "Device ID of the instance.",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"filter"},
			},
			"status":        computedString("Instance status."),
			"vcpu":          computedInt("Number of vCPUs."),
			"memory":        computedInt("Memory in MB."),
			"facility_code": computedString("Facility."),
			"template_id":   computedString("Template used for the install."),
			"iso_id":        computedString("Attached ISO, if any."),
			"cloud_init":    computedString("Cloud-init script in effect."),
			"ssh_key_ids":   computedIntList("Authorized SSH key IDs."),
			"nics": {
				Description: "Network interfaces.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        vpsNicSchema(),
			},
		},
	}
}

func setVpsDataSource(d *schema.ResourceData, vps swagger.Vps) diag.Diagnostics {
	nics, err := flattenVpsNics(vps.Nics)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := setResourceData(d, map[string]interface{}{
		"device_id":     int(vps.DeviceId),
		"status":        vps.Status,
		"vcpu":          int(vps.Vcpu),
		"memory":        int(vps.Memory),
		"facility_code": vps.FacilityCode,
		"template_id":   vps.TemplateId,
		"iso_id":        vps.IsoId,
		"cloud_init":    vps.CloudInit,
		"ssh_key_ids":   flattenInt32List(vps.SshKeyIds),
		"nics":          nics,
	}); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprint(vps.DeviceId))
	return nil
}

func dataSourceVpsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	if v, ok := d.GetOk("device_id"); ok {
		vps, _, err := hv.client.VPSApi.GetVpsIdResource(hv.auth, int32(v.(int)), nil)
		if err != nil {
			return apiErrorDiag(err, "GET /vps/%d", v.(int))
		}
		return setVpsDataSource(d, vps)
	}

	instances, _, err := hv.client.VPSApi.GetVpsResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /vps")
	}

	// Filter on the flat attributes, then set the full object (nics included).
	list := make([]map[string]interface{}, 0, len(instances))
	for i, vps := range instances {
		list = append(list, map[string]interface{}{
			"_index":        i,
			"device_id":     int(vps.DeviceId),
			"status":        vps.Status,
			"vcpu":          int(vps.Vcpu),
			"memory":        int(vps.Memory),
			"facility_code": vps.FacilityCode,
			"template_id":   vps.TemplateId,
			"iso_id":        vps.IsoId,
		})
	}
	if len(list) == 0 {
		return diag.Errorf("no VPS instances found")
	}
	item, err := doFiltering(d, list, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if item == nil {
		return diag.Errorf("no VPS instance matched the given filters")
	}

	return setVpsDataSource(d, instances[item["_index"].(int)])
}

func dataSourceVpsEvents() *schema.Resource {
	return &schema.Resource{
		Description: "Lists the event log of a VPS instance.",
		ReadContext: dataSourceVpsEventsRead,
		Schema: map[string]*schema.Schema{
			"device_id": {
				Description: "Device ID of the instance.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"events": {
				Description: "Events, as returned by the API.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_id":   computedString("Unique ID of the event."),
						"type":       computedString("Event type."),
						"level":      computedString("INFO, WARNING or ERROR."),
						"state":      computedString("Scheduled, Started or Completed."),
						"created_at": computedString("When the event was recorded."),
					},
				},
			},
		},
	}
}

func dataSourceVpsEventsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	deviceId := d.Get("device_id").(int)

	events, _, err := hv.client.VPSApi.GetVpsIdEventsResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		return apiErrorDiag(err, "GET /vps/%d/events", deviceId)
	}

	list := make([]interface{}, 0, len(events))
	for _, e := range events {
		list = append(list, map[string]interface{}{
			"event_id":   e.EventId,
			"type":       e.Type_,
			"level":      e.Level,
			"state":      e.State,
			"created_at": timeString(e.CreatedAt),
		})
	}
	if err := d.Set("events", list); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprint(deviceId))
	return nil
}
