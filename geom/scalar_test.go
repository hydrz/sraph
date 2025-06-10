package geom

import (
	"math"
	"testing"
)

// Test NewScalar generic function
func TestNewScalar(t *testing.T) {
	t.Run("I32", func(t *testing.T) {
		val := NewScalar(I32(42))
		if val != I32(42) {
			t.Errorf("Expected I32(42), got %v", val)
		}
	})

	t.Run("F32", func(t *testing.T) {
		val := NewScalar(F32(3.14))
		if val != F32(3.14) {
			t.Errorf("Expected F32(3.14), got %v", val)
		}
	})
}

// I32 Tests
func TestI32_Equal(t *testing.T) {
	tests := []struct {
		name     string
		a, b     I32
		expected bool
	}{
		{"equal positive", I32(42), I32(42), true},
		{"equal negative", I32(-42), I32(-42), true},
		{"equal zero", I32(0), I32(0), true},
		{"not equal", I32(42), I32(43), false},
		{"different signs", I32(42), I32(-42), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScalarEqual(tt.a, tt.b); got != tt.expected {
				t.Errorf("I32.Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestI32_ToFloat64(t *testing.T) {
	tests := []struct {
		input    I32
		expected float64
	}{
		{I32(0), 0.0},
		{I32(42), 42.0},
		{I32(-42), -42.0},
		{I32(2147483647), 2147483647.0},   // max int32
		{I32(-2147483648), -2147483648.0}, // min int32
	}

	for _, tt := range tests {
		if got := tt.input.ToFloat64(); got != tt.expected {
			t.Errorf("I32(%d).ToFloat64() = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestI32_String(t *testing.T) {
	tests := []struct {
		input    I32
		expected string
	}{
		{I32(0), "0"},
		{I32(42), "42"},
		{I32(-42), "-42"},
		{I32(math.MaxInt32), "2147483647"},  // max int32
		{I32(math.MinInt32), "-2147483648"}, // min int32
	}

	for _, tt := range tests {
		if got := tt.input.String(); got != tt.expected {
			t.Errorf("I32(%d).String() = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

// F32 Tests
func TestF32_Equal(t *testing.T) {
	tests := []struct {
		name     string
		a, b     F32
		expected bool
	}{
		{"equal exact", F32(1.0), F32(1.0), true},
		{"not equal", F32(1.0), F32(1.5), false},
		{"equal within epsilon", F32(1.0), F32(1.0001), true},
		{"not equal within epsilon", F32(1.0), F32(1.001), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScalarEqual(tt.a, tt.b); got != tt.expected {
				t.Errorf("F32.Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestF32_ToFloat64(t *testing.T) {
	tests := []struct {
		input    F32
		expected float64
	}{
		{F32(0.0), 0.0},
		{F32(3.14), 3.140000104904175}, // float32 precision
		{F32(-2.718), -2.7179999351501465},
	}

	for _, tt := range tests {
		if got := tt.input.ToFloat64(); got != tt.expected {
			t.Errorf("F32(%f).ToFloat64() = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestF32_String(t *testing.T) {
	tests := []struct {
		input    F32
		expected string
	}{
		{F32(0.0), "0"},
		{F32(3.14), "3.14"},
		{F32(-2.718), "-2.718"},
		{F32(1000000), "1000000"},
		{F32(math.NaN()), "NaN"},
		{F32(math.Inf(1)), "+Inf"},
		{F32(math.Inf(-1)), "-Inf"},
		{F32(math.SmallestNonzeroFloat32), "0.000000000000000000000000000000000000000000001"}, // smallest non-zero float32
	}

	for _, tt := range tests {
		if got := tt.input.String(); got != tt.expected {
			t.Errorf("F32(%f).String() = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

// F64 Tests
func TestF64_Equal(t *testing.T) {
	tests := []struct {
		name     string
		a, b     F64
		expected bool
	}{
		{"equal exact", F64(1.0), F64(1.0), true},
		{"equal within epsilon", F64(1.0), F64(1.0000005), true},
		{"not equal", F64(1.0), F64(1.5), false},
		{"zero values", F64(0.0), F64(0.0), true},
		{"epsilon difference", F64(1.0), F64(1.000001), true},
		{"not epsilon difference", F64(1.0), F64(1.00001), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScalarEqual(tt.a, tt.b); got != tt.expected {
				t.Errorf("F64.Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestF64_ToFloat64(t *testing.T) {
	tests := []struct {
		input    F64
		expected float64
	}{
		{F64(0.0), 0.0},
		{F64(math.Pi), math.Pi},
		{F64(-math.E), -math.E},
		{F64(math.Inf(1)), math.Inf(1)},
		{F64(math.Inf(-1)), math.Inf(-1)},
	}

	for _, tt := range tests {
		got := tt.input.ToFloat64()
		if math.IsInf(float64(tt.input), 0) {
			if !math.IsInf(got, int(math.Copysign(1, float64(tt.input)))) {
				t.Errorf("F64(%f).ToFloat64() = %v, want %v", tt.input, got, tt.expected)
			}
		} else if got != tt.expected {
			t.Errorf("F64(%f).ToFloat64() = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestF64_String(t *testing.T) {
	if got := F64(math.NaN()).String(); got != "NaN" {
		t.Errorf("F64(NaN).String() = %q, want %q", got, "NaN")
	}
}

// FixedI26_6 Tests
func TestFixedI26_6_Equal(t *testing.T) {
	tests := []struct {
		name     string
		a, b     FixedI26_6
		expected bool
	}{
		{"equal zero", FixedI26_6(0), FixedI26_6(0), true},
		{"equal positive", FixedI26_6(1 << 6), FixedI26_6(64), true},
		{"not equal", FixedI26_6(1 << 6), FixedI26_6(1 << 5), false},
		{"equal negative", FixedI26_6(-1 << 6), FixedI26_6(-64), true},
		{"small difference", FixedI26_6(1 << 6), FixedI26_6(1<<6 + 1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScalarEqual(tt.a, tt.b); got != tt.expected {
				t.Errorf("FixedI26_6.Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFixedI26_6_ToFloat64(t *testing.T) {
	tests := []struct {
		input    FixedI26_6
		expected float64
	}{
		{FixedI26_6(0), 0.0},
		{FixedI26_6(1 << 6), 1.0},
		{FixedI26_6(1<<6 + 1), 1.015625},                // 1.0 + 1/64
		{FixedI26_6(33554431 << 6), 33554431.0},         // max value for FixedI26_6
		{FixedI26_6(-33554432 << 6), -33554432.0},       // min value for FixedI26_6
		{FixedI26_6(33554431<<6 + 63), 33554431.984375}, // max value with fractional part
	}

	for _, tt := range tests {
		if got := tt.input.ToFloat64(); math.Abs(got-tt.expected) > 1e-10 {
			t.Errorf("FixedI26_6(%d).ToFloat64() = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestFixedI26_6_String(t *testing.T) {
	tests := []struct {
		input    FixedI26_6
		expected string
	}{
		{FixedI26_6(0), "0"},
		{FixedI26_6(1 << 6), "1"},                         // 1.0
		{FixedI26_6(1<<6 + 1), "1.015625"},                // 1.0 + 1/64
		{FixedI26_6(33554431 << 6), "33554431"},           // max value for FixedI26_6
		{FixedI26_6(-33554432 << 6), "-33554432"},         // min value for FixedI26_6
		{FixedI26_6(33554431<<6 + 63), "33554431.984375"}, // max value with fractional part
	}

	for _, tt := range tests {
		if got := tt.input.String(); got != tt.expected {
			t.Errorf("FixedI26_6(%d).String() = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

// Benchmarks
func Benchmark_ToFloat64(b *testing.B) {
	a := F32(3.14)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.ToFloat64()
	}
}

func Benchmark_Equal(b *testing.B) {
	a, other := F32(3.14), F32(3.14)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ScalarEqual(a, other)
	}
}
