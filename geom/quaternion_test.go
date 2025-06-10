package geom

import (
	"math"
	"testing"
	"time"
)

func TestNewQuaternion(t *testing.T) {
	q := NewQuaternion(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	if q.X() != 1.0 || q.Y() != 2.0 || q.Z() != 3.0 || q.W() != 4.0 {
		t.Errorf("NewQuaternion() components mismatch")
	}
}

func TestNewQuaternionFromAxisAngle(t *testing.T) {
	axis := NewVector3(Float32(0.0), Float32(0.0), Float32(1.0)) // Z-axis
	angle := Radians(math.Pi / 2)                                // 90 degrees

	q := NewQuaternionFromAxisAngle(axis, angle)

	// For 90 degree rotation around Z-axis: q = (0, 0, sin(π/4), cos(π/4)) = (0, 0, √2/2, √2/2)
	expected := Float32(math.Sqrt(2.0) / 2.0)
	if !Equal(q.X(), Float32(0.0)) || !Equal(q.Y(), Float32(0.0)) ||
		!Equal(q.Z(), expected) || !Equal(q.W(), expected) {
		t.Errorf("NewQuaternionFromAxisAngle() = (%v, %v, %v, %v), want (0, 0, %v, %v)",
			q.X(), q.Y(), q.Z(), q.W(), expected, expected)
	}
}

func TestQuaternionArithmetic(t *testing.T) {
	q1 := NewQuaternion(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	q2 := NewQuaternion(Float32(5.0), Float32(6.0), Float32(7.0), Float32(8.0))

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
		result := q1.Scale(Float32(2.0))
		if result.X() != 2.0 || result.Y() != 4.0 || result.Z() != 6.0 || result.W() != 8.0 {
			t.Errorf("Scale() components mismatch")
		}
	})

	t.Run("Neg", func(t *testing.T) {
		result := q1.Neg()
		if result.X() != -1.0 || result.Y() != -2.0 || result.Z() != -3.0 || result.W() != -4.0 {
			t.Errorf("Neg() components mismatch")
		}
	})
}

func TestQuaternionMultiplication(t *testing.T) {
	// Test Hamilton product with identity quaternion
	identity := NewQuaternion(Float32(0.0), Float32(0.0), Float32(0.0), Float32(1.0))
	q := NewQuaternion(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))

	result := q.Mul(identity)
	if !result.Equal(q) {
		t.Errorf("Quaternion multiplication with identity failed")
	}
}

func TestQuaternionLength(t *testing.T) {
	q := NewQuaternion(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	length := q.Length()
	expected := Float32(math.Sqrt(30.0)) // sqrt(1+4+9+16) = sqrt(30)
	if !Equal(length, expected) {
		t.Errorf("Length() = %v, want %v", length, expected)
	}
}

func TestQuaternionDotProduct(t *testing.T) {
	q1 := NewQuaternion(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	q2 := NewQuaternion(Float32(5.0), Float32(6.0), Float32(7.0), Float32(8.0))

	dot := q1.Dot(q2)
	expected := Float32(70.0) // 1*5 + 2*6 + 3*7 + 4*8 = 70
	if !Equal(dot, expected) {
		t.Errorf("Dot() = %v, want %v", dot, expected)
	}
}

func TestQuaternionNormalize(t *testing.T) {
	t.Run("Normal quaternion", func(t *testing.T) {
		q := NewQuaternion(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
		normalized := q.Normalize()

		// Should have length 1
		length := normalized.Length()
		if !Equal(length, Float32(1.0)) {
			t.Errorf("Normalized quaternion length = %v, want 1.0", length)
		}
	})

	t.Run("Zero quaternion", func(t *testing.T) {
		zero := NewQuaternion(Float32(0.0), Float32(0.0), Float32(0.0), Float32(0.0))
		normalized := zero.Normalize()

		// Should return zero quaternion
		if !Equal(normalized.X(), Float32(0.0)) || !Equal(normalized.Y(), Float32(0.0)) ||
			!Equal(normalized.Z(), Float32(0.0)) || !Equal(normalized.W(), Float32(0.0)) {
			t.Errorf("Normalized zero quaternion should be zero")
		}
	})
}

func TestQuaternionInvert(t *testing.T) {
	q := NewQuaternion(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	inv := q.Invert()

	// q * q^-1 should be close to identity
	identity := q.Mul(inv)
	identityNormalized := identity.Normalize()

	// Check if it's close to identity quaternion (0, 0, 0, 1) or (0, 0, 0, -1)
	if !Equal(identityNormalized.X(), Float32(0.0)) || !Equal(identityNormalized.Y(), Float32(0.0)) ||
		!Equal(identityNormalized.Z(), Float32(0.0)) {
		t.Errorf("Quaternion inversion failed")
	}
}

func TestQuaternionEqual(t *testing.T) {
	q1 := NewQuaternion(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	q2 := NewQuaternion(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	q3 := NewQuaternion(Float32(1.0), Float32(2.0), Float32(3.0), Float32(5.0))

	if !q1.Equal(q2) {
		t.Error("Equal quaternions should be equal")
	}
	if q1.Equal(q3) {
		t.Error("Different quaternions should not be equal")
	}
}

func TestQuaternionSlerp(t *testing.T) {
	q1 := NewQuaternion(Float32(0.0), Float32(0.0), Float32(0.0), Float32(1.0)) // Identity
	q2 := NewQuaternion(Float32(0.0), Float32(0.0), Float32(1.0), Float32(0.0)) // 180° rotation around Z

	// Test interpolation at t=0 (should be q1)
	t0 := time.Unix(0, 0)
	result := q1.Slerp(q2, t0)

	// The result should be close to q1
	if !Equal(result.X(), q1.X()) || !Equal(result.Y(), q1.Y()) ||
		!Equal(result.Z(), q1.Z()) || !Equal(result.W(), q1.W()) {
		t.Errorf("Slerp at t=0 should return first quaternion")
	}
}

func TestQuaternionRotateVector3(t *testing.T) {
	// Test 90 degree rotation around Z axis
	axis := NewVector3(Float32(0.0), Float32(0.0), Float32(1.0))
	angle := Radians(math.Pi / 2)
	q := NewQuaternionFromAxisAngle(axis, angle)

	// Rotate vector (1, 0, 0) around Z axis by 90 degrees
	v := NewVector3(Float32(1.0), Float32(0.0), Float32(0.0))
	rotated := q.RotateVector3(v)

	// Should result in approximately (0, 1, 0)
	if !Equal(rotated.X(), Float32(0.0)) || !Equal(rotated.Y(), Float32(1.0)) || !Equal(rotated.Z(), Float32(0.0)) {
		t.Errorf("RotateVector3() = (%v, %v, %v), want (0, 1, 0)", rotated.X(), rotated.Y(), rotated.Z())
	}
}

func TestQuaternionString(t *testing.T) {
	q := NewQuaternion(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	str := q.String()
	expected := "Quaternion(1, 2, 3, 4)"
	if str != expected {
		t.Errorf("String() = %v, want %v", str, expected)
	}
}

func BenchmarkQuaternionOperations(b *testing.B) {
	q1 := NewQuaternion(Float32(1.0), Float32(2.0), Float32(3.0), Float32(4.0))
	q2 := NewQuaternion(Float32(5.0), Float32(6.0), Float32(7.0), Float32(8.0))
	v := NewVector3(Float32(1.0), Float32(0.0), Float32(0.0))

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
