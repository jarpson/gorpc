package gorpc

import (
	"net"
)

// Checker: check recv package
// input: []byte of recv data
// output:
//    return 1: check result(<0:package err, 0: not over, >0:ok and len of package
//    return 2: full package len(0: unknown) include head of package
// notice: if input []byte is empty, return 0, length_of_head
type Checker func([]byte) (int, int)

// default checker :no check package
var DefaultChecker Checker = func(d []byte) (int, int) { return len(d), len(d) }

// recv all data
func RecvAll(c net.Conn, check Checker) (code int, recv []byte, err error) {
	if check == nil {
		check = DefaultChecker
	}

	var start, status int
	_, headlen := check([]byte{})
	if headlen <= 0 { // if no package_head
		headlen = INIT_PACKATE_LEN
	}
	buf := GetBufN(headlen)
	max := len(buf)
	for {
		buflen := len(buf)
		var recvlen int
		if max > start { // if we know package len
			recvlen, err = c.Read(buf[start:max])
		} else {
			recvlen, err = c.Read(buf[start:])
		}
		if err != nil { // over time  or connection closed
			code = ERR_RECV
			return
		}
		start += recvlen
		status, max = check(buf[:start])
		if status < 0 {
			code, err = ERR_CHECK, ERROR_PACK_CHECK
			return
		}
		if status > 0 {
			code, err = ERR_OK, nil
			recv = buf[:status]
			return // success
		}

		if max > MAX_PACKATE_LEN {
			code, err = ERR_CHECK, ERROR_PACK_LONG
		}
		// stat == 0 : go on reading

		if max > buflen { // full package len > cur buffer len
			buf = ResizeBuf(buf, max) // expand buffer to max
		} else if max == 0 && start == buflen { // cannot read head && read buffer full
			buf = ResizeBuf(buf, 2*buflen) // expand buffer to 2 * buffer
		}
	}
	code, err = ERR_RECV, ERROR_UNKNOW
	return
}

// send all data
// input:
//		@c connector
//		@data data for send
// return: code, err
func SendAll(c net.Conn, data []byte) (int, error) {
	var sum int
	total := len(data)
	for sum < total {
		slen, err := c.Write(data)
		if err != nil {
			return ERR_SEND, err
		}
		sum += slen
	}
	return ERR_OK, nil
}
