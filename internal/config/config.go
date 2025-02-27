package config

// Importação de pacotes necessários.
import (
	"flag" // Pacote para manipulação de flags de linha de comando.
	"fmt"  // Pacote para formatação de strings e mensagens de erro.
	"time" // Pacote para manipulação de tempo e durações.
)

// Estrutura Config armazena todas as configurações necessárias para o teste de estresse.
type Config struct {
	ServerAddress string        // Endereço do servidor RPC (ex: "localhost:1234").
	TotalRequests int           // Número total de requisições a serem enviadas.
	Concurrency   int           // Número de workers concorrentes (goroutines).
	RPCMethod     string        // Método RPC a ser chamado (ex: "Arithmetic.Multiply").
	Timeout       time.Duration // Timeout para as conexões com o servidor.
	Duration      time.Duration // Duração total do teste (opcional, sobrescreve TotalRequests).
	PayloadFile   string        // Caminho para um arquivo JSON com payload customizado (opcional).
}

// Função LoadConfig carrega as configurações a partir de flags de linha de comando.
func LoadConfig() *Config {
	// Cria uma nova instância da estrutura Config.
	cfg := &Config{}

	// Define as flags de linha de comando e as associa aos campos da estrutura Config.
	flag.StringVar(&cfg.ServerAddress, "server", "localhost:1234", "Endereço do servidor RPC")
	flag.IntVar(&cfg.TotalRequests, "requests", 1000, "Número total de requisições")
	flag.IntVar(&cfg.Concurrency, "concurrency", 50, "Número de workers concorrentes")
	flag.StringVar(&cfg.RPCMethod, "method", "Arithmetic.Multiply", "Método RPC a ser chamado")
	flag.DurationVar(&cfg.Timeout, "timeout", 30*time.Second, "Timeout das conexões")
	flag.DurationVar(&cfg.Duration, "duration", 0, "Duração do teste (sobrescreve requests)")
	flag.StringVar(&cfg.PayloadFile, "payload", "", "Arquivo JSON com payload customizado")

	// Processa as flags fornecidas na linha de comando.
	flag.Parse()

	// Retorna a estrutura Config preenchida com os valores das flags.
	return cfg
}

// Método Validate verifica se as configurações carregadas são válidas.
func (c *Config) Validate() error {
	// Verifica se o número de requisições é maior que zero.
	if c.TotalRequests < 1 {
		return fmt.Errorf("número de requisições deve ser maior que zero")
	}

	// Verifica se o número de workers concorrentes é maior que zero.
	if c.Concurrency < 1 {
		return fmt.Errorf("concorrência deve ser maior que zero")
	}

	// Verifica se o endereço do servidor foi fornecido.
	if c.ServerAddress == "" {
		return fmt.Errorf("endereço do servidor não pode estar vazio")
	}

	if c.RPCMethod == "" {
		return fmt.Errorf("método RPC não pode ser vazio")
	}

	// Retorna nil se todas as validações forem bem-sucedidas.
	return nil
}
