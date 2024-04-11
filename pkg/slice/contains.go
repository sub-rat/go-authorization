package slice

import (
	"encoding/json"
)

func Contains(items []interface{}, item interface{}) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}
	return false
}

func ContainsInt(items []int, item int) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}
	return false
}

func ContainsInt64(items []int64, item int64) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}
	return false
}

func ContainsString(items []string, item string) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}
	return false
}

func StructToSlice[T comparable, I comparable](item I) ([]T, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	var da []T
	err = json.Unmarshal(data, &da)
	if err != nil {
		return nil, err
	}
	return da, nil
}
