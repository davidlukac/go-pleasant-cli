module github.com/davidlukac/go-pleasant-cli

go 1.16

require (
	github.com/davidlukac/go-pleasant-vault-client v0.0.0-20210727160429-5df1ebb67bd6
	github.com/dchest/uniuri v0.0.0-20200228104902-7aecb25e1fe5
	github.com/mikefarah/yq/v4 v4.10.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	golang.org/x/crypto v0.17.0
	gopkg.in/yaml.v2 v2.4.0
)

// Local development:
//replace github.com/davidlukac/go-pleasant-vault-client v0.0.0-20210727160429-5df1ebb67bd6 => ../go-pleasant-vault-client
