package pflag_test

import (
	"fmt"
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/require"
	"net"
	"os"
	"testing"
)

func setUpIP(ip *net.IP) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.IPVar(ip, "address", net.ParseIP("0.0.0.0"), "IP Address")
	return f
}

func TestIP(t *testing.T) {
	testCases := []struct {
		input    string
		success  bool
		expected string
	}{
		{"0.0.0.0", true, "0.0.0.0"},
		{" 0.0.0.0 ", true, "0.0.0.0"},
		{"1.2.3.4", true, "1.2.3.4"},
		{"127.0.0.1", true, "127.0.0.1"},
		{"255.255.255.255", true, "255.255.255.255"},
		//{"", true, "0.0.0.0"},
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			var addr net.IP
			f := setUpIP(&addr)

			arg := fmt.Sprintf("--address=%s", tt.input)
			err := f.Parse([]string{arg})
			require.NoError(t, err)
			ip, err := f.GetIP("address")
			require.NoError(t, err)
			require.Equal(t, tt.expected, ip.String())
			require.True(t, tt.success)
		})
	}
}

func TestIP_Failures(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"0", ""},
		{"localhost", ""},
		{"0.0.0", ""},
		{"0.0.0.", ""},
		{"0.0.0.0.", ""},
		{"0.0.0.256", ""},
		{"0 . 0 . 0 . 0", ""},
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			var addr net.IP
			f := setUpIP(&addr)

			arg := fmt.Sprintf("--address=%s", tt.input)
			err := f.Parse([]string{arg})
			require.Error(t, err)
		})
	}
}
