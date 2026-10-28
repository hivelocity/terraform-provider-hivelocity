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

func resourceContact() *schema.Resource {
	return &schema.Resource{
		Description: "Manages an additional contact (user) on your account. Grant it permissions with " +
			"`hivelocity_contact_permissions`.",
		CreateContext: resourceContactCreate,
		ReadContext:   resourceContactRead,
		UpdateContext: resourceContactUpdate,
		DeleteContext: resourceContactDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"full_name": {
				Description: "Full name of the contact.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"email": {
				Description: "Email address, used as the login.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"phone": {
				Description: "Phone number.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"description": {
				Description: "Free-form description of the contact's role.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"active": {
				Description: "Whether the contact can log in. Defaults to true.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"password": {
				Description: "Initial password. Only used on creation; later changes are applied through " +
					"the contact password endpoint and require `current_password`.",
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"current_password": {
				Description: "Current password, required only when changing `password` of an existing contact.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"contact_id":    computedInt("Unique ID of the contact."),
			"client_id":     computedInt("Owning client ID."),
			"is_client":     computedBool("True for the primary account contact."),
			"referral_code": computedString("Referral code of the contact."),
		},
	}
}

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func resourceContactCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := swagger.ContactCreate{
		FullName:    d.Get("full_name").(string),
		Email:       d.Get("email").(string),
		Phone:       d.Get("phone").(string),
		Description: d.Get("description").(string),
		Active:      boolToInt32(d.Get("active").(bool)),
		Password:    d.Get("password").(string),
	}
	contact, _, err := hv.client.ContactApi.PostContactResource(hv.auth, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /contact")
	}

	d.SetId(fmt.Sprint(contact.ContactId))
	log.Printf("[INFO] Created contact %s", d.Id())

	return resourceContactRead(ctx, d, m)
}

func resourceContactRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	contactId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	contact, resp, err := hv.client.ContactApi.GetContactIdResource(hv.auth, int32(contactId), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] Contact %d not found, removing from state", contactId)
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /contact/%d", contactId)
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"full_name":     contact.FullName,
		"email":         contact.Email,
		"phone":         contact.Phone,
		"description":   contact.Description,
		"active":        contact.Active != 0,
		"contact_id":    int(contact.ContactId),
		"client_id":     int(contact.ClientId),
		"is_client":     contact.IsClient,
		"referral_code": contact.ReferralCode,
	}))
}

func resourceContactUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	contactId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("full_name", "email", "phone", "description", "active") {
		payload := swagger.ContactUpdate{
			FullName:    d.Get("full_name").(string),
			Email:       d.Get("email").(string),
			Phone:       d.Get("phone").(string),
			Description: d.Get("description").(string),
			Active:      boolToInt32(d.Get("active").(bool)),
		}
		if _, _, err := hv.client.ContactApi.PutContactIdResource(hv.auth, int32(contactId), payload, nil); err != nil {
			return apiErrorDiag(err, "PUT /contact/%d", contactId)
		}
	}

	if d.HasChange("password") && d.Get("password").(string) != "" {
		current := d.Get("current_password").(string)
		if current == "" {
			return diag.Errorf("current_password must be set to change the password of contact %d", contactId)
		}
		payload := swagger.Password{Password: current, NewPassword: d.Get("password").(string)}
		if _, err := hv.client.ContactApi.PutPasswordResource(hv.auth, int32(contactId), payload); err != nil {
			return apiErrorDiag(err, "PUT /contact/%d/password", contactId)
		}
	}

	return resourceContactRead(ctx, d, m)
}

func resourceContactDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	contactId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting contact %d", contactId)
	if resp, err := hv.client.ContactApi.DeleteContactIdResource(hv.auth, int32(contactId)); err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /contact/%d", contactId)
	}

	d.SetId("")
	return nil
}

func dataSourceContact() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up one of your account contacts, typically by `email`.",
		ReadContext: dataSourceContactRead,
		Schema: map[string]*schema.Schema{
			"filter":        dataSourceFiltersSchema(),
			"first":         dataSourceFilterFirstSchema(),
			"contact_id":    computedInt("Unique ID of the contact."),
			"client_id":     computedInt("Owning client ID."),
			"full_name":     computedString("Full name."),
			"email":         computedString("Email address."),
			"phone":         computedString("Phone number."),
			"description":   computedString("Description."),
			"active":        computedInt("1 when the contact is active, 0 otherwise."),
			"is_client":     computedBool("True for the primary account contact."),
			"referral_code": computedString("Referral code."),
		},
	}
}

func dataSourceContactRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	contacts, _, err := hv.client.ContactApi.GetContactResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /contact")
	}

	return readFilteredList(d, contacts, dataSourceContact().Schema, "contact_id", nil)
}
