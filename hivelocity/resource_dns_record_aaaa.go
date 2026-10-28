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

func resourceDnsRecordAAAA() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an AAAA record in a Hivelocity-hosted zone. Import with `<domain_id>:<record_id>`.",
		CreateContext: resourceDnsRecordAAAACreate,
		ReadContext:   resourceDnsRecordAAAARead,
		UpdateContext: resourceDnsRecordAAAAUpdate,
		DeleteContext: resourceDnsRecordAAAADelete,
		Importer: &schema.ResourceImporter{
			StateContext: dnsRecordImport,
		},
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Description: "ID of the zone (`hivelocity_dns_domain.domain_id`).",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Record name, for example `www`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"address": {
				Description:      "IPv6 address the record resolves to.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateIPv6,
			},
			"ttl": {
				Description: "Time to live of the record, in seconds.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"record_id": computedInt("Unique ID of the record."),
			"type":      computedString("Record type as reported by the API."),
		},
	}
}

// dnsRecordImport handles "<domain_id>:<record_id>" import IDs for the record
// resources that are addressed by numeric IDs.
func dnsRecordImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	domainId, recordId, err := parseTwoPartId(d.Id(), "<domain_id>:<record_id>")
	if err != nil {
		return nil, err
	}
	if err := d.Set("domain_id", domainId); err != nil {
		return nil, err
	}
	d.SetId(fmt.Sprint(recordId))
	return []*schema.ResourceData{d}, nil
}

func resourceDnsRecordAAAACreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	domainId := int32(d.Get("domain_id").(int))

	payload := swagger.AaaaRecordCreate{
		Name:    d.Get("name").(string),
		Address: d.Get("address").(string),
		Ttl:     int32(d.Get("ttl").(int)),
	}
	record, _, err := hv.client.DomainsApi.PostAaaaRecordResource(hv.auth, domainId, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /domains/%d/aaaa-record", domainId)
	}

	d.SetId(fmt.Sprint(record.Id))
	log.Printf("[INFO] Created AAAA record %s in domain %d", d.Id(), domainId)

	return resourceDnsRecordAAAARead(ctx, d, m)
}

func resourceDnsRecordAAAARead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	domainId := int32(d.Get("domain_id").(int))

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	record, resp, err := hv.client.DomainsApi.GetAaaaRecordIdResource(hv.auth, domainId, int32(recordId), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] AAAA record %d not found, removing from state", recordId)
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /domains/%d/aaaa-record/%d", domainId, recordId)
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"domain_id": int(record.DomainId),
		"name":      record.Name,
		"address":   record.Address,
		"ttl":       int(record.Ttl),
		"record_id": int(record.Id),
		"type":      record.Type_,
	}))
}

func resourceDnsRecordAAAAUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	domainId := int32(d.Get("domain_id").(int))

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload := swagger.AaaaRecordUpdate{
		Name:    d.Get("name").(string),
		Address: d.Get("address").(string),
		Ttl:     int32(d.Get("ttl").(int)),
	}
	if _, _, err := hv.client.DomainsApi.PutAaaaRecordIdResource(hv.auth, domainId, int32(recordId), payload, nil); err != nil {
		return apiErrorDiag(err, "PUT /domains/%d/aaaa-record/%d", domainId, recordId)
	}

	return resourceDnsRecordAAAARead(ctx, d, m)
}

func resourceDnsRecordAAAADelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	domainId := int32(d.Get("domain_id").(int))

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting AAAA record %d", recordId)
	if resp, err := hv.client.DomainsApi.DeleteAaaaRecordIdResource(hv.auth, domainId, int32(recordId)); err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /domains/%d/aaaa-record/%d", domainId, recordId)
	}

	d.SetId("")
	return nil
}
