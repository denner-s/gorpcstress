package main

import (
	"log"     // Para registro de logs
	"net"     // Para operações de rede TCP
	"net/rpc" // Para implementação do servidor RPC
)

// Args define a estrutura dos parâmetros de entrada das operações
type Args struct {
	A, B int // Operandos para as operações matemáticas
}

// Reply define a estrutura da resposta das operações
type Reply struct {
	Result int // Resultado da operação solicitada
}

// Arithmetic é a estrutura que implementa os métodos RPC
type Arithmetic struct{}

// Multiply implementa a operação de multiplicação via RPC
// Método exportado deve ser em maiúsculo e seguir a assinatura:
// func (t *T) MethodName(args *ArgType, reply *ReplyType) error
func (t *Arithmetic) Multiply(args *Args, reply *Reply) error {
	log.Printf("Recebida requisição: %d * %d", args.A, args.B)
	reply.Result = args.A * args.B
	return nil
}

func main() {
	// 1. Criação e registro do serviço RPC
	arithService := new(Arithmetic)

	// Registra explicitamente o serviço com o nome "Arithmetic"
	// Isso evita problemas de namespace e torna as chamadas mais precisas
	err := rpc.RegisterName("Arithmetic", arithService)
	if err != nil {
		log.Fatal("Falha ao registrar o serviço RPC:", err)
	}

	// 2. Configuração do listener TCP
	// Usamos 127.0.0.1 explicitamente para forçar IPv4
	listener, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("Erro ao iniciar listener:", err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Printf("⚠️ Erro ao fechar listener: %v", err)
		}
	}(listener) // Garante o fechamento adequado ao final

	log.Println("✅ Servidor RPC iniciado em 127.0.0.1:1234")

	// 3. Loop principal de aceitação de conexões
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Loga o erro mas mantém o servidor rodando
			log.Printf("⚠️ Erro na conexão: %v", err)
			continue
		}

		// 4. Tratamento concorrente da conexão
		go func(conn net.Conn) {
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {
					log.Printf("⚠️ Erro ao fechar conexão: %v", err)
				}
			}(conn) // Garante fechamento da conexão
			log.Printf("🔌 Nova conexão de %s", conn.RemoteAddr())

			// Serve a conexão usando o pacote RPC
			rpc.ServeConn(conn)
		}(conn)
	}
}
