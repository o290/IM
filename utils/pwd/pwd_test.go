package pwd

import (
	"fmt"
	"testing"
)

func TestHashPwd(t *testing.T) {
	hash := HashPwd("123456")
	fmt.Println(hash)
}
func TestCheckPwd(t *testing.T) {
	ok := CheckPwd("$2a$04$zXedBYXzmLBrJEW9.lVE3eTWpmxlXlfOAdZyKuA2i.ZMpufnltQF.", "1234567")
	fmt.Println(ok)
}
