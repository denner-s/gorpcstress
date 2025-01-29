package report

import (
	"fmt"
	"time"

	"github.com/denner-s/gorpcstress/internal/metrics"
)

func GenerateReport(m metrics.Metrics) {
	fmt.Println("\n=== Relatório do Teste de Estresse ===")

	printGeneralInfo(m)
	printThroughput(m)
	printLatencyMetrics(m)
}

func printGeneralInfo(m metrics.Metrics) {
	totalDuration := m.EndTime.Sub(m.StartTime)

	fmt.Printf("Tempo total de execução:\t %v\n", totalDuration.Round(time.Millisecond))
	fmt.Printf("Requisições totais:\t\t %d\n", m.TotalRequests)
	fmt.Printf("Requisições com erro:\t\t %d (%.2f%%)\n",
		m.Errors, errorRate(m))
}

func errorRate(m metrics.Metrics) float64 {
	return float64(m.Errors) / float64(m.TotalRequests) * 100
}

func printThroughput(m metrics.Metrics) {
	duration := m.EndTime.Sub(m.StartTime).Seconds()
	rps := float64(m.TotalRequests) / duration
	rpm := rps * 60

	fmt.Println("\nThroughput:")
	fmt.Printf("Requests por segundo (RPS):\t %.2f\n", rps)
	fmt.Printf("Requests por minuto (RPM):\t %.2f\n", rpm)
}

func printLatencyMetrics(m metrics.Metrics) {
	if len(m.Durations) == 0 {
		fmt.Println("\nSem métricas de latência (todas requisições falharam)")
		return
	}

	fmt.Println("\nLatência (microssegundos):")
	fmt.Printf("Média:\t\t %v\n", averageDuration(m.Durations).Round(time.Microsecond))
	fmt.Printf("Min:\t\t %v\n", minDuration(m.Durations).Round(time.Microsecond))
	fmt.Printf("Max:\t\t %v\n", maxDuration(m.Durations).Round(time.Microsecond))
	fmt.Printf("p50 (mediana):\t %v\n", percentile(m.Durations, 0.5))
	fmt.Printf("p90:\t\t %v\n", percentile(m.Durations, 0.9))
	fmt.Printf("p99:\t\t %v\n", percentile(m.Durations, 0.99))
}

// ... funções auxiliares para cálculos de latência ...

func averageDuration(durations []time.Duration) time.Duration {
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

func minDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	min := durations[0]
	for _, d := range durations {
		if d < min {
			min = d
		}
	}
	return min
}

func maxDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	max := durations[0]
	for _, d := range durations {
		if d > max {
			max = d
		}
	}
	return max
}

func percentile(durations []time.Duration, p float64) time.Duration {
	return metrics.NewCollector().CalculatePercentile(p)
}
