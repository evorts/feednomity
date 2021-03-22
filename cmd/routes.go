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
			Pattern: "/adm",
			Handler: middleware.WithMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					middleware.WithProtection(http.HandlerFunc(handler.Dashboard)),
					map[string]interface{}{
						"logger": lib.logger,
						"view":   lib.view,
					},
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/login",
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
			Pattern: "/logout",
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
		{
			Pattern: "/reload",
			Handler: middleware.WithMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(handler.Reload),
					map[string]interface{}{
						"logger": lib.logger,
						"view":   lib.view,
						"config": lib.config,
					},
				),
			),
			AdminOnly: true,
		},
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
		{
			Pattern: "/api/login",
			Handler: middleware.WithCors(
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.LoginAPI),
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
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.FeedbackSubmissionAPI),
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
	routing = append(routing, routesAssessments(lib)...)
	routing = append(routing, routesDistribution(lib)...)
	routing = append(routing, routesObjects(lib)...)
	routing = append(routing, routesLink(lib)...)
	routing = append(routing, routesGroups(lib)...)
	routing = append(routing, routesQuestions(lib)...)

	reqio.NewRoutes(routing).ExecRoutes(o)
}

func routesAssessments(lib *library) []reqio.Route {
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
		{
			Pattern: "/api/360/submission",
			Handler: middleware.WithCors(
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.Review360SubmissionAPI),
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
}

func routesDistribution(lib *library) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/distribution/publish",
			Handler: middleware.WithCors(
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.LinksBlastAPI),
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

func routesObjects(lib *library) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/objects/list",
			Handler: middleware.WithCors(
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(handler.ObjectListAPI),
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
			Pattern: "/api/objects/create",
			Handler: middleware.WithCors(
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.ObjectsCreateAPI),
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
			Pattern: "/api/objects/update",
			Handler: middleware.WithCors(
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(handler.ObjectUpdateAPI),
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
			Pattern: "/api/objects/remove",
			Handler: middleware.WithCors(
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(handler.ObjectsRemoveAPI),
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

func routesLink(lib *library) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/links",
			Handler: middleware.WithCors(
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(handler.LinksAPI),
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
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.LinksCreateAPI),
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
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(handler.LinkUpdateAPI),
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
			Pattern: "/api/links/remove",
			Handler: middleware.WithCors(
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(handler.LinksRemoveAPI),
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

func routesGroups(lib *library) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/groups",
			Handler: middleware.WithCors(
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(handler.GroupsAPI),
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
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.GroupsCreateAPI),
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
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(handler.GroupUpdateAPI),
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
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(handler.GroupsRemoveAPI),
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

func routesQuestions(lib *library) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/questions",
			Handler: middleware.WithCors(
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(handler.QuestionsAPI),
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
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.QuestionCreateAPI),
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
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(handler.QuestionUpdateAPI),
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
				lib.config.GetConfig().App.Cors.AllowedMethods,
				lib.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(handler.QuestionRemoveAPI),
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
