package importer

var WordPressImporter = &_TypechoImporter{
	ImporterInfo: ImporterInfo{
		Name: "wordpress",
		Desc: "从 WordPress 导入数据",
		Note: "",
	},
}

type _WordPressImporter struct {
	ImporterInfo
}

func (i _WordPressImporter) Run(basic BasicParams, payload []string) {
}
