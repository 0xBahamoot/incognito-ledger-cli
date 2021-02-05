package main

import (
	"io/ioutil"
)

// type TxJson struct {
// 	Account    string      `json:"account"`
// 	PrivateKey string      `json:"privatekey"`
// 	TxType     string      `json:"type"`
// 	TxParams   interface{} `json:"params"`
// }

func RequestCreateTx(txjsonFile string, n *NanoS) error {
	data, err := ioutil.ReadFile(txjsonFile)
	if err != nil {
		panic(err)
	}
	return requestCreateTx(data, n)
}
