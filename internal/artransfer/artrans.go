package artransfer

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ArtalkJS/Artalk/internal/db"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/cheggaaa/pb/v3"
)

var ArtransImporter = &_ArtransImporter{
	ImporterInfo: ImporterInfo{
		Name: "artrans",
		Desc: "Import from Artrans",
		Note: "",
	},
}

type _ArtransImporter struct {
	ImporterInfo
}

func (imp *_ArtransImporter) Run(basic *BasicParams, payload []string) {
	// 读取文件
	jsonStr, jErr := JsonFileReady(payload)
	if jErr != nil {
		logFatal(jErr)
		return
	}

	ImportArtransByStr(basic, jsonStr)
}

func ImportArtransByStr(basic *BasicParams, str string) {
	// 解析内容
	comments := []entity.Artran{}
	dErr := JsonDecodeFAS(str, &comments)
	if dErr != nil {
		logFatal(dErr)
		return
	}

	ImportArtrans(basic, comments)
}

func ImportArtrans(basic *BasicParams, srcComments []entity.Artran) {
	if len(srcComments) == 0 {
		logFatal(i18n.T("No comment"))
		return
	}

	if basic.TargetSiteUrl != "" && !utils.ValidateURL(basic.TargetSiteUrl) {
		logFatal(i18n.T("Invalid {{name}}", map[string]interface{}{"name": i18n.T("Target Site") + " " + "URL"}))
		return
	}

	// 汇总
	print("# " + i18n.T("Please review") + ":\n\n")

	// 第一条评论
	PrintEncodeData(i18n.T("First comment"), srcComments[0])

	showTSiteName := basic.TargetSiteName
	showTSiteUrl := basic.TargetSiteUrl
	if showTSiteName == "" {
		showTSiteName = i18n.T("Unspecified")

	}
	if showTSiteUrl == "" {
		showTSiteUrl = i18n.T("Unspecified")
	}

	// 目标站点名和目标站点URL都不为空，才开启 URL 解析器
	showUrlResolver := "off"
	if basic.UrlResolver {
		showUrlResolver = "on"
	}
	// if basic.TargetSiteName != "" && basic.TargetSiteUrl != "" {
	// 	basic.UrlResolver = true
	// 	showUrlResolver = "on"
	// }

	PrintTable([][]interface{}{
		{i18n.T("Target Site") + " " + i18n.T("Name"), showTSiteName},
		{i18n.T("Target Site") + " URL", showTSiteUrl},
		{i18n.T("Comment count"), fmt.Sprintf("%d", len(srcComments))},
		{i18n.T("URL Resolver"), showUrlResolver},
	})

	print("\n")

	// 确认开始
	if !Confirm(i18n.T("Confirm to continue?")) {
		os.Exit(0)
	}

	// 准备导入评论
	print("\n")

	importComments := []entity.Comment{}
	srcIdToIndexMap := map[string]uint{} // 源 ID 映射表 srcID => index
	createdDates := map[int]time.Time{}
	updatedDates := map[int]time.Time{}

	// 解析 comments
	for i, c := range srcComments {
		srcIdToIndexMap[c.ID] = uint(i + 1) // 防 0 出没
	}

	for i, c := range srcComments {
		siteName := c.SiteName
		siteUrls := c.SiteUrls

		if basic.TargetSiteName != "" {
			siteName = basic.TargetSiteName
		}
		if basic.TargetSiteUrl != "" {
			siteUrls = basic.TargetSiteUrl
		}

		// 准备 site
		site, sErr := SiteReady(siteName, siteUrls)
		if sErr != nil {
			logFatal(sErr)
			return
		}

		// 准备 user
		user := query.FindCreateUser(c.Nick, c.Email, c.Link)
		if !user.IsAdmin {
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
				query.UpdateUser(&user)
			}
		}

		// 准备 page
		nPageKey := c.PageKey
		if basic.UrlResolver { // 使用 URL 解析器
			splittedURLs := utils.SplitAndTrimSpace(basic.TargetSiteUrl, ",")
			nPageKey = UrlResolverGetPageKey(splittedURLs[0], c.PageKey)
		}

		page := query.FindCreatePage(nPageKey, c.PageTitle, site.Name)

		adminOnlyVal := c.PageAdminOnly == utils.ToString(true)
		if page.AdminOnly != adminOnlyVal {
			page.AdminOnly = adminOnlyVal
			query.UpdatePage(&page)
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
		createdDates[i] = ParseDate(c.CreatedAt)
		if c.UpdatedAt != "" {
			updatedDates[i] = ParseDate(c.UpdatedAt)
		} else {
			updatedDates[i] = ParseDate(c.CreatedAt)
		}

		importComments = append(importComments, nComment)
	}

	println(i18n.T("Saving") + "...")

	// Batch Insert
	// @link https://gorm.io/docs/create.html#Batch-Insert
	db.DB().CreateInBatches(&importComments, 100)

	// ID 变更映射表 index => new_db_id
	indexToDbIdMap := map[uint]uint{}
	for i, savedComment := range importComments {
		indexToDbIdMap[uint(i+1)] = savedComment.ID
	}

	// 进度条
	var bar *pb.ProgressBar
	if HttpOutput == nil {
		bar = pb.StartNew(len(srcComments))
	}

	total := len(srcComments)

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

		db.DB().Model(&savedComment).Updates(updateData)

		// Vote 重建 (伪投票)
		if savedComment.VoteUp > 0 {
			for i := 0; i < savedComment.VoteUp; i++ {
				query.NewVote(savedComment.ID, entity.VoteTypeCommentUp, 0, "", "")
			}
		}
		if savedComment.VoteDown > 0 {
			for i := 0; i < savedComment.VoteDown; i++ {
				query.NewVote(savedComment.ID, entity.VoteTypeCommentDown, 0, "", "")
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

	logInfo(i18n.T("{{count}} items imported", map[string]interface{}{"count": len(srcComments)}))
}
