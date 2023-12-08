package graphql_api

import (
	"fmt"
)

func getURL(host string) string {
	return fmt.Sprintf("%s/v2/query", host)
}
