package importer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Map = map[string]interface{}

var Supports = []interface{}{
	TypechoImporter,
	WordPressImporter,
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

		fmt.Print("\n")
		PrintTable([]table.Row{
			{"数据搬家 - 导入"},
			{strings.ToUpper(name)},
			{desc},
			{note},
		})
		fmt.Print("\n")

		//t1 := time.Now()
		r.MethodByName("Run").Call([]reflect.Value{
			reflect.ValueOf(basic),
			reflect.ValueOf(payload),
		})
		//elapsed := time.Since(t1)
		fmt.Print("\n")
		logrus.Info("导入执行结束") //，耗时: ", elapsed)
		return
	}

	logrus.Fatal("不支持该数据类型导入")
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

type _getParamsTo struct {
	To func(variables map[string]*string)
}

func GetParamsFrom(payload []string) _getParamsTo {
	a := _getParamsTo{}
	a.To = func(variables map[string]*string) {
		for _, pVal := range payload {
			for fromName, toVar := range variables {
				if !strings.HasPrefix(pVal, fromName+":") {
					continue
				}

				*toVar = strings.TrimPrefix(pVal, fromName+":")
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

type BasicParams struct {
	TargetSiteName string `json:"target_site_name"`
	TargetSiteUrl  string `json:"target_site_url"`
}

func GetBasicParamsFrom(payload []string) *BasicParams {
	basic := BasicParams{}
	GetParamsFrom(payload).To(map[string]*string{
		"target_site_name": &basic.TargetSiteName,
		"target_site_url":  &basic.TargetSiteUrl,
	})

	return &basic
}

func RequiredBasicTargetSite(basic *BasicParams) {
	if basic.TargetSiteName == "" {
		logrus.Fatal("请附带参数 `target_site_name:<站点名称>`")
	}
	if basic.TargetSiteUrl == "" {
		logrus.Fatal("请附带参数 `target_site_url:<站点根目录 URL>`")
	}
	if !lib.ValidateURL(basic.TargetSiteUrl) {
		logrus.Fatal("参数 `target_site_url:<站点根目录 URL>` 必须为 URL 格式")
	}
}

func DbReady(payload []string) *gorm.DB {
	var host, port, dbName, user, password, dbType, dbFile string
	GetParamsFrom(payload).To(map[string]*string{
		"host":     &host,
		"port":     &port,
		"db_name":  &dbName,
		"user":     &user,
		"password": &password,
		"db_type":  &dbType,
		"db_file":  &dbFile,
	})

	if dbType == "" {
		dbType = string(config.TypeMySql)
	}

	// sqlite
	if strings.EqualFold(dbType, string(config.TypeSQLite)) {
		if dbFile == "" {
			logrus.Fatal("SQLite 数据库：请传递参数 `db_file` 指定数据文件路径")
		}

		db, err := lib.OpenSQLite(dbFile)
		if err != nil {
			logrus.Fatal("SQLite 打开失败 ", err)
		}
		return db
	}

	dsn, err := lib.GetDsn(config.DBType(dbType), host, port, dbName, user, password)
	if err != nil {
		logrus.Fatal("数据库连接 DSN 生成错误 ", err)
	}

	db, err := lib.OpenDB(config.DBType(dbType), dsn)
	if err != nil {
		logrus.Fatal("数据库连接失败 ", err)
	}

	logrus.Info("数据库连接成功")

	return db
}

// 站点准备
func SiteReady(basic *BasicParams) model.Site {
	site := model.FindSite(basic.TargetSiteName)
	if site.IsEmpty() {
		// 创建新站点
		site = model.Site{}
		site.Name = basic.TargetSiteName
		site.Urls = basic.TargetSiteUrl
		err := lib.DB.Create(&site).Error
		if err != nil {
			logrus.Fatal("站点创建失败")
		}
	} else {
		sic := site.ToCooked()

		// 加 URL
		urlExist := false
		for _, url := range sic.Urls {
			if url == site.Urls {
				urlExist = true
				break
			}
		}

		if !urlExist {
			urls := []string{}
			urls = append(urls, basic.TargetSiteUrl) // prepend
			urls = append(urls, sic.Urls...)
			site.Urls = strings.Join(urls, ",")
			err := lib.DB.Save(&site).Error
			if err != nil {
				logrus.Fatal("站点数据更新失败")
			}
		}
	}

	return site
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
		logrus.Error("ugh: parseVersion(%s): error=%s", s, err)
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

func PrintTable(rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	for _, r := range rows {
		t.AppendRow(r)
	}

	tStyle := table.StyleLight
	tStyle.Options.SeparateRows = true
	t.SetStyle(tStyle)

	t.Render()
}

func Confirm(s string) bool {
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}

		resp := strings.ToLower(strings.TrimSpace(res))
		if resp == "y" || resp == "yes" {
			return true
		} else if resp == "n" || resp == "no" {
			return false
		}
	}
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
