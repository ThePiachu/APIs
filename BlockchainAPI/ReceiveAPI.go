package BlockchainAPI

import(
	"appengine"
)

//http://blockchain.info/api/api_receive


var APIURL string = "https://blockchain.info/api/"

func CallAPI(c appengine.Context, guid string, function string, parameters map[string]string) map[string]interface{} {
	return Call(c, APIURL, guid, function, parameters)
}

func CreateReceivingAddress(c appengine.Context, address string, callbackURL string) map[string]interface{} {
	params := make(map[string]string)
	params["address"] = address;
	params["method"] = "create";
	params["callback"] = callbackURL;
	return CallAPI(c, Guid, "receive", params)
}