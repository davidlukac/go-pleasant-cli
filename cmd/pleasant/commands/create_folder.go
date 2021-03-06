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
	"strings"

	"github.com/spf13/cobra"
)

// createFolderCmd represents the folder command
var createFolderCmd = &cobra.Command{
	Use:   "folder",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("create folder called with %s\n", args[0])
		pathParts := strings.Split(args[0], "/")
		internal.CreateFolders(pathParts, "")
		lastFolderId := internal.GetFolderIdFromPath(pathParts)
		fmt.Println(lastFolderId)
	},
}

func init() {
	createCmd.AddCommand(createFolderCmd)
}
