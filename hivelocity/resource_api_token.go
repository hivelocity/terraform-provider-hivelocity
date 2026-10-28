package hivelocity

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiToken() *schema.Resource {
	return &schema.Resource{
		Description: "Creates an API token for the user owning the provider's API key. The token value is " +
			"only returned on creation and is stored in state (sensitive). Only `name` can be changed " +
			"in place; other changes replace the token.",
		CreateContext: resourceApiTokenCreate,
		ReadContext:   resourceApiTokenRead,
		UpdateContext: resourceApiTokenUpdate,
		DeleteContext: resourceApiTokenDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Display name of the token.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"expires_at": {
				Description: "Expiration of the token in RFC 3339 format (UTC), for example `2027-01-01T00:00:00Z`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"allowed_ips": {
				Description: "IP addresses or CIDRs the token may be used from. Unrestricted when omitted.",
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"allowed_ips_name": {
				Description: "Label for the IP restriction.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"api_token_id": computedString("Unique ID of the token."),
			"token": {
				Description: "The token value. Only available on creation.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"time_created":   computedString("When the token was created."),
			"time_updated":   computedString("When the token was last updated."),
			"time_last_used": computedString("When the token was last used."),
		},
	}
}

// apiToken mirrors the token endpoints' JSON. The generated client types
// timeExpiration as time.Time while the API documents it as a string, so the
// token endpoints are called directly.
type apiToken struct {
	ApiTokenId     string `json:"apiTokenId,omitempty"`
	Name           string `json:"name,omitempty"`
	Token          string `json:"token,omitempty"`
	TimeExpiration string `json:"timeExpiration,omitempty"`
	TimeCreated    string `json:"timeCreated,omitempty"`
	TimeUpdated    string `json:"timeUpdated,omitempty"`
	TimeLastUsed   string `json:"timeLastUsed,omitempty"`
	IpAddresses    *struct {
		Name string   `json:"name,omitempty"`
		Ips  []string `json:"ips"`
	} `json:"ipAddresses,omitempty"`
}

func resourceApiTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := map[string]interface{}{
		"name":           d.Get("name").(string),
		"timeExpiration": d.Get("expires_at").(string),
	}
	if ips := expandStringList(d.Get("allowed_ips")); len(ips) > 0 {
		payload["ipAddresses"] = map[string]interface{}{
			"name": d.Get("allowed_ips_name").(string),
			"ips":  ips,
		}
	}

	var token apiToken
	if _, err := hv.rawJSON(ctx, "POST", "/token/", payload, &token); err != nil {
		return diag.FromErr(err)
	}
	if token.ApiTokenId == "" {
		return diag.Errorf("POST /token/ did not return an apiTokenId")
	}

	d.SetId(token.ApiTokenId)
	if err := d.Set("token", token.Token); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created API token %s", d.Id())

	return resourceApiTokenRead(ctx, d, m)
}

func resourceApiTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	body, _, err := hv.rawRequest(ctx, "GET", "/token/", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	var tokens []apiToken
	if err := json.Unmarshal(body, &tokens); err != nil {
		return diag.Errorf("decoding GET /token/ response: %s", err)
	}

	for _, t := range tokens {
		if t.ApiTokenId != d.Id() {
			continue
		}
		values := map[string]interface{}{
			"api_token_id":   t.ApiTokenId,
			"name":           t.Name,
			"expires_at":     t.TimeExpiration,
			"time_created":   t.TimeCreated,
			"time_updated":   t.TimeUpdated,
			"time_last_used": t.TimeLastUsed,
		}
		if t.IpAddresses != nil {
			values["allowed_ips"] = flattenStringList(t.IpAddresses.Ips)
			values["allowed_ips_name"] = t.IpAddresses.Name
		}
		// The list endpoint only returns the last characters of the token;
		// keep the full value captured at creation.
		return diag.FromErr(setResourceData(d, values))
	}

	log.Printf("[WARN] API token %s not found, removing from state", d.Id())
	d.SetId("")
	return nil
}

func resourceApiTokenUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	if d.HasChange("name") {
		payload := map[string]interface{}{
			"apiTokenId": d.Id(),
			"name":       d.Get("name").(string),
		}
		if _, err := hv.rawJSON(ctx, "PUT", "/token/", payload, nil); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceApiTokenRead(ctx, d, m)
}

func resourceApiTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Deleting API token %s", d.Id())
	payload := map[string]interface{}{"apiTokenId": d.Id()}
	if _, resp, err := hv.rawRequest(ctx, "DELETE", "/token/", payload); err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("deleting API token %s: %w", d.Id(), err))
	}

	d.SetId("")
	return nil
}
