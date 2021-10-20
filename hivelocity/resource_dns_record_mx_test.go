package hivelocity

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccHivelocityDNSRecordMXBasic(t *testing.T) {
	rDomain := acctest.RandString(6) + ".com"
	rSub := acctest.RandString(4) + "." + rDomain
	name := "hivelocity_dns_record_mx.my-record"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHivelocityDNSRecordMxDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHivelocityDNSRecordMxBase(rDomain, rSub, "mail1.domain.com", "10", "3600"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rSub),
					resource.TestCheckResourceAttr(name, "exchange", "mail1.domain.com"),
					resource.TestCheckResourceAttr(name, "preference", "10"),
					resource.TestCheckResourceAttr(name, "ttl", "3600"),
				),
			},
			{
				Config: testAccHivelocityDNSRecordMxBase(rDomain, rSub, "mail2.domain.com", "20", "900"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rSub),
					resource.TestCheckResourceAttr(name, "exchange", "mail2.domain.com"),
					resource.TestCheckResourceAttr(name, "preference", "20"),
					resource.TestCheckResourceAttr(name, "ttl", "900"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					for _, rs := range s.RootModule().Resources {
						if rs.Type == "hivelocity_dns_record_mx" {
							return rs.Primary.Attributes["domain_id"] + ":" + rs.Primary.ID, nil
						}
					}
					return "", fmt.Errorf("record resource not found in state")
				},
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "0:0",
				ExpectError:       regexp.MustCompile("could not import record: 404 NOT FOUND"),
			},
		},
	})
}

func TestAccHivelocityDNSRecordMXAtSymbol(t *testing.T) {
	rDomain := acctest.RandString(6) + ".com"
	name := "hivelocity_dns_record_mx.my-record"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHivelocityDNSRecordMxDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHivelocityDNSRecordMxBase(rDomain, "@", "mail1.domain.com", "10", "3600"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rDomain),
					resource.TestCheckResourceAttr(name, "exchange", "mail1.domain.com"),
					resource.TestCheckResourceAttr(name, "preference", "10"),
					resource.TestCheckResourceAttr(name, "ttl", "3600"),
				),
			},
			{
				Config: testAccHivelocityDNSRecordMxBase(rDomain, rDomain, "mail1.domain.com", "10", "3600"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rDomain),
					resource.TestCheckResourceAttr(name, "exchange", "mail1.domain.com"),
					resource.TestCheckResourceAttr(name, "preference", "10"),
					resource.TestCheckResourceAttr(name, "ttl", "3600"),
				),
			},
		},
	})
}

func testAccCheckHivelocityDNSRecordMxDestroy(s *terraform.State) error {
	time.Sleep(1 * time.Second)
	hv := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "hivelocity_dns_domain" {
			domainId, _ := strconv.Atoi(rs.Primary.ID)
			if _, _, err := hv.client.DomainsApi.GetDomainIdResource(hv.auth, int32(domainId), nil); err == nil {
				return fmt.Errorf("domain still exists : %s", rs.Primary.ID)
			}
		} else if rs.Type == "hivelocity_dns_record_mx" {
			recordId, _ := strconv.Atoi(rs.Primary.ID)
			domainId, _ := strconv.Atoi(rs.Primary.Attributes["domain_id"])
			if _, _, err := hv.client.DomainsApi.GetMxRecordIdResource(hv.auth, int32(domainId), int32(recordId), nil); err == nil {
				return fmt.Errorf("record still exists : %s", rs.Primary.ID)
			}
		}

	}
	return nil
}

func testAccHivelocityDNSRecordMxBase(domain, name, exchange, preference, ttl string) string {
	time.Sleep(1 * time.Second)
	return testAccHivelocityDNSDomainBase(domain) + fmt.Sprintf(`
		resource "hivelocity_dns_record_mx" "my-record" {
			domain_id = hivelocity_dns_domain.my-domain.id
			name       = "%s"
			exchange   = "%s"
			preference = %s
			ttl        = %s
		}`, name, exchange, preference, ttl)
}
