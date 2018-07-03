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


func TestBase58Encode(t *testing.T) {
	//b:=[]byte{byte(0x00),byte(0x00),byte(0x00),byte(0x00),byte('1')}
	b:=[]byte("dasasff12324")
	fmt.Printf("b[%v]\n",b)
	res:=Base58Encode(b)
	fmt.Printf("%v，len[%d]\n",res,len(res))

	deres:=Base58Decode(res)
	fmt.Printf("result:%v，len[%d]",deres,len(deres))

}