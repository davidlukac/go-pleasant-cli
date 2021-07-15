package internal

import "github.com/davidlukac/go-pleasant-vault-client/pkg/client"

// GetEntry
// Get secret entry for given ID.
func GetEntry(id string) *client.Secret {
	vault := GetVault()
	secret := vault.GetSecret(id)
	return &secret
}

// CreateEntry
// Create new Password Server entry
func CreateEntry(entry *client.Secret) *client.Secret {
	vault := GetVault()
	newEntry := vault.CreateEntry(entry)
	return newEntry
}

// EntryExistsByName checks whether an entry
func EntryExistsByName(name string, parentId string) bool {
	vault := GetVault()
	folder := vault.GetFolder(parentId)

	for _, c := range folder.Credentials {
		if c.Name == name {
			return true
		}
	}

	return false
}

// GetEntryIdByName returns entry ID if it exists, otherwise empty string.
func GetEntryIdByName(name string, parentId string) string {
	vault := GetVault()
	folder := vault.GetFolder(parentId)

	for _, c := range folder.Credentials {
		if c.Name == name {
			return c.Id
		}
	}

	return ""
}
