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

// 图片上传目录路由重写路径
const IMG_UPLOAD_PUBLIC_PATH = "/static/images"
