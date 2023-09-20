package template

import (
	"github.com/ArtalkJS/Artalk/internal/entity"
)

var _ RenderStrategy = (*EmailRenderStrategy)(nil)

type EmailRenderStrategy struct {
}

func NewEmailRenderStrategy() *EmailRenderStrategy {
	return &EmailRenderStrategy{}
}

func (r *EmailRenderStrategy) Render(tpl string, p tplParams, notify *entity.Notify, extra notifyExtraData) string {
	return renderCommon(tpl, p)
}
