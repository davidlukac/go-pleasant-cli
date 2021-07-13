module github.com/davidlukac/go-pleasant-cli

go 1.16

require (
	github.com/davidlukac/go-pleasant-vault-client v0.0.0-20201104101430-cf5b96afb3ed
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
)

replace github.com/davidlukac/go-pleasant-vault-client v0.0.0-20201104101430-cf5b96afb3ed => ../go-pleasant-vault-client
