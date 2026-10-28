package hivelocity

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

type DeviceMetadata struct {
	isReloadFlag bool
	spsStatus    string
}

func (m *DeviceMetadata) isReload() bool {
	return m.isReloadFlag || m.spsStatus == "Reloading"
}

// Return device tags
func getTags(d *schema.ResourceData, deviceKey string) []string {
	key := fmt.Sprintf("%stags", deviceKey)
	var tags []string

	for _, v := range d.Get(key).([]interface{}) {
		tags = append(tags, v.(string))
	}

	return tags
}

func waitForDevices(ctx context.Context, timeout time.Duration, hv *Client, orderId int32, newDevices []swagger.BareMetalDevice) (interface{}, error) {
	waitForDevice := &retry.StateChangeConf{
		Pending: []string{
			"waiting",
		},
		Target: []string{
			"ok",
		},
		Refresh: func() (interface{}, string, error) {
			devices, _, err := hv.client.BareMetalDevicesApi.GetBareMetalDeviceResource(hv.auth, nil)
			if err != nil {
				return 0, "", formatSwaggerError(err, "GET /bare-metal-devices/")
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
	return waitForDevice.WaitForStateContext(ctx)
}

func waitForOrder(ctx context.Context, timeout time.Duration, hv *Client, orderId int32) (string, error) {
	waitForOrder := &retry.StateChangeConf{
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
				return "", "", formatSwaggerError(err, "/order/%d", orderId)
			}
			return resp, resp.Status, nil
		},
		Timeout:                   timeout,
		Delay:                     5 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 1,
	}
	status, err := waitForOrder.WaitForStateContext(ctx)
	if statusStr, ok := status.(string); ok {
		return statusStr, err
	}
	return "", err
}

func waitForDevicePowerOff(ctx context.Context, d *schema.ResourceData, hv *Client, deviceId int32) (interface{}, error) {
	waitForDevice := &retry.StateChangeConf{
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
	return waitForDevice.WaitForStateContext(ctx)
}

func (hv *Client) waitForDevice(ctx context.Context, deviceId int32, timeout time.Duration) (string, error) {
	waitForDevice := &retry.StateChangeConf{
		Pending: []string{"waiting"},
		Target:  []string{"ok"},
		Refresh: func() (interface{}, string, error) {
			metadata, err := hv.getDeviceMetadata(deviceId)
			if err != nil {
				return nil, "", err
			}
			switch metadata.spsStatus {
			case "Failed":
				return metadata, "failed", nil
			case "InUse":
				return metadata, "ok", nil
			default:
				return metadata, "waiting", nil
			}
		},
		Timeout:                   timeout,
		Delay:                     30 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 1,
		NotFoundChecks:            360, // 1h timeout / 10s delay between requests
	}
	ret, err := waitForDevice.WaitForStateContext(ctx)
	if err != nil {
		return "", err
	}
	metadata := ret.(*DeviceMetadata)

	return metadata.spsStatus, nil
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
	stateDevice["period"] = device.Period
	stateDevice["custom_ipxe_script_url"] = device.CustomIPXEScriptURL
	stateDevice["custom_ipxe_script_contents"] = device.CustomIPXEScriptContents
	stateDevice["app_name"] = device.AppName

	return stateDevice
}

// Creates a list of devices
func createDevices(ctx context.Context, hv *Client, orderGroupId int32, devicesToCreate []map[string]interface{}) error {
	// Use the batch endpoint to create the Bare Metal Devices all at once
	var bareMetalDevicesPayload []swagger.BareMetalDeviceCreate

	for _, device := range devicesToCreate {
		var tags []string
		for _, tag := range device["tags"].([]interface{}) {
			tags = append(tags, tag.(string))
		}

		bareMetalDeviceCreate := swagger.BareMetalDeviceCreate{
			ProductId:                int32(device["product_id"].(int)),
			Hostname:                 device["hostname"].(string),
			OsName:                   device["os_name"].(string),
			LocationName:             device["location_name"].(string),
			Script:                   device["script"].(string),
			Period:                   device["period"].(string),
			PublicSshKeyId:           int32(device["public_ssh_key_id"].(int)),
			ForceDeviceId:            int32(device["force_device_id"].(int)),
			IgnitionId:               int32(device["ignition_id"].(int)),
			CustomIPXEScriptURL:      device["custom_ipxe_script_url"].(string),
			CustomIPXEScriptContents: device["custom_ipxe_script_contents"].(string),
			ForceCustomIPXEDownload:  device["force_custom_ipxe_download"].(bool),
			Tags:                     tags,
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
	if status, err := waitForOrder(ctx, BareMetalDeviceTimeout, hv, orderId); err != nil {
		if status == "cancelled" {
			return fmt.Errorf("your deployment (order %d) has been 'cancelled'. Please contact Hivelocity"+
				" support if you believe this is a mistake.\n\n %s",
				orderId, err)
		}
		return err
	}

	// Wait for devices to finish provisioning
	_, err = waitForDevices(ctx, BareMetalDeviceTimeout, hv, orderId, bareMetalDeviceResponse.Devices)
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

func isPowerOffError(err error) bool {
	if formattedErr, ok := err.(FormattedSwaggerError); ok {
		err = formattedErr.origError
	}

	if err, ok := err.(swagger.GenericSwaggerError); ok {
		contents := string(err.Body())
		return strings.Contains(strings.ToLower(contents), "power off device")
	}
	return false
}

func (hv *Client) doPowerOff(ctx context.Context, deviceId int32, d *schema.ResourceData) error {
	devicePower, _, err := hv.client.DeviceApi.GetPowerResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		myErr, _ := err.(swagger.GenericSwaggerError)
		return fmt.Errorf("GET /device/%s/power failed! (%s)\n\n %s", fmt.Sprint(deviceId), err, myErr.Body())
	}

	if fmt.Sprint(devicePower.PowerStatus) == "ON" {
		if _, _, err := hv.client.DeviceApi.PostPowerResource(hv.auth, int32(deviceId), "shutdown", nil); err != nil {
			myErr, _ := err.(swagger.GenericSwaggerError)
			return fmt.Errorf("POST /device/%s/power failed! (%s)\n\n %s", fmt.Sprint(deviceId), err, myErr.Body())
		}

		// Power status will transition to PENDING, then OFF
		if _, err := waitForDevicePowerOff(ctx, d, hv, int32(deviceId)); err != nil {
			return fmt.Errorf("error powering off device %s. The Hivelocity team will investigate:\n\n%s", fmt.Sprint(deviceId), err)
		}
	}
	return nil
}

type FormattedSwaggerError struct {
	origError error
	msg       string
}

func (e FormattedSwaggerError) Error() string {
	return e.msg
}

// formatSwaggerError takes an error that could be a swagger error and formats appropriately. This is used for
// error handling around direct swagger client calls, so helper functions that may use the swagger client can continue
// to return `error` without juggling `diag.Diagnostics` that may or may not be needed.
func formatSwaggerError(err error, msg string, a ...interface{}) FormattedSwaggerError {
	fmtStr := fmt.Sprintf("%s (%s)", fmt.Sprintf(msg, a...), err)

	if swagError, ok := err.(swagger.GenericSwaggerError); ok {
		fmtStr += "\n" + string(swagError.Body())
	}

	return FormattedSwaggerError{
		origError: err,
		msg:       fmtStr,
	}
}

func isMetadataValueTruthyString(value string) bool {
	switch value {
	case "":
	case "0":
	case "false":
		return false
	}
	return true
}

func (hv *Client) getDevice(deviceId int32) (*swagger.DeviceDump, error) {
	deviceDump, _, err := hv.client.DeviceApi.GetDeviceIdResource(hv.auth, deviceId, nil)
	if err != nil {
		return nil, formatSwaggerError(err, "GET /device/%d", deviceId)
	}
	return &deviceDump, nil
}

func (hv *Client) getDeviceMetadata(deviceId int32) (*DeviceMetadata, error) {
	var m map[string]interface{}

	if deviceDump_, err := hv.getDevice(deviceId); err == nil {
		m = (*deviceDump_.Metadata).(map[string]interface{})
	} else {
		return nil, err
	}

	isReloadFlag := false
	spsStatus := ""

	if reloadStr, ok := m["is_reload"].(string); ok {
		isReloadFlag = isMetadataValueTruthyString(reloadStr)
	}

	if spsStatusStr, ok := m["sps_status"].(string); ok {
		spsStatus = spsStatusStr
	}

	return &DeviceMetadata{
		isReloadFlag: isReloadFlag,
		spsStatus:    spsStatus,
	}, nil
}

func getBareMetalUpdatePayloadFromGroup(d *schema.ResourceData, deviceKey string) swagger.BareMetalDeviceUpdate {
	var ignitionId int32 = 0
	var publicSshKeyId int32 = 0
	var script string

	if value := fieldGet("script", d, deviceKey); value != nil {
		script = value.(string)
	}
	if value := fieldGet("ignition_id", d, deviceKey); value != nil {
		ignitionId = int32(value.(int))
	}
	if value := fieldGet("public_ssh_key_id", d, deviceKey); value != nil {
		publicSshKeyId = int32(value.(int))
	}

	var ipxeURL, ipxeContents string
	if value := fieldGet("custom_ipxe_script_url", d, deviceKey); value != nil {
		ipxeURL = value.(string)
	}
	if value := fieldGet("custom_ipxe_script_contents", d, deviceKey); value != nil {
		ipxeContents = value.(string)
	}

	return swagger.BareMetalDeviceUpdate{
		Hostname:                 fieldGet("hostname", d, deviceKey).(string),
		OsName:                   fieldGet("os_name", d, deviceKey).(string),
		Tags:                     getTags(d, deviceKey),
		Script:                   script,
		IgnitionId:               ignitionId,
		PublicSshKeyId:           publicSshKeyId,
		CustomIPXEScriptURL:      ipxeURL,
		CustomIPXEScriptContents: ipxeContents,
	}
}

func getBareMetalUpdatePayload(d *schema.ResourceData) swagger.BareMetalDeviceUpdate {
	return getBareMetalUpdatePayloadFromGroup(d, "")
}

// VlanPorts describes a VLAN and ports for a specific device
type VlanPorts struct {
	deviceId int32
	vlan     *swagger.Vlan
	ports    map[int32]*swagger.DevicePort
}

func (vp *VlanPorts) getPortIdListWithoutDevicePorts() []int32 {
	newList := make([]int32, 0, len(vp.vlan.PortIds))
	for _, portId := range vp.vlan.PortIds {
		if _, ok := vp.ports[portId]; !ok {
			newList = append(newList, portId)
		}
	}
	return newList
}

// removeVlanPorts removes all ports from designated VLANs in vlanPorts
func (hv *Client) removeVlanPorts(
	ctx context.Context,
	vlanPorts *map[int32]VlanPorts,
	timeout time.Duration,
) error {
	for _, entry := range *vlanPorts {
		newPortIds := entry.getPortIdListWithoutDevicePorts()

		if err := hv.updateVlanMembers(ctx, entry.vlan.VlanId, newPortIds, entry.vlan.IpIds, timeout); err != nil {
			return err
		}
	}
	return nil
}

// restoreVlanPorts restores ports that were removed with removeVlanPorts
func (hv *Client) restoreVlanPorts(
	ctx context.Context,
	vlanPorts *map[int32]VlanPorts,
	timeout time.Duration,
) error {
	for _, entry := range *vlanPorts {
		if err := hv.updateVlanMembers(ctx, entry.vlan.VlanId, entry.vlan.PortIds, entry.vlan.IpIds, timeout); err != nil {
			return err
		}
	}
	return nil
}

// getVlanIdToVlanPortMap returns a list of VLANs and ports that belong to one specific device rather than all devices
// for the client.
func (hv *Client) getVlanIdToVlanPortMap(deviceId int32) (map[int32]VlanPorts, error) {
	vlanIdToPorts := make(map[int32]VlanPorts)

	// Grab all ports for device
	portMap := make(map[int32]*swagger.DevicePort)
	ports, _, err := hv.client.DeviceApi.GetDevicePortResource(hv.auth, deviceId, nil)
	if err != nil {
		return nil, formatSwaggerError(err, "GET /device/%d/ports", deviceId)
	}

	// Map ports by their ID
	for _, port := range ports {
		portMap[port.PortId] = &port
	}

	// Get all vlans for client
	vlans, _, err := hv.client.VLANApi.GetVlanResource(hv.auth, nil)
	if err != nil {
		return nil, formatSwaggerError(err, "GET /vlan")
	}

	// Record ports for every matching port ID
	for _, vlan := range vlans {
		for _, portId := range vlan.PortIds {
			if port, ok := portMap[portId]; ok {
				if _, ok := vlanIdToPorts[vlan.VlanId]; !ok {
					vlanIdToPorts[vlan.VlanId] = VlanPorts{
						deviceId: deviceId,
						vlan:     &vlan,
						ports:    make(map[int32]*swagger.DevicePort),
					}
				}

				vlanIdToPorts[vlan.VlanId].ports[portId] = port
			}
		}
	}

	return vlanIdToPorts, nil
}

// Update/Reload a device, making sure to remove device ports from any VLAN and turning off device if necessary
func (hv *Client) _updateDevice(
	ctx context.Context,
	deviceId int32,
	updatePayload swagger.BareMetalDeviceUpdate,
	d *schema.ResourceData,
) error {
	_internalUpdateDevice := func() error {
		if _, _, err := hv.client.BareMetalDevicesApi.PutBareMetalDeviceIdResource(
			hv.auth,
			deviceId,
			updatePayload,
			nil,
		); err != nil {
			return formatSwaggerError(err, "PUT /bare-metal-devices/%d", deviceId)
		}
		return nil
	}
	requiresPowerOff := false

	// Check if device has ports in any VLANs
	vlanPortResults, err := hv.getVlanIdToVlanPortMap(deviceId)
	if err != nil {
		return err
	}

	// Remove device's ports from all assigned vlans
	if len(vlanPortResults) > 0 {
		if err := hv.removeVlanPorts(ctx, &vlanPortResults, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	// Attempt update, if it errors due to power on, then power off first
	if err := _internalUpdateDevice(); err != nil {
		if !isPowerOffError(err) {
			return err
		} else {
			requiresPowerOff = true
		}
	}

	// If a reload is required, it's necessary to turn the device off first
	if requiresPowerOff {
		if err := hv.doPowerOff(ctx, deviceId, d); err != nil {
			return err
		}

		// Try update again
		if err := _internalUpdateDevice(); err != nil {
			return err
		}
	}

	// Query if device is reloading, if so then wait for reload
	metadata, err := hv.getDeviceMetadata(deviceId)
	if err != nil {
		return err
	}

	// If device metadata indicates a reload is taking place, then wait
	if metadata.isReload() {
		if _, err := hv.waitForDevice(ctx, deviceId, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("error reloading device %s. The Hivelocity team will investigate:\n\n%s", fmt.Sprint(deviceId), err)
		}
	}

	// Restore VLAN ports if any
	if len(vlanPortResults) > 0 {
		if err := hv.restoreVlanPorts(ctx, &vlanPortResults, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	return nil
}

// Delete a device
func deleteDevice(hv *Client, deviceId int) error {
	httpResponse, err := hv.client.BareMetalDevicesApi.DeleteBareMetalDeviceIdResource(hv.auth, int32(deviceId))
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
