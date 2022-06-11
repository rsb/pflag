package pflag_test

import (
	"fmt"
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

func setUpF64SFlagSet(f64sp *[]float64) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.Float64SliceVar(f64sp, "f64s", []float64{}, "Command separated list!")
	return f
}

func setUpF64SFlagSetWithDefault(f64sp *[]float64) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.Float64SliceVar(f64sp, "f64s", []float64{0.0, 1.0}, "Command separated list!")
	return f
}

func TestEmptyF64S(t *testing.T) {
	var f64s []float64
	f := setUpF64SFlagSet(&f64s)
	err := f.Parse([]string{})
	require.NoError(t, err)

	getF64S, err := f.GetFloat64Slice("f64s")
	require.NoError(t, err)
	require.Len(t, getF64S, 0)
}

func TestF64S(t *testing.T) {
	var f64s []float64
	f := setUpF64SFlagSet(&f64s)

	vals := []string{"1.0", "2.0", "4.0", "3.0"}
	arg := fmt.Sprintf("--f64s=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range f64s {
		d, err := strconv.ParseFloat(vals[i], 64)
		require.NoError(t, err)
		require.Equal(t, d, v)
	}

	getF64S, err := f.GetFloat64Slice("f64s")
	require.NoError(t, err)

	for i, v := range getF64S {
		d, err := strconv.ParseFloat(vals[i], 64)
		require.NoError(t, err)
		require.Equal(t, d, v)
	}
}

func TestF64SDefault(t *testing.T) {
	var f64s []float64
	f := setUpF64SFlagSetWithDefault(&f64s)

	vals := []string{"0.0", "1.0"}

	err := f.Parse([]string{})
	require.NoError(t, err)
	for i, v := range f64s {
		d, err := strconv.ParseFloat(vals[i], 64)
		require.NoError(t, err)
		require.Equal(t, d, v)
	}

	getF64S, err := f.GetFloat64Slice("f64s")
	require.NoError(t, err)
	for i, v := range getF64S {
		d, err := strconv.ParseFloat(vals[i], 64)
		require.NoError(t, err)
		require.Equal(t, d, v)
	}
}

func TestF64SWithDefault(t *testing.T) {
	var f64s []float64
	f := setUpF64SFlagSetWithDefault(&f64s)

	vals := []string{"1.0", "2.0"}
	arg := fmt.Sprintf("--f64s=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range f64s {
		d, err := strconv.ParseFloat(vals[i], 64)
		require.NoError(t, err)
		require.Equal(t, d, v)
	}

	getF64S, err := f.GetFloat64Slice("f64s")
	require.NoError(t, err)
	for i, v := range getF64S {
		d, err := strconv.ParseFloat(vals[i], 64)
		require.NoError(t, err)
		require.Equal(t, d, v)
	}
}

func TestF64SAsSliceValue(t *testing.T) {
	var f64s []float64
	f := setUpF64SFlagSet(&f64s)

	in := []string{"1.0", "2.0"}
	argfmt := "--f64s=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)

	f.VisitAll(func(f *pflag.Flag) {
		if val, ok := f.Value.(pflag.SliceValue); ok {
			_ = val.Replace([]string{"3.1"})
		}
	})
	require.Len(t, f64s, 1)
	require.Equal(t, 3.1, f64s[0])
}

func TestF64SCalledTwice(t *testing.T) {
	var f64s []float64
	f := setUpF64SFlagSet(&f64s)

	in := []string{"1.0,2.0", "3.0"}
	expected := []float64{1.0, 2.0, 3.0}
	argfmt := "--f64s=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)
	for i, v := range f64s {
		require.Equal(t, expected[i], v)
	}
}
