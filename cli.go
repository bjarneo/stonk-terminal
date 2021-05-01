package main

import (
	"encoding/json"
	"flag"
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
	RegularMarketDayHigh       float64
	RegularMarketDayLow        float64
	Bid                        float64
	Ask                        float64
	BidSize                    int
	AskSize                    int
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

func getSymbols(separator string) string {
	return strings.Join(os.Args[1:], separator)
}

func getQuote(symbols string) []Quote {
	// https://tutorialedge.net/golang/consuming-restful-api-with-go/
	api := "https://query1.finance.yahoo.com/v7/finance/quote?corsDomain=finance.yahoo.com&symbols=" + symbols
	resp, err := http.Get(api)

	if err != nil {
		log.Fatal(err)
	}

	respData, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()

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

func printTitle() {
	fmt.Println("\033]0;STONKS: " + strings.ToUpper(getSymbols(" ")) + "\007")
}

func getPostPreMarket(preMarket float64, postMarket float64, marketPrice float64) string {
	postPreMarketPrice := preMarket
	if preMarket == 0 {
		postPreMarketPrice = postMarket
	}
	postPreMarketPriceStr := fmt.Sprintf("%.2f", postPreMarketPrice)

	// set red green normal for the premarket price
	if postPreMarketPrice >= marketPrice {
		postPreMarketPriceStr = pterm.LightGreen(postPreMarketPriceStr)
	} else {
		postPreMarketPriceStr = pterm.LightRed(postPreMarketPriceStr)
	}

	return postPreMarketPriceStr
}

func getPostPreMarketChange(postMarketChange float64, preMarketChange float64) string {
	if postMarketChange != 0 {
		return fmt.Sprintf(" (%.2f)", postMarketChange)
	}

	if preMarketChange != 0 {
		return fmt.Sprintf(" (%.2f)", preMarketChange)
	}

	return ""
}

func printTable(quote []Quote) {
	table := pterm.TableData{{"Symbol", "Prev Price", "Price", "%", "PPP", "Low", "High", "Bid", "Ask", "State", "Curr", "Exch"}}

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
				fmt.Sprintf("%.2f", elem.RegularMarketDayLow),
				fmt.Sprintf("%.2f", elem.RegularMarketDayHigh),
				fmt.Sprintf("%.2f (%d)", elem.Bid, elem.BidSize),
				fmt.Sprintf("%.2f (%d)", elem.Ask, elem.AskSize),
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

func printIntro() {
	clear()

	s, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("S", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("tonks", pterm.NewStyle(pterm.FgLightMagenta))).Srender()

	pterm.DefaultCenter.Println(s)

	time.Sleep(time.Duration(3) * time.Second)
}

func run() {
	printIntro()

	interval := flag.Int("i", 5, "interval set to refetch stock data")

	flag.Parse()

	printTitle()

	for {
		time.Sleep(time.Duration(*interval) * time.Second)

		quote := getQuote(getSymbols(","))

		clear()

		printTable(quote)
		printFooter()
	}
}

func main() {
	run()
}
