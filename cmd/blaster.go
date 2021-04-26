package cmd

import (
	"context"
	"fmt"
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/pkg/cli"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/mailer"
	"github.com/robfig/cron/v3"
	"sync"
)

var (
	l logger.IManager
	cfg config.IManager
	ds database.IManager
	initLib sync.Once
)
var Blaster = &cli.Command{
	Description: "Mail blaster command line, cron job style",
	Run: func(cmd *cli.Command, args []string) {
		fmt.Println("mail blaster command running...")
		initLib.Do(instantiateLib)
		c := cron.New()
		_, err2 := c.AddFunc(cfg.GetConfig().CronJobs.Blaster.Schedule, runCronBlaster)
		if err2 != nil {
			c.Stop()
			return
		}
		c.Run()
	},
}

func instantiateLib()  {
	var err error
	l = logger.NewLogger()
	cfg, err = config.NewConfig("config.main.yml", "config.yml").Initiate()
	if err != nil {
		l.Fatal("error reading configuration")
		return
	}
	ds = database.NewDB(
		cfg.GetConfig().DB.Dsn,
		cfg.GetConfig().DB.MaxConnectionLifetime,
		cfg.GetConfig().DB.MaxIdleConnection,
		cfg.GetConfig().DB.MaxOpenConnection,
		false,
	)
	ds.MustConnect(context.Background())
	_ = mailer.NewSendInBlue(
		cfg.GetConfig().Mailer.Providers.Get("send_in_blue").ApiUrl,
		cfg.GetConfig().Mailer.Providers.Get("send_in_blue").ApiKey,
	)
}

func runCronBlaster()  {
	initLib.Do(instantiateLib)
	_ = distribution.NewDistributionDomain(ds)

	fmt.Println("cron email blaster")
}