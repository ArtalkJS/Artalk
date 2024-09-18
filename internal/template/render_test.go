package template_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/template"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestNewRenderer(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("DefaultTemplateRender", func(t *testing.T) {
		customTplFile := fmt.Sprintf("%s/%s", t.TempDir(), "test_tpl_file")
		customTplContent := "[{{site_name}}] You got a reply from @{{reply_nick}}: {{reply_content}}"
		_ = os.WriteFile(customTplFile, []byte(customTplContent), 0644)
		defer os.Remove(customTplFile)

		tests := []struct {
			name                  string
			renderType            template.RenderType
			defaultTemplateLoader func() template.TemplateLoader
			expectedDefaultResult string
		}{
			{
				name:                  "EmailDefaultTplFile",
				renderType:            template.TYPE_EMAIL,
				defaultTemplateLoader: func() template.TemplateLoader { return template.NewFileLoader("") },
				expectedDefaultResult: "",
			},
			{
				name:                  "NotifyDefaultTplFile",
				renderType:            template.TYPE_NOTIFY,
				defaultTemplateLoader: func() template.TemplateLoader { return template.NewFileLoader("") },
				expectedDefaultResult: "",
			},
			{
				name:                  "CustomEmailTplByFileLoader",
				renderType:            template.TYPE_EMAIL,
				defaultTemplateLoader: func() template.TemplateLoader { return template.NewFileLoader(customTplFile) },
				expectedDefaultResult: "[Site A] You got a reply from @admin: <p>Hello Artalk, 你好 Artalk!</p>",
			},
			{
				name:                  "CustomNotifyTplByFileLoader",
				renderType:            template.TYPE_NOTIFY,
				defaultTemplateLoader: func() template.TemplateLoader { return template.NewFileLoader(customTplFile) },
				expectedDefaultResult: "[Site A] You got a reply from @admin: Hello Artalk, 你好 Artalk!",
			},
		}

		for _, tt := range tests {
			t.Run(string(tt.renderType), func(t *testing.T) {
				renderer := template.NewRenderer(app.Dao(), tt.renderType, tt.defaultTemplateLoader())
				tNotify := app.Dao().FindNotify(1000, 1000)

				// Test render default tpl with template loader
				result := renderer.Render(&tNotify)

				assert.NotEmpty(t, result, "default tpl should not be empty")
				if tt.expectedDefaultResult != "" {
					assert.Equal(t, tt.expectedDefaultResult, strings.TrimSpace(result), "default tpl should be rendered")
				}
			})
		}
	})

	t.Run("CustomTemplateRender", func(t *testing.T) {
		tests := []struct {
			name                 string
			renderType           template.RenderType
			customTpl            string
			expectedCustomResult string
		}{
			{
				name:                 "CustomTplByString",
				renderType:           template.TYPE_EMAIL,
				customTpl:            "[{{site_name}}] You got a reply from @{{reply_nick}}",
				expectedCustomResult: "[Site A] You got a reply from @admin",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				renderer := template.NewRenderer(app.Dao(), tt.renderType, nil)
				tNotify := app.Dao().FindNotify(1000, 1000)

				// Test custom tpl
				result := renderer.Render(&tNotify, tt.customTpl)
				assert.Equal(t, tt.expectedCustomResult, result, "custom tpl should be rendered")
			})
		}
	})
}
