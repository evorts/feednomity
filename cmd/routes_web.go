package cmd

import (
	"github.com/evorts/feednomity/handler"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/jwe"
	"github.com/evorts/feednomity/pkg/logger"
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
	config config.IManager,
	session session.IManager,
	aes crypt.ICryptAES,
	hash crypt.ICryptHash,
	view view.ITemplateManager,
	jwe jwe.IManager,
) {
	// serving assets
	fs := http.FileServer(http.Dir(cfg.GetConfig().App.AssetDirectory))
	o.Handle("/assets/", http.StripPrefix("/assets", fs))
	// serving pages
	routes := []reqio.Route{
		{
			Pattern: "/ping",
			Handler: middleware.WithMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(handler.Ping),
					map[string]interface{}{
						"view": view,
					},
				),
			),
		},
	}
	routes = append(routes, routesWebDashboard(acl, logger, session, hash, view)...)
	routes = append(routes, routesWebConsumers(acl, logger, session, hash, view)...)
	reqio.NewRoutes(routes).ExecRoutes(o)
}
