package template

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
)

var _ RenderStrategy = (*EmailRender)(nil)

type EmailRender struct {
}

func NewEmailRenderer() *EmailRender {
	return &EmailRender{}
}

func (r *EmailRender) Render(tpl string, p tplParams, notify *entity.Notify, extra notifyExtraData) string {
	return renderCommon(tpl, p)
}
