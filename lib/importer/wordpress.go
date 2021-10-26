package importer

var WordPressImporter = &_WordPressImporter{
	ImporterInfo: ImporterInfo{
		Name: "wordpress",
		Desc: "从 WordPress 导入数据",
		Note: "",
	},
}

type _WordPressImporter struct {
	ImporterInfo
}

func (imp *_WordPressImporter) Run(basic *BasicParams, payload []string) {
}
