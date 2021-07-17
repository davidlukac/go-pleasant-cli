// Package internal
package internal

// Partial implementation of representation of Kubernetes Opaque Secret objects.

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Metadata struct {
	Name string `yaml:"name"`
}

type KubernetesOpaqueSecret struct {
	Type       string            `yaml:"type"`
	Metadata   Metadata          `yaml:"metadata"`
	StringData map[string]string `yaml:"stringData"`
}

// ReadKubernetesOpaqueSecret reads file on given path into a Kubernetes Opaque Secret object.
func ReadKubernetesOpaqueSecret(path string) KubernetesOpaqueSecret {
	var secret KubernetesOpaqueSecret

	_, err := os.Stat(path)
	if err != nil {
		panic(fmt.Sprintf("Provided file doesn't exist! %s", err))
	}

	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to open file %s! %s", path, err))
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(fmt.Sprintf("Failed to close the file %s! %s", path, err))
		}
	}(file)

	secretYaml, err := ioutil.ReadAll(file)
	err = yaml.Unmarshal(secretYaml, &secret)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal YAML file! %s", err))
	}

	if secret.Type != "Opaque" {
		panic(fmt.Sprintf("Secret type is not 'Opaque' (%s instead)!", secret.Type))
	}

	return secret
}

// OpaqueSecretToPatch generate a JSON patch string from Kubernetes Opaque Secret object.
func OpaqueSecretToPatch(s *KubernetesOpaqueSecret) string {
	var patchMap = make(map[string]interface{})

	patchMap["Name"] = s.Metadata.Name
	if dataUsername, ok := s.StringData["username"]; ok {
		patchMap["Username"] = dataUsername
	}
	if dataPassword, ok := s.StringData["password"]; ok {
		patchMap["Password"] = dataPassword
	}

	customUserFields := map[string]string{}

	for k, v := range s.StringData {
		// Skip username and password, because they don't belong to customer fields.
		if k == "username" || k == "password" {
			continue
		}
		customUserFields[k] = v
	}
	patchMap["CustomUserFields"] = customUserFields

	patch, err := json.Marshal(&patchMap)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal %s into JSON! %s", patchMap, err))
	}

	return string(patch)
}
