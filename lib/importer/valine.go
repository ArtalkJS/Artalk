package importer

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
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
}

func (imp *_ValineImporter) Run(basic *BasicParams, payload []string) {
	RequiredBasicTargetSite(basic)

	// 读取文件
	jsonStr := JsonFileReady(payload)

	// 解析 Valine JSON
	var vComments []ValineCommentFAS
	JsonDecodeFAS(jsonStr, &vComments)

	// Valine 数据格式转 Artrans
	tp := []model.Artran{}
	for _, vc := range vComments {
		tp = append(tp, model.Artran{
			ID:          vc.ObjectId,
			Rid:         vc.Rid,
			Content:     vc.Comment,
			UA:          vc.UA,
			IP:          vc.IP,
			IsCollapsed: lib.ToString(false),
			IsPending:   lib.ToString(vc.IsSpam == "true"),
			CreatedAt:   vc.CreatedAt,
			UpdatedAt:   vc.UpdatedAt,
			Nick:        vc.Nick,
			Email:       vc.Mail,
			Link:        vc.Link,
			PageKey:     vc.Url,
			SiteName:    basic.TargetSiteName,
			SiteUrls:    basic.TargetSiteUrl,
		})
	}

	ImportArtrans(basic, tp)
}

// ValineCommentFAS (FieldAllStr)
type ValineCommentFAS struct {
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
