package entity

type CookedSite struct {
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
	Urls     []string `json:"urls"`
	UrlsRaw  string   `json:"urls_raw"`
	FirstUrl string   `json:"first_url"`
}
