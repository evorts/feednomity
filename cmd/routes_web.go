package cmd

import (
	"github.com/evorts/feednomity/handler"
	"github.com/evorts/feednomity/pkg/middleware"
	"github.com/evorts/feednomity/pkg/reqio"
	"net/http"
)

func routesWebDashboard(lib *library) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/adm",
			Handler: middleware.WithMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(handler.AdminGate),
					map[string]interface{}{
						"logger": lib.logger,
						"sm":     lib.session,
					},
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/adm/dashboard",
			Handler: middleware.WithProtection(middleware.ProtectionLib{
				Acl:  lib.acl,
				Sm:   lib.session,
				View: lib.view,
			}, middleware.ProtectionArgs{
				Path:           "/adm/dashboard",
				Method:         http.MethodGet,
			}, middleware.WithInjection(
				http.HandlerFunc(handler.Dashboard),
				map[string]interface{}{
					"logger": lib.logger,
					"view":   lib.view,
					"sm":     lib.session,
					"db":     lib.db,
				},
			)),
			AdminOnly: true,
		},
		{
			Pattern: "/adm/users",
			Handler: middleware.WithProtection(middleware.ProtectionLib{
				Acl:  lib.acl,
				Sm:   lib.session,
				View: lib.view,
			}, middleware.ProtectionArgs{
				Path:           "/adm/users",
				Method:         http.MethodGet,
			}, middleware.WithInjection(
				http.HandlerFunc(handler.Users),
				map[string]interface{}{
					"logger": lib.logger,
					"view":   lib.view,
					"sm":     lib.session,
					"db":     lib.db,
				},
			)),
			AdminOnly: true,
		},
		{
			Pattern: "/adm/objects",
			Handler: middleware.WithProtection(middleware.ProtectionLib{
				Acl:  lib.acl,
				Sm:   lib.session,
				View: lib.view,
			}, middleware.ProtectionArgs{
				Path:           "/adm/objects",
				Method:         http.MethodGet,
			}, middleware.WithInjection(
				http.HandlerFunc(handler.Objects),
				map[string]interface{}{
					"logger": lib.logger,
					"view":   lib.view,
					"sm":     lib.session,
					"db":     lib.db,
				},
			)),
			AdminOnly: true,
		},
		{
			Pattern: "/adm/login",
			Handler: middleware.WithMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(handler.Login),
					map[string]interface{}{
						"logger": lib.logger,
						"view":   lib.view,
						"sm":     lib.session,
						"hash":   lib.hash,
					},
				),
			),
			AdminOnly: false,
		},
		{
			Pattern: "/adm/logout",
			Handler: middleware.WithMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(handler.Logout),
					map[string]interface{}{
						"sm": lib.session,
					},
				),
			),
			AdminOnly: true,
		},
	}
}

func routesWebAssessments(lib *library) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/page/360/review",
			Handler: middleware.WithMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(handler.Form360),
					map[string]interface{}{
						"logger": lib.logger,
						"view":   lib.view,
						"db":     lib.db,
						"sm":     lib.session,
						"hash":   lib.hash,
					},
				),
			),
			AdminOnly: false,
		},
	}
}
