package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

func (n *NanoS) GetVersion() (version string, err error) {
	resp, err := n.Exchange(cmdGetVersion, 0, 0, nil)
	if err != nil {
		return "", err
	} else if len(resp) != 3 {
		fmt.Printf("%v\n", resp)
		return "", errors.New("version has wrong length")
	}
	return fmt.Sprintf("v%d.%d.%d", resp[0], resp[1], resp[2]), nil
}

func (n *NanoS) TrustHost() error {
	resp, err := n.Exchange(cmdTrustHost, 0, 0, nil)
	if err != nil {
		return err
	}
	_ = resp
	return nil
}

func (n *NanoS) GetAddress() (addr string, err error) {
	// encIndex := make([]byte, 4)
	// binary.LittleEndian.PutUint32(encIndex, index)

	resp, err := n.Exchange(cmdGetAddress, 0, 0, nil)
	if err != nil {
		return
	}
	fmt.Printf("address %v\n", string(resp[:]))
	addr = string(resp[:])
	return
}

func (n *NanoS) GetPrivateKey() (priv string, err error) {
	resp, err := n.Exchange(cmdGetPrivateKey, 0, 0, nil)
	if err != nil {
		return
	}
	fmt.Println("privatekey:", string(resp))
	priv = string(resp)
	return
}

func (n *NanoS) GetViewKey() (string, error) {
	resp, err := n.Exchange(cmdGetViewKey, 0, 0, nil)
	if err != nil {
		return "", err
	}
	fmt.Printf("viewkey: %v\n", resp)
	return hex.EncodeToString(resp), nil
}

func (n *NanoS) GetOTAKey() (string, error) {
	resp, err := n.Exchange(cmdGetOTAKey, 0, 0, nil)
	if err != nil {
		return "", err
	}
	fmt.Printf("ota: %v\n", resp)
	return hex.EncodeToString(resp), nil
}

func (n *NanoS) GetValidatorKey() error {
	resp, err := n.Exchange(cmdGetValidatorKey, 0, 0, nil)
	if err != nil {
		return err
	}
	_ = resp
	return nil
}

func (n *NanoS) SwitchKey() error {
	buf := new(bytes.Buffer)

	bs, _ := hex.DecodeString("")
	buf.Write(bs)
	resp, err := n.Exchange(cmdSwitchKey, 0, 0, buf.Next(255))
	if err != nil {
		return err
	}
	_ = resp
	return nil
}

func (n *NanoS) GenKeyImage(coinPubkey string, encryptKm string) (string, error) {
	buf := new(bytes.Buffer)
	bs, err := hex.DecodeString(encryptKm)
	if err != nil {
		panic(err)
	}
	buf.Write(bs)
	bs1, err := hex.DecodeString(coinPubkey)
	if err != nil {
		panic(err)
	}
	buf.Write(bs1)

	resp, err := n.Exchange(cmdKeyImage, 0, 0, buf.Next(255))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(resp), nil
}

//This func contain a set of commands for ledger
func (n *NanoS) GenerateAlpha(alphaLength int) error {
	_, err := n.Exchange(cmdGenAlpha, byte(alphaLength), 0, nil)
	if err != nil {
		return err
	}
	return nil
}

func (n *NanoS) CalculateFirstC(Rpi [][]byte, PedComG []byte) ([]byte, error) {
	var result []byte
	buf := new(bytes.Buffer)
	for i := 0; i < len(Rpi)-1; i++ {
		buf.Reset()
		buf.Write(Rpi[i])
		resp, err := n.Exchange(cmdCalculateC, 0, byte(i), buf.Next(255))
		if err != nil {
			return nil, err
		}
		result = append(result, resp...)
		fmt.Println("CalculateFirstC", i, "success")
	}

	buf.Reset()
	buf.Write(PedComG)
	resp, err := n.Exchange(cmdCalculateC, 1, byte(len(Rpi)-1), buf.Next(255))
	if err != nil {
		return nil, err
	}
	result = append(result, resp...)
	return result, nil
}

func (n *NanoS) CalculateR(coinLength int, cPi []byte) ([][]byte, error) {
	buf := new(bytes.Buffer)
	new_rPi := make([][]byte, coinLength)
	for idx := 0; idx < coinLength; idx++ {
		buf.Reset()
		buf.Write(cPi)
		resp, err := n.Exchange(cmdCalculateR, 0, byte(idx), buf.Next(255))
		if err != nil {
			return nil, err
		}
		fmt.Println("CalculateR success", resp)
		new_rPi[idx] = resp
	}
	return new_rPi, nil
}

func (n *NanoS) GenCoinPrivateKey(coinsH [][]byte) error {
	buf := new(bytes.Buffer)
	for idx, coinH := range coinsH {
		buf.Reset()
		buf.Write(coinH)
		if idx == len(coinsH)-1 { //add sumRand privkey
			_, err := n.Exchange(cmdGenCoinPrivateKey, 1, byte(idx), buf.Next(255))
			if err != nil {
				return err
			}
		} else {
			_, err := n.Exchange(cmdGenCoinPrivateKey, 0, byte(idx), buf.Next(255))
			if err != nil {
				return err
			}
		}
		fmt.Println("GenCoinPrivateKey success", idx)
	}
	return nil
}

func (n *NanoS) SignSchnorr(pedRandom []byte, pedPrivate []byte, randomness []byte, message []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	if pedRandom == nil {
		pedRandom = make([]byte, 32)
	}
	buf.Write(pedRandom)
	buf.Write(pedPrivate)
	buf.Write(randomness)
	buf.Write(message)

	resp, err := n.Exchange(cmdSignSchnorr, 0, 0, buf.Next(255))
	if err != nil {
		return nil, err
	}
	fmt.Printf("sig: %s\n", hex.EncodeToString(resp))

	return resp, nil
}
