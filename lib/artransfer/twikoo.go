package artransfer

import (
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
)

var TwikooImporter = &_TwikooImporter{
	ImporterInfo: ImporterInfo{
		Name: "twikoo",
		Desc: "从 Twikoo 导入数据",
		Note: "",
	},
}

type _TwikooImporter struct {
	ImporterInfo
	TargetSite model.Site
}

func (imp *_TwikooImporter) Run(basic *BasicParams, payload []string) {
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

	// 解析 Twikoo JSON
	var tComments []TwikooCommentFAS
	dErr := JsonDecodeFAS(jsonStr, &tComments)
	if dErr != nil {
		logFatal(dErr)
		return
	}

	// twikoo 转 ArtalkTransferParcel
	// @see https://github.com/imaegoo/twikoo/blob/c528c94105449c6b10c63bded6f813ceaee4bf74/src/vercel/api/index.js#L1155
	// rid 对应 _id @see https://github.com/imaegoo/twikoo/blob/c528c94105449c6b10c63bded6f813ceaee4bf74/src/vercel/api/index.js#L343
	tp := []model.Artran{}
	for _, tc := range tComments {
		tp = append(tp, model.Artran{
			ID:          tc.ID,
			Rid:         tc.Rid,
			Content:     tc.Comment,
			UA:          tc.UA,
			IP:          tc.IP,
			IsCollapsed: lib.ToString(false),
			IsPending:   lib.ToString(tc.IsSpam == "true"),
			CreatedAt:   tc.Created,
			UpdatedAt:   tc.Updated,
			Nick:        tc.Nick,
			Email:       tc.Mail,
			Link:        tc.Link,
			PageKey:     tc.Url,
			SiteName:    basic.TargetSiteName,
			SiteUrls:    basic.TargetSiteUrl,
		})
	}

	ImportArtrans(basic, tp)
}

// TwikooCommentFAS (Fields All String type)
type TwikooCommentFAS struct {
	ID      string `json:"_id"`
	Uid     string `json:"uid"`
	Nick    string `json:"nick"`
	Mail    string `json:"mail"`
	MailMd5 string `json:"mailMd5"`
	Link    string `json:"link"`
	UA      string `json:"ua"`
	IP      string `json:"ip"`
	Master  string `json:"master"`
	Url     string `json:"url"`
	Href    string `json:"href"`
	Comment string `json:"comment"`
	Pid     string `json:"pid"`
	Rid     string `json:"rid"`
	IsSpam  string `json:"isSpam"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}
