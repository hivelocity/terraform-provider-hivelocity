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

func resourceDnsDomain() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a DNS zone hosted on Hivelocity's name servers. Records are managed with the " +
			"`hivelocity_dns_record_*` resources.",
		CreateContext: resourceDnsDomainCreate,
		ReadContext:   resourceDnsDomainRead,
		DeleteContext: resourceDnsDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Domain name of the zone, for example `example.com`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"domain_id":  computedInt("Unique ID of the domain."),
			"directs_to": computedString("Where the domain directs to, as reported by the API."),
		},
	}
}

func resourceDnsDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := swagger.DomainCreate{Name: d.Get("name").(string)}
	domain, _, err := hv.client.DomainsApi.PostDomainResource(hv.auth, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /domains")
	}

	log.Printf("[INFO] Created domain ID: %d", domain.DomainId)
	d.SetId(fmt.Sprint(domain.DomainId))

	return resourceDnsDomainRead(ctx, d, m)
}

func resourceDnsDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	domainId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	domain, resp, err := hv.client.DomainsApi.GetDomainIdResource(hv.auth, int32(domainId), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] Domain %d not found, removing from state", domainId)
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /domains/%d", domainId)
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"name":       domain.Name,
		"domain_id":  int(domain.DomainId),
		"directs_to": domain.DirectsTo,
	}))
}

func resourceDnsDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	domainId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting domain ID: %d", domainId)
	if resp, err := hv.client.DomainsApi.DeleteDomainIdResource(hv.auth, int32(domainId)); err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /domains/%d", domainId)
	}

	d.SetId("")
	return nil
}

func dataSourceDnsDomain() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up one of your DNS zones, typically by `name`.",
		ReadContext: dataSourceDnsDomainRead,
		Schema: map[string]*schema.Schema{
			"filter":     dataSourceFiltersSchema(),
			"first":      dataSourceFilterFirstSchema(),
			"domain_id":  computedInt("Unique ID of the domain."),
			"name":       computedString("Domain name."),
			"directs_to": computedString("Where the domain directs to, as reported by the API."),
		},
	}
}

func dataSourceDnsDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	domains, _, err := hv.client.DomainsApi.GetDomainResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /domains")
	}

	return readFilteredList(d, domains, dataSourceDnsDomain().Schema, "domain_id", nil)
}
