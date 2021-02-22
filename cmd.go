package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/incognitochain/incognito-chain/privacy/operation"
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

func (n *NanoS) GetAddress(index uint32) (addr string, err error) {
	encIndex := make([]byte, 4)
	binary.LittleEndian.PutUint32(encIndex, index)

	resp, err := n.Exchange(cmdGetAddress, 0, 0, nil)
	if err != nil {
		return
	}
	fmt.Printf("address %v\n", string(resp[:]))
	addr = string(resp[:])
	return
}

func (n *NanoS) GetPrivateKey(index uint32) (priv string, err error) {
	encIndex := make([]byte, 4)
	binary.LittleEndian.PutUint32(encIndex, index)

	resp, err := n.Exchange(cmdGetPrivateKey, 0, p2DisplayAddress, encIndex)
	if err != nil {
		return
	}
	fmt.Println("privatekey:", string(resp))
	priv = string(resp)
	return
}

func (n *NanoS) GetViewKey() error {
	resp, err := n.Exchange(cmdGetViewKey, 0, 0, nil)
	if err != nil {
		return err
	}
	fmt.Printf("viewkey: %v", resp)
	// fmt.Println("viewkey:", string(resp[:]))
	return nil
}

func (n *NanoS) GetOTAKey() error {
	resp, err := n.Exchange(cmdGetOTAKey, 0, 0, nil)
	if err != nil {
		return err
	}
	fmt.Printf("ota: %v\n", resp)
	// fmt.Println("viewkey:", string(resp[:]))
	return nil
}

func (n *NanoS) ImportPrivateKey() error {
	buf := new(bytes.Buffer)
	// 000100000020812566598706f6f772fa0ec67e5efaac12c85a64b730518077a432fd3cb97a8c20063632b2a159e45002394460aee02de54d2b8926d236f45be2e077dcc81d0d04

	bs, _ := hex.DecodeString("00000000000214666ccc56b88d4d8d3f5fae61f1f06d9620327fe259157272016dfe54ef6fef20a408be78955356a9d3aef1729c6d83d32f91ea84cf21a974d2d9d791d71e1c06")
	buf.Write(bs)
	// buf.WriteString("111111bgk2j6vZQvzq8tkonDLLXEvLkMwBMn5BoLXLpf631boJnJ1UgJnLBzXe4qSMXGJAKw1LdKmfWZDNkhd24gkb2oqbs4q9UgjJZDvq")

	resp, err := n.Exchange(cmdImportPrivateKey, 0, 0, buf.Next(255))
	if err != nil {
		return err
	}
	_ = resp
	return nil
}

func (n *NanoS) GenKeyImage(coinPubkey string, encryptKm string) (string, error) {
	buf := new(bytes.Buffer)
	// bs, _ := hex.DecodeString("c4541151e39bb43e7b00ad6a1d999d609f5939ca622a9db7b7391c5190eea909")
	bs, err := hex.DecodeString(encryptKm)
	if err != nil {
		panic(err)
	}
	buf.Write(bs)
	// bs1, _ := hex.DecodeString("17fd6aff8fecd18243af1a83dab0e47ca5fafec256ba497b3136a6b3f68eecb1")
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

// var tmpAlpha = []*operation.Scalar{}

//This func contain a set of commands for ledger
func (n *NanoS) GenerateAlpha(alphaLength int) error {
	// tmpAlpha = []*operation.Scalar{}
	resp, err := n.Exchange(cmdGenAlpha, byte(alphaLength), 0, nil)
	if err != nil {
		return err
	}
	for i := 0; i < alphaLength; i++ {
		alpha := operation.Scalar{}
		alpha.FromBytesS(resp[i*32 : (i+1)*32])
		fmt.Println("alpha", alpha.ToBytesS())
		// tmpAlpha = append(tmpAlpha, &alpha)
	}
	// buf := new(bytes.Buffer)
	// for i := 0; i < alphaLength; i += 1 {
	// 	alpha := operation.RandomScalar()
	// 	tmpAlpha = append(tmpAlpha, alpha)
	// 	buf.Write(alpha.ToBytesS())
	// }
	// _, err := n.Exchange(cmdGenAlpha, byte(alphaLength), 0, buf.Next(255))
	// if err != nil {
	// 	return err
	// }
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
func (n *NanoS) GetValidatorKey() error {
	resp, err := n.Exchange(cmdGetValidatorKey, 0, 0, nil)
	if err != nil {
		return err
	}
	_ = resp
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

func (n *NanoS) TrustHost() error {
	resp, err := n.Exchange(cmdTrustHost, 0, 0, nil)
	if err != nil {
		return err
	}
	_ = resp
	return nil
}

func updateBalanceFlow(account string) (int, error) {
	nanos, err := OpenNanoS()
	if err != nil {
		log.Println("This cmd require connected to ledger device")
		log.Fatalln("Couldn't open device:", err)
	}
	var coinUpdated int
	keyimages, err := getEncryptKeyImages(account)
	if err != nil {
		return 0, err
	}
	fmt.Println(keyimages)

	err = nanos.TrustHost()
	if err != nil {
		return 0, err
	}
	decryptedKeyimages := make(map[string]map[string]string)
	tokenIDs := []string{}
	for tokenID, coinList := range keyimages {
		tokenIDs = append(tokenIDs, tokenID)
		decryptedKeyimages[tokenID] = make(map[string]string)
		for coinPk, km := range coinList {
			dekm, err := nanos.GenKeyImage(coinPk, km)
			if err != nil {
				panic(err)
			}
			decryptedKeyimages[tokenID][coinPk] = dekm
			fmt.Println("decryptedKeyimages[coinPk]", coinPk, dekm)
		}
	}

	for _, tokenID := range tokenIDs {
		e := submitKeyimages(tokenID, "testacc", decryptedKeyimages[tokenID])
		if err != nil {
			panic(e)
		}
	}

	return coinUpdated, nil
}
