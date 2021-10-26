package importer

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/elliotchance/phpserialize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

const (
	// 已测试可用版本
	TypechoTestedVerMain      = "1.1"
	TypechoTestedVerSub       = "17.10.30"
	TypechoRewritePostDefault = "/archives/{cid}/"
	TypechoRewritePageDefault = "/{slug}.html"
)

var TypechoImporter = &_TypechoImporter{
	ImporterInfo: ImporterInfo{
		Name: "typecho",
		Desc: "从 Typecho 导入数据",
		Note: fmt.Sprintf("测试兼容：Typecho %s/%s", TypechoTestedVerMain, TypechoTestedVerSub),
	},
}

type _TypechoImporter struct {
	ImporterInfo
	TargetSite model.Site

	Basic         *BasicParams
	SrcVersion    string
	SrcSiteName   string
	SrcSiteURL    string
	Comments      []TypechoComment
	Options       []TypechoOption
	Contents      []TypechoContent
	Metas         []TypechoMeta
	Relationships []TypechoRelationship

	DbPrefix string

	RewritePost string
	RewritePage string

	OptionRoutingTable map[string]TypechoRoute
}

// Typecho 升级相关的代码 @see https://github.com/typecho/typecho/blob/64b8e686885d8ab4c7f0cdc3d6dc2d99fa48537c/var/Utils/Upgrade.php
// 路由 @see https://github.com/typecho/typecho/blob/530312443142577509df88ce88cf3274fac9b8c4/var/Widget/Options/Permalink.php#L319
// DB @see https://github.com/typecho/typecho/blob/6558fd5e030a950335d53038f82728b06ad6c32d/install/Mysql.sql
func (imp *_TypechoImporter) Run(basic *BasicParams, payload []string) {
	// Ready
	imp.Basic = basic
	typechoDB := DbReady(payload)

	GetParamsFrom(payload).To(map[string]*string{
		"prefix":       &imp.DbPrefix,
		"rewrite_post": &imp.RewritePost,
		"rewrite_page": &imp.RewritePage,
	})

	if imp.DbPrefix == "" {
		imp.DbPrefix = "typecho_"
	}

	// Load Options
	tbOptions := imp.DbPrefix + "options"
	typechoDB.Raw("SELECT * FROM " + tbOptions).Scan(&imp.Options)
	logrus.Info(fmt.Sprintf("从数据表 `%s` 获取 %d 条记录", tbOptions, len(imp.Options)))

	for _, opt := range imp.Options {
		switch opt.Name {
		case "generator":
			imp.SrcVersion = opt.Value
		case "title":
			imp.SrcSiteName = opt.Value
		case "siteUrl":
			imp.SrcSiteURL = opt.Value
		}
	}

	// 若未指定 target site 则使用数据库数据
	if basic.TargetSiteName == "" && imp.SrcSiteName != "" {
		basic.TargetSiteName = imp.SrcSiteName
	}
	if basic.TargetSiteUrl == "" && imp.SrcSiteURL != "" && lib.ValidateURL(imp.SrcSiteURL) {
		basic.TargetSiteUrl = imp.SrcSiteURL
	}
	RequiredBasicTargetSite(basic)

	// 检查数据源版本号
	imp.VersionCheck()

	// 重写规则
	imp.RewriteRuleReady()

	// load Metas
	tbMetas := imp.DbPrefix + "metas"
	typechoDB.Raw("SELECT * FROM " + tbMetas).Scan(&imp.Metas)
	logrus.Info(fmt.Sprintf("从数据表 `%s` 获取 %d 条记录", tbMetas, len(imp.Metas)))

	// load Relationships
	tbRelationships := imp.DbPrefix + "relationships"
	typechoDB.Raw("SELECT * FROM " + tbRelationships).Scan(&imp.Relationships)
	logrus.Info(fmt.Sprintf("从数据表 `%s` 获取 %d 条记录", tbRelationships, len(imp.Relationships)))

	// 获取 contents
	tbContents := imp.DbPrefix + "contents"
	typechoDB.Raw(fmt.Sprintf("SELECT * FROM %s ORDER BY created ASC", tbContents)).Scan(&imp.Contents)
	logrus.Info(fmt.Sprintf("从数据表 `%s` 获取 %d 条记录", tbContents, len(imp.Contents)))

	// 获取 comments
	tbComments := imp.DbPrefix + "comments"
	typechoDB.Raw(fmt.Sprintf("SELECT * FROM %s ORDER BY created ASC", tbComments)).Scan(&imp.Comments)
	logrus.Info(fmt.Sprintf("从数据表 `%s` 获取 %d 条记录", tbComments, len(imp.Comments)))

	// 导入前参数汇总
	fmt.Print("\n")

	fmt.Print("# 请过目：\n\n")

	// 显示第一条数据
	pageKeyEgPost := ""
	pageKeyEgPage := ""
	for _, c := range imp.Contents {
		if c.Type == "post" {
			firstPost := fmt.Sprintf("[第一篇文章]\n\n"+
				"    %#v\n\n", imp.Contents[0])

			fmt.Print(HideJsonLongText("Text", firstPost))
			pageKeyEgPost = imp.GetNewPageKey(c)
			fmt.Printf(" -> 生成 PageKey: %#v\n\n", pageKeyEgPost)
			break
		}
	}

	for _, c := range imp.Contents {
		if c.Type == "page" {
			pageKeyEgPage = imp.GetNewPageKey(c)
			break
		}
	}

	if len(imp.Comments) > 0 {
		fmt.Printf("[第一条评论]\n\n"+
			"    %#v\n\n", imp.Comments[0])
	}

	PrintTable([]table.Row{
		{"[基本信息]", "读取数据", "导入目标"},
		{"站点名称", fmt.Sprintf("%#v", imp.SrcSiteName), fmt.Sprintf("%#v", imp.Basic.TargetSiteName)},
		{"BaseURL", fmt.Sprintf("%#v", imp.SrcSiteURL), fmt.Sprintf("%#v", imp.Basic.TargetSiteUrl)},
		{"版本号", fmt.Sprintf("%v", imp.SrcVersion), fmt.Sprintf("ArtalkGo %v", lib.Version+"/"+lib.CommitHash+"")},
		{"数量统计", fmt.Sprintf("评论: %d", len(imp.Comments)), fmt.Sprintf("页面: %d", len(imp.Contents))},
	})

	PrintTable([]table.Row{
		{"[重写规则]", "用于生成 pageKey (评论页面唯一标识)", "生成示例"},
		{"文章页面", fmt.Sprintf("%#v", imp.RewritePost), pageKeyEgPost},
		{"独立页面", fmt.Sprintf("%#v", imp.RewritePage), pageKeyEgPage},
	})

	fmt.Print("\n")
	fmt.Println("若以上内容不符合预期，请向我们反馈：https://artalk.js.org/guide/transfer.html")
	fmt.Print("\n")

	// 确认开始
	if !Confirm("确认开始导入吗？") {
		os.Exit(0)
	}

	// 准备导入评论
	fmt.Println()

	// 准备新的 site
	imp.TargetSite = SiteReady(basic)

	// 开始执行导入
	imp.ImportComments()
}

// 导入评论
func (imp *_TypechoImporter) ImportComments() {
	contents := imp.Contents
	comments := imp.Comments

	fmt.Print("\n====================================\n\n")
	logrus.Info(fmt.Sprintf("[开始导入] 共 %d 个页面，%d 条评论", len(contents), len(comments)))
	fmt.Print("\n")

	siteName := imp.Basic.TargetSiteName

	idChanges := map[uint]uint{} // comment_id: old => new

	// 遍历 Contents
	for _, c := range contents {
		// 创建页面
		pageKey := imp.GetNewPageKey(c)
		page := model.FindCreatePage(pageKey, c.Title, imp.Basic.TargetSiteName)

		// 查询评论
		commentTotal := 0
		for _, co := range comments {
			if co.Cid != c.Cid {
				continue
			}

			// 创建 user
			user := model.FindCreateUser(co.Author, co.Mail)

			if co.Url != "" {
				user.Link = co.Url
			}
			model.UpdateUser(&user)

			// 创建新 comment 实例
			nComment := model.Comment{
				Content: co.Text,

				PageKey:  page.Key,
				SiteName: siteName,

				UserID: user.ID,
				UA:     co.Agent,
				IP:     co.Ip,

				Rid: uint(co.Parent),

				IsCollapsed: false,
				IsPending:   co.Status != "approved",
			}

			// 日期恢复
			createdVal := fmt.Sprintf("%v", co.Created)
			nComment.CreatedAt = ParseDate(createdVal)
			nComment.UpdatedAt = ParseDate(createdVal)

			// 保存到数据库
			err := lib.DB.Create(&nComment).Error
			if err != nil {
				logrus.Error(fmt.Sprintf("评论源 ID:%d 保存失败", co.Coid))
				continue
			}

			idChanges[uint(co.Coid)] = nComment.ID
			commentTotal++
		}

		fmt.Printf("+ [%-3d] 条评论 <- [%5d] %-30s | %#v\n", commentTotal, c.Cid, c.Slug, c.Title)
	}

	// reply id 重建
	for _, newId := range idChanges {
		nComment := model.FindComment(newId, siteName)
		if nComment.Rid == 0 {
			continue
		}
		if newId, isExist := idChanges[nComment.Rid]; isExist {
			nComment.Rid = newId
			err := lib.DB.Save(&nComment).Error
			if err != nil {
				logrus.Error(fmt.Sprintf("[rid 更新] new_id:%d new_rid:%d", nComment.ID, newId), err)
			}
		}
	}

	fmt.Print("\n")
	logrus.Info("RID 重构完毕")
}

// 获取新的 PageKey (根据重写规则)
func (imp *_TypechoImporter) GetNewPageKey(content TypechoContent) string {
	date := ParseDate(fmt.Sprintf("%v", content.Created))

	// 替换内容制作
	replaces := map[string]string{
		"cid":      fmt.Sprintf("%v", content.Cid),
		"slug":     content.Slug,
		"category": imp.GetContentCategory(content.Cid),
		"year":     fmt.Sprintf("%v", date.Local().Year()),
		"month":    fmt.Sprintf("%v", date.Local().Month()),
		"day":      fmt.Sprintf("%v", date.Local().Day()),
	}

	baseUrl := imp.Basic.TargetSiteUrl
	baseUrl = strings.TrimSuffix(baseUrl, "/")

	rewriteRule := imp.RewritePost
	if strings.HasPrefix(content.Type, "post") {
		rewriteRule = imp.RewritePost
	} else if strings.HasPrefix(content.Type, "page") {
		rewriteRule = imp.RewritePage
	}

	rewriteRule = "/" + strings.TrimPrefix(rewriteRule, "/")
	return baseUrl + imp.ReplaceAllBracketed(rewriteRule, replaces)
}

// 替换字符串
// @note 特别注意：replaces 的 key 外侧无 “[]” 或 “{}”
func (imp *_TypechoImporter) ReplaceAllBracketed(data string, replaces map[string]string) string {
	r := regexp.MustCompile(`(\[|\{)\s*(.*?)\s*(\]|\})`) // 同时支持 {} 和 []
	return r.ReplaceAllStringFunc(data, func(m string) string {
		key := r.FindStringSubmatch(m)[2]
		if val, isExist := replaces[key]; isExist {
			return val
		} else {
			logrus.Error(fmt.Sprintf("[重写规则] \"%s\" 变量无效", key))
		}
		return m
	})
}

// 版本过旧检测
func (imp *_TypechoImporter) VersionCheck() {
	r := regexp.MustCompile(`Typecho[\s]*(.+)\/(.+)`)
	group := r.FindStringSubmatch(imp.SrcVersion)
	if len(group) < 2 {
		logrus.Warn(`无法确认您的 Typecho 版本号："` + fmt.Sprintf("%#v", imp.SrcVersion) + `"`)
		return
	}

	verMain := ParseVersion(group[1])
	verSub := ParseVersion(group[2])

	if verMain < ParseVersion(TypechoTestedVerMain) ||
		verSub < ParseVersion(TypechoTestedVerSub) {
		fmt.Print("\n")
		logrus.Warn(fmt.Sprintf("Typecho 当前版本 \"%s\" 旧于 \"Typecho %s/%s\"",
			imp.SrcVersion, TypechoTestedVerMain, TypechoTestedVerSub))
		logrus.Warn("不确定导入是否能够成功导入，您可以选择升级 Typecho: http://docs.typecho.org/upgrade")
	}
}

// TypechoRoute 路由
type TypechoRoute struct {
	URL    string `json:"url"`
	Format string `json:"format"`
	Params string `json:"params"`
	Regx   string `json:"regx"`
	Action string `json:"action"`
	Widget string `json:"widget"`
}

// 获取 typecho_options 表里面的 name:routingTable option，解析 option.value。
// option.value 数据为 PHP 序列化后的结果。
// @see https://www.php.net/manual/en/function.serialize.php
// @see https://www.php.net/manual/en/function.unserialize.php
func (imp *_TypechoImporter) GetOptionRoutingTable() (map[string]TypechoRoute, error) {
	// 仅获取一次
	if imp.OptionRoutingTable != nil {
		return imp.OptionRoutingTable, nil
	}

	dataStr := ""
	for _, opt := range imp.Options {
		if opt.Name == "routingTable" {
			dataStr = opt.Value
			break
		}
	}

	if dataStr == "" {
		return map[string]TypechoRoute{}, errors.New("`routingTable` Not Found in Options")
	}

	// Unmarshal
	var data map[interface{}]interface{}
	err := phpserialize.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return map[string]TypechoRoute{}, err
	}

	// interface{} to struct
	routingTable := map[string]TypechoRoute{}
	for k, v := range data {
		var r TypechoRoute
		mapstructure.Decode(v, &r)
		routingTable[fmt.Sprintf("%v", k)] = r
	}

	imp.OptionRoutingTable = routingTable // 就不再反复获取了

	return routingTable, nil
}

// 获取一个 Route
func (imp *_TypechoImporter) GetOptionRoute(name string) (TypechoRoute, error) {
	routingTable, err := imp.GetOptionRoutingTable()
	if err != nil {
		return TypechoRoute{}, err
	}

	for n, route := range routingTable {
		if n == name {
			return route, nil
		}
	}

	return TypechoRoute{}, errors.New(`Route Name "` + name + `" Not Found`)
}

// 重写路径获取
func (imp *_TypechoImporter) RewriteRuleReady() {
	check := func(nameText string, routeName string, field *string, defaultVal string) {
		if *field == "" {
			// 从数据库获取
			route, err := imp.GetOptionRoute(routeName)
			if err != nil || route.URL == "" {
				if err != nil {
					logrus.Error(err)
				}

				*field = defaultVal
				logrus.Error("[重写规则] \"" + nameText + "\" 无法从数据库读取，将使用默认值 \"" + imp.RewritePost + "\"")
				return
			}

			readRule := route.URL

			// 将 [year:digital:4] 变为 [year]
			r := regexp.MustCompile(`\[(([a-zA-Z0-9]+):?(.*?))\]`)
			subMatches := r.FindAllStringSubmatch(readRule, -1)
			replaces := map[string]string{}
			for _, m := range subMatches {
				if len(m) < 1 {
					continue
				}

				// 0 => [year:digital:4]
				// 1 => year:digital:4
				// 2 => year
				// 3 => digital:4
				replaces[m[1]] = "[" + m[2] + "]"
			}
			readRule = imp.ReplaceAllBracketed(readRule, replaces) // replaces keys 无括号

			// 保存从数据库读取到的 rule
			*field = readRule
			logrus.Info("重写规则 \"" + nameText + "\" 读取成功")
		} else {
			logrus.Info("[重写规则] 自定义 \""+nameText+"\" 规则：", fmt.Sprintf("%#v", *field))
		}
	}

	check("文章页", "post", &imp.RewritePost, TypechoRewritePostDefault)
	check("独立页面", "page", &imp.RewritePage, TypechoRewritePageDefault)
}

// 获取 Content 的 Metas
func (imp *_TypechoImporter) GetContentMetas(cid int) []TypechoMeta {
	metaIds := []int{}
	for _, rela := range imp.Relationships {
		if rela.Cid == cid {
			metaIds = append(metaIds, rela.Mid)
		}
	}

	metas := []TypechoMeta{}
	for _, m := range imp.Metas {
		isNeed := false
		for _, id := range metaIds {
			if id == m.Mid {
				isNeed = true
				break
			}
		}

		if isNeed {
			metas = append(metas, m)
		}
	}

	return metas
}

// 获取 Content 的分类
func (imp *_TypechoImporter) GetContentCategory(cid int) string {
	metas := imp.GetContentMetas(cid)
	for _, m := range metas {
		if m.Type == "category" {
			return m.Slug
		}
	}

	return ""
}

// Typecho 评论数据表
type TypechoComment struct {
	Coid     int    `gorm:"column:coid"` // comment_id
	Cid      int    `gorm:"column:cid"`  // content_id
	Created  int    `gorm:"column:created"`
	Author   string `gorm:"column:author"`
	AuthorId int    `gorm:"column:authorId"`
	OwnerId  int    `gorm:"column:ownerId"`
	Mail     string `gorm:"column:mail"`
	Url      string `gorm:"column:url"`
	Ip       string `gorm:"column:ip"`
	Agent    string `gorm:"column:agent"`
	Text     string `gorm:"column:text"`
	Type     string `gorm:"column:type"`
	Status   string `gorm:"column:status"`
	Parent   int    `gorm:"column:parent"`
	Stars    int    `gorm:"column:stars"`
	Notify   int    `gorm:"column:notify"`
	Likes    int    `gorm:"column:likes"`
	Dislikes int    `gorm:"column:dislikes"`
}

// Typecho 内容数据表 (Type: post, page)
type TypechoContent struct {
	Cid          int    `gorm:"column:cid"`
	Title        string `gorm:"column:title"`
	Slug         string `gorm:"column:slug"`
	Created      int    `gorm:"column:created"`
	Modified     int    `gorm:"column:modified"`
	Text         string `gorm:"column:text"`
	Order        int    `gorm:"column:order"`
	AuthorId     int    `gorm:"column:authorId"`
	Template     string `gorm:"column:template"`
	Type         string `gorm:"column:type"`
	Status       string `gorm:"column:status"`
	Password     string `gorm:"column:password"`
	CommentsNum  int    `gorm:"column:commentsNum"`
	AllowComment string `gorm:"column:allowComment"`
	AllowPing    string `gorm:"column:allowPing"`
	AllowFeed    string `gorm:"column:allowFeed"`
	Parent       int    `gorm:"column:parent"`
	Views        int    `gorm:"column:views"`
	ViewsNum     int    `gorm:"column:viewsNum"`
	LikesNum     int    `gorm:"column:likesNum"`
	WordCount    int    `gorm:"column:wordCount"`
	Likes        int    `gorm:"column:likes"`
}

// Typecho 配置数据表
type TypechoOption struct {
	Name  string `gorm:"column:name"`
	User  int    `gorm:"column:user"`
	Value string `gorm:"column:value"`
}

// Typecho 关联表 (Content => Metas)
type TypechoRelationship struct {
	Cid int `gorm:"column:cid"` // content_id
	Mid int `gorm:"column:mid"` // meta_id
}

// Typecho 附加字段数据表
type TypechoMeta struct {
	Mid         int    `gorm:"column:mid"`
	Name        string `gorm:"column:name"`
	Slug        string `gorm:"column:slug"`
	Type        string `gorm:"column:type"`
	Description string `gorm:"column:description"`
}
