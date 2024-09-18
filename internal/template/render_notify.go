package template

import (
	"github.com/artalkjs/artalk/v2/internal/entity"
)

var _ RenderStrategy = (*NotifyRender)(nil)

type NotifyRender struct {
}

func NewNotifyRenderer() *NotifyRender {
	return &NotifyRender{}
}

func (r *NotifyRender) Render(tpl string, p tplParams, notify *entity.Notify, extra notifyExtraData) string {
	p.Content = handleEmoticonsImgTagsForNotify(extra.to.ContentRaw)
	p.ReplyContent = handleEmoticonsImgTagsForNotify(extra.from.ContentRaw)

	return renderCommon(tpl, p)
}
