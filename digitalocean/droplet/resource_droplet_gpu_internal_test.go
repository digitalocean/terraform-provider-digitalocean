package droplet

import (
	"testing"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
)

func TestSetDropletAttributes_GPUPartitionMode(t *testing.T) {
	d := ResourceDigitalOceanDroplet().TestResourceData()

	droplet := &godo.Droplet{
		ID:               123,
		Name:             "gpu",
		Region:           &godo.Region{Slug: "atl1"},
		Size:             &godo.Size{Slug: "gpu-mi350x8-2304gb"},
		Networks:         &godo.Networks{},
		GPUPartitionMode: godo.GPUPartitionModeDPXNPS2,
	}

	if err := setDropletAttributes(d, droplet); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	assert.Equal(t, "PARTITION_MODE_DPX_NPS2", d.Get("gpu_partition_mode"))
}

func TestSetDropletAttributes_GPUPartitionModePreservedOnEmptyRead(t *testing.T) {
	// The API does not return gpu_partition_mode when reading a Droplet, so a
	// refresh must not clobber the configured value with an empty string.
	d := ResourceDigitalOceanDroplet().TestResourceData()
	if err := d.Set("gpu_partition_mode", godo.GPUPartitionModeDPXNPS2); err != nil {
		t.Fatalf("unexpected error seeding state: %s", err)
	}

	droplet := &godo.Droplet{
		ID:               123,
		Name:             "gpu",
		Region:           &godo.Region{Slug: "atl1"},
		Size:             &godo.Size{Slug: "gpu-mi350x8-2304gb"},
		Networks:         &godo.Networks{},
		GPUPartitionMode: "",
	}

	if err := setDropletAttributes(d, droplet); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	assert.Equal(t, "PARTITION_MODE_DPX_NPS2", d.Get("gpu_partition_mode"))
}

func TestDropletResource_GPUPartitionModeValidation(t *testing.T) {
	validate := ResourceDigitalOceanDroplet().Schema["gpu_partition_mode"].ValidateFunc
	if validate == nil {
		t.Fatal("expected a ValidateFunc on gpu_partition_mode")
	}

	for _, valid := range []string{
		godo.GPUPartitionModeSPXNPS1,
		godo.GPUPartitionModeDPXNPS2,
	} {
		_, errs := validate(valid, "gpu_partition_mode")
		assert.Emptyf(t, errs, "expected %q to be valid", valid)
	}

	_, errs := validate("PARTITION_MODE_BOGUS", "gpu_partition_mode")
	assert.NotEmpty(t, errs, "expected an unknown partition mode to be rejected")
}

func TestExpandDropletCreateRequest_GPUPartitionMode(t *testing.T) {
	d := ResourceDigitalOceanDroplet().TestResourceData()
	d.Set("image", "gpu-amd-base")
	d.Set("name", "gpu")
	d.Set("region", "atl1")
	d.Set("size", "gpu-mi350x8-2304gb")
	d.Set("gpu_partition_mode", godo.GPUPartitionModeDPXNPS2)

	opts, err := expandDropletCreateRequest(d)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	assert.Equal(t, "PARTITION_MODE_DPX_NPS2", opts.GPUPartitionMode)
}

func TestExpandDropletCreateRequest_GPUPartitionModeOmitted(t *testing.T) {
	// When the mode is not configured it must not be sent, so the server applies
	// its default (a full GPU).
	d := ResourceDigitalOceanDroplet().TestResourceData()
	d.Set("image", "ubuntu-22-04-x64")
	d.Set("name", "web")
	d.Set("region", "atl1")
	d.Set("size", "s-1vcpu-1gb")

	opts, err := expandDropletCreateRequest(d)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	assert.Equal(t, "", opts.GPUPartitionMode)
}

func TestFlattenDigitalOceanDroplet_GPUPartitionModeEmptyOnRead(t *testing.T) {
	// The API does not return gpu_partition_mode when reading a Droplet, so the
	// data source currently exposes it as an empty string.
	droplet := godo.Droplet{
		ID:               1,
		Name:             "gpu",
		Region:           &godo.Region{Slug: "atl1"},
		Size:             &godo.Size{Slug: "gpu-mi350x8-2304gb"},
		Image:            &godo.Image{Slug: "gpu-amd-base"},
		Networks:         &godo.Networks{},
		GPUPartitionMode: "",
	}

	flattened, err := flattenDigitalOceanDroplet(droplet, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	assert.Equal(t, "", flattened["gpu_partition_mode"])
}
