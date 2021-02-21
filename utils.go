package main

import (
	"log"
	"math"
	"strconv"

	"github.com/incognitochain/incognito-chain/privacy/operation"
)

func GetShardIDFromLastByte(b byte) byte {
	return byte(int(b) % 8)
}

func parseIndex(s string) uint32 {
	index, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		log.Fatalln("Couldn't parse index:", err)
	} else if index > math.MaxUint32 {
		log.Fatalf("Index too large (max %v)", math.MaxUint32)
	}
	return uint32(index)
}

func byteArrayToScalarArray(bytes []byte) []*operation.Scalar {
	var result []*operation.Scalar
	maxlen := len(bytes) / 32
	for i := 0; i < maxlen; i++ {
		sc := operation.Scalar{}
		sc.FromBytesS(bytes[i*32 : (i+1)*32])
		result = append(result, &sc)
	}
	return result
}
