package hivelocity

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceVlan() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a VLAN in a facility and the device ports (or bonds) and IP assignments attached " +
			"to it. Port and IP changes are applied asynchronously by the network automation; the provider waits " +
			"for the resulting task to finish.",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceVlanCreate,
		ReadContext:   resourceVlanRead,
		UpdateContext: resourceVlanUpdate,
		DeleteContext: resourceVlanDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"facility_code": {
				Description: "Facility where the VLAN is created, for example `NYC1`. Only devices from the " +
					"same facility can be attached.",
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Description: "`public` VLANs can have IP assignments routed to them and be reachable from the " +
					"internet. `private` VLANs never carry public IPs. All VLAN traffic is subject to bandwidth " +
					"billing except private VLAN traffic on unbonded devices.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
			},
			"port_ids": {
				Description: "IDs of device ports or bonds to include in this VLAN.",
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
			},
			"ip_ids": {
				Description: "IDs of IP assignments routed to this VLAN (public VLANs only).",
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
			},
			"vlan_tag": {
				Description: "The VLAN tag from the switch. Use this value when configuring your OS interfaces.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"tag_id": {
				Description: "Alias of `vlan_tag`, kept for backwards compatibility.",
				Deprecated:  "Use vlan_tag instead.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"automated": {
				Description: "True when the VLAN can be automated through the API. If false, contact support " +
					"to enable automation.",
				Type:     schema.TypeBool,
				Computed: true,
			},
			"q_in_q": {
				Description: "True when the VLAN is configured in Q-in-Q mode. Automation is not supported on " +
					"Q-in-Q VLANs.",
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// vlanUpdatePayload is sent through rawJSON so that empty lists are serialized
// as `[]` (the generated model drops them with omitempty, but the API requires
// `portIds`).
type vlanUpdatePayload struct {
	PortIds []int32 `json:"portIds"`
	IpIds   []int32 `json:"ipIds"`
}

// updateVlanMembers sets the full list of ports and IPs on a VLAN and waits for
// the network task to complete.
func (hv *Client) updateVlanMembers(ctx context.Context, vlanId int32, portIds, ipIds []int32, timeout time.Duration) error {
	payload := vlanUpdatePayload{PortIds: portIds, IpIds: ipIds}
	if payload.PortIds == nil {
		payload.PortIds = []int32{}
	}
	if payload.IpIds == nil {
		payload.IpIds = []int32{}
	}

	var task swagger.NetworkTaskDump
	if _, err := hv.rawJSON(ctx, "PUT", fmt.Sprintf("/vlan/%d", vlanId), payload, &task); err != nil {
		return err
	}

	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, timeout); err != nil {
		return fmt.Errorf("waiting for VLAN %d update: %w", vlanId, err)
	}
	return nil
}

func resourceVlanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := swagger.VlanCreate{
		FacilityCode: d.Get("facility_code").(string),
		Type_:        d.Get("type").(string),
	}

	vlan, _, err := hv.client.VLANApi.PostVlanResource(hv.auth, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /vlan")
	}

	log.Printf("[INFO] Created VLAN ID: %d", vlan.VlanId)
	d.SetId(fmt.Sprint(vlan.VlanId))

	portIds := expandInt32List(d.Get("port_ids"))
	ipIds := expandInt32List(d.Get("ip_ids"))

	if len(portIds) > 0 || len(ipIds) > 0 {
		if err := hv.updateVlanMembers(ctx, vlan.VlanId, portIds, ipIds, d.Timeout(schema.TimeoutCreate)); err != nil {
			// Roll back the empty VLAN so a retry starts clean.
			diags := diag.FromErr(err)
			if _, delErr := hv.client.VLANApi.DeleteVlanIdResource(hv.auth, vlan.VlanId); delErr != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("VLAN %d was created but could not be cleaned up", vlan.VlanId),
					Detail:   formatSwaggerError(delErr, "DELETE /vlan/%d", vlan.VlanId).Error(),
				})
			} else {
				d.SetId("")
			}
			return diags
		}
	}

	return resourceVlanRead(ctx, d, m)
}

func resourceVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	vlanId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vlan, response, err := hv.client.VLANApi.GetVlanIdResource(hv.auth, int32(vlanId), nil)
	if err != nil {
		if isNotFound(response) {
			log.Printf("[WARN] VLAN %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /vlan/%d", vlanId)
	}

	values := map[string]interface{}{
		"facility_code": vlan.FacilityCode,
		"type":          vlan.Type_,
		"port_ids":      SetFromInt32List(vlan.PortIds),
		"ip_ids":        SetFromInt32List(vlan.IpIds),
		"vlan_tag":      int(vlan.VlanTag),
		"tag_id":        int(vlan.VlanTag),
		"automated":     vlan.Automated,
		"q_in_q":        vlan.QInQ,
	}
	if err := setResourceData(d, values); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVlanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	vlanId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("port_ids", "ip_ids") {
		log.Printf("[INFO] Updating members of VLAN ID: %s", d.Id())
		portIds := expandInt32List(d.Get("port_ids"))
		ipIds := expandInt32List(d.Get("ip_ids"))
		if err := hv.updateVlanMembers(ctx, int32(vlanId), portIds, ipIds, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVlanRead(ctx, d, m)
}

func resourceVlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	vlanId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	vlan, response, err := hv.client.VLANApi.GetVlanIdResource(hv.auth, int32(vlanId), nil)
	if err != nil {
		if isNotFound(response) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /vlan/%d", vlanId)
	}

	// Ports and IPs must be detached before the VLAN can be removed.
	if len(vlan.PortIds) > 0 || len(vlan.IpIds) > 0 {
		log.Printf("[INFO] Clearing ports and IPs before deleting VLAN %d", vlanId)
		task, _, err := hv.client.VLANApi.PostVlanIdClearResource(hv.auth, int32(vlanId), nil)
		if err != nil {
			return apiErrorDiag(err, "POST /vlan/%d/clear", vlanId)
		}
		if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO] Deleting VLAN ID: %d", vlanId)
	if response, err := hv.client.VLANApi.DeleteVlanIdResource(hv.auth, int32(vlanId)); err != nil {
		if isNotFound(response) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /vlan/%d", vlanId)
	}

	d.SetId("")
	return nil
}
