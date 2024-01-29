package notify_pusher

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
)

func (pusher *NotifyPusher) checkNeedSendEmailToUser(comment *entity.Comment, parentComment *entity.Comment) bool {
	// 自己回复自己，不提醒
	if comment.UserID == parentComment.UserID {
		return false
	}

	// 待审状态评论回复不邮件通知 (管理员审核通过后才发送)
	if comment.IsPending {
		return false
	}

	// 对方个人设定关闭邮件接收
	if !pusher.dao.FetchUserForComment(parentComment).ReceiveEmail {
		return false
	}

	// 对方是管理员，但是管理员邮件接收关闭 (用于开启多元推送后禁用邮件通知管理员)
	if pusher.dao.FetchUserForComment(parentComment).IsAdmin && !pusher.conf.Email.Enabled {
		return false
	}

	return true
}

func (pusher *NotifyPusher) checkNeedSendEmailToAdmin(comment *entity.Comment, parentComment *entity.Comment, admin *entity.User) bool {
	// 配置文件关闭管理员邮件接收
	if !pusher.conf.Email.Enabled {
		return false
	}

	// 待审评论不发送通知
	if comment.IsPending && !pusher.conf.NotifyPending {
		return false
	}

	// 管理员自己回复自己，不提醒
	if comment.UserID == admin.ID {
		return false
	}

	// 用户回复对象是该管理员，不提醒
	// (避免当 NoiseModeOn = true 时，重复发送)
	if parentComment.UserID == admin.ID {
		return false
	}

	// 管理员评论不回复给其他管理员
	if pusher.dao.FetchUserForComment(comment).IsAdmin {
		return false
	}

	// 只发送给对应站点管理员
	// TODO admin.SiteNames had removed, so temporarily disabled
	// if admin.SiteNames != "" && !slices.Contains(pusher.dao.CookUser(admin).SiteNames, comment.SiteName) {
	// 	return false
	// }

	// 该管理员单独设定关闭接收邮件
	if !admin.ReceiveEmail {
		return false
	}

	return true
}

func (pusher *NotifyPusher) checkNeedMultiPush(comment *entity.Comment, pComment *entity.Comment) bool {
	isRootComment := pComment == nil || pComment.IsEmpty()

	// 忽略来自管理员的评论
	coUser := pusher.dao.FetchUserForComment(comment)
	if coUser.IsAdmin {
		return false
	}

	// 待审评论不发送通知
	if comment.IsPending && !pusher.conf.NotifyPending {
		return false
	}

	// 非嘈杂模式
	if !pusher.conf.NoiseMode {
		// 如果不是 root 评论 且 回复目标不是管理员，直接忽略
		if !isRootComment && !pusher.dao.FetchUserForComment(pComment).IsAdmin {
			return false
		}
	}

	return true
}
