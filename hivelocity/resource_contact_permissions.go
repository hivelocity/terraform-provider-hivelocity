package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContactPermissions() *schema.Resource {
	return &schema.Resource{
		Description: "Sets the full list of permissions granted to a contact. Available permission names " +
			"come from the `hivelocity_permissions` data source. Destroying the resource revokes all " +
			"permissions of the contact. The resource ID is the contact ID.",
		CreateContext: resourceContactPermissionsCreate,
		ReadContext:   resourceContactPermissionsRead,
		UpdateContext: resourceContactPermissionsUpdate,
		DeleteContext: resourceContactPermissionsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"contact_id": {
				Description: "ID of the contact.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"permissions": {
				Description: "Permission names granted to the contact.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// permissionNames extracts permission names from the loosely specified
// permission endpoints, which may return a list of strings, a list of objects
// with a `name`, or an object keyed by name.
func permissionNames(body []byte) ([]string, error) {
	var names []string

	var asStrings []string
	if err := json.Unmarshal(body, &asStrings); err == nil {
		names = asStrings
	} else {
		var asObjects []map[string]interface{}
		if err := json.Unmarshal(body, &asObjects); err == nil {
			for _, o := range asObjects {
				if n, ok := o["name"].(string); ok {
					names = append(names, n)
				} else if n, ok := o["permission"].(string); ok {
					names = append(names, n)
				}
			}
		} else {
			var asMap map[string]interface{}
			if err := json.Unmarshal(body, &asMap); err != nil {
				return nil, fmt.Errorf("unrecognised permission response: %s", truncate(string(body), 200))
			}
			if list, ok := asMap["permissions"]; ok {
				data, _ := json.Marshal(list)
				return permissionNames(data)
			}
			for k := range asMap {
				names = append(names, k)
			}
		}
	}

	sort.Strings(names)
	if names == nil {
		names = []string{}
	}
	return names, nil
}

func (hv *Client) setContactPermissions(ctx context.Context, contactId int, permissions []string) error {
	payload := map[string]interface{}{
		"userId":      contactId,
		"permissions": permissions,
	}
	if _, err := hv.rawJSON(ctx, "POST", "/permission/contact", payload, nil); err != nil {
		return err
	}
	return nil
}

func resourceContactPermissionsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	contactId := d.Get("contact_id").(int)

	log.Printf("[INFO] Setting permissions of contact %d", contactId)
	if err := hv.setContactPermissions(ctx, contactId, expandStringList(d.Get("permissions"))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(contactId))
	return resourceContactPermissionsRead(ctx, d, m)
}

func resourceContactPermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	contactId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body, resp, err := hv.rawRequest(ctx, "GET", fmt.Sprintf("/permission/contact/%d", contactId), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] Contact %d not found, removing permissions from state", contactId)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	names, err := permissionNames(body)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"contact_id":  contactId,
		"permissions": names,
	}))
}

func resourceContactPermissionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	contactId := d.Get("contact_id").(int)

	if d.HasChange("permissions") {
		if err := hv.setContactPermissions(ctx, contactId, expandStringList(d.Get("permissions"))); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceContactPermissionsRead(ctx, d, m)
}

func resourceContactPermissionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	contactId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Revoking all permissions of contact %d", contactId)
	if err := hv.setContactPermissions(ctx, contactId, []string{}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func dataSourcePermissions() *schema.Resource {
	return &schema.Resource{
		Description: "Lists every permission name that can be granted to a contact.",
		ReadContext: dataSourcePermissionsRead,
		Schema: map[string]*schema.Schema{
			"permissions": computedStringList("Available permission names."),
		},
	}
}

func dataSourcePermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	body, _, err := hv.rawRequest(ctx, "GET", "/permission/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	names, err := permissionNames(body)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("permissions", names); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("permissions")
	return nil
}

func dataSourceUserPermissions() *schema.Resource {
	return &schema.Resource{
		Description: "Lists the permissions of the user owning the API key.",
		ReadContext: dataSourceUserPermissionsRead,
		Schema: map[string]*schema.Schema{
			"permissions": computedStringList("Permission names of the current user."),
		},
	}
}

func dataSourceUserPermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	body, _, err := hv.rawRequest(ctx, "GET", "/permission/user", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	names, err := permissionNames(body)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("permissions", names); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("user-permissions")
	return nil
}
