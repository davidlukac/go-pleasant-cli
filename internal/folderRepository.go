package internal

import (
	"fmt"
	"github.com/davidlukac/go-pleasant-vault-client/pkg/client"
	log "github.com/sirupsen/logrus"
)

func GetRoot() string {
	vault := GetVault()
	return vault.GetRootFolder()
}

func FoldersExist(pathParts []string, parentId string) bool {
	vault := GetVault()

	if len(parentId) == 0 {
		parentId = vault.GetRootFolder()
	}

	parent := vault.GetFolder(parentId)

	for _, f := range parent.Children {
		if f.Name == pathParts[0] {
			if len(pathParts) > 1 {
				return FoldersExist(pathParts[1:], f.Id)
			} else {
				return true
			}
		}
	}

	return false
}

// GetFolder
// Return client.Folder object for given name and parent folder ID.
func GetFolder(name string, parentId string) *client.Folder {
	vault := GetVault()

	if len(parentId) == 0 {
		parentId = vault.GetRootFolder()
	}

	parent := vault.GetFolder(parentId)

	for _, f := range parent.Children {
		if f.Name == name {
			return &f
		}
	}

	return nil
}

// CreateFolders
// Create non-existing parts of the path delimited by '/' from provided root ID (root by default).
func CreateFolders(pathParts []string, parentId string) {
	vault := GetVault()

	if len(parentId) == 0 {
		parentId = vault.GetRootFolder()
	}

	if len(pathParts) > 0 {
		if false == FoldersExist(pathParts[0:1], parentId) {
			folder := client.Folder{Name: pathParts[0], ParentId: parentId}
			newFolder := vault.CreateFolder(&folder)
			log.Info(fmt.Sprintf("%s -> %s", newFolder.Name, newFolder.Id))
		}
		if len(pathParts) > 1 {
			folder := GetFolder(pathParts[0], parentId)
			CreateFolders(pathParts[1:], folder.Id)
		}
	}

	return
}
