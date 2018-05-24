// go rpc config

package gorpc

import (
	"strconv"
)

// configure for init server/client
// you can use "github.com/vaughan0/go-ini"
type Configure interface {
	Get(section, name string) (string, bool)
}

// extend configure: get number value, set default value
type ConfigureWape struct {
	Configure
	section string
}

// create ConfigureWape by Configure
// input:
//	conf: Configure
//	section: section of ini file
// output: ConfigureWape
func NewConfigureWape(conf Configure, section string) *ConfigureWape {
	return &ConfigureWape{conf, section}
}

// Get string config with default value
func (m *ConfigureWape) GetDefaultString(name, def string) string {
	str, ok := m.Get(m.section, name)
	if !ok {
		return def
	}
	return str
}

// Get int config with default value
func (m *ConfigureWape) GetDefaultInt(name string, def int) int {
	str, ok := m.Get(m.section, name)
	if !ok {
		return def
	}
	if v, err := strconv.Atoi(str); err == nil {
		return v
	}
	return def
}

// Get uint32 config with default value
func (m *ConfigureWape) GetDefaultUint32(name string, def uint32) uint32 {
	str, ok := m.Get(m.section, name)
	if !ok {
		return def
	}
	if v, err := strconv.ParseUint(str, 10, 64); err == nil {
		return uint32(v)
	}
	return def
}
