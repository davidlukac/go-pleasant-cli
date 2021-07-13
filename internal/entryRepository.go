package internal

import "github.com/davidlukac/go-pleasant-vault-client/pkg/client"

// GetEntry
// Get secret entry for given ID.
func GetEntry(id string) *client.Secret {
	vault := GetVault()
	secret := vault.GetSecret(id)
	return &secret
}
