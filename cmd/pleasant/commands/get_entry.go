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
)

var LinkToEntryFlag = false

// getEntryCmd Get an Entry from the Password Server.
var getEntryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Get an Entry from the Password Server",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var id string
		var entry *client.Secret

		log.Infoln(fmt.Sprintf("get entry called with %s", args[0]))

		if EntryIsPathNameFolder {
			id = internal.GetEntryIdForPath(args[0])
		} else {
			id = args[0]
		}

		if LinkToEntryFlag {
			fmt.Println(internal.GetLinkToEntry(id))
		} else {
			entry = internal.GetEntry(id)
			fmt.Println(entry)
		}

	},
}

func init() {
	getCmd.AddCommand(getEntryCmd)

	getEntryCmd.Flags().BoolVarP(&EntryIsPathNameFolder, "folder", "f", false, "When set to true, entry is not parsed as UUID but as folder path and entry name instead.")
	getEntryCmd.Flags().BoolVarP(&LinkToEntryFlag, "link", "l", false, "Print link to entry instead of the entry itself.")
}
