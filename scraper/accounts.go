package scraper

import (
	"log"
	"pumpFunAccountScraper/api"
	"pumpFunAccountScraper/storage"
	"sync"
	"time"
)

type Queue[T any] struct {
	elements []T
	mu       sync.Mutex
}

func (q *Queue[T]) Enqueue(value T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.elements = append(q.elements, value)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.elements) == 0 {
		var zeroValue T
		return zeroValue, false
	}

	value := q.elements[0]
	q.elements = q.elements[1:]
	return value, true
}

type AccountResult struct {
	Tokens  *[]api.Token
	Account string
}

type Task struct {
	Function func(string)
	Argument string
}

type AccountScraper struct {
	result   chan<- AccountResult
	LastFind struct {
		Date time.Time
		mu   sync.Mutex
	}
	accountsCache storage.Cache
	tokensCache   storage.Cache
	workers       struct {
		workersNum  int
		workersChan chan Task
		queue       Queue[Task]
	}
}

func NewAccountScrapper(result chan<- AccountResult, workersNum int, accountsCache storage.Cache, tokensCache storage.Cache) *AccountScraper {
	scraper := &AccountScraper{result, struct {
		Date time.Time
		mu   sync.Mutex
	}{time.Now(), sync.Mutex{}}, accountsCache, tokensCache, struct {
		workersNum  int
		workersChan chan Task
		queue       Queue[Task]
	}{workersNum, make(chan Task), Queue[Task]{elements: []Task{}}}}
	scraper.startWorkers()
	go scraper.handleQueue()
	return scraper
}

func (a *AccountScraper) handleQueue() {
	for {
		task, exist := a.workers.queue.Dequeue()
		if exist {
			a.workers.workersChan <- task
		}
	}
}

func (a *AccountScraper) registerLastFind() {
	a.LastFind.mu.Lock()
	defer a.LastFind.mu.Unlock()
	a.LastFind.Date = time.Now()
}

func (a *AccountScraper) worker() {
	for work := range a.workers.workersChan {
		work.Function(work.Argument)
	}
}

func (a *AccountScraper) startWorkers() {
	for i := 0; i < a.workers.workersNum; i++ {
		go a.worker()
	}
}

func (a *AccountScraper) ScrapeAccount(account string) {

	allow, err := a.accountAllowed(account)
	if err != nil {
		log.Println(err)
		return
	}

	if !allow {
		return
	}

	tokens, err := api.GetAccountTokens(account)
	if err != nil {
		log.Println(err)
		return
	}

	err = a.accountsCache.AddKey(account, "1")
	if err != nil {
		log.Println(err)
		return
	}

	a.result <- AccountResult{tokens, account}

	for _, token := range *tokens {
		a.workers.queue.Enqueue(Task{a.ScrapeToken, token.Mint})
	}
}

func (a *AccountScraper) ScrapeToken(token string) {

	allow, err := a.tokenAllowed(token)
	if err != nil {
		log.Println(err)
		return
	}

	if !allow {
		return
	}

	trades, err := api.GetTokenTrades(token)
	if err != nil {
		log.Println(err)
		return
	}

	err = a.tokensCache.AddKey(token, "1")
	if err != nil {
		log.Println(err)
		return
	}

	for _, trade := range *trades {
		a.workers.queue.Enqueue(Task{a.ScrapeAccount, trade.User})
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
