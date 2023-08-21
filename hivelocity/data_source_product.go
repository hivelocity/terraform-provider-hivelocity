package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

// Calling the API endpoint directly due to the dynamic types causing issues with the generated client
func dataSourceProductRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*Client)
	if !ok {
		return diag.Errorf("Failed to get configuration")
	}
	apiEndpoint := c.ApiUrl + "/inventory/product"

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("X-API-KEY", c.ApiKey)
	req.Header.Add("Referer", c.Referer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return diag.Errorf("Failed to call the API! (%s)", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("API Response:", string(body))
	fmt.Fprintf(os.Stderr, "API Response: %s\n", string(body))

	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("API returned non-200 status code: %d. Response: %s", resp.StatusCode, body)
	}

	var nestedResponse map[string][]map[string]interface{}
	errNested := json.Unmarshal(body, &nestedResponse)

	var flatResponse []map[string]interface{}
	errFlat := json.Unmarshal(body, &flatResponse)

	if errNested != nil && errFlat != nil {
		return diag.Errorf("Failed to parse API response")
	}

	var products []map[string]interface{}

	if errNested == nil {
		for _, v := range nestedResponse {
			products = append(products, v...)
		}
	} else {
		products = flatResponse
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
