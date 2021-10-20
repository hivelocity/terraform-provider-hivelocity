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

func resourceDnsRecordPTR() *schema.Resource {
	return &schema.Resource{
		Description: "Resource used to manage DNS PTR Records",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceDnsRecordPtrCreate,
		ReadContext:   resourceDnsRecordPtrRead,
		UpdateContext: resourceDnsRecordPtrUpdate,
		DeleteContext: resourceDnsRecordPtrDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Description: "The IP address for this record",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The domain name for this record",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ttl": {
				Description: "The time to live for this record",
				Type:        schema.TypeInt,
				Required:    true,
			},
		},
	}
}

func resourceDnsRecordPtrCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	records, _, err := hv.client.DomainsApi.GetPtrRecordResource(hv.auth, nil)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /domains/ptr failed! (%s)\n\n %s", err, myErr.Body())
	}

	recordFound := false
	address := d.Get("address").(string)

	for _, record := range records {
		if record.Address == address {
			payload := swagger.PtrRecordUpdate{
				Name: d.Get("name").(string),
				Ttl:  int32(d.Get("ttl").(int)),
			}
			ptrRecord, _, err := hv.client.DomainsApi.PutPtrRecordIdResource(hv.auth, record.Id, payload, nil)
			if err != nil {
				myErr := err.(swagger.GenericSwaggerError)
				return diag.Errorf("PUT /domains/ptr/%d failed! (%s)\n\n %s", record.Id, err, myErr.Body())
			}

			d.SetId(fmt.Sprint(ptrRecord.Id))

			recordFound = true
			break
		}
	}

	if !recordFound {
		return diag.Errorf("Could not find PTR record for address '%s'", address)
	}

	log.Printf("[INFO] Created PTR Record ID: %s", d.Id())

	return resourceDnsRecordPtrRead(ctx, d, m)
}

func resourceDnsRecordPtrRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	record, response, err := hv.client.DomainsApi.GetPtrRecordIdResource(hv.auth, int32(recordId), nil)
	if err != nil {
		// If resource was deleted outside terraform, remove it from state and exit gracefully
		if response.StatusCode == 404 {
			log.Printf("[WARN] PTR Record ID not found: (%s)", d.Id())
			d.SetId("")
			return nil
		}
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /domains/ptr/%d failed! (%s)\n\n %s", recordId, err, myErr.Body())
	}

	d.Set("name", record.Name)
	d.Set("address", record.Address)
	d.Set("ttl", int(record.Ttl))

	return nil
}

func resourceDnsRecordPtrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Updating PTR Record ID: %s", d.Id())

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload := swagger.PtrRecordUpdate{
		Name: d.Get("name").(string),
		Ttl:  int32(d.Get("ttl").(int)),
	}
	record, _, err := hv.client.DomainsApi.PutPtrRecordIdResource(hv.auth, int32(recordId), payload, nil)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /domains/ptr/%d failed! (%s)\n\n %s", record.Id, err, myErr.Body())
	}

	d.SetId(fmt.Sprint(record.Id))

	return resourceDnsRecordPtrRead(ctx, d, m)
}

func resourceDnsRecordPtrDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting PTR Record ID: %s", d.Id())

	// Delete resource from state
	d.SetId("")

	return nil
}
