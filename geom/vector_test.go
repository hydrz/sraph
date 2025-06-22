package geom

import (
	"math"
	"testing"
)

func TestVector3Arithmetic(t *testing.T) {
	v1 := vector3[Scalar]{1.0, 2.0, 3.0}
	v2 := vector3[Scalar]{4.0, 5.0, 6.0}

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
		if !NearlyEqual(result.X(), Scalar(4.0)) || !NearlyEqual(result.Y(), Scalar(2.5)) || !NearlyEqual(result.Z(), Scalar(2.0)) {
			t.Errorf("Div() = (%v, %v, %v), want (4.0, 2.5, 2.0)", result.X(), result.Y(), result.Z())
		}
	})

	t.Run("Scale", func(t *testing.T) {
		result := v1.Scale(Scalar(2.0))
		if result.X() != 2.0 || result.Y() != 4.0 || result.Z() != 6.0 {
			t.Errorf("Scale(2.0) = (%v, %v, %v), want (2.0, 4.0, 6.0)", result.X(), result.Y(), result.Z())
		}
	})
}

func TestVector3VectorOperations(t *testing.T) {
	v := vector3[Scalar]{3.0, 4.0, 0.0}

	t.Run("Length", func(t *testing.T) {
		length := v.Length()
		expected := Scalar(5.0) // sqrt(3^2 + 4^2 + 0^2) = 5
		if !NearlyEqual(length, expected) {
			t.Errorf("Length() = %v, want %v", length, expected)
		}
	})

	t.Run("Normalize", func(t *testing.T) {
		normalized := v.Normalize()

		// Should have length 1
		length := normalized.Length()
		if !NearlyEqual(length, Scalar(1.0)) {
			t.Errorf("Normalized vector length = %v, want 1.0", length)
		}

		// Check components
		expectedX := Scalar(0.6) // 3/5
		expectedY := Scalar(0.8) // 4/5
		if !NearlyEqual(normalized.X(), expectedX) || !NearlyEqual(normalized.Y(), expectedY) || !NearlyEqual(normalized.Z(), Scalar(0.0)) {
			t.Errorf("Normalize() = (%v, %v, %v), want (%v, %v, 0.0)", normalized.X(), normalized.Y(), normalized.Z(), expectedX, expectedY)
		}
	})

	t.Run("Zero vector normalize", func(t *testing.T) {
		zero := vector3[Scalar]{0.0, 0.0, 0.0}
		normalized := zero.Normalize()

		// Should return zero vector
		if !NearlyEqual(normalized.X(), Scalar(0.0)) || !NearlyEqual(normalized.Y(), Scalar(0.0)) || !NearlyEqual(normalized.Z(), Scalar(0.0)) {
			t.Errorf("Normalized zero vector should be zero")
		}
	})
}

func TestVector3DotProduct(t *testing.T) {
	v1 := vector3[Scalar]{1.0, 2.0, 3.0}
	v2 := vector3[Scalar]{4.0, 5.0, 6.0}

	dot := v1.Dot(v2)
	expected := Scalar(32.0) // 1*4 + 2*5 + 3*6 = 32
	if !NearlyEqual(dot, expected) {
		t.Errorf("Dot() = %v, want %v", dot, expected)
	}
}

func TestVector3CrossProduct(t *testing.T) {
	v1 := vector3[Scalar]{1.0, 0.0, 0.0}
	v2 := vector3[Scalar]{0.0, 1.0, 0.0}

	cross := v1.Cross(v2)
	// (1,0,0) × (0,1,0) = (0,0,1)
	if !NearlyEqual(cross.X(), Scalar(0.0)) || !NearlyEqual(cross.Y(), Scalar(0.0)) || !NearlyEqual(cross.Z(), Scalar(1.0)) {
		t.Errorf("Cross() = (%v, %v, %v), want (0.0, 0.0, 1.0)", cross.X(), cross.Y(), cross.Z())
	}
}

func TestVector3MathFunctions(t *testing.T) {
	v := vector3[Scalar]{1.7, 2.3, 3.9}

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
		negV := vector3[Scalar]{-1.0, -2.0, -3.0}
		result := negV.Abs()
		if result.X() != 1.0 || result.Y() != 2.0 || result.Z() != 3.0 {
			t.Errorf("Abs() = (%v, %v, %v), want (1.0, 2.0, 3.0)", result.X(), result.Y(), result.Z())
		}
	})
}

func TestVector3Lerp(t *testing.T) {
	v1 := vector3[Scalar]{0.0, 0.0, 0.0}
	v2 := vector3[Scalar]{10.0, 20.0, 30.0}

	tests := []struct {
		t        Scalar
		expected vector3[Scalar]
	}{
		{Scalar(0.0), v1},
		{Scalar(1.0), v2},
		{Scalar(0.5), vector3[Scalar]{5.0, 10.0, 15.0}},
		{Scalar(0.25), vector3[Scalar]{2.5, 5.0, 7.5}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := v1.Lerp(v2, tt.t)
			if !NearlyEqual(result.X(), tt.expected.X()) ||
				!NearlyEqual(result.Y(), tt.expected.Y()) ||
				!NearlyEqual(result.Z(), tt.expected.Z()) {
				t.Errorf("Lerp(%v) = (%v, %v, %v), want (%v, %v, %v)",
					tt.t, result.X(), result.Y(), result.Z(),
					tt.expected.X(), tt.expected.Y(), tt.expected.Z())
			}
		})
	}
}

func TestVector3Combine(t *testing.T) {
	v1 := vector3[Scalar]{1.0, 2.0, 3.0}
	v2 := vector3[Scalar]{4.0, 5.0, 6.0}

	result := v1.Combine(v2, Scalar(2.0))
	// v1 + v2 * 2.0 = (1,2,3) + (8,10,12) = (9,12,15)
	if result.X() != 9.0 || result.Y() != 12.0 || result.Z() != 15.0 {
		t.Errorf("Combine() = (%v, %v, %v), want (9.0, 12.0, 15.0)", result.X(), result.Y(), result.Z())
	}
}

func TestVector3Equal(t *testing.T) {
	v1 := vector3[Scalar]{1.0, 2.0, 3.0}
	v2 := vector3[Scalar]{1.0, 2.0, 3.0}
	v3 := vector3[Scalar]{1.0, 2.0, 4.0}

	if !v1.Equal(v2) {
		t.Error("Eq vectors should be Eq")
	}
	if v1.Equal(v3) {
		t.Error("Different vectors should not be Eq")
	}
}

func TestVector3String(t *testing.T) {
	v := vector3[Scalar]{1.0, 2.0, 3.0}
	str := v.String()
	expected := "(1, 2, 3)"
	if str != expected {
		t.Errorf("String() = %v, want %v", str, expected)
	}
}

func TestVector4Length(t *testing.T) {
	v := vector4[Scalar]{1.0, 2.0, 3.0, 4.0}
	length := v.Length()
	expected := Scalar(math.Sqrt(30.0)) // sqrt(1^2 + 2^2 + 3^2 + 4^2) = sqrt(30)
	if !NearlyEqual(length, expected) {
		t.Errorf("Length() = %v, want %v", length, expected)
	}
}

func TestVector4IsFinite(t *testing.T) {
	finite := vector4[Scalar]{1.0, 2.0, 3.0, 4.0}
	infinite := vector4[Scalar]{Scalar(math.Inf(1)), 2.0, 3.0, 4.0}

	if !finite.IsFinite() {
		t.Error("Finite vector should return true for IsFinite()")
	}
	if infinite.IsFinite() {
		t.Error("Infinite vector should return false for IsFinite()")
	}
}

func TestVector4String(t *testing.T) {
	v := vector4[Scalar]{1.0, 2.0, 3.0, 4.0}
	str := v.String()
	expected := "(1, 2, 3, 4)"
	if str != expected {
		t.Errorf("String() = %v, want %v", str, expected)
	}
}

func BenchmarkVector3Operations(b *testing.B) {
	v1 := vector3[Scalar]{1.0, 2.0, 3.0}
	v2 := vector3[Scalar]{4.0, 5.0, 6.0}

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
