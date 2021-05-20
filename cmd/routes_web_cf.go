package cmd

import (
	"github.com/evorts/feednomity/handler/hcf"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/jwe"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/middleware"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func routesWebConsumers(
	accessControl acl.IManager,
	logger logger.IManager,
	session session.IManager,
	hash crypt.ICryptHash,
	view view.ITemplateManager,
	jwx jwe.IManager,
	ds database.IManager,
	cfg config.IManager,
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/mbr/login",
			Handler: middleware.WithWebMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(hcf.Login),
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
			Pattern: "/mbr/logout",
			Handler: middleware.WithWebMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(hcf.Logout),
					map[string]interface{}{
						"logger": logger,
						"sm":     session,
					},
				),
			),
		},
		{
			Pattern: "/mbr/link/",
			Handler: middleware.WithWebMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(hcf.Link),
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
			Pattern: "/mbr/review/list",
			Handler: middleware.WithSessionProtection(
				session, view, accessControl, jwx, cfg,
				middleware.WithInjection(
					http.HandlerFunc(hcf.ReviewListing),
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
			Pattern: "/mbr/reviews/",
			Handler: middleware.WithSessionProtection(
				session, view, accessControl, jwx, cfg,
				middleware.WithInjection(
					http.HandlerFunc(hcf.ReviewDetail),
					map[string]interface{}{
						"logger": logger,
						"view":   view,
						"sm":     session,
						"hash":   hash,
						"db":     ds,
						"cfg":    cfg,
					},
				),
			),
		},
	}
}
