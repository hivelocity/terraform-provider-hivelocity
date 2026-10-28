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

func resourceWebhook() *schema.Resource {
	return &schema.Resource{
		Description: "Registers a webhook that is called when an account event fires. The list of events is " +
			"available from the `hivelocity_webhook_events` data source.",
		CreateContext: resourceWebhookCreate,
		ReadContext:   resourceWebhookRead,
		UpdateContext: resourceWebhookUpdate,
		DeleteContext: resourceWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Display name of the webhook.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"event": {
				Description: "Name of the event that triggers the webhook.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"url": {
				Description: "URL that receives the event payload.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"headers": {
				Description: "Extra HTTP headers sent with the request, for example an authorization token.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"webhook_id": computedInt("Unique ID of the webhook."),
		},
	}
}

func webhookHeaders(d *schema.ResourceData) *interface{} {
	var headers interface{} = expandStringMap(d.Get("headers"))
	return &headers
}

func resourceWebhookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	payload := swagger.WebhookCreate{
		Name:    d.Get("name").(string),
		Event:   d.Get("event").(string),
		Url:     d.Get("url").(string),
		Headers: webhookHeaders(d),
	}
	webhook, _, err := hv.client.WebhookApi.PostWebhookResource(hv.auth, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /webhooks")
	}

	d.SetId(fmt.Sprint(webhook.WebhookId))
	log.Printf("[INFO] Created webhook %s", d.Id())

	return resourceWebhookRead(ctx, d, m)
}

func resourceWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	webhookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	webhook, resp, err := hv.client.WebhookApi.GetWebhookIdResource(hv.auth, int32(webhookId), nil)
	if err != nil {
		if isNotFound(resp) {
			log.Printf("[WARN] Webhook %d not found, removing from state", webhookId)
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /webhooks/%d", webhookId)
	}

	var headers interface{}
	if webhook.Headers != nil {
		headers = *webhook.Headers
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"name":       webhook.Name,
		"event":      webhook.Event,
		"url":        webhook.Url,
		"headers":    stringMapFromInterface(headers),
		"webhook_id": int(webhook.WebhookId),
	}))
}

func resourceWebhookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	webhookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload := swagger.WebhookUpdate{
		Name:    d.Get("name").(string),
		Event:   d.Get("event").(string),
		Url:     d.Get("url").(string),
		Headers: webhookHeaders(d),
	}
	if _, _, err := hv.client.WebhookApi.PutWebhookIdResource(hv.auth, int32(webhookId), payload, nil); err != nil {
		return apiErrorDiag(err, "PUT /webhooks/%d", webhookId)
	}

	return resourceWebhookRead(ctx, d, m)
}

func resourceWebhookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	webhookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting webhook %d", webhookId)
	if resp, err := hv.client.WebhookApi.DeleteWebhookIdResource(hv.auth, int32(webhookId)); err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /webhooks/%d", webhookId)
	}

	d.SetId("")
	return nil
}

func dataSourceWebhookEvents() *schema.Resource {
	return &schema.Resource{
		Description: "Lists the events a `hivelocity_webhook` can subscribe to.",
		ReadContext: dataSourceWebhookEventsRead,
		Schema: map[string]*schema.Schema{
			"events": {
				Description: "Available webhook events.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":        computedString("Event name to use in `hivelocity_webhook.event`."),
						"description": computedString("What the event means."),
					},
				},
			},
		},
	}
}

func dataSourceWebhookEventsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	events, _, err := hv.client.WebhookApi.GetWebhookEventResource(hv.auth, nil)
	if err != nil {
		return apiErrorDiag(err, "GET /webhooks/events")
	}

	list := make([]interface{}, 0, len(events))
	for _, e := range events {
		list = append(list, map[string]interface{}{"name": e.Name, "description": e.Desc})
	}
	if err := d.Set("events", list); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("webhook-events")
	return nil
}
