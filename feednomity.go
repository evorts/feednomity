package main

import (
	"context"
	"fmt"
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/template"
	"net/http"
	"strconv"
	"time"
)

type commands struct {
	db      database.IManager
	acl     acl.IManager
	logger  logger.IManager
	config  config.IManager
	session session.IManager
	aes     crypt.ICryptAES
	hash    crypt.ICryptHash
	view    template.IManager
}

func main() {
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
		logging.Fatal("error initialize access control")
		return
	}
	sm := session.NewSession(
		time.Duration(cfg.GetConfig().App.SessionExpiration),
		time.Duration(30),
		session.Cookie{
			Name:     "feednomid",
			Domain:   cfg.GetConfig().App.CookieDomain,
			HttpOnly: false,
			Path:     "/",
			Persist:  false,
			SameSite: 0,
			Secure:   false,
		},
	)
	tm, _ := template.NewTemplates(cfg.GetConfig().App.TemplateDirectory, map[string]interface{}{
		"CopyrightYear": strconv.Itoa(time.Now().Year()),
		"FavIcon":       cfg.GetConfig().App.Logo.FavIcon,
		"LogoUrl":       cfg.GetConfig().App.Logo.Url,
		"LogoAlt":       cfg.GetConfig().App.Logo.Alt,
	}).LoadTemplates()
	o := http.NewServeMux()
	routes(o, &commands{
		ds, accessControl,logging, cfg, sm,
		aesCryptic, crypt.NewHashEncryption(cfg.GetConfig().App.HashSalt), tm,
	})
	logging.Log("started", "Dashboard app started.")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetConfig().App.Port), sm.LoadAndSave(o)); err != nil {
		logging.Fatal(err)
	}
}
