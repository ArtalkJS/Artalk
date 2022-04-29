package model

func MigrateModels() {
	// Migrate the schema
	DB().AutoMigrate(&Site{}, &Page{}, &User{},
		&Comment{}, &Notify{}, &Vote{}) // 注意表的创建顺序，因为有关联字段
}
