package impl

import "reflect"

func isPointer(entity interface{}) bool {
	if reflect.TypeOf(entity).Kind() != reflect.Ptr {
		return false
	}
	return true
}
