package partnerinterconnectattachment

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
	resource.AddTestSweepers("digitalocean_partner_interconnect_attachment", &resource.Sweeper{
		Name: "digitalocean_partner_interconnect_attachment",
		F:    sweepPartnerInterconnectAttachment,
	})
}

func sweepPartnerInterconnectAttachment(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()
	ctx := context.Background()

	opt := &godo.ListOptions{PerPage: 200}
	partnerInterconnectAttachments, _, err := client.PartnerInterconnectAttachments.List(ctx, opt)
	if err != nil {
		return err
	}

	for _, p := range partnerInterconnectAttachments {
		if strings.HasPrefix(p.Name, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Destroying Partner Interconnect Attachment %s", p.Name)
			resp, err := client.PartnerInterconnectAttachments.Delete(ctx, p.ID)
			if err != nil {
				if resp.StatusCode == http.StatusForbidden {
					log.Printf("[DEBUG] Skipping Partner Interconnect Attachment %s; still contains resources", p.Name)
				} else {
					return err
				}
			}
			log.Printf("[DEBUG] Waiting for Partner Interconnect Attachment (%s) to be deleted", p.Name)
			stateConf := &retry.StateChangeConf{
				Pending:    []string{"DELETING"},
				Target:     []string{http.StatusText(http.StatusNotFound)},
				Refresh:    partnerInterconnectAttachmentStateRefreshFunc(client, p.ID),
				Timeout:    10 * time.Minute,
				MinTimeout: 2 * time.Second,
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("error waiting for Partner Interconnect Attachment (%s) to be deleted: %s", p.Name, err)
			}
		}
	}

	return nil
}
