package main

import (
	"fmt"
	"log"

	"github.com/denner-s/gorpcstress/internal/config"
	"github.com/denner-s/gorpcstress/internal/metrics"
	"github.com/denner-s/gorpcstress/internal/runner"
	"github.com/denner-s/gorpcstress/pkg/report"
)

func main() {
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuração inválida: %v", err)
	}

	collector := metrics.NewCollector()
	stressRunner := runner.NewStressRunner(cfg, collector)

	fmt.Printf("Iniciando teste de estresse...\nServidor: %s\nRequisições: %d\nConcorrência: %d\n\n",
		cfg.ServerAddress, cfg.TotalRequests, cfg.Concurrency)

	stressRunner.Run()

	report.GenerateReport(collector.GetMetrics())
}
