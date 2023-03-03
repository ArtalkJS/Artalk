package query

import (
	"reflect"
	"testing"

	"github.com/ArtalkJS/Artalk/internal/entity"
)

func Test_GetPageAccessibleURL(t *testing.T) {
	tests := []struct {
		name     string
		pageKey  string
		siteUrls string
		want     string
	}{
		{name: "相对路径", pageKey: "/abcd.html", siteUrls: "https://qwqaq.com", want: "https://qwqaq.com/abcd.html"},
		{name: "绝对路径", pageKey: "https://xxx.com/abcd.html", siteUrls: "https://qwqaq.com", want: "https://xxx.com/abcd.html"},
		{name: "未设置站点 URL", pageKey: "/abcd.html", siteUrls: "", want: "/abcd.html"},
		{name: "使用第一个站点 URL", pageKey: "/abcd.html", siteUrls: "https://first_url.com,https://second_url.com", want: "https://first_url.com/abcd.html"},
		{name: "复杂路径解析", pageKey: "/sub-folder/xxx?abc=test#test", siteUrls: "https://qwqaq.com", want: "https://qwqaq.com/sub-folder/xxx?abc=test#test"},
		{name: "域名子目录", pageKey: "test/1.html", siteUrls: "https://qwqaq.com/sub-folder/", want: "https://qwqaq.com/sub-folder/test/1.html"},
		{name: "相对路径，上一级目录", pageKey: "../test/1.html", siteUrls: "https://qwqaq.com/sub-folder/abc/", want: "https://qwqaq.com/sub-folder/test/1.html"},
		{name: "相对路径，末尾带斜杠", pageKey: "/slash-test/", siteUrls: "https://qwqaq.com", want: "https://qwqaq.com/slash-test/"},
		{name: "相对路径，末尾无斜杠", pageKey: "/slash-test", siteUrls: "https://qwqaq.com", want: "https://qwqaq.com/slash-test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := entity.Site{}
			s.ID = uint(1010)
			s.Urls = tt.siteUrls
			p := &entity.Page{
				Key:  tt.pageKey,
				Site: s,
			}

			if got := GetPageAccessibleURL(p); got != tt.want {
				t.Errorf("GetPageAccessibleURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pageExtractFromHTML(t *testing.T) {
	html := []byte(`
<!DOCTYPE html>
<html lang="zh-cn">
  <head>
	<Title> 2333ABC  </Title>
    <link rel="canonical" href="https://qwqaq.com/red/">
    <meta name="robots" content="noindex">
    <meta charset="utf-8">
    <meta http-equiv="refresh" content="0; url=https://qwqaq.com/red/"><"">>
  </head>
</html>
`)

	type args struct {
		html []byte
	}
	tests := []struct {
		name     string
		args     args
		wantData pageExtractData
	}{
		{name: "test title and redirect url extraction", args: args{html: html}, wantData: pageExtractData{
			Title:       "2333ABC",
			RedirectURL: "https://qwqaq.com/red/",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotData := pageExtractFromHTML(tt.args.html); !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("pageExtractFromHTML() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}
