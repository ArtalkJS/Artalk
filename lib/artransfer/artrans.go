package artransfer

import (
	"fmt"
	"os"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/cheggaaa/pb/v3"
	"github.com/jedib0t/go-pretty/v6/table"
)

var ArtransImporter = &_ArtransImporter{
	ImporterInfo: ImporterInfo{
		Name: "artrans",
		Desc: "从 Artrans 导入数据",
		Note: "",
	},
}

type _ArtransImporter struct {
	ImporterInfo
}

func (imp *_ArtransImporter) Run(basic *BasicParams, payload []string) {
	err := RequiredBasicTargetSite(basic)
	if err != nil {
		logFatal(err)
		return
	}

	// 读取文件
	jsonStr, jErr := JsonFileReady(payload)
	if jErr != nil {
		logFatal(jErr)
		return
	}

	basic.UrlResolver = false

	ImportArtransByStr(basic, jsonStr)
}

func ImportArtransByStr(basic *BasicParams, str string) {
	// 解析内容
	comments := []model.Artran{}
	dErr := JsonDecodeFAS(str, &comments)
	if dErr != nil {
		logFatal(dErr)
		return
	}

	ImportArtrans(basic, comments)
}

func ImportArtrans(basic *BasicParams, comments []model.Artran) {
	if len(comments) == 0 {
		logFatal("未读取到任何一条评论")
		return
	}

	// 汇总
	print("# 请过目：\n\n")

	// 第一条评论
	PrintEncodeData("第一条评论", comments[0])

	showTSiteName := basic.TargetSiteName
	showTSiteUrl := basic.TargetSiteUrl
	if showTSiteName == "" {
		showTSiteName = "未指定"
	}
	if showTSiteUrl == "" {
		showTSiteUrl = "未指定"
	}

	PrintTable([]table.Row{
		{"目标站点名", showTSiteName},
		{"目标站点 URL", showTSiteUrl},
		{"评论数量", len(comments)},
	})

	print("\n")

	// 确认开始
	if !Confirm("确认开始导入吗？") {
		os.Exit(0)
	}

	// 准备导入评论
	print("\n")

	// 执行导入
	idMap := map[string]int{}    // ID 映射表 object_id => id
	idChanges := map[uint]uint{} // ID 变更表 original_id => new_db_id

	// 生成 ID 映射表
	id := 1
	for _, c := range comments {
		idMap[c.ID] = id
		id++
	}

	// 进度条
	bar := pb.StartNew(len(comments))

	// 遍历导入 comments
	for _, c := range comments {
		// 准备 site
		site, sErr := SiteReady(c.SiteName, c.SiteUrls)
		if sErr != nil {
			logFatal(sErr)
			return
		}

		// 准备 user
		user := model.FindCreateUser(c.Nick, c.Email, c.Link)
		if c.Password != "" {
			user.Password = c.Password
		}
		if c.BadgeName != "" {
			user.BadgeName = c.BadgeName
		}
		if c.BadgeColor != "" {
			user.BadgeColor = c.BadgeColor
		}
		model.UpdateUser(&user)

		// 准备 page
		nPageKey := c.PageKey
		if basic.UrlResolver { // 使用 URL 解析器
			nPageKey = UrlResolverGetPageKey(basic.TargetSiteUrl, c.PageKey)
		}

		page := model.FindCreatePage(nPageKey, c.PageTitle, c.SiteName)
		page.AdminOnly = c.PageAdminOnly == lib.ToString(true)
		model.UpdatePage(&page)

		// 创建新 comment 实例
		nComment := model.Comment{
			Rid: uint(idMap[c.Rid]),

			Content: c.Content,

			UA: c.UA,
			IP: c.IP,

			IsCollapsed: c.IsCollapsed == lib.ToString(true),
			IsPending:   c.IsPending == lib.ToString(true),

			UserID:   user.ID,
			PageKey:  page.Key,
			SiteName: site.Name,
		}

		// 保存到数据库
		dErr := lib.DB.Create(&nComment).Error
		if dErr != nil {
			logError(fmt.Sprintf("评论源 ID:%s 保存失败", c.ID))
			continue
		}

		// 日期恢复
		// @see https://gorm.io/zh_CN/docs/conventions.html#CreatedAt
		lib.DB.Model(&nComment).Update("CreatedAt", ParseDate(c.CreatedAt))
		lib.DB.Model(&nComment).Update("UpdatedAt", ParseDate(c.UpdatedAt))

		idChanges[uint(idMap[c.ID])] = nComment.ID

		bar.Increment()
	}
	bar.Finish()
	logInfo(fmt.Sprintf("导入 %d 条数据", len(comments)))

	// reply id 重建
	RebuildRid(idChanges)
}
