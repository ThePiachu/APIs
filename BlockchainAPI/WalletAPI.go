package BlockchainAPI

import(
	"mymath"
	"net/url"
	"appengine"
	"appengine/urlfetch"
	"io/ioutil"
	"encoding/json"
	"time"
	"net/http"
	"github.com/ThePiachu/Go/Log"
)

//http://blockchain.info/api/blockchain_wallet_api

var MerchantURL string = "https://blockchain.info/merchant/"

func CallMerchant(c appengine.Context, guid string, function string, parameters map[string]string) map[string]interface{} {
	return Call(c, MerchantURL, guid, function, parameters)
}

func Call(c appengine.Context, callURL string, guid string, function string, parameters map[string]string) map[string]interface{} {
	TimeoutDuration, err := time.ParseDuration("60s");
	if err!=nil {
		Log.Infof(c, "BlockchainAPI - Call - error - %s", err.Error())
	}
    tr := urlfetch.Transport{Context: c, Deadline: TimeoutDuration}
	client := http.Client{Transport: &tr}
	values := url.Values{}
	//turning input map into url.Values
	for k, v := range parameters {
		values.Set(k, v)
	}
	//sending request
	resp, err:=client.PostForm(callURL+guid+"/"+function, values)
	if err != nil {
		Log.Errorf(c, "BlockchainAPI post error: %s", err)
		return nil
	}
	defer resp.Body.Close()
	//reading response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Errorf(c, "BlockchainAPI read error: could not read body: %s", err)
		return nil
	}
	result := make(map[string]interface{})
	//unmarshalling JSON response
    err = json.Unmarshal(body, &result)
    if err != nil {
        Log.Infof(c, "Unmarshal: %v", err)
        Log.Infof(c, "%s", body)
    	return nil
    }
    return result
}

func Payment(c appengine.Context, to string, amount int64) map[string]interface{} {
	params := getBasicMap()
	params["to"] = to
	params["amount"] = mymath.Int642Str(amount)
	
	return CallMerchant(c, Guid, "payment", params)
}

func Payment2(c appengine.Context, to string, amount int64, from string) map[string]interface{} {
	params := getBasicMap()
	params["to"] = to
	params["amount"] = mymath.Int642Str(amount)
	params["from"] = from
	
	return CallMerchant(c, Guid, "payment", params)
}

func Payment3(c appengine.Context, to string, amount int64, from string, shared bool) map[string]interface{} {
	params := getBasicMap()
	params["to"] = to
	params["amount"] = mymath.Int642Str(amount)
	params["from"] = from
	if shared {
		params["shared"] = "true"
	} else {
		params["shared"] = "false"
	}
	
	return CallMerchant(c, Guid, "payment", params)
}

func Payment4(c appengine.Context, to string, amount int64, from string, shared bool, fee int64) map[string]interface{} {
	params := getBasicMap()
	params["to"] = to
	params["amount"] = mymath.Int642Str(amount)
	params["from"] = from
	if shared {
		params["shared"] = "true"
	} else {
		params["shared"] = "false"
	}
	params["fee"] = mymath.Int642Str(fee)
	
	return CallMerchant(c, Guid, "payment", params)
}

func Payment5(c appengine.Context, to string, amount int64, from string, shared bool, fee int64, note string) map[string]interface{} {
	params := getBasicMap()
	params["to"] = to
	params["amount"] = mymath.Int642Str(amount)
	params["from"] = from
	if shared {
		params["shared"] = "true"
	} else {
		params["shared"] = "false"
	}
	params["fee"] = mymath.Int642Str(fee)
	params["note"] = note
	return CallMerchant(c, Guid, "payment", params)
}

func SendMany(c appengine.Context, recipients map[string]int64) map[string]interface{} {
	params := getBasicMap()
	params["recipients"] = mymath.MapStringInt64ToString(recipients)
	return CallMerchant(c, Guid, "sendmany", params)
}

func SendMany2(c appengine.Context, recipients map[string]int64, from string) map[string]interface{} {
	params := getBasicMap()
	params["recipients"] = mymath.MapStringInt64ToString(recipients)
	params["from"] = from
	return CallMerchant(c, Guid, "sendmany", params)
}

func SendMany3(c appengine.Context, recipients map[string]int64, from string, shared bool) map[string]interface{} {
	params := getBasicMap()
	params["recipients"] = mymath.MapStringInt64ToString(recipients)
	params["from"] = from
	if shared {
		params["shared"] = "true"
	} else {
		params["shared"] = "false"
	}
	return CallMerchant(c, Guid, "sendmany", params)
}

func SendMany4(c appengine.Context, recipients map[string]int64, from string, shared bool, fee int64) map[string]interface{} {
	params := getBasicMap()
	params["recipients"] = mymath.MapStringInt64ToString(recipients)
	params["from"] = from
	if shared {
		params["shared"] = "true"
	} else {
		params["shared"] = "false"
	}
	params["fee"] = mymath.Int642Str(fee)
	return CallMerchant(c, Guid, "sendmany", params)
}

func SendMany5(c appengine.Context, recipients map[string]int64, from string, shared bool, fee int64, note string) map[string]interface{} {
	params := getBasicMap()
	params["recipients"] = mymath.MapStringInt64ToString(recipients)
	params["from"] = from
	if shared {
		params["shared"] = "true"
	} else {
		params["shared"] = "false"
	}
	params["fee"] = mymath.Int642Str(fee)
	params["note"] = note
	return CallMerchant(c, Guid, "sendmany", params)
}

func Balance(c appengine.Context) map[string]interface{} {
	params := getBasicMap()
	return CallMerchant(c, Guid, "balance", params)
}

func ListAddresses(c appengine.Context) map[string]interface{} {
	params := getBasicMap()
	return CallMerchant(c, Guid, "list", params)
}

func ListAddressesWithConfirmations(c appengine.Context, confirmations int) map[string]interface{} {
	params := getBasicMap()
	params["confirmations"] = mymath.Int2Str(confirmations)
	return CallMerchant(c, Guid, "list", params)
}

func GetAddressBalance(c appengine.Context, address string, confirmations int) map[string]interface{} {
	params := getBasicMap()
	params["address"] = address
	params["confirmations"] = mymath.Int2Str(confirmations)
	return CallMerchant(c, Guid, "address_balance", params)
}

func GenerateNewAddress(c appengine.Context) map[string]interface{} {
	params := getBasicMap()
	return CallMerchant(c, Guid, "new_address", params)
}

func GenerateNewAddressWithLabel(c appengine.Context, label string) map[string]interface{} {
	params := getBasicMap()
	params["label"]=label
	return CallMerchant(c, Guid, "new_address", params)
}

func ArchiveAddress(c appengine.Context, address string) map[string]interface{} {
	params := getBasicMap()
	params["address"] = address;
	return CallMerchant(c, Guid, "archive_address", params)
}

func UnarchiveAddress(c appengine.Context, address string) map[string]interface{} {
	params := getBasicMap()
	params["address"] = address;
	return CallMerchant(c, Guid, "unarchive_address", params)
}

func AutoConsolidate(c appengine.Context, days int) map[string]interface{} {
	params := getBasicMap()
	params["days"]=mymath.Int2Str(days)
	return CallMerchant(c, Guid, "auto_consolidate", params)
}

func getBasicMap() map[string]string {
	answer := make(map[string]string)
	answer["password"] = MainPassword
	if SecondPassword != "" {
		answer["second_password"] = SecondPassword
	}
	return answer
}
