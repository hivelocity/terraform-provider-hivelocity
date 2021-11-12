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
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   false,
				Description: descriptions["api_url"],
				DefaultFunc: schema.EnvDefaultFunc("HIVELOCITY_API_URL", "https://core.hivelocity.net/api/v2"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hivelocity_bare_metal_device":    dataSourceBareMetalDevice(),
			"hivelocity_device":               dataSourceDevice(),
			"hivelocity_product":              dataSourceProduct(),
			"hivelocity_device_initial_creds": dataSourceDeviceInitialCreds(),
			"hivelocity_ssh_key":              dataSourceSshKey(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"hivelocity_bare_metal_device": resourceBareMetalDevice(true),
			"hivelocity_ssh_key":           resourceSSHKey(),
			"hivelocity_vlan":              resourceVlan(),
			"hivelocity_order_group":       resourceOrderGroup(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_key": "Your API Key from the https://my.hivelocity.net portal.",
		"api_url": "The API instance to communicate with defaults to https://core.hivelocity.net/api/v2",
	}

}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	apiURL := d.Get("api_url").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	config := Config{
		ApiKey: apiKey,
		ApiUrl: apiURL,
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
