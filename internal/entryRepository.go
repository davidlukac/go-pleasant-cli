package internal

import (
	"encoding/json"
	"fmt"
	"github.com/davidlukac/go-pleasant-vault-client/pkg/client"
	"regexp"
	"strings"
)

// GetEntry
// Get secret entry for given ID.
func GetEntry(id string) *client.Secret {
	vault := GetVault()
	secret := vault.GetSecret(id)
	return &secret
}

// GetEntryWithPassword return Entry enriched with password (which is by default not present).
func GetEntryWithPassword(id string) *client.Secret {
	vault := GetVault()
	entry := vault.GetEntryWithPassword(id)

	return entry
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
			return c.ID
		}
	}

	return ""
}

// PatchEntry modifies existing entry with provided JSON fields and returns updated entry.
func PatchEntry(id string, jsonPatch string) *client.Secret {
	vault := GetVault()

	entry := GetEntry(id)
	if entry == nil {
		panic(fmt.Sprintf("Entry %s doesn't exist!", id))
	}

	err := json.Unmarshal([]byte(jsonPatch), &entry)
	if err != nil {
		panic(fmt.Sprintf("Failed to read provided string into JSON and apply it to an Entry! %s", err))
	}

	vault.PatchEntry(id, jsonPatch)

	entry = GetEntry(id)

	return entry
}

// GetEntryIdForPath returns Entry ID for a provided path, i.e. /foo/bar/entry-name.
func GetEntryIdForPath(path string) string {
	originalPath := path
	vault := GetVault()
	root := vault.GetRootFolder()

	// Sanitize the provided path.
	path = strings.Trim(path, "/")
	re := regexp.MustCompile("[ ]*/[ ]*")
	path = re.ReplaceAllString(path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 1 {
		panic(fmt.Sprintf("Provided path %s is invalid! Expecting /foo/bar/entry-name.", originalPath))
	}

	entryName := pathParts[len(pathParts)-1]
	folders := pathParts[:len(pathParts)-1]
	if false == FoldersExist(folders, root) {
		panic(fmt.Sprintf("Folders on provided path %s don't exist!", originalPath))
	}

	parentId := GetFolderIdFromPath(folders)
	if false == EntryExistsByName(entryName, parentId) {
		panic(fmt.Sprintf("Entry for provided path %s doesn't exist!", originalPath))
	}

	entryId := GetEntryIdByName(entryName, parentId)

	return entryId
}

// GetLinkToEntry returns link to given Entry in the Web UI.
func GetLinkToEntry(id string) string {
	return fmt.Sprintf("%s/WebClient/Main?itemId=%s", strings.TrimSuffix(GetVault().URL, "/"), id)
}
