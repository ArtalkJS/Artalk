package entity

type CookedUser struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Link         string   `json:"link"`
	BadgeName    string   `json:"badge_name"`
	BadgeColor   string   `json:"badge_color"`
	IsAdmin      bool     `json:"is_admin"`
	SiteNames    []string `json:"site_names"`
	SiteNamesRaw string   `json:"site_names_raw"`
	ReceiveEmail bool     `json:"receive_email"`
}
