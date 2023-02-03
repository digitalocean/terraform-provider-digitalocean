package acceptance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestResourceInstanceState is a resource.TestCheckFunc to test if a resource is
// present in the instance state.
func TestResourceInstanceState(name string, check func(*terraform.InstanceState) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := s.RootModule()
		if rs, ok := m.Resources[name]; ok {
			is := rs.Primary
			if is == nil {
				return fmt.Errorf("No primary instance: %s", name)
			}

			return check(is)
		} else {
			return fmt.Errorf("Not found: %s", name)
		}

	}
}
