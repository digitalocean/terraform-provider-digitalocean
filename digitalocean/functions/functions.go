package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

const (
	// triggerTypeScheduled is currently the only trigger type supported by
	// the DigitalOcean Functions API.
	triggerTypeScheduled = "SCHEDULED"
)

func functionsService(meta interface{}) godo.FunctionsService {
	return meta.(*config.CombinedConfig).GodoClient().Functions
}

func isNotFound(resp *godo.Response, err error) bool {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}
	if errResp, ok := err.(*godo.ErrorResponse); ok && errResp.Response != nil {
		return errResp.Response.StatusCode == http.StatusNotFound
	}
	return false
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

// makeTriggerID builds the Terraform ID of a trigger. Trigger names are only
// unique within a namespace, so the namespace is part of the ID.
func makeTriggerID(namespaceID, name string) string {
	return fmt.Sprintf("%s/%s", namespaceID, name)
}

// parseTriggerID accepts both the "<namespace_id>/<name>" form stored in state
// and the "<namespace_id>,<name>" form used by other composite imports in this
// provider.
func parseTriggerID(id string) (string, string, error) {
	for _, sep := range []string{",", "/"} {
		if parts := strings.SplitN(id, sep, 2); len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			return parts[0], parts[1], nil
		}
	}
	return "", "", fmt.Errorf("invalid functions trigger ID %q: expected <namespace_id>,<trigger_name>", id)
}

// expandTriggerBody converts the JSON object configured in Terraform into the
// payload map sent to the API. An empty string means no body.
func expandTriggerBody(raw string) (map[string]interface{}, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var body map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		return nil, fmt.Errorf("body must be a JSON object: %w", err)
	}
	if body == nil {
		return nil, fmt.Errorf("body must be a JSON object, got null")
	}
	return body, nil
}

func flattenTriggerBody(body map[string]interface{}) (string, error) {
	if len(body) == 0 {
		return "", nil
	}
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func validateTriggerBody(v interface{}, k string) ([]string, []error) {
	if _, err := expandTriggerBody(v.(string)); err != nil {
		return nil, []error{fmt.Errorf("%q: %w", k, err)}
	}
	return nil, nil
}

// suppressEquivalentJSON hides diffs between JSON documents that only differ
// in whitespace or key order.
func suppressEquivalentJSON(_, old, new string, _ *schema.ResourceData) bool {
	if old == new {
		return true
	}
	if strings.TrimSpace(old) == "" || strings.TrimSpace(new) == "" {
		return false
	}
	o, err := structure.NormalizeJsonString(old)
	if err != nil {
		return false
	}
	n, err := structure.NormalizeJsonString(new)
	if err != nil {
		return false
	}
	return o == n
}
