package hivelocity

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEffectiveIgnition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEffectiveIgnitionRead,
		Schema: map[string]*schema.Schema{
			"device_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"contents": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEffectiveIgnitionRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	hv, _ := m.(*Client)

	deviceId := int32(d.Get("device_id").(int))

	contents, err := hv.getEffectiveIgnition(deviceId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("contents", contents)
	d.SetId(fmt.Sprintf("%d", deviceId))
	return nil
}

func (hv *Client) getEffectiveIgnition(deviceId int32) (string, error) {
	ign, _, err := hv.client.DeviceApi.GetEffectiveIgnitionIdResource(hv.auth, deviceId, nil)
	if err != nil {
		return "", formatSwaggerError(err, "/device/%d/ignition", deviceId)
	} else {
		return ign.Contents, nil
	}
}
