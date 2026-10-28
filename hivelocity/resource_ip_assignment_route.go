package hivelocity

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceIpAssignmentRoute() *schema.Resource {
	return &schema.Resource{
		Description: "Statically routes an existing IP assignment (subnet) to a next-hop IP, typically an " +
			"address already configured on one of your devices. Destroying the resource clears the " +
			"assignment's routing. New IP space itself is requested from Hivelocity support (or " +
			"`POST /ip/`, which opens a ticket) and then looked up with the `hivelocity_ip_assignment` " +
			"data source.",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceIpAssignmentRouteCreate,
		ReadContext:   resourceIpAssignmentRouteRead,
		UpdateContext: resourceIpAssignmentRouteUpdate,
		DeleteContext: resourceIpAssignmentRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"assignment_id": {
				Description: "ID of the IP assignment to route.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"next_hop_ip": {
				Description: "IPv4 or IPv6 address traffic for the assignment is routed to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"subnet":        computedString("CIDR of the assignment."),
			"facility_code": computedString("Facility of the assignment."),
			"version":       computedInt("IP version (4 or 6)."),
			"device_id":     computedInt("Device currently receiving the assignment's traffic, if any."),
			"port_id":       computedInt("Port currently receiving the assignment's traffic, if any."),
			"vlan_id":       computedInt("VLAN currently receiving the assignment's traffic, if any."),
		},
	}
}

func (hv *Client) routeIpAssignment(ctx context.Context, assignmentId int32, nextHop string, timeout time.Duration) error {
	payload := swagger.IpAssignmentPut{NextHopIp: nextHop}
	task, _, err := hv.client.IPAssignmentApi.PutIpAssignmentIdResource(hv.auth, assignmentId, payload, nil)
	if err != nil {
		return formatSwaggerError(err, "PUT /ip/%d", assignmentId)
	}
	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, timeout); err != nil {
		return fmt.Errorf("waiting for IP assignment %d routing: %w", assignmentId, err)
	}
	return nil
}

func resourceIpAssignmentRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	assignmentId := int32(d.Get("assignment_id").(int))

	log.Printf("[INFO] Routing IP assignment %d to %s", assignmentId, d.Get("next_hop_ip").(string))
	if err := hv.routeIpAssignment(ctx, assignmentId, d.Get("next_hop_ip").(string), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(assignmentId))
	return resourceIpAssignmentRouteRead(ctx, d, m)
}

func resourceIpAssignmentRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	assignmentId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ip, resp, err := hv.client.IPAssignmentApi.GetIpAssignmentIdResource(hv.auth, int32(assignmentId), nil)
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /ip/%d", assignmentId)
	}

	if ip.NextHopIp == "" {
		log.Printf("[WARN] IP assignment %d is no longer statically routed, removing from state", assignmentId)
		d.SetId("")
		return nil
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"assignment_id": assignmentId,
		"next_hop_ip":   ip.NextHopIp,
		"subnet":        ip.Subnet,
		"facility_code": ip.FacilityCode,
		"version":       int(ip.Version),
		"device_id":     int(ip.DeviceId),
		"port_id":       int(ip.PortId),
		"vlan_id":       int(ip.VlanId),
	}))
}

func resourceIpAssignmentRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	assignmentId := int32(d.Get("assignment_id").(int))

	if d.HasChange("next_hop_ip") {
		if err := hv.routeIpAssignment(ctx, assignmentId, d.Get("next_hop_ip").(string), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIpAssignmentRouteRead(ctx, d, m)
}

func resourceIpAssignmentRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	assignmentId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Clearing routing of IP assignment %d", assignmentId)
	task, resp, err := hv.client.IPAssignmentApi.PostIpAssignmentIdClearResource(hv.auth, int32(assignmentId), nil)
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "POST /ip/%d/clear", assignmentId)
	}

	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
