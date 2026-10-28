package hivelocity

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Requires the name of a valid webhook event, for example:
//
//	HIVELOCITY_TEST_WEBHOOK_EVENT=device.power
func TestAccHivelocityWebhook_basic(t *testing.T) {
	testAccRequireEnv(t, "HIVELOCITY_TEST_WEBHOOK_EVENT")
	name := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig(name, "https://example.com/hook"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hivelocity_webhook.test", "name", name),
					resource.TestCheckResourceAttr("hivelocity_webhook.test", "url", "https://example.com/hook"),
					resource.TestCheckResourceAttr("hivelocity_webhook.test", "headers.X-Test", "1"),
					resource.TestCheckResourceAttrSet("hivelocity_webhook.test", "webhook_id"),
				),
			},
			{
				Config: testAccWebhookConfig(name, "https://example.com/hook2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hivelocity_webhook.test", "url", "https://example.com/hook2"),
				),
			},
			{
				ResourceName:      "hivelocity_webhook.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccWebhookConfig(name, url string) string {
	return fmt.Sprintf(`
resource "hivelocity_webhook" "test" {
  name  = %q
  event = %q
  url   = %q
  headers = {
    X-Test = "1"
  }
}
`, name, testAccEnv("HIVELOCITY_TEST_WEBHOOK_EVENT"), url)
}
