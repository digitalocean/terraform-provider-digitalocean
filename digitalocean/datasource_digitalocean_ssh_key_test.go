package digitalocean

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"golang.org/x/crypto/ssh"
)

func TestAccDataSourceDigitalOceanSSHKey_Basic(t *testing.T) {
	var key godo.Key
	keyName := fmt.Sprintf("foo-%s", acctest.RandString(10))

	pubKey, err := testAccGenerateDataSourceDigitalOceanSSHKeyPublic()
	if err != nil {
		t.Fatalf("Unable to generate public key: %v", err)
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanSSHKeyConfig_basic, keyName, pubKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanSSHKeyExists("data.digitalocean_ssh_key.foobar", &key),
					resource.TestCheckResourceAttr(
						"data.digitalocean_ssh_key.foobar", "name", keyName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_ssh_key.foobar", "public_key", pubKey),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanSSHKeyExists(n string, key *godo.Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ssh key ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		foundKey, _, err := client.Keys.GetByID(context.Background(), id)

		if err != nil {
			return err
		}

		if foundKey.ID != id {
			return fmt.Errorf("Key not found")
		}

		*key = *foundKey

		return nil
	}
}

func testAccGenerateDataSourceDigitalOceanSSHKeyPublic() (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", fmt.Errorf("Unable to generate key: %v", err)
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", fmt.Errorf("Unable to generate key: %v", err)
	}

	return strings.TrimSpace(string(ssh.MarshalAuthorizedKey(publicKey))), nil
}

const testAccCheckDataSourceDigitalOceanSSHKeyConfig_basic = `
resource "digitalocean_ssh_key" "foo" {
  name = "%s"
  public_key = "%s"
}

data "digitalocean_ssh_key" "foobar" {
  name = "${digitalocean_ssh_key.foo.name}"
}`
