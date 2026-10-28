package hivelocity

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func testFilterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"first":  dataSourceFilterFirstSchema(),
	}
}

func testFilterData(t *testing.T, first bool, filters ...map[string]interface{}) *schema.ResourceData {
	t.Helper()
	raw := map[string]interface{}{"first": first}
	if len(filters) > 0 {
		set := make([]interface{}, 0, len(filters))
		for _, f := range filters {
			set = append(set, f)
		}
		raw["filter"] = set
	}
	return schema.TestResourceDataRaw(t, testFilterSchema(), raw)
}

func testItems() []map[string]interface{} {
	return []map[string]interface{}{
		{"product_id": 1, "stock": "available", "data_center": "TPA1", "tags": []interface{}{"a", "b"}},
		{"product_id": 2, "stock": "limited", "data_center": "TPA1", "tags": []interface{}{"c"}},
		{"product_id": 3, "stock": "unavailable", "data_center": "NYC1", "tags": []interface{}{}},
	}
}

func TestDoFiltering_noFilterFirst(t *testing.T) {
	item, err := doFiltering(testFilterData(t, true), testItems(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if item["product_id"] != 1 {
		t.Fatalf("expected first item, got %v", item)
	}
}

func TestDoFiltering_noFilterNotFirst(t *testing.T) {
	_, err := doFiltering(testFilterData(t, false), testItems(), nil)
	if err == nil || !strings.Contains(err.Error(), "found 3 matches") {
		t.Fatalf("expected multi-match error, got %v", err)
	}
}

func TestDoFiltering_matchesStringAndInt(t *testing.T) {
	d := testFilterData(t, true,
		map[string]interface{}{"name": "data_center", "values": []interface{}{"TPA1"}},
		map[string]interface{}{"name": "product_id", "values": []interface{}{"2"}},
	)
	item, err := doFiltering(d, testItems(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if item["product_id"] != 2 {
		t.Fatalf("expected product 2, got %v", item)
	}
}

func TestDoFiltering_matchesListValues(t *testing.T) {
	d := testFilterData(t, true,
		map[string]interface{}{"name": "tags", "values": []interface{}{"c"}},
	)
	item, err := doFiltering(d, testItems(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if item["product_id"] != 2 {
		t.Fatalf("expected product 2, got %v", item)
	}
}

func TestDoFiltering_noMatch(t *testing.T) {
	d := testFilterData(t, true,
		map[string]interface{}{"name": "stock", "values": []interface{}{"nope"}},
	)
	item, err := doFiltering(d, testItems(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if item != nil {
		t.Fatalf("expected no match, got %v", item)
	}
}

func TestDoFiltering_unknownFilterName(t *testing.T) {
	d := testFilterData(t, true,
		map[string]interface{}{"name": "produc_id", "values": []interface{}{"1"}},
	)
	_, err := doFiltering(d, testItems(), nil)
	if err == nil || !strings.Contains(err.Error(), "unknown filter name(s) produc_id") {
		t.Fatalf("expected unknown filter error, got %v", err)
	}
	if !strings.Contains(err.Error(), "product_id") {
		t.Fatalf("expected valid names in error, got %v", err)
	}
}

func TestDoFiltering_defaultFilters(t *testing.T) {
	item, err := doFiltering(testFilterData(t, true), testItems(), []filter{{name: "data_center", values: []string{"NYC1"}}})
	if err != nil {
		t.Fatal(err)
	}
	if item["product_id"] != 3 {
		t.Fatalf("expected product 3, got %v", item)
	}
}

func TestDoFiltering_emptyItems(t *testing.T) {
	if _, err := doFiltering(testFilterData(t, true), nil, nil); err == nil {
		t.Fatal("expected error for empty item list")
	}
}
