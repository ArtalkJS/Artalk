package artransfer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/i18n"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/ArtalkJS/Artalk/internal/utils"
	"github.com/araddon/dateparse"
)

func RunImportArtrans(payload []string) {
	basic := GetBasicParamsFrom(payload)

	name := ArtransImporter.ImporterInfo.Name
	desc := ArtransImporter.ImporterInfo.Desc
	note := ArtransImporter.ImporterInfo.Note

	print("\n")
	tableData := [][]interface{}{
		{"Artransfer - Import"},
		{strings.ToUpper(name)},
		{desc},
	}
	if note != "" {
		tableData = append(tableData, []interface{}{note})
	}
	PrintTable(tableData)
	print("\n")

	//t1 := time.Now()
	ArtransImporter.Run(basic, payload)
	//elapsed := time.Since(t1)

	print("\n")
	logInfo(i18n.T("Import complete")) //，耗时: ", elapsed)
}

type ImporterInfo struct {
	Name string
	Desc string
	Note string
}

func GetImporterInfo(instance interface{}) ImporterInfo {
	var info ImporterInfo
	j, _ := json.Marshal(instance)
	json.Unmarshal(j, &info)
	return info
}

type BasicParams struct {
	TargetSiteName string
	TargetSiteUrl  string

	UrlResolver bool
}

func GetBasicParamsFrom(payload []string) *BasicParams {
	basic := BasicParams{}

	basic.UrlResolver = false // 默认关闭

	GetParamsFrom(payload).To(map[string]interface{}{
		"t_name":         &basic.TargetSiteName,
		"t_url":          &basic.TargetSiteUrl,
		"t_url_resolver": &basic.UrlResolver,
	})

	if !basic.UrlResolver {
		logWarn("Target site URL resolver disabled")
	}

	return &basic
}

func RequiredBasicTargetSite(basic *BasicParams) error {
	if basic.TargetSiteName == "" {
		return errors.New(i18n.T("{{name}} is required", map[string]interface{}{"name": "t_name:<Target Site Name>"}))
	}
	if basic.TargetSiteUrl == "" {
		return errors.New(i18n.T("{{name}} is required", map[string]interface{}{"name": "t_url:<Target Site Root URL>"}))
	}
	if !utils.ValidateURL(basic.TargetSiteUrl) {
		return errors.New("invalid URL for parameter `t_url:<Target Site Root URL>`")
	}

	return nil
}

// 站点准备
func SiteReady(tSiteName string, tSiteUrls string) (entity.Site, error) {
	site := query.FindSite(tSiteName)
	if site.IsEmpty() {
		// 创建新站点
		site = entity.Site{}
		site.Name = tSiteName
		site.Urls = tSiteUrls
		err := query.CreateSite(&site)
		if err != nil {
			return entity.Site{}, errors.New("failed to create site")
		}
	} else {
		// 追加 URL
		siteCooked := query.CookSite(&site)

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
			err := query.UpdateSite(&site)
			if err != nil {
				return entity.Site{}, errors.New("update site data failed")
			}
		}
	}

	return site, nil
}

func JsonFileReady(payload []string) (string, error) {
	var jsonFile, jsonData string
	GetParamsFrom(payload).To(map[string]interface{}{
		"json_file": &jsonFile,
		"json_data": &jsonData,
	})

	// 直接给 JSON 内容，不去读取文件
	if jsonData != "" {
		return jsonData, nil
	}

	if jsonFile == "" {

		return "", errors.New(i18n.T("{{name}} is required", map[string]interface{}{"name": "json_file:<JSON file path>"}))
	}
	if _, err := os.Stat(jsonFile); errors.Is(err, os.ErrNotExist) {
		return "", errors.New(i18n.T("{{name}} not found", map[string]interface{}{"name": i18n.T("File")}))
	}

	buf, err := os.ReadFile(jsonFile)
	if err != nil {
		return "", errors.New("file open failed" + ": " + err.Error())
	}

	return string(buf), nil
}

// PageKey (commentUrlVal 不确定是否为完整 URL 还是一个 path)
//
// @examples
// ("https://github.com", "/1.html")                => "https://github.com/1.html"
// ("https://github.com", "https://xxx.com/1.html") => "https://github.com/1.html"
// ("https://github.com/", "/1.html")               => "https://github.com/1.html"
// ("", "/1.html")                                  => "/1.html"
// ("", "https://xxx.com/1.html")                   => "https://xxx.com/1.html"
// ("https://github.com/233", "/1/")                => "https://github.com/1/"
func UrlResolverGetPageKey(baseUrlRaw string, commentUrlRaw string) string {
	if baseUrlRaw == "" {
		return commentUrlRaw
	}

	baseUrl, err := url.Parse(baseUrlRaw)
	if err != nil {
		return commentUrlRaw
	}

	commentUrl, err := url.Parse(commentUrlRaw)
	if err != nil {
		return commentUrlRaw
	}

	// "https://artalk.js.org/guide/describe.html?233" => "/guide/describe.html?233"
	commentUrl.Scheme = ""
	commentUrl.Host = ""

	// 解决拼接路径中的相对地址，例如：https://atk.xxx/abc/../artalk => https://atk.xxx/artalk
	url := baseUrl.ResolveReference(commentUrl)

	return url.String()
}

func ParseDate(s string) time.Time {
	denverLoc, _ := time.LoadLocation(config.Instance.TimeZone) // 时区
	time.Local = denverLoc
	// TODO should be restricted to using only the RFC3339 standard time format
	t, _ := dateparse.ParseIn(s, denverLoc)

	return t
}

type _getParamsTo struct {
	To func(variables map[string]interface{})
}

func GetParamsFrom(payload []string) _getParamsTo {
	a := _getParamsTo{}
	a.To = func(variables map[string]interface{}) {
		for _, pVal := range payload {
			for fromName, toVar := range variables {
				if !strings.HasPrefix(pVal, fromName+":") {
					continue
				}

				valStr := strings.TrimPrefix(pVal, fromName+":")

				switch reflect.ValueOf(toVar).Interface().(type) {
				case *string:
					*toVar.(*string) = valStr
				case *bool:
					*toVar.(*bool) = strings.EqualFold(valStr, "true")
				case *int:
					num, err := strconv.Atoi(valStr)
					if err != nil {
						*toVar.(*int) = num
					}
				}
				break
			}
		}
	}
	return a
}

func GetArrayParamsFrom(payload []string, key string) []string {
	arr := []string{}
	for _, pVal := range payload {
		if strings.HasPrefix(pVal, key+":") {
			arr = append(arr, strings.TrimPrefix(pVal, key+":"))
		}
	}

	return arr
}

func CheckIfJsonArr(str string) bool {
	x := bytes.TrimSpace([]byte(str))
	return len(x) > 0 && x[0] == '['
}

func CheckIfJsonObj(str string) bool {
	x := bytes.TrimSpace([]byte(str))
	return len(x) > 0 && x[0] == '{'
}

func TryConvertLineJsonToArr(str string) (string, error) {
	// 尝试将一行一行的 Obj 转成 Arr
	arrTmp := []map[string]interface{}{}
	for _, line := range strings.Split(strings.TrimSpace(str), "\n") {
		var tmp map[string]interface{}
		err := json.Unmarshal([]byte(line), &tmp)
		if err != nil {
			return "", err
		}
		arrTmp = append(arrTmp, tmp)
	}
	r, err := json.Marshal(arrTmp)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

// Json Decode (FAS: Fields All String Type)
// 解析 json 为字段全部是 string 类型的 struct
func JsonDecodeFAS(str string, fasStructure interface{}) error {
	if !CheckIfJsonArr(str) {
		var err error
		str, err = TryConvertLineJsonToArr(str)
		if err != nil {
			return errors.New("JSON of array type is required: " + err.Error())
		}
	}

	err := json.Unmarshal([]byte(utils.JsonObjInArrAnyStr(str)), fasStructure) // lib.ToString()
	if err != nil {
		return errors.New("failed to parse JSON: " + err.Error())
	}

	return nil
}

func HideJsonLongText(key string, text string) string {
	r := regexp.MustCompile(key + `:"(.+?)"`)
	sm := r.FindStringSubmatch(text)
	postText := ""
	if len(sm) > 0 {
		postText = sm[1]
	}

	text = r.ReplaceAllString(text, fmt.Sprintf(key+": <!-- 省略 %d 个字符 -->", utf8.RuneCountInString(postText)))
	return text
}

// 传入 ID 变更表 (原始ID => 数据库已存在记录的ID) rid 将根据此替换
func RebuildRid(idChanges map[uint]uint) {
	for _, newId := range idChanges {
		nComment := query.FindComment(newId)
		if nComment.Rid == 0 {
			continue
		}
		if newId, isExist := idChanges[nComment.Rid]; isExist {
			nComment.Rid = newId
			err := query.UpdateComment(&nComment)
			if err != nil {
				logError(fmt.Sprintf("[rid 更新] new_id:%d new_rid:%d", nComment.ID, newId), err)
			}
		}
	}

	print("\n")
	logInfo("RID 重构完毕")
}
