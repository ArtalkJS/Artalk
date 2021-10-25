package importer

var WordPressImporter = &_TypechoImporter{
	Importer: Importer{
		Name: "WordPress",
		Desc: "",
	},
}

type _WordPressImporter struct {
	Importer
}

func (i _WordPressImporter) Run(basic BasicParams, payload []string) {
}
