package maps

import (
	"fmt"
	"testing"
)

// 定义嵌套结构体
type InnerStruct struct {
	InnerField string `tag:"inner_field"`
}

// 定义包含指针和嵌套结构体的结构体
type OuterStruct struct {
	PointerField *InnerStruct `tag:"pointer_field"`
	//NestedStruct InnerStruct  `tag:"nested_struct"`
	RegularField int `tag:"regular_field"`
}

func TestRefToMap(t *testing.T) {
	inner := InnerStruct{
		InnerField: "inner_value",
	}
	outer := OuterStruct{
		PointerField: &inner,
		//NestedStruct: inner,
		RegularField: 42,
	}

	result := RefToMap(outer, "tag")
	// 输出转换后的映射
	for k, v := range result {
		fmt.Printf("Key: %s, Value: %v\n", k, v)
	}
}
