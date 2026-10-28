package hivelocity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVlan() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up one of your VLANs, for example by `vlan_tag` or `facility_code`.",
		ReadContext: dataSourceVlanRead,
		Schema: map[string]*schema.Schema{
			"filter":        dataSourceFiltersSchema(),
			"first":         dataSourceFilterFirstSchema(),
			"vlan_id":       computedInt("Unique ID of the VLAN."),
			"vlan_tag":      computedInt("VLAN tag on the switch."),
			"facility_code": computedString("Facility of the VLAN."),
			"type":          computedString("`private` or `public`."),
			"port_ids":      computedIntList("Ports and bonds in the VLAN."),
			"ip_ids":        computedIntList("IP assignments routed to the VLAN."),
			"automated":     computedBool("True when the VLAN can be automated via the API."),
			"q_in_q":        computedBool("True when the VLAN is in Q-in-Q mode."),
		},
	}
}

func dataSourceVlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	vlans, _, err := hv.client.VLANApi.GetVlanResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /vlan")
	}

	return readFilteredList(d, vlans, dataSourceVlan().Schema, "vlan_id", nil)
}
