package artransfer

import "net/url"

// Page key may be a relative path or a full URL,
// this function will resolve the full URL based on the baseURL.
func getResolvedPageKey(baseURL string, pageKey string) string {
	if baseURL == "" {
		return pageKey
	}

	baseUrl, err := url.Parse(baseURL)
	if err != nil {
		return pageKey
	}

	commentUrl, err := url.Parse(pageKey)
	if err != nil {
		return pageKey
	}

	// "https://artalk.js.org/guide/describe.html?233" => "/guide/describe.html?233"
	commentUrl.Scheme = ""
	commentUrl.Host = ""

	// Remove dots and slashes from the relative path, e.g. "https://atk.xxx/abc/../artalk" => "https://atk.xxx/artalk"
	url := baseUrl.ResolveReference(commentUrl)

	return url.String()
}
