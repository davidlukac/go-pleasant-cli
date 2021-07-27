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
	log "github.com/sirupsen/logrus"
	"strings"

	"github.com/spf13/cobra"
)

var targetFolderFlag string
var targetFolderIsPathFlag = false
var createTargetFolderFlag = false

// duplicateEntryCmd represents the `password-server duplicate entry` command
var duplicateEntryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Create a copy of an Entry (secret).",
	Long: `Take provided UUID of an Entry and make a copy. If no further information is provided, the entry is copied
within the same folder, with a suffix '- Copy' in the name. If a new parent folder UUID is provided, the entry is copied
into the new parent, without the suffix.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln(fmt.Sprintf("... duplicate entry called with %s", args))

		originEntryID := args[0]
		entry := internal.GetEntryWithPassword(originEntryID)
		attachments := internal.GetEntryAttachments(originEntryID)

		// Set up the new entry:
		if len(targetFolderFlag) > 0 {
			// Either set the target folder as ID for the given path OR from the input.
			if targetFolderIsPathFlag {
				var targetFolderID string

				root := internal.GetRoot()
				// Parse the target folder.
				pathParts, err := internal.ParsePath(targetFolderFlag)
				if err != nil {
					panic(err)
				}

				targetFolderExists := internal.FoldersExist(pathParts, root)

				if targetFolderExists {
					targetFolderID = internal.GetFolderIdFromPath(pathParts)
				} else {
					if createTargetFolderFlag {
						internal.CreateFolders(pathParts, root)
						targetFolderID = internal.GetFolderIdFromPath(pathParts)
					} else {
						panic(fmt.Sprintf("Target folder /%s doesn't exist!", strings.Join(pathParts, "/")))
					}
				}

				entry.GroupID = targetFolderID
			} else {
				entry.GroupID = targetFolderFlag

			}
		} else {
			entry.Name = fmt.Sprintf("%s - Copy", entry.Name)
		}

		// Reset the Entry ID, because we're about to clone it.
		entry.ID = ""
		entry = internal.CreateEntry(entry)
		for _, attachment := range attachments {
			attachment.ID = ""
			attachment.EntryID = entry.ID
			internal.UploadAttachment(entry.ID, attachment)
		}

		log.Infoln(fmt.Sprintf("Entry copied as %s", entry.ID))
	},
}

func init() {
	duplicateCmd.AddCommand(duplicateEntryCmd)

	duplicateEntryCmd.Flags().StringVarP(&targetFolderFlag, "target-folder", "t", "", "Target folder UUID for the entry duplicate.")
	duplicateEntryCmd.Flags().BoolVarP(&targetFolderIsPathFlag, "target-folder-path", "p", false, "Parse the target folder flag as path.")
	duplicateEntryCmd.Flags().BoolVarP(&createTargetFolderFlag, "create-target-folder", "c", false,
		"When used together with 'target-folder-path' flag, the target folder is created, if it doesn't exist (or it's parts).")
}
