package hivelocity

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

// Test basic behaviour for the vlan resource
func TestAccHivelocityVlan_basic(t *testing.T) {
	var vlan swagger.Vlan
	name := "hivelocity_vlan.test_vlan"

	// Test data is still hardcoded, which is not ideal - needs improvement
	testFacilityCode := os.Getenv("HIVELOCITY_TEST_FACILITY_CODE")
	testPortId := os.Getenv("HIVELOCITY_TEST_PORT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVlanResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVlanResource(testFacilityCode, testPortId),
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckVlanResourceExists(name, &vlan),
					// verify remote values
					testAccCheckExampleVlanValues(&vlan, testPortId),
					// verify local values
					resource.TestCheckResourceAttrSet(name, "port_ids.#"),
				),
			},
		},
	})
}

func testAccVlanResource(facilityCode string, portId string) string {
	time.Sleep(1 * time.Second)

	return fmt.Sprintf(`
		resource "hivelocity_vlan" "test_vlan" {
			facility_code = "%s"
			type = "private"
			port_ids = [%s]
		}`, facilityCode, portId)
}

// testAccCheckVlanResourceDestroy verifies the Vlan has been destroyed
func testAccCheckVlanResourceDestroy(s *terraform.State) error {
	// retrieve the connection established in Provider configuration
	hv := testAccProvider.Meta().(*Client)

	// loop through the resources in state, verifying each vlan is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hivelocity_vlan" {
			continue
		}

		vlanId, _ := strconv.Atoi(rs.Primary.ID)
		_, response, err := hv.client.VLANApi.GetVlanIdResource(hv.auth, int32(vlanId), nil)
		if response.StatusCode == 404 {
			continue
		}
		if err == nil {
			return fmt.Errorf("Vlan (%s) still exists.", rs.Primary.ID)
		}

		return err
	}

	return nil
}

// Queries the API and returns the matching Vlan object
func testAccCheckVlanResourceExists(name string, vlan *swagger.Vlan) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("VLAN ID is not set")
		}

		vlanId, _ := strconv.Atoi(rs.Primary.ID)
		hv := testAccProvider.Meta().(*Client)
		vlanResponse, _, err := hv.client.VLANApi.GetVlanIdResource(hv.auth, int32(vlanId), nil)
		if err != nil {
			return fmt.Errorf("Error getting vlan: %s", err)
		}

		*vlan = vlanResponse

		return nil
	}
}

func testAccCheckExampleVlanValues(vlan *swagger.Vlan, testPortId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] Checking values: %#v", vlan)

		portId, _ := strconv.Atoi(testPortId)

		if len(vlan.PortIds) != 1 || vlan.PortIds[0] != int32(portId) {
			return fmt.Errorf("bad ports list: %#v", vlan.PortIds)
		}

		return nil
	}
}
