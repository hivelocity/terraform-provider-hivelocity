package hivelocity

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceDnsRecordAAAA() *schema.Resource {
	return &schema.Resource{
		Description: "Resource used to manage DNS AAAA Records",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceDnsRecordAAAACreate,
		ReadContext:   resourceDnsRecordAAAARead,
		UpdateContext: resourceDnsRecordAAAAUpdate,
		DeleteContext: resourceDnsRecordAAAADelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDnsRecordAAAAImport,
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
			"address": {
				Description: "The IP address for this record",
				Type:        schema.TypeString,
				Required:    true,
				ValidateDiagFunc: func(val interface{}, cty cty.Path) diag.Diagnostics {
					ip := net.ParseIP(val.(string))
					if ip == nil || ip.To16() == nil {
						return diag.Errorf("Invalid IP address: '%s'", val.(string))
					}
					return nil
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return net.ParseIP(old).Equal(net.ParseIP(new))
				},
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

func resourceDnsRecordAAAACreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	domainId := d.Get("domain_id").(int)

	payload := swagger.AaaaRecordCreate{
		Name:    d.Get("name").(string),
		Address: d.Get("address").(string),
		Ttl:     int32(d.Get("ttl").(int)),
	}

	record, _, err := hv.client.DomainsApi.PostAaaaRecordResource(hv.auth, int32(domainId), payload, nil)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("POST /domains/%d/aaaa-record failed! (%s)\n\n %s", domainId, err, myErr.Body())
	}

	d.SetId(fmt.Sprint(record.Id))

	log.Printf("[INFO] Created AAAA Record ID: %s", d.Id())

	return resourceDnsRecordAAAARead(ctx, d, m)
}

func resourceDnsRecordAAAARead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	record, response, err := hv.client.DomainsApi.GetAaaaRecordIdResource(hv.auth, int32(domainId), int32(recordId), nil)
	if err != nil {
		// If resource was deleted outside terraform, remove it from state and exit gracefully
		if response.StatusCode == 404 {
			log.Printf("[WARN] AAAA Record ID not found: (%s)", d.Id())
			d.SetId("")
			return nil
		}
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /domains/%d/aaaa-record/%d failed! (%s)\n\n %s", domainId, recordId, err, myErr.Body())
	}

	d.Set("name", record.Name)
	d.Set("address", record.Address)
	d.Set("ttl", int(record.Ttl))
	d.Set("domain", domain.Name)

	return nil
}

func resourceDnsRecordAAAAUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Updating AAAA Record ID: %s", d.Id())

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	domainId := d.Get("domain_id").(int)

	payload := swagger.AaaaRecordUpdate{
		Name:    d.Get("name").(string),
		Address: d.Get("address").(string),
		Ttl:     int32(d.Get("ttl").(int)),
	}
	record, _, err := hv.client.DomainsApi.PutAaaaRecordIdResource(hv.auth, int32(domainId), int32(recordId), payload, nil)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /domains/%d/aaaa-record/%d failed! (%s)\n\n %s", domainId, recordId, err, myErr.Body())
	}
	d.SetId(fmt.Sprint(record.Id))

	return resourceDnsRecordAAAARead(ctx, d, m)
}

func resourceDnsRecordAAAADelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Deleting AAAA Record ID: %s", d.Id())

	recordId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	domainId := d.Get("domain_id").(int)

	response, err := hv.client.DomainsApi.DeleteAaaaRecordIdResource(hv.auth, int32(domainId), int32(recordId))
	if err != nil {
		// If resource was deleted outside terraform, remove it from state and exit gracefully
		if response.StatusCode == 404 {
			log.Printf("[WARN] AAAA Record ID not found: (%s)", d.Id())
			d.SetId("")
			return nil
		}
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /domains/%d/aaaa-record/%d failed! (%s)\n\n %s", domainId, recordId, err, myErr.Body())
	}

	// Delete resource from state
	d.SetId("")

	return nil
}

func resourceDnsRecordAAAAImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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

	if _, _, err := hv.client.DomainsApi.GetAaaaRecordIdResource(hv.auth, int32(domainId), int32(recordId), nil); err != nil {
		return nil, fmt.Errorf("could not import record: %v", err)
	}

	d.SetId(fmt.Sprint(recordId))
	d.Set("domain_id", domainId)

	return []*schema.ResourceData{d}, nil
}
