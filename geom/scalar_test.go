package geom

import (
	"math"
	"testing"
)

func TestScalar_New(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected interface{}
	}{
		{"Int32", I32(42), I32(42)},
		{"Float32", F32(3.14), F32(3.14)},
		{"Float64", F64(2.718), F64(2.718)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case I32:
				result := New(v)
				if result != tt.expected {
					t.Errorf("New() = %v, want %v", result, tt.expected)
				}
			case F32:
				result := New(v)
				if result != tt.expected {
					t.Errorf("New() = %v, want %v", result, tt.expected)
				}
			case F64:
				result := New(v)
				if result != tt.expected {
					t.Errorf("New() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestScalar_Radians(t *testing.T) {
	tests := []struct {
		name     string
		radians  Radians[F32]
		expected Degrees[F32]
	}{
		{"Zero", NewRadians(F32(0)), NewDegrees(F32(0))},
		{"Pi", NewRadians(F32(math.Pi)), NewDegrees(F32(180))},
		{"Half Pi", NewRadians(F32(math.Pi / 2)), NewDegrees(F32(90))},
		{"Two Pi", NewRadians(F32(2 * math.Pi)), NewDegrees(F32(360))},
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
		degrees  Degrees[F32]
		expected Radians[F32]
	}{
		{"Zero", NewDegrees(F32(0)), NewRadians(F32(0))},
		{"180", NewDegrees(F32(180)), NewRadians(F32(math.Pi))},
		{"90", NewDegrees(F32(90)), NewRadians(F32(math.Pi / 2))},
		{"360", NewDegrees(F32(360)), NewRadians(F32(2 * math.Pi))},
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

func TestScalar_Max(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() bool
	}{
		{"Int32", func() bool {
			result := Max[I32]()
			return result == I32(math.MaxInt32)
		}},
		{"Float32", func() bool {
			result := Max[F32]()
			return result == F32(math.MaxFloat32)
		}},
		{"Float64", func() bool {
			result := Max[F64]()
			return result == F64(math.MaxFloat64)
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.testFunc() {
				t.Errorf("Max() test failed for %s", tt.name)
			}
		})
	}
}

func TestScalar_Equal(t *testing.T) {
	tests := []struct {
		name     string
		a, b     interface{}
		expected bool
	}{
		{"Float32 Equal", F32(1.0), F32(1.0), true},
		{"Float32 Close", F32(1.0), F32(1.0 + Epsilon32/2), true},
		{"Float32 Not Equal", F32(1.0), F32(2.0), false},
		{"Float64 Equal", F64(1.0), F64(1.0), true},
		{"Float64 Close", F64(1.0), F64(1.0 + Epsilon64/2), true},
		{"Float64 Not Equal", F64(1.0), F64(2.0), false},
		{"Int32 Equal", I32(42), I32(42), true},
		{"Int32 Not Equal", I32(42), I32(43), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch a := tt.a.(type) {
			case F32:
				result = Equal(a, tt.b.(F32))
			case F64:
				result = Equal(a, tt.b.(F64))
			case I32:
				result = Equal(a, tt.b.(I32))
			}
			if result != tt.expected {
				t.Errorf("Equal(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
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

func TestScalar_ScalarTypes(t *testing.T) {
	t.Run("Int32", func(t *testing.T) {
		i := I32(42)
		if i.Float64() != 42.0 {
			t.Errorf("Int32.Float64() = %v, want 42.0", i.Float64())
		}
		if i.String() != "42" {
			t.Errorf("Int32.String() = %v, want '42'", i.String())
		}
	})

	t.Run("Float32", func(t *testing.T) {
		f := F32(3.14)
		if math.Abs(f.Float64()-3.14) > 1e-6 {
			t.Errorf("Float32.Float64() = %v, want 3.14", f.Float64())
		}
		if f.Value() != 3.14 {
			t.Errorf("Float32.Value() = %v, want 3.14", f.Value())
		}
	})

	t.Run("Float64", func(t *testing.T) {
		f := F64(2.718)
		if f.Float64() != 2.718 {
			t.Errorf("Float64.Float64() = %v, want 2.718", f.Float64())
		}
	})

	t.Run("Fixed26_6", func(t *testing.T) {
		// Test Fixed26_6: 1.25 = 1<<6 + 1<<4 = 64 + 16 = 80
		f := I26_6(80)
		expected := 1.25
		if math.Abs(f.Float64()-expected) > 1e-10 {
			t.Errorf("Fixed26_6.Float64() = %v, want %v", f.Float64(), expected)
		}
	})
}

func BenchmarkScalar_Equal(b *testing.B) {
	f1 := F64(1.0)
	f2 := F64(1.0000001)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(f1, f2)
	}
}

func BenchmarkScalar_IsFinite(b *testing.B) {
	f := F64(1.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsFinite(f)
	}
}
