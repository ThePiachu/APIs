package BlockchainAPI


var RPCUser string
var RPCPassword string

var Guid string
var MainPassword string
var SecondPassword string

func SetupWallet(guid string, mainPassword string) {
	Guid = guid
	MainPassword = mainPassword
	RPCUser=guid
	RPCPassword=mainPassword
}

func SetupWalletWithSecondPassword(guid string, mainPassword string, secondPassword string) {
	Guid = guid
	MainPassword = mainPassword
	SecondPassword = secondPassword
	RPCUser=guid
	RPCPassword=mainPassword
}
