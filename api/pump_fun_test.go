package api

import (
	"testing"
)

func TestGetAccountTokens(t *testing.T) {
	tokens, err := GetAccountTokens("12gUbZZoTNycRvvwFxCBzJqkYjEzGi5XCf8iFa6KqqcB")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(*tokens) == 0 {
		t.Errorf("Expected tokens is more than zero but got zero")
	}
}

func TestGetTokenTrades(t *testing.T) {
	trades, err := GetTokenTrades("Df6yfrKC8kZE3KNkrHERKzAetSxbrWeniQfyJY4Jpump")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(*trades) == 0 {
		t.Errorf("Expected trades is more than zero but got zero")
	}
}
