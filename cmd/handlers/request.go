package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ZondaF12/crypto-bot/config"
)

const (
	API_ENDPOINT = "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest?symbol=%s&convert=%s"
)

type CoinResponse struct {
	Data map[string][]CoinData
}

type CoinData struct {
	Id    int
	Name  string
	Slug  string
	Quote map[string]CoinQuote
}

type CoinQuote struct {
	Price              float32
	Percent_change_24h float32
	Percent_change_7d  float32
}

func doRequest(url string) ([]byte, error) {
	env, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	fmt.Println("Fetching latest prices...")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-CMC_PRO_API_KEY", env.CMC_API_KEY)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func FetchPrice(symbol string, currencyCode string) CoinResponse {
	body, err := doRequest(fmt.Sprintf(API_ENDPOINT, symbol, currencyCode))
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	var coinResponse CoinResponse
	if err = json.Unmarshal(body, &coinResponse); err != nil {
		fmt.Printf("Error: %v", err)
	}

	return coinResponse
}
