package pkg

import (
	"errors"
	"reflect"
	"strings"
)

func ArrayContains(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func ArrayJoins(array []string, sep string, format func(string) string) string {
	for k, v := range array {
		if format != nil {
			array[k] = format(v)
		}
	}
	return strings.Join(array, sep)
}

func ArrayParamJoins(array [][]string) string {
	var temp []string
	for _, v := range array {
		if len(v) != 2 {
			panic("参数错误")
		}
		temp = append(temp, LineToLowCamel(v[0])+" "+v[1])
	}
	return strings.Join(temp, ",")
}

type QuickInArray struct {
	haystackMap map[interface{}]bool
}

func NewQuickInArray(haystack interface{}) (*QuickInArray, error) {
	// 判断haystack是否是切片
	haystackValue := reflect.ValueOf(haystack)
	if haystackValue.Kind() != reflect.Slice {
		return nil, errors.New("haystack is not slice")
	}
	// 将haystack转换为map
	haystackMap := make(map[interface{}]bool)
	for i := 0; i < haystackValue.Len(); i++ {
		haystackMap[haystackValue.Index(i).Interface()] = true
	}

	return &QuickInArray{
		haystackMap: haystackMap,
	}, nil
}

func (q *QuickInArray) InArray(needle interface{}) bool {
	_, ok := q.haystackMap[needle]
	return ok
}
