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

func TestAccHivelocityDNSDomainBasic(t *testing.T) {
	rString := acctest.RandString(6) + ".com"
	newDomain := acctest.RandString(6) + ".com"
	name := "hivelocity_dns_domain.my-domain"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHivelocityDNSDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHivelocityDNSDomainBase(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rString),
				),
			},
			{
				Config: testAccHivelocityDNSDomainBase(newDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", newDomain),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "0",
				ExpectError:       regexp.MustCompile("Cannot import non-existent remote object"),
			},
		},
	})
}

func testAccCheckHivelocityDNSDomainDestroy(s *terraform.State) error {
	time.Sleep(1 * time.Second)
	hv := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hivelocity_dns_domain" {
			continue
		}

		domainId, _ := strconv.Atoi(rs.Primary.ID)
		_, _, err := hv.client.DomainsApi.GetDomainIdResource(hv.auth, int32(domainId), nil)
		if err == nil {
			return fmt.Errorf("domain still exists : %s", rs.Primary.ID)
		}

	}
	return nil
}

func testAccHivelocityDNSDomainBase(domain string) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf(`
		resource "hivelocity_dns_domain" "my-domain" {
			name = "%s"
		}`, domain)
}
