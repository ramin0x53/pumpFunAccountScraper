package main

import (
	"fmt"
	"pumpFunAccountScraper/scraper"
	"pumpFunAccountScraper/storage"
)

func main() {
	resultChannel := make(chan scraper.AccountResult, 100000)
	tokensCache := storage.NewRedisCache(1, "127.0.0.1:6379")
	accountsCache := storage.NewRedisCache(2, "127.0.0.1:6379")

	scraper := scraper.NewAccountScrapper(resultChannel, 5, accountsCache, tokensCache)
	go scraper.ScrapeToken("J19SvmVqkn6V7kJa9R2XEfpZGb6nxdxDE9TgPpqGpump")

	// for range resultChannel {
	for res := range resultChannel {
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++")
		fmt.Println(res.Account)
	}
}
