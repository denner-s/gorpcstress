package metrics

// Importação de pacotes necessários.
import (
	"sort" // Pacote para ordenação de slices.
	"time" // Pacote para manipulação de tempo e durações.
)

// Estrutura Result armazena o resultado de uma requisição individual.
type Result struct {
	Duration time.Duration // Duração da requisição.
	Error    error         // Erro (se houver) durante a requisição.
}

// Estrutura Collector gerencia a coleta de métricas de todas as requisições.
type Collector struct {
	metrics Metrics // Armazena as métricas coletadas.
}

// Estrutura Metrics armazena os dados agregados das requisições.
type Metrics struct {
	TotalRequests int             // Número total de requisições.
	Errors        int             // Número de requisições que falharam.
	Durations     []time.Duration // Lista de durações das requisições bem-sucedidas.
	StartTime     time.Time       // Timestamp de início da coleta de métricas.
	EndTime       time.Time       // Timestamp de término da coleta de métricas.
}

// Função NewCollector cria e inicializa uma nova instância de Collector.
func NewCollector() *Collector {
	return &Collector{
		metrics: Metrics{
			Durations: make([]time.Duration, 0), // Inicializa a lista de durações vazia.
		},
	}
}

// Método RecordResult registra o resultado de uma requisição no coletor.
func (c *Collector) RecordResult(result Result) {
	c.metrics.TotalRequests++ // Incrementa o contador de requisições totais.

	if result.Error != nil {
		c.metrics.Errors++ // Incrementa o contador de erros se houver um erro.
	} else {
		c.metrics.Durations = append(c.metrics.Durations, result.Duration) // Adiciona a duração à lista de durações.
	}
}

// Método GetMetrics retorna as métricas coletadas.
func (c *Collector) GetMetrics() Metrics {
	metrics := c.metrics
	// Garante duração mínima de 1 nanossegundo para evitar divisão por zero
	if metrics.EndTime.Before(metrics.StartTime.Add(1 * time.Nanosecond)) {
		metrics.EndTime = metrics.StartTime.Add(1 * time.Nanosecond)
	}
	return metrics
}

// Método CalculatePercentile calcula o percentil das durações das requisições.
func (c *Collector) CalculatePercentile(p float64) time.Duration {
	if len(c.metrics.Durations) == 0 {
		return 0 // Retorna 0 se não houver durações registradas.
	}

	// Ordena as durações em ordem crescente.
	sort.Slice(c.metrics.Durations, func(i, j int) bool {
		return c.metrics.Durations[i] < c.metrics.Durations[j]
	})

	// Calcula o índice correspondente ao percentil desejado.
	index := int(float64(len(c.metrics.Durations)) * p)
	if index >= len(c.metrics.Durations) {
		index = len(c.metrics.Durations) - 1 // Garante que o índice não ultrapasse o tamanho da lista.
	}
	return c.metrics.Durations[index] // Retorna a duração no índice calculado.
}

// Método AverageDuration calcula a duração média das requisições bem-sucedidas.
func (c *Collector) AverageDuration() time.Duration {
	if len(c.metrics.Durations) == 0 {
		return 0 // Retorna 0 se não houver durações registradas.
	}

	// Soma todas as durações.
	var total time.Duration
	for _, d := range c.metrics.Durations {
		total += d
	}

	// Retorna a média (soma das durações dividida pelo número de durações).
	return total / time.Duration(len(c.metrics.Durations))
}
