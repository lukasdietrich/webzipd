package routing

import (
	"fmt"
	"net/http"

	"github.com/lukasdietrich/webzipd/internal/render"
)

type mux struct {
	strategy func(*http.Request) (namespace string)
	renderer *render.Renderer
}

func NewMux(renderer *render.Renderer, routemode string) (http.Handler, error) {
	strategy := getStrategyByRoutemode(routemode)
	if strategy == nil {
		return nil, fmt.Errorf("unknown routemode %q", routemode)
	}

	return &mux{renderer: renderer, strategy: strategy}, nil
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.renderer.Render(w, r, m.strategy(r))
}
