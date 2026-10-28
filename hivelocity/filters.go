package hivelocity

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Filter the results. Every filter block must match (logical AND); a block matches when the " +
			"attribute named by `name` equals any of `values`.",
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Description: "Attribute name to filter on",
					Type:        schema.TypeString,
					Required:    true,
				},

				"values": {
					Description: "Accepted values for the attribute",
					Type:        schema.TypeList,
					Required:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func dataSourceFilterFirstSchema() *schema.Schema {
	return &schema.Schema{
		Description: "When true (the default) return the first match instead of erroring when more than one " +
			"item matches the filters.",
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}
}

type filter struct {
	name   string
	values []string
}

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
			return nil, fmt.Errorf("cannot filter on %v", values)
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
		if searchValue == nil {
			continue
		}
		// recurse when search value is an array
		rt := reflect.TypeOf(searchValue)
		if rt.Kind() == reflect.Array || rt.Kind() == reflect.Slice {
			if nested, ok := searchValue.([]interface{}); ok {
				if matches(filterValues, nested) {
					return true
				}
				continue
			}
		}
		stringSearchValue := fmt.Sprint(searchValue)
		for _, filterValue := range filterValues {
			if filterValue == stringSearchValue {
				return true
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

// validateFilterNames makes sure every filter refers to an attribute that exists
// on at least one item, so typos produce an error instead of silently matching
// nothing.
func validateFilterNames(filters []filter, items []map[string]interface{}) error {
	known := map[string]bool{}
	for _, item := range items {
		for k := range item {
			known[k] = true
		}
	}

	var unknown []string
	for _, f := range filters {
		if !known[f.name] {
			unknown = append(unknown, f.name)
		}
	}
	if len(unknown) == 0 {
		return nil
	}

	valid := make([]string, 0, len(known))
	for k := range known {
		valid = append(valid, k)
	}
	sort.Strings(valid)
	sort.Strings(unknown)

	return fmt.Errorf("unknown filter name(s) %s; valid names are: %s",
		strings.Join(unknown, ", "), strings.Join(valid, ", "))
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

		if err := validateFilterNames(filter_params, items); err != nil {
			return nil, err
		}
	}

	// Fallback to default filters if any
	if len(filter_params) == 0 && len(defaultFilters) > 0 {
		filter_params = defaultFilters
	}

	if len(filter_params) > 0 {
		for _, f := range filter_params {
			out, _ := json.Marshal(f.values)
			log.Printf("[DEBUG] Filtering on name %v with values %v", f.name, string(out))
		}

		for _, item := range items {
			if matchFilters(filter_params, item) {
				filteredItems = append(filteredItems, item)
			}
		}
	} else {
		log.Printf("[DEBUG] No filtering is done")
		filteredItems = items
	}

	log.Printf("[DEBUG] Filter result count %v", len(filteredItems))

	first := d.Get("first")
	if (first == nil || !d.Get("first").(bool)) && len(filteredItems) != 1 {
		return nil, fmt.Errorf("found %d matches. set first = true or modify your filters", len(filteredItems))
	}

	if len(filteredItems) == 0 {
		return nil, nil
	}

	return filteredItems[0], nil
}
