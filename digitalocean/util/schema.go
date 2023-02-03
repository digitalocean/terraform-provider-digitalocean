package util

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SetResourceDataFromMap sets a *schema.ResourceData from a map.
func SetResourceDataFromMap(d *schema.ResourceData, m map[string]interface{}) error {
	for key, value := range m {
		if err := d.Set(key, value); err != nil {
			return fmt.Errorf("Unable to set `%s` attribute: %s", key, err)
		}
	}
	return nil
}
