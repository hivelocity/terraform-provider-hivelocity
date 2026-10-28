package hivelocity

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Creates a throwaway zone plus A, AAAA and MX records in it.
func TestAccHivelocityDns_basic(t *testing.T) {
	zone := fmt.Sprintf("%s.example", acctest.RandomWithPrefix("tf-acc"))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsConfig(zone, "192.0.2.10", 3600),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hivelocity_dns_domain.test", "name", zone),
					resource.TestCheckResourceAttrSet("hivelocity_dns_domain.test", "domain_id"),
					resource.TestCheckResourceAttr("hivelocity_dns_record_a.test", "id", zone+"/www"),
					resource.TestCheckTypeSetElemAttr("hivelocity_dns_record_a.test", "addresses.*", "192.0.2.10"),
					resource.TestCheckResourceAttr("hivelocity_dns_record_aaaa.test", "address", "2001:db8::10"),
					resource.TestCheckResourceAttr("hivelocity_dns_record_mx.test", "preference", "10"),
					resource.TestCheckResourceAttrSet("hivelocity_dns_record_mx.test", "record_id"),
				),
			},
			{
				Config: testAccDnsConfig(zone, "192.0.2.11", 300),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("hivelocity_dns_record_a.test", "addresses.*", "192.0.2.11"),
					resource.TestCheckResourceAttr("hivelocity_dns_record_a.test", "ttl", "300"),
					resource.TestCheckResourceAttr("hivelocity_dns_record_mx.test", "ttl", "300"),
				),
			},
			{
				ResourceName:      "hivelocity_dns_domain.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "hivelocity_dns_record_a.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsConfig(zone, address string, ttl int) string {
	return fmt.Sprintf(`
resource "hivelocity_dns_domain" "test" {
  name = %[1]q
}

resource "hivelocity_dns_record_a" "test" {
  domain    = hivelocity_dns_domain.test.name
  name      = "www"
  addresses = [%[2]q]
  ttl       = %[3]d
}

resource "hivelocity_dns_record_aaaa" "test" {
  domain_id = hivelocity_dns_domain.test.domain_id
  name      = "www"
  address   = "2001:db8::10"
  ttl       = %[3]d
}

resource "hivelocity_dns_record_mx" "test" {
  domain_id  = hivelocity_dns_domain.test.domain_id
  name       = "@"
  exchange   = "mail.%[1]s"
  preference = 10
  ttl        = %[3]d
}
`, zone, address, ttl)
}
