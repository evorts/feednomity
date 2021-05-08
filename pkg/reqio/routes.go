package reqio

import "net/http"

type Route struct {
	Pattern   string
	Handler   http.Handler
}

type manager struct {
	routes []Route
}

type IRoutes interface {
	GetRoutes() []Route
	ExecRoutes(server *http.ServeMux)
}

func NewRoutes(routes []Route) IRoutes {
	return &manager{routes: routes}
}

func (m *manager) ExecRoutes(server *http.ServeMux) {
	for _, route := range m.GetRoutes() {
		server.Handle(route.Pattern, route.Handler)
	}
}

func (m *manager) GetRoutes() []Route {
	return m.routes
}