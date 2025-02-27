package rpcclient

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

// Args define os parâmetros de entrada para operações RPC.
// Campos exportados (maiúsculos) garantem serialização correta.
type Args struct {
	A, B int // Operandos para cálculos (ex: multiplicação)
}

// Reply define a estrutura de resposta do servidor RPC.
type Reply struct {
	Result int // Resultado da operação
}

// Client encapsula uma conexão RPC com timeout.
type Client struct {
	*rpc.Client
	Timeout time.Duration // Tempo máximo para conexão/chamadas
}

// NewClient estabelece uma conexão com o servidor RPC.
// - `serverAddress`: Endereço no formato "host:porta"
// - `timeout`: Tempo máximo de espera por conexão
func NewClient(serverAddress string, timeout time.Duration) (*Client, error) {
	conn, err := net.DialTimeout("tcp", serverAddress, timeout)
	if err != nil {
		return nil, fmt.Errorf("falha na conexão: %w", err) // Erro detalhado
	}
	return &Client{
		Client:  rpc.NewClient(conn),
		Timeout: timeout,
	}, nil
}

// Call executa uma chamada RPC com controle de timeout.
// - Usa goroutine + channel para evitar bloqueio indefinido
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	done := make(chan error, 1)
	go func() { done <- c.Client.Call(serviceMethod, args, reply) }()

	select {
	case err := <-done:
		return err // Retorna erro imediatamente se houver
	case <-time.After(c.Timeout):
		return fmt.Errorf("timeout após %v", c.Timeout) // Erro customizado
	}
}
