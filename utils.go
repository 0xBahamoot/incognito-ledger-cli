package main

import (
	"log"
	"math"
	"strconv"
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
