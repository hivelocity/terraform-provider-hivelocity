package hivelocity

import (
	"reflect"
	"testing"
	"time"
)

func TestParseTwoPartId(t *testing.T) {
	a, b, err := parseTwoPartId("12:34", "<a>:<b>")
	if err != nil || a != 12 || b != 34 {
		t.Fatalf("got %d %d %v", a, b, err)
	}
	for _, bad := range []string{"12", "12:", ":34", "a:b", "12:34:56x"} {
		if _, _, err := parseTwoPartId(bad, "<a>:<b>"); err == nil {
			t.Errorf("expected error for %q", bad)
		}
	}
}

func TestParseDnsRecordAId(t *testing.T) {
	domain, name, err := parseDnsRecordAId("example.com/www")
	if err != nil || domain != "example.com" || name != "www" {
		t.Fatalf("got %q %q %v", domain, name, err)
	}
	if _, _, err := parseDnsRecordAId("example.com"); err == nil {
		t.Fatal("expected error without record name")
	}
}

func TestParseVpsNetworkAttachmentId(t *testing.T) {
	dev, net, err := parseVpsNetworkAttachmentId("12:net-abc")
	if err != nil || dev != 12 || net != "net-abc" {
		t.Fatalf("got %d %q %v", dev, net, err)
	}
	if _, _, err := parseVpsNetworkAttachmentId("x:net"); err == nil {
		t.Fatal("expected error for non-numeric device id")
	}
}

func TestPermissionNames(t *testing.T) {
	cases := map[string][]string{
		`["b", "a"]`:                               {"a", "b"},
		`[{"name": "b"}, {"name": "a"}]`:           {"a", "b"},
		`{"permissions": ["z", "y"]}`:              {"y", "z"},
		`{"device.view": true, "device.power": 1}`: {"device.power", "device.view"},
		`[]`: {},
	}
	for body, want := range cases {
		got, err := permissionNames([]byte(body))
		if err != nil {
			t.Fatalf("%s: %v", body, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got %v want %v", body, got, want)
		}
	}
	if _, err := permissionNames([]byte(`"nope"`)); err == nil {
		t.Fatal("expected error for unrecognised shape")
	}
}

func TestTimeString(t *testing.T) {
	if got := timeString(time.Time{}); got != "" {
		t.Fatalf("zero time should be empty, got %q", got)
	}
	ts := time.Date(2026, 1, 2, 3, 4, 5, 0, time.FixedZone("x", 3600))
	if got := timeString(ts); got != "2026-01-02T02:04:05Z" {
		t.Fatalf("got %q", got)
	}
}

func TestFlattenVpsNics(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{"nicId": "n1", "networkId": "net1", "ipAddress": "10.0.0.2", "vlanTag": float64(7)},
	}
	nics, err := flattenVpsNics(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(nics) != 1 {
		t.Fatalf("expected 1 nic, got %d", len(nics))
	}
	nic := nics[0].(map[string]interface{})
	if nic["nic_id"] != "n1" || nic["network_id"] != "net1" || nic["ip_address"] != "10.0.0.2" || nic["vlan_tag"] != 7 {
		t.Fatalf("unexpected nic %v", nic)
	}
}

func TestApiTimeString(t *testing.T) {
	if got := apiTimeString(nil); got != "" {
		t.Fatalf("nil: %q", got)
	}
	if got := apiTimeString("2027-01-01T00:00:00Z"); got != "2027-01-01T00:00:00Z" {
		t.Fatalf("string: %q", got)
	}
	if got := apiTimeString(float64(1767225600)); got != "2026-01-01T00:00:00Z" {
		t.Fatalf("number: %q", got)
	}
}

func TestStringMapFromInterface(t *testing.T) {
	got := stringMapFromInterface(map[string]interface{}{"a": "x", "b": 2})
	if got["a"] != "x" || got["b"] != "2" {
		t.Fatalf("got %v", got)
	}
	if got := stringMapFromInterface(nil); len(got) != 0 {
		t.Fatalf("nil should yield empty map, got %v", got)
	}
}
