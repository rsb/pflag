package pflag_test

import (
	"fmt"
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

func setUpBSFlagSet(bsp *[]bool) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.BoolSliceVar(bsp, "bs", []bool{}, "Command separated list!")
	return f
}

func setUpBSFlagSetWithDefault(bsp *[]bool) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.BoolSliceVar(bsp, "bs", []bool{false, true}, "Command separated list!")
	return f
}

func TestEmptyBS(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSet(&bs)
	err := f.Parse([]string{})
	require.NoError(t, err)

	getBS, err := f.GetBoolSlice("bs")
	require.NoError(t, err)
	require.Len(t, getBS, 0)
}

func TestBS(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSet(&bs)

	vals := []string{"1", "F", "TRUE", "0"}
	arg := fmt.Sprintf("--bs=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range bs {
		b, err := strconv.ParseBool(vals[i])
		require.NoError(t, err)
		require.Equal(t, b, v)
	}
	getBS, err := f.GetBoolSlice("bs")
	require.NoError(t, err)

	for i, v := range getBS {
		b, err := strconv.ParseBool(vals[i])
		require.NoError(t, err)
		require.Equal(t, b, v)
	}
}

func TestBSDefault(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSetWithDefault(&bs)

	vals := []string{"false", "T"}

	err := f.Parse([]string{})
	require.NoError(t, err)

	for i, v := range bs {
		b, err := strconv.ParseBool(vals[i])
		require.NoError(t, err)
		require.Equal(t, b, v)
	}

	getBS, err := f.GetBoolSlice("bs")
	require.NoError(t, err)

	for i, v := range getBS {
		b, err := strconv.ParseBool(vals[i])
		require.NoError(t, err)
		require.Equal(t, b, v)
	}
}

func TestBSWithDefault(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSetWithDefault(&bs)

	vals := []string{"FALSE", "1"}
	arg := fmt.Sprintf("--bs=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range bs {
		b, err := strconv.ParseBool(vals[i])
		require.NoError(t, err)
		require.Equal(t, b, v)
	}

	getBS, err := f.GetBoolSlice("bs")
	require.NoError(t, err)

	for i, v := range getBS {
		b, err := strconv.ParseBool(vals[i])
		require.NoError(t, err)
		require.Equal(t, b, v)
	}
}

func TestBSCalledTwice(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSet(&bs)

	in := []string{"T,F", "T"}
	expected := []bool{true, false, true}
	argfmt := "--bs=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)

	for i, v := range bs {
		require.Equal(t, expected[i], v)
	}
}

func TestBSAsSliceValue(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSet(&bs)

	in := []string{"true", "false"}
	argfmt := "--bs=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)

	f.VisitAll(func(f *pflag.Flag) {
		if val, ok := f.Value.(pflag.SliceValue); ok {
			_ = val.Replace([]string{"false"})
		}
	})

	require.Len(t, bs, 1)
	require.False(t, bs[0])
}

func TestBSBadQuoting(t *testing.T) {

	tests := []struct {
		Want    []bool
		FlagArg []string
	}{
		{
			Want:    []bool{true, false, true},
			FlagArg: []string{"1", "0", "true"},
		},
		{
			Want:    []bool{true, false},
			FlagArg: []string{"True", "F"},
		},
		{
			Want:    []bool{true, false},
			FlagArg: []string{"T", "0"},
		},
		{
			Want:    []bool{true, false},
			FlagArg: []string{"1", "0"},
		},
		{
			Want:    []bool{true, false, false},
			FlagArg: []string{"true,false", "false"},
		},
		{
			Want:    []bool{true, false, false, true, false, true, false},
			FlagArg: []string{`"true,false,false,1,0,     T"`, " false "},
		},
		{
			Want:    []bool{false, false, true, false, true, false, true},
			FlagArg: []string{`"0, False,  T,false  , true,F"`, "true"},
		},
	}

	for _, test := range tests {

		var bs []bool
		f := setUpBSFlagSet(&bs)

		err := f.Parse([]string{fmt.Sprintf("--bs=%s", strings.Join(test.FlagArg, ","))})
		require.NoError(t, err)

		for j, b := range bs {
			require.Equal(t, test.Want[j], b)
		}
	}
}
