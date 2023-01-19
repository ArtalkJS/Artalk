package artransfer

var ExampleImporter = &_ExampleImporter{
	ImporterInfo: ImporterInfo{
		Name: "example",
		Desc: "Import data from <EXAMPLE>",
		Note: "",
	},
}

type _ExampleImporter struct {
	ImporterInfo
}

func (imp *_ExampleImporter) Run(basic *BasicParams, payload []string) {
}
