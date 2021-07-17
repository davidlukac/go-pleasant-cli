package internal

import (
	"fmt"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"strings"
)

// UpdateStringData in YAML file provided by path with string data in updated keys from provided secret.
func UpdateStringData(path string, s *KubernetesOpaqueSecret, updatedKeys []string) {
	var completedSuccessfully bool

	writeInPlaceHandler := yqlib.NewWriteInPlaceHandler(path)
	out, err := writeInPlaceHandler.CreateTempFile()
	if err != nil {
		panic(fmt.Sprintf("Unable to create a tmp file for in-place YAML update! %s", err))
	}

	defer func() {
		writeInPlaceHandler.FinishWriteInPlace(completedSuccessfully)
	}()

	printer := yqlib.NewPrinter(out, false, true, false, 2, true)
	allAtOnceEvaluator := yqlib.NewAllAtOnceEvaluator()

	query := constructUpdateQueryForCustomFields(s, updatedKeys)

	err = allAtOnceEvaluator.EvaluateFiles(query, []string{path}, printer)

	completedSuccessfully = err == nil

	if err != nil || completedSuccessfully == false {
		panic("YQ no matches found")
	}
}

// constructUpdateQueryForCustomFields constructs YQ query from list of updated keys and the secret.
func constructUpdateQueryForCustomFields(s *KubernetesOpaqueSecret, keys []string) string {
	var queryParts []string
	var query string

	for _, key := range keys {
		queryParts = append(queryParts, fmt.Sprintf(`.stringData.%s |= "%s"`, key, s.StringData[key]))
	}

	query = strings.Join(queryParts, " | ")

	return query
}
