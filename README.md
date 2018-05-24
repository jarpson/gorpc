# gorpc
go rpc with connection pool and server frame.
go 版本rpc框架，包括连接池、客户端、服务端
可以配置tcp/udp/unixsocket，设置长短连接，及连接保持时间
使用方法（见test目录下的样例）：
服务端
```go
//server
package main

import (
    "github.com/jarpson/gorpc"
    "log"

// handle function
func echo(r *gorpc.Request, data []byte) (code int, rsp []byte) {
	fmt.Printf("recv: %v,[%s]\n", r.Addr, data)
	return 0, data
}

func main() {
	srv := gorpc.Server{}
  
  // direct set
	srv.SetAddr("tcp", "127.0.0.1:8081")
	srv.SetKeepTime(20000)
  //
  // you can also load from config, see: https://github.com/jarpson/gorpc/blob/master/test/cfg_rpcsrv.go
  // srv.Load(cfg, "testsrv")
  
  // cheker packet
	srv.SetChecker(gorpc.DefaultChecker)
  
  regist hande function
	srv.SetHandler(echo)
  
  //start server
	log.Println(srv.Serve())
}
```

客户端
```go
// client
package main

import (
	gorpc ".."
	"log"
)

func main() {
  // init client api
	route, _ := gorpc.NewAddrRoute("tcp", "127.0.0.1:8081", 20, 20000, 0)
	cli := gorpc.NewRpcCli(route)
  // 
  // you can also load from config, see : https://github.com/jarpson/gorpc/blob/master/test/cfg_rpccli.go
  // cli.Load(cfg, "testapi")
  
  // send && recv with connection pool
	code, recv, addr, err := cli.SendAndRecv([]byte("hello world\n"), gorpc.DefaultChecker)
	log.Printf("code:%d,addr:%v,err:%v,len:%d,[%s]\n", code, addr, err, len(recv), recv)
}

```
配置方式如下：
```ini
[testapi]
; api request addr
addr=127.0.0.1:8081

; bind nettype: tcp/udp/unix/ip
nettype=tcp

; free connection keep time
keeptime=2000

; overtime
timeout=100

; retry if fail
retry=0

[testsrv]
; server bind addr
bind=127.0.0.1:8081

; server bind type
nettype=tcp

; free connection keep time
keeptime=2000
```
