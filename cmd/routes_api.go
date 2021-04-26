package cmd

import (
	"github.com/evorts/feednomity/handler"
	"github.com/evorts/feednomity/pkg/middleware"
	"github.com/evorts/feednomity/pkg/reqio"
	"net/http"
)

func routesApi(lib *library) []reqio.Route {
	apiRoutes := []reqio.Route{
		{
			Pattern: "/api/reload",
			Handler: middleware.WithProtection(middleware.ProtectionLib{
				Acl:  lib.acl,
				Sm:   lib.session,
				View: lib.view,
			}, middleware.ProtectionArgs{
				Path:           "/api/reload",
				Method:         http.MethodPost,
				AllowedMethods: lib.config.GetConfig().App.Cors.AllowedMethods,
				AllowedOrigins: lib.config.GetConfig().App.Cors.AllowedOrigins,
				RenderType:     "json",
			}, http.HandlerFunc(handler.ApiReload)),
			AdminOnly: true,
		},
		{
			Pattern: "/api/login",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiLogin),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: false,
		},
		{
			Pattern: "/api/feedbacks",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiFeedbackSubmission),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/360/submission",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.Api360Submission),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"db":     lib.db,
							"sm":     lib.session,
							"hash":   lib.hash,
						},
					),
				),
			),
			AdminOnly: false,
		},
	}
	apiRoutes = append(apiRoutes, routesApiDistribution(lib)...)
	apiRoutes = append(apiRoutes, routesApiUsers(lib)...)
	apiRoutes = append(apiRoutes, routesApiLink(lib)...)
	apiRoutes = append(apiRoutes, routesApiGroups(lib)...)
	apiRoutes = append(apiRoutes, routesApiQuestions(lib)...)
	return apiRoutes
}

func routesApiDistribution(lib *library) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/distribution/publish",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiLinksBlast),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
							"mm":     lib.mm,
						},
					),
				),
			),
			AdminOnly: true,
		},
	}
}

func routesApiUsers(lib *library) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/users/list",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiUsersList),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/users/create",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiUserCreate),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/users/update",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiUserUpdate),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/users/delete",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiUsersDelete),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
	}
}

func routesApiLink(lib *library) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/links/list",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiLinksList),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/links/create",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiLinksCreate),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
							"aes":    lib.aes,
							"cfg":    lib.config,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/links/update",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiLinkUpdate),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
							"aes":    lib.aes,
							"cfg":    lib.config,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/links/delete",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiLinksDelete),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
	}
}

func routesApiGroups(lib *library) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/groups/list",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiGroupsList),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/groups/create",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiGroupsCreate),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/groups/update",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiGroupUpdate),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/groups/remove",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiGroupsDelete),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
	}
}

func routesApiQuestions(lib *library) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/questions",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiQuestions),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/questions/create",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiQuestionCreate),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/questions/update",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiQuestionUpdate),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/questions/remove",
			Handler: middleware.WithCors(
				lib.view,
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(handler.ApiQuestionRemove),
						map[string]interface{}{
							"logger": lib.logger,
							"view":   lib.view,
							"sm":     lib.session,
							"hash":   lib.hash,
							"db":     lib.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
	}
}
