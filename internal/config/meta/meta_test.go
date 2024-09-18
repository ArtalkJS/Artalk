package meta_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/config/meta"
	"github.com/artalkjs/artalk/v2/test"
)

func Test_GetOptionsMetaData(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	result, err := meta.GetOptionsMetaData(config.Template("zh-CN"))
	if err != nil {
		t.Error(err)
	}

	if len(result) == 0 {
		t.Error("should get some metadata")
	}

	// output json with pretty format indent
	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))
}
