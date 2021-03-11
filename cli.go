package main

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
)


type Quote struct {
	Exchange string
	PreMarketPrice float64
	Symbol string
	MarketState string
}

type Result struct {
	Result []Quote `json:"result"`
}

type QuoteResponse struct {
	QuoteResponse Result `json:"quoteResponse"`	
}

func getSymbols() string {
	return strings.Join(os.Args[1:], ",")
}

func getQuote(symbols string) {
	api := "https://query1.finance.yahoo.com/v7/finance/quote?lang=en-US&region=US&corsDomain=finance.yahoo.com&symbols=gme,tsla"
	resp, err := http.Get(api)

	if err != nil {
		log.Fatal(err)
	}

	respData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var quote QuoteResponse

	json.Unmarshal(respData, &quote)

	fmt.Println(quote.QuoteResponse.Result)
}

func main() {
	getQuote(getSymbols())
}
