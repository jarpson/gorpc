package gorpc

import (
	"net"
	//"time"
)

type Server struct {
	bindAddr net.Addr
	
}


// tcp/unix 
func (m * Server) serveStream(nettype, addr string) {
	//l, err := net.Listen(nettype, addr)
	
	// TODO 
}
