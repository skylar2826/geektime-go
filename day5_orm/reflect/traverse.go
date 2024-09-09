package reflect

import "reflect"

func IterateArrayOrSlice(entity any) ([]any, error) {
	val := reflect.ValueOf(entity)

	res := make([]any, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		res = append(res, val.Index(i).Interface())
	}
	return res, nil
}

func IterateMap(entity any) (map[string]any, error) {
	val := reflect.ValueOf(entity)

	res := make(map[string]any, val.Len())

	for _, key := range val.MapKeys() {
		res[key.Interface().(string)] = val.MapIndex(key).Interface()
	}

	return res, nil
}
