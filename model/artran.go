package model

// 数据行囊 (n 个 Artran 组成一个 Artrans)
// Fields All String type (FAS)
type Artran struct {
	ID  string `json:"id"`
	Rid string `json:"rid"`

	Content string `json:"content"`

	UA          string `json:"ua"`
	IP          string `json:"ip"`
	IsCollapsed string `json:"is_collapsed"` // bool => string "true" or "false"
	IsPending   string `json:"is_pending"`   // bool
	IsPinned    string `json:"is_pinned"`    // bool

	// vote
	VoteUp   string `json:"vote_up"`
	VoteDown string `json:"vote_down"`

	// date
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	// user
	Nick       string `json:"nick"`
	Email      string `json:"email"`
	Link       string `json:"link"`
	BadgeName  string `json:"badge_name"`
	BadgeColor string `json:"badge_color"`

	// page
	PageKey       string `json:"page_key"`
	PageTitle     string `json:"page_title"`
	PageAdminOnly string `json:"page_admin_only"` // bool

	// site
	SiteName string `json:"site_name"`
	SiteUrls string `json:"site_urls"`
}
