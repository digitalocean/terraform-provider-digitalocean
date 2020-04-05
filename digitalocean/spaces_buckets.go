package digitalocean

import "strings"

func spacesBucketSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Bucket name",
		},
		"urn": {
			Type:        schema.TypeString,
			Description: "the uniform resource name for the bucket",
		},
		"region": {
			Type:        schema.TypeString,
			Description: "Bucket region",
			Default:     "nyc3",
			StateFunc: func(val interface{}) string {
				// DO API V2 region slug is always lowercase
				return strings.ToLower(val.(string))
			},
		},
		"acl": {
			Type:        schema.TypeString,
			Description: "Canned ACL applied on bucket creation",
		},
		"bucket_domain_name": {
			Type:        schema.TypeString,
			Description: "The FQDN of the bucket",
		},
	}
}
