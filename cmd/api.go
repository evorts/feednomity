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
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

var Api = &cli.Command{
	Description: "API Application",
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
		key := jwe.Key{Value: cfg.GetConfig().Jwe.Key}
		pk, errJwe := key.GetPrivate()
		if pk == nil || errJwe != nil {
			logging.Log("fatal_error", "error initialize jwe")
			logging.Fatal(errJwe)
			return
		}
		jwx := jwe.NewJWE(pk, cfg.GetConfig().Jwe.Expire)
		o := http.NewServeMux()
		routingApi(
			o, cfg, view.NewJsonManager(), ds, accessControl, jwx,
			crypt.NewHashEncryption(cfg.GetConfig().App.HashSalt), aesCryptic, logging,
		)
		logging.Log("started", "API Started.")
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetConfig().App.PortApi), o); err != nil {
			logging.Fatal(err)
		}
	},
}
