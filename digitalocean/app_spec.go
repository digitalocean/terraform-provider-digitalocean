package digitalocean

import (
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appSpecSchema returns map[string]*schema.Schema for the App Specification.
// Set isResource to true in order to return a schema with additional attributes
// appropriate for a resource or false for one used with a data-source.
func appSpecSchema(isResource bool) map[string]*schema.Schema {
	spec := map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(2, 32),
			Description:  "The name of the app. Must be unique across all apps in the same account.",
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The slug for the DigitalOcean data center region hosting the app",
		},
		"domain": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     appSpecDomainSchema(),
		},
		"domains": {
			Type:       schema.TypeSet,
			Optional:   true,
			Computed:   true,
			Elem:       &schema.Schema{Type: schema.TypeString},
			Deprecated: "This attribute has been replaced by `domain` which supports additional functionality.",
		},
		"service": {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			Elem:     appSpecServicesSchema(),
		},
		"static_site": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     appSpecStaticSiteSchema(),
		},
		"worker": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     appSpecWorkerSchema(),
		},
		"job": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     appSpecJobSchema(),
		},
		"database": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     appSpecDatabaseSchema(),
		},
		"env": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     appSpecEnvSchema(),
			Set:      schema.HashResource(appSpecEnvSchema()),
		},
	}

	if isResource {
		spec["domain"].ConflictsWith = []string{"spec.0.domains"}
	}

	return spec
}

func appSpecDomainSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The hostname for the domain.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DEFAULT",
					"PRIMARY",
					"ALIAS",
				}, false),
				Description: "The type of the domain.",
			},
			"wildcard": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the domain includes all sub-domains, in addition to the given domain.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If the domain uses DigitalOcean DNS and you would like App Platform to automatically manage it for you, set this to the name of the domain on your account.",
			},
		},
	}
}

func appSpecGitSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"repo_clone_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The clone URL of the repo.",
		},
		"branch": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the branch to use.",
		},
	}
}

func appSpecGitServiceSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"repo": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the repo in the format `owner/repo`.",
		},
		"branch": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the branch to use.",
		},
		"deploy_on_push": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to automatically deploy new commits made to the repo",
		},
	}
}

func appSpecGitHubSourceSchema() map[string]*schema.Schema {
	return appSpecGitServiceSourceSchema()
}

func appSpecGitLabSourceSchema() map[string]*schema.Schema {
	return appSpecGitServiceSourceSchema()
}

func appSpecImageSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"registry_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"UNSPECIFIED",
				"DOCKER_HUB",
				"DOCR",
			}, false),
			Description: "The registry type.",
		},
		"registry": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The registry name. Must be left empty for the DOCR registry type.",
		},
		"repository": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The repository name.",
		},
		"tag": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The repository tag. Defaults to latest if not provided.",
		},
	}
}

func appSpecEnvSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the environment variable.",
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The value of the environment variable.",
				Sensitive:   true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "RUN_AND_BUILD_TIME",
				ValidateFunc: validation.StringInSlice([]string{
					"UNSET",
					"RUN_TIME",
					"BUILD_TIME",
					"RUN_AND_BUILD_TIME",
				}, false),
				Description: "The visibility scope of the environment variable.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GENERAL",
					"SECRET",
				}, false),
				Description: "The type of the environment variable.",
				// The API does not always return `"type":"GENERAL"` when set.
				// As being unset and being set to `GENERAL` are functionally,
				// the same, we can safely ignore the diff.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "GENERAL" && old == ""
				},
			},
		},
	}
}

func appSpecRouteSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Path specifies an route by HTTP path prefix. Paths must start with / and must be unique within the app.",
		},
	}
}

func appSpecHealthCheckSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"http_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The route path used for the HTTP health check ping.",
		},
		"initial_delay_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The number of seconds to wait before beginning health checks.",
		},
		"period_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The number of seconds to wait between health checks.",
		},
		"timeout_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The number of seconds after which the check times out.",
		},
		"success_threshold": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The number of successful health checks before considered healthy.",
		},
		"failure_threshold": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The number of failed health checks before considered unhealthy.",
		},
	}
}

func appSpecComponentBase() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the component",
		},
		"git": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: appSpecGitSourceSchema(),
			},
		},
		"github": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: appSpecGitHubSourceSchema(),
			},
		},
		"gitlab": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: appSpecGitLabSourceSchema(),
			},
		},
		"dockerfile_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The path to a Dockerfile relative to the root of the repo. If set, overrides usage of buildpacks.",
		},
		"build_command": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An optional build command to run while building this component from source.",
		},
		"env": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     appSpecEnvSchema(),
			Set:      schema.HashResource(appSpecEnvSchema()),
		},
		"source_dir": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An optional path to the working directory to use for the build.",
		},
		"environment_slug": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An environment slug describing the type of this app.",
		},
	}
}

func appSpecServicesSchema() *schema.Resource {
	serviceSchema := map[string]*schema.Schema{
		"run_command": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "An optional run command to override the component's default.",
		},
		"http_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The internal port on which this service's run command will listen.",
		},
		"instance_size_slug": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The instance size to use for this component.",
		},
		"instance_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "The amount of instances that this component should be scaled to.",
		},
		"health_check": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: appSpecHealthCheckSchema(),
			},
		},
		"image": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: appSpecImageSourceSchema(),
			},
		},
		"routes": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: appSpecRouteSchema(),
			},
		},
		"internal_ports": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
		},
	}

	for k, v := range appSpecComponentBase() {
		serviceSchema[k] = v
	}

	return &schema.Resource{
		Schema: serviceSchema,
	}
}

func appSpecStaticSiteSchema() *schema.Resource {
	staticSiteSchema := map[string]*schema.Schema{
		"output_dir": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An optional path to where the built assets will be located, relative to the build context. If not set, App Platform will automatically scan for these directory names: `_static`, `dist`, `public`.",
		},
		"index_document": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the index document to use when serving this static site.",
		},
		"error_document": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the error document to use when serving this static site.",
		},
		"catchall_document": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the document to use as the fallback for any requests to documents that are not found when serving this static site.",
		},
		"routes": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: appSpecRouteSchema(),
			},
		},
	}

	for k, v := range appSpecComponentBase() {
		staticSiteSchema[k] = v
	}

	return &schema.Resource{
		Schema: staticSiteSchema,
	}
}

func appSpecWorkerSchema() *schema.Resource {
	workerSchema := map[string]*schema.Schema{
		"run_command": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An optional run command to override the component's default.",
		},
		"image": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: appSpecImageSourceSchema(),
			},
		},
		"instance_size_slug": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The instance size to use for this component.",
		},
		"instance_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "The amount of instances that this component should be scaled to.",
		},
	}

	for k, v := range appSpecComponentBase() {
		workerSchema[k] = v
	}

	return &schema.Resource{
		Schema: workerSchema,
	}
}

func appSpecJobSchema() *schema.Resource {
	jobSchema := map[string]*schema.Schema{
		"run_command": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An optional run command to override the component's default.",
		},
		"image": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: appSpecImageSourceSchema(),
			},
		},
		"instance_size_slug": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The instance size to use for this component.",
		},
		"instance_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "The amount of instances that this component should be scaled to.",
		},
		"kind": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "UNSPECIFIED",
			ValidateFunc: validation.StringInSlice([]string{
				"UNSPECIFIED",
				"PRE_DEPLOY",
				"POST_DEPLOY",
				"FAILED_DEPLOY",
			}, false),
			Description: "The type of job and when it will be run during the deployment process.",
		},
	}

	for k, v := range appSpecComponentBase() {
		jobSchema[k] = v
	}

	return &schema.Resource{
		Schema: jobSchema,
	}
}

func appSpecDatabaseSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the component",
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"UNSET",
					"MYSQL",
					"PG",
					"REDIS",
				}, false),
				Description: "The database engine to use.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The version of the database engine.",
			},
			"production": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether this is a production or dev database.",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the underlying DigitalOcean DBaaS cluster. This is required for production databases. For dev databases, if cluster_name is not set, a new cluster will be provisioned.",
			},
			"db_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the MySQL or PostgreSQL database to configure.",
			},
			"db_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the MySQL or PostgreSQL user to configure.",
			},
		},
	}
}

func expandAppSpec(config []interface{}) *godo.AppSpec {
	if len(config) == 0 || config[0] == nil {
		return &godo.AppSpec{}
	}
	appSpecConfig := config[0].(map[string]interface{})

	appSpec := &godo.AppSpec{
		Name:        appSpecConfig["name"].(string),
		Region:      appSpecConfig["region"].(string),
		Services:    expandAppSpecServices(appSpecConfig["service"].([]interface{})),
		StaticSites: expandAppSpecStaticSites(appSpecConfig["static_site"].([]interface{})),
		Workers:     expandAppSpecWorkers(appSpecConfig["worker"].([]interface{})),
		Jobs:        expandAppSpecJobs(appSpecConfig["job"].([]interface{})),
		Databases:   expandAppSpecDatabases(appSpecConfig["database"].([]interface{})),
		Envs:        expandAppEnvs(appSpecConfig["env"].(*schema.Set).List()),
	}

	// Prefer the `domain` block over `domains` if it is set.
	domainConfig := appSpecConfig["domain"].([]interface{})
	if len(domainConfig) > 0 {
		appSpec.Domains = expandAppSpecDomains(domainConfig)
	} else {
		appSpec.Domains = expandAppDomainSpec(appSpecConfig["domains"].(*schema.Set).List())
	}

	return appSpec
}

func flattenAppSpec(d *schema.ResourceData, spec *godo.AppSpec) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	if spec != nil {

		r := make(map[string]interface{})
		r["name"] = (*spec).Name
		r["region"] = (*spec).Region

		if len((*spec).Domains) > 0 {
			r["domains"] = flattenAppDomainSpec((*spec).Domains)
			if _, ok := d.GetOk("spec.0.domain"); ok {
				r["domain"] = flattenAppSpecDomains((*spec).Domains)
			}
		}

		if len((*spec).Services) > 0 {
			r["service"] = flattenAppSpecServices((*spec).Services)
		}

		if len((*spec).StaticSites) > 0 {
			r["static_site"] = flattenAppSpecStaticSites((*spec).StaticSites)
		}

		if len((*spec).Workers) > 0 {
			r["worker"] = flattenAppSpecWorkers((*spec).Workers)
		}

		if len((*spec).Jobs) > 0 {
			r["job"] = flattenAppSpecJobs((*spec).Jobs)
		}

		if len((*spec).Databases) > 0 {
			r["database"] = flattenAppSpecDatabases((*spec).Databases)
		}

		if len((*spec).Envs) > 0 {
			r["env"] = flattenAppEnvs((*spec).Envs)
		}

		result = append(result, r)
	}

	return result
}

// expandAppDomainSpec has been deprecated in favor of expandAppSpecDomains.
func expandAppDomainSpec(config []interface{}) []*godo.AppDomainSpec {
	appDomains := make([]*godo.AppDomainSpec, 0, len(config))

	for _, rawDomain := range config {
		domain := &godo.AppDomainSpec{
			Domain: rawDomain.(string),
		}

		appDomains = append(appDomains, domain)
	}

	return appDomains
}

func expandAppSpecDomains(config []interface{}) []*godo.AppDomainSpec {
	appDomains := make([]*godo.AppDomainSpec, 0, len(config))

	for _, rawDomain := range config {
		domain := rawDomain.(map[string]interface{})

		d := &godo.AppDomainSpec{
			Domain:   domain["name"].(string),
			Type:     godo.AppDomainSpecType(domain["type"].(string)),
			Wildcard: domain["wildcard"].(bool),
			Zone:     domain["zone"].(string),
		}

		appDomains = append(appDomains, d)
	}

	return appDomains
}

// flattenAppDomainSpec has been deprecated in favor of flattenAppSpecDomains
func flattenAppDomainSpec(spec []*godo.AppDomainSpec) *schema.Set {
	result := schema.NewSet(schema.HashString, []interface{}{})

	for _, domain := range spec {
		result.Add(domain.Domain)
	}

	return result
}

func flattenAppSpecDomains(domains []*godo.AppDomainSpec) []map[string]interface{} {
	result := make([]map[string]interface{}, len(domains))

	for i, d := range domains {
		r := make(map[string]interface{})

		r["name"] = d.Domain
		r["type"] = string(d.Type)
		r["wildcard"] = d.Wildcard
		r["zone"] = d.Zone

		result[i] = r
	}

	return result
}

func expandAppGitHubSourceSpec(config []interface{}) *godo.GitHubSourceSpec {
	gitHubSourceConfig := config[0].(map[string]interface{})

	gitHubSource := &godo.GitHubSourceSpec{
		Repo:         gitHubSourceConfig["repo"].(string),
		Branch:       gitHubSourceConfig["branch"].(string),
		DeployOnPush: gitHubSourceConfig["deploy_on_push"].(bool),
	}

	return gitHubSource
}

func flattenAppGitHubSourceSpec(spec *godo.GitHubSourceSpec) []interface{} {
	result := make([]interface{}, 0)

	if spec != nil {

		r := make(map[string]interface{})
		r["repo"] = (*spec).Repo
		r["branch"] = (*spec).Branch
		r["deploy_on_push"] = (*spec).DeployOnPush

		result = append(result, r)
	}

	return result
}

func expandAppGitLabSourceSpec(config []interface{}) *godo.GitLabSourceSpec {
	gitLabSourceConfig := config[0].(map[string]interface{})

	gitLabSource := &godo.GitLabSourceSpec{
		Repo:         gitLabSourceConfig["repo"].(string),
		Branch:       gitLabSourceConfig["branch"].(string),
		DeployOnPush: gitLabSourceConfig["deploy_on_push"].(bool),
	}

	return gitLabSource
}

func flattenAppGitLabSourceSpec(spec *godo.GitLabSourceSpec) []interface{} {
	result := make([]interface{}, 0)

	if spec != nil {

		r := make(map[string]interface{})
		r["repo"] = (*spec).Repo
		r["branch"] = (*spec).Branch
		r["deploy_on_push"] = (*spec).DeployOnPush

		result = append(result, r)
	}

	return result
}

func expandAppGitSourceSpec(config []interface{}) *godo.GitSourceSpec {
	gitSourceConfig := config[0].(map[string]interface{})

	gitSource := &godo.GitSourceSpec{
		Branch:       gitSourceConfig["branch"].(string),
		RepoCloneURL: gitSourceConfig["repo_clone_url"].(string),
	}

	return gitSource
}

func flattenAppGitSourceSpec(spec *godo.GitSourceSpec) []interface{} {
	result := make([]interface{}, 0)

	if spec != nil {

		r := make(map[string]interface{})
		r["branch"] = (*spec).Branch
		r["repo_clone_url"] = (*spec).RepoCloneURL

		result = append(result, r)
	}

	return result
}

func expandAppImageSourceSpec(config []interface{}) *godo.ImageSourceSpec {
	imageSourceConfig := config[0].(map[string]interface{})

	imageSource := &godo.ImageSourceSpec{
		RegistryType: godo.ImageSourceSpecRegistryType(imageSourceConfig["registry_type"].(string)),
		Registry:     imageSourceConfig["registry"].(string),
		Repository:   imageSourceConfig["repository"].(string),
		Tag:          imageSourceConfig["tag"].(string),
	}

	return imageSource
}

func flattenAppImageSourceSpec(spec *godo.ImageSourceSpec) []interface{} {
	result := make([]interface{}, 0)

	if spec != nil {
		r := make(map[string]interface{})
		r["registry_type"] = string((*spec).RegistryType)
		r["registry"] = (*spec).Registry
		r["repository"] = (*spec).Repository
		r["tag"] = (*spec).Tag

		result = append(result, r)
	}

	return result
}

func expandAppEnvs(config []interface{}) []*godo.AppVariableDefinition {
	appEnvs := make([]*godo.AppVariableDefinition, 0, len(config))

	for _, rawEnv := range config {
		env := rawEnv.(map[string]interface{})

		e := &godo.AppVariableDefinition{
			Value: env["value"].(string),
			Scope: godo.AppVariableScope(env["scope"].(string)),
			Key:   env["key"].(string),
			Type:  godo.AppVariableType(env["type"].(string)),
		}

		appEnvs = append(appEnvs, e)
	}

	return appEnvs
}

func flattenAppEnvs(appEnvs []*godo.AppVariableDefinition) *schema.Set {
	result := schema.NewSet(schema.HashResource(appSpecEnvSchema()), []interface{}{})

	for _, env := range appEnvs {
		r := make(map[string]interface{})
		r["value"] = env.Value
		r["scope"] = string(env.Scope)
		r["key"] = env.Key
		r["type"] = string(env.Type)

		result.Add(r)

		setFunc := schema.HashResource(appSpecEnvSchema())
		log.Printf("[DEBUG] App env hash for %s: %d", r["key"], setFunc(r))
	}

	return result
}

func expandAppHealthCheck(config []interface{}) *godo.AppServiceSpecHealthCheck {
	healthCheckConfig := config[0].(map[string]interface{})

	healthCheck := &godo.AppServiceSpecHealthCheck{
		HTTPPath:            healthCheckConfig["http_path"].(string),
		InitialDelaySeconds: int32(healthCheckConfig["initial_delay_seconds"].(int)),
		PeriodSeconds:       int32(healthCheckConfig["period_seconds"].(int)),
		TimeoutSeconds:      int32(healthCheckConfig["timeout_seconds"].(int)),
		SuccessThreshold:    int32(healthCheckConfig["success_threshold"].(int)),
		FailureThreshold:    int32(healthCheckConfig["failure_threshold"].(int)),
	}

	return healthCheck
}

func flattenAppHealthCheck(check *godo.AppServiceSpecHealthCheck) []interface{} {
	result := make([]interface{}, 0)

	if check != nil {

		r := make(map[string]interface{})
		r["http_path"] = check.HTTPPath
		r["initial_delay_seconds"] = check.InitialDelaySeconds
		r["period_seconds"] = check.PeriodSeconds
		r["timeout_seconds"] = check.TimeoutSeconds
		r["success_threshold"] = check.SuccessThreshold
		r["failure_threshold"] = check.FailureThreshold

		result = append(result, r)
	}

	return result
}

func expandAppInternalPorts(config []interface{}) []int64 {
	appInternalPorts := make([]int64, len(config))

	for i, v := range config {
		appInternalPorts[i] = int64(v.(int))
	}

	return appInternalPorts
}

func expandAppRoutes(config []interface{}) []*godo.AppRouteSpec {
	appRoutes := make([]*godo.AppRouteSpec, 0, len(config))

	for _, rawRoute := range config {
		route := rawRoute.(map[string]interface{})

		r := &godo.AppRouteSpec{
			Path: route["path"].(string),
		}

		appRoutes = append(appRoutes, r)
	}

	return appRoutes
}

func flattenAppServiceInternalPortsSpec(internalPorts []int64) *schema.Set {
	result := schema.NewSet(schema.HashInt, []interface{}{})

	for _, internalPort := range internalPorts {
		result.Add(int(internalPort))
	}

	return result
}

func flattenAppRoutes(routes []*godo.AppRouteSpec) []interface{} {
	result := make([]interface{}, 0)

	for _, route := range routes {
		r := make(map[string]interface{})
		r["path"] = route.Path

		result = append(result, r)
	}

	return result
}

func expandAppSpecServices(config []interface{}) []*godo.AppServiceSpec {
	appServices := make([]*godo.AppServiceSpec, 0, len(config))

	for _, rawService := range config {
		service := rawService.(map[string]interface{})

		s := &godo.AppServiceSpec{
			Name:             service["name"].(string),
			RunCommand:       service["run_command"].(string),
			BuildCommand:     service["build_command"].(string),
			HTTPPort:         int64(service["http_port"].(int)),
			DockerfilePath:   service["dockerfile_path"].(string),
			Envs:             expandAppEnvs(service["env"].(*schema.Set).List()),
			InstanceSizeSlug: service["instance_size_slug"].(string),
			InstanceCount:    int64(service["instance_count"].(int)),
			SourceDir:        service["source_dir"].(string),
			EnvironmentSlug:  service["environment_slug"].(string),
		}

		github := service["github"].([]interface{})
		if len(github) > 0 {
			s.GitHub = expandAppGitHubSourceSpec(github)
		}

		gitlab := service["gitlab"].([]interface{})
		if len(gitlab) > 0 {
			s.GitLab = expandAppGitLabSourceSpec(gitlab)
		}

		git := service["git"].([]interface{})
		if len(git) > 0 {
			s.Git = expandAppGitSourceSpec(git)
		}

		image := service["image"].([]interface{})
		if len(image) > 0 {
			s.Image = expandAppImageSourceSpec(image)
		}

		routes := service["routes"].([]interface{})
		if len(routes) > 0 {
			s.Routes = expandAppRoutes(routes)
		}

		checks := service["health_check"].([]interface{})
		if len(checks) > 0 {
			s.HealthCheck = expandAppHealthCheck(checks)
		}

		internalPorts := service["internal_ports"].(*schema.Set).List()
		if len(internalPorts) > 0 {
			s.InternalPorts = expandAppInternalPorts(internalPorts)
		}

		appServices = append(appServices, s)
	}

	return appServices
}

func flattenAppSpecServices(services []*godo.AppServiceSpec) []map[string]interface{} {
	result := make([]map[string]interface{}, len(services))

	for i, s := range services {
		r := make(map[string]interface{})

		r["name"] = s.Name
		r["run_command"] = s.RunCommand
		r["build_command"] = s.BuildCommand
		r["github"] = flattenAppGitHubSourceSpec(s.GitHub)
		r["gitlab"] = flattenAppGitLabSourceSpec(s.GitLab)
		r["internal_ports"] = flattenAppServiceInternalPortsSpec(s.InternalPorts)
		r["git"] = flattenAppGitSourceSpec(s.Git)
		r["image"] = flattenAppImageSourceSpec(s.Image)
		r["http_port"] = int(s.HTTPPort)
		r["routes"] = flattenAppRoutes(s.Routes)
		r["dockerfile_path"] = s.DockerfilePath
		r["env"] = flattenAppEnvs(s.Envs)
		r["health_check"] = flattenAppHealthCheck(s.HealthCheck)
		r["instance_size_slug"] = s.InstanceSizeSlug
		r["instance_count"] = int(s.InstanceCount)
		r["source_dir"] = s.SourceDir
		r["environment_slug"] = s.EnvironmentSlug

		result[i] = r
	}

	return result
}

func expandAppSpecStaticSites(config []interface{}) []*godo.AppStaticSiteSpec {
	appSites := make([]*godo.AppStaticSiteSpec, 0, len(config))

	for _, rawSite := range config {
		site := rawSite.(map[string]interface{})

		s := &godo.AppStaticSiteSpec{
			Name:             site["name"].(string),
			BuildCommand:     site["build_command"].(string),
			DockerfilePath:   site["dockerfile_path"].(string),
			Envs:             expandAppEnvs(site["env"].(*schema.Set).List()),
			SourceDir:        site["source_dir"].(string),
			OutputDir:        site["output_dir"].(string),
			IndexDocument:    site["index_document"].(string),
			ErrorDocument:    site["error_document"].(string),
			CatchallDocument: site["catchall_document"].(string),
			EnvironmentSlug:  site["environment_slug"].(string),
		}

		github := site["github"].([]interface{})
		if len(github) > 0 {
			s.GitHub = expandAppGitHubSourceSpec(github)
		}

		gitlab := site["gitlab"].([]interface{})
		if len(gitlab) > 0 {
			s.GitLab = expandAppGitLabSourceSpec(gitlab)
		}

		git := site["git"].([]interface{})
		if len(git) > 0 {
			s.Git = expandAppGitSourceSpec(git)
		}

		routes := site["routes"].([]interface{})
		if len(routes) > 0 {
			s.Routes = expandAppRoutes(routes)
		}

		appSites = append(appSites, s)
	}

	return appSites
}

func flattenAppSpecStaticSites(sites []*godo.AppStaticSiteSpec) []map[string]interface{} {
	result := make([]map[string]interface{}, len(sites))

	for i, s := range sites {
		r := make(map[string]interface{})

		r["name"] = s.Name
		r["build_command"] = s.BuildCommand
		r["github"] = flattenAppGitHubSourceSpec(s.GitHub)
		r["gitlab"] = flattenAppGitLabSourceSpec(s.GitLab)
		r["git"] = flattenAppGitSourceSpec(s.Git)
		r["routes"] = flattenAppRoutes(s.Routes)
		r["dockerfile_path"] = s.DockerfilePath
		r["env"] = flattenAppEnvs(s.Envs)
		r["source_dir"] = s.SourceDir
		r["output_dir"] = s.OutputDir
		r["index_document"] = s.IndexDocument
		r["error_document"] = s.ErrorDocument
		r["catchall_document"] = s.CatchallDocument
		r["environment_slug"] = s.EnvironmentSlug

		result[i] = r
	}

	return result
}

func expandAppSpecWorkers(config []interface{}) []*godo.AppWorkerSpec {
	appWorkers := make([]*godo.AppWorkerSpec, 0, len(config))

	for _, rawWorker := range config {
		worker := rawWorker.(map[string]interface{})

		s := &godo.AppWorkerSpec{
			Name:             worker["name"].(string),
			RunCommand:       worker["run_command"].(string),
			BuildCommand:     worker["build_command"].(string),
			DockerfilePath:   worker["dockerfile_path"].(string),
			Envs:             expandAppEnvs(worker["env"].(*schema.Set).List()),
			InstanceSizeSlug: worker["instance_size_slug"].(string),
			InstanceCount:    int64(worker["instance_count"].(int)),
			SourceDir:        worker["source_dir"].(string),
			EnvironmentSlug:  worker["environment_slug"].(string),
		}

		github := worker["github"].([]interface{})
		if len(github) > 0 {
			s.GitHub = expandAppGitHubSourceSpec(github)
		}

		gitlab := worker["gitlab"].([]interface{})
		if len(gitlab) > 0 {
			s.GitLab = expandAppGitLabSourceSpec(gitlab)
		}

		git := worker["git"].([]interface{})
		if len(git) > 0 {
			s.Git = expandAppGitSourceSpec(git)
		}

		image := worker["image"].([]interface{})
		if len(image) > 0 {
			s.Image = expandAppImageSourceSpec(image)
		}

		appWorkers = append(appWorkers, s)
	}

	return appWorkers
}

func flattenAppSpecWorkers(workers []*godo.AppWorkerSpec) []map[string]interface{} {
	result := make([]map[string]interface{}, len(workers))

	for i, w := range workers {
		r := make(map[string]interface{})

		r["name"] = w.Name
		r["run_command"] = w.RunCommand
		r["build_command"] = w.BuildCommand
		r["github"] = flattenAppGitHubSourceSpec(w.GitHub)
		r["gitlab"] = flattenAppGitLabSourceSpec(w.GitLab)
		r["git"] = flattenAppGitSourceSpec(w.Git)
		r["image"] = flattenAppImageSourceSpec(w.Image)
		r["dockerfile_path"] = w.DockerfilePath
		r["env"] = flattenAppEnvs(w.Envs)
		r["instance_size_slug"] = w.InstanceSizeSlug
		r["instance_count"] = int(w.InstanceCount)
		r["source_dir"] = w.SourceDir
		r["environment_slug"] = w.EnvironmentSlug

		result[i] = r
	}

	return result
}

func expandAppSpecJobs(config []interface{}) []*godo.AppJobSpec {
	appJobs := make([]*godo.AppJobSpec, 0, len(config))

	for _, rawJob := range config {
		job := rawJob.(map[string]interface{})

		s := &godo.AppJobSpec{
			Name:             job["name"].(string),
			RunCommand:       job["run_command"].(string),
			BuildCommand:     job["build_command"].(string),
			DockerfilePath:   job["dockerfile_path"].(string),
			Envs:             expandAppEnvs(job["env"].(*schema.Set).List()),
			InstanceSizeSlug: job["instance_size_slug"].(string),
			InstanceCount:    int64(job["instance_count"].(int)),
			SourceDir:        job["source_dir"].(string),
			EnvironmentSlug:  job["environment_slug"].(string),
			Kind:             godo.AppJobSpecKind(job["kind"].(string)),
		}

		github := job["github"].([]interface{})
		if len(github) > 0 {
			s.GitHub = expandAppGitHubSourceSpec(github)
		}

		gitlab := job["gitlab"].([]interface{})
		if len(gitlab) > 0 {
			s.GitLab = expandAppGitLabSourceSpec(gitlab)
		}

		git := job["git"].([]interface{})
		if len(git) > 0 {
			s.Git = expandAppGitSourceSpec(git)
		}

		image := job["image"].([]interface{})
		if len(image) > 0 {
			s.Image = expandAppImageSourceSpec(image)
		}

		appJobs = append(appJobs, s)
	}

	return appJobs
}

func flattenAppSpecJobs(jobs []*godo.AppJobSpec) []map[string]interface{} {
	result := make([]map[string]interface{}, len(jobs))

	for i, j := range jobs {
		r := make(map[string]interface{})

		r["name"] = j.Name
		r["run_command"] = j.RunCommand
		r["build_command"] = j.BuildCommand
		r["github"] = flattenAppGitHubSourceSpec(j.GitHub)
		r["gitlab"] = flattenAppGitLabSourceSpec(j.GitLab)
		r["git"] = flattenAppGitSourceSpec(j.Git)
		r["image"] = flattenAppImageSourceSpec(j.Image)
		r["dockerfile_path"] = j.DockerfilePath
		r["env"] = flattenAppEnvs(j.Envs)
		r["instance_size_slug"] = j.InstanceSizeSlug
		r["instance_count"] = int(j.InstanceCount)
		r["source_dir"] = j.SourceDir
		r["environment_slug"] = j.EnvironmentSlug
		r["kind"] = string(j.Kind)

		result[i] = r
	}

	return result
}

func expandAppSpecDatabases(config []interface{}) []*godo.AppDatabaseSpec {
	appDatabases := make([]*godo.AppDatabaseSpec, 0, len(config))

	for _, rawDatabase := range config {
		db := rawDatabase.(map[string]interface{})

		s := &godo.AppDatabaseSpec{
			Name:        db["name"].(string),
			Engine:      godo.AppDatabaseSpecEngine(db["engine"].(string)),
			Version:     db["version"].(string),
			Production:  db["production"].(bool),
			ClusterName: db["cluster_name"].(string),
			DBName:      db["db_name"].(string),
			DBUser:      db["db_user"].(string),
		}

		appDatabases = append(appDatabases, s)
	}

	return appDatabases
}

func flattenAppSpecDatabases(databases []*godo.AppDatabaseSpec) []map[string]interface{} {
	result := make([]map[string]interface{}, len(databases))

	for i, db := range databases {
		r := make(map[string]interface{})

		r["name"] = db.Name
		r["engine"] = db.Engine
		r["version"] = db.Version
		r["production"] = db.Production
		r["cluster_name"] = db.ClusterName
		r["db_name"] = db.DBName
		r["db_user"] = db.DBUser

		result[i] = r
	}

	return result
}
