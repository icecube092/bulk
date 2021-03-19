package bulk

import "gorm.io/gorm"

type Bulker interface {
	Insert(entity interface{}) error
	Wait()
}

type BulkConfig struct {
	Db           *gorm.DB
	Size         int
	Models       []interface{}
	ModelsSave   []interface{}
	TickEverySec int
	FileName     string
}
