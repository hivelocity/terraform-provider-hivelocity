package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func dataSourceProduct() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProductRead,
		Schema: map[string]*schema.Schema{
			"filter":                         dataSourceFiltersSchema(),
			"first":                          dataSourceFilterFirstSchema(),
			"product_id":                     {Type: schema.TypeInt, Computed: true},
			"product_name":                   {Type: schema.TypeString, Computed: true},
			"data_center":                    {Type: schema.TypeString, Computed: true},
			"core":                           {Type: schema.TypeBool, Computed: true},
			"edge":                           {Type: schema.TypeBool, Computed: true},
			"stock":                          {Type: schema.TypeString, Computed: true},
			"product_on_sale":                {Type: schema.TypeBool, Computed: true},
			"product_hourly_price":           {Type: schema.TypeFloat, Computed: true},
			"product_monthly_price":          {Type: schema.TypeFloat, Computed: true},
			"product_quarterly_price":        {Type: schema.TypeFloat, Computed: true},
			"product_semi_annually_price":    {Type: schema.TypeFloat, Computed: true},
			"product_annually_price":         {Type: schema.TypeFloat, Computed: true},
			"product_biennial_price":         {Type: schema.TypeFloat, Computed: true},
			"product_triennial_price":        {Type: schema.TypeFloat, Computed: true},
			"product_display_price":          {Type: schema.TypeFloat, Computed: true},
			"product_original_price":         {Type: schema.TypeFloat, Computed: true},
			"hourly_location_premium":        {Type: schema.TypeFloat, Computed: true},
			"monthly_location_premium":       {Type: schema.TypeFloat, Computed: true},
			"quarterly_location_premium":     {Type: schema.TypeFloat, Computed: true},
			"semi_annually_location_premium": {Type: schema.TypeFloat, Computed: true},
			"annually_location_premium":      {Type: schema.TypeFloat, Computed: true},
			"biennial_location_premium":      {Type: schema.TypeFloat, Computed: true},
			"triennial_location_premium":     {Type: schema.TypeFloat, Computed: true},
			"product_bandwidth":              {Type: schema.TypeString, Computed: true},
			"product_cpu":                    {Type: schema.TypeString, Computed: true},
			"product_cpu_cores":              {Type: schema.TypeString, Computed: true},
			"product_drive":                  {Type: schema.TypeString, Computed: true},
			"product_memory":                 {Type: schema.TypeString, Computed: true},
			"product_gpu":                    {Type: schema.TypeString, Computed: true},
			"product_disabled_billing_periods": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

/* TODO FIX - CURRENTLY BROKEN */
func dataSourceProductRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	inventory, _, err := hv.client.InventoryApi.GetStockResource(hv.auth, nil)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /inventory/product failed! (%s)\n\n %s", err, myErr.Body())
	}

	jsonProductInfo, err := json.Marshal(inventory)
	if err != nil {
		return diag.FromErr(err)
	}

	products := make([]map[string]interface{}, 0)
	err = json.Unmarshal(jsonProductInfo, &products)
	if err != nil {
		return diag.FromErr(err)
	}

	products = convertKeysOfList(products)

	product, err := doFiltering(d, products, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if product == nil {
		var diags diag.Diagnostics
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Product not found.",
			Detail:   "No products matched the search.",
		})

		return diags
	}

	for k, v := range product {
		d.Set(k, v)
	}

	d.SetId(fmt.Sprint(product["product_id"]))
	return nil
}
