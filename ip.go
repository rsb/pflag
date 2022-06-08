package pflag

import (
	"fmt"
	"net"
	"strings"
)

// -- net.IP value
type ipValue net.IP

func newIPValue(val net.IP, p *net.IP) *ipValue {
	*p = val
	return (*ipValue)(p)
}

func (i *ipValue) String() string { return net.IP(*i).String() }
func (i *ipValue) Set(s string) error {
	ip := net.ParseIP(strings.TrimSpace(s))
	if ip == nil {
		return fmt.Errorf("failed to parse IP: %q", s)
	}
	*i = ipValue(ip)
	return nil
}

func (i *ipValue) Type() string {
	return "ip"
}

func ipConv(sval string) (interface{}, error) {
	ip := net.ParseIP(sval)
	if ip != nil {
		return ip, nil
	}
	return nil, fmt.Errorf("invalid string being converted to IP address: %s", sval)
}

// GetIP return the net.IP value of a flag with the given name
func (f *FlagSet) GetIP(name string) (net.IP, error) {
	val, err := f.getFlagType(name, "ip", ipConv)
	if err != nil {
		return nil, err
	}
	return val.(net.IP), nil
}
