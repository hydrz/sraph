package geom

import (
	"math"
	"testing"
)

func TestNewMatrix(t *testing.T) {
	m := NewMatrix[F32]()

	// Should be identity matrix
	if !m.IsIdentity() {
		t.Error("NewMatrix should create identity matrix")
	}
}

func TestMatrix_At(t *testing.T) {
	m := NewMatrix[F32]()

	// Test identity matrix values
	if m.At(0, 0) != 1.0 || m.At(1, 1) != 1.0 || m.At(2, 2) != 1.0 || m.At(3, 3) != 1.0 {
		t.Error("Identity matrix diagonal should be 1.0")
	}

	if m.At(0, 1) != 0.0 || m.At(1, 0) != 0.0 {
		t.Error("Identity matrix off-diagonal should be 0.0")
	}
}

func TestMatrix_Set(t *testing.T) {
	m := NewMatrix[F32]()
	m.Set(0, 1, F32(5.0))

	if m.At(0, 1) != 5.0 {
		t.Errorf("Set(0, 1, 5.0) failed, got %v", m.At(0, 1))
	}
}

func TestMatrix_Arithmetic(t *testing.T) {
	m1 := NewMatrix[F32]()
	m2 := NewMatrix[F32]()

	// Set some test values
	m1.Set(0, 0, F32(2.0))
	m2.Set(0, 0, F32(3.0))

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
	identity := NewMatrix[F32]()
	m := NewMatrix[F32]()
	m.Set(0, 3, F32(5.0)) // Translation

	// Multiply by identity should return original
	result := m.Mul(identity)
	if result.At(0, 3) != 5.0 {
		t.Errorf("Matrix multiplication with identity failed")
	}
}

func TestMatrix_Properties(t *testing.T) {
	identity := NewMatrix[F32]()

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
		if !Equal(det, F32(1.0)) {
			t.Errorf("Identity matrix determinant = %v, want 1.0", det)
		}
	})
}

func TestMatrix_Transformations(t *testing.T) {
	t.Run("Translation", func(t *testing.T) {
		m := NewMatrix[F32]()
		translation := NewVector3[F32](1.0, 2.0, 3.0)
		result := m.Translate(translation)

		if result.At(0, 3) != 1.0 || result.At(1, 3) != 2.0 || result.At(2, 3) != 3.0 {
			t.Error("Translation matrix creation failed")
		}
	})

	t.Run("Scale", func(t *testing.T) {
		m := NewMatrix[F32]()
		scale := NewVector3[F32](2.0, 3.0, 4.0)
		result := m.Scale(scale)

		if result.At(0, 0) != 2.0 || result.At(1, 1) != 3.0 || result.At(2, 2) != 4.0 {
			t.Error("Scale matrix creation failed")
		}
	})

	t.Run("RotationX", func(t *testing.T) {
		m := NewMatrix[F32]()
		angle := NewRadians[F32](math.Pi / 2) // 90 degrees
		result := m.RotateX(angle)

		// Check rotation matrix properties
		if !Equal(result.At(0, 0), F32(1.0)) {
			t.Error("X rotation should not affect X axis")
		}
		if !Equal(result.At(1, 1), F32(0.0)) {
			t.Error("90 degree X rotation Y component should be 0")
		}
	})
}

func TestMatrix_Inverse(t *testing.T) {
	identity := NewMatrix[F32]()

	inv, err := identity.Inverse()
	if err != nil {
		t.Errorf("Identity matrix should be invertible: %v", err)
	}

	if !inv.IsIdentity() {
		t.Error("Inverse of identity should be identity")
	}
}

func TestMatrix_Transpose(t *testing.T) {
	m := NewMatrix[F32]()
	m.Set(0, 1, F32(5.0))

	transposed := m.Transpose()
	if transposed.At(1, 0) != 5.0 {
		t.Error("Transpose failed")
	}
}

func TestMatrix_BasisVectors(t *testing.T) {
	identity := NewMatrix[F32]()
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
	m := NewMatrix[F32]()
	m.Set(0, 0, F32(2.0))
	m.Set(1, 1, F32(3.0))
	m.Set(2, 2, F32(4.0))

	scale := m.GetScale()
	if !Equal(scale.X(), F32(2.0)) || !Equal(scale.Y(), F32(3.0)) || !Equal(scale.Z(), F32(4.0)) {
		t.Errorf("Scale() = (%v, %v, %v), want (2.0, 3.0, 4.0)", scale.X(), scale.Y(), scale.Z())
	}
}

func TestMatrix_TransformPoint(t *testing.T) {
	m := NewMatrix[F32]()
	point := NewPoint[F32](1, 2)

	// Identity transform should return same point
	result := m.TransformPoint(point)
	if !Equal(result.X(), point.X()) || !Equal(result.Y(), point.Y()) {
		t.Errorf("Identity transform failed")
	}
}

func TestMatrix_TransformVector(t *testing.T) {
	m := NewMatrix[F32]()
	v3 := NewVector3[F32](1, 2, 3)

	// Identity transform should return same vector
	result := m.TransformVector3D(v3)
	if !Equal(result.X(), v3.X()) || !Equal(result.Y(), v3.Y()) || !Equal(result.Z(), v3.Z()) {
		t.Errorf("Identity vector transform failed")
	}
}

func TestMatrix_Analysis(t *testing.T) {
	identity := NewMatrix[F32]()

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
	m := NewMatrix[F32]()
	d := m.Decompose()

	// Check identity decomposition
	if !Equal(d.Translation.X(), F32(0.0)) ||
		!Equal(d.Translation.Y(), F32(0.0)) ||
		!Equal(d.Translation.Z(), F32(0.0)) {
		t.Error("Identity matrix should have zero translation")
	}

	if !Equal(d.Scale.X(), F32(1.0)) ||
		!Equal(d.Scale.Y(), F32(1.0)) ||
		!Equal(d.Scale.Z(), F32(1.0)) {
		t.Error("Identity matrix should have unit scale")
	}
}

func TestMatrix_CosSin(t *testing.T) {
	m := NewMatrix[F32]()

	// Test 90 degrees
	cos, sin := m.CosSin(NewRadians[F32](math.Pi / 2))
	if !Equal(cos, F32(0.0)) || !Equal(sin, F32(1.0)) {
		t.Errorf("CosSin(π/2) = (%v, %v), want (0.0, 1.0)", cos, sin)
	}

	// Test 0 degrees
	cos, sin = m.CosSin(NewRadians[F32](0))
	if !Equal(cos, F32(1.0)) || !Equal(sin, F32(0.0)) {
		t.Errorf("CosSin(0) = (%v, %v), want (1.0, 0.0)", cos, sin)
	}
}

func TestMatrix_LookAt(t *testing.T) {
	m := NewMatrix[F32]()
	position := NewVector3[F32](0.0, 0.0, 1.0)
	target := NewVector3[F32](0.0, 0.0, 0.0)
	up := NewVector3[F32](0.0, 1.0, 0.0)

	view := m.LookAt(position, target, up)

	// Check that the matrix is not identity
	if view.IsIdentity() {
		t.Error("LookAt matrix should not be identity")
	}
}

func TestMatrix_Orthographic(t *testing.T) {
	m := NewMatrix[F32]()
	size := NewSize[F32](800.0, 600.0)

	ortho := m.Orthographic(size)

	// Check that the matrix is not identity
	if ortho.IsIdentity() {
		t.Error("Orthographic matrix should not be identity")
	}
}

func TestMatrix_QuaternionRotation(t *testing.T) {
	m := NewMatrix[F32]()
	quat := NewQuaternion[F32](0, 0, 0, 1)

	rotMatrix := m.RotateQuaternion(quat)

	// Identity quaternion should produce identity rotation matrix
	if !rotMatrix.IsIdentity() {
		t.Error("Identity quaternion should produce identity rotation matrix")
	}
}

func TestMatrix_AxisAngleRotation(t *testing.T) {
	m := NewMatrix[F32]()
	axis := NewVector3[F32](0.0, 0.0, 1.0)
	angle := NewRadians[F32](0.0)

	rotMatrix := m.RotateAxisAngle(angle, axis)

	// Zero angle should produce identity matrix
	if !rotMatrix.IsIdentity() {
		t.Error("Zero angle rotation should produce identity matrix")
	}
}

func TestMatrix_Equal(t *testing.T) {
	m1 := NewMatrix[F32]()
	m2 := NewMatrix[F32]()

	if !m1.Equal(m2) {
		t.Error("Two identity matrices should be equal")
	}

	m2.Set(0, 0, F32(2.0))
	if m1.Equal(m2) {
		t.Error("Modified matrix should not equal identity")
	}
}

func TestMatrix_ChainedTransformations(t *testing.T) {
	t.Run("TranslateScaleRotate", func(t *testing.T) {
		m := NewMatrix[F32]()

		// Apply transformations in sequence: translate -> scale -> rotate
		translation := NewVector3[F32](1, 2, 3)
		scale := NewVector3[F32](2, 3, 4)

		result := m.Translate(translation).Scale(scale).RotateZ(NewRadians[F32](math.Pi / 4))

		// Test that the final matrix is not identity
		if result.IsIdentity() {
			t.Error("Chained transformations should not result in identity matrix")
		}

		// Test that translation component is preserved correctly
		if !Equal(result.At(0, 3), F32(1.0)) {
			t.Errorf("Expected translation X to be 1.0, got %v", result.At(0, 3))
		}
	})

	t.Run("InverseTransformChain", func(t *testing.T) {
		m := NewMatrix[F32]()

		// Create a complex transformation
		original := m.Translate(NewVector3[F32](5.0, 10.0, 15.0)).
			Scale(NewVector3[F32](2.0, 3.0, 4.0)).
			RotateY(NewRadians[F32](math.Pi / 6))

		// Get inverse
		inverse, err := original.Inverse()
		if err != nil {
			t.Fatalf("Failed to compute inverse: %v", err)
		}

		// Multiply original by inverse should give identity
		identity := original.Mul(inverse)
		if !identity.IsIdentity() {
			t.Error("Original * Inverse should equal identity matrix")
		}
	})

	t.Run("MultipleRotations", func(t *testing.T) {
		m := NewMatrix[F32]()

		// Apply multiple rotations
		result := m.RotateX(NewRadians[F32](math.Pi / 4)).
			RotateY(NewRadians[F32](math.Pi / 3)).
			RotateZ(NewRadians[F32](math.Pi / 6))

		// Check that it's still a valid rotation matrix (determinant should be 1)
		det := result.Determinant()
		if !Equal(det, F32(1.0)) {
			t.Errorf("Rotation matrix determinant should be 1.0, got %v", det)
		}

		// Check that the matrix is orthogonal (transpose == inverse)
		transpose := result.Transpose()
		inverse, _ := result.Inverse()
		if !transpose.Equal(inverse) {
			t.Error("Rotation matrix should be orthogonal (transpose == inverse)")
		}
	})
}

func TestMatrix_BoundaryConditions(t *testing.T) {
	t.Run("ZeroScale", func(t *testing.T) {
		m := NewMatrix[F32]()
		zeroScale := NewVector3[F32](0.0, 0.0, 0.0)

		result := m.Scale(zeroScale)

		// Zero scale should make matrix non-invertible
		if result.IsInvertible() {
			t.Error("Zero scale matrix should not be invertible")
		}

		// Determinant should be 0
		if !Equal(result.Determinant(), F32(0.0)) {
			t.Error("Zero scale matrix determinant should be 0")
		}
	})

	t.Run("VeryLargeValues", func(t *testing.T) {
		m := NewMatrix[F32]()
		largeScale := NewVector3[F32](1e6, 1e6, 1e6)

		result := m.Scale(largeScale)

		// Should still be finite
		if !result.IsFinite() {
			t.Error("Large scale matrix should still be finite")
		}

		// Should still be invertible
		if !result.IsInvertible() {
			t.Error("Large scale matrix should still be invertible")
		}
	})

	t.Run("VerySmallValues", func(t *testing.T) {
		m := NewMatrix[F32]()
		smallScale := NewVector3[F32](1e-6, 1e-6, 1e-6)

		result := m.Scale(smallScale)

		// Should still be finite and invertible
		if !result.IsFinite() || !result.IsInvertible() {
			t.Error("Small scale matrix should be finite and invertible")
		}
	})

	t.Run("PiRotations", func(t *testing.T) {
		m := NewMatrix[F32]()

		// Test common angles
		angles := []F32{0, math.Pi / 6, math.Pi / 4, math.Pi / 3, math.Pi / 2, math.Pi, 2 * math.Pi}

		for _, angle := range angles {
			result := m.RotateZ(NewRadians[F32](angle))

			if !result.IsFinite() {
				t.Errorf("Rotation by %v radians should be finite", angle)
			}

			// Check determinant is 1 (within tolerance)
			det := result.Determinant()
			if math.Abs(float64(det)-1.0) > 1e-6 {
				t.Errorf("Rotation determinant should be 1.0, got %v for angle %v", det, angle)
			}
		}
	})
}

func TestMatrix_NumericalStability(t *testing.T) {
	t.Run("RepeatedOperations", func(t *testing.T) {
		m := NewMatrix[F32]()
		result := m

		// Apply small rotation 1000 times
		smallAngle := NewRadians[F32](0.001) // Very small angle
		for i := 0; i < 1000; i++ {
			result = result.RotateZ(smallAngle)
		}

		// Should be equivalent to single rotation of 1.0 radian
		expected := m.RotateZ(NewRadians[F32](1.0))

		// Check if results are approximately equal (allowing for numerical error)
		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				diff := math.Abs(float64(result.At(row, col) - expected.At(row, col)))
				if diff > 0.01 { // Allow for some numerical error
					t.Errorf("Repeated small rotations differ from single rotation at [%d][%d]: diff=%v", row, col, diff)
				}
			}
		}
	})

	t.Run("NearSingularMatrix", func(t *testing.T) {
		m := NewMatrix[F32]()

		// Create nearly singular matrix (very small scale in one dimension)
		nearSingular := m.Scale(NewVector3[F32](1.0, 1e-10, 1.0))

		// Should still be technically invertible
		if !nearSingular.IsInvertible() {
			t.Error("Nearly singular matrix should still be invertible")
		}

		// But inverse might be numerically unstable
		inverse, err := nearSingular.Inverse()
		if err != nil {
			t.Errorf("Failed to invert nearly singular matrix: %v", err)
		}

		// Check if inverse operation is stable
		if !inverse.IsFinite() {
			t.Error("Inverse of nearly singular matrix should be finite")
		}
	})
}

func TestMatrix_RealWorldScenarios(t *testing.T) {
	t.Run("CameraTransform", func(t *testing.T) {
		m := NewMatrix[F32]()

		// Simulate camera positioning
		position := NewVector3[F32](10.0, 5.0, 15.0)
		target := NewVector3[F32](0.0, 0.0, 0.0)
		up := NewVector3[F32](0.0, 1.0, 0.0)

		viewMatrix := m.LookAt(position, target, up)

		// View matrix should not be identity
		if viewMatrix.IsIdentity() {
			t.Error("LookAt matrix should not be identity")
		}

		// Should be invertible (for world-to-view conversion)
		if !viewMatrix.IsInvertible() {
			t.Error("View matrix should be invertible")
		}

		// Create projection matrix
		size := NewSize[F32](800.0, 600.0)
		projMatrix := m.Orthographic(size)

		// Combine view and projection
		mvp := projMatrix.Mul(viewMatrix)

		if !mvp.IsFinite() {
			t.Error("Model-View-Projection matrix should be finite")
		}
	})

	t.Run("ObjectHierarchy", func(t *testing.T) {
		// Simulate parent-child object hierarchy
		parentTransform := NewMatrix[F32]().
			Translate(NewVector3[F32](5, 0, 0)).
			RotateY(NewRadians[F32](math.Pi / 4))

		childLocalTransform := NewMatrix[F32]().
			Translate(NewVector3[F32](2, 1, 0)).
			Scale(NewVector3[F32](0.5, 0.5, 0.5))

		// Child world transform = parent * child_local
		childWorldTransform := parentTransform.Mul(childLocalTransform)

		// Test that child inherits parent's transformation
		if childWorldTransform.IsIdentity() {
			t.Error("Child world transform should not be identity")
		}

		// Child should have combined scale
		childScale := childWorldTransform.GetScale()
		if !Equal(childScale.X(), F32(0.5)) {
			t.Errorf("Child should inherit scale, got %v", childScale.X())
		}
	})

	t.Run("AnimationInterpolation", func(t *testing.T) {
		// Test matrix interpolation scenario
		start := NewMatrix[F32]().
			Translate(NewVector3[F32](0, 0, 0))

		end := NewMatrix[F32]().
			Translate(NewVector3[F32](10, 5, 0)).
			RotateZ(NewRadians[F32](math.Pi))

		// Simple linear interpolation of individual elements (not ideal, but for testing)
		tt := F32(0.5) // 50% interpolation

		var lerped matrix[F32]
		startMatrix := start.(*matrix[F32])
		endMatrix := end.(*matrix[F32])

		for i := 0; i < 16; i++ {
			lerped[i] = startMatrix[i]*(1-tt) + endMatrix[i]*tt
		}

		// Interpolated matrix should be finite
		if !lerped.IsFinite() {
			t.Error("Interpolated matrix should be finite")
		}

		// Should be somewhere between start and end positions
		translation := NewVector3(lerped[12], lerped[13], lerped[14])
		if translation.X() < 4.0 || translation.X() > 6.0 {
			t.Errorf("Interpolated translation X should be around 5.0, got %v", translation.X())
		}
	})
}

func TestMatrix_ErrorHandling(t *testing.T) {
	t.Run("SingularMatrixInverse", func(t *testing.T) {
		m := NewMatrix[F32]()

		// Create singular matrix (determinant = 0)
		singular := m.Scale(NewVector3[F32](1.0, 0.0, 1.0))

		_, err := singular.Inverse()
		if err == nil {
			t.Error("Expected error when inverting singular matrix")
		}
	})
}

func TestMatrix_SpecialCases(t *testing.T) {
	t.Run("IdentityProperties", func(t *testing.T) {
		identity := NewMatrix[F32]()

		// Identity * Identity = Identity
		result := identity.Mul(identity)
		if !result.IsIdentity() {
			t.Error("Identity * Identity should equal Identity")
		}

		// Identity + Identity should have diagonal elements = 2
		sum := identity.Add(identity)
		if !Equal(sum.At(0, 0), F32(2.0)) {
			t.Error("Identity + Identity should have diagonal elements = 2")
		}
	})

	t.Run("TransformationPreservation", func(t *testing.T) {
		// Test that transformations preserve certain properties
		m := NewMatrix[F32]()

		// Uniform scaling should preserve angles
		uniform := m.Scale(NewVector3[F32](2.0, 2.0, 2.0))

		v1 := NewVector3[F32](1.0, 0.0, 0.0)
		v2 := NewVector3[F32](0.0, 1.0, 0.0)

		transformed_v1 := uniform.TransformVector3D(v1)
		transformed_v2 := uniform.TransformVector3D(v2)

		// Angle between transformed vectors should be same as original (90 degrees)
		originalDot := v1.Dot(v2) // Should be 0
		transformedDot := transformed_v1.Normalize().Dot(transformed_v2.Normalize())

		if !Equal(originalDot, transformedDot) {
			t.Error("Uniform scaling should preserve angles")
		}
	})

	t.Run("AxisAlignedDetection", func(t *testing.T) {
		m := NewMatrix[F32]()

		// Pure translation should be axis-aligned
		translated := m.Translate(NewVector3[F32](5.0, 3.0, 1.0))
		if !translated.IsAxisAligned() {
			t.Error("Pure translation should be axis-aligned")
		}

		// Pure scale should be axis-aligned
		scaled := m.Scale(NewVector3[F32](2.0, 3.0, 0.5))
		if !scaled.IsAxisAligned() {
			t.Error("Pure scale should be axis-aligned")
		}

		// Rotation should not be axis-aligned (except for multiples of 90 degrees)
		rotated := m.RotateZ(NewRadians[F32](math.Pi / 6)) // 30 degrees
		if rotated.IsAxisAligned() {
			t.Error("30 degree rotation should not be axis-aligned")
		}

		// 90 degree rotation should be axis-aligned
		rotated90 := m.RotateZ(NewRadians[F32](math.Pi / 2))
		if !rotated90.IsAxisAligned() {
			t.Error("90 degree rotation should be axis-aligned")
		}
	})
}

func BenchmarkMatrixOperations(b *testing.B) {
	m1 := NewMatrix[F32]()
	m2 := NewMatrix[F32]()

	// Add some values to make operations more meaningful
	m1.Set(0, 3, F32(1.0))
	m2.Set(1, 3, F32(2.0))

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

	// Benchmark
	b.Run("ChainedTransforms", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m1.Translate(NewVector3[F32](1.0, 2.0, 3.0)).
				Scale(NewVector3[F32](2.0, 2.0, 2.0)).
				RotateZ(NewRadians[F32](math.Pi / 4))
		}
	})

	b.Run("VectorTransform", func(b *testing.B) {
		v := NewVector3[F32](1.0, 2.0, 3.0)
		for i := 0; i < b.N; i++ {
			m1.TransformVector3D(v)
		}
	})

	b.Run("MatrixDecomposition", func(b *testing.B) {
		complex := m1.Translate(NewVector3[F32](5.0, 3.0, 1.0)).
			Scale(NewVector3[F32](2.0, 1.5, 0.8)).
			RotateY(NewRadians[F32](math.Pi / 3))

		for i := 0; i < b.N; i++ {
			complex.Decompose()
		}
	})

	b.Run("SingleMatrixMul", func(b *testing.B) {
		base := NewMatrix[F32]()
		transform := base.Translate(NewVector3[F32](1.0, 2.0, 3.0)).
			Scale(NewVector3[F32](2.0, 2.0, 2.0))

		for i := 0; i < b.N; i++ {
			base.Mul(transform)
		}
	})

	b.Run("ChainedOperations", func(b *testing.B) {
		base := NewMatrix[F32]()

		for i := 0; i < b.N; i++ {
			base.Translate(NewVector3[F32](1.0, 2.0, 3.0)).
				Scale(NewVector3[F32](2.0, 2.0, 2.0))
		}
	})
}

func BenchmarkMatrixComparison(b *testing.B) {
	// Compare different ways of achieving the same transformation

	b.Run("SingleMatrixMul", func(b *testing.B) {
		base := NewMatrix[F32]()
		transform := base.Translate(NewVector3[F32](1.0, 2.0, 3.0)).
			Scale(NewVector3[F32](2.0, 2.0, 2.0))

		for i := 0; i < b.N; i++ {
			base.Mul(transform)
		}
	})

	b.Run("ChainedOperations", func(b *testing.B) {
		base := NewMatrix[F32]()

		for i := 0; i < b.N; i++ {
			base.Translate(NewVector3[F32](1.0, 2.0, 3.0)).
				Scale(NewVector3[F32](2.0, 2.0, 2.0))
		}
	})
}
