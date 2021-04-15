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
	Symbol                     string
	RegularMarketPrice         float64
	RegularMarketChange        float64
	RegularMarketChangePercent float64
	PreMarketPrice             float64
	PreMarketChange            float64
	PostMarketPrice            float64
	PostMarketChange           float64
	RegularMarketPreviousClose float64
	MarketState                string
	Currency                   string
	Exchange                   string
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
	api := "https://query1.finance.yahoo.com/v7/finance/quote?corsDomain=finance.yahoo.com&symbols=" + symbols
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

/*
"postMarketChangePercent": -0.21913771,
"postMarketTime": 1617827573,
"postMarketPrice": 177.58,
"postMarketChange": -0.3899994,
"regularMarketChange": -6.53,
"regularMarketChangePercent": -3.53929,
"regularMarketTime": 1617825602,
"regularMarketPrice": 177.97,
"regularMarketDayHigh": 184.46,
"regularMarketDayRange": "176.11 - 184.46",
"regularMarketDayLow": 176.11,
"regularMarketVolume": 4602621,
"regularMarketPreviousClose": 184.5,
*/

func getPostPreMarket(preMarket float64, postMarket float64, marketPrice float64) string {
	postPreMarketPrice := preMarket
	if preMarket == 0 {
		postPreMarketPrice = postMarket
	}
	postPreMarketPriceStr := fmt.Sprintf("%.2f", postPreMarketPrice)

	// set red green normal for the premarket price
	if postPreMarketPrice >= marketPrice {
		postPreMarketPriceStr = pterm.LightGreen(postPreMarketPriceStr)
	} else if postPreMarketPrice == 0.00 {
		postPreMarketPriceStr = pterm.Normal(postPreMarketPriceStr)
	} else {
		postPreMarketPriceStr = pterm.LightRed(postPreMarketPriceStr)
	}

	return postPreMarketPriceStr
}

func getPostPreMarketChange(postMarketChange float64, preMarketChange float64) string {
	if postMarketChange != 0.00 {
		return fmt.Sprintf(" (%.2f)", postMarketChange)
	}

	if preMarketChange != 0.00 {
		return fmt.Sprintf(" (%.2f)", preMarketChange)
	}

	return ""
}

func printTable(quote []Quote) {
	table := pterm.TableData{{"Symbol", "Previous Price", "Price", "%", "PPP", "State", "Currency", "Exchange"}}

	for _, elem := range quote {
		regularMarketPreviousClose := elem.RegularMarketPreviousClose
		marketPrice := elem.RegularMarketPrice
		marketPriceChange := fmt.Sprintf(" (%.2f)", elem.RegularMarketChange)
		marketPriceChangePercent := fmt.Sprintf("%.2f", elem.RegularMarketChangePercent)
		regularMarketPreviousCloseStr := fmt.Sprintf("%.2f", regularMarketPreviousClose)
		marketPriceStr := fmt.Sprintf("%.2f", marketPrice)

		// Set red green for the market price text
		if marketPrice >= regularMarketPreviousClose {
			marketPriceStr = pterm.LightGreen(marketPriceStr)
		} else {
			marketPriceStr = pterm.LightRed(marketPriceStr)
		}

		table = append(
			table,
			[]string{
				elem.Symbol,
				regularMarketPreviousCloseStr,
				marketPriceStr + marketPriceChange,
				marketPriceChangePercent,
				getPostPreMarket(elem.PreMarketPrice, elem.PostMarketPrice, marketPrice) + getPostPreMarketChange(elem.PostMarketChange, elem.PreMarketChange),
				elem.MarketState,
				elem.Currency,
				elem.Exchange},
		)
	}

	pterm.DefaultTable.WithHasHeader().WithData(table).Render()
}

func printFooter() {
	pterm.Println(pterm.Gray("Made with <3 by github.com/bjarneo"))
}

func main() {
	for {
		time.Sleep(time.Second * 5)

		quote := getQuote(getSymbols())

		clear()

		printTable(quote)
		printFooter()
	}

}
