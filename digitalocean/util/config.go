package util

import (
	"fmt"
	"os"
)

// GetMultiEnvVar is a helper function that returns the value of the first
// environment variable in the given list that returns a non-empty value.
// It replaces the SDK's schema.MultiEnvDefaultFunc for use with the plugin framework
func GetMultiEnvVar(envVars ...string) (string, error) {
	for _, val := range envVars {
		if v := os.Getenv(val); v != "" {
			return v, nil
		}
	}
	return "", fmt.Errorf("unable to retrieve any value from env vars: %v", envVars)
}
