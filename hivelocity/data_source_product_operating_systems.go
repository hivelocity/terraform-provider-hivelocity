package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func buildProductOperatingSystemSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": &schema.Schema{
			Type: schema.TypeInt,
			Computed: true,
		},
		"name": &schema.Schema{
			Type: schema.TypeString,
			Computed: true,
		},
		"monthly_price": &schema.Schema{
			Type: schema.TypeFloat,
			Computed: true,
		},
		"currency": &schema.Schema{
			Type: schema.TypeString,
			Computed: true,
		},
		"tags": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"expressions": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func dataSourceProductOperatingSystem() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProductOperatingSystemRead,
		Schema: map[string]*schema.Schema{
			"product_id": &schema.Schema{
				Type: schema.TypeInt,
				Required: true,
			},
			"product_operating_systems": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: buildProductOperatingSystemSchema(),
				},
			},
		},
	}
}

func dataSourceProductOperatingSystemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hv, _ := m.(*Client)

	initialProductId := d.Get("product_id").(int)
	productId := int32(initialProductId)

	hivelocityProductOperatingSystems, _, err := hv.client.ProductApi.GetProductOperatingSystemsResource(hv.auth, productId, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonProductOperatingSystems, err := json.Marshal(hivelocityProductOperatingSystems)
	if err != nil {
		return diag.FromErr(err)
	}

	productOperatingSystems := make([]map[string]interface{}, 0)
	err = json.Unmarshal(jsonProductOperatingSystems, &productOperatingSystems)
	if err != nil {
		return diag.FromErr(err)
	}

	productOperatingSystems = convertKeysOfList(productOperatingSystems)

	if err := d.Set("product_operating_systems", productOperatingSystems); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(productId))

	return diags
}
