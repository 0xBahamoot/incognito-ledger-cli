package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func submitKeyimages(tokenID string, account string, kms map[string]string) error {
	var reqBody struct {
		Account   string
		Keyimages map[string]map[string]string
	}
	reqKms := make(map[string]map[string]string)
	for coinpub, km := range kms {
		if _, ok := reqKms[tokenID]; !ok {
			reqKms[tokenID] = make(map[string]string)
		}
		reqKms[tokenID][coinpub] = km
	}
	reqBody.Account = account
	reqBody.Keyimages = reqKms

	reqBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "http://"+COINDAEMONADDR+"/submitkeyimages", bytes.NewBuffer(reqBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return nil
}

func getEncryptKeyImages(accountName string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)
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

func getAccountBalance(accountName string) (map[string]uint64, error) {
	var result struct {
		Address string
		Balance map[string]uint64
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
		return nil, err
	}
	return result.Balance, nil
}

func importAccount() error {
	return nil
}
