module main

go 1.13

require (
	github.com/0xkumi/incognito-dev-framework v0.0.0-20210127050012-71404660c31d
	github.com/gorilla/websocket v1.4.2
	github.com/incognitochain/incognito-chain v0.0.0-20201229061112-f61b51f89ddd
	github.com/zondax/hid v0.9.0
	golang.org/x/text v0.3.2
	lukechampine.com/flagg v1.1.1
)

// replace github.com/incognitochain/incognito-chain => /Users/truonglamchau/go/src/github.com/incognitochain/incognito-chain

replace github.com/incognitochain/incognito-chain => /home/lam/go/src/github.com/incognitochain/incognito-chain
