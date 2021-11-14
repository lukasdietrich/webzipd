package routing

import (
	"net/http"
	"path"
	"strings"
)

func getStrategyByRoutemode(routemode string) func(*http.Request) string {
	switch routemode {
	case "hostname":
		return routeByHostname

	case "path":
		return routeByPath
	}

	return nil
}

func routeByHostname(r *http.Request) string {
	return trimNamespace(r.Host)
}

func routeByPath(r *http.Request) string {
	namespace, filename := path.Split(r.URL.Path)
	r.URL.Path = filename

	return trimNamespace(namespace)
}

func trimNamespace(namespace string) string {
	return strings.Trim(namespace, "/")
}
