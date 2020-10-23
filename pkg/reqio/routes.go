package reqio

import "net/http"

type Route struct {
	Pattern   string
	Handler   http.Handler
	MemberOnly bool
}

type manager struct {
	routes []Route
}

type IRoutes interface {
	GetRoutes() []Route
	GetMemberOnlyRoutes() (routes []Route)
	GetPublicRoutes() (routes []Route)
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

func (m *manager) GetMemberOnlyRoutes() (routes []Route) {
	routes = make([]Route, 0)
	for _, route := range m.routes {
		if !route.MemberOnly {
			continue
		}
		routes = append(routes, route)
	}
	return
}

func (m *manager) GetPublicRoutes() (routes []Route) {
	routes = make([]Route, 0)
	for _, route := range m.routes {
		if route.MemberOnly {
			continue
		}
		routes = append(routes, route)
	}
	return
}