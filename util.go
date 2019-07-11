package util

import (	
	"reflect"	
)

func InSlice(slice interface{}, elem interface{}) bool {
	if reflect.TypeOf(slice).Kind() == reflect.Slice && reflect.TypeOf(slice).Elem() == reflect.TypeOf(elem) {
		valSlice := reflect.ValueOf(slice)
		valElem := reflect.ValueOf(elem)
		for i := 0; i < valSlice.Len(); i++ {
			if reflect.Int <= valElem.Kind() && valElem.Kind() <= reflect.Int64 && valSlice.Index(i).Int() == valElem.Int() ||
				reflect.Uint <= valElem.Kind() && valElem.Kind() <= reflect.Uint64 && valSlice.Index(i).Uint() == valElem.Uint() ||
				(valElem.Kind() == reflect.Float32 || valElem.Kind() == reflect.Float64) && valSlice.Index(i).Float() == valElem.Float() ||
				valElem.Kind() == reflect.String && valSlice.Index(i).String() == valElem.String() {
				return true
			}
		}
	}
	return false
}

func InSlice2(slice interface{}, equal func(i int) bool) bool {
	if reflect.TypeOf(slice).Kind() == reflect.Slice {
		valSlice := reflect.ValueOf(slice)
		for i := 0; i < valSlice.Len(); i++ {
			if equal(i) {
				return true
			}
		}
	}
	return false
}

func RemoveSliceElem(slice interface{}, elem interface{}, all bool) interface{} {
	if reflect.TypeOf(slice).Kind() == reflect.Slice && reflect.TypeOf(slice).Elem() == reflect.TypeOf(elem) {
		valSlice := reflect.ValueOf(slice)
		valElem := reflect.ValueOf(elem)
		for i := 0; i < valSlice.Len(); i++ {
			if reflect.Int <= valElem.Kind() && valElem.Kind() <= reflect.Int64 && valSlice.Index(i).Int() == valElem.Int() ||
				reflect.Uint <= valElem.Kind() && valElem.Kind() <= reflect.Uint64 && valSlice.Index(i).Uint() == valElem.Uint() ||
				(valElem.Kind() == reflect.Float32 || valElem.Kind() == reflect.Float64) && valSlice.Index(i).Float() == valElem.Float() ||
				valElem.Kind() == reflect.String && valSlice.Index(i).String() == valElem.String() {
				valSlice = reflect.AppendSlice(valSlice.Slice(0, i), valSlice.Slice(i+1, valSlice.Len()))
				if all {
					i--
				} else {
					break
				}
			}
		}
		return valSlice.Interface()
	}
	return slice
}

func RemoveSliceElem2(slice interface{}, equal func(i int) bool, all bool) interface{} {
	if reflect.TypeOf(slice).Kind() == reflect.Slice {
		valSlice := reflect.ValueOf(slice)
		for i := 0; i < valSlice.Len(); i++ {
			if equal(i) {
				valSlice = reflect.AppendSlice(valSlice.Slice(0, i), valSlice.Slice(i+1, valSlice.Len()))
				if all {
					i--
				} else {
					break
				}
			}
		}
		return valSlice.Interface()
	}
	return slice
}

func UniqueSlice(slice interface{}) interface{} {
	//slice提前排序
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return slice
	}
	valSlice := reflect.ValueOf(slice)
	if valSlice.Len() < 2 {
		return slice
	}
	ret := reflect.MakeSlice(reflect.TypeOf(slice), 0, 0)
	for i := 0; i < valSlice.Len(); i++ {
		if i > 0 && valSlice.Index(i-1).Interface() == valSlice.Index(i).Interface() {
			continue
		}
		ret = reflect.Append(ret, valSlice.Index(i))
	}
	return ret.Interface()
}

func CopyMap(m interface{}) interface{} {
	if reflect.TypeOf(m).Kind() != reflect.Map {
		return m
	}
	valMap := reflect.ValueOf(m)
	ret := reflect.MakeMap(reflect.TypeOf(m))
	for _, key := range valMap.MapKeys() {
		ret.SetMapIndex(key, valMap.MapIndex(key))
	}
	return ret.Interface()
}