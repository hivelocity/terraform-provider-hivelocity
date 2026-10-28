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
		Description: "Look up a product (server plan) and its per-location stock and pricing from the public " +
			"inventory. Filter on any attribute, for example `product_memory`, `data_center` or `stock`.",
		ReadContext: dataSourceProductRead,
		Schema: map[string]*schema.Schema{
			"filter":                         dataSourceFiltersSchema(),
			"first":                          dataSourceFilterFirstSchema(),
			"product_id":                     {Type: schema.TypeInt, Computed: true, Description: "The unique ID of this product."},
			"product_name":                   {Type: schema.TypeString, Computed: true, Description: "The unique name of this product."},
			"data_center":                    {Type: schema.TypeString, Computed: true, Description: "A facility code. For example `NYC1`."},
			"core":                           {Type: schema.TypeBool, Computed: true, Description: "True if core site."},
			"edge":                           {Type: schema.TypeBool, Computed: true, Description: "True if edge site."},
			"stock":                          {Type: schema.TypeString, Computed: true, Description: "available|limited|unavailable."},
			"is_vps":                         {Type: schema.TypeBool, Computed: true, Description: "True if this is a VPS product."},
			"is_vds":                         {Type: schema.TypeBool, Computed: true, Description: "True if this VPS product is VDS type."},
			"is_sps":                         {Type: schema.TypeBool, Computed: true, Description: "True if this is an SPS (instant) product."},
			"is_disk_upgradable":             {Type: schema.TypeBool, Computed: true, Description: "True if the product has multiple disk upgrade options."},
			"product_new":                    {Type: schema.TypeBool, Computed: true, Description: "True if the product is newly released."},
			"location_option_id":             {Type: schema.TypeInt, Computed: true, Description: "Location option ID for the data center."},
			"oos_ordering_available":         {Type: schema.TypeBool, Computed: true, Description: "True if out-of-stock ordering is available at this location."},
			"oos_lead_time_days":             {Type: schema.TypeInt, Computed: true, Description: "Lead time in days for out-of-stock ordering."},
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
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Orders for these billing periods will fail.",
			},
		},
	}
}

// The inventory endpoint is called directly: its swagger definition (`Inventory`)
// has no properties, so the generated client decodes it into an empty struct.
// The response is either a flat list of products or a map of location -> list.
func dataSourceProductRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, ok := m.(*Client)
	if !ok {
		return diag.Errorf("Failed to get configuration")
	}

	body, _, err := hv.rawRequest(ctx, "GET", "/inventory/product", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	products, err := parseInventoryProducts(body)
	if err != nil {
		return diag.FromErr(err)
	}

	products = convertKeysOfList(products)
	product, err := doFiltering(d, products, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if product == nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Product not found.",
			Detail:   "No products matched the search.",
		}}
	}

	for k, v := range filterNonSchemaKeys(product, dataSourceProduct().Schema) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprint(product["product_id"]))
	return nil
}

// parseInventoryProducts accepts both response shapes the inventory endpoint
// has used: a flat JSON array of products, or an object keyed by location whose
// values are arrays of products.
func parseInventoryProducts(body []byte) ([]map[string]interface{}, error) {
	var flatResponse []map[string]interface{}
	if err := json.Unmarshal(body, &flatResponse); err == nil {
		return flatResponse, nil
	}

	var nestedResponse map[string][]map[string]interface{}
	if err := json.Unmarshal(body, &nestedResponse); err == nil {
		var products []map[string]interface{}
		for _, v := range nestedResponse {
			products = append(products, v...)
		}
		return products, nil
	}

	return nil, fmt.Errorf("failed to parse inventory response: %s", truncate(string(body), 200))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
