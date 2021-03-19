package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x = 0
	arr := []int{1}
	fmt.Println(arr)
	e := reflect.TypeOf(x).String()
	fmt.Println(e)
	ae := reflect.TypeOf(arr[0]).String()
	fmt.Println(ae)
	fmt.Printf("%T\n", reflect.ValueOf(arr))
	models := &[]*int{}
	var xx int
	xx = 1
	mPtr := reflect.ValueOf(models)
	m := mPtr.Elem()
	m.Set(reflect.Append(m, reflect.ValueOf(&xx)))
	fmt.Println(models)
}

func unknown() interface{} {
	return []int{}

}