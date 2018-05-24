package gorpc

import (
	"errors"
	"net"
	"time"
)

// for client get server addr by router
type Router interface {
	// GetAddr: get addr for call
	// input:@retry cur retry times, @begintime req time
	// output: address, timeout, reportTime(for report), routetype(for report), err
	GetAddr(retry int, beginTime time.Time) (net.Addr, time.Duration, time.Time, int, error)

	// get socket conf, output: keeptime(ms), retyrtimes, err
	//	keeptime: 0: short connection, other: free conntion keep time
	//	retrytimes:if net error, retrytimes
	GetConf() (uint32, int, error)

	// Report: report result
	// input:@addr address get by GetAddr, 
	//		@routetype get by GetAddr, 
	//		@code return for call
	//		@beginTime req time
	Report(addr net.Addr, routetype, code int, reportTime time.Time)

	// Load From configure
	Load(cfg Configure, section string) error
}

// get net.Addr by nettype andy address
func getaddr(nettype, addr string) (raddr net.Addr, err error) {
	switch nettype {
	case "tcp", "tcp4", "tcp6":
		raddr, err = net.ResolveTCPAddr(nettype, addr)
	case "udp", "udp4", "udp6":
		raddr, err = net.ResolveUDPAddr(nettype, addr)
	case "unix", "unixgram", "unixpacket":
		raddr, err = net.ResolveUnixAddr(nettype, addr)
	case "ip", "ip4", "ip6":
		raddr, err = net.ResolveIPAddr(nettype, addr)
	default: // error
		err = errors.New("no such nettype")
	}
	return

}

type AddrRoute struct {
	rddr	net.Addr
	timeout time.Duration // ms
	keeptime uint32 // ms
	retry	int
}

func (m * AddrRoute) GetAddr(retry int, beginTime time.Time) (net.Addr, time.Duration, time.Time, int, error) {
	return m.rddr, m.timeout, beginTime, 0, nil
}

func (m * AddrRoute) GetConf() (uint32, int, error) {
	return m.keeptime, m.retry, nil
}

func (m * AddrRoute) Report(addr net.Addr, routetype, code int, beginTime time.Time) {
}

// Load From configure
func (m * AddrRoute) Load(conf Configure, apiname string) (err error) {
	cfg := NewConfigureWape(conf, apiname)

	timeout := cfg.GetDefaultUint32("timeout", 100) // ms
	m.timeout = time.Duration(timeout) * time.Millisecond

	m.keeptime = cfg.GetDefaultUint32("keeptime", 30000) // ms
	m.retry = cfg.GetDefaultInt("retry", 1)

	addr := cfg.GetDefaultString("addr", "")
	nettype := cfg.GetDefaultString("nettype", "tcp")
	m.rddr, err = getaddr(nettype, addr)
	return
}
func NewAddrRoute(nettype, addr string, timeout, keeptime uint32, retry int) (Router, error) {
	raddr := &AddrRoute{timeout:(time.Duration(timeout) * time.Millisecond), keeptime:keeptime, retry:retry}
	var err error
	raddr.rddr, err = getaddr(nettype, addr)
	return raddr, err
}
