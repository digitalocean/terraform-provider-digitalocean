package util

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CaseSensitive implements a schema.SchemaDiffSuppressFunc that ignores case
func CaseSensitive(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

// Float32Precision suppresses diffs caused by float64 config values being
// stored/returned as float32 by the API (e.g. 0.001 vs 0.0010000000474974513).
func Float32Precision(_, old, new string, _ *schema.ResourceData) bool {
	oldF, errOld := strconv.ParseFloat(old, 64)
	newF, errNew := strconv.ParseFloat(new, 64)
	if errOld != nil || errNew != nil {
		return false
	}
	return float32(oldF) == float32(newF)
}

// NormalizeFloat32 converts an API float32 into the shortest float64 that
// round-trips as that float32 (e.g. float32(0.001) -> 0.001 instead of
// 0.0010000000474974513). Use this when writing float32 API values into
// TypeFloat state so plans/outputs stay stable.
func NormalizeFloat32(f float32) float64 {
	s := strconv.FormatFloat(float64(f), 'g', -1, 32)
	out, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return float64(f)
	}
	return out
}

// SetFloat32Attribute sets a TypeFloat attribute from an API *float32 value,
// normalizing float32 expansion so state matches typical HCL literals.
func SetFloat32Attribute(d *schema.ResourceData, key string, apiVal *float32) {
	if apiVal == nil {
		return
	}
	d.Set(key, NormalizeFloat32(*apiVal))
}
