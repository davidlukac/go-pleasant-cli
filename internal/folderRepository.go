package internal

import (
	"errors"
	"fmt"
	"github.com/davidlukac/go-pleasant-vault-client/pkg/client"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

// GetRoot
// Return UUID of the root folder.
func GetRoot() string {
	vault := GetVault()
	return vault.GetRootFolder()
}

// FoldersExist checks for existence of provided folders in an array in provided parent. If parent is not provided,
// root is assumed. Path parts represent folders in each other.
func FoldersExist(pathParts []string, parentId string) bool {
	vault := GetVault()

	if len(parentId) == 0 {
		parentId = vault.GetRootFolder()
	}

	parent := vault.GetFolder(parentId)

	for _, f := range parent.Children {
		if f.Name == pathParts[0] {
			if len(pathParts) > 1 {
				return FoldersExist(pathParts[1:], f.ID)
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

func GetFolderForId(id string) *client.Folder {
	vault := GetVault()

	return vault.GetFolder(id)
}

// CreateFolders
// Create non-existing parts of the path delimited by '/' from provided root ID (root by default).
func CreateFolders(pathParts []string, parentId string) {
	var currentFolder *client.Folder
	var folderId string

	vault := GetVault()

	if len(parentId) == 0 {
		parentId = vault.GetRootFolder()
	}

	if len(pathParts) > 0 {
		if FoldersExist(pathParts[0:1], parentId) {
			currentFolder = GetFolder(pathParts[0], parentId)
		} else {
			folderToCreate := client.Folder{Name: pathParts[0], ParentID: parentId}
			currentFolder = vault.CreateFolder(&folderToCreate)
			log.Info(fmt.Sprintf("%s -> %s", currentFolder.Name, currentFolder.ID))
		}

		folderId = currentFolder.ID

		if len(pathParts) > 1 {
			CreateFolders(pathParts[1:], folderId)
		}
	}

	return
}

// GetFolderIdFromPath returns folder ID of the last folder in provided path parts. The path is assumed to start at root.
func GetFolderIdFromPath(pathParts []string) string {
	rootId := GetRoot()
	exists := FoldersExist(pathParts, rootId)

	if exists == false {
		panic(fmt.Sprintf("Provided path /%s doesn't exists - can't retrieve the ID!", strings.Join(pathParts, "/")))
	}

	parentId := rootId
	for _, f := range pathParts {
		folder := GetFolder(f, parentId)
		parentId = folder.ID
	}

	return parentId
}

// ParsePath parses input string that should represent a path in the Password Server structure. The path should be along
// following style of a path:
// - foo/bar/baz
// - /foo/bar/baz/
// - foo / bar / baz
// The result is an array of path parts or an error if parsing was unsuccessful.
func ParsePath(path string) ([]string, error) {
	var err error = nil
	var pathParts []string = nil

	originalPath := path

	// Remove any leading and trailing spaces and slashes.
	path = strings.Trim(path, " /")
	// Remove spaces before and after '/' separators.
	re := regexp.MustCompile("[ ]*/[ ]*")
	path = re.ReplaceAllString(path, "/")
	pathParts = strings.Split(path, "/")

	if len(pathParts) < 1 {
		err = errors.New(fmt.Sprintf("Provided path %s is invalid! Expecting something along of '/foo/bar/baz'.", originalPath))
		pathParts = nil
	}

	return pathParts, err
}
