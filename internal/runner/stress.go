package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/denner-s/gorpcstress/internal/config"
	"github.com/denner-s/gorpcstress/internal/metrics"
	"github.com/denner-s/gorpcstress/pkg/rpcclient"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// StressRunner gerencia toda a execução do teste de carga RPC
type StressRunner struct {
	cfg         *config.Config     // Configurações do teste
	metrics     *metrics.Collector // Coletor de métricas de desempenho
	payloadData *rpcclient.Args    // Dados do payload para as chamadas RPC
}

// Run inicia e controla o fluxo principal do teste de carga
func (sr *StressRunner) Run() {
	var wg sync.WaitGroup
	results := make(chan metrics.Result, sr.cfg.Concurrency*2) // Canal bufferizado para resultados
	done := make(chan struct{})                                // Canal para sinalização de término

	// Goroutine para coletar resultados de forma assíncrona
	go sr.collectResults(results, done)

	// Seleciona o modo de operação baseado na configuração
	if sr.cfg.Duration > 0 {
		sr.runDurationMode(time.Now(), &wg, results) // Modo de execução contínua por tempo
	} else {
		sr.runRequestMode(&wg, results) // Modo de número fixo de requisições
	}

	// Espera a conclusão de todas as goroutines
	wg.Wait()
	close(results) // Fecha o canal de resultados
	<-done         // Aguarda a finalização do processamento
}

// NewStressRunner é o construtor que inicializa o testador de carga
func NewStressRunner(cfg *config.Config, collector *metrics.Collector) *StressRunner {
	runner := &StressRunner{
		cfg:     cfg,
		metrics: collector,
	}

	// Carrega payload personalizado ou usa valores padrão
	if cfg.PayloadFile != "" {
		runner.loadPayload()
	} else {
		// Valores padrão que correspondem ao exemplo do servidor
		runner.payloadData = &rpcclient.Args{A: 5, B: 3}
	}

	return runner
}

// loadPayload carrega dados de chamada de um arquivo JSON
func (sr *StressRunner) loadPayload() {
	file, err := os.Open(sr.cfg.PayloadFile)
	if err != nil {
		log.Fatalf("Falha ao abrir arquivo de payload: %v", err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Printf("Erro ao fechar arquivo: %v", err)
		}
	}(file)

	sr.payloadData = &rpcclient.Args{}
	if err := json.NewDecoder(file).Decode(sr.payloadData); err != nil {
		log.Fatalf("Erro na decodificação do JSON: %v", err)
	}
}

// runDurationMode executa o teste continuamente por um período específico
func (sr *StressRunner) runDurationMode(start time.Time, wg *sync.WaitGroup, results chan<- metrics.Result) {
	ticker := time.NewTicker(time.Second / time.Duration(sr.cfg.Concurrency))
	defer ticker.Stop()

	// Loop enquanto estiver dentro da duração configurada
	for time.Since(start) < sr.cfg.Duration {
		<-ticker.C // Controla a taxa de requisições
		wg.Add(1)
		go func() {
			defer wg.Done()
			sr.runWorker(1, results) // Executa 1 requisição por goroutine
		}()
	}
}

// runRequestMode distribui requisições fixas entre workers
func (sr *StressRunner) runRequestMode(wg *sync.WaitGroup, results chan<- metrics.Result) {
	// Distribui requisições igualmente entre workers
	base, remaining := distributeRequests(sr.cfg.Concurrency, sr.cfg.TotalRequests)

	for i := 0; i < sr.cfg.Concurrency; i++ {
		reqCount := base
		if i < remaining {
			reqCount++ // Distribui requisições extras
		}

		wg.Add(1)
		go func(count int) {
			defer wg.Done()
			sr.runWorker(count, results) // Executa lote de requisições
		}(reqCount)
	}
}

// collectResults processa e armazena os resultados das requisições
func (sr *StressRunner) collectResults(results <-chan metrics.Result, done chan<- struct{}) {
	defer close(done)
	for res := range results {
		sr.metrics.RecordResult(res) // Registra no coletor de métricas
	}
}

// distributeRequests calcula a distribuição de carga entre workers
func distributeRequests(workers, total int) (base int, remaining int) {
	base = total / workers
	remaining = total % workers
	return
}

// runWorker executa um lote de requisições RPC
func (sr *StressRunner) runWorker(requests int, results chan<- metrics.Result) {
	client, err := rpcclient.NewClient(sr.cfg.ServerAddress, sr.cfg.Timeout)
	if err != nil {
		log.Printf("Falha na conexão RPC: %v", err)
		sr.sendConnectionErrors(requests, results, err)
		return
	}
	defer func(client *rpcclient.Client) {
		if err := client.Close(); err != nil {
			log.Printf("Erro ao fechar cliente: %v", err)
		}
	}(client)

	// Executa o número especificado de requisições
	for i := 0; i < requests; i++ {
		start := time.Now()
		var reply rpcclient.Reply

		// Chamada RPC principal
		err := client.Call(sr.cfg.RPCMethod, sr.payloadData, &reply)
		duration := time.Since(start)

		// Cria resultado com análise de erro
		results <- metrics.Result{
			Duration: duration,
			Error:    sr.analyzeError(err, &reply),
		}
	}
}

// analyzeError processa e classifica erros da chamada RPC
func (sr *StressRunner) analyzeError(err error, reply *rpcclient.Reply) error {
	if err != nil {
		return categorizeError(err) // Classifica erros de rede
	}

	// Verificação rigorosa do resultado
	expected := sr.payloadData.A * sr.payloadData.B
	if reply.Result != expected {
		return fmt.Errorf("resultado incorreto: esperado %d, recebido %d", expected, reply.Result)
	}
	return nil
}

// categorizeError classifica os tipos de erro para relatórios
func categorizeError(err error) error {
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return fmt.Errorf("timeout: %w", err)
		}
		return fmt.Errorf("erro de rede: %w", err)
	}
	return err
}

// sendConnectionErrors registra falhas de conexão para todas as requisições afetadas
func (sr *StressRunner) sendConnectionErrors(requests int, results chan<- metrics.Result, connErr error) {
	for i := 0; i < requests; i++ {
		results <- metrics.Result{
			Error: fmt.Errorf("falha na conexão: %w", connErr),
		}
	}
}
