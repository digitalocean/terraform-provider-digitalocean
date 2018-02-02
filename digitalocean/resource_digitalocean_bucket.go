package digitalocean

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/minio/minio-go"
)

func resourceDigitalOceanBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanBucketCreate,
		Read:   resourceDigitalOceanBucketRead,
		Delete: resourceDigitalOceanBucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bucket name",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bucket endpoint",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bucket region",
			},
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Spaces Access Key",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Spaces Secret Key",
			},
		},
	}
}

func resourceDigitalOceanBucketCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := minio.New(
		d.Get("endpoint").(string),
		d.Get("access_key").(string),
		d.Get("secret_key").(string),
		true,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = client.MakeBucket(d.Get("name").(string), d.Get("region").(string))
	if err != nil {
		return fmt.Errorf("Error creating bucket: %s", err)
	}
	log.Println("Bucket created")

	d.SetId(d.Get("name").(string))
	return resourceDigitalOceanBucketRead(d, meta)
}

func resourceDigitalOceanBucketRead(d *schema.ResourceData, meta interface{}) error {
	client, err := minio.New(
		d.Get("endpoint").(string),
		d.Get("access_key").(string),
		d.Get("secret_key").(string),
		true,
	)
	if err != nil {
		log.Fatal(err)
	}

	found, err := client.BucketExists(d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error finding bucket: %s", err)
	}

	if found {
		log.Println("Bucket found.")
	} else {
		log.Println("Bucket not found.")
	}

	d.Set("name", d.Get("name").(string))

	return nil
}

func resourceDigitalOceanBucketDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := minio.New(
		d.Get("endpoint").(string),
		d.Get("access_key").(string),
		d.Get("secret_key").(string),
		true,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = client.RemoveBucket(d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error deleting bucket: %s", err)
	}
	log.Println("Bucket destroyed")

	d.SetId("")
	return nil
}
