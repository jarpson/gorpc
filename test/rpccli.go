package main

import (
	gorpc ".."
	"fmt"
	//"time"
)

func main() {
	route, _ := gorpc.NewAddrRoute("tcp", "127.0.0.1:8081", 20, 20000, 0)
	cli := gorpc.NewRpcCli(route)
	code, recv, addr, err := cli.SendAndRecv([]byte("hello world\n"), gorpc.DefaultChecker)
	fmt.Printf("code:%d,addr:%v,err:%v,len:%d,[%s]\n", code, addr, err, len(recv), recv)
}
