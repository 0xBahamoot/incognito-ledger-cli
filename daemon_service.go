package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
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

type LedgerRequest struct {
	Cmd  string
	Data []byte
}

func requestCreateTx(txjson []byte, n *NanoS) error {
	c, _, err := websocket.DefaultDialer.Dial("ws://"+COINDAEMONADDR+"/createtx", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	sendMsgCh := make(chan []byte)
	done := make(chan struct{})

	go func() {
		sendMsgCh <- txjson
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var req LedgerRequest
			err = json.Unmarshal(message, &req)
			if err != nil {
				log.Println("read:", err)
				return
			}
			//TODO
			switch req.Cmd {
			case "signschnorr":
			case "createringsig":
			case "result":
			default:
				log.Println("unknown command")
			}
		}
	}()
	for {
		select {
		case <-done:
			return nil
		case msg := <-sendMsgCh:
			err := c.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write:", err)
				return err
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}
