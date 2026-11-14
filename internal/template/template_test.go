package template

import (
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCommentTemplateFields(t *testing.T) {
	params := tplParams{
		Comment: entity.EmailCommentDTO{
			CommentDTO: entity.CommentDTO{ID: 42, Content: "original"},
			Content:    "<p>Rendered</p>",
			ContentRaw: "**original**",
			Nick:       "Alice & Bob",
			Page:       entity.PageDTO{Title: "Page title"},
			Site:       entity.SiteDTO{Name: "Site name"},
		},
	}

	result := renderCommon(
		"{{comment.id}}|{{comment.content}}|{{comment.content_raw}}|{{comment.nick}}|{{comment.page.title}}|{{comment.site.name}}",
		params,
	)
	assert.Equal(t, "42|<p>Rendered</p>|**original**|Alice &amp; Bob|Page title|Site name", result)
}

func TestHandleEmoticonsImgTagsForNotify(t *testing.T) {
	tests := []struct {
		name string
		give string
		want string
	}{
		{name: "Test1", give: `<img src="https://link.jpg" atk-emoticon="testA">`, want: `[testA]`},
		{
			name: "Test2",
			give: `6Rs8OXba\n9u5PiJYiAf   \n<img src="https://link.jpg" atk-emoticon="testA">\n1sKyAAIxssts36Rs8OXba9u5PiJYiAf  \n<img src="https://link.jpg" atk-emoticon="testB">6Rs8OXba9u5PiJYiAf<img src="https://link.jpg" atk-emoticon="">`,
			want: `6Rs8OXba\n9u5PiJYiAf   \n[testA]\n1sKyAAIxssts36Rs8OXba9u5PiJYiAf  \n[testB]6Rs8OXba9u5PiJYiAf[表情]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleEmoticonsImgTagsForNotify(tt.give); got != tt.want {
				t.Errorf("HandleEmoticonsImgTagsForNotify() = %v, want %v", got, tt.want)
			}
		})
	}
}
