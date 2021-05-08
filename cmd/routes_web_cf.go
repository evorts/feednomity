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
			Pattern: "/page/360/review",
			Handler: middleware.WithMethodFilter(
				http.MethodGet,
				middleware.WithInjection(
					http.HandlerFunc(hcf.Form360),
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
