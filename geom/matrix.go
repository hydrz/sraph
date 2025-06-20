package geom

import (
	"fmt"
	"math"
)

// Matrix represents a 4x4 matrix using column-major storage.
//
// All methods that use normalized device coordinates (NDC) with Matrix
// follow the DirectX/Vulkan convention:
//
//   - Left-handed coordinate system:
//     X increases to the right, Y increases upward, Z increases away from the viewer.
//     Positive rotation is clockwise about the rotation axis.
//   - NDC bounds:
//     Lower-left corner:  (-1.0, -1.0)
//     Upper-right corner: ( 1.0,  1.0)
//     Z range (visible):   0.0 (near) to 1.0 (far)
//   - NDC origin is at (0.0, 0.0, 0.5) (center of the NDC cube).
//
// NOTE: This differs from OpenGL, which uses a right-handed system and NDC z in [-1, 1].
//
// For more details, see:
//
//	https://learn.microsoft.com/zh-cn/windows/win32/learnwin32/appendix--matrix-transforms
//	https://www.opengl-tutorial.org/beginners-tutorials/tutorial-3-matrices/
type Matrix[T TScalar] [16]T

// NewMatrix creates a new identity matrix.
func NewMatrix[T TScalar]() Matrix[T] {
	return Matrix[T]{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// At returns the value at the specified row and column (0-based).
func (m Matrix[T]) At(row int, col int) T {
	// Column-major order: index = col*4 + row
	return m[col*4+row]
}

// Set sets the value at the specified row and column (0-based).
func (m *Matrix[T]) Set(row int, col int, value T) {
	m[col*4+row] = value
}

// Add returns the sum of this matrix and another matrix (element-wise addition).
func (m Matrix[T]) Add(other Matrix[T]) Matrix[T] {
	o := other
	return Matrix[T]{
		m[0] + o[0], m[1] + o[1], m[2] + o[2], m[3] + o[3],
		m[4] + o[4], m[5] + o[5], m[6] + o[6], m[7] + o[7],
		m[8] + o[8], m[9] + o[9], m[10] + o[10], m[11] + o[11],
		m[12] + o[12], m[13] + o[13], m[14] + o[14], m[15] + o[15],
	}
}

// Sub returns the difference of this matrix and another matrix (element-wise subtraction).
func (m Matrix[T]) Sub(other Matrix[T]) Matrix[T] {
	o := other
	return Matrix[T]{
		m[0] - o[0], m[1] - o[1], m[2] - o[2], m[3] - o[3],
		m[4] - o[4], m[5] - o[5], m[6] - o[6], m[7] - o[7],
		m[8] - o[8], m[9] - o[9], m[10] - o[10], m[11] - o[11],
		m[12] - o[12], m[13] - o[13], m[14] - o[14], m[15] - o[15],
	}
}

// Mul returns the product of this matrix and another matrix (matrix multiplication).
func (m Matrix[T]) Mul(o Matrix[T]) Matrix[T] {
	return Matrix[T]{
		m[0]*o[0] + m[4]*o[1] + m[8]*o[2] + m[12]*o[3],
		m[1]*o[0] + m[5]*o[1] + m[9]*o[2] + m[13]*o[3],
		m[2]*o[0] + m[6]*o[1] + m[10]*o[2] + m[14]*o[3],
		m[3]*o[0] + m[7]*o[1] + m[11]*o[2] + m[15]*o[3],
		m[0]*o[4] + m[4]*o[5] + m[8]*o[6] + m[12]*o[7],
		m[1]*o[4] + m[5]*o[5] + m[9]*o[6] + m[13]*o[7],
		m[2]*o[4] + m[6]*o[5] + m[10]*o[6] + m[14]*o[7],
		m[3]*o[4] + m[7]*o[5] + m[11]*o[6] + m[15]*o[7],
		m[0]*o[8] + m[4]*o[9] + m[8]*o[10] + m[12]*o[11],
		m[1]*o[8] + m[5]*o[9] + m[9]*o[10] + m[13]*o[11],
		m[2]*o[8] + m[6]*o[9] + m[10]*o[10] + m[14]*o[11],
		m[3]*o[8] + m[7]*o[9] + m[11]*o[10] + m[15]*o[11],
		m[0]*o[12] + m[4]*o[13] + m[8]*o[14] + m[12]*o[15],
		m[1]*o[12] + m[5]*o[13] + m[9]*o[14] + m[13]*o[15],
		m[2]*o[12] + m[6]*o[13] + m[10]*o[14] + m[14]*o[15],
		m[3]*o[12] + m[7]*o[13] + m[11]*o[14] + m[15]*o[15],
	}
}

// Eq checks if this matrix is Eq to another matrix (all elements Eq).
// Uses exact floating-point comparison, which may not be suitable for computed matrices.
// Consider using approximate comparison for floating-point matrices.
func (m Matrix[T]) Eq(other Matrix[T]) bool {
	o := other
	for i := 0; i < 16; i++ {
		if !NearlyEqual(m[i], o[i]) {
			return false
		}
	}
	return true
}

// IsFinite returns true if all matrix comp1nts are finite.
// Checks for NaN (Not a Number) and infinite values in all matrix elements.
// This is important for numerical stability and error detection.
func (m Matrix[T]) IsFinite() bool {
	for i := 0; i < 16; i++ {
		if !IsFinite(m[i]) {
			return false
		}
	}
	return true
}

// IsAffine returns true if the matrix is an affine transformation matrix (last row is [0 0 0 1]).
// Affine transformations preserve parallel lines and ratios of distances along lines.
// They include translation, rotation, scaling, and shearing, but not perspective projection.
func (m Matrix[T]) IsAffine() bool {
	return NearlyEqual(m[2], 0) && NearlyEqual(m[3], 0) && NearlyEqual(m[6], 0) && NearlyEqual(m[7], 0) &&
		NearlyEqual(m[8], 0) && NearlyEqual(m[9], 0) && NearlyEqual(m[10], 1) && NearlyEqual(m[11], 0) &&
		NearlyEqual(m[14], 0) && NearlyEqual(m[15], 1)
}

// IsIdentity returns true if the matrix is an identity matrix.
// The identity matrix has no effect when applied as a transformation:
//   - Translations remain unchanged.
//   - Rotations are 0.
//   - Scales are uniform with factor 1.
//   - Shears are 0.
func (m Matrix[T]) IsIdentity() bool {

	return NearlyEqual(m[0], 1) && NearlyEqual(m[1], 0) && NearlyEqual(m[2], 0) && NearlyEqual(m[3], 0) &&
		NearlyEqual(m[4], 0) && NearlyEqual(m[5], 1) && NearlyEqual(m[6], 0) && NearlyEqual(m[7], 0) &&
		NearlyEqual(m[8], 0) && NearlyEqual(m[9], 0) && NearlyEqual(m[10], 1) && NearlyEqual(m[11], 0) &&
		NearlyEqual(m[12], 0) && NearlyEqual(m[13], 0) && NearlyEqual(m[14], 0) && NearlyEqual(m[15], 1)
}

// IsInvertible returns true if the matrix is invertible (determinant != 0).
// A matrix is invertible if and only if its determinant is non-0.
// Non-invertible matrices are also called singular or degenerate matrices.
func (m Matrix[T]) IsInvertible() bool { return m.Determinant() != 0 }

// Determinant returns the determinant of the matrix.
// The determinant is a scalar value that represents the scaling factor of the transformation.
// It also indicates whether the transformation preserves orientation (det > 0) or reverses it (det < 0).
func (m Matrix[T]) Determinant() T {
	// Using the same algorithm as C++ implementation
	a00, a01, a02, a03 := m[0], m[1], m[2], m[3]
	a10, a11, a12, a13 := m[4], m[5], m[6], m[7]
	a20, a21, a22, a23 := m[8], m[9], m[10], m[11]
	a30, a31, a32, a33 := m[12], m[13], m[14], m[15]

	b00 := a00*a11 - a01*a10
	b01 := a00*a12 - a02*a10
	b02 := a00*a13 - a03*a10
	b03 := a01*a12 - a02*a11
	b04 := a01*a13 - a03*a11
	b05 := a02*a13 - a03*a12
	b06 := a20*a31 - a21*a30
	b07 := a20*a32 - a22*a30
	b08 := a20*a33 - a23*a30
	b09 := a21*a32 - a22*a31
	b10 := a21*a33 - a23*a31
	b11 := a22*a33 - a23*a32

	return b00*b11 - b01*b10 + b02*b09 + b03*b08 - b04*b07 + b05*b06
}

// HasPerspective returns true if the matrix contains a 3D perspective comp1nt (non0 last row except [0 0 0 1]).
// Perspective transformations cause parallel lines to converge at vanishing points.
// The bottom row [a b c d] where a≠0 or b≠0 or c≠0 or d≠1 indicates perspective.
func (m Matrix[T]) HasPerspective() bool {
	return !NearlyEqual(m[3], 0) || !NearlyEqual(m[7], 0) || !NearlyEqual(m[11], 0) || !NearlyEqual(m[15], 1)
}

// HasPerspective2D returns true if the matrix contains a 2D perspective comp1nt.
// 2D perspective affects the homogeneous coordinate in 2D transformations.
// This is less common than 3D perspective but used in some 2D effects.
func (m Matrix[T]) HasPerspective2D() bool {
	return !NearlyEqual(m[3], 0) || !NearlyEqual(m[7], 0) || !NearlyEqual(m[15], 1)
}

// HasTranslation returns true if the matrix contains a translation comp1nt (non0 last column except [0 0 0 1]).
// Translation moves points by adding a vector to their coordinates.
// If a matrix has translation, it cannot be purely rotational or uniform scaling.
func (m Matrix[T]) HasTranslation() bool {
	return !NearlyEqual(m[12], 0) || !NearlyEqual(m[13], 0)
}

// IsAxisAligned returns true if the matrix is axis-aligned (no rotation or shear).
// An axis-aligned matrix has its basis vectors aligned with the coordinate axes.
// This is important for culling, bounding volume computation, and some optimizations.
func (m Matrix[T]) IsAxisAligned() bool {
	if m.HasPerspective() {
		return false
	}

	// Check if all three basis vectors are aligned to an axis
	v := [9]bool{
		!NearlyEqual(m[0], 0), !NearlyEqual(m[1], 0), !NearlyEqual(m[2], 0),
		!NearlyEqual(m[4], 0), !NearlyEqual(m[5], 0), !NearlyEqual(m[6], 0),
		!NearlyEqual(m[8], 0), !NearlyEqual(m[9], 0), !NearlyEqual(m[10], 0),
	}

	bti := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}

	// Check if all three basis vectors are aligned to an axis
	if (bti(v[0])+bti(v[1])+bti(v[2]) != 1) ||
		(bti(v[3])+bti(v[4])+bti(v[5]) != 1) ||
		(bti(v[6])+bti(v[7])+bti(v[8]) != 1) {
		return false
	}

	// Ensure that n1 of the basis vectors overlap
	if (bti(v[0])+bti(v[3])+bti(v[6]) != 1) ||
		(bti(v[1])+bti(v[4])+bti(v[7]) != 1) ||
		(bti(v[2])+bti(v[5])+bti(v[8]) != 1) {
		return false
	}

	return true
}

// IsAxisAligned2D returns true if the matrix is axis-aligned in 2D.
// Similar to IsAxisAligned, but only checks the X and Y axes.
func (m Matrix[T]) IsAxisAligned2D() bool {
	if m.HasPerspective2D() {
		return false
	}

	if NearlyEqual(m[1], 0) && NearlyEqual(m[4], 0) {
		return true
	}
	if NearlyEqual(m[0], 0) && NearlyEqual(m[5], 0) {
		return true
	}
	return false
}

// IsTranslationOnly returns true if the matrix contains only translation (no scale, rotation, shear, or perspective).
func (m Matrix[T]) IsTranslationOnly() bool {
	// This checks if the matrix is of the form:
	// [ 1 0 0 tx ]
	// [ 0 1 0 ty ]
	// [ 0 0 1 tz ]
	// [ 0 0 0 1 ]

	return NearlyEqual(m[0], 1) && NearlyEqual(m[1], 0) && NearlyEqual(m[2], 0) && NearlyEqual(m[3], 0) &&
		NearlyEqual(m[4], 0) && NearlyEqual(m[5], 1) && NearlyEqual(m[6], 0) && NearlyEqual(m[7], 0) &&
		NearlyEqual(m[8], 0) && NearlyEqual(m[9], 0) && NearlyEqual(m[10], 1) && NearlyEqual(m[11], 0) &&
		NearlyEqual(m[15], 1)
}

// IsTranslationScaleOnly returns true if the matrix contains only translation and scale.
func (m Matrix[T]) IsTranslationScaleOnly() bool {
	// This checks if the matrix is of the form:
	// [ sx 0 0 tx ]
	// [ 0 sy 0 ty ]
	// [ 0 0 sz tz ]
	// [ 0 0 0 1 ]
	// where sx, sy, sz are scale factors.

	return !NearlyEqual(m[0], 0) && NearlyEqual(m[1], 0) && NearlyEqual(m[2], 0) && NearlyEqual(m[3], 0) &&
		NearlyEqual(m[4], 0) && !NearlyEqual(m[5], 0) && NearlyEqual(m[6], 0) && NearlyEqual(m[7], 0) &&
		NearlyEqual(m[8], 0) && NearlyEqual(m[9], 0) && !NearlyEqual(m[10], 0) && NearlyEqual(m[11], 0) &&
		NearlyEqual(m[15], 1)
}

// Transpose returns the transpose of the matrix.
func (m Matrix[T]) Transpose() Matrix[T] {
	return Matrix[T]{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
}

// Invert returns the inverse of the matrix, or an error if not invertible.
// The inverse matrix undoes the transformation of the original matrix:
//   - Translations are negated.
//   - Rotations are inverted.
//   - Scales are inverted (1/sx, 1/sy, 1/sz).
//   - Shears are inverted.
func (m Matrix[T]) Invert() Matrix[T] {
	tmp := Matrix[T]{
		m[5]*m[10]*m[15] - m[5]*m[11]*m[14] - m[9]*m[6]*m[15] + m[9]*m[7]*m[14] + m[13]*m[6]*m[11] - m[13]*m[7]*m[10],
		-m[1]*m[10]*m[15] + m[1]*m[11]*m[14] + m[9]*m[2]*m[15] - m[9]*m[3]*m[14] - m[13]*m[2]*m[11] + m[13]*m[3]*m[10],
		m[1]*m[6]*m[15] - m[1]*m[7]*m[14] - m[5]*m[2]*m[15] + m[5]*m[3]*m[14] + m[13]*m[2]*m[7] - m[13]*m[3]*m[6],
		-m[1]*m[6]*m[11] + m[1]*m[7]*m[10] + m[5]*m[2]*m[11] - m[5]*m[3]*m[10] - m[9]*m[2]*m[7] + m[9]*m[3]*m[6],
		-m[4]*m[10]*m[15] + m[4]*m[11]*m[14] + m[8]*m[6]*m[15] - m[8]*m[7]*m[14] - m[12]*m[6]*m[11] + m[12]*m[7]*m[10],
		m[0]*m[10]*m[15] - m[0]*m[11]*m[14] - m[8]*m[2]*m[15] + m[8]*m[3]*m[14] + m[12]*m[2]*m[11] - m[12]*m[3]*m[10],
		-m[0]*m[6]*m[15] + m[0]*m[7]*m[14] + m[4]*m[2]*m[15] - m[4]*m[3]*m[14] - m[12]*m[2]*m[7] + m[12]*m[3]*m[6],
		m[0]*m[6]*m[11] - m[0]*m[7]*m[10] - m[4]*m[2]*m[11] + m[4]*m[3]*m[10] + m[8]*m[2]*m[7] - m[8]*m[3]*m[6],
		m[4]*m[9]*m[15] - m[4]*m[11]*m[13] - m[8]*m[5]*m[15] + m[8]*m[7]*m[13] + m[12]*m[5]*m[11] - m[12]*m[7]*m[9],
		-m[0]*m[9]*m[15] + m[0]*m[11]*m[13] + m[8]*m[1]*m[15] - m[8]*m[3]*m[13] - m[12]*m[1]*m[11] + m[12]*m[3]*m[9],
		m[0]*m[5]*m[15] - m[0]*m[7]*m[13] - m[4]*m[1]*m[15] + m[4]*m[3]*m[13] + m[12]*m[1]*m[7] - m[12]*m[3]*m[5],
		-m[0]*m[5]*m[11] + m[0]*m[7]*m[9] + m[4]*m[1]*m[11] - m[4]*m[3]*m[9] - m[8]*m[1]*m[7] + m[8]*m[3]*m[5],
		-m[4]*m[9]*m[14] + m[4]*m[10]*m[13] + m[8]*m[5]*m[14] - m[8]*m[6]*m[13] - m[12]*m[5]*m[10] + m[12]*m[6]*m[9],
		m[0]*m[9]*m[14] - m[0]*m[10]*m[13] - m[8]*m[1]*m[14] + m[8]*m[2]*m[13] + m[12]*m[1]*m[10] - m[12]*m[2]*m[9],
		-m[0]*m[5]*m[14] + m[0]*m[6]*m[13] + m[4]*m[1]*m[14] - m[4]*m[2]*m[13] - m[12]*m[1]*m[6] + m[12]*m[2]*m[5],
		m[0]*m[5]*m[10] - m[0]*m[6]*m[9] - m[4]*m[1]*m[10] + m[4]*m[2]*m[9] + m[8]*m[1]*m[6] - m[8]*m[2]*m[5],
	}

	det := m[0]*tmp[0] + m[1]*tmp[4] + m[2]*tmp[8] + m[3]*tmp[12]

	if det == 0 {
		return Matrix[T]{}
	}

	det = 1 / det

	return Matrix[T]{
		tmp[0] * det, tmp[1] * det, tmp[2] * det, tmp[3] * det,
		tmp[4] * det, tmp[5] * det, tmp[6] * det, tmp[7] * det,
		tmp[8] * det, tmp[9] * det, tmp[10] * det, tmp[11] * det,
		tmp[12] * det, tmp[13] * det, tmp[14] * det, tmp[15] * det,
	}
}

// To3x3 returns the 3x3 affine transformation matrix (upper-left 3x3 block).
// This is useful for extracting the 3x3 rotation+scaling matrix from a 4x4 matrix.
func (m Matrix[T]) To3x3() Matrix[T] {
	return Matrix[T]{
		m[0], m[1], 0, m[3],
		m[4], m[5], 0, m[7],
		0, 0, 1, 0,
		m[12], m[13], 0, m[15],
	}
}

// ToColumnMajor returns the matrix in column-major order.
// This is the default storage order for matrices in this package.
func (m Matrix[T]) ToColumnMajor() Matrix[T] {
	return Matrix[T]{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
}

// ToRowMajor returns the matrix in row-major order.
// This is useful for interoperability with other systems that use row-major order.
func (m Matrix[T]) ToRowMajor() Matrix[T] {
	return Matrix[T]{
		m[0], m[1], m[2], m[3],
		m[4], m[5], m[6], m[7],
		m[8], m[9], m[10], m[11],
		m[12], m[13], m[14], m[15],
	}
}

// Base returns Matrix without its `w` components (without translation).
func (m Matrix[T]) Base() Matrix[T] {
	return Matrix[T]{
		m[0], m[1], m[2], 0,
		m[4], m[5], m[6], 0,
		m[8], m[9], m[10], 0,
		0, 0, 0, 1,
	}
}

// BasisPoints returns the basis vectors as points, with translation removed.
// This extracts the rotated and scaled coordinate axes from the matrix.
// Useful for determining the orientation and scaling of objects.
func (m Matrix[T]) BasisPoints() []Point[T] {
	return []Point[T]{
		{m[0], m[1]}, // X basis
		{m[4], m[5]}, // Y basis
		{m[8], m[9]}, // Z basis (projected to 2D)
	}
}

// BasisVectors returns the basis vectors for the X, Y, and Z axes as Vector3.
// The basis vectors represent the transformed coordinate axes.
// Index 0: X-axis basis vector (first column of upper-left 3x3)
// Index 1: Y-axis basis vector (second column of upper-left 3x3)
// Index 2: Z-axis basis vector (third column of upper-left 3x3)
func (m Matrix[T]) BasisVectors() []Vector3[T] {
	return []Vector3[T]{
		{m[0], m[1], m[2]},  // X basis
		{m[4], m[5], m[6]},  // Y basis
		{m[8], m[9], m[10]}, // Z basis
	}
}

// MaxBasisLengthXY returns the maximum length of the basis vectors in the XY plane.
func (m Matrix[T]) MaxBasisLengthXY() T {
	// The full basis computation requires computing the squared scaling factor
	// for translate/scale only matrices. This substantially limits the range of
	// precision for small and large scales. Instead, check for the common cases
	// and directly return the max scaling factor.
	if NearlyEqual(m[1], 0) && NearlyEqual(m[4], 0) {
		return max(Abs(m[0]), Abs(m[5]))
	}

	return T(math.Sqrt(ToFloat64(
		max(m[0]*m[0]+m[1]*m[1], m[4]*m[4]+m[5]*m[5]),
	)))
}

// BasisX returns the X-axis basis vector of the matrix.
func (m Matrix[T]) BasisX() Vector3[T] { return Vector3[T]{m[0], m[1], m[2]} }

// BasisY returns the Y-axis basis vector of the matrix.
func (m Matrix[T]) BasisY() Vector3[T] { return Vector3[T]{m[4], m[5], m[6]} }

// BasisZ returns the Z-axis basis vector of the matrix.
func (m Matrix[T]) BasisZ() Vector3[T] { return Vector3[T]{m[8], m[9], m[10]} }

// GetScale returns the scale factors along each axis as a Vector3.
// Extracts the scaling comp1nt from the transformation matrix.
// For non-uniform scaling or matrices with rotation, this returns the
// length of each basis vector.
func (m Matrix[T]) GetScale() Vector3[T] {
	basisX := Vector3[T]{m[0], m[1], m[2]}
	basisY := Vector3[T]{m[4], m[5], m[6]}
	basisZ := Vector3[T]{m[8], m[9], m[10]}
	return Vector3[T]{basisX.Length(), basisY.Length(), basisZ.Length()}
}

// GetDirectionScale returns the scale factor along the given direction vector.
// This computes how much the matrix scales a vector in the specified direction.
// Useful for anisotropic filtering and determining scaling in arbitrary directions.
func (m Matrix[T]) GetDirectionScale(dir Vector3[T]) T {
	return 1 / (m.Base().Invert().transformVector3D(dir.Normalize())).Length() * dir.Length()
}

// Scale scales the matrix by a given vector or scalar.
//
// Typical v types include:
// - [Scalar]
// - [Vector2]
// - [Vector3]
func (m Matrix[T]) Scale(s any) Matrix[T] {
	var v Vector3[T]
	switch s := s.(type) {
	case T:
		v = Vector3[T]{s, s, s}
	case Vector2[T]:
		v = Vector3[T]{s.X, s.Y, 1}
	case Vector3[T]:
		v = s
	default:
		panic("unsupported type for Scale: must be Scalar, Vector2, or Vector3")
	}
	sm := Matrix[T]{
		v.X, 0, 0, 0,
		0, v.Y, 0, 0,
		0, 0, v.Z, 0,
		0, 0, 0, 1,
	}
	return m.Mul(sm)
}

// Translate moves the matrix by a given vector.
//
// Typical t types include:
// - [Scalar]
// - [Vector2]
// - [Vector3]
func (m Matrix[T]) Translate(t any) Matrix[T] {
	var v Vector3[T]
	switch t := t.(type) {
	case T:
		v = Vector3[T]{t, t, t}
	case Vector2[T]:
		v = Vector3[T]{t.X, t.Y, 0}
	case Vector3[T]:
		v = t
	default:
		panic("unsupported type for Translate: must be Scalar, Vector2, or Vector3")
	}
	tm := Matrix[T]{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		v.X, v.Y, v.Z, 1,
	}
	return m.Mul(tm)
}

// Skew creates a shear matrix.
//
// Typical v types include:
// - [Scalar]
// - [Vector2]
// - [Vector3]
func (m Matrix[T]) Skew(s any) Matrix[T] {
	var v Vector3[T]
	switch s := s.(type) {
	case T:
		v = Vector3[T]{s, s, s}
	case Vector2[T]:
		v = Vector3[T]{s.X, s.Y, 0}
	case Vector3[T]:
		v = s
	default:
		panic("unsupported type for Skew: must be Scalar, Vector2, or Vector3")
	}

	sm := Matrix[T]{
		1, v.Y, 0, 0,
		v.X, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	return m.Mul(sm)
}

// RotateX creates a rotation matrix around the X axis.
func (m Matrix[T]) RotateX(angle Radians) Matrix[T] {
	cos, sin := m.CosSin(angle)
	rot := Matrix[T]{
		1, 0, 0, 0,
		0, cos, sin, 0,
		0, -sin, cos, 0,
		0, 0, 0, 1,
	}
	return m.Mul(rot)
}

// RotateY creates a rotation matrix around the Y axis.
func (m Matrix[T]) RotateY(angle Radians) Matrix[T] {
	cos, sin := m.CosSin(angle)
	rot := Matrix[T]{
		cos, 0, -sin, 0,
		0, 1, 0, 0,
		sin, 0, cos, 0,
		0, 0, 0, 1,
	}
	return m.Mul(rot)
}

// RotateZ creates a rotation matrix around the Z axis.
func (m Matrix[T]) RotateZ(angle Radians) Matrix[T] {
	cos, sin := m.CosSin(angle)
	rot := Matrix[T]{
		cos, sin, 0, 0,
		-sin, cos, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	return m.Mul(rot)
}

// Rotate creates a rotation matrix around an arbitrary axis using Rodrigues' formula.
//
// The resulting matrix rotates points around the specified axis by the given angle.
// The axis vector should be normalized before calling this function.
func (m Matrix[T]) Rotate(angle Radians, axis Vector3[T]) Matrix[T] {
	v := axis.Normalize()
	cos, sin := m.CosSin(angle)
	cosp := 1 - cos
	rm := Matrix[T]{
		cos + cosp*v.X*v.X,
		cosp*v.X*v.Y + v.Z*sin,
		cosp*v.X*v.Z - v.Y*sin,
		0,
		cosp*v.X*v.Y - v.Z*sin,
		cos + cosp*v.Y*v.Y,
		cosp*v.Y*v.Z + v.X*sin,
		0,
		cosp*v.X*v.Z + v.Y*sin,
		cosp*v.Y*v.Z - v.X*sin,
		cos + cosp*v.Z*v.Z,
		0,
		0, 0, 0, 1,
	}
	return m.Mul(rm)
}

// RotateQuat creates a rotation matrix from a quaternion.
//
// The quaternion should be normalized before calling this function.
// The resulting matrix applies the rotation defined by the quaternion to points in 3D space.
func (m Matrix[T]) RotateQuat(quat Quaternion[T]) Matrix[T] {
	x, y, z, w := quat.X, quat.Y, quat.Z, quat.W
	rm := Matrix[T]{
		1 - 2*(y*y+z*z), 2 * (x*y + z*w), 2 * (x*z - y*w), 0,
		2 * (x*y - z*w), 1 - 2*(x*x+z*z), 2 * (y*z + x*w), 0,
		2 * (x*z + y*w), 2 * (y*z - x*w), 1 - 2*(x*x+y*y), 0,
		0, 0, 0, 1,
	}
	return m.Mul(rm)
}

// Orthographic creates an orthographic projection matrix based on the current transformation.
// The result maps the rectangle [0, width] x [0, height] to normalized device coordinates (NDC) [-1, 1].
func (m Matrix[T]) Orthographic(size Size[T]) Matrix[T] {
	om := Matrix[T]{
		2 / size.Width, 0, 0, 0,
		0, 2 / size.Height, 0, 0,
		0, 0, 1, 0,
		-1, 1, T(1) / T(2), 1,
	}
	return m.Mul(om)
}

// LookAt creates a view matrix for a camera based on the current transformation.
func (m Matrix[T]) LookAt(position, target, up Vector3[T]) Matrix[T] {
	forward := target.Sub(position).Normalize()
	right := up.Cross(forward)
	upNorm := forward.Cross(right)

	lm := Matrix[T]{
		right.X, upNorm.X, forward.X, 0,
		right.Y, upNorm.Y, forward.Y, 0,
		right.Z, upNorm.Z, forward.Z, 0,
		-right.Dot(position), -upNorm.Dot(position), -forward.Dot(position), 1,
	}
	return m.Mul(lm)
}

// Perspective creates a perspective projection matrix.
func (m Matrix[T]) Perspective(fovY Radians, aspectRatio, zNear, zFar T) Matrix[T] {
	height := T(math.Tan(ToFloat64(fovY) * 0.5))
	width := height * aspectRatio

	pm := Matrix[T]{
		1 / width, 0, 0, 0,
		0, 1 / height, 0, 0,
		0, 0, zFar / (zFar - zNear), 1,
		0, 0, -(zFar * zNear) / (zFar - zNear), 0,
	}
	return m.Mul(pm)
}

// PerspectiveSize creates a perspective projection matrix from size.
func (m Matrix[T]) PerspectiveSize(fovY Radians, size Size[T], zNear, zFar T) Matrix[T] {
	aspectRatio := size.Width / size.Height
	return m.Perspective(fovY, aspectRatio, zNear, zFar)
}

// Transform transforms a geometry using the matrix.
//
// This applies the matrix to the geometry's types include:
// - [Point]  transforms a 2D point, suitable for 2D affine or projection transformations.
// - [Vector3] fill in the homogeneous components, and finally perform perspective division.
// - [Vector4] transforms a 4D vector, useful for homogeneous coordinates in 3D transformations and projections.
// - [Quad] transforms a quadrilateral by transforming each of its four points.
func (m Matrix[T]) Transform(t any) any {
	switch t := t.(type) {
	case Point[T]:
		return m.transformPoint(t)
	case Vector3[T]:
		return m.transformVector3D(t)
	case Vector4[T]:
		return m.transformVector4D(t)
	case Quad[T]:
		return Quad[T]{
			m.transformPoint(t[0]),
			m.transformPoint(t[1]),
			m.transformPoint(t[2]),
			m.transformPoint(t[3]),
		}
	default:
		panic(fmt.Sprintf("unsupported geometry type for Transform: %T", t))
	}
}

// TransformDirection do a linear transformation without perspective division.
// This is useful for transforming directions or normals.
//
// It applies the matrix to the geometry's types include:
// - [Vector2]
// - [Vector3]
// - [Vector4]
func (m Matrix[T]) TransformDirection(v any) any {
	switch v := v.(type) {
	case Vector2[T]:
		return Vector2[T]{X: v.X*m[0] + v.Y*m[4], Y: v.X*m[1] + v.Y*m[5]}
	case Vector3[T]:
		return Vector3[T]{
			X: v.X*m[0] + v.Y*m[4] + v.Z*m[8],
			Y: v.X*m[1] + v.Y*m[5] + v.Z*m[9],
			Z: v.X*m[2] + v.Y*m[6] + v.Z*m[10],
		}
	case Vector4[T]:
		return Vector4[T]{
			X: v.X*m[0] + v.Y*m[4] + v.Z*m[8],
			Y: v.X*m[1] + v.Y*m[5] + v.Z*m[9],
			Z: v.X*m[2] + v.Y*m[6] + v.Z*m[10],
			W: v.W,
		}
	default:
		panic(fmt.Sprintf("unsupported geometry type for TransformDirection: %T", v))
	}
}

// TransformHomogenous extends a 2D point to 3D homogeneous coordinates and transforms it.
func (m Matrix[T]) TransformHomogenous(p Point[T]) Vector3[T] {
	return Vector3[T]{
		X: p.X*m[0] + p.Y*m[4] + m[12],
		Y: p.X*m[1] + p.Y*m[5] + m[13],
		Z: p.X*m[3] + p.Y*m[7] + m[15],
	}
}

// transformPoint performs a 2D point transformation.
func (m Matrix[T]) transformPoint(p Point[T]) Point[T] {
	w := p.X*m[3] + p.Y*m[7] + m[15]
	r := Point[T]{
		X: p.X*m[0] + p.Y*m[4] + m[12],
		Y: p.X*m[1] + p.Y*m[5] + m[13],
	}
	if w != 0 {
		w = 1 / w
	}
	return Point[T]{r.X * w, r.Y * w}
}

// transformVector3D performs a 3D vector transformation.
// Automatically fill in the homogeneous components, and finally perform perspective division.
func (m Matrix[T]) transformVector3D(v Vector3[T]) Vector3[T] {
	w := v.X*m[3] + v.Y*m[7] + v.Z*m[11] + m[15]
	r := Vector3[T]{
		v.X*m[0] + v.Y*m[4] + v.Z*m[8] + m[12],
		v.X*m[1] + v.Y*m[5] + v.Z*m[9] + m[13],
		v.X*m[2] + v.Y*m[6] + v.Z*m[10] + m[14],
	}
	if w != 0 {
		w = 1 / w
	}
	return r.Scale(w)
}

// transformVector4D performs a 4D vector transformation.
// This is used for homogeneous coordinates in 3D transformations and projections.
func (m Matrix[T]) transformVector4D(v Vector4[T]) Vector4[T] {
	return Vector4[T]{
		v.X*m[0] + v.Y*m[4] + v.Z*m[8] + v.W*m[12],
		v.X*m[1] + v.Y*m[5] + v.Z*m[9] + v.W*m[13],
		v.X*m[2] + v.Y*m[6] + v.Z*m[10] + v.W*m[14],
		v.X*m[3] + v.Y*m[7] + v.Z*m[11] + v.W*m[15],
	}
}

// CosSin returns the cosine and sine of the given angle in radians.
// This is a helper function for rotation operations.
func (m Matrix[T]) CosSin(angle Radians) (cos, sin T) {
	sinVal := T(math.Sin(ToFloat64(angle)))
	if math.Abs(ToFloat64(sinVal)) == 1.0 {
		// 90 or 270 degrees
		return T(0), sinVal
	}

	cosVal := math.Cos(ToFloat64(angle))
	if math.Abs(cosVal) == 1.0 {
		// 0 or 180 degrees
		return T(cosVal), T(0)
	}

	return T(cosVal), sinVal
}

// String returns a string representation of the matrix.
func (m Matrix[T]) String() string {
	return fmt.Sprintf(
		"[%v %v %v %v]\n[%v %v %v %v]\n[%v %v %v %v]\n[%v %v %v %v]",
		m[0], m[1], m[2], m[3],
		m[4], m[5], m[6], m[7],
		m[8], m[9], m[10], m[11],
		m[12], m[13], m[14], m[15],
	)
}

// Decompose returns the matrix decomposition into translation, scale, shear, perspective, and rotation comp1nts.
// This is useful for extracting transformation parameters from a matrix.
func (m Matrix[T]) Decompose() MatrixDecomp[T] {
	// This is a simplified version - full decomposition is complex
	// For now, return basic comp1nts
	return MatrixDecomp[T]{
		Translation: Vector3[T]{m[12], m[13], m[14]},
		Scale:       m.GetScale(),
		Shear:       Shear[T]{XY: 0, XZ: 0, YZ: 0},
		Perspective: Vector4[T]{0, 0, 0, 1},
		Rotation:    Quaternion[T]{0, 0, 0, 1},
	}
}

// MatrixFlags represents the components of a matrix decomposition.
type MatrixFlags uint32

const (
	MatrixFlagsTranslation MatrixFlags = 1 << iota
	MatrixFlagsScale
	MatrixFlagsShear
	MatrixFlagsPerspective
	MatrixFlagsRotation
)

type MatrixDecomp[T TScalar] struct {
	Translation Vector3[T]
	Scale       Vector3[T]
	Shear       Shear[T]
	Perspective Vector4[T]
	Rotation    Quaternion[T]
}

func (md *MatrixDecomp[T]) Mask() MatrixFlags {
	var mask MatrixFlags

	// Check translation
	if md.Translation.X != 0 || md.Translation.Y != 0 || md.Translation.Z != 0 {
		mask |= MatrixFlagsTranslation
	}

	// Check scale
	if md.Scale.X != 1 || md.Scale.Y != 1 || md.Scale.Z != 1 {
		mask |= MatrixFlagsScale
	}

	// Check shear
	if md.Shear.XY != 0 || md.Shear.XZ != 0 || md.Shear.YZ != 0 {
		mask |= MatrixFlagsShear
	}

	// Check perspective
	if md.Perspective.X != 0 || md.Perspective.Y != 0 || md.Perspective.Z != 0 || md.Perspective.W != 1 {
		mask |= MatrixFlagsPerspective
	}

	// Check rotation (identity quaternion has W=1, others=0)
	if md.Rotation.X != 0 || md.Rotation.Y != 0 || md.Rotation.Z != 0 || md.Rotation.W != 1 {
		mask |= MatrixFlagsRotation
	}

	return mask
}
