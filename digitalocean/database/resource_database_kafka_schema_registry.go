package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDatabaseKafkaSchemaRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseKafkaSchemaRegistryCreate,
		ReadContext:   resourceDigitalOceanDatabaseKafkaSchemaRegistryRead,
		DeleteContext: resourceDigitalOceanDatabaseKafkaSchemaRegistryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanDatabaseKafkaSchemaRegistryImport,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"subject_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schema_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AVRO",
					"JSON",
					"PROTOBUF",
				}, false),
			},
			"schema": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDigitalOceanDatabaseKafkaSchemaRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	schemaRegistry := &godo.DatabaseKafkaSchemaRegistryRequest{
		SubjectName: d.Get("subject_name").(string),
		SchemaType:  d.Get("schema_type").(string),
		Schema:      d.Get("schema").(string),
	}

	log.Printf("[INFO] Creating Kafka Schema Registry for cluster %s", clusterID)

	schemaRegistrySubject, _, err := client.Databases.CreateKafkaSchemaRegistry(ctx, clusterID, schemaRegistry)
	if err != nil {
		return diag.Errorf("Error creating Kafka Schema Registry: %s", err)
	}

	d.SetId(makeKafkaSchemaSubject(clusterID, schemaRegistrySubject.SubjectName))

	return resourceDigitalOceanDatabaseKafkaSchemaRegistryRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseKafkaSchemaRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	subjectName := d.Get("subject_name").(string)

	log.Printf("[INFO] Reading Kafka Schema Registry for cluster %s", clusterID)

	schemaRegistry, resp, err := client.Databases.GetKafkaSchemaRegistry(ctx, clusterID, subjectName)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] DigitalOcean Kafka Schema Registry (%s) not found", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving Kafka Schema Registry: %s", err)
	}

	d.Set("schema", schemaRegistry.Schema)
	d.Set("subject_name", schemaRegistry.SubjectName)
	d.Set("schema_type", schemaRegistry.SchemaType)

	return nil
}

func resourceDigitalOceanDatabaseKafkaSchemaRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	subjectName := d.Get("subject_name").(string)

	log.Printf("[INFO] Deleting Kafka Schema Registry for cluster %s", clusterID)

	_, err := client.Databases.DeleteKafkaSchemaRegistry(ctx, clusterID, subjectName)
	if err != nil {
		return diag.Errorf("Error deleting Kafka Schema Registry: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanDatabaseKafkaSchemaRegistryImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		d.Set("cluster_id", s[0])
		return []*schema.ResourceData{d}, nil
	}

	d.Set("cluster_id", d.Id())
	return []*schema.ResourceData{d}, nil
}

func makeKafkaSchemaSubject(clusterID string, name string) string {
	return fmt.Sprintf("%s/schema-registry/%s", clusterID, name)
}
