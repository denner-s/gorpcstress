package main

import (
	"log"
	"net"
	"net/rpc"
)

type Arithmetic struct{}

type Args struct {
	A, B int
}

type Reply struct {
	Result int
}

func (t *Arithmetic) Multiply(args *Args, reply *Reply) error {
	reply.Result = args.A * args.B
	return nil
}

func main() {
	arith := new(Arithmetic)
	err := rpc.Register(arith)
	if err != nil {
		return
	}

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Erro ao escutar:", err)
	}

	log.Println("Servidor RPC ouvindo na porta 1234")
	rpc.Accept(listener)
}
