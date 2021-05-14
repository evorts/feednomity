package cmd

import (
	"github.com/evorts/feednomity/handler/hapi"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/jwe"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/middleware"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func routingApi(
	o *http.ServeMux,
	cfg config.IManager,
	view view.IManager,
	db database.IManager,
	accessControl acl.IManager,
	jwx jwe.IManager,
	hash crypt.ICryptHash,
	aes crypt.ICryptAES,
	log logger.IManager,
) {
	routes := []reqio.Route{
		{
			Pattern: "/ping",
			Handler: middleware.WithMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(hapi.Ping),
					map[string]interface{}{
						"view": view,
					},
				),
			),
		},
		{
			Pattern: "/reload",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiReload),
					map[string]interface{}{
						"view":   view,
						"logger": log,
					},
				),
			),
		},
		{
			Pattern: "/feedbacks",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiFeedbackSubmission),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"hash":   hash,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/360/submission",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.Api360Submission),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
						"hash":   hash,
					},
				),
			),
		},
	}
	routes = append(routes, routesApiDistribution(cfg, view, db, accessControl, jwx, aes, log)...)
	routes = append(routes, routesApiUsers(cfg, view, db, accessControl, jwx, hash, log)...)
	routes = append(routes, routesApiLink(cfg, view, db, accessControl, jwx, hash, aes, log)...)
	routes = append(routes, routesApiGroups(cfg, view, db, accessControl, jwx, hash, log)...)
	routes = append(routes, routesApiOrganizations(cfg, view, db, accessControl, jwx, hash, log)...)
	routes = append(routes, routesApiQuestions(cfg, view, db, accessControl, jwx, hash, log)...)

	reqio.NewRoutes(routes).ExecRoutes(o)
}

func routesApiDistribution(
	cfg config.IManager,
	view view.IManager,
	db database.IManager,
	accessControl acl.IManager,
	jwx jwe.IManager,
	aes crypt.ICryptAES,
	log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/distributions/publish",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiDistributionBlast),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
						"cfg":    cfg,
					},
				),
			),
		},
		{
			Pattern: "/distributions/list",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiDistributionsList),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/distributions/create",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiDistributionsCreate),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/distributions/update",
			Handler: middleware.WithTokenProtection(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiDistributionsUpdate),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/distributions/delete",
			Handler: middleware.WithTokenProtection(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiDistributionsDelete),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/dist-objects/list",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiDistObjectsList),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/dist-objects/create",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiDistObjectsCreate),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
						"aes":    aes,
						"cfg":    cfg,
					},
				),
			),
		},
		{
			Pattern: "/dist-objects/update",
			Handler: middleware.WithTokenProtection(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiDistObjectsUpdate),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/dist-objects/delete",
			Handler: middleware.WithTokenProtection(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiDistObjectsDelete),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
	}
}

func routesApiUsers(
	cfg config.IManager,
	view view.IManager,
	db database.IManager,
	accessControl acl.IManager,
	jwx jwe.IManager,
	hash crypt.ICryptHash,
	log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/users/login",
			Handler: middleware.WithCorsProtection(
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(hapi.ApiLogin),
						map[string]interface{}{
							"logger": log,
							"view":   view,
							"hash":   hash,
							"db":     db,
							"jwx":    jwx,
						},
					),
				),
			),
		},
		{
			Pattern: "/users/list",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiUsersList),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/users/create",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiUserCreate),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/users/update",
			Handler: middleware.WithTokenProtection(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiUserUpdate),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/users/delete",
			Handler: middleware.WithTokenProtection(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiUsersDelete),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
	}
}

func routesApiLink(
	cfg config.IManager,
	view view.IManager,
	db database.IManager,
	accessControl acl.IManager,
	jwx jwe.IManager,
	hash crypt.ICryptHash,
	aes crypt.ICryptAES,
	log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/links/list",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiLinksList),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/links/create",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiLinksCreate),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
						"aes":    aes,
						"cfg":    cfg,
					},
				),
			),
		},
		{
			Pattern: "/links/update",
			Handler: middleware.WithTokenProtection(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiLinkUpdate),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"hash":   hash,
						"db":     db,
						"aes":    aes,
						"cfg":    cfg,
					},
				),
			),
		},
		{
			Pattern: "/links/delete",
			Handler: middleware.WithTokenProtection(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiLinksDelete),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
	}
}

func routesApiGroups(
	cfg config.IManager,
	view view.IManager,
	db database.IManager,
	accessControl acl.IManager,
	jwx jwe.IManager,
	hash crypt.ICryptHash,
	log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/groups/list",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiGroupsList),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"hash":   hash,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/groups/create",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiGroupsCreate),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"hash":   hash,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/groups/update",
			Handler: middleware.WithTokenProtection(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiGroupUpdate),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"hash":   hash,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/groups/delete",
			Handler: middleware.WithTokenProtection(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiGroupsDelete),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"hash":   hash,
						"db":     db,
					},
				),
			),
		},
	}
}

func routesApiOrganizations(
	cfg config.IManager,
	view view.IManager,
	db database.IManager,
	accessControl acl.IManager,
	jwx jwe.IManager,
	hash crypt.ICryptHash,
	log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/organizations/list",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiOrganizationsList),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"hash":   hash,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/organizations/create",
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiOrganizationsCreate),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"hash":   hash,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/organizations/update",
			Handler: middleware.WithTokenProtection(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiOrganizationUpdate),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"hash":   hash,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: "/organizations/delete",
			Handler: middleware.WithTokenProtection(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiOrganizationsDelete),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"hash":   hash,
						"db":     db,
					},
				),
			),
		},
	}
}

func routesApiQuestions(
	cfg config.IManager,
	view view.IManager,
	db database.IManager,
	accessControl acl.IManager,
	jwx jwe.IManager,
	hash crypt.ICryptHash,
	log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/questions",
			Handler: middleware.WithCorsProtection(
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodGet,
					middleware.WithInjection(
						http.HandlerFunc(hapi.ApiQuestions),
						map[string]interface{}{
							"logger": log,
							"view":   view,
							"hash":   hash,
							"db":     db,
						},
					),
				),
			),
		},
		{
			Pattern: "/questions/create",
			Handler: middleware.WithCorsProtection(
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPost,
					middleware.WithInjection(
						http.HandlerFunc(hapi.ApiQuestionCreate),
						map[string]interface{}{
							"logger": log,
							"view":   view,
							"hash":   hash,
							"db":     db,
						},
					),
				),
			),
		},
		{
			Pattern: "/questions/update",
			Handler: middleware.WithCorsProtection(
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodPut,
					middleware.WithInjection(
						http.HandlerFunc(hapi.ApiQuestionUpdate),
						map[string]interface{}{
							"logger": log,
							"view":   view,
							"hash":   hash,
							"db":     db,
						},
					),
				),
			),
		},
		{
			Pattern: "/questions/remove",
			Handler: middleware.WithCorsProtection(
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				middleware.WithMethodFilter(
					http.MethodDelete,
					middleware.WithInjection(
						http.HandlerFunc(hapi.ApiQuestionRemove),
						map[string]interface{}{
							"logger": log,
							"view":   view,
							"hash":   hash,
							"db":     db,
						},
					),
				),
			),
		},
	}
}
