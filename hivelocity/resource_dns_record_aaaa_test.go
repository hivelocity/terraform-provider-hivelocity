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

func TestAccHivelocityDNSRecordAAAABasic(t *testing.T) {
	rDomain := acctest.RandString(6) + ".com"
	rSub := acctest.RandString(4) + "." + rDomain
	name := "hivelocity_dns_record_aaaa.my-record"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHivelocityDNSRecordAAAADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHivelocityDNSRecordAAAABase(rDomain, rSub, "fe80::0001", "3600"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rSub),
					resource.TestCheckResourceAttr(name, "address", "fe80::1"),
					resource.TestCheckResourceAttr(name, "ttl", "3600"),
				),
			},
			{
				Config: testAccHivelocityDNSRecordAAAABase(rDomain, rSub, "fe80::0002", "900"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rSub),
					resource.TestCheckResourceAttr(name, "address", "fe80::2"),
					resource.TestCheckResourceAttr(name, "ttl", "900"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					for _, rs := range s.RootModule().Resources {
						if rs.Type == "hivelocity_dns_record_aaaa" {
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

func TestAccHivelocityDNSRecordAAAAAtSymbol(t *testing.T) {
	rDomain := acctest.RandString(6) + ".com"
	name := "hivelocity_dns_record_aaaa.my-record"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHivelocityDNSRecordAAAADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHivelocityDNSRecordAAAABase(rDomain, "@", "fe80::1", "3600"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rDomain),
					resource.TestCheckResourceAttr(name, "address", "fe80::1"),
					resource.TestCheckResourceAttr(name, "ttl", "3600"),
				),
			},
			{
				Config: testAccHivelocityDNSRecordAAAABase(rDomain, rDomain, "fe80::1", "3600"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rDomain),
					resource.TestCheckResourceAttr(name, "address", "fe80::1"),
					resource.TestCheckResourceAttr(name, "ttl", "3600"),
				),
			},
		},
	})
}

func testAccCheckHivelocityDNSRecordAAAADestroy(s *terraform.State) error {
	time.Sleep(1 * time.Second)
	hv := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "hivelocity_dns_domain" {
			domainId, _ := strconv.Atoi(rs.Primary.ID)
			if _, _, err := hv.client.DomainsApi.GetDomainIdResource(hv.auth, int32(domainId), nil); err == nil {
				return fmt.Errorf("domain still exists : %s", rs.Primary.ID)
			}
		} else if rs.Type == "hivelocity_dns_record_aaaa" {
			recordId, _ := strconv.Atoi(rs.Primary.ID)
			domainId, _ := strconv.Atoi(rs.Primary.Attributes["domain_id"])
			if _, _, err := hv.client.DomainsApi.GetAaaaRecordIdResource(hv.auth, int32(domainId), int32(recordId), nil); err == nil {
				return fmt.Errorf("record still exists : %s", rs.Primary.ID)
			}
		}

	}
	return nil
}

func testAccHivelocityDNSRecordAAAABase(domain, name, address, ttl string) string {
	time.Sleep(1 * time.Second)
	return testAccHivelocityDNSDomainBase(domain) + fmt.Sprintf(`
		resource "hivelocity_dns_record_aaaa" "my-record" {
			domain_id = hivelocity_dns_domain.my-domain.id
			name       = "%s"
			address    = "%s"
			ttl        = %s
		}`, name, address, ttl)
}
