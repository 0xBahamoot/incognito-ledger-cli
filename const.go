package main

const (
	CLI_version           = "0.5.5"
	DefaultCoinDaemonAddr = "127.0.0.1:9000"
)

var (
	compatibleLedgerVersions = []int{}
	compatibleDaemonVersions = []int{}
)

const (
	cmdGetVersion      = 0x01
	cmdGetAddress      = 0x02
	cmdGetViewKey      = 0x03
	cmdGetPrivateKey   = 0x04
	cmdSwitchKey       = 0x05
	cmdGetOTAKey       = 0x06
	cmdGetValidatorKey = 0x07
	cmdKeyImage        = 0x10
	// gen ring sig cmds set
	cmdGenAlpha          = 0x21
	cmdCalculateC        = 0x22
	cmdCalculateR        = 0x23
	cmdGenCoinPrivateKey = 0x24

	cmdSignSchnorr = 0x40
	cmdTrustHost   = 0x60

	p1First = 0x00
	p1More  = 0x80

	p2DisplayAddress = 0x00
	p2DisplayPubkey  = 0x01
	p2DisplayHash    = 0x00
	p2SignHash       = 0x01

	cla = 0xE0
)
