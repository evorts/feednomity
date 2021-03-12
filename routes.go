package main

import (
	"github.com/evorts/feednomity/handler"
	"github.com/evorts/feednomity/pkg/middleware"
	"github.com/evorts/feednomity/pkg/reqio"
	"net/http"
)

func routes(o *http.ServeMux, cmd *commands) {
	// serving assets
	fs := http.FileServer(http.Dir(cmd.config.GetConfig().App.AssetDirectory))
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
						"logger": cmd.logger,
						"view":   cmd.view,
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
						"logger": cmd.logger,
						"view":   cmd.view,
						"sm":     cmd.session,
						"hash":   cmd.hash,
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
						"sm": cmd.session,
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
						"logger": cmd.logger,
						"view":   cmd.view,
						"config": cmd.config,
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
						"view": cmd.view,
					},
				),
			),
			AdminOnly: false,
		},
		{
			Pattern: "/api/login",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.LoginAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: false,
		},
		{
			Pattern: "/api/feedbacks",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.FeedbackSubmissionAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
	}
	routing = append(routing, routesAssessments(cmd)...)
	routing = append(routing, routesObjects(cmd)...)
	routing = append(routing, routesLink(cmd)...)
	routing = append(routing, routesGroups(cmd)...)
	routing = append(routing, routesQuestions(cmd)...)

	reqio.NewRoutes(routing).ExecRoutes(o)
}

func routesAssessments(cmd *commands) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/page/360/review",
			Handler: middleware.WithMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(handler.Form360),
					map[string]interface{}{
						"logger": cmd.logger,
						"view":   cmd.view,
						"db":     cmd.db,
						"sm":     cmd.session,
						"hash":   cmd.hash,
					},
				),
			),
			AdminOnly: false,
		},
		{
			Pattern: "/api/360/submission",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.Review360SubmissionAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"db":     cmd.db,
							"sm":     cmd.session,
							"hash":   cmd.hash,
						},
					),
				),
			),
			AdminOnly: false,
		},
	}
}

func routesObjects(cmd *commands) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/objects/list",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(handler.ObjectListAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/objects/create",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.ObjectsCreateAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/objects/update",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(handler.ObjectUpdateAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/objects/remove",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(handler.ObjectsRemoveAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
	}
}

func routesLink(cmd *commands) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/links",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(handler.LinksAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/links/create",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.LinksCreateAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
							"aes":    cmd.aes,
							"cfg":    cmd.config,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/links/update",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(handler.LinkUpdateAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
							"aes":    cmd.aes,
							"cfg":    cmd.config,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/links/remove",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(handler.LinksRemoveAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/links/blast",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.LinksBlastAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
							"aes":    cmd.aes,
							"cfg":    cmd.config,
						},
					),
				),
			),
			AdminOnly: true,
		},
	}
}

func routesGroups(cmd *commands) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/groups",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(handler.GroupsAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/groups/create",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.GroupsCreateAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/groups/update",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(handler.GroupUpdateAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/groups/remove",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(handler.GroupsRemoveAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
	}
}

func routesQuestions(cmd *commands) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/api/questions",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(handler.QuestionsAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/questions/create",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(handler.QuestionCreateAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/questions/update",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(handler.QuestionUpdateAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
		{
			Pattern: "/api/questions/remove",
			Handler: middleware.WithCors(
				cmd.config.GetConfig().App.Cors.AllowedMethods,
				cmd.config.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(handler.QuestionRemoveAPI),
						map[string]interface{}{
							"logger": cmd.logger,
							"view":   cmd.view,
							"sm":     cmd.session,
							"hash":   cmd.hash,
							"db":     cmd.db,
						},
					),
				),
			),
			AdminOnly: true,
		},
	}
}
