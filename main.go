package main

import (
	"fmt"
	"log"

	"github.com/0xkumi/incognito-dev-framework/account"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/privacy"
	"github.com/incognitochain/incognito-chain/privacy/operation"
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
	privUsage         = ``
	viewKeyUsage      = ``
	importPrivUsage   = ``
	getOTAKeyUsage    = ``
	getValidatorUsage = ``
	genKeyImageUsage  = ``
	genRingSigUsage   = ``
	signSchnorrUsage  = ``
	trustHostUsage    = ``

	listAccountUsage   = ``
	getBalanceUsage    = ``
	updateBalanceUsage = ``
	createTxUsage      = ``
)

func main() {
	log.SetFlags(0)
	rootCmd := flagg.Root
	rootCmd.Usage = flagg.SimpleUsage(rootCmd, rootUsage)

	versionCmd := flagg.New("version", versionUsage)
	addrCmd := flagg.New("addr", addrUsage)
	privCmd := flagg.New("priv", privUsage)
	getViewKeyCmd := flagg.New("view", viewKeyUsage)
	importPrivateKeyCmd := flagg.New("importpriv", importPrivUsage)
	getOTAKeyCmd := flagg.New("ota", getOTAKeyUsage)
	getValidatorCmd := flagg.New("getvalidator", getValidatorUsage)
	genKeyImageCmd := flagg.New("genkeyimage", genKeyImageUsage)
	signSchnorrCmd := flagg.New("signschnorr", signSchnorrUsage)
	trustHostCmd := flagg.New("trust", trustHostUsage)
	listAccountCmd := flagg.New("listaccount", listAccountUsage)
	getBalanceCmd := flagg.New("getbalance", getBalanceUsage)
	updateBalanceCmd := flagg.New("updatebalance", updateBalanceUsage)
	createTxCmd := flagg.New("createtx", createTxUsage)
	cmd := flagg.Parse(flagg.Tree{
		Cmd: rootCmd,
		Sub: []flagg.Tree{
			{Cmd: versionCmd},
			{Cmd: addrCmd},
			{Cmd: privCmd},
			{Cmd: getViewKeyCmd},
			{Cmd: importPrivateKeyCmd},
			{Cmd: genKeyImageCmd},
			{Cmd: getValidatorCmd},
			{Cmd: getOTAKeyCmd},
			{Cmd: signSchnorrCmd},
			{Cmd: trustHostCmd},
			{Cmd: listAccountCmd},
			{Cmd: getBalanceCmd},
			{Cmd: updateBalanceCmd},
			{Cmd: createTxCmd},
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
	case addrCmd:
		if len(args) != 1 {
			addrCmd.Usage()
			return
		}
		addr, err := nanos.GetAddress(parseIndex(args[0]))
		if err != nil {
			log.Fatalln("Couldn't get address:", err)
		}
		fmt.Println(addr)
	case privCmd:
		if len(args) != 1 {
			privCmd.Usage()
			return
		}
		priv, err := nanos.GetPrivateKey(parseIndex(args[0]))
		if err != nil {
			log.Fatalln("Couldn't get address:", err)
		}
		fmt.Println(priv)
	case getViewKeyCmd:
		err := nanos.GetViewKey()
		if err != nil {
			log.Fatalln(err)
		}
	case importPrivateKeyCmd:
		err := nanos.ImportPrivateKey()
		if err != nil {
			log.Fatalln(err)
		}
	case getOTAKeyCmd:
		err := nanos.GetValidatorKey()
		if err != nil {
			log.Fatalln(err)
		}
	case getValidatorCmd:
		err := nanos.GetValidatorKey()
		if err != nil {
			log.Fatalln(err)
		}
	case signSchnorrCmd:
		r := new(privacy.Scalar).FromUint64(1)
		pedRandom := operation.PedCom.G[operation.PedersenRandomnessIndex].GetKey()
		pedPrivate := operation.PedCom.G[operation.PedersenPrivateKeyIndex].GetKey()
		message := "testasfdgtestasfdgfgfdgfgtestasfdgfgfdgtestasfdgfgfdgtestasfdgfgfdgfdg"
		hash := common.HashH([]byte(message))
		resp, err := nanos.SignSchnorr(pedRandom[:], pedPrivate[:], r.ToBytesS(), hash.Bytes())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("resp", resp)

		acc0, _ := account.NewAccountFromPrivatekey("111111bgk2j6vZQvzq8tkonDLLXEvLkMwBMn5BoLXLpf631boJnPDGEQMGvA1pRfT71Crr7MM2ShvpkxCBWBL2icG22cXSpcKybKCQmaxa")
		//verify
		verifyKey := new(privacy.SchnorrPublicKey)
		metaSigPublicKey, err := new(privacy.Point).FromBytesS(acc0.Keyset.PaymentAddress.Pk)
		if err != nil {
			panic(err)
		}
		verifyKey.Set(metaSigPublicKey)

		signature := new(privacy.SchnSignature)
		if err := signature.SetBytes(resp); err != nil {
			panic(err)
		}
		fmt.Println("verify sig", verifyKey.Verify(signature, hash.Bytes()))

	case trustHostCmd:
		err := nanos.TrustHost()
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
		result, err := updateBalanceFlow(account)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(result)
	case createTxCmd:
		txjsonLink := args[0]
		result, err := requestCreateTx(txjsonLink)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(result)
	}
}
