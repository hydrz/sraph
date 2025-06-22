package geom

import (
	"fmt"
	"math"
)

// Matrix represents a 4x4 transformation matrix using column-major storage.
// Each Matrix instance defines a linear transformation in 3D space that can include
// translation, rotation, scaling, shearing, and perspective projection.
//
// Matrix follows the DirectX/Vulkan convention for normalized device coordinates (NDC):
//   - Left-handed coordinate system with X right, Y up, Z away from viewer
//   - Positive rotation is clockwise about the rotation axis
//   - NDC bounds: x,y ∈ [-1,1], z ∈ [0,1] where 0 is near plane, 1 is far plane
//   - NDC origin at (0, 0, 0.5) representing the center of the NDC cube
//
// This differs from OpenGL which uses right-handed coordinates and z ∈ [-1,1].
//
// The zero value is not a valid transformation matrix. Use NewMatrix to create
// an identity matrix for initialization.
//
// For more details, see:
//
//	https://learn.microsoft.com/zh-cn/windows/win32/learnwin32/appendix--matrix-transforms
//	https://www.opengl-tutorial.org/beginners-tutorials/tutorial-3-matrices/
type Matrix [16]Scalar

// NewMatrix returns an identity matrix that applies no transformation.
func NewMatrix() Matrix {
	return Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// At returns the matrix element at the specified row and column (0-based indices).
func (m Matrix) At(row int, col int) Scalar {
	// Column-major order: index = col*4 + row
	return m[col*4+row]
}

// Set modifies the matrix element at the specified row and column (0-based indices).
// Note: This creates a copy, so modifications won't persist. Use a pointer for persistent changes.
func (m *Matrix) Set(row int, col int, value Scalar) {
	m[col*4+row] = value
}

// Add returns the element-wise sum of two matrices.
func (m Matrix) Add(other Matrix) Matrix {
	o := other
	return Matrix{
		m[0] + o[0], m[1] + o[1], m[2] + o[2], m[3] + o[3],
		m[4] + o[4], m[5] + o[5], m[6] + o[6], m[7] + o[7],
		m[8] + o[8], m[9] + o[9], m[10] + o[10], m[11] + o[11],
		m[12] + o[12], m[13] + o[13], m[14] + o[14], m[15] + o[15],
	}
}

// Sub returns the element-wise difference of two matrices.
func (m Matrix) Sub(other Matrix) Matrix {
	o := other
	return Matrix{
		m[0] - o[0], m[1] - o[1], m[2] - o[2], m[3] - o[3],
		m[4] - o[4], m[5] - o[5], m[6] - o[6], m[7] - o[7],
		m[8] - o[8], m[9] - o[9], m[10] - o[10], m[11] - o[11],
		m[12] - o[12], m[13] - o[13], m[14] - o[14], m[15] - o[15],
	}
}

// Mul returns the matrix product of two matrices (standard matrix multiplication).
func (m Matrix) Mul(other Matrix) Matrix {
	o := other
	return Matrix{
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

// Equal reports whether two matrices are approximately equal within floating-point tolerance.
func (m Matrix) Equal(other Matrix) bool {
	o := other
	for i := 0; i < 16; i++ {
		if !NearlyEqual(m[i], o[i]) {
			return false
		}
	}
	return true
}

// IsFinite reports whether all matrix elements are finite numbers.
func (m Matrix) IsFinite() bool {
	for i := 0; i < 16; i++ {
		if !IsFinite(m[i]) {
			return false
		}
	}
	return true
}

// IsAffine reports whether the matrix represents an affine transformation.
func (m Matrix) IsAffine() bool {
	return NearlyEqual(m[2], 0) && NearlyEqual(m[3], 0) && NearlyEqual(m[6], 0) && NearlyEqual(m[7], 0) &&
		NearlyEqual(m[8], 0) && NearlyEqual(m[9], 0) && NearlyEqual(m[10], 1) && NearlyEqual(m[11], 0) &&
		NearlyEqual(m[14], 0) && NearlyEqual(m[15], 1)
}

// IsIdentity reports whether the matrix is an identity matrix.
func (m Matrix) IsIdentity() bool {
	return NearlyEqual(m[0], 1) && NearlyEqual(m[1], 0) && NearlyEqual(m[2], 0) && NearlyEqual(m[3], 0) &&
		NearlyEqual(m[4], 0) && NearlyEqual(m[5], 1) && NearlyEqual(m[6], 0) && NearlyEqual(m[7], 0) &&
		NearlyEqual(m[8], 0) && NearlyEqual(m[9], 0) && NearlyEqual(m[10], 1) && NearlyEqual(m[11], 0) &&
		NearlyEqual(m[12], 0) && NearlyEqual(m[13], 0) && NearlyEqual(m[14], 0) && NearlyEqual(m[15], 1)
}

// IsInvertible reports whether the matrix can be inverted.
func (m Matrix) IsInvertible() bool { return m.Determinant() != 0 }

// Determinant returns the determinant of the matrix.
func (m Matrix) Determinant() Scalar {
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

// HasPerspective reports whether the matrix contains perspective projection.
func (m Matrix) HasPerspective() bool {
	return !NearlyEqual(m[3], 0) || !NearlyEqual(m[7], 0) || !NearlyEqual(m[11], 0) || !NearlyEqual(m[15], 1)
}

// HasPerspective2D reports whether the matrix contains 2D perspective transformation.
func (m Matrix) HasPerspective2D() bool {
	return !NearlyEqual(m[3], 0) || !NearlyEqual(m[7], 0) || !NearlyEqual(m[15], 1)
}

// HasTranslation reports whether the matrix contains translation components.
func (m Matrix) HasTranslation() bool {
	return !NearlyEqual(m[12], 0) || !NearlyEqual(m[13], 0)
}

// IsAxisAligned reports whether the matrix transformation is axis-aligned.
func (m Matrix) IsAxisAligned() bool {
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

	// Ensure that none of the basis vectors overlap
	if (bti(v[0])+bti(v[3])+bti(v[6]) != 1) ||
		(bti(v[1])+bti(v[4])+bti(v[7]) != 1) ||
		(bti(v[2])+bti(v[5])+bti(v[8]) != 1) {
		return false
	}

	return true
}

// IsAxisAligned2D reports whether the matrix is axis-aligned in the XY plane.
func (m Matrix) IsAxisAligned2D() bool {
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

// IsTranslationOnly reports whether the matrix contains only translation.
func (m Matrix) IsTranslationOnly() bool {
	return NearlyEqual(m[0], 1) && NearlyEqual(m[1], 0) && NearlyEqual(m[2], 0) && NearlyEqual(m[3], 0) &&
		NearlyEqual(m[4], 0) && NearlyEqual(m[5], 1) && NearlyEqual(m[6], 0) && NearlyEqual(m[7], 0) &&
		NearlyEqual(m[8], 0) && NearlyEqual(m[9], 0) && NearlyEqual(m[10], 1) && NearlyEqual(m[11], 0) &&
		NearlyEqual(m[15], 1)
}

// IsTranslationScaleOnly reports whether the matrix contains only translation and uniform scaling.
func (m Matrix) IsTranslationScaleOnly() bool {
	return !NearlyEqual(m[0], 0) && NearlyEqual(m[1], 0) && NearlyEqual(m[2], 0) && NearlyEqual(m[3], 0) &&
		NearlyEqual(m[4], 0) && !NearlyEqual(m[5], 0) && NearlyEqual(m[6], 0) && NearlyEqual(m[7], 0) &&
		NearlyEqual(m[8], 0) && NearlyEqual(m[9], 0) && !NearlyEqual(m[10], 0) && NearlyEqual(m[11], 0) &&
		NearlyEqual(m[15], 1)
}

// Transpose returns the transpose of the matrix (rows become columns).
func (m Matrix) Transpose() Matrix {
	return Matrix{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
}

// Invert returns the inverse matrix that undoes this transformation.
func (m Matrix) Invert() Matrix {
	tmp := Matrix{
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
		return Matrix{}
	}

	det = 1 / det

	return Matrix{
		tmp[0] * det, tmp[1] * det, tmp[2] * det, tmp[3] * det,
		tmp[4] * det, tmp[5] * det, tmp[6] * det, tmp[7] * det,
		tmp[8] * det, tmp[9] * det, tmp[10] * det, tmp[11] * det,
		tmp[12] * det, tmp[13] * det, tmp[14] * det, tmp[15] * det,
	}
}

// To3x3 returns a 3x3 matrix containing the upper-left 3x3 portion.
func (m Matrix) To3x3() Matrix {
	return Matrix{
		m[0], m[1], 0, m[3],
		m[4], m[5], 0, m[7],
		0, 0, 1, 0,
		m[12], m[13], 0, m[15],
	}
}

// ToColumnMajor returns the matrix elements in column-major order.
func (m Matrix) ToColumnMajor() Matrix {
	return Matrix{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
}

// ToRowMajor returns the matrix elements in row-major order.
func (m Matrix) ToRowMajor() Matrix {
	return Matrix{
		m[0], m[1], m[2], m[3],
		m[4], m[5], m[6], m[7],
		m[8], m[9], m[10], m[11],
		m[12], m[13], m[14], m[15],
	}
}

// Base returns the matrix with translation components removed.
func (m Matrix) Base() Matrix {
	return Matrix{
		m[0], m[1], m[2], 0,
		m[4], m[5], m[6], 0,
		m[8], m[9], m[10], 0,
		0, 0, 0, 1,
	}
}

// BasisPoints returns the 2D basis vectors as points.
func (m Matrix) BasisPoints() []Point {
	return []Point{
		NewPoint(m[0], m[1]), // X basis
		NewPoint(m[4], m[5]), // Y basis
		NewPoint(m[8], m[9]), // Z basis (projected to 2D)
	}
}

// BasisVectors returns the 3D basis vectors representing the transformed coordinate axes.
func (m Matrix) BasisVectors() []Vector3 {
	return []Vector3{
		NewVector3(m[0], m[1], m[2]),  // X basis
		NewVector3(m[4], m[5], m[6]),  // Y basis
		NewVector3(m[8], m[9], m[10]), // Z basis
	}
}

// MaxBasisLengthXY returns the maximum scaling factor in the XY plane.
func (m Matrix) MaxBasisLengthXY() Scalar {
	if NearlyEqual(m[1], 0) && NearlyEqual(m[4], 0) {
		return max(Abs(m[0]), Abs(m[5]))
	}

	return ToScalar(math.Sqrt(ToFloat64(
		max(m[0]*m[0]+m[1]*m[1], m[4]*m[4]+m[5]*m[5]),
	)))
}

// BasisX returns the transformed X-axis basis vector.
func (m Matrix) BasisX() Vector3 { return NewVector3(m[0], m[1], m[2]) }

// BasisY returns the transformed Y-axis basis vector.
func (m Matrix) BasisY() Vector3 { return NewVector3(m[4], m[5], m[6]) }

// BasisZ returns the transformed Z-axis basis vector.
func (m Matrix) BasisZ() Vector3 { return NewVector3(m[8], m[9], m[10]) }

// GetScale returns the scaling factors along each axis.
func (m Matrix) GetScale() Vector3 {
	basisX := NewVector3(m[0], m[1], m[2])
	basisY := NewVector3(m[4], m[5], m[6])
	basisZ := NewVector3(m[8], m[9], m[10])
	return NewVector3(basisX.Length(), basisY.Length(), basisZ.Length())
}

// GetDirectionScale returns the scaling factor applied to vectors in the specified direction.
func (m Matrix) GetDirectionScale(dir Vector3) Scalar {
	return 1 / dir.Normalize().Transform(m.Base().Invert()).Length() * dir.Length()
}

// Scale returns a new matrix with scaling transformation applied.
func (m Matrix) Scale(s any) Matrix {
	var v Vector3
	switch s := s.(type) {
	case Scalar:
		v = NewVector3(s, s, s)
	case Vector2:
		v = NewVector3(s.X(), s.Y(), 1)
	case Vector3:
		v = s
	default:
		panic("unsupported type for Scale: must be Scalar, Vector2, or Vector3")
	}
	sm := Matrix{
		v.X(), 0, 0, 0,
		0, v.Y(), 0, 0,
		0, 0, v.Z(), 0,
		0, 0, 0, 1,
	}
	return m.Mul(sm)
}

// Translate returns a new matrix with translation transformation applied.
func (m Matrix) Translate(t any) Matrix {
	var v Vector3
	switch t := t.(type) {
	case Scalar:
		v = NewVector3(t, t, t)
	case Vector2:
		v = NewVector3(t.X(), t.Y(), 0)
	case Vector3:
		v = t
	default:
		panic("unsupported type for Translate: must be Scalar, Vector2, or Vector3")
	}
	tm := Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		v.X(), v.Y(), v.Z(), 1,
	}
	return m.Mul(tm)
}

// Skew returns a new matrix with shear transformation applied.
func (m Matrix) Skew(s any) Matrix {
	var v Vector3
	switch s := s.(type) {
	case Scalar:
		v = NewVector3(s, s, s)
	case Vector2:
		v = NewVector3(s.X(), s.Y(), 0)
	case Vector3:
		v = s
	default:
		panic("unsupported type for Skew: must be Scalar, Vector2, or Vector3")
	}

	sm := Matrix{
		1, v.Y(), 0, 0,
		v.X(), 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	return m.Mul(sm)
}

// RotateX returns a new matrix with rotation around the X-axis applied.
func (m Matrix) RotateX(angle Radians) Matrix {
	cos, sin := m.CosSin(angle)
	rot := Matrix{
		1, 0, 0, 0,
		0, cos, sin, 0,
		0, -sin, cos, 0,
		0, 0, 0, 1,
	}
	return m.Mul(rot)
}

// RotateY returns a new matrix with rotation around the Y-axis applied.
func (m Matrix) RotateY(angle Radians) Matrix {
	cos, sin := m.CosSin(angle)
	rot := Matrix{
		cos, 0, -sin, 0,
		0, 1, 0, 0,
		sin, 0, cos, 0,
		0, 0, 0, 1,
	}
	return m.Mul(rot)
}

// RotateZ returns a new matrix with rotation around the Z-axis applied.
func (m Matrix) RotateZ(angle Radians) Matrix {
	cos, sin := m.CosSin(angle)
	rot := Matrix{
		cos, sin, 0, 0,
		-sin, cos, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	return m.Mul(rot)
}

// Rotate returns a new matrix with rotation around an arbitrary axis applied.
func (m Matrix) Rotate(angle Radians, axis Vector3) Matrix {
	v := axis.Normalize()
	cos, sin := m.CosSin(angle)
	cosp := 1 - cos
	rm := Matrix{
		cos + cosp*v.X()*v.X(),
		cosp*v.X()*v.Y() + v.Z()*sin,
		cosp*v.X()*v.Z() - v.Y()*sin,
		0,
		cosp*v.X()*v.Y() - v.Z()*sin,
		cos + cosp*v.Y()*v.Y(),
		cosp*v.Y()*v.Z() + v.X()*sin,
		0,
		cosp*v.X()*v.Z() + v.Y()*sin,
		cosp*v.Y()*v.Z() - v.X()*sin,
		cos + cosp*v.Z()*v.Z(),
		0,
		0, 0, 0, 1,
	}
	return m.Mul(rm)
}

// RotateQuat returns a new matrix with rotation from a quaternion applied.
func (m Matrix) RotateQuat(quat Quaternion) Matrix {
	x, y, z, w := quat.X(), quat.Y(), quat.Z(), quat.W()
	rm := Matrix{
		1 - 2*(y*y+z*z), 2 * (x*y + z*w), 2 * (x*z - y*w), 0,
		2 * (x*y - z*w), 1 - 2*(x*x+z*z), 2 * (y*z + x*w), 0,
		2 * (x*z + y*w), 2 * (y*z - x*w), 1 - 2*(x*x+y*y), 0,
		0, 0, 0, 1,
	}
	return m.Mul(rm)
}

// Orthographic returns a new matrix with orthographic projection applied.
func (m Matrix) Orthographic(size Size) Matrix {
	om := Matrix{
		2 / size.Width(), 0, 0, 0,
		0, 2 / size.Height(), 0, 0,
		0, 0, 1, 0,
		-1, 1, 0.5, 1,
	}
	return m.Mul(om)
}

// LookAt returns a new matrix with view transformation applied.
func (m Matrix) LookAt(position, target, up Vector3) Matrix {
	forward := target.Sub(position).Normalize()
	right := up.Cross(forward)
	upNorm := forward.Cross(right)

	lm := Matrix{
		right.X(), upNorm.X(), forward.X(), 0,
		right.Y(), upNorm.Y(), forward.Y(), 0,
		right.Z(), upNorm.Z(), forward.Z(), 0,
		-right.Dot(position), -upNorm.Dot(position), -forward.Dot(position), 1,
	}
	return m.Mul(lm)
}

// Perspective returns a new matrix with perspective projection applied.
func (m Matrix) Perspective(fovY Radians, aspectRatio, zNear, zFar Scalar) Matrix {
	height := Scalar(math.Tan(ToFloat64(fovY) * 0.5))
	width := height * aspectRatio

	pm := Matrix{
		1 / width, 0, 0, 0,
		0, 1 / height, 0, 0,
		0, 0, zFar / (zFar - zNear), 1,
		0, 0, -(zFar * zNear) / (zFar - zNear), 0,
	}
	return m.Mul(pm)
}

// PerspectiveSize returns a new matrix with perspective projection applied using viewport size.
func (m Matrix) PerspectiveSize(fovY Radians, size Size, zNear, zFar Scalar) Matrix {
	aspectRatio := size.Width() / size.Height()
	return m.Perspective(fovY, aspectRatio, zNear, zFar)
}

// CosSin returns the cosine and sine of the given angle.
func (m Matrix) CosSin(angle Radians) (cos, sin Scalar) {
	sinVal := Scalar(math.Sin(ToFloat64(angle)))
	if math.Abs(ToFloat64(sinVal)) == 1.0 {
		// 90 or 270 degrees
		return 0, sinVal
	}

	cosVal := math.Cos(ToFloat64(angle))
	if math.Abs(cosVal) == 1.0 {
		// 0 or 180 degrees
		return ToScalar(cosVal), 0
	}

	return ToScalar(cosVal), sinVal
}

// String returns a string representation of the matrix in a 4x4 grid format.
func (m Matrix) String() string {
	return fmt.Sprintf(
		"[%v %v %v %v]\n[%v %v %v %v]\n[%v %v %v %v]\n[%v %v %v %v]",
		m[0], m[1], m[2], m[3],
		m[4], m[5], m[6], m[7],
		m[8], m[9], m[10], m[11],
		m[12], m[13], m[14], m[15],
	)
}

// Decompose extracts the transformation components from the matrix.
func (m Matrix) Decompose() MatrixDecomp {
	return MatrixDecomp{
		Translation: NewVector3(m[12], m[13], m[14]),
		Scale:       m.GetScale(),
		Shear:       NewShear(0, 0, 0),
		Perspective: NewVector4(0, 0, 0, 1),
		Rotation:    NewQuaternion(0, 0, 0, 1),
	}
}

// MatrixFlags represents bitmask flags for matrix transformation components.
type MatrixFlags uint32

const (
	MatrixFlagsTranslation MatrixFlags = 1 << iota // Matrix contains translation
	MatrixFlagsScale                               // Matrix contains scaling
	MatrixFlagsShear                               // Matrix contains shearing
	MatrixFlagsPerspective                         // Matrix contains perspective projection
	MatrixFlagsRotation                            // Matrix contains rotation
)

// MatrixDecomp represents the decomposed components of a transformation matrix.
type MatrixDecomp struct {
	Translation Vector3    // Translation vector
	Scale       Vector3    // Scale factors for each axis
	Shear       Shear      // Shear transformation parameters
	Perspective Vector4    // Perspective transformation parameters
	Rotation    Quaternion // Rotation as quaternion
}

// Mask returns a bitmask indicating which transformation components are present.
func (md MatrixDecomp) Mask() MatrixFlags {
	var mask MatrixFlags

	// Check translation
	if md.Translation.X() != 0 || md.Translation.Y() != 0 || md.Translation.Z() != 0 {
		mask |= MatrixFlagsTranslation
	}

	// Check scale
	if md.Scale.X() != 1 || md.Scale.Y() != 1 || md.Scale.Z() != 1 {
		mask |= MatrixFlagsScale
	}

	// Check shear
	if md.Shear.XY() != 0 || md.Shear.XZ() != 0 || md.Shear.YZ() != 0 {
		mask |= MatrixFlagsShear
	}

	// Check perspective
	if md.Perspective.X() != 0 || md.Perspective.Y() != 0 || md.Perspective.Z() != 0 || md.Perspective.W() != 1 {
		mask |= MatrixFlagsPerspective
	}

	// Check rotation (identity quaternion has W=1, others=0)
	if md.Rotation.X() != 0 || md.Rotation.Y() != 0 || md.Rotation.Z() != 0 || md.Rotation.W() != 1 {
		mask |= MatrixFlagsRotation
	}

	return mask
}
