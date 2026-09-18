package droplet

import (
	"testing"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
)

func TestFindIPv4AddrByType_NilSafe(t *testing.T) {
	assert.Equal(t, "", FindIPv4AddrByType(nil, "public"))
	assert.Equal(t, "", FindIPv4AddrByType(&godo.Droplet{}, "public"))
}

func TestFindIPv6AddrByType_NilSafe(t *testing.T) {
	assert.Equal(t, "", FindIPv6AddrByType(nil, "public"))
	assert.Equal(t, "", FindIPv6AddrByType(&godo.Droplet{}, "public"))
}

func TestFindIPv4AddrByType_PublicAndPrivate(t *testing.T) {
	droplet := &godo.Droplet{
		Networks: &godo.Networks{
			V4: []godo.NetworkV4{
				{Type: "public", IPAddress: "203.0.113.10"},
				{Type: "private", IPAddress: "10.0.0.5"},
			},
		},
	}

	assert.Equal(t, "203.0.113.10", FindIPv4AddrByType(droplet, "public"))
	assert.Equal(t, "10.0.0.5", FindIPv4AddrByType(droplet, "private"))
	assert.Equal(t, "", FindIPv4AddrByType(droplet, "unknown"))
}

func TestFindIPv4AddrByType_InvalidIPIgnored(t *testing.T) {
	droplet := &godo.Droplet{
		Networks: &godo.Networks{
			V4: []godo.NetworkV4{
				{Type: "public", IPAddress: "not-an-ip"},
			},
		},
	}

	assert.Equal(t, "", FindIPv4AddrByType(droplet, "public"))
}

func TestExpectPublicIPv4FromCreateOpts(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		setPublicNet   bool
		publicNetValue bool
		expectPublic   bool
	}{
		{name: "unset defaults to public", expectPublic: true},
		{name: "explicit true", setPublicNet: true, publicNetValue: true, expectPublic: true},
		{name: "explicit false", setPublicNet: true, publicNetValue: false, expectPublic: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := ResourceDigitalOceanDroplet().TestResourceData()
			d.Set("image", "ubuntu-22-04-x64")
			d.Set("name", "web")
			d.Set("region", "nyc3")
			d.Set("size", "s-1vcpu-1gb")
			if tc.setPublicNet {
				d.Set("public_networking", tc.publicNetValue)
			}

			opts, err := expandDropletCreateRequest(d)
			assert.NoError(t, err)

			expectPublicIPv4 := opts.PublicNetworking == nil || *opts.PublicNetworking
			assert.Equal(t, tc.expectPublic, expectPublicIPv4)
		})
	}
}
