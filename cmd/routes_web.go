package cmd

import (
	"github.com/evorts/feednomity/handler"
	"github.com/evorts/feednomity/handler/hcf"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/jwe"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/memory"
	"github.com/evorts/feednomity/pkg/middleware"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func routingWeb(
	o *http.ServeMux,
	acl acl.IManager,
	logger logger.IManager,
	cfg config.IManager,
	session session.IManager,
	aes crypt.ICryptAES,
	hash crypt.ICryptHash,
	view view.ITemplateManager,
	jwx jwe.IManager,
	mem memory.IManager,
	ds database.IManager,
) {
	// serving assets
	fs := http.FileServer(http.Dir(cfg.GetConfig().App.AssetDirectory))
	o.Handle("/assets/", http.StripPrefix("/assets", fs))
	// serving pages
	routes := []reqio.Route{
		{
			Pattern: "/ping",
			Handler: middleware.WithWebMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(handler.Ping),
					map[string]interface{}{
						"view": view,
					},
				),
			),
		},
		{
			Pattern: "/",
			Handler: middleware.WithInjection(
				http.HandlerFunc(hcf.Home),
				map[string]interface{}{
					"view": view,
					"sm":   session,
				},
			),
		},
		{
			Pattern: "/forgot-password",
			Handler: middleware.WithWebMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(hcf.ForgotPassword),
					map[string]interface{}{
						"logger": logger,
						"sm":     session,
						"vm":     view,
						"hash":   hash,
					},
				),
			),
		},
		{
			Pattern: "/crp/",
			Handler: middleware.WithWebMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(hcf.CreatePassword),
					map[string]interface{}{
						"logger": logger,
						"sm":     session,
						"vm":     view,
					},
				),
			),
		},
		{
			Pattern: "/chp/",
			Handler: middleware.WithSessionProtection(
				session, view, acl, jwx, cfg,
				middleware.WithInjection(
					http.HandlerFunc(hcf.ChangePassword),
					map[string]interface{}{
						"logger": logger,
						"view":   view,
						"sm":     session,
						"hash":   hash,
						"db":     ds,
					},
				),
			),
		},
		{
			Pattern: "/logout",
			Handler: middleware.WithWebMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(hcf.Logout),
					map[string]interface{}{
						"logger": logger,
						"cfg":    cfg,
						"sm":     session,
					},
				),
			),
		},
	}
	routes = append(routes, routesWebDashboard(acl, logger, session, hash, view, jwx, cfg)...)
	routes = append(routes, routesWebConsumers(acl, logger, session, hash, view, jwx, ds, cfg, mem)...)
	reqio.NewRoutes(routes).ExecRoutes(o)
}
