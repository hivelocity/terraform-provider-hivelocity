package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func ignitionConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
	        "filter": dataSourceFiltersSchema(),
		    "first":  dataSourceFilterFirstSchema(),
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"contents": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
	}
}

func dataSourceIgnitionConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIgnitionConfigRead,
		Schema:      ignitionConfigSchema(),
	}
}

func dataSourceIgnitionConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	ignitionConfigInfo, _, err := hv.client.IgnitionApi.GetIgnitionResource(hv.auth, nil)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /ignition failed! (%s)\n\n %s", err, myErr.Body())
	}

	jsonIgnitionConfigInfo, err := json.Marshal(ignitionConfigInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	ignitionConfigs := make([]map[string]interface{}, 0)
	err = json.Unmarshal(jsonIgnitionConfigInfo, &ignitionConfigs)
	if err != nil {
		return diag.FromErr(err)
	}

	ignitionConfigs = convertKeysOfList(ignitionConfigs)

	ignitionConfig, err := doFiltering(d, ignitionConfigs)
	if err != nil {
		return diag.FromErr(err)
	}

	if ignitionConfig == nil {
		return nil
	}

	for k, v := range ignitionConfig {
		d.Set(k, v)
	}

	d.SetId(fmt.Sprint(ignitionConfig["id"]))

	return nil
}
