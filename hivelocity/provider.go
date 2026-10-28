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
				Description: "Your API key from the https://my.hivelocity.net portal. Can also be set with the `HIVELOCITY_API_KEY` environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("HIVELOCITY_API_KEY", nil),
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The API endpoint to communicate with. Defaults to `https://core.hivelocity.net/api/v2`. Can also be set with the `HIVELOCITY_API_URL` environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("HIVELOCITY_API_URL", "https://core.hivelocity.net/api/v2"),
			},
			"referer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Value sent in the `Referer` header so API calls can be attributed. Defaults to `Terraform`. Can also be set with the `HIVELOCITY_REFERER` environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("HIVELOCITY_REFERER", "Terraform"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hivelocity_bare_metal_device":    dataSourceBareMetalDevice(),
			"hivelocity_device":               dataSourceDevice(),
			"hivelocity_product":              dataSourceProduct(),
			"hivelocity_device_initial_creds": dataSourceDeviceInitialCreds(),
			"hivelocity_ssh_key":              dataSourceSshKey(),
			"hivelocity_effective_ignition":   dataSourceEffectiveIgnition(),
			"hivelocity_ignition":             dataSourceIgnitionConfig(),
			"hivelocity_device_port":          dataSourceDevicePort(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"hivelocity_bare_metal_device": resourceBareMetalDevice(true),
			"hivelocity_ssh_key":           resourceSSHKey(),
			"hivelocity_vlan":              resourceVlan(),
			"hivelocity_order_group":       resourceOrderGroup(),
			"hivelocity_ignition":          resourceIgnitionConfig(true),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		ApiKey:  d.Get("api_key").(string),
		ApiUrl:  d.Get("api_url").(string),
		Referer: d.Get("referer").(string),
	}

	client, err := config.Client()
	if err != nil {
		return nil, diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Unable to create the Hivelocity client",
			Detail:   err.Error(),
		}}
	}

	return client, nil
}
