package importer

import (
	"encoding/json"

	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/sirupsen/logrus"
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
	RequiredBasicTargetSite(basic)

	// 读取文件
	jsonStr := JsonFileReady(payload)

	// 解析 JSON
	tComments, err := ParseTwikooCommentJSON(jsonStr)
	if err != nil {
		logrus.Fatal("json 解析失败：", err)
	}

	// twikoo 转 valine
	// @see https://github.com/imaegoo/twikoo/blob/c528c94105449c6b10c63bded6f813ceaee4bf74/src/vercel/api/index.js#L1155
	// rid 对应 _id @see https://github.com/imaegoo/twikoo/blob/c528c94105449c6b10c63bded6f813ceaee4bf74/src/vercel/api/index.js#L343
	vComments := []ValineComment{}
	for _, tc := range tComments {
		vComments = append(vComments, ValineComment{
			ObjectId:  tc.ID,
			Nick:      tc.Nick,
			IP:        tc.IP,
			Mail:      tc.Mail,
			MailMd5:   tc.MailMd5,
			IsSpam:    tc.IsSpam,
			UA:        tc.UA,
			Link:      tc.Link,
			Pid:       tc.Pid,
			Rid:       tc.Rid,
			Comment:   tc.Comment,
			Url:       tc.Url,
			CreatedAt: tc.Created,
			UpdatedAt: tc.Updated,
		})
	}

	ImportValine(basic, payload, vComments)
}

func ParseTwikooCommentJSON(jsonStr string) ([]TwikooComment, error) {
	var list []TwikooComment
	err := json.Unmarshal([]byte(jsonStr), &list)
	if err != nil {
		return []TwikooComment{}, err
	}
	return list, nil
}

type TwikooComment struct {
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
