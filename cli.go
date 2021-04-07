package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

type Quote struct {
	Symbol                     string
	RegularMarketPrice         float64
	PreMarketPrice             float64
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

func printTable(quote []Quote) {
	table := pterm.TableData{{"Symbol", "Previous Price", "Price", "%", "Pre Price", "State", "Currency", "Exchange"}}

	for _, elem := range quote {
		regularMarketPreviousClose := elem.RegularMarketPreviousClose
		marketPrice := elem.RegularMarketPrice
		preMarketPrice := elem.PreMarketPrice

		regularMarketPreviousCloseStr := fmt.Sprintf("%.2f", regularMarketPreviousClose)
		marketPriceStr := fmt.Sprintf("%.2f", marketPrice)
		preMarketPriceStr := fmt.Sprintf("%.2f", preMarketPrice)

		// set red green normal for the premarket price
		if preMarketPrice >= marketPrice {
			preMarketPriceStr = pterm.LightGreen(preMarketPriceStr)
		} else if preMarketPrice == 0.00 {
			preMarketPriceStr = pterm.Normal(preMarketPriceStr)
		} else {
			preMarketPriceStr = pterm.LightRed(preMarketPriceStr)
		}

		// Set red green for the market price text
		if marketPrice >= regularMarketPreviousClose {
			marketPriceStr = pterm.LightGreen(marketPriceStr)
		} else {
			marketPriceStr = pterm.LightRed(marketPriceStr)
		}

		// Add + if not is negative number
		marketPriceDiff := marketPrice - regularMarketPreviousClose
		marketPriceDiffStr := fmt.Sprintf(" (%.2f)", marketPriceDiff)
		if !math.Signbit(marketPriceDiff) {
			marketPriceDiffStr = fmt.Sprintf(" (+%.2f)", marketPriceDiff)
		}

		preMarketPriceDiff := ""
		if preMarketPrice != 0.00 {
			preMarketPriceDiff = fmt.Sprintf(" (%.2f)", preMarketPrice-regularMarketPreviousClose)
		}

		// Add + if not is negative number
		percentageDiff := (marketPrice / regularMarketPreviousClose * 100) - 100
		percentageDiffStr := fmt.Sprintf("%.2f", percentageDiff)
		if !math.Signbit(percentageDiff) {
			percentageDiffStr = fmt.Sprintf("+%.2f", percentageDiff)
		}

		table = append(
			table,
			[]string{
				elem.Symbol,
				regularMarketPreviousCloseStr,
				marketPriceStr + marketPriceDiffStr,
				percentageDiffStr,
				preMarketPriceStr + preMarketPriceDiff,
				elem.MarketState,
				elem.Currency,
				elem.Exchange},
		)
	}

	pterm.DefaultTable.WithHasHeader().WithData(table).Render()
}

func printLogo() {
	logo := `
    /$$$$$$  /$$$$$$$$/$$$$$$  /$$   /$$ /$$   /$$  /$$$$$$ 
   /$$__  $$|__  $$__/$$__  $$| $$$ | $$| $$  /$$/ /$$__  $$
  | $$  \__/   | $$ | $$  \ $$| $$$$| $$| $$ /$$/ | $$  \__/
  |  $$$$$$    | $$ | $$  | $$| $$ $$ $$| $$$$$/  |  $$$$$$ 
   \____  $$   | $$ | $$  | $$| $$  $$$$| $$  $$   \____  $$
   /$$  \ $$   | $$ | $$  | $$| $$\  $$$| $$\  $$  /$$  \ $$
  |  $$$$$$/   | $$ |  $$$$$$/| $$ \  $$| $$ \  $$|  $$$$$$/
   \______/    |__/  \______/ |__/  \__/|__/  \__/ \______/ 
                                                          
`

	pterm.Println(pterm.Green(logo))
}

func printFooter() {
	pterm.Println(pterm.Gray("Made with <3 by github.com/bjarneo"))
}

func main() {
	for {
		time.Sleep(time.Second * 5)

		quote := getQuote(getSymbols())

		clear()
		clear()

		printLogo()
		printTable(quote)
		printFooter()
	}

}
