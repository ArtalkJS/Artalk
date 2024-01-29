package notify_pusher

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
	"golang.org/x/exp/slices"
)

func (pusher *NotifyPusher) sendEmail(notify *entity.Notify) {
	if pusher.conf.EmailPush != nil {
		pusher.conf.EmailPush(notify)
	}
}

func (pusher *NotifyPusher) emailToUser(comment *entity.Comment, pComment *entity.Comment) {
	if !pusher.checkNeedSendEmailToUser(comment, pComment) {
		log.Debug("ignore email notify by pusher.checkNeedSendEmailToUser")
		return
	}

	notify := pusher.dao.FindCreateNotify(pComment.UserID, comment.ID)
	pusher.dao.NotifySetInitial(&notify)

	// 邮件通知
	pusher.sendEmail(&notify)
}

func (pusher *NotifyPusher) emailToAdmins(comment *entity.Comment, pComment *entity.Comment) {
	toAddrSent := []string{} // 记录已发送的收件人地址（避免重复发送）
	for _, admin := range pusher.dao.GetAllAdmins() {
		// 该管理员地址已曾发送，避免重复发送
		if slices.Contains(toAddrSent, admin.Email) {
			continue
		}
		toAddrSent = append(toAddrSent, admin.Email)

		if !pusher.checkNeedSendEmailToAdmin(comment, pComment, &admin) {
			log.Debug("ignore email notify by pusher.checkNeedSendEmailToAdmin")
			continue
		}

		notify := pusher.dao.FindCreateNotify(admin.ID, comment.ID)
		pusher.dao.NotifySetInitial(&notify)

		// 发送邮件给管理员
		pusher.sendEmail(&notify)
	}
}
