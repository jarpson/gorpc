package main

import (
	"fmt"
	gorpc ".."
)

func echo (r *gorpc.Request, data []byte) (code int, rsp []byte) {
	fmt.Printf("recv: %v,[%s]\n", r.Addr, data)
	return 0, data
}

func main() {
	srv := gorpc.Server{}
	srv.SetAddr("tcp", "127.0.0.1:8081")
	srv.SetKeepTime(20000)
	srv.SetChecker(gorpc.DefaultChecker)
	srv.SetHandler(echo)
	fmt.Println(srv.Serve())
}
