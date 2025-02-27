package report

// Importação de pacotes necessários.
import (
	"fmt"  // Pacote para formatação de strings.
	"sort" // Pacote para ordenação de slices.
	"time" // Pacote para manipulação de tempo e durações.

	// Dependência interna do projeto.
	"github.com/denner-s/gorpcstress/internal/metrics" // Métricas coletadas durante o teste.
)

// Função GenerateReport gera e exibe um relatório de desempenho com base nas métricas coletadas.
func GenerateReport(m metrics.Metrics) {
	fmt.Println("\n=== Relatório do Teste de Estresse ===")

	// Exibe informações gerais sobre o teste.
	printGeneralInfo(m)

	// Exibe métricas de throughput (taxa de transferência).
	printThroughput(m)

	// Exibe métricas de latência.
	printLatencyMetrics(m)
}

// Função printGeneralInfo exibe informações gerais sobre o teste.
func printGeneralInfo(m metrics.Metrics) {
	// Calcula a duração total do teste.
	totalDuration := m.EndTime.Sub(m.StartTime)

	// Exibe a duração total, o número total de requisições e o número de erros.
	fmt.Printf("Tempo total de execução:\t %v\n", totalDuration.Round(time.Millisecond))
	fmt.Printf("Requisições totais:\t\t %d\n", m.TotalRequests)
	fmt.Printf("Requisições com erro:\t\t %d (%.2f%%)\n",
		m.Errors, errorRate(m))
}

// Função errorRate calcula a taxa de erro em porcentagem.
func errorRate(m metrics.Metrics) float64 {
	return float64(m.Errors) / float64(m.TotalRequests) * 100
}

// Função printThroughput exibe métricas de throughput (requisições por segundo e por minuto).
func printThroughput(m metrics.Metrics) {
	// Converter duração para segundos com precisão
	duration := m.EndTime.Sub(m.StartTime)
	seconds := duration.Seconds()

	var rps, rpm float64
	if seconds > 0 {
		rps = float64(m.TotalRequests) / seconds
		rpm = rps * 60
	} else {
		rps = 0
		rpm = 0
	}

	fmt.Println("\nThroughput:")
	fmt.Printf("Requests por segundo (RPS):\t %.2f\n", rps)
	fmt.Printf("Requests por minuto (RPM):\t %.2f\n", rpm)
}

// Função percentile calcula o percentil das durações das requisições.
func percentile(durations []time.Duration, p float64) time.Duration {
	if len(durations) == 0 {
		return 0 // Retorna 0 se não houver durações registradas.
	}

	// Ordena as durações em ordem crescente.
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	// Calcula o índice correspondente ao percentil desejado.
	index := int(float64(len(durations)-1) * p) // Corrigido para usar len-1
	return durations[index]                     // Retorna a duração no índice calculado.
}

// Função printLatencyMetrics exibe métricas de latência.
func printLatencyMetrics(m metrics.Metrics) {
	// Verifica se há durações registradas.
	if len(m.Durations) == 0 {
		fmt.Println("\nSem métricas de latência (todas requisições falharam)")
		return
	}

	// Cria uma cópia ordenada das durações.
	sorted := make([]time.Duration, len(m.Durations))
	copy(sorted, m.Durations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	// Exibe as métricas de latência.
	fmt.Println("\nLatência (microssegundos):")
	fmt.Printf("Média:\t\t %v\n", averageDuration(sorted).Round(time.Microsecond))
	fmt.Printf("Min:\t\t %v\n", sorted[0].Round(time.Microsecond))
	fmt.Printf("Max:\t\t %v\n", sorted[len(sorted)-1].Round(time.Microsecond))
	fmt.Printf("p50 (mediana):\t %v\n", percentile(sorted, 0.5))
	fmt.Printf("p90:\t\t %v\n", percentile(sorted, 0.9))
	fmt.Printf("p99:\t\t %v\n", percentile(sorted, 0.99))
}

// Função averageDuration calcula a duração média das requisições.
func averageDuration(durations []time.Duration) time.Duration {
	var total time.Duration
	for _, d := range durations {
		total += d // Soma todas as durações.
	}
	return total / time.Duration(len(durations)) // Retorna a média.
}
