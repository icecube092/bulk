package impl

import (
	"bulk"
	"fmt"
	"testing"
	"time"
)

type Model struct {
	X int
	Y string
}

func TestBulk_Insert(t *testing.T) {
	b := New(&bulk.BulkConfig{
		Db:         nil,
		Size:       1,
		Models:     []interface{}{&Model{}},
		ModelsSave: []interface{}{&[]*Model{}},
		TickEverySec: 1,
		FileName: "test.json",
	})
	b.Wait()
	m := &Model{X: 1, Y: "bla"}
	err := b.Insert(m)
	fmt.Println(err)
	time.Sleep(time.Second * 1)
}