package email

import "testing"

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
			if got := HandleEmoticonsImgTagsForNotify(tt.give); got != tt.want {
				t.Errorf("HandleEmoticonsImgTagsForNotify() = %v, want %v", got, tt.want)
			}
		})
	}
}
