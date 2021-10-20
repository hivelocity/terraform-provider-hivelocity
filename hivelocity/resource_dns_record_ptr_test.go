package hivelocity

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHivelocityDNSRecordPTRBasic(t *testing.T) {
	rDomain := acctest.RandString(6) + ".com"
	name := "hivelocity_dns_record_ptr.my-record"
	ip := os.Getenv("HIVELOCITY_TEST_DEVICE_IP")
	if ip == "" {
		t.Fatal("HIVELOCITY_TEST_DEVICE_IP must be set for this acceptance test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHivelocityDNSRecordPtrBase(rDomain, rDomain, ip, "3600"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "address", ip),
					resource.TestCheckResourceAttr(name, "name", rDomain),
					resource.TestCheckResourceAttr(name, "ttl", "3600"),
				),
			},
			{
				Config: testAccHivelocityDNSRecordPtrBase(rDomain, "new."+rDomain, ip, "900"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "address", ip),
					resource.TestCheckResourceAttr(name, "name", "new."+rDomain),
					resource.TestCheckResourceAttr(name, "ttl", "900"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccHivelocityDNSRecordPtrBase(domain, name, address, ttl string) string {
	time.Sleep(1 * time.Second)
	return testAccHivelocityDNSRecordABase(domain, name, address, ttl) + fmt.Sprintf(`
		resource "hivelocity_dns_record_ptr" "my-record" {
			address    = "%s"
			name       = "%s"
			ttl        = %s

			depends_on = [
				hivelocity_dns_record_a.my-record
			]
		}`, address, name, ttl)
}
