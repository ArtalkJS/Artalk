package importer

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ArtalkJS/ArtalkGo/config"
	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/araddon/dateparse"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Map = map[string]interface{}

type Importer struct {
	Name string
	Desc string
	Run  func(basic BasicParams, payload []string)
}

var Supports = []interface{}{
	TypechoImporter,
	WordPressImporter,
}

func GetSupportTypes() []string {
	types := []string{}
	for _, i := range Supports {
		r := reflect.ValueOf(i)
		f := reflect.Indirect(r).FieldByName("Name")
		types = append(types, f.String())
	}

	return types
}

func RunByName(dataType string, payload []string) {
	basic := GetBasicParamsFrom(payload)
	for _, i := range Supports {
		r := reflect.ValueOf(i)
		name := reflect.Indirect(r).FieldByName("Name").String()
		desc := reflect.Indirect(r).FieldByName("Desc").String()
		if strings.EqualFold(name, dataType) {
			fmt.Print("\n* * *\n\n")
			fmt.Print(" [数据导入] " + name + "\n\n")
			logrus.Info(desc)
			fmt.Print("\n* * *\n\n")

			t1 := time.Now()
			r.MethodByName("Run").Call([]reflect.Value{
				reflect.ValueOf(basic),
				reflect.ValueOf(payload),
			})
			elapsed := time.Since(t1)
			fmt.Print("\n")
			logrus.Info("导出执行结束，耗时: ", elapsed)
		}
	}
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

type BasicParams struct {
	TargetSiteName string `json:"target_site_name"`
	TargetSiteUrl  string `json:"target_site_url"`
}

func GetBasicParamsFrom(payload []string) BasicParams {
	basic := BasicParams{}
	GetParamsFrom(payload).To(map[string]*string{
		"target_site_name": &basic.TargetSiteName,
		"target_site_url":  &basic.TargetSiteUrl,
	})

	if basic.TargetSiteName == "" {
		logrus.Fatal("参数 `target_site_name:<站点名称>` 不能为空")
	}
	if basic.TargetSiteUrl == "" {
		logrus.Fatal("参数 `target_site_url:<站点根目录 URL>` 不能为空")
	}
	if !lib.ValidateURL(basic.TargetSiteUrl) {
		logrus.Fatal("参数 `target_site_url:<站点根目录 URL>` 必须为 URL 格式")
	}

	return basic
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
	return db
}

// 站点准备
func SiteReady(basic BasicParams) model.Site {
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
