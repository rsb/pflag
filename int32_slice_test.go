package pflag_test

import (
	"fmt"
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

func setUpI32SFlagSet(isp *[]int32) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.Int32SliceVar(isp, "is", []int32{}, "Command separated list!")
	return f
}

func setUpI32SFlagSetWithDefault(isp *[]int32) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.Int32SliceVar(isp, "is", []int32{0, 1}, "Command separated list!")
	return f
}

func TestEmptyI32S(t *testing.T) {
	var is []int32
	f := setUpI32SFlagSet(&is)
	err := f.Parse([]string{})
	require.NoError(t, err)

	getI32S, err := f.GetInt32Slice("is")
	require.NoError(t, err)
	require.Len(t, getI32S, 0)
}

func TestI32S(t *testing.T) {
	var is []int32
	f := setUpI32SFlagSet(&is)

	vals := []string{"1", "2", "4", "3"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range is {
		d64, err := strconv.ParseInt(vals[i], 0, 32)
		require.NoError(t, err)
		d := int32(d64)
		require.Equal(t, d, v)
	}

	getI32S, err := f.GetInt32Slice("is")
	require.NoError(t, err)

	for i, v := range getI32S {
		d64, err := strconv.ParseInt(vals[i], 0, 32)
		require.NoError(t, err)

		d := int32(d64)
		require.Equal(t, d, v)
	}
}

func TestI32SDefault(t *testing.T) {
	var is []int32
	f := setUpI32SFlagSetWithDefault(&is)

	vals := []string{"0", "1"}

	err := f.Parse([]string{})
	require.NoError(t, err)

	for i, v := range is {
		d64, err := strconv.ParseInt(vals[i], 0, 32)
		require.NoError(t, err)

		d := int32(d64)
		require.Equal(t, d, v)
	}

	getI32S, err := f.GetInt32Slice("is")
	require.NoError(t, err)

	for i, v := range getI32S {
		d64, err := strconv.ParseInt(vals[i], 0, 32)
		require.NoError(t, err)

		d := int32(d64)
		require.Equal(t, d, v)
	}
}

func TestI32SWithDefault(t *testing.T) {
	var is []int32
	f := setUpI32SFlagSetWithDefault(&is)

	vals := []string{"1", "2"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range is {
		d64, err := strconv.ParseInt(vals[i], 0, 32)
		require.NoError(t, err)
		d := int32(d64)
		require.Equal(t, d, v)
	}

	getI32S, err := f.GetInt32Slice("is")
	require.NoError(t, err)

	for i, v := range getI32S {
		d64, err := strconv.ParseInt(vals[i], 0, 32)
		require.NoError(t, err)

		d := int32(d64)
		require.Equal(t, d, v)
	}
}

func TestI32SAsSliceValue(t *testing.T) {
	var i32s []int32
	f := setUpI32SFlagSet(&i32s)

	in := []string{"1", "2"}
	argfmt := "--is=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)

	f.VisitAll(func(f *pflag.Flag) {
		if val, ok := f.Value.(pflag.SliceValue); ok {
			_ = val.Replace([]string{"3"})
		}
	})
	require.Len(t, i32s, 1)
	require.Equal(t, int32(3), i32s[0])
}

func TestI32SCalledTwice(t *testing.T) {
	var is []int32
	f := setUpI32SFlagSet(&is)

	in := []string{"1,2", "3"}
	expected := []int32{1, 2, 3}
	argfmt := "--is=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)

	for i, v := range is {
		require.Equal(t, expected[i], v)
	}
}
