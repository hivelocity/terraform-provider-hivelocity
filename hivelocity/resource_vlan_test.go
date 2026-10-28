package hivelocity

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Requires a facility code and the ID of a device port in that facility:
//
//	HIVELOCITY_TEST_FACILITY_CODE=TPA1 HIVELOCITY_TEST_PORT_ID=12345
func TestAccHivelocityVlan_basic(t *testing.T) {
	testAccRequireEnv(t, "HIVELOCITY_TEST_FACILITY_CODE", "HIVELOCITY_TEST_PORT_ID")
	facility := os.Getenv("HIVELOCITY_TEST_FACILITY_CODE")
	portId := os.Getenv("HIVELOCITY_TEST_PORT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVlanConfig(facility, "private", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hivelocity_vlan.test", "facility_code", facility),
					resource.TestCheckResourceAttr("hivelocity_vlan.test", "type", "private"),
					resource.TestCheckResourceAttrSet("hivelocity_vlan.test", "vlan_tag"),
					resource.TestCheckResourceAttr("hivelocity_vlan.test", "port_ids.#", "0"),
				),
			},
			{
				Config: testAccVlanConfig(facility, "private", portId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hivelocity_vlan.test", "port_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("hivelocity_vlan.test", "port_ids.*", portId),
				),
			},
			{
				ResourceName:      "hivelocity_vlan.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVlanConfig(facility, vlanType, portId string) string {
	ports := "[]"
	if portId != "" {
		ports = fmt.Sprintf("[%s]", portId)
	}
	return fmt.Sprintf(`
resource "hivelocity_vlan" "test" {
  facility_code = %q
  type          = %q
  port_ids      = %s
}
`, facility, vlanType, ports)
}
