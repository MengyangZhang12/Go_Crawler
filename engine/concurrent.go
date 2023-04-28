package engine

import "fmt"

type ConcurrentEngine struct {
	Scheduler        Scheduler
	WorkerCount      int
	ItemChan         chan Item
	RequestProcessor Processor
}

type Processor func(Request) (ParseResult, error)

type Scheduler interface {
	ReadyNotifier
	Submit(Request)
	WorkerChan() chan Request
	Run()
}

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParseResult)
	e.Scheduler.Run()

	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}
	for _, r := range seeds {
		if isDuplicate(r.Url) {
			fmt.Printf("Duplicate request: %s\n", r.Url)
			continue
		}
		e.Scheduler.Submit(r)
	}
	for {
		result := <-out
		for _, item := range result.Items {
			go func() {
				e.ItemChan <- item
			}()
		}
		for _, request := range result.Requests {
			if isDuplicate(request.Url) {
				fmt.Printf("Duplicate request: %s\n", request.Url)
				continue
			}
			e.Scheduler.Submit(request)
		}
	}
}

var visitedUrls = make(map[string]bool)

func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}
	visitedUrls[url] = true
	return false
}

func (e *ConcurrentEngine) createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			// tell scheduler worker ready
			ready.WorkerReady(in)
			request := <-in
			result, err := e.RequestProcessor(request) //call rpc
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
