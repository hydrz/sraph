package geom

import (
	"math"
	"testing"
)

func TestQuaternion_NewQuaternionFromAxisAngle(t *testing.T) {
	axis := NewVector3[Scalar](0.0, 0.0, 1.0) // Z-axis
	angle := Radians(math.Pi / 2)             // 90 degrees

	q := NewQuaternionFromAxisAngle(axis, angle)

	// For 90 degree rotation around Z-axis: q = (0, 0, sin(π/4), cos(π/4)) = (0, 0, √2/2, √2/2)
	expected := Scalar(math.Sqrt(2.0) / 2.0)
	if !NearlyEqual(q.X(), 0.0) || !NearlyEqual(q.Y(), 0.0) ||
		!NearlyEqual(q.Z(), expected) || !NearlyEqual(q.W(), expected) {
		t.Errorf("NewQuaternionFromAxisAngle() = (%v, %v, %v, %v), want (0, 0, %v, %v)",
			q.X(), q.Y(), q.Z(), q.W(), expected, expected)
	}
}

func TestQuaternion_QuaternionArithmetic(t *testing.T) {
	q1 := NewQuaternion(1.0, 2.0, 3.0, 4.0)
	q2 := NewQuaternion(5.0, 6.0, 7.0, 8.0)

	t.Run("Add", func(t *testing.T) {
		result := q1.Add(q2)
		if result.X() != 6.0 || result.Y() != 8.0 || result.Z() != 10.0 || result.W() != 12.0 {
			t.Errorf("Add() components mismatch")
		}
	})

	t.Run("Sub", func(t *testing.T) {
		result := q2.Sub(q1)
		if result.X() != 4.0 || result.Y() != 4.0 || result.Z() != 4.0 || result.W() != 4.0 {
			t.Errorf("Sub() components mismatch")
		}
	})

	t.Run("Scale", func(t *testing.T) {
		result := q1.Scale(2.0)
		if result.X() != 2.0 || result.Y() != 4.0 || result.Z() != 6.0 || result.W() != 8.0 {
			t.Errorf("Scale() components mismatch")
		}
	})

	t.Run("Neg", func(t *testing.T) {
		result := q1.Neg()
		if result.X() != Scalar(-1.0) || result.Y() != Scalar(-2.0) || result.Z() != Scalar(-3.0) || result.W() != Scalar(-4.0) {
			t.Errorf("Neg() components mismatch")
		}
	})
}

func TestQuaternion_QuaternionMultiplication(t *testing.T) {
	// Test Hamilton product with identity quaternion
	identity := NewQuaternion(0.0, 0.0, 0.0, 1.0)
	q := NewQuaternion(1.0, 2.0, 3.0, 4.0)

	result := q.Mul(identity)
	if !result.Equal(q) {
		t.Errorf("Quaternion multiplication with identity failed")
	}
}

func TestQuaternion_QuaternionLength(t *testing.T) {
	q := NewQuaternion(1.0, 2.0, 3.0, 4.0)
	length := q.Length()
	expected := Scalar(math.Sqrt(30.0)) // sqrt(1+4+9+16) = sqrt(30)
	if !NearlyEqual(length, expected) {
		t.Errorf("Length() = %v, want %v", length, expected)
	}
}

func TestQuaternion_QuaternionDotProduct(t *testing.T) {
	q1 := NewQuaternion(1.0, 2.0, 3.0, 4.0)
	q2 := NewQuaternion(5.0, 6.0, 7.0, 8.0)

	dot := q1.Dot(q2)
	expected := 70.0 // 1*5 + 2*6 + 3*7 + 4*8 = 70
	if !NearlyEqual(dot, expected) {
		t.Errorf("Dot() = %v, want %v", dot, expected)
	}
}

func TestQuaternion_QuaternionNormalize(t *testing.T) {
	t.Run("Normal quaternion", func(t *testing.T) {
		q := NewQuaternion(1.0, 2.0, 3.0, 4.0)
		normalized := q.Normalize()

		// Should have length 1
		length := normalized.Length()
		if !NearlyEqual(length, 1.0) {
			t.Errorf("Normalized quaternion length = %v, want 1.0", length)
		}
	})

	t.Run("Zero quaternion", func(t *testing.T) {
		zero := NewQuaternion(0.0, 0.0, 0.0, 0.0)
		normalized := zero.Normalize()

		// Should return zero quaternion
		if !NearlyEqual(normalized.X(), 0.0) || !NearlyEqual(normalized.Y(), 0.0) ||
			!NearlyEqual(normalized.Z(), 0.0) || !NearlyEqual(normalized.W(), 0.0) {
			t.Errorf("Normalized zero quaternion should be zero")
		}
	})
}

func TestQuaternion_QuaternionInvert(t *testing.T) {
	q := NewQuaternion(1.0, 2.0, 3.0, 4.0)
	inv := q.Invert()

	// q * q^-1 should be close to identity
	identity := q.Mul(inv)
	identityNormalized := identity.Normalize()

	// Check if it's close to identity quaternion (0, 0, 0, 1) or (0, 0, 0, -1)
	if !NearlyEqual(identityNormalized.X(), 0.0) || !NearlyEqual(identityNormalized.Y(), 0.0) ||
		!NearlyEqual(identityNormalized.Z(), 0.0) {
		t.Errorf("Quaternion inversion failed")
	}
}

func TestQuaternion_QuaternionEqual(t *testing.T) {
	q1 := NewQuaternion(1.0, 2.0, 3.0, 4.0)
	q2 := NewQuaternion(1.0, 2.0, 3.0, 4.0)
	q3 := NewQuaternion(1.0, 2.0, 3.0, 5.0)

	if !q1.Equal(q2) {
		t.Error("Eq quaternions should be Eq")
	}
	if q1.Equal(q3) {
		t.Error("Different quaternions should not be Eq")
	}
}

func TestQuaternion_QuaternionSlerp(t *testing.T) {
	q1 := NewQuaternionFromAxisAngle(NewVector3[Scalar](0.0, 0.0, 1.0), Radians(0.0))       // 0° rotation around Z
	q2 := NewQuaternionFromAxisAngle(NewVector3[Scalar](0.0, 0.0, 1.0), Radians(math.Pi/4)) // 45° rotation around Z

	q3 := q1.Slerp(q2, 0.5) // Slerp at t=0.5

	expected := NewQuaternionFromAxisAngle(NewVector3[Scalar](0.0, 0.0, 1.0), Radians(math.Pi/8)) // 22.5° rotation around Z

	if !q3.Equal(expected) {
		t.Errorf("Slerp() = %v, want %v", q3, expected)
	}
}

func TestQuaternion_QuaternionRotateVector3(t *testing.T) {
	// Test 90 degree rotation around Z axis
	axis := NewVector3[Scalar](0.0, 0.0, 1.0)
	angle := Radians(math.Pi / 2)
	q := NewQuaternionFromAxisAngle(axis, angle)

	// Rotate vector (1, 0, 0) around Z axis by 90 degrees
	v := NewVector3[Scalar](1.0, 0.0, 0.0)
	rotated := q.RotateVector3(v)

	// Should result in approximately (0, 1, 0)
	if !NearlyEqual(rotated.X(), 0.0) || !NearlyEqual(rotated.Y(), 1.0) || !NearlyEqual(rotated.Z(), 0.0) {
		t.Errorf("RotateVector3() = (%v, %v, %v), want (0, 1, 0)", rotated.X(), rotated.Y(), rotated.Z())
	}
}

func TestQuaternion_QuaternionString(t *testing.T) {
	q := NewQuaternion(1.0, 2.0, 3.0, 4.0)
	str := q.String()
	expected := "(1, 2, 3, 4)"
	if str != expected {
		t.Errorf("String() = %v, want %v", str, expected)
	}
}

func BenchmarkQuaternionOperations(b *testing.B) {
	q1 := NewQuaternion(1.0, 2.0, 3.0, 4.0)
	q2 := NewQuaternion(5.0, 6.0, 7.0, 8.0)
	v := NewVector3[Scalar](1.0, 0.0, 0.0)

	b.Run("Mul", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			q1.Mul(q2)
		}
	})

	b.Run("Normalize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			q1.Normalize()
		}
	})

	b.Run("RotateVector3", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			q1.RotateVector3(v)
		}
	})
}
