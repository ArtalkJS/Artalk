package config

// 默认配置文件名
var CONF_DEFAULT_FILENAMES = [...]string{
	"artalk.yml",
	"artalk-go.yml", // for 向下兼容
	"config.yml",
	"conf.yml",

	// .yaml extension
	"artalk.yaml",
	"artalk-go.yaml",
	"config.yml",
	"conf.yml",
}

// 所有站点
const ATK_SITE_ALL = "__ATK_SITE_ALL"

// Cookie 键
const COOKIE_KEY_ATK_AUTH = "ATK_AUTH"

// ctx keys
const CTX_KEY_ATK_SITE_ID = "atk_site_id"
const CTX_KEY_ATK_SITE_NAME = "atk_site_name"
const CTX_KEY_ATK_SITE_ALL = "atk_site_all"

// 图片上传目录路由重写路径
const IMG_UPLOAD_PUBLIC_PATH = "/static/images"
