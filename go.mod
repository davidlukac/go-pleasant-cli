module github.com/davidlukac/go-pleasant-cli

go 1.16

require (
	github.com/davidlukac/go-pleasant-vault-client v0.0.0-20201104101430-cf5b96afb3ed
	github.com/dchest/uniuri v0.0.0-20200228104902-7aecb25e1fe5
	github.com/mikefarah/yq/v4 v4.10.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/davidlukac/go-pleasant-vault-client v0.0.0-20201104101430-cf5b96afb3ed => ../go-pleasant-vault-client
