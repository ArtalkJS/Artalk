package dao

import "github.com/artalkjs/artalk/v2/internal/utils"

// TrustedOrigins is shared by browser-origin checks and page-title fetching.
// Both configuration entries and administrator-managed site URLs grant trust.
func (dao *Dao) TrustedOrigins(trustedDomains []string) []string {
	urls := append([]string{}, trustedDomains...)
	for _, site := range dao.FindAllSitesCooked() {
		urls = append(urls, site.Urls...)
	}
	origins := make([]string, 0, len(urls))
	for _, rawURL := range urls {
		if origin, err := utils.URLOrigin(rawURL); err == nil {
			origins = append(origins, origin)
		}
	}
	return origins
}
