package geom

import (
	"math"
	"testing"
)

func TestScalar_Radians(t *testing.T) {
	tests := []struct {
		name     string
		radians  Radians
		expected Degrees
	}{
		{"Zero", Radians(0.0), Degrees(0.0)},
		{"Pi", Radians(math.Pi), Degrees(180.0)},
		{"Half Pi", Radians(math.PiOver2), Degrees(90.0)},
		{"Two Pi", Radians(2 * math.Pi), Degrees(360.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.radians.Degrees()
			if math.Abs(result.Float64()-tt.expected.Float64()) > 1e-10 {
				t.Errorf("Radians.Degrees() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestScalar_Degrees(t *testing.T) {
	tests := []struct {
		name     string
		degrees  Degrees
		expected Radians
	}{
		{"Zero", Degrees(0.0), Radians(0.0)},
		{"180", Degrees(180.0), Radians(math.Pi)},
		{"90", Degrees(90.0), Radians(math.PiOver2)},
		{"360", Degrees(360.0), Radians(2 * math.Pi)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.degrees.Radians()
			if math.Abs(result.Float64()-tt.expected.Float64()) > 1e-10 {
				t.Errorf("Degrees.Radians() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestScalar_Eq(t *testing.T) {
	tests := []struct {
		name     string
		a, b     interface{}
		expected bool
	}{
		{"Float32 Eq", F32(1.0), F32(1.0), true},
		{"Float32 Close", F32(1.0), F32(1.0 + Epsilon32/2), true},
		{"Float32 Not Eq", F32(1.0), F32(2.0), false},
		{"Float64 Eq", F64(1.0), F64(1.0), true},
		{"Float64 Close", F64(1.0), F64(1.0 + Epsilon64/2), true},
		{"Float64 Not Eq", F64(1.0), F64(2.0), false},
		{"Int32 Eq", I32(42), I32(42), true},
		{"Int32 Not Eq", I32(42), I32(43), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch a := tt.a.(type) {
			case F32:
				result = ScalarEq(a, tt.b.(F32))
			case F64:
				result = ScalarEq(a, tt.b.(F64))
			case I32:
				result = ScalarEq(a, tt.b.(I32))
			}
			if result != tt.expected {
				t.Errorf("Eq(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestScalar_IsFinite(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"Float32 Finite", F32(1.0), true},
		{"Float32 NaN", F32(float32(math.NaN())), false},
		{"Float32 Inf", F32(float32(math.Inf(1))), false},
		{"Float64 Finite", F64(1.0), true},
		{"Float64 NaN", F64(math.NaN()), false},
		{"Float64 Inf", F64(math.Inf(1)), false},
		{"Int32 Always Finite", I32(42), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch v := tt.value.(type) {
			case F32:
				result = IsFinite(v)
			case F64:
				result = IsFinite(v)
			case I32:
				result = IsFinite(v)
			}
			if result != tt.expected {
				t.Errorf("IsFinite(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestScalar_Clamp(t *testing.T) {
	tests := []struct {
		name     string
		value    F32
		min      F32
		max      F32
		expected F32
	}{
		{"Within Range", F32(5.0), F32(1.0), F32(10.0), F32(5.0)},
		{"Below Min", F32(0.5), F32(1.0), F32(10.0), F32(1.0)},
		{"Above Max", F32(15.0), F32(1.0), F32(10.0), F32(10.0)},
		{"At Min", F32(1.0), F32(1.0), F32(10.0), F32(1.0)},
		{"At Max", F32(10.0), F32(1.0), F32(10.0), F32(10.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Clamp(tt.value, tt.min, tt.max)
			if result != tt.expected {
				t.Errorf("Clamp(%v, %v, %v) = %v, want %v", tt.value, tt.min, tt.max, result, tt.expected)
			}
		})
	}
}

func BenchmarkScalar_Eq(b *testing.B) {
	f1 := F64(1.0)
	f2 := F64(1.0000001)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ScalarEq(f1, f2)
	}
}
