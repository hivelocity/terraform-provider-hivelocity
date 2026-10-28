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

func resourceDnsRecordPTR() *schema.Resource {
	return &schema.Resource{
		Description: "Sets the reverse DNS (PTR) name of one of your IP addresses. PTR records exist for " +
			"the lifetime of the IP assignment and cannot be created or deleted through the API, so this " +
			"resource adopts the existing record (look its ID up with the `hivelocity_dns_record_ptr` " +
			"data source) and destroying it leaves the current name in place.",
		CreateContext: resourceDnsRecordPTRCreate,
		ReadContext:   resourceDnsRecordPTRRead,
		UpdateContext: resourceDnsRecordPTRUpdate,
		DeleteContext: resourceDnsRecordPTRDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"record_id": {
				Description: "ID of the PTR record.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Host name the IP address resolves to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ttl": {
				Description: "Time to live of the record, in seconds.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"address":   computedString("IP address of the record."),
			"domain_id": computedInt("ID of the reverse zone the record belongs to."),
			"disabled":  computedBool("True when the record is disabled."),
			"type":      computedString("Record type as reported by the API."),
		},
	}
}

func (hv *Client) updatePtrRecord(d *schema.ResourceData, recordId int32) error {
	payload := swagger.PtrRecordUpdate{Name: d.Get("name").(string)}
	if v, ok := d.GetOk("ttl"); ok {
		payload.Ttl = int32(v.(int))
	}
	if _, _, err := hv.client.DomainsApi.PutPtrRecordIdResource(hv.auth, recordId, payload, nil); err != nil {
		return formatSwaggerError(err, "PUT /domains/ptr/%d", recordId)
	}
	return nil
}

func resourceDnsRecordPTRCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	recordId := int32(d.Get("record_id").(int))

	if err := hv.updatePtrRecord(d, recordId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(recordId))
	return resourceDnsRecordPTRRead(ctx, d, m)
}

func resourceDnsRecordPTRRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	record, resp, err := hv.client.DomainsApi.GetPtrRecordIdResource(hv.auth, int32(recordId), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] PTR record %d not found, removing from state", recordId)
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /domains/ptr/%d", recordId)
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"record_id": recordId,
		"name":      record.Name,
		"ttl":       int(record.Ttl),
		"address":   record.Address,
		"domain_id": int(record.DomainId),
		"disabled":  record.Disabled,
		"type":      record.Type_,
	}))
}

func resourceDnsRecordPTRUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("name", "ttl") {
		if err := hv.updatePtrRecord(d, int32(recordId)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDnsRecordPTRRead(ctx, d, m)
}

func resourceDnsRecordPTRDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// PTR records cannot be deleted through the API; the current name is left in place.
	log.Printf("[INFO] Removing PTR record %s from state without changing it", d.Id())
	d.SetId("")
	return nil
}

func dataSourceDnsRecordPTR() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up the PTR (reverse DNS) record of one of your IP addresses, typically by `address`.",
		ReadContext: dataSourceDnsRecordPTRRead,
		Schema: map[string]*schema.Schema{
			"filter":    dataSourceFiltersSchema(),
			"first":     dataSourceFilterFirstSchema(),
			"record_id": computedInt("Unique ID of the record."),
			"address":   computedString("IP address of the record."),
			"name":      computedString("Host name the address resolves to."),
			"ttl":       computedInt("Time to live, in seconds."),
			"domain_id": computedInt("ID of the reverse zone."),
			"disabled":  computedBool("True when the record is disabled."),
			"type":      computedString("Record type."),
		},
	}
}

func dataSourceDnsRecordPTRRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	records, _, err := hv.client.DomainsApi.GetPtrRecordResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /domains/ptr")
	}

	list, err := toGenericList(records)
	if err != nil {
		return diag.FromErr(err)
	}
	// Expose the record's `id` under the clearer `record_id` name.
	for _, item := range list {
		item["record_id"] = item["id"]
	}
	if len(list) == 0 {
		return diag.Errorf("no PTR records found")
	}

	item, err := doFiltering(d, list, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if item == nil {
		return diag.Errorf("no PTR record matched the given filters")
	}

	return setGenericItem(d, item, dataSourceDnsRecordPTR().Schema, "record_id")
}
