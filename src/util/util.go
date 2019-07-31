package util

import (	
	"reflect"	
	"time"
	"math/rand"	
)

func InSlice(slice interface{}, elem interface{}) bool {
	if reflect.TypeOf(slice).Kind() != reflect.Slice || reflect.TypeOf(slice).Elem() != reflect.TypeOf(elem) {
		return false				
	}	
	if t := reflect.TypeOf(elem).Kind(); !(t >= reflect.Int && t <= reflect.Int64 || t >= reflect.Uint && t <= reflect.Uint64 ||
		t == reflect.Float32 || t == reflect.Float64 || t == reflect.String) {
		return false
	}
	valSlice := reflect.ValueOf(slice)
	valElem := reflect.ValueOf(elem)
	for i := 0; i < valSlice.Len(); i++ {
		if valSlice.Index(i).Interface() == valElem.Interface() {
			return true
		}
	}
	return false
}

func InSlice2(slice interface{}, equal func(i int) bool) bool {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return false
	}
	valSlice := reflect.ValueOf(slice)
	for i := 0; i < valSlice.Len(); i++ {
		if equal(i) {
			return true
		}
	}
	return false
}

func RemoveSliceElem(slice interface{}, elem interface{}, all bool) interface{} {
	if reflect.TypeOf(slice).Kind() != reflect.Slice || reflect.TypeOf(slice).Elem() != reflect.TypeOf(elem) {
		return slice				
	}	
	if t := reflect.TypeOf(elem).Kind(); !(t >= reflect.Int && t <= reflect.Int64 || t >= reflect.Uint && t <= reflect.Uint64 ||
		t == reflect.Float32 || t == reflect.Float64 || t == reflect.String) {
		return slice
	}
	valSlice := reflect.ValueOf(slice)
	valElem := reflect.ValueOf(elem)
	for i := 0; i < valSlice.Len(); i++ {
		if valSlice.Index(i).Interface() == valElem.Interface() {
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

func RemoveSliceElem2(slice interface{}, equal func(i int) bool, all bool) interface{} {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return slice	
	}
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

func UniqueSlice(slice interface{}, bSort bool) interface{} {	
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return slice
	}	
	if t := reflect.TypeOf(slice).Elem().Kind(); !(t >= reflect.Int && t <= reflect.Int64 || t >= reflect.Uint && t <= reflect.Uint64 ||
		t == reflect.Float32 || t == reflect.Float64 || t == reflect.String) {
		return slice
	}
	valSlice := reflect.ValueOf(slice)
	if valSlice.Len() < 2 {
		return slice
	}	
	for i := 0; i < valSlice.Len(); i++ {
		bDel := false
		if bSort {
			if i > 0 && valSlice.Index(i-1).Interface() == valSlice.Index(i).Interface() {
				bDel = true
			}
		} else {			
			for j := 0; j < i; j++ {
				if valSlice.Index(j).Interface() == valSlice.Index(i).Interface() {
					bDel = true
					break
				}
			}			
		}				
		if bDel {
			valSlice = reflect.AppendSlice(valSlice.Slice(0, i), valSlice.Slice(i+1, valSlice.Len()))
			i--
		}
	}
	return valSlice.Interface()
}

func UniqueSlice2(slice interface{}, equal func(i, j int) bool,  bSort bool) interface{} {	
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return slice
	}		
	valSlice := reflect.ValueOf(slice)
	if valSlice.Len() < 2 {
		return slice
	}
	for i := 0; i < valSlice.Len(); i++ {
		bDel := false
		if bSort {
			if i > 0 && equal(i-1, i) {
				bDel = true
			}
		} else {			
			for j := 0; j < i; j++ {
				if equal(j, i) {
					bDel = true
					break
				}
			}			
		}				
		if bDel {
			valSlice = reflect.AppendSlice(valSlice.Slice(0, i), valSlice.Slice(i+1, valSlice.Len()))
			i--
		}
	}
	return valSlice.Interface()	
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

func GetRandom(start int32, end int32) int32 {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	if start > end {
		return 0
	}
	return start + r.Int31n(end-start+1)
}

func GetNormRandom(start int32, end int32) int32 {
	if start > end {
		return 0
	} else if start == end {
		return start
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		ret := int32(r.NormFloat64()*float64(end-start)/6 + float64(start+end)/2)
		if ret >= start && ret <= end {
			return ret
		}
	}
	return 0
}