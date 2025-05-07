package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"strings"
)

func InitList(list []string, key string) bool {
	for _, s := range list {
		if s == key {
			return true
		}
	}
	return false
}

// InitListByRegex  正则表达式
func InitListByRegex(list []string, key string) (ok bool) {
	for _, s := range list {
		regex, err := regexp.Compile(s)
		if err != nil {
			logx.Error(err)
			return
		}
		if regex.MatchString(key) {
			return true
		}
	}
	return false
}

func MD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	cipherStr := h.Sum(nil)
	fmt.Println(cipherStr)
	return hex.EncodeToString(cipherStr)
}
func GetFilePrefix(fileName string) (prefix string) {
	nameList := strings.Split(fileName, ".")
	for i := 0; i < len(nameList)-1; i++ {
		if i == len(nameList)-2 {
			prefix += nameList[i]
			continue
		} else {
			prefix += nameList[i] + "."
		}
	}
	return
}

// DeduplicationList  去重
// [T string | int | uint | uint32] 泛型类约束，只有string | int | uint | uint32可调用
func DeduplicationList[T string | int | uint | uint32](req []T) (response []T) {
	i32Map := make(map[T]bool)
	//遍历req并标记
	for _, i32 := range req {
		if !i32Map[i32] {
			i32Map[i32] = true
		}
	}

	//收集去重后的元素
	for key, _ := range i32Map {
		response = append(response, key)
	}
	return
}
