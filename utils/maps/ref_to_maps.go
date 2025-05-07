package maps

import (
	"reflect"
)

// RefToMap 将一个结构体类型的数据转换为一个 map[string]any 类型的映射
// data 要转换的数据 tag 筛选结构体字段的标签名
func RefToMap(data any, tag string) map[string]any {
	maps := map[string]any{}
	//1.获取反射信息
	//获取data的反射类型
	t := reflect.TypeOf(data)
	//获取data反射值
	v := reflect.ValueOf(data)

	//2.遍历结构体字段
	//NumField返回结构体t的字段数量
	for i := 0; i < t.NumField(); i++ {
		//Field获取结构体的第 i 个字段的类型信息
		field := t.Field(i)
		//找当前字段是否有指定的 tag 标签
		getTag, ok := field.Tag.Lookup(tag)
		if !ok {
			continue
		}
		//检查字段值是否是零值
		val := v.Field(i)
		if val.IsZero() {
			continue
		}
		//如果当前字段是结构体类型，说明是一个嵌套的结构体，递归调用
		if field.Type.Kind() == reflect.Struct {
			newMaps := RefToMap(val.Elem().Interface(), tag)
			maps[getTag] = newMaps
			continue
		}
		//如果当前字段类型是指针类型
		if field.Type.Kind() == reflect.Ptr {
			if field.Type.Elem().Kind() == reflect.Struct {
				newMaps := RefToMap(val.Elem().Interface(), tag)
				maps[getTag] = newMaps
				continue
			}
			//获取指针指向的值的反射表示，再将其转换为 interface{} 类型的值。
			maps[getTag] = val.Elem().Interface()
			continue
		}
		//如果是其他类型
		//返回 val 所代表的实际值，以 interface{} 类型呈现
		maps[getTag] = val.Interface()
	}
	return maps
}
