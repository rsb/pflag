package pflag_test

import (
	"fmt"
	"github.com/rsb/pflag"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

func setUpF32SFlagSet(f32sp *[]float32) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.Float32SliceVar(f32sp, "f32s", []float32{}, "Command separated list!")
	return f
}

func setUpF32SFlagSetWithDefault(f32sp *[]float32) *pflag.FlagSet {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.Float32SliceVar(f32sp, "f32s", []float32{0.0, 1.0}, "Command separated list!")
	return f
}

func TestEmptyF32S(t *testing.T) {
	var f32s []float32
	f := setUpF32SFlagSet(&f32s)
	err := f.Parse([]string{})
	require.NoError(t, err)

	getF32S, err := f.GetFloat32Slice("f32s")
	require.NoError(t, err)
	require.Len(t, getF32S, 0)
}

func TestF32S(t *testing.T) {
	var f32s []float32
	f := setUpF32SFlagSet(&f32s)

	vals := []string{"1.0", "2.0", "4.0", "3.0"}
	arg := fmt.Sprintf("--f32s=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range f32s {
		d64, err := strconv.ParseFloat(vals[i], 32)
		require.NoError(t, err)

		d := float32(d64)
		require.Equal(t, d, v)
	}

	getF32S, err := f.GetFloat32Slice("f32s")
	require.NoError(t, err)

	for i, v := range getF32S {
		d64, err := strconv.ParseFloat(vals[i], 32)
		require.NoError(t, err)

		d := float32(d64)
		require.Equal(t, d, v)
	}
}

func TestF32SDefault(t *testing.T) {
	var f32s []float32
	f := setUpF32SFlagSetWithDefault(&f32s)

	vals := []string{"0.0", "1.0"}

	err := f.Parse([]string{})
	require.NoError(t, err)

	for i, v := range f32s {
		d64, err := strconv.ParseFloat(vals[i], 32)
		require.NoError(t, err)

		d := float32(d64)
		require.Equal(t, d, v)
	}

	getF32S, err := f.GetFloat32Slice("f32s")
	require.NoError(t, err)

	for i, v := range getF32S {
		d64, err := strconv.ParseFloat(vals[i], 32)
		require.NoError(t, err)

		d := float32(d64)
		require.Equal(t, d, v)
	}
}

func TestF32SWithDefault(t *testing.T) {
	var f32s []float32
	f := setUpF32SFlagSetWithDefault(&f32s)

	vals := []string{"1.0", "2.0"}
	arg := fmt.Sprintf("--f32s=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	require.NoError(t, err)

	for i, v := range f32s {
		d64, err := strconv.ParseFloat(vals[i], 32)
		require.NoError(t, err)

		d := float32(d64)
		require.Equal(t, d, v)
	}

	getF32S, err := f.GetFloat32Slice("f32s")
	require.NoError(t, err)

	for i, v := range getF32S {
		d64, err := strconv.ParseFloat(vals[i], 32)
		require.NoError(t, err)

		d := float32(d64)
		require.Equal(t, d, v)
	}
}

func TestF32SAsSliceValue(t *testing.T) {
	var f32s []float32
	f := setUpF32SFlagSet(&f32s)

	in := []string{"1.0", "2.0"}
	argfmt := "--f32s=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)

	f.VisitAll(func(f *pflag.Flag) {
		if val, ok := f.Value.(pflag.SliceValue); ok {
			_ = val.Replace([]string{"3.1"})
		}
	})
	require.Len(t, f32s, 1)
	require.Equal(t, float32(3.1), f32s[0])
}

func TestF32SCalledTwice(t *testing.T) {
	var f32s []float32
	f := setUpF32SFlagSet(&f32s)

	in := []string{"1.0,2.0", "3.0"}
	expected := []float32{1.0, 2.0, 3.0}
	argfmt := "--f32s=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	require.NoError(t, err)

	for i, v := range f32s {
		require.Equal(t, expected[i], v)
	}
}