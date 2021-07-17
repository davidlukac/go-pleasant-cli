package internal

import (
	"fmt"
	"github.com/davidlukac/go-pleasant-vault-client/pkg/client"
	"github.com/dchest/uniuri"
	"reflect"
	"strings"
)

const RandomToken = "$RANDOM$"
const ReferenceQueryStartToken = "$REF"
const ReferenceQueryEndToken = "REF$"

// RandomizeData iterates of StringData in the Secret object and replaces 'random' tokens with random strings and
// return list of modified keys.
func RandomizeData(s *KubernetesOpaqueSecret) []string {
	var updatedKeys []string

	for k, v := range s.StringData {
		if v == RandomToken {
			v = uniuri.NewLen(20)
			s.StringData[k] = v
			updatedKeys = append(updatedKeys, k)
		}
	}

	return updatedKeys
}

// ResolveReferences takes an Opaque Secret object, presumably loaded from a YAML file. This might contain references
// to existing values in provided Entry objects.
// Example
// ...
// stringData:
//   another-password-field: $REF.ENTRY.Password.REF$
//
// will look up Password field in the Entry object.
func ResolveReferences(s *KubernetesOpaqueSecret, e *client.Secret) []string {
	var updatedKeys []string

	for k, v := range s.StringData {
		if strings.HasPrefix(v, ReferenceQueryStartToken) && strings.HasSuffix(v, ReferenceQueryEndToken) {
			entityType, query := parseReferenceQuery(v)
			if entityType == reflect.TypeOf(client.Secret{}) {
				// Get value of the referenced field.
				rValue := reflect.ValueOf(e)
				rField := reflect.Indirect(rValue).FieldByName(query)
				s.StringData[k] = string(rField.String())

				updatedKeys = append(updatedKeys, k)
			}
		}
	}

	return updatedKeys
}

// Parse reference query - $REF.ENTITY.QUERY.REF$ - and returns entity type and query part.
func parseReferenceQuery(queryString string) (reflect.Type, string) {
	if false == strings.HasPrefix(queryString, ReferenceQueryStartToken) || false == strings.HasSuffix(queryString, ReferenceQueryEndToken) {
		panic(fmt.Sprintf("Invalid queryString string %s!", queryString))
	}

	queryParts := strings.Split(queryString, ".")
	if len(queryParts) < 4 {
		panic(fmt.Sprintf("Invalid queryString string %s!", queryString))
	}

	entityTypeStr := queryParts[1]
	query := queryParts[2]

	if strings.ToLower(entityTypeStr) == "entry" {
		return reflect.TypeOf(client.Secret{}), query
	} else {
		panic(fmt.Sprintf("Unsupported entity type %s!", entityTypeStr))
	}
}
