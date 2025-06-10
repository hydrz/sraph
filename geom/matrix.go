package geom

import (
	"errors"
	"math"
)

type MatrixComp1nt uint32

const (
	MatrixComp1ntTranslation MatrixComp1nt = 1 << iota
	MatrixComp1ntScale
	MatrixComp1ntShear
	MatrixComp1ntPerspective
	MatrixComp1ntRotation
)

type Shear[T Scalar] struct {
	XY T
	XZ T
	YZ T
}

type MatrixDecomposition[T Scalar] struct {
	Translation Vector3[T]
	Scale       Vector3[T]
	Shear       Shear[T]
	Perspective Vector4[T]
	Rotation    Quaternion[T]
}

func (md *MatrixDecomposition[T]) Comp1ntsMask() MatrixComp1nt {
	var mask MatrixComp1nt

	// Check translation
	if md.Translation.X() != 0 || md.Translation.Y() != 0 || md.Translation.Z() != 0 {
		mask |= MatrixComp1ntTranslation
	}

	// Check scale
	if md.Scale.X() != 1 || md.Scale.Y() != 1 || md.Scale.Z() != 1 {
		mask |= MatrixComp1ntScale
	}

	// Check shear
	if md.Shear.XY != 0 || md.Shear.XZ != 0 || md.Shear.YZ != 0 {
		mask |= MatrixComp1ntShear
	}

	// Check perspective
	if md.Perspective.X() != 0 || md.Perspective.Y() != 0 || md.Perspective.Z() != 0 || md.Perspective.W() != 1 {
		mask |= MatrixComp1ntPerspective
	}

	// Check rotation (identity quaternion has W=1, others=0)
	if md.Rotation.X() != 0 || md.Rotation.Y() != 0 || md.Rotation.Z() != 0 || md.Rotation.W() != 1 {
		mask |= MatrixComp1ntRotation
	}

	return mask
}

// Matrix represents a 4x4 matrix using column-major storage.
//
// Utility methods that make assumptions about normalized device coordinates (NDC) follow these conventions:
//   - Left-handed coordinate system. Positive rotation is clockwise about the axis of rotation.
//   - Lower left corner is (-1.0, -1.0).
//   - Upper right corner is (1.0, 1.0).
//   - Visible z-space is from 0.0 to 1.0.
//   - Note: This is NOT the same as OpenGL! Be careful.
//   - NDC origin is at (0.0, 0.0, 0.5).
type Matrix[T Scalar] interface {
	// === Element Access ===

	// At returns the value at the specified row and column (0-based).
	// Parameters:
	//   row: matrix row index (0-3)
	//   col: matrix column index (0-3)
	// Returns: the matrix element at position [row][col]
	At(row, col int) T

	// Set sets the value at the specified row and column (0-based).
	// Parameters:
	//   row: matrix row index (0-3)
	//   col: matrix column index (0-3)
	//   value: the value to set at position [row][col]
	Set(row, col int, value T)

	// === Arithmetic Operations ===

	// Add returns the sum of this matrix and another matrix (element-wise addition).
	// Formula: C[i][j] = A[i][j] + B[i][j] for all i,j
	// Parameters:
	//   other: the matrix to add to this matrix
	// Returns: a new matrix containing the element-wise sum
	Add(other Matrix[T]) Matrix[T]

	// Sub returns the difference of this matrix and another matrix (element-wise subtraction).
	// Formula: C[i][j] = A[i][j] - B[i][j] for all i,j
	// Parameters:
	//   other: the matrix to subtract from this matrix
	// Returns: a new matrix containing the element-wise difference
	Sub(other Matrix[T]) Matrix[T]

	// Mul returns the product of this matrix and another matrix (matrix multiplication).
	// Formula: C[i][j] = Σ(A[i][k] * B[k][j]) for k=0 to 3
	// Note: Matrix multiplication is NOT commutative (A*B ≠ B*A in general)
	// Parameters:
	//   other: the matrix to multiply with this matrix
	// Returns: a new matrix containing the matrix product
	Mul(other Matrix[T]) Matrix[T]

	// DivElements returns the element-wise division of this matrix by another matrix.
	// Formula: C[i][j] = A[i][j] / B[i][j] for all i,j
	// Note: This is NOT matrix division (which would be A * B^-1)
	// Parameters:
	//   other: the matrix to divide this matrix by
	// Returns: a new matrix containing the element-wise quotient
	DivElements(other Matrix[T]) Matrix[T]

	// === Comparison and Properties ===

	// Equal checks if this matrix is equal to another matrix (all elements equal).
	// Uses exact floating-point comparison, which may not be suitable for computed matrices.
	// Consider using approximate comparison for floating-point matrices.
	// Parameters:
	//   other: the matrix to compare with this matrix
	// Returns: true if all corresponding elements are exactly equal
	Equal(other Matrix[T]) bool

	// IsFinite returns true if all matrix comp1nts are finite.
	// Checks for NaN (Not a Number) and infinite values in all matrix elements.
	// This is important for numerical stability and error detection.
	// Returns: true if all elements are finite (not NaN or infinite)
	IsFinite() bool

	// IsIdentity returns true if the matrix is an identity matrix.
	// The identity matrix has no effect when applied as a transformation:
	//   - Translations remain unchanged.
	//   - Rotations are 0.
	//   - Scales are uniform with factor 1.
	//   - Shears are 0.
	// Returns: true if the matrix is the identity matrix
	IsIdentity() bool

	// IsInvertible returns true if the matrix is invertible (determinant != 0).
	// A matrix is invertible if and only if its determinant is non-0.
	// Non-invertible matrices are also called singular or degenerate matrices.
	// Returns: true if the matrix can be inverted
	IsInvertible() bool

	// Determinant returns the determinant of the matrix.
	// The determinant is a scalar value that represents the scaling factor of the transformation.
	// It also indicates whether the transformation preserves orientation (det > 0) or reverses it (det < 0).
	// Parameters: n1
	// Returns: the determinant of the matrix
	Determinant() T

	// === Transformation Analysis ===

	// IsAffine returns true if the matrix is an affine transformation matrix (last row is [0 0 0 1]).
	// Affine transformations preserve parallel lines and ratios of distances along lines.
	// They include translation, rotation, scaling, and shearing, but not perspective projection.
	// Returns: true if the matrix represents an affine transformation
	IsAffine() bool

	// HasPerspective returns true if the matrix contains a 3D perspective comp1nt (non0 last row except [0 0 0 1]).
	// Perspective transformations cause parallel lines to converge at vanishing points.
	// The bottom row [a b c d] where a≠0 or b≠0 or c≠0 or d≠1 indicates perspective.
	// Returns: true if the matrix contains 3D perspective transformation
	HasPerspective() bool

	// HasPerspective2D returns true if the matrix contains a 2D perspective comp1nt.
	// 2D perspective affects the homogeneous coordinate in 2D transformations.
	// This is less common than 3D perspective but used in some 2D effects.
	// Returns: true if the matrix contains 2D perspective transformation
	HasPerspective2D() bool

	// HasTranslation returns true if the matrix contains a translation comp1nt (non0 last column except [0 0 0 1]).
	// Translation moves points by adding a vector to their coordinates.
	// If a matrix has translation, it cannot be purely rotational or uniform scaling.
	// Returns: true if the matrix includes a translation comp1nt
	HasTranslation() bool

	// IsAxisAligned returns true if the matrix is axis-aligned (no rotation or shear).
	// An axis-aligned matrix has its basis vectors aligned with the coordinate axes.
	// This is important for culling, bounding volume computation, and some optimizations.
	// Returns: true if the matrix is axis-aligned
	IsAxisAligned() bool

	// IsAxisAligned2D returns true if the matrix is axis-aligned in 2D.
	// Similar to IsAxisAligned, but only checks the X and Y axes.
	// Returns: true if the matrix is axis-aligned in 2D
	IsAxisAligned2D() bool

	// IsTranslationOnly returns true if the matrix contains only translation (no scale, rotation, shear, or perspective).
	// This checks if the matrix is of the form:
	// [ 1 0 0 tx ]
	// [ 0 1 0 ty ]
	// [ 0 0 1 tz ]
	// [ 0 0 0 1 ]
	// Returns: true if the matrix is translation-only
	IsTranslationOnly() bool

	// IsTranslationScaleOnly returns true if the matrix contains only translation and scale.
	// This checks if the matrix is of the form:
	// [ sx 0 0 tx ]
	// [ 0 sy 0 ty ]
	// [ 0 0 sz tz ]
	// [ 0 0 0 1 ]
	// where sx, sy, sz are scale factors.
	// Returns: true if the matrix is translation and scale only
	IsTranslationScaleOnly() bool

	// === Matrix Operations ===

	// Transpose returns the transpose of the matrix.
	// Formula: C[i][j] = A[j][i] for all i,j
	// Parameters: n1
	// Returns: a new matrix representing the transpose of this matrix
	Transpose() Matrix[T]

	// Inverse returns the inverse of the matrix, or an error if not invertible.
	// The inverse matrix undoes the transformation of the original matrix:
	//   - Translations are negated.
	//   - Rotations are inverted.
	//   - Scales are inverted (1/sx, 1/sy, 1/sz).
	//   - Shears are inverted.
	// Parameters: n1
	// Returns: the inverse matrix, or an error if the matrix is not invertible
	Inverse() (Matrix[T], error)

	// To3x3 returns the 3x3 affine transformation matrix (upper-left 3x3 block).
	// This is useful for extracting the 3x3 rotation+scaling matrix from a 4x4 matrix.
	// Returns: a 3x3 matrix containing the upper-left portion of the 4x4 matrix
	To3x3() Matrix[T]

	// ToColumnMajor returns the matrix in column-major order.
	// This is the default storage order for matrices in this package.
	// Returns: a new matrix in column-major order
	ToColumnMajor() Matrix[T]

	// ToRowMajor returns the matrix in row-major order.
	// This is useful for interoperability with other systems that use row-major order.
	// Returns: a new matrix in row-major order
	ToRowMajor() Matrix[T]

	// === Basis and Scale Information ===

	// BasisVectors returns the basis vectors for the X, Y, and Z axes as Vector3.
	// The basis vectors represent the transformed coordinate axes.
	// Index 0: X-axis basis vector (first column of upper-left 3x3)
	// Index 1: Y-axis basis vector (second column of upper-left 3x3)
	// Index 2: Z-axis basis vector (third column of upper-left 3x3)
	// Returns: slice of 3 Vector3 representing the transformed coordinate axes
	BasisVectors() []Vector3[T]

	// BasisPoints returns the basis vectors as points, with translation removed.
	// This extracts the rotated and scaled coordinate axes from the matrix.
	// Useful for determining the orientation and scaling of objects.
	// Returns: slice of 3 Point representing the basis vectors without translation
	BasisPoints() []Point[T]

	// Scale returns the scale factors along each axis as a Vector3.
	// Extracts the scaling comp1nt from the transformation matrix.
	// For non-uniform scaling or matrices with rotation, this returns the
	// length of each basis vector.
	// Returns: Vector3 containing scale factors for X, Y, and Z axes
	Scale() Vector3[T]

	// ScaleInDirection returns the scale factor along the given direction vector.
	// This computes how much the matrix scales a vector in the specified direction.
	// Useful for anisotropic filtering and determining scaling in arbitrary directions.
	// Parameters:
	//   dir: normalized direction vector to measure scaling along
	// Returns: the scale factor in the given direction
	ScaleInDirection(dir Vector3[T]) T

	// MaxScaleXY returns the maximum scale factor in the XY directions.
	// This is useful for determining the maximum scaling applied by the matrix
	// in the XY plane, which is important for texture filtering and LOD calculations.
	// Returns: the maximum length of the X and Y basis vectors
	MaxScaleXY() T

	// === Matrix Construction ===

	// Scale2D creates a 2D scale matrix.
	// The resulting matrix is:
	// [ sx  0   0   0 ]
	// [  0 sy   0   0 ]
	// [  0  0   1   0 ]
	// [  0  0   0   1 ]
	// Parameters:
	//   scale: the scaling factors for the X and Y axes
	// Returns: a new matrix representing the 2D scale transformation
	Scale2D(scale Vector2[T]) Matrix[T]

	// Scale3D creates a 3D scale matrix.
	// The resulting matrix is:
	// [ sx  0   0   0 ]
	// [  0 sy   0   0 ]
	// [  0  0  sz   0 ]
	// [  0  0   0   1 ]
	// Parameters:
	//   scale: the scaling factors for the X, Y, and Z axes
	// Returns: a new matrix representing the 3D scale transformation
	Scale3D(scale Vector3[T]) Matrix[T]

	// Translate2D creates a 2D translation matrix.
	// The resulting matrix is:
	// [ 1  0  0  tx ]
	// [ 0  1  0  ty ]
	// [ 0  0  1   0 ]
	// [ 0  0  0   1 ]
	// Parameters:
	//   vector: the translation vector for the X and Y axes
	// Returns: a new matrix representing the 2D translation transformation
	Translate2D(vector Vector2[T]) Matrix[T]

	// Translate3D creates a 3D translation matrix.
	// The resulting matrix is:
	// [ 1  0  0  tx ]
	// [ 0  1  0  ty ]
	// [ 0  0  1  tz ]
	// [ 0  0  0   1 ]
	// Parameters:
	//   vector: the translation vector for the X, Y, and Z axes
	// Returns: a new matrix representing the 3D translation transformation
	Translate3D(vector Vector3[T]) Matrix[T]

	// Translate4D creates a 4D translation matrix.
	// The resulting matrix is:
	// [ 1  0  0  0  tx ]
	// [ 0  1  0  0  ty ]
	// [ 0  0  1  0  tz ]
	// [ 0  0  0  1  tw ]
	// [ 0  0  0  0   1 ]
	// Parameters:
	//   vector: the translation vector for the X, Y, Z, and W axes
	// Returns: a new matrix representing the 4D translation transformation
	Translate4D(vector Vector4[T]) Matrix[T]

	// RotateX creates a rotation matrix around the X axis.
	// The resulting matrix is:
	// [ 1    0       0    0 ]
	// [ 0  cosθ   -sinθ  0 ]
	// [ 0  sinθ    cosθ  0 ]
	// [ 0    0       0   1 ]
	// Parameters:
	//   angle: the rotation angle in radians
	// Returns: a new matrix representing the rotation around the X axis
	RotateX(angle Radians) Matrix[T]

	// RotateY creates a rotation matrix around the Y axis.
	// The resulting matrix is:
	// [ cosθ  0  sinθ  0 ]
	// [   0   1   0    0 ]
	// [ -sinθ 0  cosθ  0 ]
	// [   0   0   0    1 ]
	// Parameters:
	//   angle: the rotation angle in radians
	// Returns: a new matrix representing the rotation around the Y axis
	RotateY(angle Radians) Matrix[T]

	// RotateZ creates a rotation matrix around the Z axis.
	// The resulting matrix is:
	// [ cosθ  -sinθ  0  0 ]
	// [ sinθ   cosθ  0  0 ]
	// [   0      0   1  0 ]
	// [   0      0   0  1 ]
	// Parameters:
	//   angle: the rotation angle in radians
	// Returns: a new matrix representing the rotation around the Z axis
	RotateZ(angle Radians) Matrix[T]

	// RotateAxisAngle creates a rotation matrix around an arbitrary axis using Rodrigues' formula.
	// Formula: R = I + sin(θ)K + (1-cos(θ))K^2, where K is the cross-product matrix of the axis.
	// Parameters:
	//   angle: the rotation angle in radians
	//   axis: the axis of rotation (must be a unit vector)
	// Returns: a new matrix representing the rotation around the given axis
	RotateAxisAngle(angle Radians, axis Vector3[T]) Matrix[T]

	// RotateQuaternion creates a rotation matrix from a quaternion.
	// See quaternion to matrix conversion formulas.
	// Parameters:
	//   quat: the rotation quaternion
	// Returns: a new matrix representing the rotation
	RotateQuaternion(quat Quaternion[T]) Matrix[T]

	// Orthographic creates an orthographic projection matrix for the given size.
	// The resulting matrix is typically:
	// [ 2/w   0    0   0 ]
	// [  0   2/h   0   0 ]
	// [  0    0   1/d  0 ]
	// [  0    0    0   1 ]
	// Parameters:
	//   size: the size of the view volume (width, height, depth)
	// Returns: a new matrix representing the orthographic projection
	Orthographic(size Size[T]) Matrix[T]

	// LookAt creates a view matrix for a camera at position looking at target with up vector.
	// The resulting matrix orients the camera in 3D space.
	// Parameters:
	//   position: the position of the camera
	//   target: the point the camera is looking at
	//   up: the up vector, which should be perpendicular to the view direction
	// Returns: a new matrix representing the view transformation
	LookAt(position, target, up Vector3[T]) Matrix[T]

	// === Vector Transformation ===

	// TransformPoint transforms a point using the matrix in homogeneous coordinates.
	// Formula: p' = M * p, where p is a 4D point (x, y, z, w).
	// Parameters:
	//   point: the point in homogeneous coordinates to transform
	// Returns: the transformed point in homogeneous coordinates
	TransformPoint(point Point[T]) Point[T]

	// TransformVector2D transforms a 2D direction vector (ignoring translation).
	// Formula: v' = M * v, where v is a direction vector (z=0, w=0).
	// Parameters:
	//   vector: the 2D direction vector to transform
	// Returns: the transformed 2D direction vector
	TransformVector2D(vector Vector2[T]) Vector2[T]

	// TransformVector3D transforms a 3D direction vector (ignoring translation).
	// Formula: v' = M * v, where v is a direction vector (w=0).
	// Parameters:
	//   vector: the 3D direction vector to transform
	// Returns: the transformed 3D direction vector
	TransformVector3D(vector Vector3[T]) Vector3[T]

	// TransformVector4D transforms a 4D direction vector (ignoring translation).
	// Formula: v' = M * v, where v is a direction vector.
	// Parameters:
	//   vector: the 4D direction vector to transform
	// Returns: the transformed 4D direction vector
	TransformVector4D(vector Vector4[T]) Vector4[T]

	// === Advanced Operations ===

	// Decompose returns the matrix decomposition into translation, scale, shear, perspective, and rotation comp1nts.
	// This is useful for extracting transformation parameters from a matrix.
	// Returns: a MatrixDecomposition struct containing the extracted comp1nts
	Decompose() MatrixDecomposition[T]

	// CosSin returns the cosine and sine of the given angle in radians.
	// This is a helper function for rotation operations.
	// Parameters:
	//   angle: the angle in radians
	// Returns: the cosine and sine of the angle
	CosSin(angle Radians) (cos, sin T)
}

// NewMatrix creates a new empty matrix with default values.
func NewMatrix[T Scalar]() Matrix[T] {
	return &matrix[T]{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

type matrix[T Scalar] [16]T

// At returns the value at the specified row and column (0-based).
func (m *matrix[T]) At(row int, col int) T {
	// Column-major order: index = col*4 + row
	return m[col*4+row]
}

// Set sets the value at the specified row and column (0-based).
func (m *matrix[T]) Set(row int, col int, value T) {
	m[col*4+row] = value
}

// Add implements Matrix.
func (m *matrix[T]) Add(other Matrix[T]) Matrix[T] {
	result := &matrix[T]{}
	for i := 0; i < 16; i++ {
		result[i] = m[i] + other.(*matrix[T])[i]
	}
	return result
}

// Sub implements Matrix.
func (m *matrix[T]) Sub(other Matrix[T]) Matrix[T] {
	result := &matrix[T]{}
	for i := 0; i < 16; i++ {
		result[i] = m[i] - other.(*matrix[T])[i]
	}
	return result
}

// Mul implements Matrix.
func (m *matrix[T]) Mul(other Matrix[T]) Matrix[T] {
	o := other.(*matrix[T])
	result := &matrix[T]{}

	// Matrix multiplication using column-major order
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			var sum T
			for k := 0; k < 4; k++ {
				sum += m[k*4+row] * o[col*4+k]
			}
			result[col*4+row] = sum
		}
	}
	return result
}

// DivElements implements Matrix.
func (m *matrix[T]) DivElements(other Matrix[T]) Matrix[T] {
	result := &matrix[T]{}
	o := other.(*matrix[T])
	for i := 0; i < 16; i++ {
		result[i] = m[i] / o[i]
	}
	return result
}

// Equal implements Matrix.
func (m *matrix[T]) Equal(other Matrix[T]) bool {
	o := other.(*matrix[T])
	for i := 0; i < 16; i++ {
		if !Equal(m[i], o[i]) {
			return false
		}
	}
	return true
}

// IsFinite implements Matrix.
func (m *matrix[T]) IsFinite() bool {
	for i := 0; i < 16; i++ {
		if !IsFinite(m[i]) {
			return false
		}
	}
	return true
}

// IsIdentity implements Matrix.
func (m *matrix[T]) IsIdentity() bool {

	return Equal(m[0], 1) && Equal(m[1], 0) && Equal(m[2], 0) && Equal(m[3], 0) &&
		Equal(m[4], 0) && Equal(m[5], 1) && Equal(m[6], 0) && Equal(m[7], 0) &&
		Equal(m[8], 0) && Equal(m[9], 0) && Equal(m[10], 1) && Equal(m[11], 0) &&
		Equal(m[12], 0) && Equal(m[13], 0) && Equal(m[14], 0) && Equal(m[15], 1)
}

// IsInvertible implements Matrix.
func (m *matrix[T]) IsInvertible() bool {

	return m.Determinant() != 0
}

// Determinant implements Matrix.
func (m *matrix[T]) Determinant() T {
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

// IsAffine implements Matrix.
func (m *matrix[T]) IsAffine() bool {
	return Equal(m[2], 0) && Equal(m[3], 0) && Equal(m[6], 0) && Equal(m[7], 0) &&
		Equal(m[8], 0) && Equal(m[9], 0) && Equal(m[10], 1) && Equal(m[11], 0) &&
		Equal(m[14], 0) && Equal(m[15], 1)
}

// HasPerspective implements Matrix.
func (m *matrix[T]) HasPerspective() bool {
	return !Equal(m[3], 0) || !Equal(m[7], 0) || !Equal(m[11], 0) || !Equal(m[15], 1)
}

// HasPerspective2D implements Matrix.
func (m *matrix[T]) HasPerspective2D() bool {
	return !Equal(m[3], 0) || !Equal(m[7], 0) || !Equal(m[15], 1)
}

// HasTranslation implements Matrix.
func (m *matrix[T]) HasTranslation() bool {
	return !Equal(m[12], 0) || !Equal(m[13], 0)
}

// IsAxisAligned implements Matrix.
func (m *matrix[T]) IsAxisAligned() bool {
	if m.HasPerspective() {
		return false
	}

	// Check if all three basis vectors are aligned to an axis
	v := [9]bool{
		!Equal(m[0], 0), !Equal(m[1], 0), !Equal(m[2], 0),
		!Equal(m[4], 0), !Equal(m[5], 0), !Equal(m[6], 0),
		!Equal(m[8], 0), !Equal(m[9], 0), !Equal(m[10], 0),
	}

	// Check if all three basis vectors are aligned to an axis
	if (boolToInt(v[0])+boolToInt(v[1])+boolToInt(v[2]) != 1) ||
		(boolToInt(v[3])+boolToInt(v[4])+boolToInt(v[5]) != 1) ||
		(boolToInt(v[6])+boolToInt(v[7])+boolToInt(v[8]) != 1) {
		return false
	}

	// Ensure that n1 of the basis vectors overlap
	if (boolToInt(v[0])+boolToInt(v[3])+boolToInt(v[6]) != 1) ||
		(boolToInt(v[1])+boolToInt(v[4])+boolToInt(v[7]) != 1) ||
		(boolToInt(v[2])+boolToInt(v[5])+boolToInt(v[8]) != 1) {
		return false
	}

	return true
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// IsAxisAligned2D implements Matrix.
func (m *matrix[T]) IsAxisAligned2D() bool {
	if m.HasPerspective2D() {
		return false
	}

	if Equal(m[1], 0) && Equal(m[4], 0) {
		return true
	}
	if Equal(m[0], 0) && Equal(m[5], 0) {
		return true
	}
	return false
}

// IsTranslationOnly implements Matrix.
func (m *matrix[T]) IsTranslationOnly() bool {

	return Equal(m[0], 1) && Equal(m[1], 0) && Equal(m[2], 0) && Equal(m[3], 0) &&
		Equal(m[4], 0) && Equal(m[5], 1) && Equal(m[6], 0) && Equal(m[7], 0) &&
		Equal(m[8], 0) && Equal(m[9], 0) && Equal(m[10], 1) && Equal(m[11], 0) &&
		Equal(m[15], 1)
}

// IsTranslationScaleOnly implements Matrix.
func (m *matrix[T]) IsTranslationScaleOnly() bool {

	return !Equal(m[0], 0) && Equal(m[1], 0) && Equal(m[2], 0) && Equal(m[3], 0) &&
		Equal(m[4], 0) && !Equal(m[5], 0) && Equal(m[6], 0) && Equal(m[7], 0) &&
		Equal(m[8], 0) && Equal(m[9], 0) && !Equal(m[10], 0) && Equal(m[11], 0) &&
		Equal(m[15], 1)
}

// Transpose implements Matrix.
func (m *matrix[T]) Transpose() Matrix[T] {
	result := &matrix[T]{}
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			result[row*4+col] = m[col*4+row]
		}
	}
	return result
}

// Inverse implements Matrix.
func (m *matrix[T]) Inverse() (Matrix[T], error) {
	// Using the same algorithm as C++ implementation
	tmp := &matrix[T]{
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
		return nil, errors.New("matrix is not invertible")
	}

	invDet := T(1) / det
	result := &matrix[T]{}
	for i := 0; i < 16; i++ {
		result[i] = tmp[i] * invDet
	}
	return result, nil
}

// To3x3 implements Matrix.
func (m *matrix[T]) To3x3() Matrix[T] {
	result := &matrix[T]{
		m[0], m[1], 0, m[3],
		m[4], m[5], 0, m[7],
		0, 0, 1, 0,
		m[12], m[13], 0, m[15],
	}
	return result
}

// ToColumnMajor implements Matrix.
func (m *matrix[T]) ToColumnMajor() Matrix[T] {
	// Already in column-major order
	result := &matrix[T]{}
	copy(result[:], m[:])
	return result
}

// ToRowMajor implements Matrix.
func (m *matrix[T]) ToRowMajor() Matrix[T] {
	result := &matrix[T]{}
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			result[row*4+col] = m[col*4+row]
		}
	}
	return result
}

// BasisVectors implements Matrix.
func (m *matrix[T]) BasisVectors() []Vector3[T] {
	return []Vector3[T]{
		NewVector3(m[0], m[1], m[2]),  // X basis
		NewVector3(m[4], m[5], m[6]),  // Y basis
		NewVector3(m[8], m[9], m[10]), // Z basis
	}
}

// BasisPoints implements Matrix.
func (m *matrix[T]) BasisPoints() []Point[T] {
	return []Point[T]{
		NewPoint(m[0], m[1]), // X basis
		NewPoint(m[4], m[5]), // Y basis
		NewPoint(m[8], m[9]), // Z basis (projected to 2D)
	}
}

// Scale implements Matrix.
func (m *matrix[T]) Scale() Vector3[T] {
	basisX := NewVector3(m[0], m[1], m[2])
	basisY := NewVector3(m[4], m[5], m[6])
	basisZ := NewVector3(m[8], m[9], m[10])
	return NewVector3(basisX.Length(), basisY.Length(), basisZ.Length())
}

// ScaleInDirection implements Matrix.
func (m *matrix[T]) ScaleInDirection(dir Vector3[T]) T {
	normalized := dir.Normalize()
	basis := m.To3x3()
	inv, err := basis.Inverse()
	if err != nil {
		return T(1) // fallback
	}
	transformed := inv.TransformVector3D(normalized)
	return T(1) / transformed.Length() * dir.Length()
}

// MaxScaleXY implements Matrix.
func (m *matrix[T]) MaxScaleXY() T {
	// Check for common case of axis-aligned transformation

	if Equal(m[1], 0) && Equal(m[4], 0) {
		absX := m[0]
		if absX < 0 {
			absX = -absX
		}
		absY := m[5]
		if absY < 0 {
			absY = -absY
		}
		if absX > absY {
			return absX
		}
		return absY
	}

	// General case
	scaleX := T(math.Sqrt(float64(m[0]*m[0] + m[1]*m[1])))
	scaleY := T(math.Sqrt(float64(m[4]*m[4] + m[5]*m[5])))
	if scaleX > scaleY {
		return scaleX
	}
	return scaleY
}

// Scale2D implements Matrix.
func (m *matrix[T]) Scale2D(scale Vector2[T]) Matrix[T] {
	return &matrix[T]{
		scale.X(), 0, 0, 0,
		0, scale.Y(), 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Scale3D implements Matrix.
func (m *matrix[T]) Scale3D(scale Vector3[T]) Matrix[T] {
	return &matrix[T]{
		scale.X(), 0, 0, 0,
		0, scale.Y(), 0, 0,
		0, 0, scale.Z(), 0,
		0, 0, 0, 1,
	}
}

// Translate2D implements Matrix.
func (m *matrix[T]) Translate2D(vector Vector2[T]) Matrix[T] {
	return &matrix[T]{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		vector.X(), vector.Y(), 0, 1,
	}
}

// Translate3D implements Matrix.
func (m *matrix[T]) Translate3D(vector Vector3[T]) Matrix[T] {
	return &matrix[T]{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		vector.X(), vector.Y(), vector.Z(), 1,
	}
}

// Translate4D implements Matrix.
func (m *matrix[T]) Translate4D(vector Vector4[T]) Matrix[T] {
	return &matrix[T]{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		vector.X(), vector.Y(), vector.Z(), vector.W(),
	}
}

// RotateX implements Matrix.
func (m *matrix[T]) RotateX(angle Radians) Matrix[T] {
	cos, sin := m.CosSin(angle)
	return &matrix[T]{
		1, 0, 0, 0,
		0, cos, sin, 0,
		0, -sin, cos, 0,
		0, 0, 0, 1,
	}
}

// RotateY implements Matrix.
func (m *matrix[T]) RotateY(angle Radians) Matrix[T] {
	cos, sin := m.CosSin(angle)
	return &matrix[T]{
		cos, 0, -sin, 0,
		0, 1, 0, 0,
		sin, 0, cos, 0,
		0, 0, 0, 1,
	}
}

// RotateZ implements Matrix.
func (m *matrix[T]) RotateZ(angle Radians) Matrix[T] {
	cos, sin := m.CosSin(angle)
	return &matrix[T]{
		cos, sin, 0, 0,
		-sin, cos, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// RotateAxisAngle implements Matrix.
func (m *matrix[T]) RotateAxisAngle(angle Radians, axis Vector3[T]) Matrix[T] {
	v := axis.Normalize()
	cos, sin := m.CosSin(angle)
	cosp := T(1) - cos

	return &matrix[T]{
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
}

// RotateQuaternion implements Matrix.
func (m *matrix[T]) RotateQuaternion(quat Quaternion[T]) Matrix[T] {
	x, y, z, w := quat.X(), quat.Y(), quat.Z(), quat.W()
	return &matrix[T]{
		T(1) - T(2)*(y*y+z*z), T(2) * (x*y + z*w), T(2) * (x*z - y*w), 0,
		T(2) * (x*y - z*w), T(1) - T(2)*(x*x+z*z), T(2) * (y*z + x*w), 0,
		T(2) * (x*z + y*w), T(2) * (y*z - x*w), T(1) - T(2)*(x*x+y*y), 0,
		0, 0, 0, 1,
	}
}

// Orthographic implements Matrix.
func (m *matrix[T]) Orthographic(size Size[T]) Matrix[T] {
	// Per NDC assumptions: scale and translate to NDC space
	scaleX := T(2) / size.Width()
	scaleY := -T(2) / size.Height()

	return &matrix[T]{
		scaleX, 0, 0, 0,
		0, scaleY, 0, 0,
		0, 0, 0, 0,
		-1, 1, 1 / 2, 1,
	}
}

// LookAt implements Matrix.
func (m *matrix[T]) LookAt(position, target, up Vector3[T]) Matrix[T] {
	forward := target.Sub(position).Normalize()
	right := up.Cross(forward)
	upNorm := forward.Cross(right)

	return &matrix[T]{
		right.X(), upNorm.X(), forward.X(), 0,
		right.Y(), upNorm.Y(), forward.Y(), 0,
		right.Z(), upNorm.Z(), forward.Z(), 0,
		-right.Dot(position), -upNorm.Dot(position), -forward.Dot(position), 1,
	}
}

// TransformPoint implements Matrix.
func (m *matrix[T]) TransformPoint(point Point[T]) Point[T] {
	x, y := point.X(), point.Y()
	w := x*m[3] + y*m[7] + m[15]
	resultX := x*m[0] + y*m[4] + m[12]
	resultY := x*m[1] + y*m[5] + m[13]

	if w != 0 {
		w = T(1) / w
	}
	return NewPoint(resultX*w, resultY*w)
}

// TransformVector2D implements Matrix.
func (m *matrix[T]) TransformVector2D(vector Vector2[T]) Vector2[T] {
	x, y := vector.X(), vector.Y()
	return NewVector2(
		x*m[0]+y*m[4],
		x*m[1]+y*m[5],
	)
}

// TransformVector3D implements Matrix.
func (m *matrix[T]) TransformVector3D(vector Vector3[T]) Vector3[T] {
	x, y, z := vector.X(), vector.Y(), vector.Z()
	return NewVector3(
		x*m[0]+y*m[4]+z*m[8],
		x*m[1]+y*m[5]+z*m[9],
		x*m[2]+y*m[6]+z*m[10],
	)
}

// TransformVector4D implements Matrix.
func (m *matrix[T]) TransformVector4D(vector Vector4[T]) Vector4[T] {
	x, y, z, w := vector.X(), vector.Y(), vector.Z(), vector.W()
	return NewVector4(
		x*m[0]+y*m[4]+z*m[8]+w*m[12],
		x*m[1]+y*m[5]+z*m[9]+w*m[13],
		x*m[2]+y*m[6]+z*m[10]+w*m[14],
		x*m[3]+y*m[7]+z*m[11]+w*m[15],
	)
}

// Decompose implements Matrix.
func (m *matrix[T]) Decompose() MatrixDecomposition[T] {
	// This is a simplified version - full decomposition is complex
	// For now, return basic comp1nts
	return MatrixDecomposition[T]{
		Translation: NewVector3(m[12], m[13], m[14]),
		Scale:       m.Scale(),
		Shear:       Shear[T]{XY: 0, XZ: 0, YZ: 0},
		Perspective: NewVector4(T(0), T(0), T(0), T(1)),
		Rotation:    NewQuaternion(T(0), T(0), T(0), T(1)),
	}
}

// CosSin implements Matrix.
func (m *matrix[T]) CosSin(angle Radians) (cos, sin T) {
	sinVal := T(math.Sin(float64(angle)))
	if math.Abs(float64(sinVal)) == 1.0 {
		// 90 or 270 degrees
		return T(0), sinVal
	}

	cosVal := T(math.Cos(float64(angle)))
	if math.Abs(float64(cosVal)) == 1.0 {
		// 0 or 180 degrees
		return cosVal, T(0)
	}

	return cosVal, sinVal
}
