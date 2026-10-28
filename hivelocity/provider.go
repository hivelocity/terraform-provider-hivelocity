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
			// Devices and provisioning
			"hivelocity_bare_metal_device":         dataSourceBareMetalDevice(),
			"hivelocity_device":                    dataSourceDevice(),
			"hivelocity_device_initial_creds":      dataSourceDeviceInitialCreds(),
			"hivelocity_device_hardware":           dataSourceDeviceHardware(),
			"hivelocity_device_ipmi":               dataSourceDeviceIpmi(),
			"hivelocity_device_tags":               dataSourceDeviceTags(),
			"hivelocity_effective_ignition":        dataSourceEffectiveIgnition(),
			"hivelocity_ignition":                  dataSourceIgnitionConfig(),
			"hivelocity_ssh_key":                   dataSourceSshKey(),
			"hivelocity_product":                   dataSourceProduct(),
			"hivelocity_product_operating_systems": dataSourceProductOperatingSystems(),
			"hivelocity_product_options":           dataSourceProductOptions(),
			// Networking
			"hivelocity_device_port":   dataSourceDevicePort(),
			"hivelocity_device_ips":    dataSourceDeviceIps(),
			"hivelocity_ip_assignment": dataSourceIpAssignment(),
			"hivelocity_vlan":          dataSourceVlan(),
			// DNS
			"hivelocity_dns_domain":     dataSourceDnsDomain(),
			"hivelocity_dns_record_ptr": dataSourceDnsRecordPTR(),
			// VPS
			"hivelocity_vps":          dataSourceVps(),
			"hivelocity_vps_events":   dataSourceVpsEvents(),
			"hivelocity_vps_volume":   dataSourceVpsVolume(),
			"hivelocity_vps_snapshot": dataSourceVpsSnapshot(),
			"hivelocity_vps_iso":      dataSourceVpsIso(),
			"hivelocity_vps_network":  dataSourceVpsNetwork(),
			// Account
			"hivelocity_contact":            dataSourceContact(),
			"hivelocity_permissions":        dataSourcePermissions(),
			"hivelocity_user_permissions":   dataSourceUserPermissions(),
			"hivelocity_webhook_events":     dataSourceWebhookEvents(),
			"hivelocity_service":            dataSourceService(),
			"hivelocity_services":           dataSourceServices(),
			"hivelocity_service_types":      dataSourceServiceTypes(),
			"hivelocity_order":              dataSourceOrder(),
			"hivelocity_orders":             dataSourceOrders(),
			"hivelocity_maintenance_events": dataSourceMaintenanceEvents(),
		},
		ResourcesMap: map[string]*schema.Resource{
			// Devices and provisioning
			"hivelocity_bare_metal_device": resourceBareMetalDevice(true),
			"hivelocity_order_group":       resourceOrderGroup(),
			"hivelocity_device_tags":       resourceDeviceTags(),
			"hivelocity_ssh_key":           resourceSSHKey(),
			"hivelocity_ignition":          resourceIgnitionConfig(true),
			// Networking
			"hivelocity_vlan":                resourceVlan(),
			"hivelocity_device_port_bond":    resourceDevicePortBond(),
			"hivelocity_device_port":         resourceDevicePort(),
			"hivelocity_ip_assignment_route": resourceIpAssignmentRoute(),
			"hivelocity_null_route":          resourceNullRoute(),
			// DNS
			"hivelocity_dns_domain":      resourceDnsDomain(),
			"hivelocity_dns_record_a":    resourceDnsRecordA(),
			"hivelocity_dns_record_aaaa": resourceDnsRecordAAAA(),
			"hivelocity_dns_record_mx":   resourceDnsRecordMX(),
			"hivelocity_dns_record_ptr":  resourceDnsRecordPTR(),
			// VPS
			"hivelocity_vps":                    resourceVps(),
			"hivelocity_vps_volume":             resourceVpsVolume(),
			"hivelocity_vps_snapshot":           resourceVpsSnapshot(),
			"hivelocity_vps_snapshot_schedule":  resourceVpsSnapshotSchedule(),
			"hivelocity_vps_iso":                resourceVpsIso(),
			"hivelocity_vps_network":            resourceVpsNetwork(),
			"hivelocity_vps_network_attachment": resourceVpsNetworkAttachment(),
			// Account
			"hivelocity_webhook":             resourceWebhook(),
			"hivelocity_contact":             resourceContact(),
			"hivelocity_contact_permissions": resourceContactPermissions(),
			"hivelocity_api_token":           resourceApiToken(),
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
