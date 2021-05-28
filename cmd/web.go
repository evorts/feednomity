package cmd

import (
	"context"
	"fmt"
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/cli"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/jwe"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/memory"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"strconv"
	"time"
)

var Web = &cli.Command{
	Description: "Web Application",
	Run: func(cmd *cli.Command, args []string) {
		logging := logger.NewLogger()
		cfg, err := config.NewConfig("config.main.yml", "config.yml").Initiate()
		if err != nil {
			logging.Fatal("error reading configuration")
			return
		}
		aesCryptic := crypt.NewCryptAES(cfg.GetConfig().App.AESSalt)
		if _, err = aesCryptic.Initialize(); err != nil {
			logging.Fatal("error initialize cryptic modules")
			return
		}
		mem := memory.NewRedisStorage(
			cfg.GetConfig().Memory.Get("redis").Address,
			cfg.GetConfig().Memory.Get("redis").Password,
			cfg.GetConfig().Memory.Get("redis").Db,
		)
		mem.MustConnect(context.Background())
		ds := database.NewDB(
			cfg.GetConfig().DB.Dsn,
			cfg.GetConfig().DB.MaxConnectionLifetime,
			cfg.GetConfig().DB.MaxIdleConnection,
			cfg.GetConfig().DB.MaxOpenConnection,
			false,
		)
		ds.MustConnect(context.Background())
		defer func() {
			//_ = ds.Close(context.Background())
		}()
		accessControl := acl.NewACLManager(users.NewUserDomain(ds), users.NewUserAccessDomain(ds))
		if err2 := accessControl.Populate(); err2 != nil {
			logging.Log("fatal_error", "error initialize access control")
			logging.Fatal(err2.Error())
			return
		}
		sm := session.NewSession(
			time.Duration(cfg.GetConfig().App.SessionExpiration),
			time.Duration(30),
			session.Cookie{
				Name:     "feednomid",
				Domain:   cfg.GetConfig().App.CookieDomain,
				HttpOnly: true,
				Path:     "/",
				Persist:  true,
				SameSite: http.SameSiteStrictMode,
				Secure:   cfg.GetConfig().App.CookieSecure == 1,
			},
		)
		key := jwe.Key{Value: cfg.GetConfig().Jwe.Key}
		pk, errJwe := key.GetPrivate()
		if pk == nil || errJwe != nil {
			logging.Log("fatal_error", "error initialize jwe")
			logging.Fatal(errJwe)
			return
		}
		jwx := jwe.NewJWE(pk, cfg.GetConfig().Jwe.Expire)
		tm, _ := view.NewTemplateManager(cfg.GetConfig().App.TemplateDirectory, map[string]interface{}{
			"CopyrightYear": strconv.Itoa(time.Now().Year()),
			"FavIcon":       cfg.GetConfig().App.Logo.FavIcon,
			"LogoUrl":       cfg.GetConfig().App.Logo.Url,
			"LogoAlt":       cfg.GetConfig().App.Logo.Alt,
			"ApiBaseUrl":    cfg.GetConfig().App.BaseUrlApi,
		}).LoadTemplates()
		o := http.NewServeMux()
		routingWeb(
			o, accessControl, logging, cfg, sm,
			aesCryptic, crypt.NewHashEncryption(cfg.GetConfig().App.HashSalt),
			tm, jwx, mem, ds,
		)
		logging.Log("started", "Web Application Started.")
		if err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetConfig().App.Port), sm.LoadAndSave(o)); err != nil {
			logging.Fatal(err)
		}
	},
}
