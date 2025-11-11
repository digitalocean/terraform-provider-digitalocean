package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanAppCreate,
		ReadContext:   resourceDigitalOceanAppRead,
		UpdateContext: resourceDigitalOceanAppUpdate,
		DeleteContext: resourceDigitalOceanAppDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"spec": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "A DigitalOcean App Platform Spec",
				Elem: &schema.Resource{
					Schema: appSpecSchema(true),
				},
			},

			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed attributes
			"default_ingress": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default URL to access the App",
			},

			"dedicated_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The dedicated egress IP addresses associated with the app.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "The IP address of the dedicated egress IP.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "The ID of the dedicated egress IP.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "The status of the dedicated egress IP: 'UNKNOWN', 'ASSIGNING', 'ASSIGNED', or 'REMOVED'",
						},
					},
				},
			},

			"live_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"live_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// TODO: The full Deployment should be a data source, not a resource
			// specify the app id for the active deployment, include a deployment
			// id for a specific one
			"active_deployment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID the App's currently active deployment",
			},

			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The uniform resource identifier for the app",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of when the App was last updated",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of when the App was created",
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceDigitalOceanAppCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	appCreateRequest := &godo.AppCreateRequest{}
	appCreateRequest.Spec = expandAppSpec(d.Get("spec").([]interface{}))

	if v, ok := d.GetOk("project_id"); ok {
		appCreateRequest.ProjectID = v.(string)
	}

	log.Printf("[DEBUG] App create request: %#v", appCreateRequest)
	app, _, err := client.Apps.Create(context.Background(), appCreateRequest)
	if err != nil {
		return diag.Errorf("Error creating App: %s", err)
	}

	d.SetId(app.ID)
	log.Printf("[DEBUG] Waiting for app (%s) deployment to become active", app.ID)
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForAppDeployment(client, app.ID, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] App created, ID: %s", d.Id())
	err = syncAppAlertDestinations(d, client, app.ID)
	if err != nil {
		return diag.Errorf("Error updating alert destination: %s", err)
	}

	return resourceDigitalOceanAppRead(ctx, d, meta)
}

type alertMatchCallback = func(a, b *godo.AppAlert) bool

func resourceDigitalOceanAlertDestinationUpdate(appID string, existingAlerts, schemaAlerts []*godo.AppAlert, client *godo.Client, match alertMatchCallback) error {
	for _, schemaAlert := range schemaAlerts {
		for _, alert := range existingAlerts {
			if !match(schemaAlert, alert) {
				continue
			}
			log.Printf("[DEBUG] Updating alert destination: %#v", alert)
			_, _, err := client.Apps.UpdateAlertDestinations(context.Background(), appID, alert.ID, &godo.AlertDestinationUpdateRequest{
				Emails:        schemaAlert.Emails,
				SlackWebhooks: schemaAlert.SlackWebhooks,
			})

			if err != nil {
				return err
			}
			log.Printf("[INFO] Alert updated: %s", alert.ID)
		}
	}

	return nil
}

func resourceDigitalOceanAppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	app, resp, err := client.Apps.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] App (%s) was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading App: %s", err)
	}

	d.SetId(app.ID)
	d.Set("default_ingress", app.DefaultIngress)
	d.Set("live_url", app.LiveURL)
	d.Set("live_domain", app.LiveDomain)
	d.Set("updated_at", app.UpdatedAt.UTC().String())
	d.Set("created_at", app.CreatedAt.UTC().String())
	d.Set("urn", app.URN())
	d.Set("project_id", app.ProjectID)

	if app.DedicatedIps != nil {
		d.Set("dedicated_ips", appDedicatedIps(d, app))
	}

	if err := d.Set("spec", flattenAppSpec(d, app.Spec)); err != nil {
		return diag.Errorf("Error setting app spec: %#v", err)
	}

	if app.ActiveDeployment != nil {
		d.Set("active_deployment_id", app.ActiveDeployment.ID)
	} else {
		deploymentWarning := diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("No active deployment found for app: %s (%s)", app.Spec.Name, app.ID),
		}
		d.Set("active_deployment_id", "")
		return diag.Diagnostics{deploymentWarning}
	}

	return nil
}

func appDedicatedIps(d *schema.ResourceData, app *godo.App) []interface{} {
	remote := make([]interface{}, 0, len(app.DedicatedIps))
	for _, change := range app.DedicatedIps {
		rawChange := map[string]interface{}{
			"ip":     change.Ip,
			"id":     change.ID,
			"status": change.Status,
		}
		remote = append(remote, rawChange)
	}
	return remote
}

func resourceDigitalOceanAppUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if d.HasChange("spec") {
		appUpdateRequest := &godo.AppUpdateRequest{}
		appUpdateRequest.Spec = expandAppSpec(d.Get("spec").([]interface{}))

		app, _, err := client.Apps.Update(context.Background(), d.Id(), appUpdateRequest)
		if err != nil {
			return diag.Errorf("Error updating app (%s): %s", d.Id(), err)
		}

		log.Printf("[DEBUG] Waiting for app (%s) deployment to become active", app.ID)
		timeout := d.Timeout(schema.TimeoutCreate)
		err = waitForAppDeployment(client, app.ID, timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[INFO] Updated app (%s)", app.ID)

		err = syncAppAlertDestinations(d, client, app.ID)
		if err != nil {
			return diag.Errorf("Error updating alert destination: %s", err)
		}
	}

	return resourceDigitalOceanAppRead(ctx, d, meta)
}

func resourceDigitalOceanAppDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	log.Printf("[INFO] Deleting App: %s", d.Id())
	_, err := client.Apps.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deletingApp: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForAppDeployment(client *godo.Client, id string, timeout time.Duration) error {
	tickerInterval := 10 //10s
	timeoutSeconds := int(timeout.Seconds())
	n := 0

	var deploymentID string
	ticker := time.NewTicker(time.Duration(tickerInterval) * time.Second)
	for range ticker.C {
		if n*tickerInterval > timeoutSeconds {
			ticker.Stop()
			break
		}

		if deploymentID == "" {
			// The InProgressDeployment is generally not known and returned as
			// part of the initial response to the request. For config updates
			// (as opposed to updates to the app's source), the "deployment"
			// can complete before the first time we poll the app. We can not
			// know if the InProgressDeployment has not started or if it has
			// already completed. So instead we need to list all of the
			// deployments for the application.
			opts := &godo.ListOptions{PerPage: 2}
			deployments, _, err := client.Apps.ListDeployments(context.Background(), id, opts)
			if err != nil {
				return fmt.Errorf("Error trying to read app deployment state: %s", err)
			}

			// We choose the most recent deployment. Note that there is a possibility
			// that the deployment has not been created yet. If that is true,
			// we will do the wrong thing here and test the status of a previously
			// completed deployment and exit. However there is no better way to
			// correlate a deployment with the request that triggered it.
			if len(deployments) > 0 {
				deploymentID = deployments[0].ID
			}
		} else {
			deployment, _, err := client.Apps.GetDeployment(context.Background(), id, deploymentID)
			if err != nil {
				ticker.Stop()
				return fmt.Errorf("Error trying to read app deployment state: %s", err)
			}

			allSuccessful := true
			for _, step := range deployment.Progress.Steps {
				if step.Status != godo.DeploymentProgressStepStatus_Success {
					allSuccessful = false
					break
				}
			}

			if allSuccessful {
				ticker.Stop()
				return nil
			}

			if deployment.Progress.ErrorSteps > 0 {
				ticker.Stop()
				return fmt.Errorf("error deploying app (%s) (deployment ID: %s):\n%s", id, deployment.ID, godo.Stringify(deployment.Progress))
			}

			log.Printf("[DEBUG] Waiting for app (%s) deployment (%s) to become active. Phase: %s (%d/%d)",
				id, deployment.ID, deployment.Phase, deployment.Progress.SuccessSteps, deployment.Progress.TotalSteps)
		}

		n++
	}

	return fmt.Errorf("timeout waiting for app (%s) deployment", id)
}

func mapAppLevelAlerts(d *schema.ResourceData) []*godo.AppAlert {
	spec, ok := d.GetOk("spec")
	if !ok {
		return nil
	}

	specList := spec.([]interface{})
	if len(specList) == 0 {
		return nil
	}
	specMap, ok := specList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	alert, ok := specMap["alert"]
	if !ok {
		return nil
	}

	alertList, ok := alert.([]interface{})
	if !ok {
		return nil
	}

	var alertDetails []*godo.AppAlert

	for _, alertItem := range alertList {
		alertData, ok := alertItem.(map[string]interface{})
		if !ok {
			continue
		}

		alertDetail := &godo.AppAlert{
			Spec:          &godo.AppAlertSpec{},
			Emails:        []string{},
			SlackWebhooks: []*godo.AppAlertSlackWebhook{},
		}

		for key, value := range alertData {
			switch key {
			case "rule":
				if rule, ok := value.(string); ok {
					alertDetail.Spec.Rule = godo.AppAlertSpecRule(rule)
				}
			case "destinations":
				if destinations, ok := value.(*schema.Set); ok {
					emails, slackWebhooks := extractEmailsAndSlackDestinations(destinations)

					if len(emails) == 0 && len(slackWebhooks) == 0 {
						continue
					}

					alertDetail.Emails = emails
					alertDetail.SlackWebhooks = slackWebhooks
				}
			}
		}

		alertDetails = append(alertDetails, alertDetail)
	}

	return alertDetails
}

func syncAppAlertDestinations(d *schema.ResourceData, client *godo.Client, appID string) error {
	appConfigurationAlert := mapAppLevelAlerts(d)
	computeComponentAlerts := mapComputeComponentAlert(d)

	hasAppLevelAlerts := len(appConfigurationAlert) > 0 && (len(appConfigurationAlert[0].Emails) > 0 || len(appConfigurationAlert[0].SlackWebhooks) > 0)
	hasComputeComponentAlerts := len(computeComponentAlerts) > 0 && (len(computeComponentAlerts[0].Emails) > 0 || len(computeComponentAlerts[0].SlackWebhooks) > 0)

	var alerts []*godo.AppAlert
	if hasAppLevelAlerts || hasComputeComponentAlerts {
		appAlerts, _, err := client.Apps.ListAlerts(context.Background(), appID)
		if err != nil {
			return err
		}
		alerts = appAlerts
	}

	if hasAppLevelAlerts {
		err := resourceDigitalOceanAlertDestinationUpdate(appID, alerts, appConfigurationAlert, client, func(a, b *godo.AppAlert) bool {
			return string(a.Spec.Rule) == string(b.Spec.Rule)
		})
		if err != nil {
			return err
		}
	}

	if hasComputeComponentAlerts {
		err := resourceDigitalOceanAlertDestinationUpdate(appID, alerts, computeComponentAlerts, client, func(a, b *godo.AppAlert) bool {
			return string(a.Spec.Rule) == string(b.Spec.Rule)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func mapComputeComponentAlert(d *schema.ResourceData) []*godo.AppAlert {
	spec, ok := d.GetOk("spec")
	if !ok {
		return nil
	}

	specList := spec.([]interface{})
	if len(specList) == 0 {
		return nil
	}
	specMap, ok := specList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	var alertDetails []*godo.AppAlert

	if serviceSpec, ok := specMap["service"]; ok {
		if services, ok := serviceSpec.([]interface{}); ok {
			alertDetails = extractAlertDetailsFromComputeComponent(services, alertDetails)
		}
	}

	if workerSpec, ok := specMap["worker"]; ok {
		if workers, ok := workerSpec.([]interface{}); ok {
			alertDetails = extractAlertDetailsFromComputeComponent(workers, alertDetails)
		}
	}

	if jobSpec, ok := specMap["job"]; ok {
		if jobs, ok := jobSpec.([]interface{}); ok {
			alertDetails = extractAlertDetailsFromComputeComponent(jobs, alertDetails)
		}
	}

	if functionSpec, ok := specMap["function"]; ok {
		if functions, ok := functionSpec.([]interface{}); ok {
			alertDetails = extractAlertDetailsFromComputeComponent(functions, alertDetails)
		}
	}

	return alertDetails
}

func extractAlertDetailsFromComputeComponent(computeComponents []interface{}, alertDetails []*godo.AppAlert) []*godo.AppAlert {
	for _, computeServiceMap := range computeComponents {
		computeServiceMapTyped, ok := computeServiceMap.(map[string]interface{})
		if !ok {
			continue
		}

		alertRaw, ok := computeServiceMapTyped["alert"].([]interface{})
		if !ok {
			continue // Skip if no alerts are defined for this compute component
		}

		for _, alert := range alertRaw {
			alertMap, ok := alert.(map[string]interface{})
			if !ok {
				continue
			}

			destinations, ok := alertMap["destinations"].(*schema.Set)
			if !ok {
				continue
			}

			emails, slackWebhooks := extractEmailsAndSlackDestinations(destinations)
			if len(emails) == 0 && len(slackWebhooks) == 0 {
				continue
			}

			appAlert := &godo.AppAlert{
				Spec:          &godo.AppAlertSpec{},
				Emails:        emails,
				SlackWebhooks: slackWebhooks,
			}
			if rule, ok := alertMap["rule"].(string); ok {
				appAlert.Spec.Rule = godo.AppAlertSpecRule(rule)
			}
			if value, ok := alertMap["value"].(float32); ok {
				appAlert.Spec.Value = value
			}
			if operator, ok := alertMap["operator"].(string); ok {
				appAlert.Spec.Operator = godo.AppAlertSpecOperator(operator)
			}
			if window, ok := alertMap["window"].(string); ok {
				appAlert.Spec.Window = godo.AppAlertSpecWindow(window)
			}

			alertDetails = append(alertDetails, appAlert)
		}
	}
	return alertDetails
}

func extractEmailsAndSlackDestinations(destinations interface{}) ([]string, []*godo.AppAlertSlackWebhook) {
	var emails []string
	var slackWebhookAlerts []*godo.AppAlertSlackWebhook

	if destinations == nil {
		return emails, slackWebhookAlerts
	}

	destinationSet, ok := destinations.(*schema.Set)
	if !ok {
		return emails, slackWebhookAlerts
	}

	for _, destinationInterface := range destinationSet.List() {
		destination, ok := destinationInterface.(map[string]interface{})
		if !ok {
			continue
		}

		if destinationEmails, ok := destination["emails"].([]interface{}); ok && len(destinationEmails) > 0 {
			for _, email := range destinationEmails {
				if emailStr, ok := email.(string); ok {
					emails = append(emails, emailStr)
				}
			}
		}

		if slackWebhooks, ok := destination["slack_webhooks"].([]interface{}); ok && len(slackWebhooks) > 0 {
			for _, slackInterface := range slackWebhooks {
				if slack, ok := slackInterface.(map[string]interface{}); ok {
					slackDetail := &godo.AppAlertSlackWebhook{}
					if channel, ok := slack["channel"].(string); ok {
						slackDetail.Channel = channel
					}
					if url, ok := slack["url"].(string); ok {
						slackDetail.URL = url
					}
					slackWebhookAlerts = append(slackWebhookAlerts, slackDetail)
				}
			}
		}
	}

	return emails, slackWebhookAlerts
}
