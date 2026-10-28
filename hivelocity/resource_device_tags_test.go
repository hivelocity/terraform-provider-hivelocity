package hivelocity

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Requires an existing device the test may re-tag:
//
//	HIVELOCITY_TEST_DEVICE_ID=12345
func TestAccHivelocityDeviceTags_basic(t *testing.T) {
	testAccRequireEnv(t, "HIVELOCITY_TEST_DEVICE_ID")
	deviceId := os.Getenv("HIVELOCITY_TEST_DEVICE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTagsConfig(deviceId, `["tf-acc", "one"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hivelocity_device_tags.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("hivelocity_device_tags.test", "tags.0", "tf-acc"),
				),
			},
			{
				Config: testAccDeviceTagsConfig(deviceId, `["tf-acc"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hivelocity_device_tags.test", "tags.#", "1"),
				),
			},
			{
				ResourceName:      "hivelocity_device_tags.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDeviceTagsConfig(deviceId, tags string) string {
	return fmt.Sprintf(`
resource "hivelocity_device_tags" "test" {
  device_id = %s
  tags      = %s
}
`, deviceId, tags)
}
