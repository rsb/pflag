package pflag

import (
	"fmt"
	"net"
	"strings"
)

// IPNet adapts net.IPNet for use as a flag.
type ipNetValue net.IPNet

func (ipnet ipNetValue) String() string {
	n := net.IPNet(ipnet)
	return n.String()
}

func (ipnet *ipNetValue) Set(value string) error {
	_, n, err := net.ParseCIDR(strings.TrimSpace(value))
	if err != nil {
		return err
	}
	*ipnet = ipNetValue(*n)
	return nil
}

func (*ipNetValue) Type() string {
	return "ipNet"
}

func newIPNetValue(val net.IPNet, p *net.IPNet) *ipNetValue {
	*p = val
	return (*ipNetValue)(p)
}

func ipNetConv(sval string) (interface{}, error) {
	_, n, err := net.ParseCIDR(strings.TrimSpace(sval))
	if err == nil {
		return *n, nil
	}
	return nil, fmt.Errorf("invalid string being converted to IPNet: %s", sval)
}
