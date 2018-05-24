package main

import (
	"fmt"
	"github.com/vaughan0/go-ini"
	gorpc ".."
	//"time"
)

func main() {
	cfg, err := ini.LoadFile("./cfg.ini")
	fmt.Println(cfg, err)
	cli := gorpc.NewRpcCli(&gorpc.AddrRoute{})
	cli.Load(cfg, "testapi")
	code, recv, addr, err := cli.SendAndRecv([]byte("hello world\n"), gorpc.DefaultChecker)
	fmt.Printf("code:%d,addr:%v,err:%v,len:%d,[%s]\n", code, addr, err, len(recv),recv)
}

