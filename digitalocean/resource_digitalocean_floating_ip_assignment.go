package digitalocean

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDigitalOceanFloatingIpAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanFloatingIpAssignmentCreate,
		Update: resourceDigitalOceanFloatingIpAssignmentUpdate,
		Read:   resourceDigitalOceanFloatingIpAssignmentRead,
		Delete: resourceDigitalOceanFloatingIpAssignmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
			},

			"droplet_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceDigitalOceanFloatingIpAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	d.SetId(d.Get("ip_address").(string))

	if v, ok := d.GetOk("droplet_id"); ok {

		log.Printf("[INFO] Assigning the Floating IP to the Droplet %d", v.(int))
		action, _, err := client.FloatingIPActions.Assign(context.Background(), d.Id(), v.(int))
		if err != nil {
			return fmt.Errorf(
				"Error Assigning FloatingIP (%s) to the droplet: %s", d.Id(), err)
		}

		_, unassignedErr := waitForFloatingIPReady(d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
		if unassignedErr != nil {
			return fmt.Errorf(
				"Error waiting for FloatingIP (%s) to be Assigned: %s", d.Id(), unassignedErr)
		}
	}

	return resourceDigitalOceanFloatingIpAssignmentRead(d, meta)
}

func resourceDigitalOceanFloatingIpAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	log.Printf("[INFO] Reading the details of the FloatingIP %s", d.Id())
	floatingIp, resp, err := client.FloatingIPs.Get(context.Background(), d.Id())
	if resp.StatusCode != 404 {
		if err != nil {
			return fmt.Errorf("Error retrieving FloatingIP: %s", err)
		}

		if floatingIp.Droplet != nil {
			log.Printf("[INFO] A droplet was detected on the FloatingIP.")
			d.Set("droplet_id", floatingIp.Droplet.ID)
		}
	} else {
		d.SetId("")
	}

	return nil
}

func resourceDigitalOceanFloatingIpAssignmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	if d.HasChange("droplet_id") {
		if v, ok := d.GetOk("droplet_id"); ok {
			log.Printf("[INFO] Assigning the Floating IP %s to the Droplet %d", d.Id(), v.(int))
			action, _, err := client.FloatingIPActions.Assign(context.Background(), d.Id(), v.(int))
			if err != nil {
				return fmt.Errorf(
					"Error Assigning FloatingIP (%s) to the droplet: %s", d.Id(), err)
			}

			_, unassignedErr := waitForFloatingIPReady(d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return fmt.Errorf(
					"Error waiting for FloatingIP (%s) to be Assigned: %s", d.Id(), unassignedErr)
			}
		} else {
			log.Printf("[INFO] Unassigning the Floating IP %s", d.Id())
			action, _, err := client.FloatingIPActions.Unassign(context.Background(), d.Id())
			if err != nil {
				return fmt.Errorf(
					"Error unassigning FloatingIP (%s): %s", d.Id(), err)
			}

			_, unassignedErr := waitForFloatingIPReady(d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return fmt.Errorf(
					"Error waiting for FloatingIP (%s) to be Unassigned: %s", d.Id(), unassignedErr)
			}
		}
	}

	return resourceDigitalOceanFloatingIpAssignmentRead(d, meta)
}

func resourceDigitalOceanFloatingIpAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	if _, ok := d.GetOk("droplet_id"); ok {
		log.Printf("[INFO] Unassigning the Floating IP from the Droplet")
		action, resp, err := client.FloatingIPActions.Unassign(context.Background(), d.Id())
		if resp.StatusCode != 422 {
			if err != nil {
				return fmt.Errorf(
					"Error unassigning FloatingIP (%s) from the droplet: %s", d.Id(), err)
			}

			_, unassignedErr := waitForFloatingIPReady(d, "completed", []string{"new", "in-progress"}, "status", meta, action.ID)
			if unassignedErr != nil {
				return fmt.Errorf(
					"Error waiting for FloatingIP (%s) to be unassigned: %s", d.Id(), unassignedErr)
			}
		} else {
			log.Printf("[DEBUG] Couldn't unassign FloatingIP (%s) from droplet, possibly out of sync: %s", d.Id(), err)
		}
	}

	d.SetId("")
	return nil
}
