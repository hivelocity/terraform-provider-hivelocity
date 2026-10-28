package hivelocity

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// testAccProviderFactories is used by acceptance tests (TF_ACC=1).
var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"hivelocity": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

// testAccPreCheck validates the necessary test API keys exist
// in the testing environment
func testAccPreCheck(t *testing.T) {
	t.Helper()
	if v := os.Getenv("HIVELOCITY_API_KEY"); v == "" {
		t.Fatal("HIVELOCITY_API_KEY must be set for acceptance tests")
	}
}

// testAccRequireEnv skips the test unless every named environment variable is set.
func testAccRequireEnv(t *testing.T, names ...string) {
	t.Helper()
	for _, name := range names {
		if os.Getenv(name) == "" {
			t.Skipf("%s must be set to run this acceptance test", name)
		}
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("provider schema failed validation: %s", err)
	}
}
