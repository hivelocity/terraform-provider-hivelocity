package hivelocity

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func dataSourceFilterFirstSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}
}

type filter struct {
	name   string
	values []string
}

func buildFilters(set *schema.Set) []filter {
	var filters []filter
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var values []string
		for _, value := range m["values"].([]interface{}) {
			values = append(values, value.(string))
		}
		filters = append(filters, filter{
			name:   m["name"].(string),
			values: values,
		})
	}
	return filters
}

func filterArrayIntersection(filter filter, arrayV []interface{}) bool {
	for _, filterValue := range filter.values {
		for _, itemValue := range arrayV {
			if filterValue == itemValue {
				return true
			}
		}
	}
	return false
}

func matches(filter filter, m map[string]interface{}) bool {
	v, ok := m[filter.name]
	if !ok {
		return false
	}

	arrayV, fail := v.([]interface{})
	if !fail {
		return filterArrayIntersection(filter, arrayV)
	}

	for _, value := range filter.values {
		if v == value {
			return true
		}
	}
	return false

}

func matchFilters(filters []filter, m map[string]interface{}) bool {
	for _, filter := range filters {
		if !matches(filter, m) {
			return false
		}
	}
	return true
}

func doFiltering(d *schema.ResourceData, items []map[string]interface{}) (map[string]interface{}, error) {
	if len(items) == 0 {
		return nil, nil
	}

	filters, filtersOk := d.GetOk("filter")
	var filteredItems []map[string]interface{}
	if filtersOk {
		f := buildFilters(filters.(*schema.Set))
		for _, product := range items {
			if matchFilters(f, product) {
				filteredItems = append(filteredItems, product)
			}
		}
	} else {
		filteredItems = items
	}

	first := d.Get("first")
	if (first == nil || !d.Get("first").(bool)) && len(filteredItems) != 1 {
		return nil, fmt.Errorf("found %s matches. set first = true or modify your filters", fmt.Sprint(len(filteredItems)))
	}
	if len(filteredItems) < 1 {
		return nil, nil
	}

	return filteredItems[0], nil
}
