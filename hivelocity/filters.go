package hivelocity

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/prometheus/common/log"
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

// TODO: Find way to access `map[string]*schema.Schema` so we can validate fields
func buildFilters(set *schema.Set) ([]filter, error) {
	var filters []filter

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var values []string
		for _, value := range m["values"].([]interface{}) {
			values = append(values, value.(string))
		}

		name := m["name"].(string)

		if name == "filter" {
			return nil, errors.New(fmt.Sprintf("Cannot filter on %v", values))
		}

		filters = append(filters, filter{
			name:   name,
			values: values,
		})
	}
	return filters, nil
}

func matches(filterValues []string, searchValues []interface{}) bool {
	// match each element, recursively matches on more arrays
	for _, searchValue := range searchValues {
		// recurse when search value is an array
		rt := reflect.TypeOf(searchValue)
		if rt.Kind() == reflect.Array {
			if matches(filterValues, searchValue.([]interface{})) {
				return true
			}
		} else {
			stringSearchValue := fmt.Sprint(searchValue)
			for _, filterValue := range filterValues {
				if filterValue == stringSearchValue {
					return true
				}
			}
		}
	}

	return false
}

func matchFilters(filters []filter, m map[string]interface{}) bool {
	for _, filter := range filters {
		filterValues := make([]interface{}, 1)
		searchValue, ok := m[filter.name]
		if !ok {
			return false
		}

		filterValues[0] = searchValue

		if !matches(filter.values, filterValues) {
			return false
		}
	}
	return true
}

func doFiltering(
	d *schema.ResourceData,
	items []map[string]interface{},
	defaultFilters []filter,
) (map[string]interface{}, error) {
	var filter_params []filter
	var filteredItems []map[string]interface{}

	if len(items) == 0 {
		return nil, errors.New("no items to filter on")
	}

	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		filter_params_, err := buildFilters(filters.(*schema.Set))
		if err != nil {
			return nil, err
		}
		filter_params = filter_params_
	}

	// Fallback to defualt filters if any
	if len(filter_params) == 0 && len(defaultFilters) > 0 {
		filter_params = defaultFilters
	}

	if len(filter_params) > 0 {
		for _, f := range filter_params {
			out, _ := json.Marshal(f.values)
			log.Infof("Filtering on name %v with values %v", f.name, string(out))
		}

		for _, item := range items {
			if matchFilters(filter_params, item) {
				filteredItems = append(filteredItems, item)
			}
		}
	} else {
		log.Infof("No filtering is done")
		filteredItems = items
	}

	log.Infof("Filter result count %v", len(filteredItems))

	first := d.Get("first")
	if (first == nil || !d.Get("first").(bool)) && len(filteredItems) != 1 {
		return nil, fmt.Errorf("found %s matches. set first = true or modify your filters", fmt.Sprint(len(filteredItems)))
	}

	if len(filteredItems) == 0 {
		return nil, nil
	}

	return filteredItems[0], nil
}
