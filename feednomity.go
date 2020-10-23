package main

import (
	"context"
	"fmt"
	"github.com/evorts/feednomity/handler"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/db"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/middleware"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/template"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type commands struct {
	db      db.IManager
	logger  logger.IManager
	config  config.IManager
	session session.IManager
	crypt   crypt.ICrypt
	hash    crypt.ICrypt
	view    template.IManager
}

func routes(o *http.ServeMux, cmd *commands) {
	// serving assets
	fs := http.FileServer(http.Dir(cmd.config.GetConfig().App.AssetDirectory))
	o.Handle("/assets/", http.StripPrefix("/assets", fs))
	// serving pages
	reqio.NewRoutes([]reqio.Route{
		{
			Pattern: "/",
			Handler: middleware.WithInjection(
				middleware.WithProtection(http.HandlerFunc(handler.Dashboard)),
				map[string]interface{}{
					"logger": cmd.logger,
					"view":   cmd.view,
				},
			),
			MemberOnly: true,
		},
		{
			Pattern: "/forms",
			Handler: middleware.WithInjection(
				http.HandlerFunc(handler.Forms),
				map[string]interface{}{
					"logger": cmd.logger,
					"view":   cmd.view,
				},
			),
			MemberOnly: false,
		},
		{
			Pattern: "/login",
			Handler: middleware.WithInjection(
				http.HandlerFunc(handler.Login),
				map[string]interface{}{
					"logger": cmd.logger,
					"view":   cmd.view,
					"sm":     cmd.session,
					"hash":   cmd.hash,
				},
			),
			MemberOnly: false,
		},
		{
			Pattern: "/logout",
			Handler: middleware.WithInjection(
				http.HandlerFunc(handler.Logout),
				map[string]interface{}{
					"sm": cmd.session,
				},
			),
			MemberOnly: true,
		},
		{
			Pattern: "/reload",
			Handler: middleware.WithInjection(
				http.HandlerFunc(handler.Reload),
				map[string]interface{}{
					"logger": cmd.logger,
					"view":   cmd.view,
					"config": cmd.config,
				},
			),
			MemberOnly: true,
		},
		{
			Pattern: "/ping",
			Handler: middleware.WithInjection(
				http.HandlerFunc(handler.Ping),
				map[string]interface{}{
					"view": cmd.view,
				},
			),
			MemberOnly: false,
		},
	}).ExecRoutes(o)
}

func main() {
	logging := logger.NewLogger()
	cfg, err := config.NewConfig("config.main.yml", "config.yml").Initiate()
	if err != nil {
		logging.Fatal("error reading configuration")
		return
	}
	ds := db.NewDB(
		cfg.GetConfig().DB.Dsn,
		cfg.GetConfig().DB.MaxConnectionLifetime,
		cfg.GetConfig().DB.MaxIdleConnection,
		cfg.GetConfig().DB.MaxOpenConnection,
	)

	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)
	go func() {
		select {
		case <-appSignal:
			stop()
		}
	}()

	ds.MustConnect(ctx)
	//goland:noinspection GoUnhandledErrorResult
	defer ds.Close(ctx)
	sm := session.NewSession(
		cfg.GetConfig().App.SessionExpiration,
		time.Duration(30),
		session.Cookie{
			Name:     "feednonimid",
			Domain:   cfg.GetConfig().App.CookieDomain,
			HttpOnly: false,
			Path:     "/",
			Persist:  false,
			SameSite: 0,
			Secure:   false,
		},
	)
	tm, _ := template.NewTemplates(cfg.GetConfig().App.TemplateDirectory, map[string]interface{}{
		"PageAttributes": map[string]interface{}{
			"Year":    strconv.Itoa(time.Now().Year()),
			"FavIcon": cfg.GetConfig().App.Logo.FavIcon,
			"LogoUrl": cfg.GetConfig().App.Logo.Url,
			"LogoAlt": cfg.GetConfig().App.Logo.Alt,
		},
	}).LoadTemplates()
	o := http.NewServeMux()
	routes(o, &commands{
		ds, logging, cfg, sm, crypt.NewCrypt(cfg.GetConfig().App.Salt),
		crypt.NewCrypt(""), tm,
	})
	logging.Log("started", "Dashboard app started.")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetConfig().App.Port), sm.LoadAndSave(o)); err != nil {
		logging.Fatal(err)
	}
}
