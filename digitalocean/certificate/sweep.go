package certificate

import (
	"context"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("digitalocean_certificate", &resource.Sweeper{
		Name: "digitalocean_certificate",
		F:    sweepCertificate,
	})

}

func sweepCertificate(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	certs, _, err := client.Certificates.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, c := range certs {
		if strings.HasPrefix(c.Name, "certificate-") {
			log.Printf("Destroying certificate %s", c.Name)

			if _, err := client.Certificates.Delete(context.Background(), c.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
