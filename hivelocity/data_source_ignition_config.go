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

	ignitionConfigInfoList, _, err := hv.client.IgnitionApi.GetIgnitionResource(hv.auth, nil)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /ignition failed! (%s)\n\n %s", err, myErr.Body())
	}

	// Normalize json string
	for i, item := range ignitionConfigInfoList {
		normStr, err := normalizeJsonString(item.Contents)
		if err != nil {
			return diag.FromErr(err)
		}
		ignitionConfigInfoList[i].Contents = normStr
	}

	jsonIgnitionConfigInfo, err := json.Marshal(ignitionConfigInfoList)
	if err != nil {
		return diag.FromErr(err)
	}

	ignitionConfigs := make([]map[string]interface{}, 0)
	if err = json.Unmarshal(jsonIgnitionConfigInfo, &ignitionConfigs); err != nil {
		return diag.FromErr(err)
	}

	ignitionConfigs = convertKeysOfList(ignitionConfigs)

	ignitionConfig, err := doFiltering(d, ignitionConfigs, nil)
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
