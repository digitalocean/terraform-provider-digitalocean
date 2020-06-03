module github.com/terraform-providers/terraform-provider-digitalocean

require (
	contrib.go.opencensus.io/exporter/ocagent v0.6.0 // indirect
	github.com/aws/aws-sdk-go v1.25.4
	github.com/digitalocean/godo v1.36.0
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/terraform-plugin-sdk v1.6.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/terraform-providers/terraform-provider-kubernetes v1.9.1-0.20191018170806-2c80accb5635
	golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	gopkg.in/yaml.v2 v2.2.4
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

replace github.com/keybase/go-crypto v0.0.0-20190523171820-b785b22cc757 => github.com/keybase/go-crypto v0.0.0-20190416182011-b785b22cc757

replace github.com/terraform-providers/terraform-provider-google v2.17.0+incompatible => github.com/terraform-providers/terraform-provider-google v1.20.1-0.20191008212436-363f2d283518

replace github.com/terraform-providers/terraform-provider-aws v2.32.0+incompatible => github.com/terraform-providers/terraform-provider-aws v1.60.1-0.20191010190908-1261a98537f2

replace github.com/terraform-providers/terraform-provider-random v2.2.1+incompatible => github.com/terraform-providers/terraform-provider-random v0.0.0-20190925210718-83518d96ae4f

go 1.13
