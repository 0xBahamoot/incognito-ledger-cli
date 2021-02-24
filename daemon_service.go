package main

import (
	"bytes"
	"encoding/hex"
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

func getDaemonVersion() (string, error) {
	resp, err := http.Get("http://" + COINDAEMONADDR + "/getaccountlist")
	if err != nil {
		log.Fatalln(err)
		return "", nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return "", nil
	}
	return string(body), nil
}

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
	fmt.Println(resp.StatusCode)
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

func importAccount(name, addr, otaKey, viewKey string, beaconHeight uint64) error {
	type API_import_account_req struct {
		AccountName    string
		PaymentAddress string
		OTAKey         string
		Viewkey        string
		BeaconHeight   uint64
	}

	reqdata := API_import_account_req{
		AccountName:    name,
		OTAKey:         otaKey,
		Viewkey:        viewKey,
		PaymentAddress: addr,
		BeaconHeight:   beaconHeight,
	}
	reqBytes, err := json.Marshal(reqdata)
	if err != nil {
		panic(err)
	}
	fmt.Println("reqBytes", string(reqBytes))
	// 12su5Urq6hucGGNEdk37RXJW1mY2LAGcrgjdYJ4uhzj9K4F47SRFSkLSzCcz7uJ2mAUTwnrA5mkaCvzobTc6ceocdNAhRQgZeveaLQmMkxJqueSYm9gKkNV39ba1CvR5n3Euig9gNLeP1TkwonfZ
	// 131iy88imFE4QUxJP8bURMkTf9B1YTwETcXWRvLRX9XmPD4ABG8wWk2YHJi5QkMURcp7HYPeRvpsD4h75mJcKaBec5A8RjQCZAm1FXn6oJZ59XW6DP54VZF
	req, err := http.NewRequest("POST", "http://"+COINDAEMONADDR+"/importaccount", bytes.NewBuffer(reqBytes))
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

type LedgerRequest struct {
	Cmd  string
	Data []byte
}

func requestCreateTx(txjsonFile string) (string, error) {
	var txID string
	data, err := ioutil.ReadFile(txjsonFile)
	if err != nil {
		panic(err)
	}
	c, _, err := websocket.DefaultDialer.Dial("ws://"+COINDAEMONADDR+"/createtx", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	sendMsgCh := make(chan []byte)
	done := make(chan struct{})
	var nanos *NanoS

	nanos, err = OpenNanoS()
	if err != nil {
		log.Println("This cmd require connected to ledger device")
		log.Fatalln("Couldn't open device:", err)
	}

	go func() {
		fmt.Println("sendMsgCh <- data", data)
		sendMsgCh <- data
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
			switch req.Cmd {
			case "signschnorr":
				type ReqStruct struct {
					PedRandom  []byte
					PedPrivate []byte
					Randomness []byte
					Message    []byte
				}
				requestData := ReqStruct{}
				err := json.Unmarshal(req.Data, &requestData)
				if err != nil {
					panic(err)
				}
				sig, err := nanos.SignSchnorr(requestData.PedRandom, requestData.PedPrivate, requestData.Randomness, requestData.Message)
				if err != nil {
					panic(err)
				}
				sendMsgCh <- sig
			case "genalpha":
				type ReqStruct struct {
					AlphaLength int
				}
				requestData := ReqStruct{}
				err := json.Unmarshal(req.Data, &requestData)
				if err != nil {
					panic(err)
				}
				fmt.Println("genalpha with AlphaLength", requestData.AlphaLength)
				err = nanos.GenerateAlpha(requestData.AlphaLength)
				if err != nil {
					panic(err)
				}
				sendMsgCh <- []byte("success")
			case "gencoinprivate":
				type ReqStruct struct {
					CoinsH [][]byte
				}
				requestData := ReqStruct{}
				err := json.Unmarshal(req.Data, &requestData)
				if err != nil {
					panic(err)
				}
				fmt.Println("gencoinprivate with CoinsH", len(requestData.CoinsH))
				err = nanos.GenCoinPrivateKey(requestData.CoinsH)
				if err != nil {
					panic(err)
				}
				sendMsgCh <- []byte("success")
			case "calculatec": // calculate 1st C
				type ReqStruct struct {
					Rpi     [][]byte
					PedComG []byte
				}
				requestData := ReqStruct{}
				err := json.Unmarshal(req.Data, &requestData)
				if err != nil {
					panic(err)
				}
				firstC, err := nanos.CalculateFirstC(requestData.Rpi, requestData.PedComG)
				if err != nil {
					panic(err)
				}
				sendMsgCh <- firstC
			case "calculater": // calculate r
				type ReqStruct struct {
					CoinLength int
					Cpi        []byte
				}
				requestData := ReqStruct{}
				err := json.Unmarshal(req.Data, &requestData)
				if err != nil {
					panic(err)
				}
				new_rPi, err := nanos.CalculateR(requestData.CoinLength, requestData.Cpi)
				if err != nil {
					panic(err)
				}
				new_rPiBytes, err := json.Marshal(new_rPi)
				if err != nil {
					panic(err)
				}
				sendMsgCh <- new_rPiBytes
			case "result":
				fmt.Println(string(req.Data), hex.EncodeToString(req.Data))
				return
			default:
				log.Println("unknown command")
			}
		}
	}()
	for {
		select {
		case <-done:
			return txID, nil
		case msg := <-sendMsgCh:
			err := c.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write:", err)
				return txID, err
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return txID, err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return txID, nil
		}
	}
}
