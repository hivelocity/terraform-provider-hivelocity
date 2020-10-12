package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProduct() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProductRead,
		Schema: map[string]*schema.Schema{
			"filter":                         dataSourceFiltersSchema(),
			"first":                          dataSourceFilterFirstSchema(),
			"product_id":                     {Type: schema.TypeInt, Computed: true},
			"product_name":                   {Type: schema.TypeString, Computed: true},
			"location":                       {Type: schema.TypeString, Computed: true},
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
		},
	}
}

func dataSourceProductRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	productInfo, _, err := hv.client.ProductApi.GetProductListResource(hv.auth, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonProductInfo, err := json.Marshal(productInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	products := make([]map[string]interface{}, 0)
	err = json.Unmarshal(jsonProductInfo, &products)
	if err != nil {
		return diag.FromErr(err)
	}

	products = convertKeysOfList(products)

	product, err := doFiltering(d, products)
	if err != nil {
		return diag.FromErr(err)
	}

	if product == nil {
		return nil
	}

	for k, v := range product {
		d.Set(k, v)
	}

	d.SetId(fmt.Sprint(product["product_id"]))

	return nil
}
