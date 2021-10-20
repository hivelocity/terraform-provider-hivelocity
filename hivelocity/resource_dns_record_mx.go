package hivelocity

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceDnsRecordMX() *schema.Resource {
	return &schema.Resource{
		Description: "Resource used to manage DNS MX Records",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceDnsRecordMXCreate,
		ReadContext:   resourceDnsRecordMXRead,
		UpdateContext: resourceDnsRecordMXUpdate,
		DeleteContext: resourceDnsRecordMXDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDnsRecordMXImport,
		},
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Description: "ID of DNS Domain associated with this record",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The domain name for this record",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					domain := d.Get("domain").(string)
					if domain == "" {
						return false
					}
					return (old == "@" && new == domain) || (new == "@" && old == domain)
				},
			},
			"exchange": {
				Description: "The mail exchange server name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"preference": {
				Description: "Preference of this record",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"ttl": {
				Description: "The time to live for this record",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"domain": {
				Description: "The domain name (zone) associated with this record",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceDnsRecordMXCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	domainId := d.Get("domain_id").(int)

	payload := swagger.MxRecordCreate{
		Name:       d.Get("name").(string),
		Exchange:   d.Get("exchange").(string),
		Preference: int32(d.Get("preference").(int)),
		Ttl:        int32(d.Get("ttl").(int)),
	}

	record, _, err := hv.client.DomainsApi.PostMxRecordResource(hv.auth, int32(domainId), payload, nil)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("POST /domains/%d/mx-record failed! (%s)\n\n %s", domainId, err, myErr.Body())
	}

	d.SetId(fmt.Sprint(record.Id))

	log.Printf("[INFO] Created MX Record ID: %s", d.Id())

	return resourceDnsRecordMXRead(ctx, d, m)
}

func resourceDnsRecordMXRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	domainId := d.Get("domain_id").(int)

	domain, response, err := hv.client.DomainsApi.GetDomainIdResource(hv.auth, int32(domainId), nil)
	if err != nil {
		// If resource was deleted outside terraform, remove it from state and exit gracefully
		if response.StatusCode == 404 {
			log.Printf("[WARN] Domain ID not found: (%s)", d.Id())
			d.SetId("")
			return nil
		}
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /domains/%d failed! (%s)\n\n %s", domainId, err, myErr.Body())
	}

	record, response, err := hv.client.DomainsApi.GetMxRecordIdResource(hv.auth, int32(domainId), int32(recordId), nil)
	if err != nil {
		// If resource was deleted outside terraform, remove it from state and exit gracefully
		if response.StatusCode == 404 {
			log.Printf("[WARN] MX Record ID not found: (%s)", d.Id())
			d.SetId("")
			return nil
		}
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /domains/%d/mx-record/%d failed! (%s)\n\n %s", domainId, recordId, err, myErr.Body())
	}

	d.Set("name", record.Name)
	d.Set("preference", int(record.Preference))
	d.Set("exchange", record.Exchange)
	d.Set("ttl", int(record.Ttl))
	d.Set("domain", domain.Name)

	return nil
}

func resourceDnsRecordMXUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Updating MX Record ID: %s", d.Id())

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	domainId := d.Get("domain_id").(int)

	payload := swagger.MxRecordUpdate{
		Name:       d.Get("name").(string),
		Preference: int32(d.Get("preference").(int)),
		Exchange:   d.Get("exchange").(string),
		Ttl:        int32(d.Get("ttl").(int)),
	}
	record, _, err := hv.client.DomainsApi.PutMxRecordIdResource(hv.auth, int32(domainId), int32(recordId), payload, nil)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /domains/%d/mx-record/%d failed! (%s)\n\n %s", domainId, recordId, err, myErr.Body())
	}
	d.SetId(fmt.Sprint(record.Id))

	return resourceDnsRecordMXRead(ctx, d, m)
}

func resourceDnsRecordMXDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Deleting MX Record ID: %s", d.Id())

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	domainId := d.Get("domain_id").(int)

	response, err := hv.client.DomainsApi.DeleteMxRecordIdResource(hv.auth, int32(domainId), int32(recordId))
	if err != nil {
		// If resource was deleted outside terraform, remove it from state and exit gracefully
		if response.StatusCode == 404 {
			log.Printf("[WARN] MX Record ID not found: (%s)", d.Id())
			d.SetId("")
			return nil
		}
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /domains/%d/mx-record/%d failed! (%s)\n\n %s", domainId, recordId, err, myErr.Body())
	}

	// Delete resource from state
	d.SetId("")

	return nil
}

func resourceDnsRecordMXImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	hv, _ := m.(*Client)
	importId := d.Id()
	parts := strings.SplitN(importId, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected domainId:recordId", importId)
	}

	domainId, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	recordId, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	if _, _, err := hv.client.DomainsApi.GetMxRecordIdResource(hv.auth, int32(domainId), int32(recordId), nil); err != nil {
		return nil, fmt.Errorf("could not import record: %v", err)
	}

	d.SetId(fmt.Sprint(recordId))
	d.Set("domain_id", domainId)

	return []*schema.ResourceData{d}, nil
}
