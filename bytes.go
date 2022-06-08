package pflag

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// BytesHex adapts []byte for use as a flag. Value of flag is HEX encoded
type bytesHexValue []byte

// String implements pflag.Value.String.
func (bytesHex bytesHexValue) String() string {
	return fmt.Sprintf("%X", []byte(bytesHex))
}

// Set implements pflag.Value.Set.
func (bytesHex *bytesHexValue) Set(value string) error {
	bin, err := hex.DecodeString(strings.TrimSpace(value))

	if err != nil {
		return err
	}

	*bytesHex = bin

	return nil
}

// Type implements pflag.Value.Type.
func (*bytesHexValue) Type() string {
	return "bytesHex"
}

func newBytesHexValue(val []byte, p *[]byte) *bytesHexValue {
	*p = val
	return (*bytesHexValue)(p)
}

func bytesHexConv(sval string) (interface{}, error) {

	bin, err := hex.DecodeString(sval)

	if err == nil {
		return bin, nil
	}

	return nil, fmt.Errorf("invalid string being converted to Bytes: %s %s", sval, err)
}
