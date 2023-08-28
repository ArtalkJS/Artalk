package artransfer

import "net/url"

// PageKey (commentUrlVal 不确定是否为完整 URL 还是一个 path)
//
// @examples
// ("https://github.com", "/1.html")                => "https://github.com/1.html"
// ("https://github.com", "https://xxx.com/1.html") => "https://github.com/1.html"
// ("https://github.com/", "/1.html")               => "https://github.com/1.html"
// ("", "/1.html")                                  => "/1.html"
// ("", "https://xxx.com/1.html")                   => "https://xxx.com/1.html"
// ("https://github.com/233", "/1/")                => "https://github.com/1/"
func urlResolverGetPageKey(baseUrlRaw string, commentUrlRaw string) string {
	if baseUrlRaw == "" {
		return commentUrlRaw
	}

	baseUrl, err := url.Parse(baseUrlRaw)
	if err != nil {
		return commentUrlRaw
	}

	commentUrl, err := url.Parse(commentUrlRaw)
	if err != nil {
		return commentUrlRaw
	}

	// "https://artalk.js.org/guide/describe.html?233" => "/guide/describe.html?233"
	commentUrl.Scheme = ""
	commentUrl.Host = ""

	// 解决拼接路径中的相对地址，例如：https://atk.xxx/abc/../artalk => https://atk.xxx/artalk
	url := baseUrl.ResolveReference(commentUrl)

	return url.String()
}
