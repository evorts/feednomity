package cmd

import (
	"fmt"
	"github.com/evorts/feednomity/handler"
	"github.com/evorts/feednomity/handler/hapi"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/jwe"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/mailer"
	"github.com/evorts/feednomity/pkg/memory"
	"github.com/evorts/feednomity/pkg/middleware"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

const apiPrefix = "/rest"
const apiVersionPrefix = "/v1"

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
	mem memory.IManager,
	mail mailer.IMailer,
) {
	routes := make([]reqio.Route, 0)
	prefix := fmt.Sprintf("%s%s", apiPrefix, apiVersionPrefix)
	routes = append(routes, routesApiMaintenance(apiPrefix, cfg, view, db, accessControl, jwx, log)...)
	routes = append(routes, routesApiFeedbacks(prefix, cfg, view, db, accessControl, jwx, aes, log)...)
	routes = append(routes, routesApiDistribution(prefix, cfg, view, db, accessControl, jwx, aes, log)...)
	routes = append(routes, routesApiUsers(prefix, cfg, view, db, accessControl, jwx, hash, log, mail, mem)...)
	routes = append(routes, routesApiLink(prefix, cfg, view, db, accessControl, jwx, hash, aes, log)...)
	routes = append(routes, routesApiGroups(prefix, cfg, view, db, accessControl, jwx, hash, log)...)
	routes = append(routes, routesApiOrganizations(prefix, cfg, view, db, accessControl, jwx, hash, log)...)
	routes = append(routes, routesApiQuestions(prefix, cfg, view, db, accessControl, jwx, hash, log)...)
	reqio.NewRoutes(routes).ExecRoutes(o)
}

func routesApiFeedbacks(
	pathPrefix string,
	cfg config.IManager, view view.IManager, db database.IManager, accessControl acl.IManager,
	jwx jwe.IManager, aes crypt.ICryptAES, log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: fmt.Sprintf("%s/reviews/list", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiReviewList),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: fmt.Sprintf("%s/reviews/detail/", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodGet,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiReviewDetail),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
						"aes":    aes,
					},
				),
			),
		},
		{
			Pattern: fmt.Sprintf("%s/reviews/submit", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiReviewSubmit),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
						"aes":    aes,
					},
				),
			),
		},
	}
}

func routesApiDistribution(
	pathPrefix string,
	cfg config.IManager, view view.IManager, db database.IManager, accessControl acl.IManager,
	jwx jwe.IManager, aes crypt.ICryptAES, log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: fmt.Sprintf("%s/distributions/publish", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/distributions/list", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/distributions/create", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/distributions/update", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/distributions/delete", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/dist-objects/list", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/dist-objects/create", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/dist-objects/update", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/dist-objects/delete", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
	pathPrefix string,
	cfg config.IManager, view view.IManager, db database.IManager, accessControl acl.IManager,
	jwx jwe.IManager, hash crypt.ICryptHash, log logger.IManager, mail mailer.IMailer, mem memory.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: fmt.Sprintf("%s/users/forgot-password", pathPrefix),
			Handler: middleware.WithFiltersForApi(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				view,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiForgotPassword),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"hash":   hash,
						"cfg":    cfg,
						"mail":   mail,
						"mem":    mem,
						"db":     db,
					},
				),
			),
		},
		{
			Pattern: fmt.Sprintf("%s/users/change-password", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiChangePassword),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
						"hash":   hash,
					},
				),
			),
		},
		{
			Pattern: fmt.Sprintf("%s/users/create-password", pathPrefix),
			Handler: middleware.WithFiltersForApi(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				view,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiCreatePassword),
					map[string]interface{}{
						"logger": log,
						"view":   view,
						"db":     db,
						"mem":    mem,
					},
				),
			),
		},
		{
			Pattern: fmt.Sprintf("%s/users/login", pathPrefix),
			Handler: middleware.WithFiltersForApi(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				view,
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
		},
		{
			Pattern: fmt.Sprintf("%s/users/list", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/users/create", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/users/update", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/users/delete", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
	pathPrefix string,
	cfg config.IManager, view view.IManager, db database.IManager, accessControl acl.IManager,
	jwx jwe.IManager, hash crypt.ICryptHash, aes crypt.ICryptAES, log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: fmt.Sprintf("%s/links/list", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/links/create", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/links/update", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/links/delete", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
	pathPrefix string,
	cfg config.IManager, view view.IManager, db database.IManager, accessControl acl.IManager,
	jwx jwe.IManager, hash crypt.ICryptHash, log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: fmt.Sprintf("%s/groups/list", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/groups/create", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/groups/update", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/groups/delete", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
	pathPrefix string,
	cfg config.IManager, view view.IManager, db database.IManager, accessControl acl.IManager,
	jwx jwe.IManager, hash crypt.ICryptHash, log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: fmt.Sprintf("%s/organizations/list", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/organizations/create", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/organizations/update", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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
			Pattern: fmt.Sprintf("%s/organizations/delete", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
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

func routesApiMaintenance(
	pathPrefix string,
	cfg config.IManager, view view.IManager, db database.IManager, accessControl acl.IManager,
	jwx jwe.IManager, log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: fmt.Sprintf("%s/ping", pathPrefix),
			Handler: middleware.WithInjection(
				http.HandlerFunc(handler.Ping),
				map[string]interface{}{
					"view": view,
				},
			),
		},
		{
			Pattern: fmt.Sprintf("%s/reload", pathPrefix),
			Handler: middleware.WithTokenProtection(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				accessControl, jwx, view,
				middleware.WithInjection(
					http.HandlerFunc(hapi.ApiReload),
					map[string]interface{}{
						"view":   view,
						"logger": log,
					},
				),
			),
		},
	}
}

func routesApiQuestions(
	pathPrefix string,
	cfg config.IManager, view view.IManager, db database.IManager, accessControl acl.IManager,
	jwx jwe.IManager, hash crypt.ICryptHash, log logger.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: fmt.Sprintf("%s/questions", pathPrefix),
			Handler: middleware.WithFiltersForApi(
				http.MethodGet,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				view,
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
		},
		{
			Pattern: fmt.Sprintf("%s/questions/create", pathPrefix),
			Handler: middleware.WithFiltersForApi(
				http.MethodPost,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				view,
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
		},
		{
			Pattern: fmt.Sprintf("%s/questions/update", pathPrefix),
			Handler: middleware.WithFiltersForApi(
				http.MethodPut,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				view,
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
		},
		{
			Pattern: fmt.Sprintf("%s/questions/remove", pathPrefix),
			Handler: middleware.WithFiltersForApi(
				http.MethodDelete,
				cfg.GetConfig().App.Cors.AllowedMethods,
				cfg.GetConfig().App.Cors.AllowedOrigins,
				view,
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
		},
	}
}
