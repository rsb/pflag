package pflag_test

import (
	"fmt"
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

func setUpI64SFlagSet(isp *[]int64) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.Int64SliceVar(isp, "is", []int64{}, "Command separated list!")
	return f
}

func setUpI64SFlagSetWithDefault(isp *[]int64) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.Int64SliceVar(isp, "is", []int64{0, 1}, "Command separated list!")
	return f
}

func TestEmptyI64S(t *testing.T) {
	var is []int64
	f := setUpI64SFlagSet(&is)
	err := f.Parse([]string{})
	require.NoError(t, err)

	getI64S, err := f.GetInt64Slice("is")
	require.NoError(t, err)

	require.Len(t, getI64S, 0)
}

func TestI64S(t *testing.T) {
	var is []int64
	f := setUpI64SFlagSet(&is)

	vals := []string{"1", "2", "4", "3"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range is {
		d, err := strconv.ParseInt(vals[i], 0, 64)
		require.NoError(t, err)
		require.Equal(t, d, v)
	}

	getI64S, err := f.GetInt64Slice("is")
	require.NoError(t, err)

	for i, v := range getI64S {
		d, err := strconv.ParseInt(vals[i], 0, 64)
		require.NoError(t, err)
		require.Equal(t, d, v)
	}
}

func TestI64SDefault(t *testing.T) {
	var is []int64
	f := setUpI64SFlagSetWithDefault(&is)

	vals := []string{"0", "1"}

	err := f.Parse([]string{})
	require.NoError(t, err)

	for i, v := range is {
		d, err := strconv.ParseInt(vals[i], 0, 64)
		require.NoError(t, err)
		require.Equal(t, d, v)
	}

	getI64S, err := f.GetInt64Slice("is")
	require.NoError(t, err)

	for i, v := range getI64S {
		d, err := strconv.ParseInt(vals[i], 0, 64)
		require.NoError(t, err)
		require.Equal(t, d, v)
	}
}

func TestI64SWithDefault(t *testing.T) {
	var is []int64
	f := setUpI64SFlagSetWithDefault(&is)

	vals := []string{"1", "2"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range is {
		d, err := strconv.ParseInt(vals[i], 0, 64)
		require.NoError(t, err)
		require.Equal(t, d, v)
	}

	getI64S, err := f.GetInt64Slice("is")
	require.NoError(t, err)
	for i, v := range getI64S {
		d, err := strconv.ParseInt(vals[i], 0, 64)
		require.NoError(t, err)
		require.Equal(t, d, v)
	}
}

func TestI64SAsSliceValue(t *testing.T) {
	var i64s []int64
	f := setUpI64SFlagSet(&i64s)

	in := []string{"1", "2"}
	argfmt := "--is=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	f.VisitAll(func(f *pflag.Flag) {
		if val, ok := f.Value.(pflag.SliceValue); ok {
			_ = val.Replace([]string{"3"})
		}
	})
	require.Len(t, i64s, 1)
	require.Equal(t, int64(3), i64s[0])
}

func TestI64SCalledTwice(t *testing.T) {
	var is []int64
	f := setUpI64SFlagSet(&is)

	in := []string{"1,2", "3"}
	expected := []int64{1, 2, 3}
	argfmt := "--is=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)

	for i, v := range is {
		require.Equal(t, expected[i], v)
	}
}
