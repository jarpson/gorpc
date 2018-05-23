package gorpc

import (
	"container/list"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const CLEAR_TIME = time.Duration(10) * time.Second

var (
	// global tcp pool
	g_addrMap = make(map[net.Addr]*tcpPool)
	g_lock = new(sync.RWMutex)
)

// free overtime tcp conn
func init() {
	go func() {
		// TODO clear g_addrMap
		fdList := list.New()
		addrList := []net.Addr{}
		for {
			time.Sleep(CLEAR_TIME)

			ts := time.Now()
			g_lock.RLock()

			// scan overtime connection and addr pool
			for k, v := range g_addrMap {
				v.LockRun(func() {
					for e := v.Clist.Back(); true; e = v.Clist.Back() {
						if e == nil {
							if v.LastOpTime.Sub(ts) > CLEAR_TIME {
								addrList = append(addrList, k)
							}
							break
						}
						if c, ok := e.Value.(connBase); ok {
							if c.CheckTimeAfter(ts) {
								fdList.PushBack(c)
								v.Clist.Remove(e)
							} else {
								break
							}
						}
					}
				})
			}
			g_lock.RUnlock()

			// clear connections
			clearFd(fdList)

			// clear addr pool
			if len(addrList) > 0 {
				g_lock.Lock()
				for _, v := range addrList {
					if tcp, ok := g_addrMap[v]; ok {
						tcp.Destory = true
						delete(g_addrMap, v)
					}
				}
				g_lock.Unlock()
				addrList = addrList[:0]
			}
		}
	}()
}

// tcp pool-node
type connBase struct {
	net.Conn
	Parent *tcpPool
	lastUsedTime time.Time
}

func NewConnBase(addr net.Addr, parent *tcpPool, dialTimeout time.Duration) (conn net.Conn, err error) {
	conn, err = net.DialTimeout(addr.Network(), addr.String(), dialTimeout)
	if err == nil {
		conn = &connBase{Conn:conn, Parent:parent}
	}
	return 
}

func (m *connBase) CheckTimeAfter(ts time.Time) bool {
	return m.lastUsedTime.After(ts)
}

// connBase.Close not close tcp connection but put back to pool list
// if pointer to pool list is nil, close connection
func (m *connBase) Close() error {
	if m.Parent == nil || m.Parent.Destory {
		m.Parent = nil
		m.Conn.Close()
		return nil
	}
	m.lastUsedTime = time.Now()
	pushfail := false
	m.Parent.LockRun(func() {
		if m.Parent == nil || m.Parent.Destory { // check again 
			pushfail = true
		} else {
			m.Parent.Clist.PushFront(m)
			m.Parent.LastOpTime = m.lastUsedTime
		}
	})
	if pushfail {
		m.Conn.Close()
	}
	return nil
}

// close connect and reset pointer to parent
func (m *connBase) Destory() {
	m.Parent = nil
	m.Close()
}

type tcpPool struct {
	addr net.Addr
	keeptime uint32  // ms

	Clist *list.List // pool list
	lock int32	// lock for Clist

	LastOpTime time.Time
	Destory bool
}

func newTcpPool(addr net.Addr, keeptime uint32) *tcpPool {
	return &tcpPool{addr:addr, keeptime:keeptime, Clist:list.New()}
}

// use min time, Millisecond
func (m *tcpPool) KeepTime(k uint32) {
	if m.keeptime > k {
		for !atomic.CompareAndSwapUint32(&m.keeptime, m.keeptime, k) {
			if m.keeptime <= k {
				break
			}
		}
	}
}

func (m *tcpPool) LockRun(f func()) {
	for  {
		if atomic.SwapInt32(&m.lock, 1) == 0 {
			f()
			break
		} else {
			runtime.Gosched()
		}
	}
	atomic.SwapInt32(&m.lock, 0)
}

// get net.Conn from tcp-pool
func (m *tcpPool) GetFd(addr net.Addr, dialTimeout time.Duration, beginTime time.Time) (conn net.Conn, err error) {
	tsAfter := beginTime.Add(-time.Duration(m.keeptime) * time.Millisecond)
	m.LockRun(func() {
		if e := m.Clist.Front(); e != nil {
			if c, ok := e.Value.(connBase); ok {
				if c.CheckTimeAfter(tsAfter) { // not over time
					conn = &c
					m.Clist.Remove(e)
				}
			}
		}
	})
	if conn == nil {
		conn, err = NewConnBase(addr, m, dialTimeout)
	}
	m.LastOpTime = beginTime
	return conn, err
}

func clearFd(fdList * list.List) {
	if fdList == nil {
		return
	}
	for e := fdList.Front(); e != nil; e = e.Next() {
		if c, ok := e.Value.(connBase); ok {
			c.Destory()
		}
	}
	fdList.Init()
}

// get rpc fd
// @addr: fd for addr
// @timeout: timeout r/w
// @keeptime: keepalive time/ms (just for tcp)
// return net.Conn, it should call CloseFd after use,unless error not nil
func GetFd(addr net.Addr, timeout time.Duration, keeptime uint32, beginTime time.Time) (conn net.Conn, err error) {
	if keeptime == 0 || (addr.Network() != "tcp" && addr.Network() != "tcp4" && addr.Network() != "tcp6") {
		conn, err = net.DialTimeout(addr.Network(), addr.String(), timeout)
	} else {
		g_lock.RLock()
		tcp, ok := g_addrMap[addr]
		g_lock.RUnlock()
		if !ok {	// not found, so add
			g_lock.Lock()
			if tcp, ok = g_addrMap[addr]; !ok { // check again
				tcp = newTcpPool(addr, keeptime)
				g_addrMap[addr] = tcp
			}
			g_lock.Unlock()
		} else {
			tcp.KeepTime(keeptime)
		}
		conn, err = tcp.GetFd(addr, timeout, beginTime)
	}
	if err == nil {
		conn.SetDeadline(beginTime.Add(timeout))
	}
	return conn, err
}

// close rpc fd
func CloseFd(c net.Conn, err error) {
	if c == nil {
		return
	}
	if err != nil {
		if cb, ok := c.(*connBase); ok {
			cb.Destory()
			return
		}
	}
	c.Close()
}

