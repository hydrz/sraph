package geom

import (
	"math"
	"testing"
)

func TestNewMatrix(t *testing.T) {
	m := NewMatrix[Float32]()

	// Should be identity matrix
	if !m.IsIdentity() {
		t.Error("NewMatrix should create identity matrix")
	}
}

func TestMatrixAt(t *testing.T) {
	m := NewMatrix[Float32]()

	// Test identity matrix values
	if m.At(0, 0) != 1.0 || m.At(1, 1) != 1.0 || m.At(2, 2) != 1.0 || m.At(3, 3) != 1.0 {
		t.Error("Identity matrix diagonal should be 1.0")
	}

	if m.At(0, 1) != 0.0 || m.At(1, 0) != 0.0 {
		t.Error("Identity matrix off-diagonal should be 0.0")
	}
}

func TestMatrixSet(t *testing.T) {
	m := NewMatrix[Float32]()
	m.Set(0, 1, Float32(5.0))

	if m.At(0, 1) != 5.0 {
		t.Errorf("Set(0, 1, 5.0) failed, got %v", m.At(0, 1))
	}
}

func TestMatrixArithmetic(t *testing.T) {
	m1 := NewMatrix[Float32]()
	m2 := NewMatrix[Float32]()

	// Set some test values
	m1.Set(0, 0, Float32(2.0))
	m2.Set(0, 0, Float32(3.0))

	t.Run("Add", func(t *testing.T) {
		result := m1.Add(m2)
		if result.At(0, 0) != 5.0 {
			t.Errorf("Add() failed, got %v, want 5.0", result.At(0, 0))
		}
	})

	t.Run("Sub", func(t *testing.T) {
		result := m2.Sub(m1)
		if result.At(0, 0) != 1.0 {
			t.Errorf("Sub() failed, got %v, want 1.0", result.At(0, 0))
		}
	})
}

func TestMatrixMultiplication(t *testing.T) {
	identity := NewMatrix[Float32]()
	m := NewMatrix[Float32]()
	m.Set(0, 3, Float32(5.0)) // Translation

	// Multiply by identity should return original
	result := m.Mul(identity)
	if result.At(0, 3) != 5.0 {
		t.Errorf("Matrix multiplication with identity failed")
	}
}

func TestMatrixProperties(t *testing.T) {
	identity := NewMatrix[Float32]()

	t.Run("IsIdentity", func(t *testing.T) {
		if !identity.IsIdentity() {
			t.Error("Identity matrix should return true for IsIdentity()")
		}
	})

	t.Run("IsFinite", func(t *testing.T) {
		if !identity.IsFinite() {
			t.Error("Identity matrix should be finite")
		}
	})

	t.Run("IsInvertible", func(t *testing.T) {
		if !identity.IsInvertible() {
			t.Error("Identity matrix should be invertible")
		}
	})

	t.Run("Determinant", func(t *testing.T) {
		det := identity.Determinant()
		if !Equal(det, Float32(1.0)) {
			t.Errorf("Identity matrix determinant = %v, want 1.0", det)
		}
	})
}

func TestMatrixTransformations(t *testing.T) {
	t.Run("Translation", func(t *testing.T) {
		m := NewMatrix[Float32]()
		translation := NewVector3(Float32(1.0), Float32(2.0), Float32(3.0))
		result := m.Translate3D(translation)

		if result.At(0, 3) != 1.0 || result.At(1, 3) != 2.0 || result.At(2, 3) != 3.0 {
			t.Error("Translation matrix creation failed")
		}
	})

	t.Run("Scale", func(t *testing.T) {
		m := NewMatrix[Float32]()
		scale := NewVector3(Float32(2.0), Float32(3.0), Float32(4.0))
		result := m.Scale3D(scale)

		if result.At(0, 0) != 2.0 || result.At(1, 1) != 3.0 || result.At(2, 2) != 4.0 {
			t.Error("Scale matrix creation failed")
		}
	})

	t.Run("RotationX", func(t *testing.T) {
		m := NewMatrix[Float32]()
		angle := Radians(math.Pi / 2) // 90 degrees
		result := m.RotateX(angle)

		// Check rotation matrix properties
		if !Equal(result.At(0, 0), Float32(1.0)) {
			t.Error("X rotation should not affect X axis")
		}
		if !Equal(result.At(1, 1), Float32(0.0)) {
			t.Error("90 degree X rotation Y component should be 0")
		}
	})
}

func TestMatrixInverse(t *testing.T) {
	identity := NewMatrix[Float32]()

	inv, err := identity.Inverse()
	if err != nil {
		t.Errorf("Identity matrix should be invertible: %v", err)
	}

	if !inv.IsIdentity() {
		t.Error("Inverse of identity should be identity")
	}
}

func TestMatrixTranspose(t *testing.T) {
	m := NewMatrix[Float32]()
	m.Set(0, 1, Float32(5.0))

	transposed := m.Transpose()
	if transposed.At(1, 0) != 5.0 {
		t.Error("Transpose failed")
	}
}

func TestMatrixBasisVectors(t *testing.T) {
	identity := NewMatrix[Float32]()
	basis := identity.BasisVectors()

	if len(basis) != 3 {
		t.Errorf("Expected 3 basis vectors, got %d", len(basis))
	}

	// Check identity matrix basis vectors
	if basis[0].X() != 1.0 || basis[0].Y() != 0.0 || basis[0].Z() != 0.0 {
		t.Error("X basis vector incorrect")
	}
	if basis[1].X() != 0.0 || basis[1].Y() != 1.0 || basis[1].Z() != 0.0 {
		t.Error("Y basis vector incorrect")
	}
	if basis[2].X() != 0.0 || basis[2].Y() != 0.0 || basis[2].Z() != 1.0 {
		t.Error("Z basis vector incorrect")
	}
}

func TestMatrixScale(t *testing.T) {
	m := NewMatrix[Float32]()
	m.Set(0, 0, Float32(2.0))
	m.Set(1, 1, Float32(3.0))
	m.Set(2, 2, Float32(4.0))

	scale := m.Scale()
	if !Equal(scale.X(), Float32(2.0)) || !Equal(scale.Y(), Float32(3.0)) || !Equal(scale.Z(), Float32(4.0)) {
		t.Errorf("Scale() = (%v, %v, %v), want (2.0, 3.0, 4.0)", scale.X(), scale.Y(), scale.Z())
	}
}

func TestMatrixTransformPoint(t *testing.T) {
	m := NewMatrix[Float32]()
	point := NewPoint(Float32(1.0), Float32(2.0))

	// Identity transform should return same point
	result := m.TransformPoint(point)
	if !Equal(result.X(), point.X()) || !Equal(result.Y(), point.Y()) {
		t.Errorf("Identity transform failed")
	}
}

func TestMatrixTransformVector(t *testing.T) {
	m := NewMatrix[Float32]()
	v3 := NewVector3(Float32(1.0), Float32(2.0), Float32(3.0))

	// Identity transform should return same vector
	result := m.TransformVector3D(v3)
	if !Equal(result.X(), v3.X()) || !Equal(result.Y(), v3.Y()) || !Equal(result.Z(), v3.Z()) {
		t.Errorf("Identity vector transform failed")
	}
}

func TestMatrixAnalysis(t *testing.T) {
	identity := NewMatrix[Float32]()

	t.Run("IsAffine", func(t *testing.T) {
		if !identity.IsAffine() {
			t.Error("Identity matrix should be affine")
		}
	})

	t.Run("HasTranslation", func(t *testing.T) {
		if identity.HasTranslation() {
			t.Error("Identity matrix should not have translation")
		}
	})

	t.Run("IsAxisAligned", func(t *testing.T) {
		if !identity.IsAxisAligned() {
			t.Error("Identity matrix should be axis aligned")
		}
	})

	t.Run("IsTranslationOnly", func(t *testing.T) {
		if !identity.IsTranslationOnly() {
			t.Error("Identity matrix should be considered translation-only")
		}
	})
}

func TestMatrixDecompose(t *testing.T) {
	m := NewMatrix[Float32]()
	decomp := m.Decompose()

	// Check identity decomposition
	if !Equal(decomp.Translation.X(), Float32(0.0)) ||
		!Equal(decomp.Translation.Y(), Float32(0.0)) ||
		!Equal(decomp.Translation.Z(), Float32(0.0)) {
		t.Error("Identity matrix should have zero translation")
	}

	if !Equal(decomp.Scale.X(), Float32(1.0)) ||
		!Equal(decomp.Scale.Y(), Float32(1.0)) ||
		!Equal(decomp.Scale.Z(), Float32(1.0)) {
		t.Error("Identity matrix should have unit scale")
	}
}

func TestMatrixCosSin(t *testing.T) {
	m := NewMatrix[Float32]()

	// Test 90 degrees
	cos, sin := m.CosSin(Radians(math.Pi / 2))
	if !Equal(cos, Float32(0.0)) || !Equal(sin, Float32(1.0)) {
		t.Errorf("CosSin(π/2) = (%v, %v), want (0.0, 1.0)", cos, sin)
	}

	// Test 0 degrees
	cos, sin = m.CosSin(Radians(0))
	if !Equal(cos, Float32(1.0)) || !Equal(sin, Float32(0.0)) {
		t.Errorf("CosSin(0) = (%v, %v), want (1.0, 0.0)", cos, sin)
	}
}

func TestMatrixLookAt(t *testing.T) {
	m := NewMatrix[Float32]()
	position := NewVector3(Float32(0.0), Float32(0.0), Float32(1.0))
	target := NewVector3(Float32(0.0), Float32(0.0), Float32(0.0))
	up := NewVector3(Float32(0.0), Float32(1.0), Float32(0.0))

	view := m.LookAt(position, target, up)

	// Check that the matrix is not identity
	if view.IsIdentity() {
		t.Error("LookAt matrix should not be identity")
	}
}

func TestMatrixOrthographic(t *testing.T) {
	m := NewMatrix[Float32]()
	size := NewSize(Float32(800.0), Float32(600.0))

	ortho := m.Orthographic(size)

	// Check that the matrix is not identity
	if ortho.IsIdentity() {
		t.Error("Orthographic matrix should not be identity")
	}
}

func TestMatrixQuaternionRotation(t *testing.T) {
	m := NewMatrix[Float32]()
	quat := NewQuaternion(Float32(0.0), Float32(0.0), Float32(0.0), Float32(1.0))

	rotMatrix := m.RotateQuaternion(quat)

	// Identity quaternion should produce identity rotation matrix
	if !rotMatrix.IsIdentity() {
		t.Error("Identity quaternion should produce identity rotation matrix")
	}
}

func TestMatrixAxisAngleRotation(t *testing.T) {
	m := NewMatrix[Float32]()
	axis := NewVector3(Float32(0.0), Float32(0.0), Float32(1.0))
	angle := Radians(0.0)

	rotMatrix := m.RotateAxisAngle(angle, axis)

	// Zero angle should produce identity matrix
	if !rotMatrix.IsIdentity() {
		t.Error("Zero angle rotation should produce identity matrix")
	}
}

func TestMatrixEqual(t *testing.T) {
	m1 := NewMatrix[Float32]()
	m2 := NewMatrix[Float32]()

	if !m1.Equal(m2) {
		t.Error("Two identity matrices should be equal")
	}

	m2.Set(0, 0, Float32(2.0))
	if m1.Equal(m2) {
		t.Error("Modified matrix should not equal identity")
	}
}

func BenchmarkMatrixOperations(b *testing.B) {
	m1 := NewMatrix[Float32]()
	m2 := NewMatrix[Float32]()

	// Add some values to make operations more meaningful
	m1.Set(0, 3, Float32(1.0))
	m2.Set(1, 3, Float32(2.0))

	b.Run("Mul", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m1.Mul(m2)
		}
	})

	b.Run("Determinant", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m1.Determinant()
		}
	})

	b.Run("Transpose", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m1.Transpose()
		}
	})

	b.Run("Inverse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m1.Inverse()
		}
	})
}
