package hivelocity

import (
	"context"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceIPAssignment() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceIPAssignmentCreate,
		ReadContext:   resourceIPAssignmentRead,
		UpdateContext: resourceIPAssignmentUpdate,
		DeleteContext: resourceIPAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"version": {
				Type:        schema.TypeInt,
				Description: "IP version of this assignment (ipv4 or ipv6)",
				Computed:    true,
			},
			"assignment_id": {
				Type:        schema.TypeInt,
				Description: "Unique ID of this assignment",
				Computed:    true,
			},
			"client_id": {
				Type:        schema.TypeInt,
				Description: "Unique ID of the assignment owner",
				Computed:    true,
			},
			"subnet": {
				Type:        schema.TypeString,
				Description: "CIDR representation of this assignment.",
				Computed:    true,
			},
			"netmask": {
				Type:        schema.TypeString,
				Description: "Netmask for this assignment.",
				Computed:    true,
			},
			"broadcast_ip": {
				Type:        schema.TypeString,
				Description: "Broadcast address for this subnet.",
				Computed:    true,
			},
			"gateway_ip": {
				Type:        schema.TypeString,
				Description: "Gateway address for this subnet.",
				Computed:    true,
			},
			"dns_resolvers": {
				Type:        schema.TypeList,
				Description: "DNS resolvers for this subnet. Only applicable to VPS instances.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"first_usable_ip": {
				Type:        schema.TypeString,
				Description: "First address available for use by Devices on this subnet.",
				Computed:    true,
			},
			"last_usable_ip": {
				Type:        schema.TypeString,
				Description: "Last address available for use by Devices on this subnet.",
				Computed:    true,
			},
			"usable_ips": {
				Type:        schema.TypeList,
				Description: "A list of usable IP addresses on assignment. It is only filled for IPv4 subnets.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"facility_code": {
				Type:        schema.TypeString,
				Description: "Facility code of this assignment.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description for this assignment.",
				Computed:    true,
			},
			"device_id": {
				Type:        schema.TypeInt,
				Description: "The device receiving traffic from the assignment.",
				Computed:    true,
			},
			"port_id": {
				Type:        schema.TypeInt,
				Description: "The port receiving traffic from the assignment.",
				Computed:    true,
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Description: "The VLAN receiving traffic from the assignment.",
				Computed:    true,
			},
			"next_hop_ip": {
				Type:        schema.TypeString,
				Description: "Next Hop IP address, if this assignment is statically routed.",
				Computed:    true,
			},
			"purpose": {
				Type:        schema.TypeString,
				Description: "List the intended use of each IP in the subnet. Required by ICANN",
				Required:    true,
			},
			"prefix_length": {
				Type:        schema.TypeInt,
				Description: "The prefix length of the subnet. For example: `/27`",
				Required:    true,
			},
		},
	}
}

func resourceIPAssignmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	hv, _ := m.(*Client)

	log.Printf("[INFO] Creating IP Assignment")
	payload := swagger.IpAssignmentRequest{
		FacilityCode: d.Get("facility_code").(string),
		Purpose:      d.Get("purpose").(string),
		PrefixLength: int32(d.Get("prefix_length").(int)),
	}

	_, err := hv.client.IPAssignmentApi.PostIpAssignmentResource(hv.auth, payload)

	if err != nil {
		d.SetId("")
		return diag.FromErr(formatSwaggerError(err, "POST /ip"))
	}

	log.Printf("[INFO] IP Assignment requested")
	return diags
}

func resourceIPAssignmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	hv, _ := m.(*Client)

	ipAssignmentId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ipResp, httpResponse, err := hv.client.IPAssignmentApi.GetIpAssignmentIdResource(hv.auth, int32(ipAssignmentId), nil)

	if err != nil {
		if httpResponse.StatusCode == 404 {
			log.Printf("[INFO] IP Assignment not found")
			d.SetId("")
			return diags
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /ip/%d failed! (%s)\n\n %s", ipAssignmentId, err, myErr.Body())
	}

	d.Set("version", ipResp.Version)
	d.Set("assignment_id", ipResp.AssignmentId)
	d.Set("client_id", ipResp.ClientId)
	d.Set("subnet", ipResp.Subnet)
	d.Set("netmask", ipResp.Netmask)
	d.Set("broadcast_ip", ipResp.BroadcastIp)
	d.Set("gateway_ip", ipResp.GatewayIp)
	d.Set("dns_resolvers", ipResp.DnsResolvers)
	d.Set("first_usable_ip", ipResp.FirstUsableIp)
	d.Set("last_usable_ip", ipResp.LastUsableIp)
	d.Set("usable_ips", ipResp.UsableIps)
	d.Set("facility_code", ipResp.FacilityCode)
	d.Set("description", ipResp.Description)
	d.Set("device_id", ipResp.DeviceId)
	d.Set("port_id", ipResp.PortId)
	d.Set("vlan_id", ipResp.VlanId)
	d.Set("next_hop_ip", ipResp.NextHopIp)

	return nil
}

func resourceIPAssignmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[INFO] Updating IP Assignment")
	var payload swagger.IpAssignmentPut
	hv, _ := m.(*Client)

	ipAssignmentId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	fieldToPayloadMap := map[string]string{
		// "version":         "Version",
		// "assignment_id":   "AssignmentId",
		// "client_id":       "ClientId",
		// "subnet":          "Subnet",
		// "netmask":         "Netmask",
		// "broadcast_ip":    "BroadcastIp",
		// "gateway_ip":      "GatewayIp",
		// "dns_resolvers":   "DnsResolvers",
		// "first_usable_ip": "FirstUsableIp",
		// "last_usable_ip":  "LastUsableIp",
		// "usable_ips":      "UsableIps",
		// "facility_code":   "FacilityCode",
		// "description":     "Description",
		// "device_id":       "DeviceId",
		// "port_id":         "PortId",
		// "vlan_id":         "VlanId",
		"next_hop_ip": "NextHopIp",
		// "purpose":         "Purpose",
		// "prefix_length":   "PrefixLength",
	}

	for schemaField, payloadField := range fieldToPayloadMap {
		if d.HasChange(schemaField) {
			if newValue, ok := d.GetOk(schemaField); ok {
				payloadValue, _ := newValue.(string)
				reflect.ValueOf(&payload).Elem().FieldByName(payloadField).SetString(payloadValue)
			}
		}
	}

	// Update IP Assignment
	task, _, err := hv.client.IPAssignmentApi.PutIpAssignmentIdResource(hv.auth, int32(ipAssignmentId), payload, nil)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /ip/%d failed! (%s)\n\n %s", ipAssignmentId, err, myErr.Body())
	}

	// Wait for IP Assignment to be updated
	log.Printf("[INFO] Waiting for IP Assignment (%d) to be updated", ipAssignmentId)
	if _, err := waitForNetworkTaskByClient(hv.auth, d.Timeout(schema.TimeoutUpdate), hv, task.TaskId); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] IP Assignment (%d) updated", ipAssignmentId)
	return resourceIPAssignmentRead(ctx, d, m)
}

func resourceIPAssignmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting IP Assignment")
	var diags diag.Diagnostics
	hv, _ := m.(*Client)

	ipAssignmentId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := hv.client.IPAssignmentApi.DeleteIpAssignmentIdResource(hv.auth, int32(ipAssignmentId))

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			log.Printf("[INFO] IP Assignment not found")
			return diags
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /ip/%d failed! (%s)\n\n %s", ipAssignmentId, err, myErr.Body())
	}

	d.SetId("")
	log.Printf("[INFO] IP Assignment (%d) deleted", ipAssignmentId)
	return diags
}
