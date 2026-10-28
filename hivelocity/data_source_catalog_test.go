package hivelocity

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Read-only data sources that only need an API key.
func TestAccHivelocityReadOnlyDataSources(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "hivelocity_webhook_events" "all" {}
data "hivelocity_service_types" "all" {}
data "hivelocity_permissions" "all" {}
data "hivelocity_user_permissions" "me" {}
data "hivelocity_maintenance_events" "all" {}
data "hivelocity_device_tags" "all" {}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hivelocity_webhook_events.all", "events.#"),
					resource.TestCheckResourceAttrSet("data.hivelocity_service_types.all", "service_types.#"),
					resource.TestCheckResourceAttrSet("data.hivelocity_permissions.all", "permissions.#"),
					resource.TestCheckResourceAttrSet("data.hivelocity_user_permissions.me", "permissions.#"),
					resource.TestCheckResourceAttrSet("data.hivelocity_maintenance_events.all", "events.#"),
					resource.TestCheckResourceAttrSet("data.hivelocity_device_tags.all", "tags.#"),
				),
			},
		},
	})
}

// Requires a product ID to inspect:
//
//	HIVELOCITY_TEST_PRODUCT_ID=580
func TestAccHivelocityProductDataSources(t *testing.T) {
	testAccRequireEnv(t, "HIVELOCITY_TEST_PRODUCT_ID")
	productId := os.Getenv("HIVELOCITY_TEST_PRODUCT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "hivelocity_product" "test" {
  first = true
  filter {
    name   = "product_id"
    values = [%[1]q]
  }
}
data "hivelocity_product_operating_systems" "test" {
  product_id = %[1]s
}
data "hivelocity_product_options" "test" {
  product_id = %[1]s
}
`, productId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.hivelocity_product.test", "product_id", productId),
					resource.TestCheckResourceAttrSet("data.hivelocity_product.test", "data_center"),
					resource.TestCheckResourceAttrSet("data.hivelocity_product_operating_systems.test", "operating_systems.0.name"),
					resource.TestCheckResourceAttrSet("data.hivelocity_product_options.test", "options.#"),
				),
			},
		},
	})
}

// Requires an existing device:
//
//	HIVELOCITY_TEST_DEVICE_ID=12345
func TestAccHivelocityDeviceDataSources(t *testing.T) {
	testAccRequireEnv(t, "HIVELOCITY_TEST_DEVICE_ID")
	deviceId := os.Getenv("HIVELOCITY_TEST_DEVICE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "hivelocity_device_port" "test" {
  device_id = %[1]s
}
data "hivelocity_device_ips" "test" {
  device_id = %[1]s
}
data "hivelocity_device_hardware" "test" {
  device_id = %[1]s
}
data "hivelocity_device" "test" {
  first = true
  filter {
    name   = "device_id"
    values = [%[1]q]
  }
}
`, deviceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hivelocity_device_port.test", "port_id"),
					resource.TestCheckResourceAttrSet("data.hivelocity_device_ips.test", "ip_assignments.#"),
					resource.TestCheckResourceAttrSet("data.hivelocity_device_hardware.test", "drive_count"),
					resource.TestCheckResourceAttr("data.hivelocity_device.test", "device_id", deviceId),
				),
			},
		},
	})
}
