package hivelocity

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

// Return device tags
func getTags(d *schema.ResourceData, deviceKey string) []string {
	key := fmt.Sprintf("%stags", deviceKey)
	var tags []string

	for _, v := range d.Get(key).([]interface{}) {
		tags = append(tags, v.(string))
	}

	return tags
}

func waitForDevices(timeout time.Duration, hv *Client, orderId int32, newDevices []swagger.BareMetalDevice) (interface{}, error) {
	waitForDevice := &resource.StateChangeConf{
		Pending: []string{
			"waiting",
		},
		Target: []string{
			"ok",
		},
		Refresh: func() (interface{}, string, error) {
			devices, _, err := hv.client.BareMetalDevicesApi.GetBareMetalDeviceResource(hv.auth, nil)
			if err != nil {
				return 0, "", err
			}

			// Look for the number of devices specified
			var devicesFound []swagger.BareMetalDevice
			for _, device := range devices {
				if device.OrderId == orderId {
					devicesFound = append(devicesFound, device)
				}
			}
			if len(devicesFound) == len(newDevices) {
				return devicesFound, "ok", nil
			}
			return nil, "waiting", nil
		},
		Timeout:                   timeout,
		Delay:                     30 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 1,
	}
	return waitForDevice.WaitForState()
}

func waitForOrder(timeout time.Duration, hv *Client, orderId int32) (interface{}, error) {
	waitForOrder := &resource.StateChangeConf{
		Pending: []string{
			"verification",
			"lead",
			"provisioning",
			"assembling",
		},
		Target: []string{
			"complete",
		},
		Refresh: func() (interface{}, string, error) {
			resp, _, err := hv.client.OrderApi.GetOrderIdResource(hv.auth, orderId, nil)
			if err != nil {
				return 0, "", err
			}
			return resp, resp.Status, nil
		},
		Timeout:                   timeout,
		Delay:                     5 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 1,
	}
	return waitForOrder.WaitForState()
}

func waitForDevicePowerOff(d *schema.ResourceData, hv *Client, deviceId int32) (interface{}, error) {
	waitForDevice := &resource.StateChangeConf{
		Pending: []string{
			"waiting",
		},
		Target: []string{
			"ok",
		},
		Refresh: func() (interface{}, string, error) {
			device, _, err := hv.client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(hv.auth, deviceId, nil)

			if err != nil {
				return 0, "", err
			}

			if device.PowerStatus == "OFF" {
				return device, "ok", nil
			}

			return nil, "waiting", nil
		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     5 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 1,
	}
	return waitForDevice.WaitForState()
}

func waitForDeviceReload(d *schema.ResourceData, hv *Client, deviceId int32) (interface{}, error) {
	waitForDevice := &resource.StateChangeConf{
		Pending: []string{"waiting"},
		Target:  []string{"ok"},
		Refresh: func() (interface{}, string, error) {
			device, _, err := hv.client.DeviceApi.GetDeviceIdResource(hv.auth, deviceId, nil)
			if err != nil {
				return 0, "", err
			}
			if device.Metadata != nil {
				metadataValue := *(device.Metadata)
				metadata := metadataValue.(map[string]interface{})
				spsStatus, ok := metadata["sps_status"]
				if ok && spsStatus == "InUse" {
					return device, "ok", nil
				}
			}
			return nil, "waiting", nil
		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     30 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 1,
		NotFoundChecks:            360, // 1h timeout / 10s delay between requests
	}
	return waitForDevice.WaitForState()
}

// Returns true if deviceId is contained in newDevices list
func deviceInList(deviceId int, newDevices []map[string]interface{}) bool {
	for _, device := range newDevices {
		if device["device_id"] == deviceId {
			return true
		}
	}

	return false
}

// Builds a map filled with Bare Metal Device fields from a response
func buildDeviceForState(device swagger.BareMetalDevice) map[string]interface{} {
	stateDevice := make(map[string]interface{})

	stateDevice["device_id"] = int(device.DeviceId)
	stateDevice["hostname"] = device.Hostname
	stateDevice["location_name"] = device.LocationName
	stateDevice["order_id"] = int(device.OrderId)
	stateDevice["os_name"] = device.OsName
	stateDevice["power_status"] = device.PowerStatus
	stateDevice["primary_ip"] = device.PrimaryIp
	stateDevice["product_id"] = int(device.ProductId)
	stateDevice["product_name"] = device.ProductName
	stateDevice["service_id"] = int(device.ServiceId)
	stateDevice["tags"] = device.Tags
	stateDevice["public_ssh_key_id"] = int(device.PublicSshKeyId)
	stateDevice["script"] = device.Script

	return stateDevice
}

// Creates a list of devices
func createDevices(hv *Client, orderGroupId int32, devicesToCreate []map[string]interface{}) error {
	// Use the batch endpoint to create the Bare Metal Devices all at once
	var bareMetalDevicesPayload []swagger.BareMetalDeviceCreate

	for _, device := range devicesToCreate {
		var tags []string
		for _, tag := range device["tags"].([]interface{}) {
			tags = append(tags, tag.(string))
		}

		bareMetalDeviceCreate := swagger.BareMetalDeviceCreate{
			ProductId:      int32(device["product_id"].(int)),
			Hostname:       device["hostname"].(string),
			OsName:         device["os_name"].(string),
			LocationName:   device["location_name"].(string),
			Script:         device["script"].(string),
			Period:         device["period"].(string),
			PublicSshKeyId: int32(device["public_ssh_key_id"].(int)),
			ForceDeviceId:  int32(device["force_device_id"].(int)),
			Tags:           tags,
		}
		bareMetalDevicesPayload = append(bareMetalDevicesPayload, bareMetalDeviceCreate)
	}

	bareMetalDeviceBatchCreatePayload := swagger.BareMetalDeviceBatchCreate{
		OrderGroupId: orderGroupId,
		Devices:      bareMetalDevicesPayload,
	}
	bareMetalDeviceResponse, _, err := hv.client.BareMetalDevicesApi.PostBareMetalDeviceBatchResource(hv.auth, bareMetalDeviceBatchCreatePayload, nil)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return fmt.Errorf("POST /bare-metal-devices/batch failed! (%s)\n\n %s", err, myErr.Body())
	}

	// The device id is returned immediately but it won't show up in GET requests
	// until the order is approved and provisioning is finished

	// Check if all devices were created
	if len(bareMetalDeviceResponse.Devices) != len(devicesToCreate) {
		return fmt.Errorf(
			"number of devices created should match the number of devices requested: %d != %d",
			len(bareMetalDeviceResponse.Devices), len(devicesToCreate))
	}

	// Assign device and order ids to each requested device
	for i, device := range bareMetalDeviceResponse.Devices {
		devicesToCreate[i]["device_id"] = int(device.DeviceId)
		devicesToCreate[i]["order_id"] = int(device.OrderId)
	}

	// Wait for for order to start provisioning
	orderId := bareMetalDeviceResponse.Devices[0].OrderId
	if _, err := waitForOrder(BareMetalDeviceTimeout, hv, orderId); err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		if strings.Contains(fmt.Sprint(err), "'cancelled'") {
			return fmt.Errorf("your deployment (order %d) has been 'cancelled'. Please contact Hivelocity support if you believe this is a mistake.\n\n %s",
				orderId, myErr.Body())
		}
		return fmt.Errorf("error provisioning order %d. The Hivelocity team will investigate:\n\n%s\n\n %s",
			orderId, err, myErr.Body())
	}

	// Wait for devices to finish provisioning
	_, err = waitForDevices(BareMetalDeviceTimeout, hv, orderId, bareMetalDeviceResponse.Devices)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return fmt.Errorf("error finding devices for order %d. The Hivelocity team will investigate:\n\n%s\n\n %s",
			orderId, err, myErr.Body())
	}

	return nil
}

func fieldGet(fieldName string, d *schema.ResourceData, deviceKey string) interface{} {
	key := fmt.Sprintf("%s%s", deviceKey, fieldName)
	return d.Get(key)
}

func fieldHasChange(fieldName string, d *schema.ResourceData, deviceKey string) bool {
	key := fmt.Sprintf("%s%s", deviceKey, fieldName)
	return d.HasChange(key)
}

// Update/Reload a device
func updateDevice(hv *Client, d *schema.ResourceData, deviceKey string, skipReload bool) error {
	deviceId := int32(fieldGet("device_id", d, deviceKey).(int))

	payload := swagger.BareMetalDeviceUpdate{
		Hostname:       fieldGet("hostname", d, deviceKey).(string),
		OsName:         fieldGet("os_name", d, deviceKey).(string),
		Script:         fieldGet("script", d, deviceKey).(string),
		PublicSshKeyId: int32(fieldGet("public_ssh_key_id", d, deviceKey).(int)),
		Tags:           getTags(d, deviceKey),
	}

	reloadRequired := false

	if fieldHasChange("hostname", d, deviceKey) {
		reloadRequired = true
	}

	if fieldHasChange("os_name", d, deviceKey) {
		reloadRequired = true
	}

	if fieldHasChange("public_ssh_key_id", d, deviceKey) {
		reloadRequired = true
	}

	if fieldHasChange("script", d, deviceKey) {
		reloadRequired = true
	}

	// Used when updating just the tags after device creation
	if skipReload {
		reloadRequired = false
	}

	// If a reload is required, it's necessary to turn the device off first
	if reloadRequired {
		devicePower, _, err := hv.client.DeviceApi.GetPowerResource(hv.auth, int32(deviceId), nil)
		if err != nil {
			myErr, _ := err.(swagger.GenericSwaggerError)
			return fmt.Errorf("GET /device/%s/power failed! (%s)\n\n %s", fmt.Sprint(deviceId), err, myErr.Body())
		}

		if devicePower.PowerStatus == "ON" {
			if _, _, err := hv.client.DeviceApi.PostPowerResource(hv.auth, int32(deviceId), "shutdown", nil); err != nil {
				myErr, _ := err.(swagger.GenericSwaggerError)
				return fmt.Errorf("POST /device/%s/power failed! (%s)\n\n %s", fmt.Sprint(deviceId), err, myErr.Body())
			}

			// Power status will transition to PENDING, then OFF
			if _, err := waitForDevicePowerOff(d, hv, int32(deviceId)); err != nil {
				return fmt.Errorf("error powering off device %s. The Hivelocity team will investigate:\n\n%s", fmt.Sprint(deviceId), err)
			}
		}
	}

	if _, _, err := hv.client.BareMetalDevicesApi.PutBareMetalDeviceIdResource(hv.auth, int32(deviceId), payload, nil); err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return fmt.Errorf("PUT /bare-metal-devices/%d failed! (%s)\n\n %s", deviceId, err, myErr.Body())
	}

	if reloadRequired {
		if _, err := waitForDeviceReload(d, hv, int32(deviceId)); err != nil {
			return fmt.Errorf("error reloading device %s. The Hivelocity team will investigate:\n\n%s", fmt.Sprint(deviceId), err)
		}
	}

	return nil
}

// Delete a device
func deleteDevice(hv *Client, deviceId int) error {
	httpResponse, err := hv.client.BareMetalDevicesApi.DeleteBareMetalDeviceIdResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			log.Printf("[WARN] Device (%d) not found", deviceId)
		} else {
			myErr, _ := err.(swagger.GenericSwaggerError)
			return fmt.Errorf(
				"DELETE /bare-metal-devices/%d failed! (%s)\n\n %s", deviceId, err, myErr.Body(),
			)
		}
	}

	return nil
}

// Query API for each Device in the Order Group and return a map indexed by each device's id
func getOrderGroupDevices(hv *Client, orderGroup swagger.OrderGroup) (map[int]swagger.BareMetalDevice, error) {
	orderGroupDevices := make(map[int]swagger.BareMetalDevice)

	for _, deviceId := range orderGroup.DeviceIds {
		deviceResponse, httpResponse, err := hv.client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(hv.auth, int32(deviceId), nil)
		if err == nil {
			orderGroupDevices[int(deviceId)] = deviceResponse
		} else {
			if httpResponse.StatusCode == 404 {
				log.Printf("[WARN] Device (%d) is no longer in OrderGroup (%d)", deviceId, orderGroup.Id)
			} else {
				myErr, _ := err.(swagger.GenericSwaggerError)
				return nil, fmt.Errorf("GET /bare-metal-devices/%d failed! (%s)\n\n %s", deviceId, err, myErr.Body())
			}
		}
	}

	return orderGroupDevices, nil
}
