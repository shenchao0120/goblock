package utils

import (
	"fmt"
	"strconv"
	"testing"
)

func TestConvIntToHex(t *testing.T) {
	a := 1000000

	aBytes := ConvIntToHex(int64(a))
	fmt.Println("result1:", aBytes)

	fmt.Println("result2:", []byte(strconv.FormatInt(int64(a), 10)))
}


func TestRemoveIntInSlice(t *testing.T) {
	s:=[]int{0}
	s = RemoveIntInSlice(s,0)
	fmt.Println(s)
}