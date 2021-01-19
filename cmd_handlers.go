package main

import (
	"bytes"
	"encoding/binary"
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

func (n *NanoS) GetAddress(index uint32) (addr string, err error) {
	encIndex := make([]byte, 4)
	binary.LittleEndian.PutUint32(encIndex, index)

	resp, err := n.Exchange(cmdGetAddress, 0, 0, nil)
	if err != nil {
		return
	}
	fmt.Printf("address %v\n", resp[:])
	fmt.Println("\n", string(resp[:]))
	// fmt.Println("address:", string(resp))
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
	// err = addr.LoadString(string(resp[32:]))
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

func (n *NanoS) GenKeyImage() error {
	buf := new(bytes.Buffer)
	bs, _ := hex.DecodeString("c4541151e39bb43e7b00ad6a1d999d609f5939ca622a9db7b7391c5190eea909")
	buf.Write(bs)
	bs1, _ := hex.DecodeString("17fd6aff8fecd18243af1a83dab0e47ca5fafec256ba497b3136a6b3f68eecb1")
	buf.Write(bs1)

	resp, err := n.Exchange(cmdKeyImage, 0, 0, buf.Next(255))
	if err != nil {
		return err
	}
	fmt.Println("I", hex.EncodeToString(resp))
	return nil
}

func (n *NanoS) GenRingSig() error {
	resp, err := n.Exchange(cmdGenRingSig, 0, 0, nil)
	if err != nil {
		return err
	}
	_ = resp
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

func (n *NanoS) SignMetadata() error {
	resp, err := n.Exchange(cmdSignMetaData, 0, 0, nil)
	if err != nil {
		return err
	}
	fmt.Printf("sig: %s\n", hex.EncodeToString(resp))
	return nil
}

func (n *NanoS) TrustDevice() error {
	resp, err := n.Exchange(cmdTrustDevice, 0, 0, nil)
	if err != nil {
		return err
	}
	_ = resp
	return nil
}

func (n *NanoS) CreateTx() error {

	return nil
}
