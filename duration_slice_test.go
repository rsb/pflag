package pflag_test

import (
	"fmt"
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func setUpDSFlagSet(dsp *[]time.Duration) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.DurationSliceVar(dsp, "ds", []time.Duration{}, "Command separated list!")
	return f
}

func setUpDSFlagSetWithDefault(dsp *[]time.Duration) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.DurationSliceVar(dsp, "ds", []time.Duration{0, 1}, "Command separated list!")
	return f
}

func TestEmptyDS(t *testing.T) {
	var ds []time.Duration
	f := setUpDSFlagSet(&ds)
	err := f.Parse([]string{})
	require.NoError(t, err)

	getDS, err := f.GetDurationSlice("ds")
	require.NoError(t, err)
	require.Len(t, getDS, 0)
}

func TestDS(t *testing.T) {
	var ds []time.Duration
	f := setUpDSFlagSet(&ds)

	vals := []string{"1ns", "2ms", "3m", "4h"}
	arg := fmt.Sprintf("--ds=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range ds {
		d, err := time.ParseDuration(vals[i])
		require.NoError(t, err)
		require.Equal(t, d, v)
	}

	getDS, err := f.GetDurationSlice("ds")
	require.NoError(t, err)

	for i, v := range getDS {
		d, err := time.ParseDuration(vals[i])
		require.NoError(t, err)
		require.Equal(t, d, v)
	}
}

func TestDSDefault(t *testing.T) {
	var ds []time.Duration
	f := setUpDSFlagSetWithDefault(&ds)

	vals := []string{"0s", "1ns"}

	err := f.Parse([]string{})
	require.NoError(t, err)

	for i, v := range ds {
		d, err := time.ParseDuration(vals[i])
		require.NoError(t, err)
		require.Equal(t, d, v)
	}

	getDS, err := f.GetDurationSlice("ds")
	require.NoError(t, err)

	for i, v := range getDS {
		d, err := time.ParseDuration(vals[i])
		require.NoError(t, err)
		require.Equal(t, d, v)
	}
}

func TestDSWithDefault(t *testing.T) {
	var ds []time.Duration
	f := setUpDSFlagSetWithDefault(&ds)

	vals := []string{"1ns", "2ns"}
	arg := fmt.Sprintf("--ds=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range ds {
		d, err := time.ParseDuration(vals[i])
		require.NoError(t, err)
		require.Equal(t, d, v)
	}

	getDS, err := f.GetDurationSlice("ds")
	require.NoError(t, err)

	for i, v := range getDS {
		d, err := time.ParseDuration(vals[i])
		require.NoError(t, err)
		require.Equal(t, d, v)
	}
}

func TestDSAsSliceValue(t *testing.T) {
	var ds []time.Duration
	f := setUpDSFlagSet(&ds)

	in := []string{"1ns", "2ns"}
	argfmt := "--ds=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)

	f.VisitAll(func(f *pflag.Flag) {
		if val, ok := f.Value.(pflag.SliceValue); ok {
			_ = val.Replace([]string{"3ns"})
		}
	})
	require.Len(t, ds, 1)
	require.Equal(t, time.Duration(3), ds[0])
}

func TestDSCalledTwice(t *testing.T) {
	var ds []time.Duration
	f := setUpDSFlagSet(&ds)

	in := []string{"1ns,2ns", "3ns"}
	expected := []time.Duration{1, 2, 3}
	argfmt := "--ds=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)
	for i, v := range ds {
		require.Equal(t, expected[i], v)
	}
}
