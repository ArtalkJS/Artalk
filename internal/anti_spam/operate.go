package anti_spam

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/log"
)

func (as AntiSpam) getEnabledCheckers() []Checker {
	checkers := []Checker{}

	akismetKey := strings.TrimSpace(as.conf.AkismetKey)
	if akismetKey != "" {
		checkers = append(checkers, NewAkismetChecker(akismetKey))
	}

	// 腾讯云
	tencentConf := as.conf.Tencent
	if tencentConf.Enabled {
		checkers = append(checkers, NewTencentChecker(
			tencentConf.SecretID, tencentConf.SecretKey, tencentConf.Region))
	}

	// 阿里云
	aliyunConf := as.conf.Aliyun
	if aliyunConf.Enabled {
		checkers = append(checkers, NewAliyunChecker(
			aliyunConf.AccessKeyID, aliyunConf.AccessKeySecret, aliyunConf.Region))
	}

	// 关键字过滤
	keywordsConf := as.conf.Keywords
	if keywordsConf.Enabled {

		var kwCheckerMode KwCheckerMode
		if as.conf.Keywords.Pending {
			kwCheckerMode = KwCheckerModeBlock
		} else {
			kwCheckerMode = KwCheckerModeReplace
		}

		checkers = append(checkers, NewKeywordsChecker(&KeywordsCheckerConf{
			Files:     as.conf.Keywords.Files,
			FileSep:   as.conf.Keywords.FileSep,
			ReplaceTo: as.conf.Keywords.ReplacTo,
			Mode:      kwCheckerMode,
			OnUpdateComment: func(commentID uint, content string) error {
				comment := as.dao.FindComment(commentID)
				comment.Content = content
				as.dao.UpdateComment(&comment)
				return nil
			},
		}))

	}

	return checkers
}

// 拦截评论
func (as AntiSpam) blockComment(checker Checker, comment *entity.Comment) {
	log.Info(LOG_TAG, fmt.Sprintf("[%s] Successful blocking of comments ID=%d CONT=%s",
		checker.Name(), comment.ID, strconv.Quote(comment.Content)))

	if comment.IsPending {
		return
	}

	comment.IsPending = true // 改为待审状态
	as.dao.UpdateComment(comment)
}

func (as AntiSpam) getCheckerParams(data *CheckData) *CheckerParams {
	user := as.dao.FetchUserForComment(data.Comment)
	siteURL := ""

	if data.Comment.SiteName != "" {
		site := as.dao.FindSite(data.Comment.SiteName)
		siteURL = as.dao.CookSite(&site).FirstUrl
	}
	if siteURL == "" { // 从 referer 中提取网站
		if pr, err := url.Parse(data.ReqReferer); err == nil && pr.Scheme != "" && pr.Host != "" {
			siteURL = fmt.Sprintf("%s://%s", pr.Scheme, pr.Host)
		}
	}

	return &CheckerParams{
		BlogURL: siteURL,

		Content:   data.Comment.Content,
		CommentID: data.Comment.ID,

		UserName:  user.Name,
		UserEmail: user.Email,
		UserID:    user.ID,
		UserIP:    data.ReqIP,
		UserAgent: data.ReqUserAgent,
	}
}

func (as AntiSpam) checkerTrigger(checker Checker, data *CheckData) {
	params := as.getCheckerParams(data)
	isPass, err := checker.Check(params)

	if err != nil { // Api 发生错误
		log.Error(LOG_TAG, fmt.Sprintf(
			"%s Interception error occurred ID=%d Err:", checker.Name(), data.Comment.ID), err)

		if as.conf.ApiFailBlock {
			as.blockComment(checker, data.Comment) // 仍然拦截
		}

		return
	}

	if !isPass { // not Pass 且 Api 未发生错误
		as.blockComment(checker, data.Comment) // 拦截评论
	}
}
