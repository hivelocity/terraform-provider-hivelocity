package hivelocity

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestConvertKeys(t *testing.T) {
	got := convertKeys(map[string]interface{}{"deviceId": 1, "primaryIp": "1.2.3.4", "already_snake": true})
	want := map[string]interface{}{"device_id": 1, "primary_ip": "1.2.3.4", "already_snake": true}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestFilterNonSchemaKeys(t *testing.T) {
	s := map[string]*schema.Schema{"keep": {Type: schema.TypeString}}
	got := filterNonSchemaKeys(map[string]interface{}{"keep": "a", "drop": "b"}, s)
	if len(got) != 1 || got["keep"] != "a" {
		t.Fatalf("unexpected result %v", got)
	}
}

func TestNormalizeJsonString(t *testing.T) {
	a, err := normalizeJsonString(`{"b":1,"a":[1,2]}`)
	if err != nil {
		t.Fatal(err)
	}
	b, err := normalizeJsonString("{\n  \"a\": [1, 2],\n  \"b\": 1\n}")
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Fatalf("expected equal normalized output:\n%s\n%s", a, b)
	}
	if _, err := normalizeJsonString("not json"); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseInventoryProducts(t *testing.T) {
	flat, err := parseInventoryProducts([]byte(`[{"product_id": 1}, {"product_id": 2}]`))
	if err != nil || len(flat) != 2 {
		t.Fatalf("flat parse failed: %v %v", flat, err)
	}
	nested, err := parseInventoryProducts([]byte(`{"TPA1": [{"product_id": 1}], "NYC1": [{"product_id": 2}, {"product_id": 3}]}`))
	if err != nil || len(nested) != 3 {
		t.Fatalf("nested parse failed: %v %v", nested, err)
	}
	if _, err := parseInventoryProducts([]byte(`{"code": 400}`)); err == nil {
		t.Fatal("expected error for error-shaped response")
	}
}

func TestExpandHelpers(t *testing.T) {
	if got := expandInt32List([]interface{}{1, 2}); !reflect.DeepEqual(got, []int32{1, 2}) {
		t.Fatalf("expandInt32List list: %v", got)
	}
	if got := expandInt32List(nil); got == nil || len(got) != 0 {
		t.Fatalf("expandInt32List nil should be empty non-nil slice: %#v", got)
	}
	if got := expandStringList(schema.NewSet(schema.HashString, []interface{}{"a"})); !reflect.DeepEqual(got, []string{"a"}) {
		t.Fatalf("expandStringList set: %v", got)
	}
	if got := flattenInt32List([]int32{3}); !reflect.DeepEqual(got, []int{3}) {
		t.Fatalf("flattenInt32List: %v", got)
	}
	if got := expandStringMap(map[string]interface{}{"k": "v"}); got["k"] != "v" {
		t.Fatalf("expandStringMap: %v", got)
	}
}
