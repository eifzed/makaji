package utility

import (
	"reflect"
)

func StringExistInSlice(item string, itemSlice []string) bool {
	for _, i := range itemSlice {
		if i == item {
			return true
		}
	}
	return false
}

// func RoleExistInSlice(item jwt.Role, itemSlice []jwt.Role) bool {
// 	for _, i := range itemSlice {
// 		if i.ID == item.ID {
// 			return true
// 		}
// 	}
// 	return false
// }

func IntExistInSlice(item int, itemSlice []int) bool {
	for _, i := range itemSlice {
		if i == item {
			return true
		}
	}
	return false
}

func Int32ExistInSlice(item int32, itemSlice []int32) bool {
	for _, i := range itemSlice {
		if i == item {
			return true
		}
	}
	return false
}

func Int64ExistInSlice(item int64, itemSlice []int64) bool {
	for _, i := range itemSlice {
		if i == item {
			return true
		}
	}
	return false
}

func ConvertSliceToSliceOfInterface(slice interface{}) []interface{} {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return []interface{}{slice}
	}

	interfaceSlice := make([]interface{}, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		interfaceSlice[i] = sliceValue.Index(i).Interface()
	}

	return interfaceSlice
}
