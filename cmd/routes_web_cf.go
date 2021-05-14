package cmd

import (
	"github.com/evorts/feednomity/handler/hcf"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/crypt"
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
) []reqio.Route {
	return []reqio.Route{
		{
			Pattern: "/mbr/login",
			Handler: middleware.WithMethodFilter(
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
			Pattern: "/mbr/link/",
			Handler: middleware.WithMethodFilter(
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
				session, view, accessControl,
				middleware.WithInjection(
					http.HandlerFunc(hcf.ReviewListing),
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
			Pattern: "/mbr/review/form/",
			Handler: middleware.WithSessionProtection(
				session, view, accessControl,
				middleware.WithInjection(
					http.HandlerFunc(hcf.Review360Form),
					map[string]interface{}{
						"logger": logger,
						"view":   view,
						"sm":     session,
						"hash":   hash,
					},
				),
			),
		},
	}
}
