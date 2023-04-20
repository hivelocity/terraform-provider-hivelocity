package hivelocity

import (
	"context"
	"fmt"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIP() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		CreateContext: resourceIPCreate,
		ReadContext:   resourceIPRead,
		UpdateContext: resourceIPUpdate,
		DeleteContext: resourceIPDelete,
		Schema: map[string]*schema.Schema{
			"facility_code": &schema.Schema{
				Description: "Facility code of this assignment.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"purpose": &schema.Schema{
				Description: "List the intended use of each IP in the subnet. It's required by ICANN.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"prefix_length": &schema.Schema{
				Description: "The prefix length of the subnet. For example: `27`.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"ip_assignment_id": &schema.Schema{
				Description: "The ID of the IP assignment.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"next_hop_ip": &schema.Schema{
				Description: "The next hop IP address of the subnet.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceIPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	hv, _ := m.(*Client)

	payload := swagger.IpAssignmentRequest{
		FacilityCode: d.Get("facility_code").(string),
		Purpose:      d.Get("purpose").(string),
		PrefixLength: d.Get("prefix_length").(int32),
	}

	_, err := hv.client.IPAssignmentApi.PostIpAssignmentResource(hv.auth, payload)

	if err != nil {
		return diag.FromErr(formatSwaggerError(err, "POST /ip/"))
	}
	log.Printf("[INFO] IP Assignment requested")
	return diags
}

func resourceIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	ipAssignmentId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ipAssignmentResponse, httpResponse, err := hv.client.IPAssignmentApi.GetIpAssignmentIdResource(hv.auth, int32(ipAssignmentId), nil)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /ip/%d failed! (%s)\n\n %s", ipAssignmentId, err, myErr.Body())
	}

	d.Set("version", ipAssignmentResponse.Version)
	d.Set("assignment_id", ipAssignmentResponse.AssignmentId)
	d.Set("client_id", ipAssignmentResponse.ClientId)
	d.Set("subnet", ipAssignmentResponse.Subnet)
	d.Set("netmask", ipAssignmentResponse.Netmask)
	d.Set("broadcast_ip", ipAssignmentResponse.BroadcastIp)
	d.Set("gateway_ip", ipAssignmentResponse.GatewayIp)
	d.Set("first_usable_ip", ipAssignmentResponse.FirstUsableIp)
	d.Set("last_usable_ip", ipAssignmentResponse.LastUsableIp)
	d.Set("usable_ips", ipAssignmentResponse.UsableIps)
	d.Set("facility_code", ipAssignmentResponse.FacilityCode)
	d.Set("description", ipAssignmentResponse.Description)
	d.Set("device_id", ipAssignmentResponse.DeviceId)
	d.Set("port_id", ipAssignmentResponse.PortId)
	d.Set("vlan_id", ipAssignmentResponse.VlanId)
	d.Set("next_hop_ip", ipAssignmentResponse.NextHopIp)

	return nil
}

func resourceIPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	ipAssignmentId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload := swagger.IpAssignmentPut{}
	if d.HasChange("next_hop_ip") {
		payload.NextHopIp = d.Get("next_hop_ip").(string)
	}

	_, _, err = hv.client.IPAssignmentApi.PutIpAssignmentIdResource(hv.auth, int32(ipAssignmentId), payload, nil)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /ip/%d failed! (%s)\n\n %s", fmt.Sprint(ipAssignmentId), err, myErr.Body())
	}

	return resourceIPRead(ctx, d, m)
}

func resourceIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	hv, _ := m.(*Client)

	ipAssignmentId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := hv.client.IPAssignmentApi.DeleteIpAssignmentIdResource(hv.auth, int32(ipAssignmentId))
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /ip/%d failed! (%s)\n\n %s", ipAssignmentId, err, myErr.Body())
	}

	d.SetId("")
	return diags
}
