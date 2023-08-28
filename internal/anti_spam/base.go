package anti_spam

import (
	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
)

const LOG_TAG = "[Spam Interception]"

type AntiSpamConf struct {
	config.ModeratorConf
	Dao *dao.Dao
}

type AntiSpam struct {
	conf *AntiSpamConf
	dao  *dao.Dao
}

func NewAntiSpam(conf *AntiSpamConf) *AntiSpam {
	return &AntiSpam{
		conf: conf,
		dao:  conf.Dao,
	}
}

type CheckData struct {
	Comment      *entity.Comment
	ReqReferer   string
	ReqIP        string
	ReqUserAgent string
}

func (as AntiSpam) CheckAndBlock(data *CheckData) {
	checkers := as.getEnabledCheckers()

	// 执行检查
	for _, checker := range checkers {
		as.checkerTrigger(checker, data)
	}
}
