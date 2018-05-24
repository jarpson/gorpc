// go rpc logger
package gorpc

// logger interface for rpc frame
type Logger interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

// default logger with nothing output
type EmptyLog struct {
}

// logger interface
func (m *EmptyLog) Debugf(format string, v ...interface{}) {
}

// logger interface
func (m *EmptyLog) Errorf(format string, v ...interface{}) {
}

var emptylog EmptyLog
