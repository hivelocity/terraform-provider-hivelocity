package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceIgnitionConfig(forceNew bool) *schema.Resource {
	return &schema.Resource{
		Description: "`hivelocity_ignition_config` holds the contents of an CoreOS ignition file used for CoreOS/Flatcar" +
			" installs.",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		CreateContext: resourceIgnitionCreate,
		ReadContext:   resourceIgnitionRead,
		DeleteContext: resourceIgnitionDelete,
		UpdateContext: resourceIgnitionUpdate,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    forceNew,
				Description: "Name of ignition file resource",
			},
			"contents": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					var dest interface{}
					warnings := make([]string, 0)
					errors := make([]error, 0)
					if err := json.Unmarshal([]byte(i.(string)), &dest); err != nil {
						errors = append(errors, err)
					}
					return warnings, errors
				},
				StateFunc: func(i interface{}) string {
					normStr, err := normalizeJsonString(i.(string))
					// If an error happened then it's probably not a json string, so just return it
					if err != nil {
						return i.(string)
					}
					return normStr
				},
				Description: "String of the JSON contents of the ignition file",
			},
		},
	}
}

func resourceIgnitionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := swagger.CreateIgnition{
		Name:     d.Get("name").(string),
		Contents: d.Get("contents").(string),
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

func _resourceIgnitionReadFromResponse(d *schema.ResourceData, ignitionId int, resp swagger.IgnitionResponse) diag.Diagnostics {
	normStr, err := normalizeJsonString(resp.Contents)
	if err != nil {
		return diag.FromErr(err)
	}

	values := map[string]interface{}{
		"name":     resp.Name,
		"contents": normStr,
	}

	for k, v := range values {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceIgnitionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	d.SetId(fmt.Sprint(IgnitionConfigResponse.Id))

	return _resourceIgnitionReadFromResponse(d, IgnitionConfigID, IgnitionConfigResponse)
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

func resourceIgnitionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	IgnitionConfigID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload := swagger.UpdateIgnition{
		Contents: d.Get("contents").(string),
	}

	ignitionResponse, response, err := hv.client.IgnitionApi.PutIgnitionResourceId(hv.auth, int32(IgnitionConfigID), payload, nil)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return diag.FromErr(err)
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /ignition/%s failed! (%s)\n\n %s", fmt.Sprint(IgnitionConfigID), err, myErr.Body())
	}

	return _resourceIgnitionReadFromResponse(d, IgnitionConfigID, ignitionResponse)
}
