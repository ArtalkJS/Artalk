package notify_launcher

import (
	"context"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib/email"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/dingding"
	"github.com/nikoksr/notify/service/line"
	"github.com/nikoksr/notify/service/slack"
	"github.com/nikoksr/notify/service/telegram"
)

var NotifyCtx = context.Background()

func Init() {
	// 初始化邮件队列
	email.InitQueue()

	// Telegram
	tgConf := config.Instance.AdminNotify.Telegram
	if tgConf.Enabled {
		telegramService, _ := telegram.New(tgConf.ApiToken)
		telegramService.AddReceivers(tgConf.Receivers...)
		notify.UseServices(telegramService)
	}

	// 钉钉
	dingTalkConf := config.Instance.AdminNotify.DingTalk
	if dingTalkConf.Enabled {
		dingTalkService := dingding.New(&dingding.Config{Token: dingTalkConf.Token, Secret: dingTalkConf.Secret})
		notify.UseServices(dingTalkService)
	}

	// Slack
	slackConf := config.Instance.AdminNotify.Slack
	if slackConf.Enabled {
		slackService := slack.New(slackConf.OauthToken)
		slackService.AddReceivers(slackConf.Receivers...)
		notify.UseServices(slackService)
	}

	// LINE
	LINEConf := config.Instance.AdminNotify.LINE
	if LINEConf.Enabled {
		lineService, _ := line.New(config.Instance.AdminNotify.LINE.ChannelSecret, config.Instance.AdminNotify.LINE.ChannelAccessToken)
		lineService.AddReceivers(LINEConf.Receivers...)
		notify.UseServices(lineService)
	}
}
