package VaultAPI

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
	"crypto/sha512"
	"bytes"
	"github.com/ThePiachu/Go/Log"
)

var URL string = "https://api.vaultofsatoshi.com"
var APIKey string
var APISecret string

func Setup(apiKey string, apiSecret string) {
	APIKey = apiKey
	APISecret = apiSecret
}

func GetNonce() int64 {
	now:=time.Now()
	return now.UnixNano()/1000
}

func GetNonceStr() string {
	nonce:=GetNonce()
	return mymath.Int642Str(nonce)
}

func GenerateSignature(secretKey string, function string, data map[string]string) string {
	values := url.Values{}
	for k, v:=range data {
		values.Set(k, v)
	}
	return GenerateSignatureFromValues(secretKey, function, values)
}

func GenerateSignatureFromValues(secretKey string, function string, values url.Values) string {
	query:=mymath.ASCII2Hex(values.Encode())
	toEncode:=mymath.ASCII2Hex(function)
	toEncode = append(toEncode, 0x00)
	toEncode = append(toEncode, query...)

	key:=mymath.ASCII2Hex(secretKey)
	
	hmacHash:=hmac.New(sha512.New, key)
	hmacHash.Write(toEncode)
	
	answer := hmacHash.Sum(nil)
	
	return mymath.Hex2Base64(mymath.ASCII2Hex(mymath.ToLower(mymath.Hex2Str(answer))))
}

func Call(c appengine.Context, callURL string, function string, apiKey string, apiSecret string, parameters map[string]string) map[string]interface{} {
	TimeoutDuration, err := time.ParseDuration("60s");
	if err!=nil {
		Log.Infof(c, "API - Call - error 1 - %s", err.Error())
		return nil
	}
    tr := urlfetch.Transport{Context: c, Deadline: TimeoutDuration}
	values := url.Values{}
	//turning input map into url.Values
	for k, v := range parameters {
		values.Set(k, v)
	}
	values.Set("nonce", GetNonceStr())
	
	signature:=GenerateSignatureFromValues(apiSecret, function, values)
	
	req, err:=http.NewRequest("POST", callURL+function, bytes.NewBufferString(values.Encode()))
	if err!=nil {
		Log.Infof(c, "API - Call - error 2 - %s", err.Error())
		return nil
	}
	//req.Form=values
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Sign", signature)
	
	Log.Infof(c, "req - %v", req)
	Log.Infof(c, "FormValue[nonce] - %v", req.FormValue("nonce"))
	//sending request
	resp, err:=tr.RoundTrip(req)
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

func InfoCurrency(c appengine.Context) map[string]interface{} {
	params := getBasicMap()
	
	return Call(c, URL, "/info/currency", APIKey, APISecret, params)
}

func InfoCurrency2(c appengine.Context, currency string) map[string]interface{} {
	params := getBasicMap()
	params["currency"]=currency
	return Call(c, URL, "/info/currency", APIKey, APISecret, params)
}

func getBasicMap() map[string]string {
	answer := make(map[string]string)
	return answer
}

func Test2(c appengine.Context) {
	secret:=mymath.Hex2Str(mymath.ASCII2Hex("ENTER_YOUR_API_SECRET_HERE"))
	function:="/info/order_detail"
	data:=map[string]string {"nonce":"1386502805898680"}
	
	signature:=GenerateSignature(secret, function, data)
	
	Log.Infof(c, "signature - %s", signature)
}

func Test(c appengine.Context) {
	toEncode:=mymath.Str2Hex("2F696E666F2F63757272656E6379006E6F6E63653D31333836353032383035383938363830")
	secretKey:=mymath.ASCII2Hex("Enter your secret key here")
	
	
	Log.Infof(c, "dataToEncode - %X", toEncode)
	Log.Infof(c, "apiSecret - %X", secretKey)
	
	
	hmacHash:=hmac.New(sha512.New, secretKey)
	hmacHash.Write(toEncode)
	
	answer := hmacHash.Sum(nil)
	
	Log.Infof(c, "answer - %X", answer)
}
