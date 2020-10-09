package hivelocity

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceProducts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProductRead,
		Schema: map[string]*schema.Schema{
			"filter":                         dataSourceFiltersSchema(),
			"first":                          {Type: schema.TypeBool, Optional: true, Default: true},
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
	filters, filtersOk := d.GetOk("filter")

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

	var filteredProducts []map[string]interface{}
	if filtersOk {
		f := buildFilters(filters.(*schema.Set))
		for _, product := range products {
			if matchFilters(f, product) {
				filteredProducts = append(filteredProducts, product)
			}
		}
	} else {
		filteredProducts = products
	}

	if !d.Get("first").(bool) && len(filteredProducts) != 1 {
		return diag.Errorf("found %s matches. set first = true or modify your filters", len(filteredProducts))
	}

	for k, v := range filteredProducts[0] {
		d.Set(k, v)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
