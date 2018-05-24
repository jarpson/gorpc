package gorpc

import "fmt"

const (
	ERR_OK         = 0
	ERR_CONF       = -1
	ERR_ROUTE      = -2
	ERR_CONN       = -3
	ERR_SEND       = -4
	ERR_RECV       = -6
	ERR_CHECK      = -7
	ERR_PACK       = -8
	ERR_UNPACK     = -9
	ERR_NODEAL_FUN = -10
)

var (
	INIT_PACKATE_LEN = 4096
	MAX_PACKATE_LEN  = 65536
)

var (
	//ERROR_OVERLOAD = fmt.Errorf("overload")
	ERROR_PACK_CHECK = fmt.Errorf("package check fail")
	ERROR_PACK_LONG  = fmt.Errorf("package too large")
	ERROR_NO_NETTYPE = fmt.Errorf("no such net")
	ERROR_UNKNOW     = fmt.Errorf("unknow error")
)
