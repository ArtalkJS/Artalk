package template

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
)

var _ RenderStrategy = (*NotifyRenderStrategy)(nil)

type NotifyRenderStrategy struct {
}

func NewNotifyRenderer() *NotifyRenderStrategy {
	return &NotifyRenderStrategy{}
}

func (r *NotifyRenderStrategy) Render(tpl string, p tplParams, notify *entity.Notify, extra notifyExtraData) string {
	p.Content = handleEmoticonsImgTagsForNotify(extra.to.ContentRaw)
	p.ReplyContent = handleEmoticonsImgTagsForNotify(extra.from.ContentRaw)

	return renderCommon(tpl, p)
}
