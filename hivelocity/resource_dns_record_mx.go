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

func resourceDnsRecordMX() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an MX record in a Hivelocity-hosted zone. Import with `<domain_id>:<record_id>`.",
		CreateContext: resourceDnsRecordMXCreate,
		ReadContext:   resourceDnsRecordMXRead,
		UpdateContext: resourceDnsRecordMXUpdate,
		DeleteContext: resourceDnsRecordMXDelete,
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
				Description: "Record name, for example `@` for the zone apex.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"exchange": {
				Description: "Mail server host name, for example `mail.example.com`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"preference": {
				Description: "Preference (priority) of the mail server; lower is preferred.",
				Type:        schema.TypeInt,
				Required:    true,
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

func resourceDnsRecordMXCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	domainId := int32(d.Get("domain_id").(int))

	payload := swagger.MxRecordCreate{
		Name:       d.Get("name").(string),
		Exchange:   d.Get("exchange").(string),
		Preference: int32(d.Get("preference").(int)),
		Ttl:        int32(d.Get("ttl").(int)),
	}
	record, _, err := hv.client.DomainsApi.PostMxRecordResource(hv.auth, domainId, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /domains/%d/mx-record", domainId)
	}

	d.SetId(fmt.Sprint(record.Id))
	log.Printf("[INFO] Created MX record %s in domain %d", d.Id(), domainId)

	return resourceDnsRecordMXRead(ctx, d, m)
}

func resourceDnsRecordMXRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	domainId := int32(d.Get("domain_id").(int))

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	record, resp, err := hv.client.DomainsApi.GetMxRecordIdResource(hv.auth, domainId, int32(recordId), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] MX record %d not found, removing from state", recordId)
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /domains/%d/mx-record/%d", domainId, recordId)
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"domain_id":  int(record.DomainId),
		"name":       record.Name,
		"exchange":   record.Exchange,
		"preference": int(record.Preference),
		"ttl":        int(record.Ttl),
		"record_id":  int(record.Id),
		"type":       record.Type_,
	}))
}

func resourceDnsRecordMXUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	domainId := int32(d.Get("domain_id").(int))

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload := swagger.MxRecordUpdate{
		Name:       d.Get("name").(string),
		Exchange:   d.Get("exchange").(string),
		Preference: int32(d.Get("preference").(int)),
		Ttl:        int32(d.Get("ttl").(int)),
	}
	if _, _, err := hv.client.DomainsApi.PutMxRecordIdResource(hv.auth, domainId, int32(recordId), payload, nil); err != nil {
		return apiErrorDiag(err, "PUT /domains/%d/mx-record/%d", domainId, recordId)
	}

	return resourceDnsRecordMXRead(ctx, d, m)
}

func resourceDnsRecordMXDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	domainId := int32(d.Get("domain_id").(int))

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting MX record %d", recordId)
	if resp, err := hv.client.DomainsApi.DeleteMxRecordIdResource(hv.auth, domainId, int32(recordId)); err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /domains/%d/mx-record/%d", domainId, recordId)
	}

	d.SetId("")
	return nil
}
