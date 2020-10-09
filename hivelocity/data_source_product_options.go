package hivelocity

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func buildProductOptionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"product_id": &schema.Schema{
			Type: schema.TypeInt,
			Required: true,
		},
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

func dataSourceProductOption() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProductOptionRead,
		Schema: map[string]*schema.Schema{
			"product_options": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: buildProductOptionSchema(),
				},
			},
		},
	}
}

func dataSourceProductOptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diag.Errorf("Product Options is not yet supported.")

	hv, _ := m.(*Client)

	initialProductId := d.Get("product_id").(int)
	productId := int32(initialProductId)

	hivelocityProductOptions, _, err := hv.client.ProductApi.GetProductOptionResource(hv.auth, productId, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonProductOptions, err := json.Marshal(hivelocityProductOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	productOptions := make([]map[string]interface{}, 0)
	err = json.Unmarshal(jsonProductOptions, &productOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("product_options", productOptions); err != nil {
		return diag.FromErr(err)
	}

	productOptions = convertKeysOfList(productOptions)

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
