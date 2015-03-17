package BlockchainAPI

import(
	"net/url"
	"appengine"
	"appengine/urlfetch"
	"io/ioutil"
	"encoding/json"
	"time"
	"net/http"
	"github.com/ThePiachu/Go/Log"
)

//https://blockchain.info/api/exchange_rates_api

func CallURL(c appengine.Context, callURL string, parameters map[string]string) map[string]interface{} {
	TimeoutDuration, err := time.ParseDuration("60s");
	if err!=nil {
		Log.Infof(c, "BlockchainAPI - CallURL - error - %s", err.Error())
	}
    tr := urlfetch.Transport{Context: c, Deadline: TimeoutDuration}
	client := http.Client{Transport: &tr}
	values := url.Values{}
	//turning input map into url.Values
	for k, v := range parameters {
		values.Set(k, v)
	}
	//sending request
	resp, err:=client.PostForm(callURL, values)
	if err != nil {
		Log.Errorf(c, "BlockchainAPI - CallURL - post error: %s", err)
		return nil
	}
	defer resp.Body.Close()
	//reading response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Errorf(c, "BlockchainAPI - CallURL - read error: could not read body: %s", err)
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

func Ticker(c appengine.Context) map[string]interface{} {
	return CallURL(c, "https://blockchain.info/ticker", nil)
}

//TODO: do ToBTC
