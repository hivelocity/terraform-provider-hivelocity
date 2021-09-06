package hivelocity

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

// testAccPreCheck validates the necessary test API keys exist
// in the testing environment
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("HIVELOCITY_API_KEY"); v == "" {
		t.Fatal("HIVELOCITY_API_KEY must be set for acceptance tests")
	}

	if v := os.Getenv("HIVELOCITY_TEST_DEVICE_ID"); v == "" {
		t.Fatal("HIVELOCITY_TEST_DEVICE_ID must be set for acceptance tests")
	}
}

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"hivelocity": testAccProvider,
	}
}
