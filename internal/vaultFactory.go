package internal

import "github.com/davidlukac/go-pleasant-vault-client/pkg/client"

func GetVault() client.Vault {
	return client.Vault{
		URL:      "https://72506700.pleasantpassworddemo.com/",
		Username: "admin",
		Password: "admin123",
	}
}
