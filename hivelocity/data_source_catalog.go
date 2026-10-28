package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

// --- product operating systems -------------------------------------------

func dataSourceProductOperatingSystems() *schema.Resource {
	return &schema.Resource{
		Description: "Lists the operating systems that can be installed on a product. Use the `name` values " +
			"for `os_name` on `hivelocity_bare_metal_device` and `hivelocity_vps`.",
		ReadContext: dataSourceProductOperatingSystemsRead,
		Schema: map[string]*schema.Schema{
			"product_id": {
				Description: "ID of the product.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"operating_systems": {
				Description: "Available operating systems.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":            computedInt("Option ID of the operating system."),
						"name":          computedString("Name to use as `os_name`."),
						"monthly_price": computedFloat("Monthly price of the OS, if any."),
						"currency":      computedString("Currency of the price."),
						"upgrade_name":  computedString("Upgrade group the option belongs to."),
						"tags":          computedStringList("Compatibility tags."),
						"cloud_init":    computedBool("True when the OS supports cloud-init scripts."),
						"post_install":  computedBool("True when the OS supports post-install scripts."),
						"ssh_key":       computedBool("True when the OS supports SSH key injection."),
					},
				},
			},
		},
	}
}

func flattenOption(o swagger.OptionDump) map[string]interface{} {
	item := map[string]interface{}{
		"id":            int(o.Id),
		"name":          o.Name,
		"monthly_price": float64(o.MonthlyPrice),
		"currency":      o.Currency,
		"upgrade_name":  o.UpgradeName,
		"tags":          flattenStringList(o.Tags),
	}
	if o.OperatingSystemFeatures != nil {
		item["cloud_init"] = o.OperatingSystemFeatures.CloudInit
		item["post_install"] = o.OperatingSystemFeatures.PostInstall
		item["ssh_key"] = o.OperatingSystemFeatures.SshKey
	}
	return item
}

func dataSourceProductOperatingSystemsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	productId := d.Get("product_id").(int)

	options, _, err := hv.client.ProductApi.GetProductOperatingSystemsResource(hv.auth, int32(productId), nil)
	if err != nil {
		return apiErrorDiag(err, "GET /product/%d/operating-systems", productId)
	}

	list := make([]interface{}, 0, len(options))
	for _, o := range options {
		list = append(list, flattenOption(o))
	}
	if err := d.Set("operating_systems", list); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprint(productId))
	return nil
}

// --- product options -----------------------------------------------------

func dataSourceProductOptions() *schema.Resource {
	return &schema.Resource{
		Description: "Lists the configurable options (upgrades, add-ons, operating systems) of a product, " +
			"grouped by category.",
		ReadContext: dataSourceProductOptionsRead,
		Schema: map[string]*schema.Schema{
			"product_id": {
				Description: "ID of the product.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"options": {
				Description: "Available options.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category":      computedString("Option category, for example `Operating System`."),
						"id":            computedInt("Option ID."),
						"name":          computedString("Option name."),
						"monthly_price": computedFloat("Monthly price of the option."),
						"currency":      computedString("Currency of the price."),
						"upgrade_name":  computedString("Upgrade group the option belongs to."),
						"tags":          computedStringList("Compatibility tags."),
					},
				},
			},
		},
	}
}

func dataSourceProductOptionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	productId := d.Get("product_id").(int)

	// The response is an object keyed by category whose values are option
	// lists; the generated model for it is empty, so decode by hand.
	var byCategory map[string][]swagger.OptionDump
	if _, err := hv.rawJSON(ctx, "GET", fmt.Sprintf("/product/%d/options", productId), nil, &byCategory); err != nil {
		return diag.FromErr(err)
	}

	categories := make([]string, 0, len(byCategory))
	for c := range byCategory {
		categories = append(categories, c)
	}
	sort.Strings(categories)

	list := make([]interface{}, 0)
	for _, c := range categories {
		for _, o := range byCategory[c] {
			item := flattenOption(o)
			delete(item, "cloud_init")
			delete(item, "post_install")
			delete(item, "ssh_key")
			item["category"] = c
			list = append(list, item)
		}
	}
	if err := d.Set("options", list); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprint(productId))
	return nil
}

// --- services ------------------------------------------------------------

func serviceAttributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service_id":    computedInt("Unique ID of the service."),
		"service_name":  computedString("Name of the service."),
		"service_type":  computedString("Type of the service."),
		"type_code":     computedString("Service type code (see `hivelocity_service_types`)."),
		"status":        computedString("Status of the service."),
		"product_id":    computedInt("Product ID of the service."),
		"product_name":  computedString("Product name of the service."),
		"order_id":      computedInt("Order the service was created from."),
		"period":        computedString("Billing period."),
		"facility_name": computedString("Facility of the service."),
		"ip_address":    computedString("Primary IP address, if any."),
		"service_cost":  computedFloat("Cost of the service."),
		"billed_price":  computedFloat("Billed price."),
		"start_date":    computedInt("Start date as a Unix timestamp."),
		"end_date":      computedInt("End date as a Unix timestamp."),
		"renew_date":    computedInt("Next renewal date as a Unix timestamp."),
		"contract_term": computedInt("Contract term in months."),
		"auto_bill":     computedBool("True when the service is billed automatically."),
	}
}

func dataSourceService() *schema.Resource {
	s := serviceAttributes()
	s["service_id"] = &schema.Schema{
		Description: "ID of the service.",
		Type:        schema.TypeInt,
		Required:    true,
	}
	return &schema.Resource{
		Description: "Returns details of one account service by ID.",
		ReadContext: dataSourceServiceRead,
		Schema:      s,
	}
}

func dataSourceServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	serviceId := d.Get("service_id").(int)

	service, _, err := hv.client.ServiceApi.GetServiceIdResource(hv.auth, int32(serviceId), nil)
	if err != nil {
		return apiErrorDiag(err, "GET /service/%d", serviceId)
	}

	item, err := toGenericMap(service)
	if err != nil {
		return diag.FromErr(err)
	}
	return setGenericItem(d, item, serviceAttributes(), "service_id")
}

func dataSourceServices() *schema.Resource {
	s := serviceAttributes()
	s["filter"] = dataSourceFiltersSchema()
	s["first"] = dataSourceFilterFirstSchema()
	s["status_filter"] = &schema.Schema{
		Description: "Only list services with this status: all, active, deactivated, pending, suspended or cancelled. Defaults to `all`.",
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "all",
	}
	s["type_code_filter"] = &schema.Schema{
		Description: "Only list services of this type code (see `hivelocity_service_types`).",
		Type:        schema.TypeString,
		Optional:    true,
	}
	return &schema.Resource{
		Description: "Finds one of your account services by filtering the service list.",
		ReadContext: dataSourceServicesRead,
		Schema:      s,
	}
}

func dataSourceServicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	opts := &swagger.ServiceApiGetServiceResourceOpts{}
	if v, ok := d.GetOk("type_code_filter"); ok {
		opts.TypeCode = optional.NewString(v.(string))
	}
	services, _, err := hv.client.ServiceApi.GetServiceResource(hv.auth, d.Get("status_filter").(string), opts)
	if err != nil {
		return apiErrorDiag(err, "GET /service")
	}

	return readFilteredList(d, services, serviceAttributes(), "service_id", nil)
}

func dataSourceServiceTypes() *schema.Resource {
	return &schema.Resource{
		Description: "Lists the service type codes used by `hivelocity_services`.",
		ReadContext: dataSourceServiceTypesRead,
		Schema: map[string]*schema.Schema{
			"service_types": {
				Description: "Available service types.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": computedString("Type code."),
						"name": computedString("Type name."),
					},
				},
			},
		},
	}
}

func dataSourceServiceTypesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	types, _, err := hv.client.ServiceApi.GetServiceTypeResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /service/types")
	}

	list := make([]interface{}, 0, len(types))
	for _, t := range types {
		list = append(list, map[string]interface{}{"code": t.Code, "name": t.Name})
	}
	if err := d.Set("service_types", list); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("service-types")
	return nil
}

// --- orders --------------------------------------------------------------

func orderAttributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"order_id": computedInt("Unique ID of the order."),
		"status":   computedString("Order status, for example `complete`."),
		"owner":    computedString("Owner of the order."),
		"total":    computedFloat("Order total."),
		"info":     computedStringMap("Additional order information as reported by the API."),
	}
}

func flattenOrder(o swagger.OrderDump) map[string]interface{} {
	var info interface{}
	if o.Info != nil {
		info = *o.Info
	}
	return map[string]interface{}{
		"order_id": int(o.OrderId),
		"status":   o.Status,
		"owner":    o.Owner,
		"total":    float64(o.Total),
		"info":     stringMapFromInterface(info),
	}
}

func dataSourceOrder() *schema.Resource {
	s := orderAttributes()
	s["order_id"] = &schema.Schema{
		Description: "ID of the order.",
		Type:        schema.TypeInt,
		Required:    true,
	}
	return &schema.Resource{
		Description: "Returns the status of one order by ID.",
		ReadContext: dataSourceOrderRead,
		Schema:      s,
	}
}

func dataSourceOrderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	orderId := d.Get("order_id").(int)

	order, _, err := hv.client.OrderApi.GetOrderIdResource(hv.auth, int32(orderId), nil)
	if err != nil {
		return apiErrorDiag(err, "GET /order/%d", orderId)
	}

	if err := setResourceData(d, flattenOrder(order)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprint(orderId))
	return nil
}

func dataSourceOrders() *schema.Resource {
	s := orderAttributes()
	s["filter"] = dataSourceFiltersSchema()
	s["first"] = dataSourceFilterFirstSchema()
	return &schema.Resource{
		Description: "Finds one of your orders by filtering the order list, for example on `status`.",
		ReadContext: dataSourceOrdersRead,
		Schema:      s,
	}
}

func dataSourceOrdersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	orders, _, err := hv.client.OrderApi.GetOrderResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /order")
	}

	list := make([]map[string]interface{}, 0, len(orders))
	for _, o := range orders {
		list = append(list, flattenOrder(o))
	}
	if len(list) == 0 {
		return diag.Errorf("no orders found")
	}
	item, err := doFiltering(d, list, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if item == nil {
		return diag.Errorf("no order matched the given filters")
	}
	return setGenericItem(d, item, orderAttributes(), "order_id")
}

// --- maintenance events --------------------------------------------------

func dataSourceMaintenanceEvents() *schema.Resource {
	return &schema.Resource{
		Description: "Lists scheduled maintenance events affecting your devices.",
		ReadContext: dataSourceMaintenanceEventsRead,
		Schema: map[string]*schema.Schema{
			"events": {
				Description: "Maintenance events.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maintenance_event_id": computedInt("Unique ID of the event."),
						"title":                computedString("Title."),
						"description":          computedString("Description."),
						"start":                computedString("Start time (UTC)."),
						"end":                  computedString("End time (UTC)."),
						"device_ids":           computedIntList("Devices affected by the event."),
					},
				},
			},
		},
	}
}

func dataSourceMaintenanceEventsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	events, _, err := hv.client.MaintenanceApi.GetMaintenanceEventClientResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /maintenance")
	}

	list := make([]interface{}, 0, len(events))
	for _, e := range events {
		list = append(list, map[string]interface{}{
			"maintenance_event_id": int(e.MaintenanceEventId),
			"title":                e.Title,
			"description":          e.Description,
			"start":                timeString(e.Start),
			"end":                  timeString(e.End),
			"device_ids":           flattenInt32List(e.DeviceIds),
		})
	}
	if err := d.Set("events", list); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("maintenance-events")
	return nil
}

// --- device hardware -----------------------------------------------------

func dataSourceDeviceHardware() *schema.Resource {
	return &schema.Resource{
		Description: "Returns the hardware inventory (drives, memory, RAID) of a device.",
		ReadContext: dataSourceDeviceHardwareRead,
		Schema: map[string]*schema.Schema{
			"device_id": {
				Description: "ID of the device.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"drive_count":            computedInt("Number of physical drives."),
			"drive_types":            computedStringList("Distinct drive types present (NVMe, SSD, HDD)."),
			"in_raid_array":          computedBool("True when drives are members of a RAID virtual disk."),
			"memory_slots_populated": computedInt("Number of populated RAM slots."),
			"drives": {
				Description: "Physical drives.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description":   computedString("Drive description."),
						"product":       computedString("Drive product name."),
						"type":          computedString("Drive type."),
						"size":          computedString("Drive size."),
						"unit_capacity": computedString("Capacity of a single unit."),
						"is_nvme":       computedBool("True for NVMe drives."),
						"is_ssd":        computedBool("True for SSDs."),
						"serials":       computedStringList("Serial numbers."),
					},
				},
			},
			"memory_sticks": {
				Description: "Installed memory modules.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description":  computedString("Module description."),
						"capacity":     computedString("Capacity."),
						"capacity_str": computedString("Human readable capacity."),
					},
				},
			},
			"virtual_disks": {
				Description: "RAID virtual disks.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"raid_type":     computedString("RAID level."),
						"size":          computedString("Size."),
						"unit_capacity": computedString("Capacity of a single unit."),
					},
				},
			},
		},
	}
}

func dataSourceDeviceHardwareRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	deviceId := d.Get("device_id").(int)

	hw, _, err := hv.client.DeviceApi.GetDeviceHardwareResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		return apiErrorDiag(err, "GET /device/%d/hardware", deviceId)
	}

	drives := make([]interface{}, 0, len(hw.Drives))
	for _, dr := range hw.Drives {
		drives = append(drives, map[string]interface{}{
			"description":   dr.Description,
			"product":       dr.Product,
			"type":          dr.Type_,
			"size":          dr.Size,
			"unit_capacity": dr.UnitCapacity,
			"is_nvme":       dr.IsNvme,
			"is_ssd":        dr.IsSsd,
			"serials":       flattenStringList(dr.Serials),
		})
	}
	sticks := make([]interface{}, 0, len(hw.MemorySticks))
	for _, s := range hw.MemorySticks {
		sticks = append(sticks, map[string]interface{}{
			"description":  s.Description,
			"capacity":     s.Capacity,
			"capacity_str": s.CapacityStr,
		})
	}
	vdisks := make([]interface{}, 0, len(hw.VirtualDisks))
	for _, v := range hw.VirtualDisks {
		vdisks = append(vdisks, map[string]interface{}{
			"raid_type":     v.RaidType,
			"size":          v.Size,
			"unit_capacity": v.UnitCapacity,
		})
	}

	if err := setResourceData(d, map[string]interface{}{
		"drive_count":            int(hw.DriveCount),
		"drive_types":            flattenStringList(hw.DriveTypes),
		"in_raid_array":          hw.InRaidArray,
		"memory_slots_populated": int(hw.MemorySlotsPopulated),
		"drives":                 drives,
		"memory_sticks":          sticks,
		"virtual_disks":          vdisks,
	}); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprint(deviceId))
	return nil
}

// --- device IPMI ---------------------------------------------------------

func dataSourceDeviceIpmi() *schema.Resource {
	return &schema.Resource{
		Description: "Returns IPMI controller information and current sensor readings of a device.",
		ReadContext: dataSourceDeviceIpmiRead,
		Schema: map[string]*schema.Schema{
			"device_id": {
				Description: "ID of the device.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"info": computedStringMap("IPMI controller details (firmware revision, manufacturer, ...)."),
			"sensors": {
				Description: "Sensor readings.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sensor_id": computedString("Sensor ID."),
						"name":      computedString("Sensor name."),
						"group":     computedString("Sensor group."),
						"reading":   computedFloat("Current reading."),
						"units":     computedString("Units of the reading."),
						"status":    computedBool("True when the sensor is within thresholds."),
					},
				},
			},
		},
	}
}

func dataSourceDeviceIpmiRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	deviceId := d.Get("device_id").(int)

	ipmi, _, err := hv.client.DeviceApi.GetIpmiInfoIdResource(hv.auth, int32(deviceId), nil)
	if err != nil {
		return apiErrorDiag(err, "GET /device/%d/ipmi", deviceId)
	}

	info := map[string]string{}
	if ipmi.Info != nil {
		data, _ := json.Marshal(ipmi.Info)
		var generic map[string]interface{}
		_ = json.Unmarshal(data, &generic)
		info = stringMapFromInterface(generic)
	}

	sensors := make([]interface{}, 0, len(ipmi.Sensors))
	for _, s := range ipmi.Sensors {
		sensors = append(sensors, map[string]interface{}{
			"sensor_id": s.SensorId,
			"name":      s.Name,
			"group":     s.Group,
			"reading":   float64(s.Reading),
			"units":     s.Units,
			"status":    s.Status,
		})
	}

	if err := setResourceData(d, map[string]interface{}{
		"info":    info,
		"sensors": sensors,
	}); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprint(deviceId))
	return nil
}
