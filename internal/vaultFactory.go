package internal

import (
	"github.com/davidlukac/go-pleasant-vault-client/pkg/client"
	"github.com/spf13/viper"
)

const PasswordServerUrl = "password_server_url"
const PasswordServerUsername = "password_server_username"
const PasswordServerPassword = "password_server_password"

func GetVault() client.Vault {
	return client.Vault{
		URL:      viper.GetString(PasswordServerUrl),
		Username: viper.GetString(PasswordServerUsername),
		Password: viper.GetString(PasswordServerPassword),
	}
}
