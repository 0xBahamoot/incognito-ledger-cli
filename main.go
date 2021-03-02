package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/0xkumi/incognito-dev-framework/account"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/privacy"
	"github.com/incognitochain/incognito-chain/privacy/operation"

	// "github.com/tendermint/tendermint/types/time"
	"lukechampine.com/flagg"
)

const (
	rootUsage = `Usage:
    incognitoledger [flags] [action]

Actions:
    addr            generate an address
    pubkey          generate a pubkey
    hash            sign a trusted hash
    txn             sign a transaction
`

	versionUsage = `Usage:
	incognitoledger version

Prints the version of the incognitoledger binary, as well as the version reported by
the Incognito Ledger Nano S app (if available).
`
	addrUsage = `Usage:
	incognitoledger addr [key index]

Generates an address using the public key with the specified index.
`
	trustHostUsage     = ``
	viewKeyUsage       = ``
	getOTAKeyUsage     = ``
	getValidatorUsage  = ``
	listAccountUsage   = ``
	getBalanceUsage    = ``
	updateBalanceUsage = ``
	createTxUsage      = ``
	importAccountUsage = ``

	privUsage        = ``
	importPrivUsage  = ``
	genKeyImageUsage = ``
	signSchnorrUsage = ``
	benchmarkUsage   = ``
)

func main() {
	log.SetFlags(0)
	rootCmd := flagg.Root
	rootCmd.Usage = flagg.SimpleUsage(rootCmd, rootUsage)

	versionCmd := flagg.New("version", versionUsage)
	addrCmd := flagg.New("addr", addrUsage)
	getViewKeyCmd := flagg.New("view", viewKeyUsage)
	getOTAKeyCmd := flagg.New("ota", getOTAKeyUsage)
	getValidatorCmd := flagg.New("getvalidator", getValidatorUsage)
	trustHostCmd := flagg.New("trust", trustHostUsage)
	listAccountCmd := flagg.New("listaccount", listAccountUsage)
	getBalanceCmd := flagg.New("getbalance", getBalanceUsage)
	updateBalanceCmd := flagg.New("updatebalance", updateBalanceUsage)
	createTxCmd := flagg.New("createtx", createTxUsage)
	importAccountCmd := flagg.New("importacc", importAccountUsage)

	benchmarkCmd := flagg.New("benchmark", benchmarkUsage)
	privCmd := flagg.New("priv", privUsage)
	importPrivateKeyCmd := flagg.New("importpriv", importPrivUsage)
	genKeyImageCmd := flagg.New("genkeyimage", genKeyImageUsage)
	signSchnorrCmd := flagg.New("signschnorr", signSchnorrUsage)

	cmd := flagg.Parse(flagg.Tree{
		Cmd: rootCmd,
		Sub: []flagg.Tree{
			// actual cmd
			{Cmd: trustHostCmd},
			{Cmd: versionCmd},
			{Cmd: addrCmd},
			{Cmd: getViewKeyCmd},
			{Cmd: getValidatorCmd},
			{Cmd: getOTAKeyCmd},
			{Cmd: listAccountCmd},
			{Cmd: getBalanceCmd},
			{Cmd: updateBalanceCmd},
			{Cmd: createTxCmd},
			{Cmd: importAccountCmd},

			// dev ledger cmd
			{Cmd: privCmd},
			{Cmd: importPrivateKeyCmd},
			{Cmd: genKeyImageCmd},
			{Cmd: signSchnorrCmd},
			{Cmd: benchmarkCmd},
		},
	})
	args := cmd.Args()
	fmt.Println("args", args)
	readConfig()
	var nanos *NanoS
	if cmd != rootCmd && cmd != versionCmd && cmd != listAccountCmd && cmd != getBalanceCmd && cmd != updateBalanceCmd && cmd != createTxCmd {
		var err error
		nanos, err = OpenNanoS()
		if err != nil {
			log.Println("This cmd require connected to ledger device")
			log.Fatalln("Couldn't open device:", err)
		}
	}

	switch cmd {
	case rootCmd:
		if len(args) != 0 {
			rootCmd.Usage()
			return
		}
		fallthrough
	case versionCmd:
		// try to get Nano S app version
		var appVersion string
		nanos, err := OpenNanoS()
		if err != nil {
			appVersion = "(could not connect to Nano S)"
		} else if appVersion, err = nanos.GetVersion(); err != nil {
			appVersion = "(could not read version from Nano S: " + err.Error() + ")"
		}

		fmt.Printf("CLI version: %s\n", CLI_version)
		fmt.Println("Nano S app version:", appVersion)
		fmt.Printf("CoinDaemon version: %s\n", CLI_version)
	case trustHostCmd:
		err := nanos.TrustHost()
		if err != nil {
			log.Fatalln(err)
		}
	case addrCmd:
		addr, err := nanos.GetAddress()
		if err != nil {
			log.Fatalln("Couldn't get address:", err)
		}
		fmt.Println(addr)
	case getViewKeyCmd:
		_, err := nanos.GetViewKey()
		if err != nil {
			log.Fatalln(err)
		}
	case getOTAKeyCmd:
		_, err := nanos.GetOTAKey()
		if err != nil {
			log.Fatalln(err)
		}
	case getValidatorCmd:
		err := nanos.GetValidatorKey()
		if err != nil {
			log.Fatalln(err)
		}
	case listAccountCmd:
		result, err := getAccountList()
		if err != nil {
			log.Fatalln(err)
		}
		for name, addr := range result {
			fmt.Printf("%s: %s", name, addr)
		}
	case getBalanceCmd:
		account := args[0]
		result, err := getAccountBalance(account)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(result)
	case updateBalanceCmd:
		account := args[0]
		result, err := requestUpdateBalance(account)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(result)
	case createTxCmd:
		t := time.Now()
		txjsonLink := args[0]
		result, err := requestCreateTx(txjsonLink)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(result)
		fmt.Println("time:", time.Since(t))
	case importAccountCmd:
		err := nanos.TrustHost()
		if err != nil {
			log.Fatalln(err)
		}
		accountName := args[0]
		beaconHeight := uint64(0)
		if len(args) == 2 {
			var err error
			beaconHeight, err = strconv.ParseUint(args[1], 0, 64)
			if err != nil {
				panic(err)
			}
		}
		viewKey, err := nanos.GetViewKey()
		if err != nil {
			panic(err)
		}
		otaKey, err := nanos.GetOTAKey()
		if err != nil {
			panic(err)
		}
		addr, err := nanos.GetAddress()
		if err != nil {
			log.Fatalln("Couldn't get address:", err)
		}
		err = importAccount(accountName, addr, otaKey, viewKey, beaconHeight)
		if err != nil {
			log.Fatalln(err)
		}

	//for dev-use only
	case importPrivateKeyCmd:
		_, err := nanos.ImportPrivateKey()
		if err != nil {
			log.Fatalln(err)
		}
	case privCmd:
		if len(args) != 1 {
			privCmd.Usage()
			return
		}
		priv, err := nanos.GetPrivateKey()
		if err != nil {
			log.Fatalln("Couldn't get address:", err)
		}
		fmt.Println(priv)
	case genKeyImageCmd:
		err := nanos.TrustHost()
		if err != nil {
			log.Fatalln(err)
		}
		result, err := nanos.GenKeyImage("17fd6aff8fecd18243af1a83dab0e47ca5fafec256ba497b3136a6b3f68eecb1", "c4541151e39bb43e7b00ad6a1d999d609f5939ca622a9db7b7391c5190eea909")
		if err != nil {
			panic(err)
		}
		_ = result
	case signSchnorrCmd:
		acc0, _ := account.NewAccountFromPrivatekey("111111bgk2j6vZQvzq8tkonDLLXEvLkMwBMn5BoLXLpf631boJnPDGEQMGvA1pRfT71Crr7MM2ShvpkxCBWBL2icG22cXSpcKybKCQmaxa")

		t := time.Now()
		r := new(privacy.Scalar).FromUint64(1)
		pedRandom := operation.PedCom.G[operation.PedersenRandomnessIndex].GetKey()
		pedPrivate := operation.PedCom.G[operation.PedersenPrivateKeyIndex].GetKey()
		message := "testasfdgtestasfdgfgfdgfgtestasfdgfgfdgtestasfdgfgfdgtestasfdgfgfdgfdg"
		hash := common.HashH([]byte(message))
		resp, err := nanos.SignSchnorr(pedRandom[:], pedPrivate[:], r.ToBytesS(), hash.Bytes())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("resp", resp, len(resp))
		fmt.Println("signSchnorr:", time.Since(t))

		//verify
		verifyKey := new(privacy.SchnorrPublicKey)
		metaSigPublicKey, err := new(privacy.Point).FromBytesS(acc0.Keyset.PaymentAddress.Pk)
		if err != nil {
			panic(err)
		}
		metaSigPublicKey.Add(metaSigPublicKey, new(operation.Point).ScalarMult(operation.PedCom.G[operation.PedersenRandomnessIndex], r))
		verifyKey.Set(metaSigPublicKey)

		signature := new(privacy.SchnSignature)
		if err := signature.SetBytes(resp); err != nil {
			panic(err)
		}
		fmt.Println("verify sig", verifyKey.Verify(signature, hash.Bytes()))
	case benchmarkCmd:
		err := nanos.TrustHost()
		if err != nil {
			log.Fatalln(err)
		}
		t := time.Now()
		for i := 0; i < 20; i++ {
			result, err := nanos.GenKeyImage("17fd6aff8fecd18243af1a83dab0e47ca5fafec256ba497b3136a6b3f68eecb1", "c4541151e39bb43e7b00ad6a1d999d609f5939ca622a9db7b7391c5190eea909")
			if err != nil {
				panic(err)
			}
			_ = result
		}
		fmt.Println("genKeyImage 20:", time.Since(t))

		t = time.Now()
		r := new(privacy.Scalar).FromUint64(1)
		pedRandom := operation.PedCom.G[operation.PedersenRandomnessIndex].GetKey()
		pedPrivate := operation.PedCom.G[operation.PedersenPrivateKeyIndex].GetKey()
		message := "testasfdgtestasfdgfgfdgfgtestasfdgfgfdgtestasfdgfgfdgtestasfdgfgfdgfdg"
		hash := common.HashH([]byte(message))
		resp, err := nanos.SignSchnorr(pedRandom[:], pedPrivate[:], r.ToBytesS(), hash.Bytes())
		if err != nil {
			log.Fatalln(err)
		}
		_ = resp
		fmt.Println("signSchnorr:", time.Since(t))
	}
}
