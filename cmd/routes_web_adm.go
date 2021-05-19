package cmd

import (
	"github.com/evorts/feednomity/handler/hadm"
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

func routesWebDashboard(
	accessControl acl.IManager,
	logger logger.IManager,
	session session.IManager,
	hash crypt.ICryptHash,
	view view.ITemplateManager,
	jwx jwe.IManager,
	cfg config.IManager,

) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/adm",
			Handler: middleware.WithWebMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(hadm.AdminGate),
					map[string]interface{}{
						"logger": logger,
						"sm":     session,
					},
				),
			),
		},
		{
			Pattern: "/adm/dashboard",
			Handler: middleware.WithSessionProtection(
				session, view, accessControl, jwx, cfg,
				middleware.WithInjection(
					http.HandlerFunc(hadm.Dashboard),
					map[string]interface{}{
						"logger": logger,
						"view":   view,
						"sm":     session,
					},
				),
			),
		},
		{
			Pattern: "/adm/users",
			Handler: middleware.WithSessionProtection(
				session, view, accessControl, jwx, cfg,
				middleware.WithInjection(
					http.HandlerFunc(hadm.Users),
					map[string]interface{}{
						"logger": logger,
						"view":   view,
						"sm":     session,
					},
				),
			),
		},
		{
			Pattern: "/adm/objects",
			Handler: middleware.WithSessionProtection(
				session, view, accessControl, jwx, cfg,
				middleware.WithInjection(
					http.HandlerFunc(hadm.Objects),
					map[string]interface{}{
						"logger": logger,
						"view":   view,
						"sm":     session,
					},
				),
			),
		},
		{
			Pattern: "/adm/login",
			Handler: middleware.WithWebMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(hadm.Login),
					map[string]interface{}{
						"logger": logger,
						"view":   view,
						"sm":     session,
						"hash":   hash,
					},
				),
			),
		},
		{
			Pattern: "/adm/logout",
			Handler: middleware.WithWebMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(hadm.Logout),
					map[string]interface{}{
						"sm": session,
					},
				),
			),
		},
	}
}
