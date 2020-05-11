package digitalocean

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testResourceInstanceState(name string, check func(*terraform.InstanceState) error) resource.TestCheckFunc {
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
