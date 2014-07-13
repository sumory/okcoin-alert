package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type TickerInfo struct {
	Buy  string
	High string
	Last string
	Low  string
	Sell string
	Vol  string
}

type Ticker struct {
	Ticker TickerInfo
}

func tickerBTC() string {
	response, _ := http.Get("https://www.okcoin.cn/api/ticker.do")
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return string(body)
	}
	return ""
}

func tickerLTC() string {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", "https://www.okcoin.cn/api/ticker.do?symbol=ltc_cny", nil)
	//request.Header.Set("Accept","*/*")

	response, _ := client.Do(request)
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return string(body)
	}
	return ""
}

func dealBTC() (string, error) {
	var result string
	var error error
	func() {
		defer func() { //捕获panic异常
			//fmt.Println("defer...")
			if err := recover(); err != nil {
				//fmt.Println(err) //panic传入的内容
				error = errors.New(err.(string))
			}
		}()

		tickerB := tickerBTC()
		var btc Ticker
		if err := json.Unmarshal([]byte(tickerB), &btc); err != nil {
			panic("error when parse BTC data")
		}
		result = btc.Ticker.Last
		error = nil

	}()

	return result, error
}

func dealLTC() (string, error) {
	var result string
	var error error
	func() {
		defer func() { //捕获panic异常
			if err := recover(); err != nil {
				//fmt.Println(err) //panic传入的内容
				error = errors.New(err.(string))
			}
		}()

		tickerL := tickerLTC()
		var ltc Ticker
		if err := json.Unmarshal([]byte(tickerL), &ltc); err != nil {
			panic("error when parse LTC data")
		}
		result = ltc.Ticker.Last
		error = nil
	}()

	return result, error
}

func SendSMS(m string) {
	fmt.Println("send sms: ", m)
}

func main() {
	go func() {

		t := time.Tick(5 * time.Second)
		count := 0

		for now := range t {
			count++
			//			if count > 3 {
			//				break
			//			}

			fmt.Println("\n", now)

			btcLast, _ := dealBTC()
			ltcLast, _ := dealLTC()
			//fmt.Println("btcLast:", btcLast)
			//fmt.Println("ltcLast:", ltcLast)

			//			if err1 == nil {
			//				//btcLastValue,_ := strconv.ParseFloat(btcLast.(string), 64)
			//				fmt.Println("btcLast:", btcLast)
			//			} else {
			//				fmt.Println("err1:", err1)
			//			}
			//
			//			if err2 == nil {
			//				ltcLastValue, _ := strconv.ParseFloat(ltcLast, 64)
			//				fmt.Println("ltcLast:", ltcLastValue)
			//			} else {
			//				fmt.Println("err2:", err2)
			//			}

			SendSMS("btcLast: " + btcLast + " ltcLast: " + ltcLast)
		}
	}()

	server := http.Server{
		Addr:        ":8000",
		ReadTimeout: 5 * time.Second,
	}

	http.HandleFunc("/ticker", func(w http.ResponseWriter, r *http.Request) {
		btcLast, _ := dealBTC()
		ltcLast, _ := dealLTC()
		w.Write([]byte("btcLast: " + btcLast + " ltcLast: " + ltcLast))
	})
	server.ListenAndServe()

}
