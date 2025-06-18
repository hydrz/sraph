package geom

import (
	"math"
	"testing"
	"time"
)

func TestQuaternion_NewQuaternionFromAxisAngle(t *testing.T) {
	axis := Vector3[F32]{0.0, 0.0, 1.0}   // Z-axis
	angle := NewRadians[F32](math.Pi / 2) // 90 degrees

	q := NewQuaternionFromAxisAngle(axis, angle)

	// For 90 degree rotation around Z-axis: q = (0, 0, sin(π/4), cos(π/4)) = (0, 0, √2/2, √2/2)
	expected := F32(math.Sqrt(2.0) / 2.0)
	if !Eq(q.X, F32(0.0)) || !Eq(q.Y, F32(0.0)) ||
		!Eq(q.Z, expected) || !Eq(q.W, expected) {
		t.Errorf("NewQuaternionFromAxisAngle() = (%v, %v, %v, %v), want (0, 0, %v, %v)",
			q.X, q.Y, q.Z, q.W, expected, expected)
	}
}

func TestQuaternion_QuaternionArithmetic(t *testing.T) {
	q1 := Quaternion[F32]{1.0, 2.0, 3.0, 4.0}
	q2 := Quaternion[F32]{5.0, 6.0, 7.0, 8.0}

	t.Run("Add", func(t *testing.T) {
		result := q1.Add(q2)
		if result.X != 6.0 || result.Y != 8.0 || result.Z != 10.0 || result.W != 12.0 {
			t.Errorf("Add() components mismatch")
		}
	})

	t.Run("Sub", func(t *testing.T) {
		result := q2.Sub(q1)
		if result.X != 4.0 || result.Y != 4.0 || result.Z != 4.0 || result.W != 4.0 {
			t.Errorf("Sub() components mismatch")
		}
	})

	t.Run("Scale", func(t *testing.T) {
		result := q1.Scale(F32(2.0))
		if result.X != 2.0 || result.Y != 4.0 || result.Z != 6.0 || result.W != 8.0 {
			t.Errorf("Scale() components mismatch")
		}
	})

	t.Run("Neg", func(t *testing.T) {
		result := q1.Neg()
		if result.X != -1.0 || result.Y != -2.0 || result.Z != -3.0 || result.W != -4.0 {
			t.Errorf("Neg() components mismatch")
		}
	})
}

func TestQuaternion_QuaternionMultiplication(t *testing.T) {
	// Test Hamilton product with identity quaternion
	identity := Quaternion[F32]{0.0, 0.0, 0.0, 1.0}
	q := Quaternion[F32]{1.0, 2.0, 3.0, 4.0}

	result := q.Mul(identity)
	if !result.Eq(q) {
		t.Errorf("Quaternion multiplication with identity failed")
	}
}

func TestQuaternion_QuaternionLength(t *testing.T) {
	q := Quaternion[F32]{1.0, 2.0, 3.0, 4.0}
	length := q.Length()
	expected := F32(math.Sqrt(30.0)) // sqrt(1+4+9+16) = sqrt(30)
	if !Eq(length, expected) {
		t.Errorf("Length() = %v, want %v", length, expected)
	}
}

func TestQuaternion_QuaternionDotProduct(t *testing.T) {
	q1 := Quaternion[F32]{1.0, 2.0, 3.0, 4.0}
	q2 := Quaternion[F32]{5.0, 6.0, 7.0, 8.0}

	dot := q1.Dot(q2)
	expected := F32(70.0) // 1*5 + 2*6 + 3*7 + 4*8 = 70
	if !Eq(dot, expected) {
		t.Errorf("Dot() = %v, want %v", dot, expected)
	}
}

func TestQuaternion_QuaternionNormalize(t *testing.T) {
	t.Run("Normal quaternion", func(t *testing.T) {
		q := Quaternion[F32]{1.0, 2.0, 3.0, 4.0}
		normalized := q.Normalize()

		// Should have length 1
		length := normalized.Length()
		if !Eq(length, F32(1.0)) {
			t.Errorf("Normalized quaternion length = %v, want 1.0", length)
		}
	})

	t.Run("Zero quaternion", func(t *testing.T) {
		zero := Quaternion[F32]{0.0, 0.0, 0.0, 0.0}
		normalized := zero.Normalize()

		// Should return zero quaternion
		if !Eq(normalized.X, F32(0.0)) || !Eq(normalized.Y, F32(0.0)) ||
			!Eq(normalized.Z, F32(0.0)) || !Eq(normalized.W, F32(0.0)) {
			t.Errorf("Normalized zero quaternion should be zero")
		}
	})
}

func TestQuaternion_QuaternionInvert(t *testing.T) {
	q := Quaternion[F32]{1.0, 2.0, 3.0, 4.0}
	inv := q.Invert()

	// q * q^-1 should be close to identity
	identity := q.Mul(inv)
	identityNormalized := identity.Normalize()

	// Check if it's close to identity quaternion (0, 0, 0, 1) or (0, 0, 0, -1)
	if !Eq(identityNormalized.X, F32(0.0)) || !Eq(identityNormalized.Y, F32(0.0)) ||
		!Eq(identityNormalized.Z, F32(0.0)) {
		t.Errorf("Quaternion inversion failed")
	}
}

func TestQuaternion_QuaternionEq(t *testing.T) {
	q1 := Quaternion[F32]{1.0, 2.0, 3.0, 4.0}
	q2 := Quaternion[F32]{1.0, 2.0, 3.0, 4.0}
	q3 := Quaternion[F32]{1.0, 2.0, 3.0, 5.0}

	if !q1.Eq(q2) {
		t.Error("Eq quaternions should be Eq")
	}
	if q1.Eq(q3) {
		t.Error("Different quaternions should not be Eq")
	}
}

func TestQuaternion_QuaternionSlerp(t *testing.T) {
	q1 := Quaternion[F32]{0.0, 0.0, 0.0, 1.0} // Identity
	q2 := Quaternion[F32]{0.0, 0.0, 1.0, 0.0} // 180° rotation around Z

	// Test interpolation at t=0 (should be q1)
	t0 := time.Unix(0, 0)
	result := q1.Slerp(q2, t0)

	// The result should be close to q1
	if !Eq(result.X, q1.X) || !Eq(result.Y, q1.Y) ||
		!Eq(result.Z, q1.Z) || !Eq(result.W, q1.W) {
		t.Errorf("Slerp at t=0 should return first quaternion")
	}
}

func TestQuaternion_QuaternionRotateVector3(t *testing.T) {
	// Test 90 degree rotation around Z axis
	axis := Vector3[F32]{0.0, 0.0, 1.0}
	angle := NewRadians[F32](math.Pi / 2)
	q := NewQuaternionFromAxisAngle(axis, angle)

	// Rotate vector (1, 0, 0) around Z axis by 90 degrees
	v := Vector3[F32]{1.0, 0.0, 0.0}
	rotated := q.RotateVector3(v)

	// Should result in approximately (0, 1, 0)
	if !Eq(rotated.X, F32(0.0)) || !Eq(rotated.Y, F32(1.0)) || !Eq(rotated.Z, F32(0.0)) {
		t.Errorf("RotateVector3() = (%v, %v, %v), want (0, 1, 0)", rotated.X, rotated.Y, rotated.Z)
	}
}

func TestQuaternion_QuaternionString(t *testing.T) {
	q := Quaternion[F32]{1.0, 2.0, 3.0, 4.0}
	str := q.String()
	expected := "(1, 2, 3, 4)"
	if str != expected {
		t.Errorf("String() = %v, want %v", str, expected)
	}
}

func BenchmarkQuaternionOperations(b *testing.B) {
	q1 := Quaternion[F32]{1.0, 2.0, 3.0, 4.0}
	q2 := Quaternion[F32]{5.0, 6.0, 7.0, 8.0}
	v := Vector3[F32]{1.0, 0.0, 0.0}

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
