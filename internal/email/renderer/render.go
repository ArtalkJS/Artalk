package renderer

import (
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
)

type RenderType string

const (
	TYPE_EMAIL  RenderType = "email"
	TYPE_NOTIFY RenderType = "notify"
)

type RenderStrategy interface {
	Render(tpl string, p tplParams, notify *entity.Notify, extra notifyExtraData) string
}

type Renderer struct {
	dao      *dao.Dao
	strategy RenderStrategy
	tpl      string
}

func (r *Renderer) Render(notify *entity.Notify, customTpl ...string) string {
	var tpl string
	if len(customTpl) > 0 {
		tpl = customTpl[0]
	} else {
		tpl = r.tpl
	}

	extra := getNotifyExtraData(r.dao, notify)
	params := getCommonParams(r.dao, notify, extra)

	return r.strategy.Render(tpl, params, notify, extra)
}

func NewRenderer(dao *dao.Dao, renderType RenderType, defaultTplName string) *Renderer {
	renderer := &Renderer{
		dao: dao,
	}

	// Retrieve template
	renderer.tpl = getTplByName(renderType, defaultTplName)

	// Render type
	var renderStrategies = map[RenderType]func() RenderStrategy{
		TYPE_EMAIL: func() RenderStrategy {
			return NewEmailRenderStrategy()
		},
		TYPE_NOTIFY: func() RenderStrategy {
			return NewNotifyRenderer()
		},
	}

	// Set render strategy
	if strategyFunc, ok := renderStrategies[renderType]; ok {
		renderer.strategy = strategyFunc()
	} else {
		renderer.strategy = NewEmailRenderStrategy() // Default
	}

	return renderer
}
