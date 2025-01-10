package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const pumpFunUrl = "https://frontend-api-v2.pump.fun"

type Token struct {
	Address   string  `json:"address"`
	Mint      string  `json:"mint"`
	Balance   int64   `json:"balance"`
	ImageUrl  string  `json:"image_url"`
	Symbol    string  `json:"symbol"`
	Name      string  `json:"name"`
	MarketCap float64 `json:"market_cap"`
	Value     float64 `json:"value"`
}

type Trade struct {
	Signature    string `json:"signature"`
	Mint         string `json:"mint"`
	SolAmount    int64  `json:"sol_amount"`
	TokenAmount  int64  `json:"token_amount"`
	IsBuy        bool   `json:"is_buy"`
	User         string `json:"user"`
	Timestamp    int64  `json:"timestamp"`
	TxIndex      int64  `json:"tx_index"`
	Username     string `json:"username"`
	ProfileImage string `json:"profile_image"`
	Slot         int64  `json:"slot"`
}

func GetAccountTokens(account string) (*[]Token, error) {
	var allTokens []Token
	limit := 100
	offset := 0
	for {
		url := fmt.Sprintf("%s/balances/%s?limit=%d&offset=%d&minBalance=-1", pumpFunUrl, account, limit, offset)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		var tokens []Token
		if err := json.Unmarshal(body, &tokens); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		if len(tokens) == 0 {
			break
		}

		allTokens = append(allTokens, tokens...)

		offset += limit
	}

	return &allTokens, nil
}

func GetTokenTrades(token string) (*[]Trade, error) {
	var allTrades []Trade
	limit := 200
	offset := 0
	for {
		url := fmt.Sprintf("%s/trades/all/%s?limit=%d&offset=%d&minimumSize=0", pumpFunUrl, token, limit, offset)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		var trades []Trade
		if err := json.Unmarshal(body, &trades); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		if len(trades) == 0 {
			break
		}

		allTrades = append(allTrades, trades...)

		offset += limit
	}
	return &allTrades, nil
}
