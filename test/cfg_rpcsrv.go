package main

import (
	gorpc ".."
	"fmt"
	"github.com/vaughan0/go-ini"
)

func echo(r *gorpc.Request, data []byte) (code int, rsp []byte) {
	fmt.Printf("recv: %v,[%s]\n", r.Addr, data)
	return 0, data
}

func main() {
	cfg, err := ini.LoadFile("./cfg.ini")
	fmt.Println(cfg, err)
	srv := gorpc.Server{}
	srv.Load(cfg, "testsrv")
	srv.SetChecker(gorpc.DefaultChecker)
	srv.SetHandler(echo)
	fmt.Println(srv.Serve())
}
