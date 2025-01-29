package runner

import (
	"log"
	"sync"
	"time"

	"github.com/denner-s/gorpcstress/internal/config"
	metrics "github.com/denner-s/gorpcstress/internal/metrics"
	"github.com/denner-s/gorpcstress/pkg/rpcclient"
)

type StressRunner struct {
	cfg     *config.Config
	metrics *metrics.Collector
}

func NewStressRunner(cfg *config.Config, collector *metrics.Collector) *StressRunner {
	return &StressRunner{
		cfg:     cfg,
		metrics: collector,
	}
}

func (sr *StressRunner) Run() {
	var wg sync.WaitGroup
	results := make(chan metrics.Result, sr.cfg.TotalRequests)
	done := make(chan struct{})

	go sr.collectResults(results, done)

	requestsPerWorker, remaining := distributeRequests(sr.cfg.Concurrency, sr.cfg.TotalRequests)

	for i := 0; i < sr.cfg.Concurrency; i++ {
		numRequests := requestsPerWorker
		if i < remaining {
			numRequests++
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			sr.runWorker(numRequests, results)
		}()
	}

	wg.Wait()
	close(results)
	<-done
}

func (sr *StressRunner) collectResults(results <-chan metrics.Result, done chan<- struct{}) {
	defer close(done)

	for res := range results {
		sr.metrics.RecordResult(res.Duration, res.Error)
	}
}

func distributeRequests(workers, total int) (base, remaining int) {
	base = total / workers
	remaining = total % workers
	return
}

func (sr *StressRunner) runWorker(requests int, results chan<- metrics.Result) {
	client, err := rpcclient.NewClient(sr.cfg.ServerAddress)
	if err != nil {
		log.Printf("Erro ao conectar: %v", err)
		sr.sendConnectionErrors(requests, results, err)
		return
	}
	defer func(client *rpcclient.Client) {
		err := client.Close()
		if err != nil {
			log.Printf("Erro ao desconectar: %v", err)
		}
	}(client)

	args := &rpcclient.Args{A: 5, B: 3}
	var reply rpcclient.Reply

	for i := 0; i < requests; i++ {
		start := time.Now()
		err := client.Call(sr.cfg.RPCMethod, args, &reply)
		duration := time.Since(start)

		results <- metrics.Result{
			Duration: duration,
			Error:    err,
		}
	}
}

func (sr *StressRunner) sendConnectionErrors(requests int, results chan<- metrics.Result, err error) {
	for i := 0; i < requests; i++ {
		results <- metrics.Result{
			Error: err,
		}
	}
}
