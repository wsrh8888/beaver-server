package maps

import "reflect"

func RefToMap(data any, tag string) map[string]any {
	maps := map[string]any{}
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		getTag, ok := field.Tag.Lookup(tag)
		if !ok {
			continue
		}
		val := v.Field(i)
		if val.IsZero() {
			continue
		}
		if field.Type.Kind() == reflect.Ptr {
			if field.Type.Elem().Kind() == reflect.Struct {
				newMaps := RefToMap(val.Elem().Interface(), tag)
				maps[getTag] = newMaps
			} else {
				maps[getTag] = val.Elem().Interface()
			}
		} else {
			maps[getTag] = val.Interface()
		}
	}
	return maps
}
