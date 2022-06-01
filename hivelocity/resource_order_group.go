package hivelocity

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceOrderGroup() *schema.Resource {
	resourceBareMetalDevice := resourceBareMetalDevice(false)

	return &schema.Resource{
		Description: "Order Group resource",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(BareMetalDeviceTimeout),
		},
		CreateContext: resourceOrderGroupCreate,
		ReadContext:   resourceOrderGroupRead,
		UpdateContext: resourceOrderGroupUpdate,
		DeleteContext: resourceOrderGroupDelete,
		CustomizeDiff: resourceOrderGroupDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of this Order Group",
				Type:        schema.TypeString,
				Required:    true,
			},
			"same_rack": {
				Description: "Devices will be placed in the same rack",
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
			},
			"bare_metal_device": {
				Description: "A list of Bare Metal Devices to deploy",
				Type:        schema.TypeList,
				Elem:        resourceBareMetalDevice,
				Required:    true,
				ConfigMode:  schema.SchemaConfigModeBlock,
			},
		},
	}
}

// Create the Order Group and all Bare Metal Devices within it
func resourceOrderGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	// Create the Order Group
	orderGroupPayload := swagger.OrderGroupCreate{
		Name:     d.Get("name").(string),
		SameRack: d.Get("same_rack").(bool),
	}
	orderGroup, _, err := hv.client.OrderGroupsApi.PostOrderGroupResource(hv.auth, orderGroupPayload, nil)
	if err != nil {
		d.SetId("")
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("POST /order-groups failed! (%s)\n\n %s", err, myErr.Body())
	}
	d.SetId(fmt.Sprint(orderGroup.Id))

	// Create all devices
	devicesToCreate := make([]map[string]interface{}, d.Get("bare_metal_device.#").(int))
	for i, device := range d.Get("bare_metal_device").([]interface{}) {
		devicesToCreate[i] = device.(map[string]interface{})
	}
	err = createDevices(hv, orderGroup.Id, devicesToCreate)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set state to the new values
	d.Set("bare_metal_device", devicesToCreate)

	// Merge any errors or warnings from Create and Read functions
	return resourceOrderGroupRead(ctx, d, m)
}

// Fetch status of Order Group and all associated Bare Metal Devices
func resourceOrderGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	orderGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Fetch status for the Order Group
	orderGroupResponse, httpResponse, err := hv.client.OrderGroupsApi.GetOrderGroupIdResource(hv.auth, int32(orderGroupId), nil)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			log.Printf("[WARN] Order Group (%d) not found", orderGroupId)
			d.SetId("") // Mark for removal from state
			return nil
		}
		myErr, _ := err.(swagger.GenericSwaggerError)
		return diag.Errorf("GET /order-groups/%d failed! (%s)\n\n %s", orderGroupId, err, myErr.Body())
	}

	// Fetch status for each Bare Metal Device and store in a map indexed by device id
	orderGroupDevices, err := getOrderGroupDevices(hv, orderGroupResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	var updatedDevicesList []map[string]interface{}
	stateDevicesList := make([]map[string]interface{}, d.Get("bare_metal_device.#").(int))
	for i, device := range d.Get("bare_metal_device").([]interface{}) {
		stateDevicesList[i] = device.(map[string]interface{})
	}

	// Try to keep the order of devices the same, so that the diff changes are minimized
	for _, stateDevice := range stateDevicesList {
		deviceId := stateDevice["device_id"].(int)
		if deviceResponse, ok := orderGroupDevices[deviceId]; ok {
			// Update fields and copy device to the new list
			stateDevice := buildDeviceForState(deviceResponse)
			updatedDevicesList = append(updatedDevicesList, stateDevice)
			// Remove device from orderGroupDevices
			delete(orderGroupDevices, deviceId)
		} else {
			log.Printf("[WARN] Device (%d) found in state but not in the API", deviceId)
		}
	}

	// Iterate over orderGroupDevices and if there's any device left add them to state list
	for deviceId, deviceResponse := range orderGroupDevices {
		log.Printf("[WARN] Device (%d) found in the API but not in state", deviceId)
		stateDevice := buildDeviceForState(deviceResponse)
		updatedDevicesList = append(updatedDevicesList, stateDevice)
	}

	d.Set("name", orderGroupResponse.Name)
	d.Set("same_rack", orderGroupResponse.SameRack)
	// The whole list of devices needs to be set at once
	d.Set("bare_metal_device", updatedDevicesList)

	return nil
}

// Detect changes to the Order Group and its Devices and update them
func resourceOrderGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	hv, _ := m.(*Client)

	orderGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Update Order Group
	if d.HasChange("name") {
		payload := swagger.OrderGroupUpdate{Name: d.Get("name").(string)}
		if _, _, err := hv.client.OrderGroupsApi.PutOrderGroupIdResource(hv.auth, int32(orderGroupId), payload, nil); err != nil {
			myErr, _ := err.(swagger.GenericSwaggerError)
			return diag.Errorf("PUT /order-groups/%d failed! (%s)\n\n %s", orderGroupId, err, myErr.Body())
		}
	}

	// Update Bare Metal Devices
	old, new := d.GetChange("bare_metal_device")
	oldDevices := make([]map[string]interface{}, len(old.([]interface{})))
	for i, device := range old.([]interface{}) {
		oldDevices[i] = device.(map[string]interface{})
	}
	currentDevices := make([]map[string]interface{}, len(new.([]interface{})))
	for i, device := range new.([]interface{}) {
		currentDevices[i] = device.(map[string]interface{})
	}

	// Delete devices that were present in the old config but are absent in the new
	for _, device := range oldDevices {
		deviceId := device["device_id"].(int)

		if !deviceInList(deviceId, currentDevices) {
			if err := deleteDevice(hv, deviceId); err == nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Device deleted (%d)", deviceId),
				})
			} else {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	}

	// Map key holds initial device position in the slice
	var devicesToCreate []map[string]interface{}

	// Create/Update devices
	for i, device := range currentDevices {
		if deviceId, ok := device["device_id"]; ok && deviceId.(int) > 0 {
			// Update/Reload existing device
			deviceKey := fmt.Sprintf("bare_metal_device.%d.", i)
			deviceId := int32(fieldGet("device_id", d, deviceKey).(int))
			updatePayload := getBareMetalUpdatePayloadFromGroup(d, deviceKey)

			if err := hv._updateDevice(deviceId, updatePayload, d); err == nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Device updated (%d)", deviceId),
				})
			} else {
				diags = append(diags, diag.FromErr(err)...)
			}
		} else {
			// New device, will be batch created later
			devicesToCreate = append(devicesToCreate, device)
		}
	}

	if len(devicesToCreate) > 0 {
		// Create all new devices and set them in the same position they were found
		err := createDevices(hv, int32(orderGroupId), devicesToCreate)
		if err == nil {
			// Warn user about newly provisioned devices
			for _, device := range devicesToCreate {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Device created (%d)", device["device_id"]),
				})
			}

			// Set state to the new values
			d.Set("bare_metal_device", currentDevices)
		} else {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Merge list of warnings and errors
	return append(diags, resourceOrderGroupRead(ctx, d, m)...)
}

// Delete all Bare Metal Devices and then the Order Group
func resourceOrderGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	// Delete all Bare Metal Devices
	bareMetalDevices := d.Get("bare_metal_device").([]interface{})
	for _, bareMetalDevice := range bareMetalDevices {
		device := bareMetalDevice.(map[string]interface{})
		deviceId := device["device_id"].(int)
		httpResponse, err := hv.client.BareMetalDevicesApi.DeleteBareMetalDeviceIdResource(hv.auth, int32(deviceId), nil)
		if err != nil {
			if httpResponse != nil && httpResponse.StatusCode == 404 {
				log.Printf("[WARN] Device (%d) not found", deviceId)
			} else {
				myErr, _ := err.(swagger.GenericSwaggerError)
				return diag.Errorf("DELETE /bare-metal-devices/%d failed! (%s)\n\n %s", deviceId, err, myErr.Body())
			}
		}
	}

	// Delete Order Group
	orderGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	httpResponse, err := hv.client.OrderGroupsApi.DeleteOrderGroupIdResource(hv.auth, int32(orderGroupId))
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			log.Printf("[WARN] Order Group (%d) not found", orderGroupId)
		} else {
			myErr, _ := err.(swagger.GenericSwaggerError)
			return diag.Errorf("DELETE /order-groups/%d failed! (%s)\n\n %s", orderGroupId, err, myErr.Body())
		}
	}

	d.SetId("")

	return nil
}

// Detect if product or location was changed in any device and force the re-creation of the Order Group
func resourceOrderGroupDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	for i := range d.Get("bare_metal_device").([]interface{}) {
		keyDeviceId := fmt.Sprintf("bare_metal_device.%d.device_id", i)
		keyLocation := fmt.Sprintf("bare_metal_device.%d.location_name", i)
		keyProduct := fmt.Sprintf("bare_metal_device.%d.product_id", i)

		if d.Get(keyDeviceId).(int) > 0 {
			if d.HasChange(keyLocation) {
				d.ForceNew(keyLocation)
			}
			if d.HasChange(keyProduct) {
				d.ForceNew(keyProduct)
			}
		}
	}

	return nil
}
