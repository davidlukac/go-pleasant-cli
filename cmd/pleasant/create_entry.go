// Package pleasant /*
package pleasant

/*
Copyright © 2021 David Lukac <david.lukac@users.noreply.github.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"fmt"
	"github.com/davidlukac/go-pleasant-cli/internal"
	"github.com/davidlukac/go-pleasant-vault-client/pkg/client"
	"github.com/dchest/uniuri"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"syscall"
)

var UsernameFlag string
var ParentFlag string
var ParentIsPathFlag = false
var RandomPasswordFlag = false

// createEntryCmd represents the entry command
var createEntryCmd = &cobra.Command{
	Use:   "entry",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var password string
		var parentId string

		log.Infoln(fmt.Sprintf("Creating password entry with name %s", args[0]))

		name := args[0]
		rootId := internal.GetRoot()

		if len(UsernameFlag) == 0 {
			UsernameFlag = name
		}

		if ParentIsPathFlag {
			pathParts := strings.Split(strings.Trim(ParentFlag, "/"), "/")
			if false == internal.FoldersExist(pathParts, rootId) {
				internal.CreateFolders(pathParts, rootId)
			}
			parentId = internal.GetFolderIdFromPath(pathParts)
		} else {
			if len(ParentFlag) == 0 {
				parentId = rootId
			} else {
				parentId = ParentFlag
			}
		}

		if internal.EntryExistsByName(name, parentId) {
			fmt.Printf("Entry with name %s already exists in folder %s. Please choose other name or use command line flags to override this check.", name, parentId)
			return
		}

		if RandomPasswordFlag {
			password = uniuri.NewLen(20)
		} else {
			fmt.Println("Please enter password:")
			passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				panic(err)
			}
			password = string(passwordBytes)
		}

		newEntry := internal.CreateEntry(&client.Secret{
			Name:             name,
			Username:         UsernameFlag,
			Password:         password,
			CustomUserFields: map[string]string{},
			GroupId:          parentId,
		})

		log.Info(fmt.Sprintf("Created new entry with ID %s", newEntry.Id))
	},
}

func init() {
	createCmd.AddCommand(createEntryCmd)

	createEntryCmd.Flags().StringVarP(&UsernameFlag, "username", "u", "", "Username")
	createEntryCmd.Flags().StringVarP(&ParentFlag, "parent", "p", "", "ID of the parent folder.")
	createEntryCmd.Flags().BoolVarP(&ParentIsPathFlag, "parent-folder", "f", false, "Parse parent as path, e.g. /foo/bar")
	createEntryCmd.Flags().BoolVarP(&RandomPasswordFlag, "random-password", "r", false, "Generate random password for the new entry.")
}
