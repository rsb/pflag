package pflag_test

import (
	goflag "flag"
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGoflags(t *testing.T) {
	goflag.String("stringFlag", "stringFlag", "stringFlag")
	goflag.Bool("boolFlag", false, "boolFlag")

	f := pflag.NewFlagSet("test", pflag.ContinueOnError)

	f.AddGoFlagSet(goflag.CommandLine)
	err := f.Parse([]string{"--stringFlag=bob", "--boolFlag"})
	require.NoError(t, err)

	getString, err := f.GetString("stringFlag")
	require.NoError(t, err)
	require.Equal(t, "bob", getString)

	getBool, err := f.GetBool("boolFlag")
	require.NoError(t, err)
	require.True(t, getBool)
	require.True(t, f.Parsed())

	// in fact it is useless. because `go test` called flag.Parse()
	require.True(t, goflag.CommandLine.Parsed())
}
