package BitstampAPI

import(
	"net/url"
	"appengine"
	"appengine/urlfetch"
	"io/ioutil"
	"encoding/json"
	"time"
	"net/http"
	"mymath"
	"crypto/hmac"
	"crypto/sha256"
	"bytes"
	"github.com/ThePiachu/Go/Log"
)

//https://www.bitstamp.net/api/


var BitstampKey string
var BitstampSecret string
var BitstampClientID string

func Setup(apiKey, apiSecret, clientID string) {
	BitstampKey = apiKey
	BitstampSecret = apiSecret
	BitstampClientID = clientID
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
	
	hmacHash:=hmac.New(sha256.New, key)
	hmacHash.Write(toEncode)
	
	answer := hmacHash.Sum(nil)
	
	return mymath.Hex2Str(answer)
}

func CallPost(c appengine.Context, callURL string, parameters map[string]string) map[string]interface{} {
	return Call(c, callURL, parameters, "POST")
}
func CallGet(c appengine.Context, callURL string, parameters map[string]string) map[string]interface{} {
	return Call(c, callURL, parameters, "GET")
}

func Call(c appengine.Context, callURL string, parameters map[string]string, method string) map[string]interface{} {
	TimeoutDuration, err := time.ParseDuration("60s");
	if err!=nil {
		Log.Infof(c, "API - Call - error 1 - %s", err.Error())
		return nil
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
		return nil
	}
	req.Header.Set("key", parameters["key"])
	req.Header.Set("signature", parameters["signature"])
	req.Header.Set("nonce", parameters["nonce"])
	
	
	Log.Infof(c, "req - %v", req)
	//sending request
	//resp, err:=client.PostForm(callURL+function, values)
	resp, err:=client.Do(req)
	if err != nil {
		Log.Errorf(c, "API post error: %s", err)
		return nil
	}
	defer resp.Body.Close()
	//reading response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Errorf(c, "API read error: could not read body: %s", err)
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

//Public functions

func Ticker(c appengine.Context) map[string]interface{} {
	params:=getEmptyMap()
	return CallGet(c, "https://www.bitstamp.net/api/ticker/", params)
}

func OrderBook(c appengine.Context) map[string]interface{} {
	return OrderBookGrouped(c, true)
}

func OrderBookGrouped(c appengine.Context, group bool) map[string]interface{} {
	params:=getEmptyMap()
	if group {
		params["group"]="1"
	} else {
		params["group"]="0"
	}
	return CallGet(c, "https://www.bitstamp.net/api/order_book/", params)
}

func Transactions(c appengine.Context) map[string]interface{} {
	return TransactionWithTimeFrame(c, "hour")
}

func TransactionWithTimeFrame(c appengine.Context, time string) map[string]interface{} {
	params:=getEmptyMap()
	params["time"]=time
	return CallGet(c, "https://www.bitstamp.net/api/transactions/", params)
}

func EurUsdConversionRate(c appengine.Context) map[string]interface{} {
	params:=getBasicMap()
	return CallGet(c, "https://www.bitstamp.net/api/eur_usd/", params)
}




//Private functions

func AccountBalance(c appengine.Context) map[string]interface{} {
	params:=getBasicMap()
	return CallPost(c, "https://www.bitstamp.net/api/balance/", params)
}

func UserTransactions(c appengine.Context) map[string]interface{} {
	return UserTransactions4(c, 0, 100, true)
}

func UserTransactions2(c appengine.Context, offset int64) map[string]interface{} {
	return UserTransactions4(c, offset, 100, true)
}

func UserTransactions3(c appengine.Context, offset int64, limit int64) map[string]interface{} {
	return UserTransactions4(c, offset, limit, true)
}

func UserTransactions4(c appengine.Context, offset int64, limit int64, descending bool) map[string]interface{} {
	params:=getBasicMap()
	params["offset"]=mymath.Int642Str(offset)
	params["limit "]=mymath.Int642Str(limit)
	if descending {
		params["sort"]="desc"
	} else {
		params["sort"]="asc"
	}
	return CallPost(c, "https://www.bitstamp.net/api/user_transactions/", params)
}

func OpenOrders(c appengine.Context) map[string]interface{} {
	params:=getBasicMap()
	return CallPost(c, "https://www.bitstamp.net/api/open_orders/", params)
}

func CancelOrder(c appengine.Context, orderID string) map[string]interface{} {
	params:=getBasicMap()
	params["id"]=orderID
	return CallPost(c, "https://www.bitstamp.net/api/cancel_order/", params)
}

func BuyLimitOrder(c appengine.Context, amount float64, price float64) map[string]interface{} {
	params:=getBasicMap()
	params["amount"]=mymath.Float642Str(mymath.TrimFloatToXDecimalDigits(amount, 8))
	params["price"]=mymath.Float642Str(price)
	return CallPost(c, "https://www.bitstamp.net/api/buy/", params)
}

func SellLimitOrder(c appengine.Context, amount float64, price float64) map[string]interface{} {
	params:=getBasicMap()
	params["amount"]=mymath.Float642Str(mymath.TrimFloatToXDecimalDigits(amount, 8))
	params["price"]=mymath.Float642Str(price)
	return CallPost(c, "https://www.bitstamp.net/api/sell/", params)
}

func WithdrawalRequests(c appengine.Context) map[string]interface{} {
	params:=getBasicMap()
	return CallPost(c, "https://www.bitstamp.net/api/withdrawal_requests/", params)
}

func BitcoinWithdrawal(c appengine.Context, amount float64, address string) map[string]interface{} {
	params:=getBasicMap()
	params["amount"]=mymath.Float642Str(mymath.TrimFloatToXDecimalDigits(amount, 8))
	params["address"]=address
	return CallPost(c, "https://www.bitstamp.net/api/bitcoin_withdrawal/", params)
}

func BitcoinDepositAddress(c appengine.Context) map[string]interface{} {
	params:=getBasicMap()
	return CallPost(c, "https://www.bitstamp.net/api/bitcoin_deposit_address/", params)
}

func UnconfirmedBitcoinDeposits(c appengine.Context) map[string]interface{} {
	params:=getBasicMap()
	return CallPost(c, "https://www.bitstamp.net/api/unconfirmed_btc/", params)
}

func RippleWithdrawal(c appengine.Context, amount float64, address string, currency string) map[string]interface{} {
	params:=getBasicMap()
	params["amount"]=mymath.Float642Str(mymath.TrimFloatToXDecimalDigits(amount, 8))
	params["address"]=address
	params["currency"]=currency
	return CallPost(c, "https://www.bitstamp.net/api/ripple_withdrawal/", params)
}

func RippleDepositAddress(c appengine.Context) map[string]interface{} {
	params:=getBasicMap()
	return CallPost(c, "https://www.bitstamp.net/api/ripple_address/", params)
}

func getEmptyMap() map[string]string {
	answer := make(map[string]string)
	return answer
}
func getBasicMap() map[string]string {
	answer := make(map[string]string)
	nonce:=GetNonceStr()
	answer["key"]=BitstampKey
	answer["signature"]=GenerateSignature(nonce, BitstampClientID, BitstampKey, BitstampSecret)
	answer["nonce"]=nonce
	return answer
}
