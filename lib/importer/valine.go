package importer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/sirupsen/logrus"
)

var ValineImporter = &_ValineImporter{
	ImporterInfo: ImporterInfo{
		Name: "valine",
		Desc: "从 Valine 导入数据",
		Note: "",
	},
}

type _ValineImporter struct {
	ImporterInfo
	TargetSite model.Site
}

func (imp *_ValineImporter) Run(basic *BasicParams, payload []string) {
	RequiredBasicTargetSite(basic)

	imp.TargetSite = SiteReady(basic)

	// 读取文件
	var jsonFile string
	GetParamsFrom(payload).To(map[string]*string{
		"json_file": &jsonFile,
	})

	if jsonFile == "" {
		logrus.Fatal("请附带参数 `json_file:<Valine 导出的 JSON 文件路径>`")
	}
	if _, err := os.Stat(jsonFile); errors.Is(err, os.ErrNotExist) {
		logrus.Fatal("文件不存在，请检查参数 `json_file` 传入路径是否正确")
	}

	buf, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		logrus.Fatal("json 文件打开失败：", err)
	}
	jsonStr := string(buf)

	// 解析 JSON
	comments, err := ParseValineCommentJSON(jsonStr)
	if err != nil {
		logrus.Fatal("json 解析失败：", err)
	}

	// 汇总
	fmt.Print("# 请过目：\n\n")

	// 第一条评论
	if len(comments) > 0 {
		fmt.Printf("[第一条评论]\n\n"+
			"    %#v\n\n", comments[0])
	}

	PrintTable([]table.Row{
		{"目标站点名", basic.TargetSiteName},
		{"目标站点 URL", basic.TargetSiteUrl},
		{"评论数量", len(comments)},
	})

	fmt.Print("\n")

	// 确认开始
	if !Confirm("确认开始导入吗？") {
		os.Exit(0)
	}

	// 准备导入评论
	fmt.Print("\n")

	ImportValineComments(basic, comments)
}

func ParseValineCommentJSON(jsonStr string) ([]ValineComment, error) {
	var list []ValineComment
	err := json.Unmarshal([]byte(jsonStr), &list)
	if err != nil {
		return []ValineComment{}, err
	}

	return list, nil
}

func ImportValineComments(basic *BasicParams, comments []ValineComment) {
	siteName := basic.TargetSiteName

	// 查找父评论
	idDict := map[string]int{}
	id := 1

	idChanges := map[uint]uint{}

	for _, c := range comments {
		idDict[c.ObjectId] = id
		id++
	}

	for _, c := range comments {
		// PageKey
		tSiteUrl := basic.TargetSiteUrl
		tSiteUrl = strings.TrimSuffix(tSiteUrl, "/") + "/"
		pageKey := tSiteUrl + strings.TrimPrefix(lib.GetUrlWithoutDomain(c.Url), "/")

		// 创建 user
		user := model.FindCreateUser(c.Nick, c.Mail)
		page := model.FindCreatePage(pageKey, "", siteName)

		if c.Link != "" {
			user.Link = c.Link
		}
		model.UpdateUser(&user)

		// 创建新 comment 实例
		nComment := model.Comment{
			Content: c.Comment,

			PageKey:  page.Key,
			SiteName: basic.TargetSiteName,

			UserID: user.ID,
			UA:     c.UA,
			IP:     c.IP,

			Rid: uint(idDict[c.Rid]),

			IsCollapsed: false,
			IsPending:   false,
		}

		// 日期恢复
		createdVal := fmt.Sprintf("%v", c.CreatedAt)
		updatedVal := fmt.Sprintf("%v", c.UpdatedAt)
		nComment.CreatedAt = ParseDate(createdVal)
		nComment.UpdatedAt = ParseDate(updatedVal)

		// 保存到数据库
		err := lib.DB.Create(&nComment).Error
		if err != nil {
			logrus.Error(fmt.Sprintf("评论源 ID:%s 保存失败", c.ObjectId))
			continue
		}

		idChanges[uint(idDict[c.ObjectId])] = nComment.ID
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

type ValineComment struct {
	ObjectId  string `json:"objectId"`
	Nick      string `json:"nick"`
	IP        string `json:"ip"`
	Mail      string `json:"mail"`
	MailMd5   string `json:"mailMd5"`
	IsSpam    string `json:"isSpam"`
	UA        string `json:"ua"`
	Link      string `json:"link"`
	Pid       string `json:"pid"`
	Rid       string `json:"rid"`
	Comment   string `json:"comment"`
	Url       string `json:"url"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
