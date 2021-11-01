module github.com/digitalocean/terraform-provider-digitalocean

require (
	github.com/aws/aws-sdk-go v1.25.4
	github.com/digitalocean/godo v1.69.1
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/go-version v1.3.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.7.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/hashstructure/v2 v2.0.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b
	golang.org/x/net v0.0.0-20211006190231-62292e806868 // indirect
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

replace github.com/keybase/go-crypto v0.0.0-20190523171820-b785b22cc757 => github.com/keybase/go-crypto v0.0.0-20190416182011-b785b22cc757

go 1.16
