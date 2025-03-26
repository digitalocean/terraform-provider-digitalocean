package partnernetworkconnect

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers("digitalocean_partner_attachment", &resource.Sweeper{
		Name: "digitalocean_partner_attachment",
		F:    sweepPartnerAttachment,
	})
}

func sweepPartnerAttachment(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()
	ctx := context.Background()

	opt := &godo.ListOptions{PerPage: 200}
	partnerAttachments, _, err := client.PartnerNetworkConnect.List(ctx, opt)
	if err != nil {
		return err
	}

	for _, p := range partnerAttachments {
		if strings.HasPrefix(p.Name, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Destroying Partner Attachment %s", p.Name)
			resp, err := client.PartnerNetworkConnect.Delete(ctx, p.ID)
			if err != nil {
				if resp.StatusCode == http.StatusForbidden {
					log.Printf("[DEBUG] Skipping Partner Attachment %s; still contains resources", p.Name)
				} else {
					return err
				}
			}
			log.Printf("[DEBUG] Waiting for Partner Attachment (%s) to be deleted", p.Name)
			stateConf := &retry.StateChangeConf{
				Pending:    []string{"DELETING"},
				Target:     []string{http.StatusText(http.StatusNotFound)},
				Refresh:    partnerAttachmentStateRefreshFunc(client, p.ID),
				Timeout:    10 * time.Minute,
				MinTimeout: 2 * time.Second,
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("error waiting for Partner Attachment (%s) to be deleted: %s", p.Name, err)
			}
		}
	}

	return nil
}
