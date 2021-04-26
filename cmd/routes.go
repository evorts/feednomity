package cmd

import (
	"github.com/evorts/feednomity/handler"
	"github.com/evorts/feednomity/pkg/middleware"
	"github.com/evorts/feednomity/pkg/reqio"
	"net/http"
)

func routes(o *http.ServeMux, lib *library) {
	// serving assets
	fs := http.FileServer(http.Dir(lib.config.GetConfig().App.AssetDirectory))
	o.Handle("/assets/", http.StripPrefix("/assets", fs))
	// serving pages
	routing := []reqio.Route{
		{
			Pattern: "/ping",
			Handler: middleware.WithMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(handler.Ping),
					map[string]interface{}{
						"view": lib.view,
					},
				),
			),
			AdminOnly: false,
		},
	}
	routing = append(routing, routesWebDashboard(lib)...)
	routing = append(routing, routesWebAssessments(lib)...)
	routing = append(routing, routesApi(lib)...)

	reqio.NewRoutes(routing).ExecRoutes(o)
}


