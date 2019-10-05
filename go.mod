module github.com/terraform-providers/terraform-provider-digitalocean

require (
	contrib.go.opencensus.io/exporter/ocagent v0.6.0 // indirect
	github.com/aws/aws-sdk-go v1.22.0
	github.com/digitalocean/godo v1.21.0
	github.com/gophercloud/gophercloud v0.4.0 // indirect
	github.com/hashicorp/terraform v0.12.8
	github.com/terraform-providers/terraform-provider-kubernetes v1.9.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/yaml.v2 v2.2.2
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

replace github.com/gophercloud/gophercloud v0.0.0-20190523203818-4885c347dcf4 => github.com/gophercloud/gophercloud v0.0.0-20190523203039-4885c347dcf4

replace github.com/keybase/go-crypto v0.0.0-20190523171820-b785b22cc757 => github.com/keybase/go-crypto v0.0.0-20190416182011-b785b22cc757

replace github.com/go-critic/go-critic v0.0.0-20181204210945-ee9bf5809ead => github.com/go-critic/go-critic v0.3.5-0.20190210220443-ee9bf5809ead

replace golang.org/x/tools v0.0.0-20190314010720-f0bfdbff1f9c => golang.org/x/tools v0.0.0-20190315214010-f0bfdbff1f9c

replace github.com/golangci/errcheck v0.0.0-20181003203344-ef45e06d44b6 => github.com/golangci/errcheck v0.0.0-20181223084120-ef45e06d44b6

replace github.com/golangci/go-tools v0.0.0-20180109140146-af6baa5dc196 => github.com/golangci/go-tools v0.0.0-20190318060251-af6baa5dc196

replace github.com/golangci/gofmt v0.0.0-20181105071733-0b8337e80d98 => github.com/golangci/gofmt v0.0.0-20181222123516-0b8337e80d98

replace github.com/golangci/gosec v0.0.0-20180901114220-66fb7fc33547 => github.com/golangci/gosec v0.0.0-20190211064107-66fb7fc33547

replace mvdan.cc/unparam v0.0.0-20190124213536-fbb59629db34 => mvdan.cc/unparam v0.0.0-20190209190245-fbb59629db34

go 1.13
