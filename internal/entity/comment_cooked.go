package entity

type CookedComment struct {
	ID             uint   `json:"id"`
	Content        string `json:"content"`
	ContentMarked  string `json:"content_marked"`
	UserID         uint   `json:"user_id"`
	Nick           string `json:"nick"`
	EmailEncrypted string `json:"email_encrypted"`
	Link           string `json:"link"`
	UA             string `json:"ua"`
	Date           string `json:"date"`
	IsCollapsed    bool   `json:"is_collapsed"`
	IsPending      bool   `json:"is_pending"`
	IsPinned       bool   `json:"is_pinned"`
	IsAllowReply   bool   `json:"is_allow_reply"`
	Rid            uint   `json:"rid"`
	BadgeName      string `json:"badge_name"`
	BadgeColor     string `json:"badge_color"`
	IP             string `json:"-"`
	IPRegion       string `json:"ip_region,omitempty"`
	Visible        bool   `json:"visible"`
	VoteUp         int    `json:"vote_up"`
	VoteDown       int    `json:"vote_down"`
	PageKey        string `json:"page_key"`
	PageURL        string `json:"page_url"`
	SiteName       string `json:"site_name"`
}
