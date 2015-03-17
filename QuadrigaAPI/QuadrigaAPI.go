package QuadrigaAPI

import(
	"net/url"
	"appengine"
	"appengine/urlfetch"
	"io/ioutil"
	"encoding/json"
	"time"
	"net/http"
	"crypto/hmac"
	"crypto/sha256"
	"bytes"
	"github.com/ThePiachu/Go/mymath"
	"github.com/ThePiachu/Go/Log"
)

//https://www.quadrigacx.com/api_info


var APIKey string
var APISecret string
var ClientID string

func Setup(apiKey, apiSecret, clientID string) {
	APIKey = apiKey
	APISecret = apiSecret
	ClientID = clientID
}

func GetNonce() int64 {
	now:=time.Now()
	return now.UnixNano()/1000
}

func GetNonceStr() string {
	nonce:=GetNonce()
	return mymath.Int642Str(nonce)
}

func GenerateSignature(nonce string, clientID string, apiKey string, secretKey string) string {
	toEncode:=mymath.ASCII2Hex(nonce+clientID+apiKey)
	key:=mymath.ASCII2Hex(secretKey)
	md5:=mymath.ASCII2Hex(mymath.ToLower(mymath.Hex2Str(mymath.MD5(key))))

	hmacHash:=hmac.New(sha256.New, md5)
	hmacHash.Write(toEncode)

	answer := hmacHash.Sum(nil)

	return mymath.Hex2Str(answer)
}

func CallPost(c appengine.Context, callURL string, parameters map[string]string) (map[string]interface{}, error) {
	resp, err:=Call(c, callURL, parameters, "POST")
	return (resp).(map[string]interface{}), err
}
func CallPostArray(c appengine.Context, callURL string, parameters map[string]string) ([]interface{}, error) {
	resp, err:=Call(c, callURL, parameters, "POST")
	return (resp).([]interface{}), err
}
func CallGet(c appengine.Context, callURL string, parameters map[string]string) (map[string]interface{}, error) {
	resp, err:=Call(c, callURL, parameters, "GET")
	return (resp).(map[string]interface{}), err
}
func CallGetArray(c appengine.Context, callURL string, parameters map[string]string) ([]interface{}, error) {
	resp, err:=Call(c, callURL, parameters, "GET")
	return (resp).([]interface{}), err
}


func Call(c appengine.Context, callURL string, parameters map[string]string, method string) (interface{}, error) {
	TimeoutDuration, err := time.ParseDuration("60s");
	if err!=nil {
		Log.Infof(c, "API - Call - error 1 - %s", err.Error())
		return nil, err
	}
	tr := urlfetch.Transport{Context: c, Deadline: TimeoutDuration, AllowInvalidServerCertificate:true}
	client := http.Client{Transport: &tr}
	values := url.Values{}
	//turning input map into url.Values
	for k, v := range parameters {
		values.Set(k, v)
	}
	req, err:=http.NewRequest(method, callURL, bytes.NewBufferString(values.Encode()))
	if err!=nil {
		Log.Infof(c, "API - Call - error 2 - %s", err.Error())
		return nil, err
	}
	req.Header.Set("key", parameters["key"])
	req.Header.Set("signature", parameters["signature"])
	req.Header.Set("nonce", parameters["nonce"])


	Log.Infof(c, "req - %v", req)
	//sending request
	resp, err:=client.Do(req)
	if err != nil {
		Log.Errorf(c, "API post error: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	//reading response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Errorf(c, "API read error: could not read body: %s", err)
		return nil, err
	}
	var result interface{}
	//unmarshalling JSON response
	err = json.Unmarshal(body, &result)
	if err != nil {
		Log.Infof(c, "Unmarshal: %v", err)
		Log.Infof(c, "%s", body)
		return nil, err
	}
	return result, nil
}

//Public functions

func Ticker(c appengine.Context, book string) (map[string]interface{}, error) {
	params:=getEmptyMap()
	params["book"] = book
	return CallGet(c, "https://api.quadrigacx.com/v2/ticker", params)
}

func OrderBook(c appengine.Context) (map[string]interface{}, error) {
	return OrderBookWithParams(c, "", true)
}

func OrderBookWithParams(c appengine.Context, book string, group bool) (map[string]interface{}, error) {
	params:=getEmptyMap()
	if group {
		params["group"]="1"
	} else {
		params["group"]="0"
	}
	if book!="" {
		params["book"] = book
	}
	return CallGet(c, "https://api.quadrigacx.com/v2/order_book", params)
}

func Transactions(c appengine.Context) ([]interface{}, error) {
	return TransactionWithTimeFrame(c, "", "")
}

func TransactionWithTimeFrame(c appengine.Context, book string, time string) ([]interface{}, error) {
	params:=getEmptyMap()
	if time!="" {
		params["time"] = time
	}
	if book!="" {
		params["book"] = book
	}
	return CallGetArray(c, "https://api.quadrigacx.com/v2/transactions", params)
}





//Private functions

func AccountBalance(c appengine.Context) (map[string]interface{}, error) {
	params:=getBasicMap()
	return CallPost(c, "https://api.quadrigacx.com/v2/balance", params)
}

func UserTransactions(c appengine.Context) ([]interface{}, error) {
	return UserTransactionsFull(c, 0, 100, true, "")
}

func UserTransactionsFull(c appengine.Context, offset int64, limit int64, descending bool, book string) ([]interface{}, error) {
	params:=getBasicMap()
	params["offset"]=mymath.Int642Str(offset)
	params["limit "]=mymath.Int642Str(limit)
	if descending {
		params["sort"]="desc"
	} else {
		params["sort"]="asc"
	}
	if book!="" {
		params["book"] = book
	}
	return CallPostArray(c, "https://api.quadrigacx.com/v2/user_transactions", params)
}

func OpenOrders(c appengine.Context, book string) ([]interface{}, error) {
	params:=getBasicMap()
	if book!="" {
		params["book"] = book
	}
	return CallPostArray(c, "https://api.quadrigacx.com/v2/open_orders", params)
}

func LookupOrder(c appengine.Context, id string) (map[string]interface{}, error) {
	params:=getBasicMap()
	params["id"] = id
	return CallPost(c, "https://api.quadrigacx.com/v2/lookup_order", params)
}

func CancelOrder(c appengine.Context, id string) (map[string]interface{}, error) {
	params:=getBasicMap()
	params["id"] = id
	return CallPost(c, "https://api.quadrigacx.com/v2/cancel_order", params)
}

func BuyLimitOrder(c appengine.Context, amount float64, price float64, book string) (map[string]interface{}, error) {
	params:=getBasicMap()
	params["amount"]=mymath.Float642Str(mymath.TrimFloatToXDecimalDigits(amount, 8))
	params["price"]=mymath.Float642Str(price)
	if book!="" {
		params["book"] = book
	}
	return CallPost(c, "https://api.quadrigacx.com/v2/buy", params)
}

func BuyMarketOrder(c appengine.Context, amount float64, book string) (map[string]interface{}, error) {
	params:=getBasicMap()
	params["amount"]=mymath.Float642Str(mymath.TrimFloatToXDecimalDigits(amount, 8))
	if book!="" {
		params["book"] = book
	}
	return CallPost(c, "https://api.quadrigacx.com/v2/buy", params)
}

func SellLimitOrder(c appengine.Context, amount float64, price float64, book string) (map[string]interface{}, error) {
	params:=getBasicMap()
	params["amount"]=mymath.Float642Str(mymath.TrimFloatToXDecimalDigits(amount, 8))
	params["price"]=mymath.Float642Str(price)
	if book!="" {
		params["book"] = book
	}
	return CallPost(c, "https://api.quadrigacx.com/v2/sell", params)
}

func SellMarketOrder(c appengine.Context, amount float64, book string) (map[string]interface{}, error) {
	params:=getBasicMap()
	params["amount"]=mymath.Float642Str(mymath.TrimFloatToXDecimalDigits(amount, 8))
	if book!="" {
		params["book"] = book
	}
	return CallPost(c, "https://api.quadrigacx.com/v2/sell", params)
}

func BitcoinDepositAddress(c appengine.Context) (string, error) {
	params:=getBasicMap()
	resp, err:=Call(c, "https://api.quadrigacx.com/v2/bitcoin_deposit_address", params, "POST")
	if err!=nil {
		return "", err
	}
	return resp.(string), err
}

func BitcoinWithdrawal(c appengine.Context, amount float64, address string) (string, error) {
	params:=getBasicMap()
	params["amount"]=mymath.Float642Str(mymath.TrimFloatToXDecimalDigits(amount, 8))
	params["address"]=address
	resp, err:=Call(c, "https://api.quadrigacx.com/v2/bitcoin_withdrawal", params, "POST")
	if err!=nil {
		return "", err
	}
	return resp.(string), err
}

func getEmptyMap() map[string]string {
	answer := make(map[string]string)
	return answer
}
func getBasicMap() map[string]string {
	answer := make(map[string]string)
	nonce:=GetNonceStr()
	answer["key"]=APIKey
	answer["signature"]=GenerateSignature(nonce, ClientID, APIKey, APISecret)
	answer["nonce"]=nonce
	return answer
}
