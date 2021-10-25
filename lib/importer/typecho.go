package importer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/sirupsen/logrus"
)

const (
	// 已测试可用版本
	TypechoTestedVerMain = "1.1"
	TypechoTestedVerSub  = "17.10.30"
)

var TypechoImporter = &_TypechoImporter{
	Importer: Importer{
		Name: "Typecho",
		Desc: fmt.Sprintf("测试兼容：Typecho %s/%s", TypechoTestedVerMain, TypechoTestedVerSub),
	},
}

type _TypechoImporter struct {
	Importer
	TargetSite model.Site

	Basic         BasicParams
	SrcVersion    string
	SrcSiteName   string
	SrcSiteURL    string
	Comments      []TypechoComment
	Options       []TypechoOption
	Contents      []TypechoContent
	Metas         []TypechoMeta
	Relationships []TypechoRelationship

	DbPrefix string

	RewriteRulePost string
	RewriteRulePage string
}

// Typecho 升级相关的代码 @see https://github.com/typecho/typecho/blob/64b8e686885d8ab4c7f0cdc3d6dc2d99fa48537c/var/Utils/Upgrade.php
// 路由 @see https://github.com/typecho/typecho/blob/530312443142577509df88ce88cf3274fac9b8c4/var/Widget/Options/Permalink.php#L319
func (imp *_TypechoImporter) Run(basic BasicParams, payload []string) {
	// Ready
	imp.Basic = basic
	typechoDB := DbReady(payload)
	imp.TargetSite = SiteReady(basic)

	GetParamsFrom(payload).To(map[string]*string{
		"prefix":            &imp.DbPrefix,
		"rewrite_rule_post": &imp.RewriteRulePost,
		"rewrite_rule_page": &imp.RewriteRulePage,
	})

	if imp.DbPrefix == "" {
		imp.DbPrefix = "typecho_"
	}

	// Load Options
	typechoDB.Raw(fmt.Sprintf("SELECT * FROM %soptions", imp.DbPrefix)).Scan(&imp.Options)

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

	fmt.Printf("- 源版本：%#v\n", imp.SrcVersion)
	fmt.Printf("- 源站点名：%#v\n", imp.SrcSiteName)
	fmt.Printf("- 源站点URL：%#v\n", imp.SrcSiteURL)

	imp.CheckVersion()

	// load Metas
	typechoDB.Raw(fmt.Sprintf("SELECT * FROM %smetas", imp.DbPrefix)).Scan(&imp.Metas)

	// load Relationships
	typechoDB.Raw(fmt.Sprintf("SELECT * FROM %srelationships", imp.DbPrefix)).Scan(&imp.Relationships)

	// 重写规则
	imp.CheckRewriteRule()

	// 开始导入评论
	fmt.Println()

	var contents []TypechoContent
	typechoDB.Raw(fmt.Sprintf("SELECT * FROM %scontents ORDER BY created ASC", imp.DbPrefix)).Scan(&contents)
	if len(contents) > 0 {
		fmt.Printf("[第一篇文章]\n\n"+
			"    %#v\n\n", contents[0])
	}
	imp.Contents = contents

	var comments []TypechoComment
	typechoDB.Raw(fmt.Sprintf("SELECT * FROM %scomments ORDER BY created ASC", imp.DbPrefix)).Scan(&comments)
	if len(contents) > 0 {
		fmt.Printf("[第一条评论]\n\n"+
			"    %#v\n\n", comments[0])
	}
	imp.Comments = comments

	// Import Data
	imp.ImportComments()
}

// 导入评论
func (imp *_TypechoImporter) ImportComments() {
	articles := imp.Contents
	comments := imp.Comments

	fmt.Print("====================================\n\n")
	logrus.Info(fmt.Sprintf("[开始导入] 共 %d 个页面，%d 条评论", len(articles), len(comments)))
	fmt.Print("\n")

	siteName := imp.Basic.TargetSiteName

	idChanges := map[uint]uint{} // comment_id: old => new

	// 遍历导入文章的评论
	for _, art := range articles {
		// 创建页面
		pageKey := imp.GetNewPageKey(art)
		page := model.FindCreatePage(pageKey, art.Title, imp.Basic.TargetSiteName)

		// 查询评论
		commentTotal := 0
		for _, co := range comments {
			if co.Cid != art.Cid {
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

		fmt.Printf("+ [%-3d] 条评论 <- [%5d] %-30s | %#v\n", commentTotal, art.Cid, art.Slug, art.Title)
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

func (imp *_TypechoImporter) GetNewPageKey(content TypechoContent) string {
	// TODO: page/post 重写规则

	baseUrl := imp.Basic.TargetSiteUrl
	baseUrl = strings.TrimSuffix(baseUrl, "/") + "/"
	return baseUrl + content.Slug + ".html"
}

// 版本国旧检测
func (imp *_TypechoImporter) CheckVersion() {
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

func (imp *_TypechoImporter) GetOptionRoutingTable() string {
	for _, opt := range imp.Options {
		if opt.Name == "routingTable" {
			return opt.Value
		}
	}

	return ""
}

// 重写路径获取
func (imp *_TypechoImporter) CheckRewriteRule() {
	if imp.RewriteRulePost == "" {

	} else {
		logrus.Info("[重写规则] 自定义 “文章” 规则：", fmt.Sprintf("%#v", imp.RewriteRulePost))
	}

	if imp.RewriteRulePage == "" {

	} else {
		logrus.Info("[重写规则] 自定义 “独立页面” 规则：", fmt.Sprintf("%#v", imp.RewriteRulePage))
	}
}

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

type TypechoOption struct {
	Name  string `gorm:"column:name"`
	User  int    `gorm:"column:user"`
	Value string `gorm:"column:value"`
}

type TypechoRelationship struct {
	Cid int `gorm:"column:cid"` // content_id
	Mid int `gorm:"column:mid"` // meta_id
}

type TypechoMeta struct {
	Mid         int    `gorm:"column:mid"`
	Name        string `gorm:"column:name"`
	Slug        string `gorm:"column:slug"`
	Type        string `gorm:"column:type"`
	Description string `gorm:"column:description"`
}
