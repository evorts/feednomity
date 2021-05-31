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
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"sync"
	"time"
)

var (
	l          logger.IManager
	cfg        config.IManager
	ds         database.IManager
	initLib    sync.Once
	m          mailer.IMailer
	today      *time.Time
	todayLimit int
)
var Blaster = &cli.Command{
	Description: "Mail blaster command line, cron job style",
	Run: func(cmd *cli.Command, args []string) {
		fmt.Println("mail blaster command running...")
		initLib.Do(instantiateLib)
		c := cron.New()
		_, err2 := c.AddFunc(cfg.GetConfig().CronJobs.Blaster.Schedule, runCronBlaster)
		if err2 != nil {
			fmt.Println(err2)
			c.Stop()
			return
		}
		c.Run()
	},
}

func instantiateLib() {
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
	m = mailer.NewSendInBlue(
		cfg.GetConfig().Mailer.Providers.Get("send_in_blue").Get("api_url"),
		cfg.GetConfig().Mailer.Providers.Get("send_in_blue").Get("api_key"),
	)
	m.SetSender(
		cfg.GetConfig().Mailer.SenderName,
		cfg.GetConfig().Mailer.SenderEmail,
	)
}

func runCronBlaster() {
	l.Log("cron_phase_start", "cron for email blasting started...")
	initLib.Do(instantiateLib)
	now := time.Now()
	if today == nil || today.Day() != now.Day() {
		today = &now
		todayLimit = cfg.GetConfig().Mailer.DailyLimit
	}
	if todayLimit-cfg.GetConfig().CronJobs.Blaster.BatchRows < 1 {
		l.Log(
			"cron_mail_blast_limit",
			fmt.Sprintf("reaching mail blast limit for today: %s", today.Format("2006-01-02")),
		)
		return
	}
	distDomain := distribution.NewDistributionDomain(ds)
	items, _, err := distDomain.FindAllQueues(context.Background(), 1, cfg.GetConfig().CronJobs.Blaster.BatchRows)
	if err != nil {
		l.Log("cron_queue_find_all_error", err)
		return
	}
	if len(items) < 1 {
		l.Log("cron_queue_find_all_empty", "no items in queue found")
		return
	}
	objectIds := make([]int64, 0)
	templates := make(map[string]string, 0)
	for _, item := range items {
		time.Sleep(10 * time.Millisecond)
		respondentName := ""
		if v, ok := item.Arguments["respondent_name"]; ok {
			respondentName = v.(string)
		}
		_, ok := templates[item.Template]
		if !ok {
			content, err2 := ioutil.ReadFile(fmt.Sprintf("%s/%s", cfg.GetConfig().App.MailTemplateDirectory, item.Template))
			if err2 != nil {
				continue
			}
			templates[item.Template] = string(content)
		}
		html := templates[item.Template]
		m.SetReplyTo("No Reply", "no-reply@evorts.com")
		var body []byte
		body, err = m.SendHtml(
			context.Background(),
			[]mailer.Target{
				{
					Name:  respondentName,
					Email: item.ToEmail,
				},
			},
			item.Subject, html,
			utils.MapStringInterface(item.Arguments).ToMapString(),
		)
		if err != nil {
			l.Log(
				"cron_mail_send_error",
				fmt.Sprintf(
					"sending mail to: %s, with subject: %s, error: %v",
					item.ToEmail,
					item.Subject, err,
				),
			)
			continue
		}
		fmt.Println(string(body))
		l.Log(
			"cron_mail_send_success",
			fmt.Sprintf(
				"sending mail to: %s, with subject: %s, id: %d",
				item.ToEmail,
				item.Subject, item.Id,
			),
		)
		objectIds = append(objectIds, item.Id)
	}
	if len(objectIds) > 0 {
		err = distDomain.UpdateObjectRetryCountByIds(context.Background(), objectIds...)
		if err != nil {
			l.Log("cron_queue_update_retry_error", err)
		}
		l.Log("cron_queue_deletion_in_progress", objectIds)
		err = distDomain.DeleteQueueByIds(context.Background(), objectIds...)
		if err != nil {
			l.Log("cron_queue_deletion_error", err)
		}
	}
	todayLimit -= len(objectIds)
	l.Log("cron_mail_daily_limit_left", fmt.Sprintf("daily limit left: %d", todayLimit))
}
