package hivelocity

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func TestMakeVlanCreatePayload(t *testing.T) {
	r := resourceVlan()
	config := map[string]interface{}{}
	d := schema.TestResourceDataRaw(t, r.Schema, config)
	d.Set("device_ids", []int{1, 2, 3})

	payload := makeVlanCreatePayload(d)
	expectedDeviceIds := []int32{1, 2, 3}
	if !reflect.DeepEqual(payload.DeviceIds, expectedDeviceIds) {
		t.Fatalf("Error matching output and expected: %#v vs %#v", payload.DeviceIds, expectedDeviceIds)
	}
}

func TestGetVlanDeviceIds(t *testing.T) {
	cases := []struct {
		vlan              swagger.Vlan
		expectedDeviceIds []int
	}{
		{
			vlan: swagger.Vlan{
				DeviceIds: []int32{4, 5, 6},
			},
			expectedDeviceIds: []int{4, 5, 6},
		},
		{
			vlan: swagger.Vlan{
				DeviceIds: []int32{},
			},
			expectedDeviceIds: []int{},
		},
	}

	for _, c := range cases {
		out := getVlanDeviceIds(&c.vlan)
		if !reflect.DeepEqual(out, c.expectedDeviceIds) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expectedDeviceIds)
		}
	}
}

func TestCompareArraysDeviceIds(t *testing.T) {
	a := []int32{3, 2, 1}
	b := []int32{1, 2, 3}

	if !arraysEqual(a, b) {
		t.Fatalf("Arrays should be equal: %#v vs %#v", a, b)
	}

	a = append(a, 4)
	if arraysEqual(a, b) {
		t.Fatalf("Arrays should be different: %#v vs %#v", a, b)
	}
}

// Test basic behaviour for the vlan resource
func TestAccHivelocityVlan_basic(t *testing.T) {
	var vlan swagger.Vlan
	name := "hivelocity_vlan.test_vlan"

	// The device id is still hardcoded, which is not ideal - needs improvement
	testDevice := os.Getenv("HIVELOCITY_TEST_DEVICE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVlanResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVlanResource(testDevice),
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckVlanResourceExists(name, &vlan),
					// verify remote values
					testAccCheckExampleVlanValues(&vlan, testDevice),
					// verify local values
					resource.TestCheckResourceAttrSet(name, "device_ids.#"),
				),
			},
		},
	})
}

func testAccVlanResource(device_ids string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf(`
		resource "hivelocity_vlan" "test_vlan" {
			device_ids = [%s]
		}`, device_ids)
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

func testAccCheckExampleVlanValues(vlan *swagger.Vlan, device string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] Checking values: %#v", vlan)

		deviceId, _ := strconv.Atoi(device)

		if len(vlan.DeviceIds) != 1 || vlan.DeviceIds[0] != int32(deviceId) {
			return fmt.Errorf("bad devices list: %#v", vlan.DeviceIds)
		}

		return nil
	}
}
