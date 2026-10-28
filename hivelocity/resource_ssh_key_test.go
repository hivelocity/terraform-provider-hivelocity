package hivelocity

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccSshPublicKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJl3dJ9r0dHcqj6v5hyJ6W1HV7n0E3Q8xJ3iZ2yqkd3T terraform-acc-test"

func TestAccHivelocitySshKey_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSshKeyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hivelocity_ssh_key.test", "name", name),
					resource.TestCheckResourceAttr("hivelocity_ssh_key.test", "public_key", testAccSshPublicKey),
					resource.TestCheckResourceAttrSet("hivelocity_ssh_key.test", "ssh_key_id"),
				),
			},
			{
				Config: testAccSshKeyConfig(name + "-renamed"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hivelocity_ssh_key.test", "name", name+"-renamed"),
				),
			},
		},
	})
}

func testAccSshKeyConfig(name string) string {
	return fmt.Sprintf(`
resource "hivelocity_ssh_key" "test" {
  name       = %q
  public_key = %q
}
`, name, testAccSshPublicKey)
}
