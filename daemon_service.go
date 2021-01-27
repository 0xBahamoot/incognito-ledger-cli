package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func getAccountList() (map[string]string, error) {
	result := make(map[string]string)
	resp, err := http.Get("http://" + COINDAEMONADDR + "/getaccountlist")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func submitKeyimages(km map[string][]byte) error {
	return nil
}

func getEncryptKeyImages(accountName string) (map[string]map[string][]byte, error) {
	result := make(map[string]map[string][]byte)
	resp, err := http.Get("http://" + COINDAEMONADDR + "/getcoinstodecrypt?account=" + accountName)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getAccountBalance(accountName string) (uint64, error) {
	var result struct {
		Address string
		Balance uint64
	}
	resp, err := http.Get("http://" + COINDAEMONADDR + "/getbalance?account=" + accountName)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}
	return result.Balance, nil
}

func importAccount() error {
	return nil
}
