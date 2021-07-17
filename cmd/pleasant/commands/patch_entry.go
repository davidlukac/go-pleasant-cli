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

var JsonFlag string
var FromKubernetesOpaqueSecretYamlPathFlag string
var RandomizeFlag = false
var UpdateKubernetesOpaqueSecretYamlFileFlag = false
var EntryIsPathNameFolder = false

// patchEntryCmd is CLI command for patching an Entry.
var patchEntryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Patch an existing Entry",
	Long: `Patch an existing Entry with a JSON.
Entry can by referenced by it's ID or by a path to it. Patching also supports a Kubernetes Opaque Secret YAML file as 
a source. If the YAML file is used, it may contain various tokens, e.g. for randomization of values or referencing
existing values.
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var entry *client.Secret
		var id string

		log.Infoln(fmt.Sprintf("Patch entry called with %s.", args[0]))

		if EntryIsPathNameFolder {
			id = internal.GetEntryIdForPath(args[0])
		} else {
			id = args[0]
		}

		entry = internal.GetEntryWithPassword(id)

		if len(FromKubernetesOpaqueSecretYamlPathFlag) > 0 {
			var updatedKeys []string

			kubeSecret := internal.ReadKubernetesOpaqueSecret(FromKubernetesOpaqueSecretYamlPathFlag)

			updatedKeys = internal.ResolveReferences(&kubeSecret, entry)

			if RandomizeFlag {
				// Randomize the tokens in Custom fields.
				updatedKeys = append(updatedKeys, internal.RandomizeData(&kubeSecret)...)
			}

			if UpdateKubernetesOpaqueSecretYamlFileFlag {
				// Also update the source YAML file with the random values.
				internal.UpdateStringData(FromKubernetesOpaqueSecretYamlPathFlag, &kubeSecret, updatedKeys)
			}

			patch := internal.OpaqueSecretToPatch(&kubeSecret)
			log.Debugln(fmt.Sprintf("Patch read from the file: %s as %s", kubeSecret, patch))
			JsonFlag = patch
		}

		entry = internal.PatchEntry(id, JsonFlag)

		log.Infoln(fmt.Sprintf("Entry %s updated: %s", id, entry))
	},
}

func init() {
	patchCmd.AddCommand(patchEntryCmd)

	patchEntryCmd.Flags().StringVarP(&JsonFlag, "json", "j", "{}", "JSON-formatted patch of an entry.")
	patchEntryCmd.Flags().StringVarP(&FromKubernetesOpaqueSecretYamlPathFlag, "from-k8s-opaque-yaml-file", "y", "", "Construct the patch from a provided path, which is expected to be a Kubernetes Opaque Secret YAML file.")
	patchEntryCmd.Flags().BoolVarP(&RandomizeFlag, "randomize", "r", false, fmt.Sprintf("Randomize values where '%s' is found.", internal.RandomToken))
	patchEntryCmd.Flags().BoolVarP(&UpdateKubernetesOpaqueSecretYamlFileFlag, "update-k8s-opaque-yaml-file", "u", false, "When used together with provided path to K8s Opaque Secret YAML file and `randomize`, the random tokens in the file are updated with generated values.")
	patchEntryCmd.Flags().BoolVarP(&EntryIsPathNameFolder, "folder", "f", false, "When set to true, entry is not parsed as UUID but as folder path and entry name instead.")
}
