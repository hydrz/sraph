package geom

import (
	"math"
	"testing"
)

func TestNewMatrix(t *testing.T) {
	m := NewMatrix()

	// Should be identity matrix
	if !m.IsIdentity() {
		t.Error("NewMatrix should create identity matrix")
	}
}

func TestMatrix_At(t *testing.T) {
	m := NewMatrix()

	// Test identity matrix values
	if m.At(0, 0) != 1.0 || m.At(1, 1) != 1.0 || m.At(2, 2) != 1.0 || m.At(3, 3) != 1.0 {
		t.Error("Identity matrix diagonal should be 1.0")
	}

	if m.At(0, 1) != 0.0 || m.At(1, 0) != 0.0 {
		t.Error("Identity matrix off-diagonal should be 0.0")
	}
}

func TestMatrix_Set(t *testing.T) {
	m := NewMatrix()
	m.Set(0, 1, 5.0)

	if m.At(0, 1) != 5.0 {
		t.Errorf("Set(0, 1, 5.0) failed, got %v", m.At(0, 1))
	}
}

func TestMatrix_Arithmetic(t *testing.T) {
	m1 := NewMatrix()
	m2 := NewMatrix()

	// Set some test values
	m1.Set(0, 0, 2.0)
	m2.Set(0, 0, 3.0)

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

func TestMatrix_Multiplication(t *testing.T) {
	identity := NewMatrix()
	m := NewMatrix()
	m.Set(0, 3, 5.0) // Translation

	// Multiply by identity should return original
	result := m.Mul(identity)
	if result.At(0, 3) != 5.0 {
		t.Errorf("Matrix multiplication with identity failed")
	}
}

func TestMatrix_Properties(t *testing.T) {
	identity := NewMatrix()

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
		if !NearlyEqual(det, 1.0) {
			t.Errorf("Identity matrix determinant = %v, want 1.0", det)
		}
	})
}

func TestMatrix_Transformations(t *testing.T) {
	t.Run("Translation", func(t *testing.T) {
		m := NewMatrix()
		translation := NewVector3(1.0, 2.0, 3.0)
		result := m.Translate(translation)

		if result.At(0, 3) != 1.0 || result.At(1, 3) != 2.0 || result.At(2, 3) != 3.0 {
			t.Error("Translation matrix creation failed")
		}
	})

	t.Run("Scale", func(t *testing.T) {
		m := NewMatrix()
		scale := NewVector3(2.0, 3.0, 4.0)
		result := m.Scale(scale)

		if result.At(0, 0) != 2.0 || result.At(1, 1) != 3.0 || result.At(2, 2) != 4.0 {
			t.Error("Scale matrix creation failed")
		}
	})

	t.Run("RotationX", func(t *testing.T) {
		m := NewMatrix()
		angle := Radians(math.Pi / 2) // 90 degrees
		result := m.RotateX(angle)

		// Check rotation matrix properties
		if !NearlyEqual(result.At(0, 0), 1.0) {
			t.Error("X rotation should not affect X axis")
		}
		if !NearlyEqual(result.At(1, 1), 0.0) {
			t.Error("90 degree X rotation Y component should be 0")
		}
	})
}

func TestMatrix_Inverse(t *testing.T) {
	identity := NewMatrix()

	inv := identity.Invert()
	// Invert() returns Matrix, not (Matrix, error)
	// So we check if the result is identity or zero matrix (for non-invertible)
	if !inv.IsIdentity() {
		t.Error("Inverse of identity should be identity")
	}
}

func TestMatrix_Transpose(t *testing.T) {
	m := NewMatrix()
	m.Set(0, 1, 5.0)

	transposed := m.Transpose()
	if transposed.At(1, 0) != 5.0 {
		t.Error("Transpose failed")
	}
}

func TestMatrix_BasisVectors(t *testing.T) {
	identity := NewMatrix()
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

func TestMatrix_Scale(t *testing.T) {
	m := NewMatrix()
	m.Set(0, 0, 2.0)
	m.Set(1, 1, 3.0)
	m.Set(2, 2, 4.0)

	scale := m.GetScale()
	if !NearlyEqual(scale.X(), 2.0) || !NearlyEqual(scale.Y(), 3.0) || !NearlyEqual(scale.Z(), 4.0) {
		t.Errorf("Scale() = (%v, %v, %v), want (2.0, 3.0, 4.0)", scale.X(), scale.Y(), scale.Z())
	}
}

func TestMatrix_TransformPoint(t *testing.T) {
	m := NewMatrix()
	point := NewPoint(1, 2)

	// Identity transform should return same point
	result := point.Transform(m)
	if !NearlyEqual(result.X(), point.X()) || !NearlyEqual(result.Y(), point.Y()) {
		t.Errorf("Identity transform failed")
	}
}

func TestMatrix_TransformVector(t *testing.T) {
	m := NewMatrix()
	v3 := NewVector3(1, 2, 3)

	// Identity transform should return same vector
	result := v3.Transform(m)
	if !NearlyEqual(result.X(), v3.X()) || !NearlyEqual(result.Y(), v3.Y()) || !NearlyEqual(result.Z(), v3.Z()) {
		t.Errorf("Identity vector transform failed")
	}
}

func TestMatrix_Analysis(t *testing.T) {
	identity := NewMatrix()

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

func TestMatrix_Decompose(t *testing.T) {
	m := NewMatrix()
	d := m.Decompose()

	// Check identity decomposition
	if !NearlyEqual(d.Translation.X(), 0.0) ||
		!NearlyEqual(d.Translation.Y(), 0.0) ||
		!NearlyEqual(d.Translation.Z(), 0.0) {
		t.Error("Identity matrix should have zero translation")
	}

	if !NearlyEqual(d.Scale.X(), 1.0) ||
		!NearlyEqual(d.Scale.Y(), 1.0) ||
		!NearlyEqual(d.Scale.Z(), 1.0) {
		t.Error("Identity matrix should have unit scale")
	}
}

func TestMatrix_CosSin(t *testing.T) {
	m := NewMatrix()

	// Test 90 degrees
	cos, sin := m.CosSin(Radians(math.Pi / 2))
	if !NearlyEqual(cos, 0.0) || !NearlyEqual(sin, 1.0) {
		t.Errorf("CosSin(π/2) = (%v, %v), want (0.0, 1.0)", cos, sin)
	}

	// Test 0 degrees
	cos, sin = m.CosSin(Radians(0))
	if !NearlyEqual(cos, 1.0) || !NearlyEqual(sin, 0.0) {
		t.Errorf("CosSin(0) = (%v, %v), want (1.0, 0.0)", cos, sin)
	}
}

func TestMatrix_LookAt(t *testing.T) {
	m := NewMatrix()
	position := NewVector3(0.0, 0.0, 1.0)
	target := NewVector3(0.0, 0.0, 0.0)
	up := NewVector3(0.0, 1.0, 0.0)

	view := m.LookAt(position, target, up)

	// Check that the matrix is not identity
	if view.IsIdentity() {
		t.Error("LookAt matrix should not be identity")
	}
}

func TestMatrix_Orthographic(t *testing.T) {
	m := NewMatrix()
	size := NewSize(800.0, 600.0)

	ortho := m.Orthographic(size)

	// Check that the matrix is not identity
	if ortho.IsIdentity() {
		t.Error("Orthographic matrix should not be identity")
	}
}

func TestMatrix_QuaternionRotation(t *testing.T) {
	m := NewMatrix()
	quat := NewQuaternion(0, 0, 0, 1)

	rotMatrix := m.RotateQuat(quat)

	// Identity quaternion should produce identity rotation matrix
	if !rotMatrix.IsIdentity() {
		t.Error("Identity quaternion should produce identity rotation matrix")
	}
}

func TestMatrix_AxisAngleRotation(t *testing.T) {
	m := NewMatrix()
	axis := NewVector3(0.0, 0.0, 1.0)
	angle := Radians(0.0)

	rotMatrix := m.Rotate(angle, axis)

	// Zero angle should produce identity matrix
	if !rotMatrix.IsIdentity() {
		t.Error("Zero angle rotation should produce identity matrix")
	}
}

func TestMatrix_Equal(t *testing.T) {
	m1 := NewMatrix()
	m2 := NewMatrix()

	if !m1.Equal(m2) {
		t.Error("Two identity matrices should be Eq")
	}

	m2.Set(0, 0, 2.0)
	if m1.Equal(m2) {
		t.Error("Modified matrix should not Eq identity")
	}
}

func TestMatrix_ChainedTransformations(t *testing.T) {
	t.Run("TranslateScaleRotate", func(t *testing.T) {
		m := NewMatrix()

		// Apply transformations in sequence: translate -> scale -> rotate
		translation := NewVector3(1, 2, 3)
		scale := NewVector3(2, 3, 4)

		result := m.Translate(translation).Scale(scale).RotateZ(Radians(math.Pi / 4))

		// Test that the final matrix is not identity
		if result.IsIdentity() {
			t.Error("Chained transformations should not result in identity matrix")
		}

		// Test that translation component is preserved correctly
		if !NearlyEqual(result.At(0, 3), 1.0) {
			t.Errorf("Expected translation X to be 1.0, got %v", result.At(0, 3))
		}
	})

	t.Run("InverseTransformChain", func(t *testing.T) {
		m := NewMatrix()

		// Create a complex transformation
		original := m.Translate(NewVector3(5.0, 10.0, 15.0)).
			Scale(NewVector3(2.0, 3.0, 4.0)).
			RotateY(Radians(math.Pi / 6))

		// Get inverse
		inverse := original.Invert()

		// Multiply original by inverse should give identity
		identity := original.Mul(inverse)
		if !identity.IsIdentity() {
			t.Error("Original * Inverse should Eq identity matrix")
		}
	})

	t.Run("MultipleRotations", func(t *testing.T) {
		m := NewMatrix()

		// Apply multiple rotations
		result := m.RotateX(Radians(math.Pi / 4)).
			RotateY(Radians(math.Pi / 3)).
			RotateZ(Radians(math.Pi / 6))

		// Check that it's still a valid rotation matrix (determinant should be 1)
		det := result.Determinant()
		if !NearlyEqual(det, 1.0) {
			t.Errorf("Rotation matrix determinant should be 1.0, got %v", det)
		}

		// Check that the matrix is orthogonal (transpose == inverse)
		transpose := result.Transpose()
		inverse := result.Invert()
		if !transpose.Equal(inverse) {
			t.Error("Rotation matrix should be orthogonal (transpose == inverse)")
		}
	})
}
