package artransfer

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ArtalkJS/Artalk/internal/dao"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/cheggaaa/pb/v3"
)

type ImportParams struct {
	TargetSiteName string `json:"target_site_name" form:"target_site_name" validate:"optional"` // The target site name
	TargetSiteUrl  string `json:"target_site_url" form:"target_site_url" validate:"optional"`   // The target site url
	UrlResolver    bool   `json:"url_resolver" form:"url_resolver" validate:"optional"`         // Enable URL resolver
	JsonFile       string `json:"json_file,omitempty" form:"json_file" validate:"optional"`     // The JSON file path
	JsonData       string `json:"json_data,omitempty" form:"json_data" validate:"optional"`     // The JSON data
	Assumeyes      bool   `json:"assumeyes" form:"assumeyes" validate:"optional"`               // Automatically answer yes for all questions.
}

func importArtrans(dao *dao.Dao, params *ImportParams, comments []*entity.Artran) {
	if len(comments) == 0 {
		logFatal(i18n.T("No comment"))
		return
	}

	if params.TargetSiteUrl != "" && !utils.ValidateURL(params.TargetSiteUrl) {
		logFatal(i18n.T("Invalid {{name}}", map[string]interface{}{"name": i18n.T("Target Site") + " " + "URL"}))
		return
	}

	// 汇总
	print("# " + i18n.T("Please review") + ":\n\n")

	// 第一条评论
	printEncodeData(i18n.T("First comment"), comments[0])

	showTSiteName := params.TargetSiteName
	showTSiteUrl := params.TargetSiteUrl
	if showTSiteName == "" {
		showTSiteName = i18n.T("Unspecified")

	}
	if showTSiteUrl == "" {
		showTSiteUrl = i18n.T("Unspecified")
	}

	// 目标站点名和目标站点URL都不为空，才开启 URL 解析器
	showUrlResolver := "off"
	if params.UrlResolver {
		showUrlResolver = "on"
	}
	// if basic.TargetSiteName != "" && basic.TargetSiteUrl != "" {
	// 	basic.UrlResolver = true
	// 	showUrlResolver = "on"
	// }

	printTable([][]interface{}{
		{i18n.T("Target Site") + " " + i18n.T("Name"), showTSiteName},
		{i18n.T("Target Site") + " URL", showTSiteUrl},
		{i18n.T("Comment count"), fmt.Sprintf("%d", len(comments))},
		{i18n.T("URL Resolver"), showUrlResolver},
	})

	print("\n")

	// 确认开始
	if !params.Assumeyes && !confirm(i18n.T("Confirm to continue?")) {
		os.Exit(0)
	}

	// 准备导入评论
	print("\n")

	importComments := []*entity.Comment{}
	srcIdToIndexMap := map[string]uint{} // 源 ID 映射表 srcID => index
	createdDates := map[int]time.Time{}
	updatedDates := map[int]time.Time{}

	// 解析 comments
	for i, c := range comments {
		srcIdToIndexMap[c.ID] = uint(i + 1) // 防 0 出没
	}

	for i, c := range comments {
		siteName := c.SiteName
		siteUrls := c.SiteUrls

		if params.TargetSiteName != "" {
			siteName = params.TargetSiteName
		}
		if params.TargetSiteUrl != "" {
			siteUrls = params.TargetSiteUrl
		}

		// 准备 site
		site, sErr := siteReady(dao, siteName, siteUrls)
		if sErr != nil {
			logFatal(sErr)
			return
		}

		// 准备 user
		user, err := dao.FindCreateUser(c.Nick, c.Email, c.Link)
		if err == nil && !user.IsAdmin {
			userModified := false
			if c.BadgeName != "" && c.BadgeName != user.BadgeName {
				user.BadgeName = c.BadgeName
				userModified = true
			}
			if c.BadgeColor != "" && c.BadgeColor != user.BadgeColor {
				user.BadgeColor = c.BadgeColor
				userModified = true
			}
			if userModified {
				dao.UpdateUser(&user)
			}
		}

		// 准备 page
		nPageKey := c.PageKey
		if params.UrlResolver { // 使用 URL 解析器
			splittedURLs := utils.SplitAndTrimSpace(params.TargetSiteUrl, ",")
			if len(splittedURLs) == 0 {
				logFatal("[URL Resolver] " + i18n.T("{{name}} cannot be empty", map[string]interface{}{"name": i18n.T("Target Site") + " " + "URL"}))
				return
			}
			nPageKey = urlResolverGetPageKey(splittedURLs[0], c.PageKey)
		}

		page := dao.FindCreatePage(nPageKey, c.PageTitle, site.Name)

		adminOnlyVal := c.PageAdminOnly == utils.ToString(true)
		if page.AdminOnly != adminOnlyVal {
			page.AdminOnly = adminOnlyVal
			dao.UpdatePage(&page)
		}

		voteUp, _ := strconv.Atoi(c.VoteUp)
		voteDown, _ := strconv.Atoi(c.VoteDown)

		// 创建新 comment 实例
		nComment := entity.Comment{
			Rid: srcIdToIndexMap[c.Rid], // [-1-] rid => index+1

			Content: c.Content,

			UA: c.UA,
			IP: c.IP,

			IsCollapsed: c.IsCollapsed == utils.ToString(true),
			IsPending:   c.IsPending == utils.ToString(true),
			IsPinned:    c.IsPinned == utils.ToString(true),

			VoteUp:   voteUp,
			VoteDown: voteDown,

			UserID:   user.ID,
			PageKey:  page.Key,
			SiteName: site.Name,
		}

		// 时间还原
		createdDates[i] = parseDate(c.CreatedAt)
		if c.UpdatedAt != "" {
			updatedDates[i] = parseDate(c.UpdatedAt)
		} else {
			updatedDates[i] = parseDate(c.CreatedAt)
		}

		importComments = append(importComments, &nComment)
	}

	println(i18n.T("Saving") + "...")

	// Batch Insert
	// @link https://gorm.io/docs/create.html#Batch-Insert
	dao.DB().CreateInBatches(&importComments, 100)

	// ID 变更映射表 index => new_db_id
	indexToDbIdMap := map[uint]uint{}
	for i, savedComment := range importComments {
		indexToDbIdMap[uint(i+1)] = savedComment.ID
	}

	// 进度条
	var bar *pb.ProgressBar
	if HttpOutput == nil {
		bar = pb.StartNew(len(comments))
	}

	total := len(comments)

	for i, savedComment := range importComments {
		// 日期恢复
		// @see https://gorm.io/zh_CN/docs/conventions.html#CreatedAt
		// @see https://github.com/go-gorm/gorm/issues/4827#issuecomment-960480148 无语...
		// TODO
		// savedComment.CreatedAt = createdDates[i] // 无效
		// savedComment.UpdatedAt = updatedDates[i]

		updateData := map[string]interface{}{
			"CreatedAt": createdDates[i],
			"UpdatedAt": updatedDates[i],
		}

		// Rid 重建
		if savedComment.Rid != 0 {
			updateData["Rid"] = indexToDbIdMap[savedComment.Rid] // [-2-] index+1 => db_new_id
		}

		dao.DB().Model(&savedComment).Updates(updateData)

		// Vote 重建 (伪投票)
		if savedComment.VoteUp > 0 {
			for i := 0; i < savedComment.VoteUp; i++ {
				dao.NewVote(savedComment.ID, entity.VoteTypeCommentUp, 0, "", "")
			}
		}
		if savedComment.VoteDown > 0 {
			for i := 0; i < savedComment.VoteDown; i++ {
				dao.NewVote(savedComment.ID, entity.VoteTypeCommentDown, 0, "", "")
			}
		}

		if bar != nil {
			bar.Increment()
		}
		if HttpOutput != nil && i%50 == 0 {
			print(fmt.Sprintf("%.0f%%... ", float64(i)/float64(total)*100))
		}
	}

	if bar != nil {
		bar.Finish()
	}
	if HttpOutput != nil {
		println()
	}

	logInfo(i18n.T("{{count}} items imported", map[string]interface{}{"count": len(comments)}))
}

// 站点准备
func siteReady(dao *dao.Dao, tSiteName string, tSiteUrls string) (*entity.Site, error) {
	site := dao.FindSite(tSiteName)
	if site.IsEmpty() {
		// 创建新站点
		site = entity.Site{}
		site.Name = tSiteName
		site.Urls = tSiteUrls
		err := dao.CreateSite(&site)
		if err != nil {
			return nil, fmt.Errorf("failed to create site")
		}
	} else {
		// 追加 URL
		siteCooked := dao.CookSite(&site)

		urlExist := func(tUrl string) bool {
			for _, u := range siteCooked.Urls {
				if u == tUrl {
					return true
				}
			}
			return false
		}

		tUrlsSpit := utils.SplitAndTrimSpace(tSiteUrls, ",")

		rUrls := []string{}
		for _, u := range tUrlsSpit {
			if !urlExist(u) {
				rUrls = append(rUrls, u) // prepend 不存在的站点
			}
		}

		if len(rUrls) > 0 {
			// 保存
			rUrls = append(rUrls, siteCooked.Urls...)
			site.Urls = strings.Join(rUrls, ",")
			err := dao.UpdateSite(&site)
			if err != nil {
				return nil, fmt.Errorf("update site data failed")
			}
		}
	}

	return &site, nil
}
