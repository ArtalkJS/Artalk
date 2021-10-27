package artransfer

var ExampleImporter = &_ExampleImporter{
	ImporterInfo: ImporterInfo{
		Name: "example",
		Desc: "从 Example 导入数据",
		Note: "",
	},
}

type _ExampleImporter struct {
	ImporterInfo
}

func (imp *_ExampleImporter) Run(basic *BasicParams, payload []string) {
}
