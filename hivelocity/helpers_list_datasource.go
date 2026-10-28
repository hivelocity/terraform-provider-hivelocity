package hivelocity

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// toGenericList marshals an API list response (a slice of generated models)
// into generic maps with snake_case keys so it can be filtered and set into
// state.
func toGenericList(items interface{}) ([]map[string]interface{}, error) {
	data, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return convertKeysOfList(list), nil
}

// toGenericMap marshals a single API model into a generic map with snake_case
// keys.
func toGenericMap(item interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return convertKeys(m), nil
}

// readFilteredList applies the data source's `filter`/`first` arguments to a
// list of API objects and copies the matching item into state, using idKey as
// the resource ID.
func readFilteredList(
	d *schema.ResourceData,
	items interface{},
	s map[string]*schema.Schema,
	idKey string,
	defaultFilters []filter,
) diag.Diagnostics {
	list, err := toGenericList(items)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(list) == 0 {
		return diag.Errorf("the API returned no items to filter on")
	}

	item, err := doFiltering(d, list, defaultFilters)
	if err != nil {
		return diag.FromErr(err)
	}
	if item == nil {
		return diag.Errorf("no item matched the given filters")
	}

	return setGenericItem(d, item, s, idKey)
}

// setGenericItem copies the schema-known keys of item into state and sets the ID.
func setGenericItem(d *schema.ResourceData, item map[string]interface{}, s map[string]*schema.Schema, idKey string) diag.Diagnostics {
	for k, v := range filterNonSchemaKeys(item, s) {
		if k == "filter" || k == "first" {
			continue
		}
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("setting %q: %s", k, err)
		}
	}

	id := fmt.Sprint(item[idKey])
	if id == "" || id == "<nil>" {
		return diag.Errorf("item has no %q to use as ID", idKey)
	}
	d.SetId(id)
	return nil
}

// computedString and friends keep data source schemas short.
func computedString(desc string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeString, Computed: true, Description: desc}
}

func computedInt(desc string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeInt, Computed: true, Description: desc}
}

func computedBool(desc string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeBool, Computed: true, Description: desc}
}

func computedFloat(desc string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeFloat, Computed: true, Description: desc}
}

func computedStringList(desc string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: desc}
}

func computedIntList(desc string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeInt}, Description: desc}
}

func computedStringMap(desc string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: desc}
}
