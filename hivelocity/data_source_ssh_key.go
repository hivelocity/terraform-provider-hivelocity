package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func sshKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter":     dataSourceFiltersSchema(),
		"first":      dataSourceFilterFirstSchema(),
		"ssh_key_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"public_key": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}


func dataSourceSshKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSshKeyRead,
		Schema:      sshKeySchema(),
	}
}

func dataSourceSshKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	sshKeyInfo, _, err := hv.client.SshKeyApi.GetSshKeyResource(hv.auth, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonSshKeyInfo, err := json.Marshal(sshKeyInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	sshKeys := make([]map[string]interface{}, 0)
	err = json.Unmarshal(jsonSshKeyInfo, &sshKeys)
	if err != nil {
		return diag.FromErr(err)
	}

	sshKeys = convertKeysOfList(sshKeys)

	sshKey, err := doFiltering(d, sshKeys)
	if err != nil {
		return diag.FromErr(err)
	}

	if sshKey == nil {
		return nil
	}

	for k, v := range sshKey {
		d.Set(k, v)	
	}

	d.SetId(fmt.Sprint(sshKey["ssh_key_id"]))

	return nil
}
