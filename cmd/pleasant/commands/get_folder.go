// Package commands /*
package commands

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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var FromPathFlag = false
var JustFolderIdFlag = false

// getFolderCmd represents the folder command
var getFolderCmd = &cobra.Command{
	Use:   "folder",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln(fmt.Sprintf("get folder called with %s", args[0]))
		if args[0] == "root" {
			rootFolderId := internal.GetRoot()
			fmt.Println(rootFolderId)
		} else {
			var folder *client.Folder

			if FromPathFlag {
				folderId := internal.GetFolderIdFromPath(strings.Split(args[0], "/"))
				folder = internal.GetFolderForId(folderId)
			} else {
				vault := internal.GetVault()
				folder = vault.GetFolder(args[0])
			}

			if JustFolderIdFlag {
				fmt.Println(folder.ID)
			} else {
				fmt.Println(folder)
			}
		}
	},
}

func init() {
	getCmd.AddCommand(getFolderCmd)

	getFolderCmd.Flags().BoolVarP(&FromPathFlag, "from-path", "p", false, "Return the deepest folder for a given path.")
	getFolderCmd.Flags().BoolVarP(&JustFolderIdFlag, "id", "i", false, "Print just folder ID.")
}
