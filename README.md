# GoRPCStress 🚀

Uma biblioteca para testes de estresse em endpoints RPC escritos em Go, com métricas detalhadas e alto desempenho.

## Recursos Principais

- ✅ Teste de carga com controle de concorrência
- ✅ Coleta detalhada de métricas de performance
- ✅ Cálculo de percentis (p50, p90, p99)
- ✅ Medição de throughput (RPS/RPM)
- ✅ Relatório completo em tempo real
- ✅ Configuração flexível via linha de comando
- ✅ Tratamento robusto de erros
- ✅ Conexões RPC otimizadas

## Instalação

```bash
go install github.com/denner-s/gorpcstress@latest
```

## Uso Básico

```bash
gorpcstress \
  -server=localhost:1234 \
  -requests=5000 \
  -concurrency=100 \
  -method=Service.Method
```

## Opções de Configuração

| Flag           | Descrição                          | Padrão               |
|----------------|------------------------------------|----------------------|
| `-server`      | Endereço do servidor RPC           | localhost:1234       |
| `-requests`    | Número total de requisições        | 1000                 |
| `-concurrency` | Número de workers concorrentes     | 50                   |
| `-method`      | Método RPC a ser testado           | Arithmetic.Multiply  |

## Exemplo de Saída

```text
=== Relatório do Teste de Estresse ===

Tempo total de execução:      2.45s
Requisições totais:           5000
Requisições com erro:         12 (0.24%)

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

## Funcionalidades Futuras (Roadmap)

- [ ] Suporte a HTTP/gRPC
- [ ] Modo de teste por duração
- [ ] Relatórios em JSON/CSV
- [ ] Monitoramento de recursos do sistema
- [ ] Validação de respostas personalizada
- [ ] Carga dinâmica com ramp-up

## Contribuição

Contribuições são bem-vindas! Siga estes passos:

1. Fork o repositório
2. Crie um branch com sua feature (`git checkout -b feature/incrivel`)
3. Commit suas mudanças (`git commit -am 'Add incrivel feature'`)
4. Push para o branch (`git push origin feature/incrivel`)
5. Abra um Pull Request

## Licença

Distribuído sob a licença MIT. Veja `LICENSE` para mais informações.

---

**Aviso**: Use com cuidado em ambientes de produção. Testes de estresse podem impactar a performance do sistema.
