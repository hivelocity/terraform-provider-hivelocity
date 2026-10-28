package hivelocity

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceDnsRecordA() *schema.Resource {
	return &schema.Resource{
		Description: "Manages an A record set in a Hivelocity-hosted zone. The API addresses A records by " +
			"zone name and record name rather than by numeric ID, and a record name can carry several " +
			"addresses (round robin). The resource ID is `<domain>/<name>`.",
		CreateContext: resourceDnsRecordACreate,
		ReadContext:   resourceDnsRecordARead,
		UpdateContext: resourceDnsRecordAUpdate,
		DeleteContext: resourceDnsRecordADelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDnsRecordAImport,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Description: "Zone name the record belongs to, for example `example.com`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Record name, for example `www` or `@` for the zone apex.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"addresses": {
				Description: "IPv4 addresses the record resolves to.",
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateIPv4,
				},
			},
			"ttl": {
				Description: "Time to live of the record, in seconds.",
				Type:        schema.TypeInt,
				Required:    true,
			},
		},
	}
}

func validateIPv4(val interface{}, path cty.Path) diag.Diagnostics {
	ip := net.ParseIP(val.(string))
	if ip == nil || ip.To4() == nil {
		return diag.Errorf("invalid IPv4 address: %q", val.(string))
	}
	return nil
}

func validateIPv6(val interface{}, path cty.Path) diag.Diagnostics {
	ip := net.ParseIP(val.(string))
	if ip == nil || ip.To4() != nil {
		return diag.Errorf("invalid IPv6 address: %q", val.(string))
	}
	return nil
}

func dnsRecordAId(domain, name string) string {
	return domain + "/" + name
}

func parseDnsRecordAId(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected ID format %q, expected <domain>/<record name>", id)
	}
	return parts[0], parts[1], nil
}

func aRecordPayload(d *schema.ResourceData) swagger.ARecord {
	return swagger.ARecord{
		Name:      d.Get("name").(string),
		Addresses: expandStringList(d.Get("addresses")),
		Ttl:       int32(d.Get("ttl").(int)),
	}
}

func resourceDnsRecordACreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	domain := d.Get("domain").(string)

	record, _, err := hv.client.DomainsApi.PostARecordResource(hv.auth, domain, aRecordPayload(d), nil)
	if err != nil {
		return apiErrorDiag(err, "POST /domains/%s/a-record", domain)
	}

	name := record.Name
	if name == "" {
		name = d.Get("name").(string)
	}
	d.SetId(dnsRecordAId(domain, name))
	log.Printf("[INFO] Created A record %s", d.Id())

	return resourceDnsRecordARead(ctx, d, m)
}

func resourceDnsRecordARead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	domain, name, err := parseDnsRecordAId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	record, resp, err := hv.client.DomainsApi.GetARecordIdResource(hv.auth, domain, name, nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] A record %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /domains/%s/a-record/%s", domain, name)
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"domain":    domain,
		"name":      name,
		"addresses": flattenStringList(record.Addresses),
		"ttl":       int(record.Ttl),
	}))
}

func resourceDnsRecordAUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	domain, name, err := parseDnsRecordAId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if _, _, err := hv.client.DomainsApi.PutARecordIdResource(hv.auth, domain, name, aRecordPayload(d), nil); err != nil {
		return apiErrorDiag(err, "PUT /domains/%s/a-record/%s", domain, name)
	}

	return resourceDnsRecordARead(ctx, d, m)
}

func resourceDnsRecordADelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	domain, name, err := parseDnsRecordAId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting A record %s", d.Id())
	if resp, err := hv.client.DomainsApi.DeleteARecordIdResource(hv.auth, domain, name); err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /domains/%s/a-record/%s", domain, name)
	}

	d.SetId("")
	return nil
}

func resourceDnsRecordAImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	domain, name, err := parseDnsRecordAId(d.Id())
	if err != nil {
		return nil, err
	}
	if err := d.Set("domain", domain); err != nil {
		return nil, err
	}
	if err := d.Set("name", name); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
