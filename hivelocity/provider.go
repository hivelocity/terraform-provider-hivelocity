package hivelocity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a schema.Provider for Hivelocity.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: descriptions["api_key"],
				DefaultFunc: schema.EnvDefaultFunc("HIVELOCITY_API_KEY", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hivelocity_devices":                   dataSourceDevices(),
			"hivelocity_products":                  dataSourceProducts(),
			"hivelocity_product_options":           dataSourceProductOption(),
			"hivelocity_product_operating_systems": dataSourceProductOperatingSystem(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"hivelocity_bare_metal_devices": resourceDevice(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_key": "Your API Key from the https://my.hivelocity.net portal.",
	}

}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	config := Config{
		ApiKey: apiKey,
	}

	client, err := config.Client()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create the Hivelocity client",
			Detail:   "Unable to create the Hivelocity client",
		})

		return nil, diags
	}

	return client, diags
}
