package gorpc

import (
	"net"
	"runtime"
	"time"
)

// rpc server info
type Request struct {
	Addr net.Addr	// remote addr
	Ext interface{} // ext msg
}

// rpc callback handle
// input:
//	r: request info
//	data: request data
// output:
//	code: return value, code < 0 is error and close connection
//	rsp: data send to client
type Handler func(r *Request, data []byte) (code int, rsp []byte)

// rpc server struct
type Server struct {
	network string
	address string
	keeptime time.Duration
	checker Checker
	handler Handler
}

// set bind addr
func (m * Server) SetAddr(network, address string) {
	m.network = network
	m.address = address
}

// set free connection keet times(ms)
func (m * Server) SetKeepTime(keepms uint32) {
	m.keeptime = time.Duration(keepms) * time.Millisecond
}

// set packet checker
func (m * Server) SetChecker(checker Checker) {
	m.checker = checker
}

// set handle function
func (m * Server) SetHandler(handler Handler) {
	m.handler = handler
}

// deal goroutines
func (m * Server) handProc(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			// log.Printf("http: panic serving %v: %v\n%s", conn.RemoteAddr(), err, buf)
			if conn != nil {
				conn.Close()
				conn = nil
			}
		}
	}()

	// spport long tcp
	for {
		now := time.Now()
		conn.SetDeadline(now.Add(m.keeptime))
		code, data, err := RecvAll(conn, m.checker)
		var rsp []byte
		if err == nil {
			r := &Request{Addr:conn.RemoteAddr()}
			code, rsp = m.handler(r, data)
		}
		if len(rsp) > 0 {
			SendAll(conn, rsp)
		}
		switch code {
			// connection was closed by client
			case ERR_RECV:

			// deal ok
			case ERR_OK:

			// others
			default:
		}
		if code < 0 {
			conn.Close()
			conn = nil
			break;
		}
	}
}


// start rpc server
func (m * Server) Serve() error {
	listiner, err := net.Listen(m.network, m.address)
	if err != nil {
		return err
	}
	for {
		conn, err := listiner.Accept()
		if err != nil {
			continue
		}
		go m.handProc(conn)
	}
	return nil
}

// Load server by conf
func (m * Server) Load(conf Configure, apiname string) {
	cfg := NewConfigureWape(conf, apiname)
	nettype := cfg.GetDefaultString("nettype", "tcp")
	addr := cfg.GetDefaultString("bind", "")
	keeptime := cfg.GetDefaultUint32("keeptime", 30000) // ms
	m.SetAddr(nettype, addr)
	m.SetKeepTime(keeptime)
}

