package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

type Quote struct {
	Symbol         string
	PreMarketPrice float64
	MarketState    string
	Exchange       string
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

func getQuote(symbols string) []Quote {
	// https://tutorialedge.net/golang/consuming-restful-api-with-go/
	api := "https://query1.finance.yahoo.com/v7/finance/quote?lang=en-US&region=US&corsDomain=finance.yahoo.com&symbols=" + symbols
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

	return quote.QuoteResponse.Result
}

func clear() {
	print("\033[H\033[2J")
}

func printTable(quote []Quote) {
	table := pterm.TableData{{"Symbol", "PreMarketPrice", "MarketState", "Exchange"}}

	for _, elem := range quote {
		table = append(
			table,
			[]string{elem.Symbol, fmt.Sprintf("%.2f", elem.PreMarketPrice), elem.MarketState, elem.Exchange},
		)
	}

	pterm.DefaultTable.WithHasHeader().WithData(table).Render()
}

func main() {
	for {
		time.Sleep(time.Second * 5)

		quote := getQuote(getSymbols())

		clear()
		/*
			for x, elem := range quote {
				fmt.Println(x)
				fmt.Println(elem)
			}
		*/

		printTable(quote)
	}

}
