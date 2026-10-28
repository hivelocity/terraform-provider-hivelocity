package hivelocity

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

func resourceNullRoute() *schema.Resource {
	return &schema.Resource{
		Description: "Null-routes (blackholes) one of your IP addresses at the network edge, for example " +
			"to stop an incoming DDoS. Destroying the resource removes the null route. The resource ID is " +
			"the IP address.",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		CreateContext: resourceNullRouteCreate,
		ReadContext:   resourceNullRouteRead,
		DeleteContext: resourceNullRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ip": {
				Description: "IP address to null-route.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"comment": {
				Description: "Free-form note explaining why the address is null-routed.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"created":   computedString("When the null route was created."),
			"client_id": computedInt("Owning client ID."),
			"admin_id":  computedInt("ID of the admin who created the route, if created by Hivelocity staff."),
		},
	}
}

func resourceNullRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)
	ip := d.Get("ip").(string)

	log.Printf("[INFO] Null-routing %s", ip)
	payload := swagger.NullRouteRequest{Ip: ip, Comment: d.Get("comment").(string)}
	task, _, err := hv.client.NetworkApi.PostNullRoutesResource(hv.auth, payload, nil)
	if err != nil {
		return apiErrorDiag(err, "POST /network/null-route")
	}

	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ip)
	return resourceNullRouteRead(ctx, d, m)
}

func resourceNullRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	route, resp, err := hv.client.NetworkApi.GetNullRoutesIdResource(hv.auth, d.Id(), nil)
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "GET /network/null-route/%s", d.Id())
	}

	return diag.FromErr(setResourceData(d, map[string]interface{}{
		"ip":        route.Ip,
		"comment":   route.Comment,
		"created":   timeString(route.Created),
		"client_id": int(route.ClientId),
		"admin_id":  int(route.AdminId),
	}))
}

func resourceNullRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hv, _ := m.(*Client)

	log.Printf("[INFO] Removing null route for %s", d.Id())
	task, resp, err := hv.client.NetworkApi.DeleteNullRoutesIdResource(hv.auth, d.Id(), nil)
	if err != nil {
		if isNotFound(resp) {
			d.SetId("")
			return nil
		}
		return apiErrorDiag(err, "DELETE /network/null-route/%s", d.Id())
	}

	if _, err := hv.waitForNetworkTask(ctx, task.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
