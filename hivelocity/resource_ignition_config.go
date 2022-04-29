package hivelocity

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceIgnitionConfig(forceNew bool) *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		CreateContext: resourceIgnitionCreate,
		ReadContext:   resourceIgnitionRead,
		DeleteContext: resourceIgnitionDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: forceNew,
			},
			"contents": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: forceNew,
			},
		},
	}
}

func resourceIgnitionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := swagger.CreateIgnition{
		Name:       d.Get("name").(string),
		Contents:   d.Get("contents").(string),
	}

	ignitionConfig, _, err := hv.client.IgnitionApi.PostIgnitionResource(hv.auth, payload, nil)
	if err != nil {
		d.SetId("")
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("POST /ignition failed! (%s)\n\n %s", err, myErr.Body())
	}

	d.SetId(fmt.Sprint(ignitionConfig.Id))

	return resourceIgnitionRead(ctx, d, m)
}

func resourceIgnitionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hv, _ := m.(*Client)

	IgnitionConfigID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	IgnitionConfigResponse, httpResponse, err := hv.client.IgnitionApi.GetIgnitionResourceId(hv.auth, int32(IgnitionConfigID), nil)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /ignition/%d failed! (%s)\n\n %s", IgnitionConfigID, err, myErr.Body())
	}

	d.Set("id", IgnitionConfigResponse.Id)
	d.Set("name", IgnitionConfigResponse.Name)
	d.Set("contents", IgnitionConfigResponse.Contents)

	return diags
}

func resourceIgnitionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hv, _ := m.(*Client)

	IgnitionConfigID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := hv.client.IgnitionApi.DeleteIgnitionResourceId(hv.auth, int32(IgnitionConfigID), nil)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /ignition/%s failed! (%s)\n\n %s", fmt.Sprint(IgnitionConfigID), err, myErr.Body())
	}

	d.SetId("")

	return diags
}
