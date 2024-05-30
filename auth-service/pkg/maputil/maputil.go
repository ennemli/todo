package maputil

import (
	"reflect"
)

func AnyKeys(m interface{}, keys ...string) bool {
	mv := reflect.ValueOf(m)
	if mv.Kind() != reflect.Map {
		return false
	}

	keySet := map[string]bool{}
	for _, k := range keys {
		keySet[k] = true
	}

	iter := mv.MapRange()
	for iter.Next() {
		k := iter.Key().String()
		if !keySet[k] {
			return false
		}
	}
	if len(keys) > 0 && mv.Len() == 0 {
		return false
	}
	return true
}
