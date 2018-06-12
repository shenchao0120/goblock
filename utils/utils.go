package utils

import (
	"bytes"
	"encoding/binary"
	"log"
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
