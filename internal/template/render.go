package template

import (
	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/internal/entity"
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
	dao        *dao.Dao
	strategy   RenderStrategy
	defaultTpl string
}

func (r *Renderer) Render(notify *entity.Notify, customTpl ...string) string {
	var tpl string
	if len(customTpl) > 0 {
		tpl = customTpl[0]
	} else {
		tpl = r.defaultTpl
	}

	extra := getNotifyExtraData(r.dao, notify)
	params := getCommonParams(r.dao, notify, extra)

	return r.strategy.Render(tpl, params, notify, extra)
}

func NewRenderer(dao *dao.Dao, renderType RenderType, defaultTemplateLoader TemplateLoader) *Renderer {
	r := &Renderer{
		dao: dao,
	}

	// load default template
	if defaultTemplateLoader != nil {
		r.defaultTpl = defaultTemplateLoader.Load(renderType)
	}

	// Render type
	var renderStrategies = map[RenderType]func() RenderStrategy{
		TYPE_EMAIL: func() RenderStrategy {
			return NewEmailRenderer()
		},
		TYPE_NOTIFY: func() RenderStrategy {
			return NewNotifyRenderer()
		},
	}

	// Set render strategy
	if strategyFunc, ok := renderStrategies[renderType]; ok {
		r.strategy = strategyFunc()
	} else {
		r.strategy = NewEmailRenderer() // Default
	}

	return r
}
