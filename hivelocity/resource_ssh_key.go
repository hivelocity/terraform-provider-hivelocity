package hivelocity

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		CreateContext: resourceSSHKeyCreate,
		ReadContext:   resourceSSHKeyRead,
		UpdateContext: resourceSSHKeyUpdate,
		DeleteContext: resourceSSHKeyDelete,
		Schema: map[string]*schema.Schema{
			"ssh_key_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSSHKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := swagger.SshKey{
		Name:      d.Get("name").(string),
		PublicKey: d.Get("public_key").(string),
	}

	sshKey, _, err := hv.client.SshKeyApi.PostSshKeyResource(hv.auth, payload, nil)
	if err != nil {
		d.SetId("")
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("POST /ssh_key failed! (%s)\n\n %s", err, myErr.Body())
	}

	d.SetId(fmt.Sprint(sshKey.SshKeyId))

	return resourceSSHKeyRead(ctx, d, m)
}

func resourceSSHKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hv, _ := m.(*Client)

	SSHKeyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	SSHKeyResponse, _, err := hv.client.SshKeyApi.GetSshKeyIdResource(hv.auth, int32(SSHKeyID), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("ssh_key_id", SSHKeyResponse.SshKeyId)
	d.Set("name", SSHKeyResponse.Name)
	d.Set("public_key", SSHKeyResponse.PublicKey)

	return diags
}

func resourceSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	SSHKeyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload := swagger.SshKeyUpdate{}

	if d.HasChange("name") {
		payload.Name = d.Get("name").(string)
	}

	if d.HasChange("public_key") {
		payload.PublicKey = d.Get("public_key").(string)
	}

	_, _, err = hv.client.SshKeyApi.PutSshKeyIdResource(hv.auth, int32(SSHKeyID), payload, nil)
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("PUT /ssh_key/%s failed! (%s)\n\n %s", fmt.Sprint(SSHKeyID), err, myErr.Body())
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
	return resourceSSHKeyRead(ctx, d, m)
}

func resourceSSHKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hv, _ := m.(*Client)

	SSHKeyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Check ssh key exists still, if not mark as already destroyed.
	_, _, err = hv.client.SshKeyApi.GetSshKeyIdResource(hv.auth, int32(SSHKeyID), nil)
	if err != nil {
		d.SetId("")
		return diags
	}

	_, err = hv.client.SshKeyApi.DeleteSshKeyIdResource(hv.auth, int32(SSHKeyID))
	if err != nil {
		myErr := err.(swagger.GenericSwaggerError)
		return diag.Errorf("DELETE /ssh_key/%s failed! (%s)\n\n %s", fmt.Sprint(SSHKeyID), err, myErr.Body())
	}

	d.SetId("")

	return diags
}
