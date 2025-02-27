# GoRPCStress 🚀

Uma biblioteca profissional para testes de estresse em endpoints RPC escritos em Go, com métricas detalhadas e alto desempenho.

## Recursos Principais

- ✅ Teste de carga com controle de concorrência
- ✅ Coleta detalhada de métricas de performance
- ✅ Cálculo de percentis (p50, p90, p99)
- ✅ Medição de throughput (RPS/RPM)
- ✅ Relatório completo em tempo real
- ✅ Configuração flexível via linha de comando
- ✅ Tratamento robusto de erros e reconexão automática
- ✅ Conexões RPC otimizadas e pool de clientes
- ✅ Validação automática de respostas
- ✅ Suporte a diferentes tipos de payloads

## Instalação

```bash
# Instalação global
go install github.com/denner-s/gorpcstress@latest

# Ou usando como módulo
go mod init seu-projeto
go get github.com/denner-s/gorpcstress
```

## Uso Básico

```bash
gorpcstress \
  -server=localhost:1234 \
  -requests=5000 \
  -concurrency=100 \
  -method=Service.Method
```

## Exemplo Completo

**Servidor de Teste (server.go):**
```go
package main

import (
	"net"
	"net/rpc"
)

type Calculator struct{}

func (c *Calculator) Multiply(args *struct{A, B int}, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func main() {
	rpc.Register(new(Calculator))
	listener, _ := net.Listen("tcp", ":1234")
	rpc.Accept(listener)
}
```

**Executando o Teste:**
```bash
./bin/gorpcstress -server=localhost:1234 \
  -requests=5000 \
  -concurrency=100 \
  -method=Calculator.Multiply
```
ou
```bash
./bin/gorpcstress   -server=localhost:1234   -method=Calculator.Multiply   -requests=1000   -concurrency=50   -timeout=5s
```

## Opções de Configuração

| Flag           | Descrição                          | Padrão               |
|----------------|------------------------------------|----------------------|
| `-server`      | Endereço do servidor RPC           | localhost:1234       |
| `-requests`    | Número total de requisições        | 1000                 |
| `-concurrency` | Número de workers concorrentes     | 50                   |
| `-method`      | Método RPC a ser testado           | Arithmetic.Multiply  |
| `-timeout`     | Timeout por requisição (opcional)  | 10s                  |

## Exemplo de Saída

```text
=== Relatório do Teste de Estresse ===

Tempo total de execução:      2.45s
Requisições totais:           5000
Requisições com erro:         23 (0.46%)
  • 10 erros de conexão
  • 5 timeouts
  • 8 respostas inválidas

Throughput:
Requests por segundo (RPS):   2040.82
Requests por minuto (RPM):    122449.22

Latência (microssegundos):
Média:        1.2ms
Min:          0.8ms
Max:          15.4ms
p50 (mediana): 1.1ms
p90:          1.5ms
p99:          3.8ms
```

## Uso Avançado

**Teste com Payload Customizado:**
1. Modifique as structs no arquivo `pkg/rpcclient/client.go`:
```go
type Args struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Reply struct {
	Total int `json:"total"`
}
```

2. Execute com:
```bash
./bin/gorpcstress -method=Calculator.Sum -requests=2000 -concurrency=75
```

**Teste de Duração:**
```bash
# Executar por 5 minutos
./bin/gorpcstress -concurrency=100 -duration=5m
```

## Solução de Problemas Comuns

**Erro: "Too many open files"**
```bash
# Aumente o limite de arquivos
ulimit -n 10000
```

**Erro: "Connection refused"**
```bash
# Verifique se o servidor está aceitando conexões
telnet localhost 1234

# Verifique firewalls
sudo ufw allow 1234/tcp
```

## Funcionalidades Futuras (Roadmap)

- [ ] Suporte a HTTP/gRPC
- [x] Modo de teste por duração
- [ ] Relatórios em JSON/CSV
- [x] Validação de respostas
- [ ] Monitoramento de recursos do sistema
- [ ] Carga dinâmica com ramp-up
- [ ] Teste distribuído em múltiplos nós
- [ ] Geração de gráficos de performance
- [ ] Suporte a payloads customizados

## Contribuição

Siga estes passos para contribuir:

1. Fork o repositório
2. Crie um branch descritivo (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -am 'Adiciona incrível funcionalidade'`)
4. Push para o branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request com detalhes das mudanças

**Requisitos para contribuição:**
- Documentação atualizada
- Código seguindo as diretrizes do [Effective Go](https://go.dev/doc/effective_go)

## Licença

Distribuído sob a licença MIT. Veja [LICENSE](LICENSE) para mais informações.

---

**Aviso Importante:**  
⚠️ Use com cuidado em ambientes de produção  
⚠️ Monitore o servidor alvo durante os testes  
⚠️ Configure timeouts adequados para suas necessidades  
⚠️ Não utilize em sistemas críticos sem autorização

**Dica profissional:**  
Para melhores resultados, execute primeiro um teste de aquecimento com 10% da carga total antes do teste principal.
