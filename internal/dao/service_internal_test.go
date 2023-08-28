package dao

import (
	"reflect"
	"testing"
)

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
