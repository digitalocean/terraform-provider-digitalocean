package app

import (
	"testing"

	"github.com/digitalocean/godo"
)

func TestFlattenAppEnvs(t *testing.T) {
	const otherCiphertext = "EV[1:b2xkZXJub25jZWZvcnRlc3Rpbmc=:b2xkZXJjaXBoZXJ0ZXh0]"

	secret := func(value string) *godo.AppVariableDefinition {
		return &godo.AppVariableDefinition{
			Key:   "TOKEN",
			Value: value,
			Scope: godo.AppVariableScope_RunTime,
			Type:  godo.AppVariableType_Secret,
		}
	}
	tracked := func(key, value, scope string) map[string]any {
		return map[string]any{
			"key":   key,
			"value": value,
			"scope": scope,
			"type":  "SECRET",
		}
	}

	cases := []struct {
		name    string
		env     *godo.AppVariableDefinition
		tracked []any
		want    string
	}{
		{
			name: "general value is stored as returned, even when it looks encrypted",
			env: &godo.AppVariableDefinition{
				Key:   "TOKEN",
				Value: testSecretCiphertext,
				Scope: godo.AppVariableScope_RunTime,
				Type:  godo.AppVariableType_General,
			},
			tracked: []any{tracked("TOKEN", "tracked-value", "RUN_TIME")},
			want:    testSecretCiphertext,
		},
		{
			name:    "encrypted secret keeps the tracked plaintext",
			env:     secret(testSecretCiphertext),
			tracked: []any{tracked("TOKEN", "tracked-plaintext", "RUN_TIME")},
			want:    "tracked-plaintext",
		},
		{
			name: "encrypted secret with nothing tracked keeps the ciphertext",
			env:  secret(testSecretCiphertext),
			want: testSecretCiphertext,
		},
		{
			name:    "encrypted secret tracked as ciphertext keeps the ciphertext returned",
			env:     secret(testSecretCiphertext),
			tracked: []any{tracked("TOKEN", otherCiphertext, "RUN_TIME")},
			want:    testSecretCiphertext,
		},
		{
			name: "encrypted secret keeps the plaintext tracked under its own scope",
			env:  secret(testSecretCiphertext),
			tracked: []any{
				tracked("TOKEN", "build-plaintext", "BUILD_TIME"),
				tracked("TOKEN", "run-plaintext", "RUN_TIME"),
			},
			want: "run-plaintext",
		},
		{
			name:    "encrypted secret tracked only under another key keeps the ciphertext",
			env:     secret(testSecretCiphertext),
			tracked: []any{tracked("OTHER_TOKEN", "tracked-plaintext", "RUN_TIME")},
			want:    testSecretCiphertext,
		},
		{
			name:    "secret returned unencrypted is stored as returned",
			env:     secret("returned-plaintext"),
			tracked: []any{tracked("TOKEN", "tracked-plaintext", "RUN_TIME")},
			want:    "returned-plaintext",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			envs := flattenAppEnvs([]*godo.AppVariableDefinition{tc.env}, tc.tracked).List()
			if len(envs) != 1 {
				t.Fatalf("flattened %d env vars, expected 1", len(envs))
			}
			if value := envs[0].(map[string]any)["value"]; value != tc.want {
				t.Errorf("value = %q, expected %q", value, tc.want)
			}
		})
	}
}
