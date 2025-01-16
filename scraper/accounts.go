package scraper

import (
	"log"
	"pumpFunAccountScraper/api"
	"pumpFunAccountScraper/storage"
	"sync"
	"time"
)

const waitSleep = 5 * time.Microsecond

type AccountScraper struct {
	LastFind struct {
		Date time.Time
		mu   sync.Mutex
	}
	accountsCache storage.Cache
	tokensCache   storage.Cache
	workers       struct {
		workersLimit int64
		workersNum   int64
		mu           sync.Mutex
	}
}

func NewAccountScrapper(workersLimit int64, accountsCache storage.Cache, tokensCache storage.Cache) *AccountScraper {
	return &AccountScraper{struct {
		Date time.Time
		mu   sync.Mutex
	}{time.Now(), sync.Mutex{}}, accountsCache, tokensCache, struct {
		workersLimit int64
		workersNum   int64
		mu           sync.Mutex
	}{workersLimit, 0, sync.Mutex{}}}
}

func (a *AccountScraper) registerLastFind() {
	a.LastFind.mu.Lock()
	defer a.LastFind.mu.Unlock()
	a.LastFind.Date = time.Now()
}

func (a *AccountScraper) incWorker() bool {
	a.workers.mu.Lock()
	defer a.workers.mu.Unlock()
	if a.checkWorker() {
		a.workers.workersNum++
		return true
	}
	return false
}

func (a *AccountScraper) wait() {
	for !a.incWorker() {
		time.Sleep(waitSleep)
	}
}

func (a *AccountScraper) decWorker() {
	a.workers.mu.Lock()
	defer a.workers.mu.Unlock()
	if a.workers.workersNum > 0 {
		a.workers.workersNum--
	}
}

func (a *AccountScraper) checkWorker() bool {
	if a.workers.workersNum < a.workers.workersLimit {
		return true
	} else {
		return false
	}
}

func (a *AccountScraper) ScrapeAccount(account string) {
	defer a.decWorker()

	allow, err := a.accountAllowed(account)
	if err != nil {
		log.Panic(err)
	}

	if !allow {
		return
	}

	tokens, err := api.GetAccountTokens(account)
	if err != nil {
		log.Panic(err)
	}

	err = a.accountsCache.AddKey(account, "1")
	if err != nil {
		log.Panic(err)
	}

	//TODO: should send a result to channel with tokens

	for _, token := range *tokens {
		a.wait()
		go a.ScrapeToken(token.Mint)
	}
}

func (a *AccountScraper) ScrapeToken(token string) {
	defer a.decWorker()

	allow, err := a.tokenAllowed(token)
	if err != nil {
		log.Panic(err)
	}

	if !allow {
		return
	}

	trades, err := api.GetTokenTrades(token)
	if err != nil {
		log.Panic(err)
	}

	err = a.tokensCache.AddKey(token, "1")
	if err != nil {
		log.Panic(err)
	}

	for _, trade := range *trades {
		a.wait()
		go a.ScrapeAccount(trade.User)
	}
}

func (a *AccountScraper) tokenAllowed(token string) (bool, error) {
	exist, err := a.tokensCache.KeyExist(token)
	if err != nil {
		return false, err
	}

	if exist {
		return false, nil
	} else {
		return true, nil
	}
}

func (a *AccountScraper) accountAllowed(account string) (bool, error) {
	exist, err := a.accountsCache.KeyExist(account)
	if err != nil {
		return false, err
	}

	if exist {
		return false, nil
	} else {
		a.registerLastFind()

		return true, nil
	}
}
