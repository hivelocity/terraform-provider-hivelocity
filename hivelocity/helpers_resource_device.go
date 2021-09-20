package hivelocity

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func getTags(d *schema.ResourceData) ([]string, error) {
	ts := d.Get("tags")
	var sts []string

	switch ts.(type) {
	case []interface{}:
		for _, v := range ts.([]interface{}) {
			sts = append(sts, v.(string))
		}
		if len(sts) == 0 {
			return []string{""}, nil
		} else {
			return sts, nil
		}
	default:
		return nil, fmt.Errorf("garbage in tags: %s", ts)
	}
}

func updateTagsForCreate(hv *Client, deviceId int32, tags []string) (*swagger.BareMetalDevice, error) {
	payload := swagger.BareMetalDeviceUpdate{Tags: tags}
	bm, _, err := hv.client.BareMetalDevicesApi.PutBareMetalDeviceIdResource(hv.auth, int32(deviceId), payload, nil)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return nil, fmt.Errorf("PUT /bare-metal-device/%s failed! (%s)\n\n %s", fmt.Sprint(deviceId), err, myErr.Body())
	}
	return &bm, nil

}

func waitForDevice(d *schema.ResourceData, hv *Client, orderId int32) (interface{}, error) {
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

			for _, device := range devices {
				if device.OrderId == orderId {
					return device, "ok", nil
				}
			}
			return nil, "waiting", nil
		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     30 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 1,
	}
	return waitForDevice.WaitForState()
}

func waitForOrder(d *schema.ResourceData, hv *Client, orderId int32) (interface{}, error) {
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
		Timeout:                   d.Timeout(schema.TimeoutCreate),
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
