package pflag_test

import (
	"fmt"
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

func setUpISFlagSet(isp *[]int) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.IntSliceVar(isp, "is", []int{}, "Command separated list!")
	return f
}

func setUpISFlagSetWithDefault(isp *[]int) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.IntSliceVar(isp, "is", []int{0, 1}, "Command separated list!")
	return f
}

func TestEmptyIS(t *testing.T) {
	var is []int
	f := setUpISFlagSet(&is)
	err := f.Parse([]string{})
	require.NoError(t, err)

	getIS, err := f.GetIntSlice("is")
	require.NoError(t, err)
	require.Len(t, getIS, 0)
}

func TestIS(t *testing.T) {
	var is []int
	f := setUpISFlagSet(&is)

	vals := []string{"1", "2", "4", "3"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range is {
		d, err := strconv.Atoi(vals[i])
		require.NoError(t, err)
		require.Equal(t, d, v)
	}

	getIS, err := f.GetIntSlice("is")
	require.NoError(t, err)

	for i, v := range getIS {
		d, err := strconv.Atoi(vals[i])
		require.NoError(t, err)
		require.Equal(t, d, v)
	}
}

func TestISDefault(t *testing.T) {
	var is []int
	f := setUpISFlagSetWithDefault(&is)

	vals := []string{"0", "1"}

	err := f.Parse([]string{})
	require.NoError(t, err)

	for i, v := range is {
		d, err := strconv.Atoi(vals[i])
		require.NoError(t, err)
		require.Equal(t, d, v)
	}

	getIS, err := f.GetIntSlice("is")
	require.NoError(t, err)

	for i, v := range getIS {
		d, err := strconv.Atoi(vals[i])
		require.NoError(t, err)
		require.Equal(t, d, v)
	}
}

func TestISWithDefault(t *testing.T) {
	var is []int
	f := setUpISFlagSetWithDefault(&is)

	vals := []string{"1", "2"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range is {
		d, err := strconv.Atoi(vals[i])
		require.NoError(t, err)
		require.Equal(t, d, v)
	}

	getIS, err := f.GetIntSlice("is")
	require.NoError(t, err)

	for i, v := range getIS {
		d, err := strconv.Atoi(vals[i])
		require.NoError(t, err)
		require.Equal(t, d, v)
	}
}

func TestISCalledTwice(t *testing.T) {
	var is []int
	f := setUpISFlagSet(&is)

	in := []string{"1,2", "3"}
	expected := []int{1, 2, 3}
	argfmt := "--is=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)

	for i, v := range is {
		require.Equal(t, expected[i], v)
	}
}
