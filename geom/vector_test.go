package geom

import (
	"math"
	"testing"
)

func TestNewVector2(t *testing.T) {
	v := NewVector2(Float32(1.0), Float32(2.0))
	if v.X() != 1.0 || v.Y() != 2.0 {
		t.Errorf("NewVector2(1.0, 2.0) = (%v, %v), want (1.0, 2.0)", v.X(), v.Y())
	}
}

func TestNewVector3(t *testing.T) {
	v := NewVector3(Float32(1.0), Float32(2.0), Float32(3.0))
	if v.X() != 1.0 || v.Y() != 2.0 || v.Z() != 3.0 {
		t.Errorf("NewVector3(1.0, 2.0, 3.0) = (%v, %v, %v), want (1.0, 2.0, 3.0)", v.X(), v.Y(), v.Z())
	}
}

func TestNewVector4(t *testing.T) {
	v := NewVector4(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	if v.X() != 1.0 || v.Y() != 2.0 || v.Z() != 3.0 || v.W() != 4.0 {
		t.Errorf("NewVector4() components mismatch")
	}
}

func TestVector3Arithmetic(t *testing.T) {
	v1 := NewVector3(Float32(1.0), Float32(2.0), Float32(3.0))
	v2 := NewVector3(Float32(4.0), Float32(5.0), Float32(6.0))

	t.Run("Add", func(t *testing.T) {
		result := v1.Add(v2)
		if result.X() != 5.0 || result.Y() != 7.0 || result.Z() != 9.0 {
			t.Errorf("Add() = (%v, %v, %v), want (5.0, 7.0, 9.0)", result.X(), result.Y(), result.Z())
		}
	})

	t.Run("Sub", func(t *testing.T) {
		result := v2.Sub(v1)
		if result.X() != 3.0 || result.Y() != 3.0 || result.Z() != 3.0 {
			t.Errorf("Sub() = (%v, %v, %v), want (3.0, 3.0, 3.0)", result.X(), result.Y(), result.Z())
		}
	})

	t.Run("Mul", func(t *testing.T) {
		result := v1.Mul(v2)
		if result.X() != 4.0 || result.Y() != 10.0 || result.Z() != 18.0 {
			t.Errorf("Mul() = (%v, %v, %v), want (4.0, 10.0, 18.0)", result.X(), result.Y(), result.Z())
		}
	})

	t.Run("Div", func(t *testing.T) {
		result := v2.Div(v1)
		if !Equal(result.X(), Float32(4.0)) || !Equal(result.Y(), Float32(2.5)) || !Equal(result.Z(), Float32(2.0)) {
			t.Errorf("Div() = (%v, %v, %v), want (4.0, 2.5, 2.0)", result.X(), result.Y(), result.Z())
		}
	})

	t.Run("Scale", func(t *testing.T) {
		result := v1.Scale(Float32(2.0))
		if result.X() != 2.0 || result.Y() != 4.0 || result.Z() != 6.0 {
			t.Errorf("Scale(2.0) = (%v, %v, %v), want (2.0, 4.0, 6.0)", result.X(), result.Y(), result.Z())
		}
	})
}

func TestVector3VectorOperations(t *testing.T) {
	v := NewVector3(Float32(3.0), Float32(4.0), Float32(0.0))

	t.Run("Length", func(t *testing.T) {
		length := v.Length()
		expected := Float32(5.0) // sqrt(3^2 + 4^2 + 0^2) = 5
		if !Equal(length, expected) {
			t.Errorf("Length() = %v, want %v", length, expected)
		}
	})

	t.Run("Normalize", func(t *testing.T) {
		normalized := v.Normalize()

		// Should have length 1
		length := normalized.Length()
		if !Equal(length, Float32(1.0)) {
			t.Errorf("Normalized vector length = %v, want 1.0", length)
		}

		// Check components
		expectedX := Float32(0.6) // 3/5
		expectedY := Float32(0.8) // 4/5
		if !Equal(normalized.X(), expectedX) || !Equal(normalized.Y(), expectedY) || !Equal(normalized.Z(), Float32(0.0)) {
			t.Errorf("Normalize() = (%v, %v, %v), want (%v, %v, 0.0)", normalized.X(), normalized.Y(), normalized.Z(), expectedX, expectedY)
		}
	})

	t.Run("Zero vector normalize", func(t *testing.T) {
		zero := NewVector3(Float32(0.0), Float32(0.0), Float32(0.0))
		normalized := zero.Normalize()

		// Should return zero vector
		if !Equal(normalized.X(), Float32(0.0)) || !Equal(normalized.Y(), Float32(0.0)) || !Equal(normalized.Z(), Float32(0.0)) {
			t.Errorf("Normalized zero vector should be zero")
		}
	})
}

func TestVector3DotProduct(t *testing.T) {
	v1 := NewVector3(Float32(1.0), Float32(2.0), Float32(3.0))
	v2 := NewVector3(Float32(4.0), Float32(5.0), Float32(6.0))

	dot := v1.Dot(v2)
	expected := Float32(32.0) // 1*4 + 2*5 + 3*6 = 32
	if !Equal(dot, expected) {
		t.Errorf("Dot() = %v, want %v", dot, expected)
	}
}

func TestVector3CrossProduct(t *testing.T) {
	v1 := NewVector3(Float32(1.0), Float32(0.0), Float32(0.0))
	v2 := NewVector3(Float32(0.0), Float32(1.0), Float32(0.0))

	cross := v1.Cross(v2)
	// (1,0,0) × (0,1,0) = (0,0,1)
	if !Equal(cross.X(), Float32(0.0)) || !Equal(cross.Y(), Float32(0.0)) || !Equal(cross.Z(), Float32(1.0)) {
		t.Errorf("Cross() = (%v, %v, %v), want (0.0, 0.0, 1.0)", cross.X(), cross.Y(), cross.Z())
	}
}

func TestVector3MathFunctions(t *testing.T) {
	v := NewVector3(Float32(1.7), Float32(2.3), Float32(3.9))

	t.Run("Floor", func(t *testing.T) {
		result := v.Floor()
		if result.X() != 1.0 || result.Y() != 2.0 || result.Z() != 3.0 {
			t.Errorf("Floor() = (%v, %v, %v), want (1.0, 2.0, 3.0)", result.X(), result.Y(), result.Z())
		}
	})

	t.Run("Ceil", func(t *testing.T) {
		result := v.Ceil()
		if result.X() != 2.0 || result.Y() != 3.0 || result.Z() != 4.0 {
			t.Errorf("Ceil() = (%v, %v, %v), want (2.0, 3.0, 4.0)", result.X(), result.Y(), result.Z())
		}
	})

	t.Run("Round", func(t *testing.T) {
		result := v.Round()
		if result.X() != 2.0 || result.Y() != 2.0 || result.Z() != 4.0 {
			t.Errorf("Round() = (%v, %v, %v), want (2.0, 2.0, 4.0)", result.X(), result.Y(), result.Z())
		}
	})

	t.Run("Abs", func(t *testing.T) {
		negV := NewVector3(Float32(-1.0), Float32(-2.0), Float32(-3.0))
		result := negV.Abs()
		if result.X() != 1.0 || result.Y() != 2.0 || result.Z() != 3.0 {
			t.Errorf("Abs() = (%v, %v, %v), want (1.0, 2.0, 3.0)", result.X(), result.Y(), result.Z())
		}
	})
}

func TestVector3MinMax(t *testing.T) {
	v1 := NewVector3(Float32(1.0), Float32(5.0), Float32(3.0))
	v2 := NewVector3(Float32(4.0), Float32(2.0), Float32(6.0))
	v3 := NewVector3(Float32(2.0), Float32(3.0), Float32(1.0))

	t.Run("Min with multiple vectors", func(t *testing.T) {
		result := v1.Min(v2, v3)
		if result.X() != 1.0 || result.Y() != 2.0 || result.Z() != 1.0 {
			t.Errorf("Min() = (%v, %v, %v), want (1.0, 2.0, 1.0)", result.X(), result.Y(), result.Z())
		}
	})

	t.Run("Max with multiple vectors", func(t *testing.T) {
		result := v1.Max(v2, v3)
		if result.X() != 4.0 || result.Y() != 5.0 || result.Z() != 6.0 {
			t.Errorf("Max() = (%v, %v, %v), want (4.0, 5.0, 6.0)", result.X(), result.Y(), result.Z())
		}
	})
}

func TestVector3Lerp(t *testing.T) {
	v1 := NewVector3(Float32(0.0), Float32(0.0), Float32(0.0))
	v2 := NewVector3(Float32(10.0), Float32(20.0), Float32(30.0))

	tests := []struct {
		t        Float32
		expected Vector3[Float32]
	}{
		{Float32(0.0), v1},
		{Float32(1.0), v2},
		{Float32(0.5), NewVector3(Float32(5.0), Float32(10.0), Float32(15.0))},
		{Float32(0.25), NewVector3(Float32(2.5), Float32(5.0), Float32(7.5))},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := v1.Lerp(v2, tt.t)
			if !Equal(result.X(), tt.expected.X()) || !Equal(result.Y(), tt.expected.Y()) || !Equal(result.Z(), tt.expected.Z()) {
				t.Errorf("Lerp(%v) = (%v, %v, %v), want (%v, %v, %v)", tt.t, result.X(), result.Y(), result.Z(), tt.expected.X(), tt.expected.Y(), tt.expected.Z())
			}
		})
	}
}

func TestVector3Combine(t *testing.T) {
	v1 := NewVector3(Float32(1.0), Float32(2.0), Float32(3.0))
	v2 := NewVector3(Float32(4.0), Float32(5.0), Float32(6.0))

	result := v1.Combine(v2, Float32(2.0))
	// v1 + v2 * 2.0 = (1,2,3) + (8,10,12) = (9,12,15)
	if result.X() != 9.0 || result.Y() != 12.0 || result.Z() != 15.0 {
		t.Errorf("Combine() = (%v, %v, %v), want (9.0, 12.0, 15.0)", result.X(), result.Y(), result.Z())
	}
}

func TestVector3Equal(t *testing.T) {
	v1 := NewVector3(Float32(1.0), Float32(2.0), Float32(3.0))
	v2 := NewVector3(Float32(1.0), Float32(2.0), Float32(3.0))
	v3 := NewVector3(Float32(1.0), Float32(2.0), Float32(4.0))

	if !v1.Equal(v2) {
		t.Error("Equal vectors should be equal")
	}
	if v1.Equal(v3) {
		t.Error("Different vectors should not be equal")
	}
}

func TestVector3String(t *testing.T) {
	v := NewVector3(Float32(1.0), Float32(2.0), Float32(3.0))
	str := v.String()
	expected := "Vector3(1, 2, 3)"
	if str != expected {
		t.Errorf("String() = %v, want %v", str, expected)
	}
}

// Vector4 tests
func TestVector4XY(t *testing.T) {
	v4 := NewVector4(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	v2 := v4.XY()

	if v2.X() != 1.0 || v2.Y() != 2.0 {
		t.Errorf("XY() = (%v, %v), want (1.0, 2.0)", v2.X(), v2.Y())
	}
}

func TestVector4Length(t *testing.T) {
	v := NewVector4(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	length := v.Length()
	expected := Float32(math.Sqrt(30.0)) // sqrt(1^2 + 2^2 + 3^2 + 4^2) = sqrt(30)
	if !Equal(length, expected) {
		t.Errorf("Length() = %v, want %v", length, expected)
	}
}

func TestVector4CrossProduct(t *testing.T) {
	v1 := NewVector4(Float32(1.0), Float32(0.0), Float32(0.0), Float32(0.0))
	v2 := NewVector4(Float32(0.0), Float32(1.0), Float32(0.0), Float32(0.0))

	// Cross product returns zero vector for 4D
	cross := v1.Cross(v2)
	if cross.X() != 0.0 || cross.Y() != 0.0 || cross.Z() != 0.0 || cross.W() != 0.0 {
		t.Errorf("Cross() should return zero vector for 4D")
	}
}

func TestVector4IsFinite(t *testing.T) {
	finite := NewVector4(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	infinite := NewVector4(Float32(math.Inf(1)), Float32(2.0), Float32(3.0), Float32(4.0))

	if !finite.IsFinite() {
		t.Error("Finite vector should return true for IsFinite()")
	}
	if infinite.IsFinite() {
		t.Error("Infinite vector should return false for IsFinite()")
	}
}

func TestVector4String(t *testing.T) {
	v := NewVector4(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	str := v.String()
	expected := "Vector4(1, 2, 3, 4)"
	if str != expected {
		t.Errorf("String() = %v, want %v", str, expected)
	}
}

func BenchmarkVector3Operations(b *testing.B) {
	v1 := NewVector3(Float32(1.0), Float32(2.0), Float32(3.0))
	v2 := NewVector3(Float32(4.0), Float32(5.0), Float32(6.0))

	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v1.Add(v2)
		}
	})

	b.Run("Cross", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v1.Cross(v2)
		}
	})

	b.Run("Normalize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v1.Normalize()
		}
	})
}
