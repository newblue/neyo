package neyo

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 将map赋值到一个struct
func ToStruct(m map[string]interface{}, val reflect.Value) {
	if m == nil {
		return
	}
	if val.Kind() == reflect.Ptr && val.Elem().Kind() != reflect.Struct {
		Log(ERROR, "Only supported struct")
		return
	}

	origin_val := val
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for k, v := range m {
		field := val.FieldByName(strings.Title(k))
		if !field.IsValid() {
			continue
		}
		if !field.CanSet() {
			Log(WARN, "CanSet = false, %s %s", k, v)
			continue
		}
		Log(DEBUG, "Found key=%s value=%s %s", k, v, field.Type().String())
		switch field.Kind() {
		case reflect.String:
			Log(DEBUG, "Ready to SetString %s", v)
			if _str, ok := v.(string); ok {
				field.SetString(_str)
			} else {
				field.SetString(fmt.Sprint("%v", v))
			}
		case reflect.Int:
			fallthrough
		case reflect.Int64:
			field.SetInt(ToInt64(v, 0))
		case reflect.Slice: // 字段总是slice
			if _strs, ok := v.([]string); ok {
				field.Set(reflect.ValueOf(_strs))
			} else if _slice, ok := v.([]interface{}); ok {
				strs := make([]string, len(_slice))
				for i, vz := range _slice {
					strs[i] = vz.(string)
				}
				field.Set(reflect.ValueOf(strs))
			} else {
				Log(WARN, "Only []string is supported yet")
			}
		case reflect.Map:
			field.Set(reflect.ValueOf(v))
		case reflect.Ptr:
			// No support yet
		case reflect.Struct:
			if field.Type().String() == "time.Time" {
				if t, ok := v.(time.Time); ok {
					field.Set(reflect.ValueOf(t))
					break
				}
			}
			v2, ok := v.(map[string]interface{})
			if !ok {
				Log(INFO, "Not a map[string]interface{} key=%s value=%s", k, v)
				return
			}
			ToStruct(v2, field)
		default:
			field.Set(reflect.ValueOf(v))
		}
	}
	_ = origin_val
}

func ToInt(v interface{}, defaultValue int) int {
	if v == nil {
		return defaultValue
	}
	if i, ok := v.(int); ok {
		return i
	}
	if i2, ok := v.(int64); ok {
		return int(i2)
	}
	str := fmt.Sprintf("%v", v)
	i, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return i
}

func ToInt64(v interface{}, defaultValue int64) int64 {
	if v == nil {
		return defaultValue
	}
	if i, ok := v.(int64); ok {
		return i
	}
	if i2, ok := v.(int); ok {
		return int64(i2)
	}
	str := fmt.Sprintf("%v", v)
	i, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		return defaultValue
	}
	return i
}
