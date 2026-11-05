package size

import (
	"testing"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
)

func TestFlattenGPUInfo_Nil(t *testing.T) {
	// Non-GPU sizes have no gpu_info; it should flatten to an empty block.
	assert.Equal(t, []interface{}{}, flattenGPUInfo(nil))
}

func TestFlattenGPUInfo_WithSupportedPartitionModes(t *testing.T) {
	info := &godo.GPUInfo{
		Count: 8,
		Model: "amd_mi350x",
		VRAM:  &godo.VRAM{Amount: 2304, Unit: "gib"},
		SupportedPartitionModes: []string{
			godo.GPUPartitionModeDPXNPS2,
			godo.GPUPartitionModeSPXNPS1,
		},
	}

	result := flattenGPUInfo(info)
	if len(result) != 1 {
		t.Fatalf("expected 1 gpu_info entry, got %d", len(result))
	}

	gpuInfo := result[0].(map[string]interface{})
	assert.Equal(t, 8, gpuInfo["count"])
	assert.Equal(t, "amd_mi350x", gpuInfo["model"])

	vram := gpuInfo["vram"].([]interface{})
	if len(vram) != 1 {
		t.Fatalf("expected 1 vram entry, got %d", len(vram))
	}
	vramMap := vram[0].(map[string]interface{})
	assert.Equal(t, 2304, vramMap["amount"])
	assert.Equal(t, "gib", vramMap["unit"])

	assert.Equal(t, []interface{}{
		"PARTITION_MODE_DPX_NPS2",
		"PARTITION_MODE_SPX_NPS1",
	}, gpuInfo["supported_partition_modes"])
}

func TestFlattenGPUInfo_EmptySupportedPartitionModes(t *testing.T) {
	// When the field is absent on the wire, godo decodes it to nil. The provider
	// should expose an empty list (not nil) so consumers can treat it uniformly.
	info := &godo.GPUInfo{
		Count: 1,
		Model: "amd_mi350x",
		VRAM:  &godo.VRAM{Amount: 288, Unit: "gib"},
	}

	result := flattenGPUInfo(info)
	if len(result) != 1 {
		t.Fatalf("expected 1 gpu_info entry, got %d", len(result))
	}

	gpuInfo := result[0].(map[string]interface{})
	assert.Equal(t, []interface{}{}, gpuInfo["supported_partition_modes"])
}

func TestFlattenGPUInfo_NilVRAM(t *testing.T) {
	result := flattenGPUInfo(&godo.GPUInfo{Count: 1, Model: "amd_mi350x"})
	if len(result) != 1 {
		t.Fatalf("expected 1 gpu_info entry, got %d", len(result))
	}

	gpuInfo := result[0].(map[string]interface{})
	assert.Equal(t, []interface{}{}, gpuInfo["vram"])
}

func TestFlattenDigitalOceanSize_WithGPUInfo(t *testing.T) {
	size := godo.Size{
		Slug: "gpu-mi350x8-2304gb",
		GPUInfo: &godo.GPUInfo{
			Count:                   8,
			Model:                   "amd_mi350x",
			SupportedPartitionModes: []string{godo.GPUPartitionModeSPXNPS1},
		},
	}

	flattened, err := flattenDigitalOceanSize(size, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	gpuInfo := flattened["gpu_info"].([]interface{})
	if len(gpuInfo) != 1 {
		t.Fatalf("expected 1 gpu_info entry, got %d", len(gpuInfo))
	}

	modes := gpuInfo[0].(map[string]interface{})["supported_partition_modes"]
	assert.Equal(t, []interface{}{"PARTITION_MODE_SPX_NPS1"}, modes)
}

func TestFlattenDigitalOceanSize_NoGPUInfo(t *testing.T) {
	flattened, err := flattenDigitalOceanSize(godo.Size{Slug: "s-1vcpu-1gb"}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assert.Equal(t, []interface{}{}, flattened["gpu_info"])
}
