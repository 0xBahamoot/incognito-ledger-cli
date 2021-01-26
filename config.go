package main

import (
	"encoding/json"
	"io/ioutil"
)

var COINDAEMONADDR string

func init() {
	COINDAEMONADDR = DefaultCoinDaemonAddr
}

func readConfig() {
	data, err := ioutil.ReadFile("./cfg.json")
	if err != nil {
		panic(err)
	}

	type ConfigJSON struct {
	}
	var cfgJson ConfigJSON
	err = json.Unmarshal(data, &cfgJson)
	if err != nil {
		panic(err)
	}
}
