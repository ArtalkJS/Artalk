package notify_pusher

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
)

// 通知发送 (from comment to parentComment)
func (pusher *NotifyPusher) Push(comment *entity.Comment, pComment *entity.Comment) {
	isRootComment := pComment == nil || pComment.IsEmpty()
	isAdminNoiseModeOn := pusher.conf.NoiseMode

	// ==============
	//  邮件回复对方
	// ==============
	if !isRootComment {
		pusher.emailToUser(comment, pComment)
	}

	// ==============
	//  邮件通知管理员
	// ==============
	if isRootComment || isAdminNoiseModeOn {
		pusher.emailToAdmins(comment, pComment)
	}

	// 管理员多元推送
	pusher.multiPush(comment, pComment)
}
