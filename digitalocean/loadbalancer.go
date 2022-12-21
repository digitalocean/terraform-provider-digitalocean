package digitalocean

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func loadbalancerStateRefreshFunc(client *godo.Client, loadbalancerId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, _, err := client.LoadBalancers.Get(context.Background(), loadbalancerId)
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in LoadbalancerStateRefreshFunc to DigitalOcean for Load Balancer '%s': %s", loadbalancerId, err)
		}

		return lb, lb.Status, nil
	}
}

func expandStickySessions(config []interface{}) *godo.StickySessions {
	stickysessionConfig := config[0].(map[string]interface{})

	stickySession := &godo.StickySessions{
		Type: stickysessionConfig["type"].(string),
	}

	if v, ok := stickysessionConfig["cookie_name"]; ok {
		stickySession.CookieName = v.(string)
	}

	if v, ok := stickysessionConfig["cookie_ttl_seconds"]; ok {
		stickySession.CookieTtlSeconds = v.(int)
	}

	return stickySession
}

func expandLBFirewall(config []interface{}) *godo.LBFirewall {
	firewallConfig := config[0].(map[string]interface{})

	firewall := &godo.LBFirewall{}

	if v, ok := firewallConfig["allow"]; ok {
		allows := make([]string, 0, len(v.([]interface{})))
		for _, val := range v.([]interface{}) {
			allows = append(allows, val.(string))
		}
		firewall.Allow = allows
	}

	if v, ok := firewallConfig["deny"]; ok {
		denies := make([]string, 0, len(v.([]interface{})))
		for _, val := range v.([]interface{}) {
			denies = append(denies, val.(string))
		}
		firewall.Deny = denies
	}

	return firewall
}

func expandHealthCheck(config []interface{}) *godo.HealthCheck {
	healthcheckConfig := config[0].(map[string]interface{})

	healthcheck := &godo.HealthCheck{
		Protocol:               healthcheckConfig["protocol"].(string),
		Port:                   healthcheckConfig["port"].(int),
		CheckIntervalSeconds:   healthcheckConfig["check_interval_seconds"].(int),
		ResponseTimeoutSeconds: healthcheckConfig["response_timeout_seconds"].(int),
		UnhealthyThreshold:     healthcheckConfig["unhealthy_threshold"].(int),
		HealthyThreshold:       healthcheckConfig["healthy_threshold"].(int),
	}

	if v, ok := healthcheckConfig["path"]; ok {
		healthcheck.Path = v.(string)
	}

	return healthcheck
}

func expandForwardingRules(client *godo.Client, config []interface{}) ([]godo.ForwardingRule, error) {
	forwardingRules := make([]godo.ForwardingRule, 0, len(config))

	for _, rawRule := range config {
		rule := rawRule.(map[string]interface{})

		r := godo.ForwardingRule{
			EntryPort:      rule["entry_port"].(int),
			EntryProtocol:  rule["entry_protocol"].(string),
			TargetPort:     rule["target_port"].(int),
			TargetProtocol: rule["target_protocol"].(string),
			TlsPassthrough: rule["tls_passthrough"].(bool),
		}

		if name, nameOk := rule["certificate_name"]; nameOk {
			certName := name.(string)
			if certName != "" {
				cert, err := findCertificateByName(client, certName)
				if err != nil {
					return nil, err
				}

				r.CertificateID = cert.ID
			}
		}

		if id, idOk := rule["certificate_id"]; idOk && r.CertificateID == "" {
			// When the certificate type is lets_encrypt, the certificate
			// ID will change when it's renewed, so we have to rely on the
			// certificate name as the primary identifier instead.
			certName := id.(string)
			if certName != "" {
				cert, err := findCertificateByName(client, certName)
				if err != nil {
					if strings.Contains(err.Error(), "not found") {
						log.Println("[DEBUG] Certificate not found looking up by name. Falling back to lookup by ID.")
						cert, _, err = client.Certificates.Get(context.Background(), certName)
						if err != nil {
							return nil, err
						}
					} else {
						return nil, err
					}
				}

				r.CertificateID = cert.ID
			}
		}

		forwardingRules = append(forwardingRules, r)

	}

	return forwardingRules, nil
}

func hashForwardingRules(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["entry_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-",
		strings.ToLower(m["entry_protocol"].(string))))
	buf.WriteString(fmt.Sprintf("%d-", m["target_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-",
		strings.ToLower(m["target_protocol"].(string))))

	if v, ok := m["certificate_id"]; ok {
		if v.(string) == "" {
			if name, nameOk := m["certificate_name"]; nameOk {
				buf.WriteString(fmt.Sprintf("%s-", name.(string)))
			}
		} else {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
	}

	if v, ok := m["tls_passthrough"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
	}

	return SDKHashString(buf.String())
}

func flattenDropletIds(list []int) *schema.Set {
	flatSet := schema.NewSet(schema.HashInt, []interface{}{})
	for _, v := range list {
		flatSet.Add(v)
	}
	return flatSet
}

func flattenHealthChecks(health *godo.HealthCheck) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	if health != nil {

		r := make(map[string]interface{})
		r["protocol"] = (*health).Protocol
		r["port"] = (*health).Port
		r["path"] = (*health).Path
		r["check_interval_seconds"] = (*health).CheckIntervalSeconds
		r["response_timeout_seconds"] = (*health).ResponseTimeoutSeconds
		r["unhealthy_threshold"] = (*health).UnhealthyThreshold
		r["healthy_threshold"] = (*health).HealthyThreshold

		result = append(result, r)
	}

	return result
}

func flattenStickySessions(session *godo.StickySessions) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	if session != nil {

		r := make(map[string]interface{})
		r["type"] = (*session).Type
		r["cookie_name"] = (*session).CookieName
		r["cookie_ttl_seconds"] = (*session).CookieTtlSeconds

		result = append(result, r)
	}

	return result
}

func flattenLBFirewall(firewall *godo.LBFirewall) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	if firewall != nil {
		r := make(map[string]interface{})
		r["allow"] = (*firewall).Allow
		r["deny"] = (*firewall).Deny

		result = append(result, r)
	}

	return result
}

func flattenForwardingRules(client *godo.Client, rules []godo.ForwardingRule) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, 1)

	for _, rule := range rules {
		r := make(map[string]interface{})

		r["entry_protocol"] = rule.EntryProtocol
		r["entry_port"] = rule.EntryPort
		r["target_protocol"] = rule.TargetProtocol
		r["target_port"] = rule.TargetPort
		r["tls_passthrough"] = rule.TlsPassthrough

		if rule.CertificateID != "" {
			// When the certificate type is lets_encrypt, the certificate
			// ID will change when it's renewed, so we have to rely on the
			// certificate name as the primary identifier instead.
			cert, _, err := client.Certificates.Get(context.Background(), rule.CertificateID)
			if err != nil {
				return nil, err
			}
			r["certificate_id"] = cert.Name
			r["certificate_name"] = cert.Name
		}

		result = append(result, r)
	}

	return result, nil
}
