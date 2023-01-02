package entity

type CookedPage struct {
	ID        uint   `json:"id"`
	AdminOnly bool   `json:"admin_only"`
	Key       string `json:"key"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	SiteName  string `json:"site_name"`
	VoteUp    int    `json:"vote_up"`
	VoteDown  int    `json:"vote_down"`
	PV        int    `json:"pv"`
}
