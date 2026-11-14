package entity

type EmailCommentDTO struct {
	CommentDTO
	Content    string  `json:"content"`
	ContentRaw string  `json:"content_raw"`
	Nick       string  `json:"nick"`
	Email      string  `json:"email"`
	IP         string  `json:"ip"`
	Datetime   string  `json:"datetime"`
	Date       string  `json:"date"`
	Time       string  `json:"time"`
	PageKey    string  `json:"page_key"`
	PageTitle  string  `json:"page_title"`
	Page       PageDTO `json:"page"`
	SiteName   string  `json:"site_name"`
	Site       SiteDTO `json:"site"`
}
