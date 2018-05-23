package gorpc

import (
	"net"
	"time"
)

type Router interface {
	// GetAddr: get addr for call
	// input:@retry cur retry times, @begintime req time
	// output: address, timeout, reportTime(for report), routetype(for report), err
	GetAddr(retry int, beginTime time.Time) (net.Addr, time.Duration, time.Time, int, error)

	// GetConf: client / server conf 
	// 1. client: get socket conf
	//		output: keeptime(ms), retyrtimes, err
	// 2. server: get server conf
	//		output: keeptimes, maxgoroutes, 
	GetConf() (uint32, int, error)

	// Report: report result
	// input:@addr address get by GetAddr, 
	//		@routetype get by GetAddr, 
	//		@code return for call
	//		@beginTime req time
	Report(addr net.Addr, routetype, code int, reportTime time.Time)
}

type AddrRoute struct {
	Addr	net.Addr
	Timeout time.Duration // ms
	KeepTime uint32 // ms
	Retry	int
}

func (m * AddrRoute) GetAddr(retry int, beginTime time.Time) (net.Addr, time.Duration, time.Time, int, error) {
	return m.Addr, m.Timeout, beginTime, 0, nil
}

func (m * AddrRoute) GetConf() (uint32, int, error) {
	return m.KeepTime, m.Retry, nil
}

func (m * AddrRoute) Report(addr net.Addr, routetype, code int, beginTime time.Time) {
}

func NewAddrRoute(nettype, addr string, timeout, keeptime uint32, retry int) (Router, error) {
	raddr := &AddrRoute{Timeout:(time.Duration(timeout) * time.Millisecond), KeepTime:keeptime, Retry:retry}
	var err error
	switch nettype {
	case "tcp", "tcp4", "tcp6":
		raddr.Addr, err = net.ResolveTCPAddr(nettype, addr)
	case "udp", "udp4", "udp6":
		raddr.Addr, err = net.ResolveUDPAddr(nettype, addr)
	case "unix", "unixgram", "unixpacket":
		raddr.Addr, err = net.ResolveUnixAddr(nettype, addr)
	default: // "ip", "ip4", "ip6":
		raddr.Addr, err = net.ResolveIPAddr(nettype, addr)
	}
	return raddr, err
}
