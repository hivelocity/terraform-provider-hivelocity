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

func resourceDnsDomain() *schema.Resource {
	return &schema.Resource{
		Description: "Resource used to manage DNS Domains",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceDnsDomainCreate,
		ReadContext:   resourceDnsDomainRead,
		DeleteContext: resourceDnsDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Domain name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceDnsDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := swagger.DomainCreate{Name: d.Get("name").(string)}

	domain, _, err := hv.client.DomainsApi.PostDomainResource(hv.auth, payload, nil)
	if err != nil {
		d.SetId("")
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("POST /domains failed! (%s)\n\n %s", err, myErr.Body())
	}

	log.Printf("[INFO] Created Domain ID: %d", domain.DomainId)

	d.SetId(fmt.Sprint(domain.DomainId))

	return resourceDnsDomainRead(ctx, d, m)
}

func resourceDnsDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	domainId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

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

	d.Set("name", domain.Name)

	return nil
}

func resourceDnsDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Deleting Domain ID: %s", d.Id())

	domainId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := hv.client.DomainsApi.DeleteDomainIdResource(hv.auth, int32(domainId))
	if err != nil {
		// If resource was deleted outside terraform, remove it from state and exit gracefully
		if response.StatusCode == 404 {
			log.Printf("[WARN] Domain ID not found: (%s)", d.Id())
			d.SetId("")
			return nil
		}
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /domains/%d failed! (%s)\n\n %s", domainId, err, myErr.Body())
	}

	// Delete resource from state
	d.SetId("")

	return nil
}
