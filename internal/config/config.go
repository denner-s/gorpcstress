package config

import (
	"flag"
	"fmt"
)

type Config struct {
	ServerAddress string
	TotalRequests int
	Concurrency   int
	RPCMethod     string
}

func LoadConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "server", "localhost:1234", "Endereço do servidor RPC")
	flag.IntVar(&cfg.TotalRequests, "requests", 1000, "Número total de requisições")
	flag.IntVar(&cfg.Concurrency, "concurrency", 50, "Número de workers concorrentes")
	flag.StringVar(&cfg.RPCMethod, "method", "Arithmetic.Multiply", "Método RPC a ser chamado")

	flag.Parse()
	return cfg
}

func (c *Config) Validate() error {
	if c.TotalRequests < 1 {
		return fmt.Errorf("número de requisições deve ser maior que zero")
	}
	if c.Concurrency < 1 {
		return fmt.Errorf("concorrência deve ser maior que zero")
	}
	if c.ServerAddress == "" {
		return fmt.Errorf("endereço do servidor não pode estar vazio")
	}
	return nil
}
