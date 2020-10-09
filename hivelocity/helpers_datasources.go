package hivelocity

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/iancoleman/strcase"
)

func convertKeys(m map[string]interface{}) map[string]interface{} {
	new := make(map[string]interface{}, 0)
	for k, v := range m {
		fixed := strcase.ToSnake(k)
		new[fixed] = v
		delete(m, k)
	}
	return new
}

func convertKeysOfList(l []map[string]interface{}) []map[string]interface{} {
	for i := 0; i < len(l); i++ {
		l[i] = convertKeys(l[i])
	}
	return l
}

func forceValuesToStrings(m map[string]interface{}) map[string]interface{} {
	new := make(map[string]interface{}, 0)
	for k, v := range m {
		str := fmt.Sprintf("%v", v)
		new[k] = str
	}
	return new
}

func forceValuesToStringOfList(l []map[string]interface{}, key string) []map[string]interface{} {
	for i := 0; i < len(l); i++ {
		if l[i][key] != nil {
			l[i][key] = forceValuesToStrings(l[i][key].(map[string]interface{}))
		}
	}
	return l
}

func getSchemaKeys(schemas map[string]*schema.Schema) []string {
	keys := make([]string, 0)
	for k, _ := range schemas {
		keys = append(keys, k)
	}
	return keys
}

func filterNonSchemaKeys(m map[string]interface{}, schema map[string]*schema.Schema) map[string]interface{} {
	keys := getSchemaKeys(schema)
	new := make(map[string]interface{}, 0)
	for k, v := range m {
		for _, a := range keys {
			if a == k {
				new[k] = v
			}
		}
	}
	return new
}

func filterNonSchemaKeysForList(l []map[string]interface{}, schema map[string]*schema.Schema) []map[string]interface{} {
	for i := 0; i < len(l); i++ {
		l[i] = filterNonSchemaKeys(l[i], schema)
	}
	return l
}

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

func matches(filter filter, m map[string]interface{}) bool {
	v, ok := m[filter.name]
	if !ok {
		return false
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
