package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"fmt"
)

// Convert int64 to []byte
func ConvIntToHex(num int64) []byte {
	buf := new(bytes.Buffer)
	//buf.Write()
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

func RemoveIntInSlice (slices []int ,a int) []int{
	for index,b:=range slices {
		if b == a {
				slices = append(slices[0:index], slices[index+1:]...)
				break
		}
	}
	return slices
}


func ContainInt(slices []int ,a int) bool{
	for _,b:=range slices {
		if b == a {
			return true
		}
	}
	return false
}
