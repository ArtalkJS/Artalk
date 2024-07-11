package notify_pusher

import (
	"context"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/dingding"
	"github.com/nikoksr/notify/service/line"
	"github.com/nikoksr/notify/service/slack"
	"github.com/nikoksr/notify/service/telegram"
)

type NotifyPusherConf struct {
	config.AdminNotifyConf
	Dao *dao.Dao

	// Provide a custom function to bridge the gap between Notify pusher and Email pusher
	EmailPush func(notify *entity.Notify) error
}

type NotifyPusher struct {
	conf   *NotifyPusherConf
	dao    *dao.Dao
	ctx    context.Context
	helper *notify.Notify
}

func NewNotifyPusher(conf *NotifyPusherConf) *NotifyPusher {
	pusher := &NotifyPusher{
		conf:   conf,
		dao:    conf.Dao,
		ctx:    context.Background(),
		helper: notify.New(),
	}

	pusher.loadHelper()

	return pusher
}

func (pusher *NotifyPusher) loadHelper() {
	var (
		helper = pusher.helper
		conf   = pusher.conf
	)

	// Telegram
	tgConf := conf.Telegram
	if tgConf.Enabled {
		if telegramService, err := telegram.New(tgConf.ApiToken); err == nil {
			telegramService.AddReceivers(tgConf.Receivers...)
			helper.UseServices(telegramService)
		} else {
			log.Error("[Notify] Telegram service init error: ", err)
		}
	}

	// 钉钉
	dingTalkConf := conf.DingTalk
	if dingTalkConf.Enabled {
		dingTalkService := dingding.New(&dingding.Config{Token: dingTalkConf.Token, Secret: dingTalkConf.Secret})
		helper.UseServices(dingTalkService)
	}

	// Slack
	slackConf := conf.Slack
	if slackConf.Enabled {
		slackService := slack.New(slackConf.OauthToken)
		slackService.AddReceivers(slackConf.Receivers...)
		helper.UseServices(slackService)
	}

	// LINE
	LINEConf := conf.LINE
	if LINEConf.Enabled {
		if lineService, err := line.New(pusher.conf.LINE.ChannelSecret, pusher.conf.LINE.ChannelAccessToken); err == nil {
			lineService.AddReceivers(LINEConf.Receivers...)
			helper.UseServices(lineService)
		} else {
			log.Error("[Notify] LINE service init error: ", err)
		}
	}
}
