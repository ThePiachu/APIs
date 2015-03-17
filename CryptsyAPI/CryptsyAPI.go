package CryptsyAPI

import(
	"net/url"
	"appengine"
	"appengine/urlfetch"
	"github.com/ThePiachu/Go/Log"
	"io/ioutil"
	"encoding/json"
	"time"
	"net/http"
	"github.com/ThePiachu/Go/mymath"
	"crypto/hmac"
	"crypto/sha512"
	"bytes"
	
)

//https://www.cryptsy.com/pages/api

var PubAPI string = "http://pubapi.cryptsy.com/api.php"
var AuthAPI string = "https://api.cryptsy.com/api"
var APIKey string
var APISecret string

func Setup(apiKey string, apiSecret string) {
	APIKey = apiKey
	APISecret = apiSecret
}


func GetNonceStr() string {
	now:=time.Now()
	nonce:=now.UnixNano()/1000
	return mymath.Int642Str(nonce)
}

func GenerateSignature(secretKey string, data map[string]string) string {
	values := url.Values{}
	for k, v:=range data {
		values.Set(k, v)
	}
	return GenerateSignatureFromValues(secretKey, values)
}

func GenerateSignatureFromValues(secretKey string, values url.Values) string {
	toEncode:=mymath.ASCII2Hex(values.Encode())
	
	key:=mymath.ASCII2Hex(secretKey)
	
	hmacHash:=hmac.New(sha512.New, key)
	hmacHash.Write(toEncode)
	
	answer := hmacHash.Sum(nil)

	return mymath.ToLower(mymath.Hex2Str(answer))
}

func PublicCall(c appengine.Context, callURL string, parameters map[string]string) (map[string]interface{}, error) {
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
	req, err:=http.NewRequest("GET", callURL+"?"+values.Encode(), nil)
	if err!=nil {
		Log.Infof(c, "API - Call - error 2 - %s", err.Error())
		return nil, err
	}
	
	
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
	result := make(map[string]interface{})
	//unmarshalling JSON response
    err = json.Unmarshal(body, &result)
    if err != nil {
        Log.Infof(c, "Unmarshal: %v", err)
        Log.Infof(c, "%s", body)
    	return nil, err
    }
    return result, nil
}








func PrivateCall(c appengine.Context, parameters map[string]string) (map[string]interface{}, error) {
	TimeoutDuration, err := time.ParseDuration("60s");
	if err!=nil {
		Log.Infof(c, "API - Call - error 1 - %s", err.Error())
		return nil, err
	}
    tr := urlfetch.Transport{Context: c, Deadline: TimeoutDuration}
	
	values := url.Values{}
	for k, v := range parameters {
		values.Set(k, v)
	}
	signature:=GenerateSignatureFromValues(APISecret, values)
	
	req, err:=http.NewRequest("POST", AuthAPI, bytes.NewBufferString(values.Encode()))
	if err!=nil {
		Log.Infof(c, "API - Call - error 2 - %s", err.Error())
		return nil, err
	}
	//req.Form=values
	req.Header.Set("Key", APIKey)
	req.Header.Set("Sign", signature)
	
	Log.Infof(c, "req - %v", req)
	client := http.Client{Transport: &tr}
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
	result := make(map[string]interface{})
	//unmarshalling JSON response
    err = json.Unmarshal(body, &result)
    if err != nil {
        Log.Infof(c, "Unmarshal: %v", err)
        Log.Infof(c, "%s", body)
    	return nil, err
    }
    return result, nil
}




//*************************** Public methods ************************

func GeneralMarketDataAll(c appengine.Context) (map[string]interface{}, error) {
	resp, err:=PublicCall(c, PubAPI, map[string]string{"method":"marketdatav2"})
	if err!=nil {
		Log.Errorf(c, "GeneralMarketDataAll error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}

func GeneralMarketDataSingle(c appengine.Context, marketID string) (map[string]interface{}, error) {
	resp, err:=PublicCall(c, PubAPI, map[string]string{"method":"singlemarketdata", "marketid":marketID})
	if err!=nil {
		Log.Errorf(c, "GeneralMarketDataSingle error 1 - %v", err)
		return nil, err
	}
	for _, v:= range resp["return"].(map[string]interface{}) {
		return v.(map[string]interface{}), nil
	}
	return resp["return"].(map[string]interface{})["markets"].(map[string]interface{}), nil
}

func GeneralOrderbookDataAll(c appengine.Context) (map[string]interface{}, error) {
	resp, err:=PublicCall(c, PubAPI, map[string]string{"method":"orderdatav2"})
	if err!=nil {
		Log.Errorf(c, "GeneralOrderbookDataAll error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}

func GeneralOrderbookDataSingle(c appengine.Context, marketID string) (map[string]interface{}, error) {
	resp, err:=PublicCall(c, PubAPI, map[string]string{"method":"singleorderdata", "marketid":marketID})
	if err!=nil {
		Log.Errorf(c, "GeneralOrderbookDataSingle error 1 - %v", err)
		return nil, err
	}
	for _, v:= range resp["return"].(map[string]interface{}) {
		return v.(map[string]interface{}), nil
	}
	return resp["return"].(map[string]interface{})["markets"].(map[string]interface{}), nil
}

//*************************** Private methods ************************

func GetInfo(c appengine.Context) (map[string]interface{}, error) {
	params:=GetBasicMap("getinfo")
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "GetInfo error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func GetMarkets(c appengine.Context) (map[string]interface{}, error) {
	params:=GetBasicMap("getmarkets")
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "GetMarkets error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func GetWalletStatus(c appengine.Context) (map[string]interface{}, error) {
	params:=GetBasicMap("getwalletstatus")
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "GetWalletStatus error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func MyTransactions(c appengine.Context) (map[string]interface{}, error) {
	params:=GetBasicMap("mytransactions")
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "MyTransactions error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func MarketTrades(c appengine.Context, marketID int) (map[string]interface{}, error) {
	params:=GetBasicMap("markettrades")
	params["marketid"]=mymath.Int2Str(marketID)
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "MarketTrades error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func MarketOrders(c appengine.Context, marketID int) (map[string]interface{}, error) {
	params:=GetBasicMap("marketorders")
	params["marketid"]=mymath.Int2Str(marketID)
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "MarketOrders error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func MyTrades(c appengine.Context, marketID int) (map[string]interface{}, error) {
	params:=GetBasicMap("mytrades")
	params["marketid"]=mymath.Int2Str(marketID)
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "MyTrades error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func MyTradesWithLimit(c appengine.Context, marketID int, limit int) (map[string]interface{}, error) {
	params:=GetBasicMap("mytrades")
	params["marketid"]=mymath.Int2Str(marketID)
	params["limit"]=mymath.Int2Str(limit)
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "MyTradesWithLimit error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func AllMyTrades(c appengine.Context) (map[string]interface{}, error) {
	params:=GetBasicMap("allmytrades")
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "AllMyTrades error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func AllMyTradesWithDates(c appengine.Context, startDate string, endDate string) (map[string]interface{}, error) {
	params:=GetBasicMap("allmytrades")
	params["startdate"]=startDate
	params["enddate"]=endDate
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "AllMyTradesWithDates error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func MyOrders(c appengine.Context, marketID int) (map[string]interface{}, error) {
	params:=GetBasicMap("myorders")
	params["marketid"]=mymath.Int2Str(marketID)
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "MyOrders error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func Depth(c appengine.Context, marketID int) (map[string]interface{}, error) {
	params:=GetBasicMap("depth")
	params["marketid"]=mymath.Int2Str(marketID)
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "Depth error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func AllMyOrders(c appengine.Context) (map[string]interface{}, error) {
	params:=GetBasicMap("allmyorders")
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "AllMyOrders error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}

func CreateOrder(c appengine.Context, marketID string, orderType string, quantity string, price string) (map[string]interface{}, error) {
	params:=GetBasicMap("createorder")
	params["marketid"]=marketID
	params["ordertype"]=orderType
	params["quantity"]=quantity
	params["price"]=price
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "CreateOrder error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func CancelOrder(c appengine.Context, orderID string) (map[string]interface{}, error) {
	params:=GetBasicMap("cancelorder")
	params["orderid"]=orderID
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "CancelOrder error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func CancelMarketOrders(c appengine.Context, marketID int) (map[string]interface{}, error) {
	params:=GetBasicMap("cancelmarketorders")
	params["marketid"]=mymath.Int2Str(marketID)
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "CancelMarketOrders error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func CancelAllOrders(c appengine.Context) (map[string]interface{}, error) {
	params:=GetBasicMap("cancelallorders")
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "CancelAllOrders error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func CalculateFees(c appengine.Context, orderType string, quantity string, price string) (map[string]interface{}, error) {
	params:=GetBasicMap("calculatefees")
	params["ordertype"]=orderType
	params["quantity"]=quantity
	params["price"]=price
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "CalculateFees error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func GenerateNewAddressWithCurrencyID(c appengine.Context, currencyID int) (map[string]interface{}, error) {
	params:=GetBasicMap("generatenewaddress")
	params["currencyid"]=mymath.Int2Str(currencyID)
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "GenerateNewAddressWithCurrencyID error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func GenerateNewAddressWithCurrencyCode(c appengine.Context, currencyCode string) (map[string]interface{}, error) {
	params:=GetBasicMap("generatenewaddress")
	params["currencycode"]=currencyCode
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "GenerateNewAddressWithCurrencyCode error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func MyTransfers(c appengine.Context) (map[string]interface{}, error) {
	params:=GetBasicMap("mytransfers")
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "MyTransfers error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func MakeWithdrawal(c appengine.Context, address string, amount string) (map[string]interface{}, error) {
	params:=GetBasicMap("makewithdrawal")
	params["address"]=address
	params["amount"]=amount
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "MakeWithdrawal error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func GetMyDepositAddresses(c appengine.Context) (map[string]interface{}, error) {
	params:=GetBasicMap("getmydepositaddresses")
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "GetMyDepositAddresses error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}
func GetOrderStatus(c appengine.Context, orderID string) (map[string]interface{}, error) {
	params:=GetBasicMap("getorderstatus")
	params["orderid"]=orderID
	resp, err:=PrivateCall(c, params)
	if err!=nil {
		Log.Errorf(c, "GetOrderStatus error 1 - %v", err)
		return nil, err
	}
	return resp, nil
}



 
 
 
func GetBasicMap(method string) map[string]string{
	answer:=map[string]string{}
	answer["method"]=method
	answer["nonce"]=GetNonceStr()
	return answer
}
