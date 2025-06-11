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
		{"Int32", Int32(42), Int32(42)},
		{"Float32", Float32(3.14), Float32(3.14)},
		{"Float64", Float64(2.718), Float64(2.718)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case Int32:
				result := New(v)
				if result != tt.expected {
					t.Errorf("New() = %v, want %v", result, tt.expected)
				}
			case Float32:
				result := New(v)
				if result != tt.expected {
					t.Errorf("New() = %v, want %v", result, tt.expected)
				}
			case Float64:
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
		radians  Radians
		expected Degrees
	}{
		{"Zero", Radians(0), Degrees(0)},
		{"Pi", Radians(math.Pi), Degrees(180)},
		{"Half Pi", Radians(math.Pi / 2), Degrees(90)},
		{"Two Pi", Radians(2 * math.Pi), Degrees(360)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.radians.Degrees()
			if math.Abs(float64(result-tt.expected)) > 1e-10 {
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
		{"Zero", Degrees(0), Radians(0)},
		{"180", Degrees(180), Radians(math.Pi)},
		{"90", Degrees(90), Radians(math.Pi / 2)},
		{"360", Degrees(360), Radians(2 * math.Pi)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.degrees.Radians()
			if math.Abs(float64(result-tt.expected)) > 1e-10 {
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
			result := Max[Int32]()
			return result == Int32(math.MaxInt32)
		}},
		{"Float32", func() bool {
			result := Max[Float32]()
			return result == Float32(math.MaxFloat32)
		}},
		{"Float64", func() bool {
			result := Max[Float64]()
			return result == Float64(math.MaxFloat64)
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
		{"Float32 Equal", Float32(1.0), Float32(1.0), true},
		{"Float32 Close", Float32(1.0), Float32(1.0 + Epsilon32/2), true},
		{"Float32 Not Equal", Float32(1.0), Float32(2.0), false},
		{"Float64 Equal", Float64(1.0), Float64(1.0), true},
		{"Float64 Close", Float64(1.0), Float64(1.0 + Epsilon64/2), true},
		{"Float64 Not Equal", Float64(1.0), Float64(2.0), false},
		{"Int32 Equal", Int32(42), Int32(42), true},
		{"Int32 Not Equal", Int32(42), Int32(43), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch a := tt.a.(type) {
			case Float32:
				result = Equal(a, tt.b.(Float32))
			case Float64:
				result = Equal(a, tt.b.(Float64))
			case Int32:
				result = Equal(a, tt.b.(Int32))
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
		{"Float32 Finite", Float32(1.0), true},
		{"Float32 NaN", Float32(float32(math.NaN())), false},
		{"Float32 Inf", Float32(float32(math.Inf(1))), false},
		{"Float64 Finite", Float64(1.0), true},
		{"Float64 NaN", Float64(math.NaN()), false},
		{"Float64 Inf", Float64(math.Inf(1)), false},
		{"Int32 Always Finite", Int32(42), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch v := tt.value.(type) {
			case Float32:
				result = IsFinite(v)
			case Float64:
				result = IsFinite(v)
			case Int32:
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
		value    Float32
		min      Float32
		max      Float32
		expected Float32
	}{
		{"Within Range", Float32(5.0), Float32(1.0), Float32(10.0), Float32(5.0)},
		{"Below Min", Float32(0.5), Float32(1.0), Float32(10.0), Float32(1.0)},
		{"Above Max", Float32(15.0), Float32(1.0), Float32(10.0), Float32(10.0)},
		{"At Min", Float32(1.0), Float32(1.0), Float32(10.0), Float32(1.0)},
		{"At Max", Float32(10.0), Float32(1.0), Float32(10.0), Float32(10.0)},
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
		i := Int32(42)
		if i.Float64() != 42.0 {
			t.Errorf("Int32.Float64() = %v, want 42.0", i.Float64())
		}
		if i.String() != "42" {
			t.Errorf("Int32.String() = %v, want '42'", i.String())
		}
	})

	t.Run("Float32", func(t *testing.T) {
		f := Float32(3.14)
		if math.Abs(f.Float64()-3.14) > 1e-6 {
			t.Errorf("Float32.Float64() = %v, want 3.14", f.Float64())
		}
		if f.Value() != 3.14 {
			t.Errorf("Float32.Value() = %v, want 3.14", f.Value())
		}
	})

	t.Run("Float64", func(t *testing.T) {
		f := Float64(2.718)
		if f.Float64() != 2.718 {
			t.Errorf("Float64.Float64() = %v, want 2.718", f.Float64())
		}
	})

	t.Run("Fixed26_6", func(t *testing.T) {
		// Test Fixed26_6: 1.25 = 1<<6 + 1<<4 = 64 + 16 = 80
		f := Fixed26_6(80)
		expected := 1.25
		if math.Abs(f.Float64()-expected) > 1e-10 {
			t.Errorf("Fixed26_6.Float64() = %v, want %v", f.Float64(), expected)
		}
	})
}

func BenchmarkScalar_Equal(b *testing.B) {
	f1 := Float64(1.0)
	f2 := Float64(1.0000001)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(f1, f2)
	}
}

func BenchmarkScalar_IsFinite(b *testing.B) {
	f := Float64(1.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsFinite(f)
	}
}
