package main

// Importação de dependências externas e internas
import (
	"fmt" // Pacote para formatação e impressão de textos
	"log" // Pacote para registro de logs

	// Dependências internas do projeto
	"github.com/denner-s/gorpcstress/internal/config"  // Manipulação de configurações
	"github.com/denner-s/gorpcstress/internal/metrics" // Coleta de métricas
	"github.com/denner-s/gorpcstress/internal/runner"  // Lógica de execução do teste de estresse
	"github.com/denner-s/gorpcstress/pkg/report"       // Geração de relatórios
)

// Função principal que será executada ao iniciar o programa
func main() {
	// Carrega a configuração do arquivo/configuração de ambiente
	cfg := config.LoadConfig()

	// Valida a configuração carregada. Se houver erro, encerra o programa com log
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuração inválida: %v", err) // log.Fatalf imprime mensagem e chama os.Exit(1)
	}

	// Cria um novo coletor de métricas para armazenar dados de desempenho
	collector := metrics.NewCollector()

	// Inicializa o executor de testes de estresse com a configuração e coletor
	stressRunner := runner.NewStressRunner(cfg, collector)

	// Exibe informações iniciais do teste formatadas
	fmt.Printf("Iniciando teste de estresse...\nServidor: %s\nRequisições: %d\nConcorrência: %d\n\n",
		cfg.ServerAddress, cfg.TotalRequests, cfg.Concurrency)

	// Executa efetivamente o teste de estresse
	stressRunner.Run()

	// Gera o relatório final com base nas métricas coletadas
	report.GenerateReport(collector.GetMetrics())
}
