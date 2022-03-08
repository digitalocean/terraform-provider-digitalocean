module github.com/digitalocean/terraform-provider-digitalocean

require (
	github.com/aws/aws-sdk-go v1.42.18
	github.com/digitalocean/godo v1.73.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/awspolicyequivalence v1.5.0
	github.com/hashicorp/go-uuid v1.0.2
	github.com/hashicorp/go-version v1.3.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.10.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/hashstructure/v2 v2.0.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

replace github.com/keybase/go-crypto v0.0.0-20190523171820-b785b22cc757 => github.com/keybase/go-crypto v0.0.0-20190416182011-b785b22cc757

go 1.16
