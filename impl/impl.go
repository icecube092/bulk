package impl

import (
	"bulk"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
	"reflect"
	"sync"
	"time"
)

type Bulk struct {
	db         *gorm.DB
	size       int
	mux        sync.Mutex
	models     []interface{}
	modelsSave []interface{}
	ticker     *time.Ticker
	file       *os.File
}

func New(config *bulk.BulkConfig) bulk.Bulker {
	var err error
	b := &Bulk{
		db:         config.Db,
		size:       config.Size,
		models:     config.Models,
		modelsSave: config.ModelsSave,
		ticker:     time.NewTicker(time.Second * time.Duration(config.TickEverySec)),
	}
	b.file, err = os.Create(config.FileName)
	err = b.check()
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func (b *Bulk) check() error { // todo check match models and modelsSave
	for _, m := range b.models {
		if reflect.TypeOf(m).Kind() == reflect.Array ||
			reflect.TypeOf(m).Kind() == reflect.Slice {
			return errors.New(fmt.Sprintf("models must be not array/slice"))
		} else if reflect.TypeOf(m).Kind() != reflect.Ptr {
			return errors.New("models must be ptr")
		}
	}
	for _, model := range b.modelsSave {
		if reflect.ValueOf(model).Elem().Kind() != reflect.Array &&
			reflect.ValueOf(model).Elem().Kind() != reflect.Slice {
			return errors.New("modelsSave must be array/slice")
		}
		if reflect.TypeOf(model).Kind() != reflect.Ptr {
			return errors.New("modelsSave must be ptr")
		}
	}
	//if b.db == nil {
	//	return errors.New("nil db")
	//}
	return nil
}

func (b *Bulk) Insert(entity interface{}) error {
	if !isPointer(entity) {
		// todo err
	}
	for _, model := range b.models {
		if reflect.TypeOf(model).String() == reflect.TypeOf(entity).String() {
			b.mux.Lock()
			defer b.mux.Unlock()
			b.stash(entity)
			return nil
		}
	}
	return errors.New("no model")
}

func (b *Bulk) stash(entity interface{}) {
	for _, models := range b.modelsSave {
		if reflect.TypeOf(models).Elem().Elem().String() == reflect.TypeOf(entity).String() {
			mPtr := reflect.ValueOf(models)
			m := mPtr.Elem()
			m.Set(reflect.Append(m, reflect.ValueOf(entity)))
		}
	}
}

func (b *Bulk) Wait() {
	go func() {
		for {
			for _, models := range b.modelsSave {
				m := reflect.ValueOf(models)
				arrayOfModels := m.Elem()
				if m.Elem().Len() >= b.size {
					fmt.Printf("%+v\n", arrayOfModels.Index(0).Elem())
					buf, err := json.Marshal(arrayOfModels.Index(0).Elem().String())
					_, err = b.file.WriteString(string(buf))
					fmt.Println(err, buf, string(buf), err)
					fmt.Println("DONE: ", reflect.TypeOf(models).String())
					fmt.Println("WRITE ",arrayOfModels.Index(0).Elem().Field(0))
					m.Elem().Set(reflect.MakeSlice(m.Type().Elem(), 0, b.size))
				}
			}
			select {
			case <-b.ticker.C:
				// todo reserve
			}
		}
	}()
}
