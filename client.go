package gorpc

import (
	"net"
	"time"
)

// rpc client
type RpcCli struct {
	Router // router, for get conf and addr
}

// create new rpc client
func NewRpcCli(route Router) *RpcCli {
	return &RpcCli{route}
}

// send request and recv full package
// input: @data: msg to send, @check: check package recved
// output:
//		@code: return code, see:define.go, if request ok return 0
//		@recv: recved data
//		@addr: request server addr
//		@err: error message
func (r *RpcCli) SendAndRecv(data []byte, check Checker) (code int, recv []byte, addr net.Addr, err error) {
	beginTime := time.Now()
	keeptime, retry, err := r.GetConf()
	if err != nil {
		code = ERR_CONF
		return
	}
	trytimes := retry + 1
	if trytimes <= 0 {
		trytimes = 1
	}
	var timeout time.Duration
	var routetype int
	var reportTime time.Time
	for try := 0; try < trytimes; try++ {
		addr, timeout, reportTime, routetype, err = r.GetAddr(try, beginTime)
		if err != nil {
			code = ERR_ROUTE
			r.Router.Report(addr, routetype, code, reportTime)
			continue
		}
		var conn net.Conn
		conn, err = GetFd(addr, timeout, keeptime, beginTime)
		if err == nil {
			defer CloseFd(conn, err)
			// send message
			code, err = SendAll(conn, data)
			if code == ERR_OK {
				// recv message
				code, recv, err = RecvAll(conn, check)
			}
		} else {
			code = ERR_CONN
		}
		r.Router.Report(addr, routetype, code, reportTime)
		if code == ERR_OK {
			break
		}
	}
	return
}

func (r *RpcCli) Load(conf Configure, apiname string) (err error) {
	if r.Router == nil {
		r.Router = &AddrRoute{}
	}
	return r.Router.Load(conf, apiname)
}
