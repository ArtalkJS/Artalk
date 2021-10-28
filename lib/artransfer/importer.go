package artransfer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/araddon/dateparse"
	"github.com/jedib0t/go-pretty/v6/table"
	"gorm.io/gorm"
)

type Map = map[string]interface{}

var Supports = []interface{}{
	ArtransImporter,
	TypechoImporter,
	ValineImporter,
	TwikooImporter,
	ArtalkV1Importer,
}

func RunByName(dataType string, payload []string) {
	basic := GetBasicParamsFrom(payload)
	for _, i := range Supports {
		r := reflect.ValueOf(i)
		name := reflect.Indirect(r).FieldByName("Name").String()
		desc := reflect.Indirect(r).FieldByName("Desc").String()
		note := reflect.Indirect(r).FieldByName("Note").String()
		if !strings.EqualFold(name, dataType) {
			continue
		}

		print("\n")
		tableData := []table.Row{
			{"数据搬家 - 导入"},
			{strings.ToUpper(name)},
			{desc},
		}
		if note != "" {
			tableData = append(tableData, table.Row{note})
		}
		PrintTable(tableData)
		print("\n")

		//t1 := time.Now()
		r.MethodByName("Run").Call([]reflect.Value{
			reflect.ValueOf(basic),
			reflect.ValueOf(payload),
		})
		//elapsed := time.Since(t1)
		print("\n")
		logInfo("导入执行结束") //，耗时: ", elapsed)
		return
	}

	logFatal("不支持该数据类型导入")
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

func GetSupportNames() []string {
	types := []string{}
	for _, i := range Supports {
		r := reflect.ValueOf(i)
		f := reflect.Indirect(r).FieldByName("Name")
		types = append(types, f.String())
	}

	return types
}

type BasicParams struct {
	TargetSiteName string
	TargetSiteUrl  string

	UrlResolver bool
}

func GetBasicParamsFrom(payload []string) *BasicParams {
	basic := BasicParams{}

	basic.UrlResolver = true // 默认开启

	GetParamsFrom(payload).To(map[string]interface{}{
		"t_name":         &basic.TargetSiteName,
		"t_url":          &basic.TargetSiteUrl,
		"t_url_resolver": &basic.UrlResolver,
	})

	if !basic.UrlResolver {
		logWarn("目标站点 URL 解析器已关闭")
	}

	return &basic
}

func RequiredBasicTargetSite(basic *BasicParams) error {
	if basic.TargetSiteName == "" {
		return errors.New("请附带参数 `t_name:<目标站点名称>`")
	}
	if basic.TargetSiteUrl == "" {
		return errors.New("请附带参数 `t_url:<目标站点根目录 URL>`")
	}
	if !lib.ValidateURL(basic.TargetSiteUrl) {
		return errors.New("参数 `t_url:<目标站点根目录 URL>` 必须为 URL 格式")
	}

	return nil
}

func DbReady(payload []string) (*gorm.DB, error) {
	var host, port, dbName, user, password, dbType, dbFile string
	GetParamsFrom(payload).To(map[string]interface{}{
		"db_host":     &host,
		"db_port":     &port,
		"db_name":     &dbName,
		"db_user":     &user,
		"db_password": &password,
		"db_type":     &dbType,
		"db_file":     &dbFile,
	})

	if dbType == "" {
		dbType = string(config.TypeMySql)
	}

	// sqlite
	if strings.EqualFold(dbType, string(config.TypeSQLite)) {
		if dbFile == "" {
			return nil, errors.New("SQLite 数据库：请传递参数 `db_file` 指定数据文件路径")
		}

		db, err := lib.OpenSQLite(dbFile)
		if err != nil {
			return nil, errors.New("SQLite 打开失败 " + err.Error())
		}
		return db, nil
	}

	dsn, err := lib.GetDsn(config.DBType(dbType), host, port, dbName, user, password)
	if err != nil {
		return nil, errors.New("数据库连接 DSN 生成错误 " + err.Error())
	}

	db, err := lib.OpenDB(config.DBType(dbType), dsn)
	if err != nil {
		return nil, errors.New("数据库连接失败 " + err.Error())
	}

	logInfo("数据库连接成功")

	return db, nil
}

// 站点准备
func SiteReady(tSiteName string, tSiteUrls string) (model.Site, error) {
	site := model.FindSite(tSiteName)
	if site.IsEmpty() {
		// 创建新站点
		site = model.Site{}
		site.Name = tSiteName
		site.Urls = tSiteUrls
		err := lib.DB.Create(&site).Error
		if err != nil {
			return model.Site{}, errors.New("站点创建失败")
		}
	} else {
		// 追加 URL
		siteCooked := site.ToCooked()

		urlExist := func(tUrl string) bool {
			for _, u := range siteCooked.Urls {
				if u == tUrl {
					return true
				}
			}
			return false
		}

		tUrlsSpit := strings.Split(tSiteUrls, ",")

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
			err := lib.DB.Save(&site).Error
			if err != nil {
				return model.Site{}, errors.New("站点数据更新失败")
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
		return "", errors.New("请附带参数 `json_file:<JSON 数据文件路径>`")
	}
	if _, err := os.Stat(jsonFile); errors.Is(err, os.ErrNotExist) {
		return "", errors.New("文件不存在，请检查参数 `json_file` 传入路径是否正确")
	}

	buf, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return "", errors.New("json 文件打开失败：" + err.Error())
	}

	return string(buf), nil
}

// PageKey (commentUrlVal 不确定是否为完整 URL 还是一个 path)
func UrlResolverGetPageKey(baseUrl string, commentUrlVal string) string {
	baseUrl = strings.TrimSuffix(baseUrl, "/") + "/"
	path := strings.TrimPrefix(lib.GetUrlWithoutDomain(commentUrlVal), "/")

	// 解决拼接路径中的相对地址，例如：https://atk.xxx/abc/../artalk => https://atk.xxx/artalk
	u, err := url.ParseRequestURI(baseUrl + path)
	if err != nil {
		logError("GetNewPageKey Error: ", err)
		return commentUrlVal
	}

	// pathIsDir := strings.HasSuffix(u.Path, "/") // path 以 / 结尾
	abs, absErr := filepath.Abs(u.Path) // 相对路径转绝对路径
	if absErr != nil {
		return u.String()
	}

	// TODO 会导致 "http://aaa.com", "/" => "http://aaa.com//" 暂时搁置
	// if pathIsDir {
	// 	u.Path = abs + "/" // 加上 "/"
	// 	// 这是一个 patch: 因为 filepath.Abs() 结果无论是 目录还是文件，都会去掉 /
	// } else {
	// 	u.Path = abs
	// }

	u.Path = abs

	return u.String()
}

func ParseVersion(s string) int64 {
	strList := strings.Split(s, ".")
	format := fmt.Sprintf("%%s%%0%ds", len(strList))
	v := ""
	for _, value := range strList {
		v = fmt.Sprintf(format, v, value)
	}
	var result int64
	var err error
	if result, err = strconv.ParseInt(v, 10, 64); err != nil {
		logError("ugh: parseVersion(%s): error=%s", s, err)
		return 0
	}
	return result
}

func ParseDate(s string) time.Time {
	denverLoc, _ := time.LoadLocation(config.Instance.TimeZone) // 时区
	time.Local = denverLoc
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

// Json Decode (FAS: Fields All String Type)
// 解析 json 为字段全部是 string 类型的 struct
func JsonDecodeFAS(str string, fasStructure interface{}) error {
	err := json.Unmarshal([]byte(lib.JsonObjInArrAnyStr(str)), fasStructure) // lib.ToString()
	if err != nil {
		return errors.New("JSON 解析失败 " + err.Error())
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
		nComment := model.FindComment(newId)
		if nComment.Rid == 0 {
			continue
		}
		if newId, isExist := idChanges[nComment.Rid]; isExist {
			nComment.Rid = newId
			err := lib.DB.Save(&nComment).Error
			if err != nil {
				logError(fmt.Sprintf("[rid 更新] new_id:%d new_rid:%d", nComment.ID, newId), err)
			}
		}
	}

	print("\n")
	logInfo("RID 重构完毕")
}
