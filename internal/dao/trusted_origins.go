package dao

import "github.com/artalkjs/artalk/v2/internal/utils"

// TrustedURLs combines configuration entries and administrator-managed site URLs.
func (dao *Dao) TrustedURLs(trustedDomains []string) []string {
	urls := append([]string{}, trustedDomains...)
	for _, site := range dao.FindAllSitesCooked() {
		urls = append(urls, site.Urls...)
	}
	return urls
}

// TrustedOrigins excludes wildcard patterns so they cannot authorize server-side requests.
func (dao *Dao) TrustedOrigins(trustedDomains []string) []string {
	urls := dao.TrustedURLs(trustedDomains)
	origins := make([]string, 0, len(urls))
	for _, rawURL := range urls {
		if origin, err := utils.URLOrigin(rawURL); err == nil {
			origins = append(origins, origin)
		}
	}
	return origins
}
