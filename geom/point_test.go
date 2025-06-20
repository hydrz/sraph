package geom

import (
	"math"
	"testing"
)

func TestPoint_Point(t *testing.T) {
	p := Point[Scalar]{X: 3.0, Y: 4.0}
	if p.X != 3.0 || p.Y != 4.0 {
		t.Errorf("Point(3.0, 4.0) = (%v, %v), want (3.0, 4.0)", p.X, p.Y)
	}
}

func TestPoint_PointArithmetic(t *testing.T) {
	p1 := Point[Scalar]{X: 1.0, Y: 2.0}
	p2 := Point[Scalar]{X: 3.0, Y: 4.0}

	t.Run("Add", func(t *testing.T) {
		result := p1.Add(p2)
		if result.X != 4.0 || result.Y != 6.0 {
			t.Errorf("Add() = (%v, %v), want (4.0, 6.0)", result.X, result.Y)
		}
	})

	t.Run("Sub", func(t *testing.T) {
		result := p2.Sub(p1)
		if result.X != 2.0 || result.Y != 2.0 {
			t.Errorf("Sub() = (%v, %v), want (2.0, 2.0)", result.X, result.Y)
		}
	})

	t.Run("Mul", func(t *testing.T) {
		result := p1.Mul(p2)
		if result.X != 3.0 || result.Y != 8.0 {
			t.Errorf("Mul() = (%v, %v), want (3.0, 8.0)", result.X, result.Y)
		}
	})

	t.Run("Div", func(t *testing.T) {
		result := p2.Div(p1)
		if result.X != 3.0 || result.Y != 2.0 {
			t.Errorf("Div() = (%v, %v), want (3.0, 2.0)", result.X, result.Y)
		}
	})
}

func TestPoint_PointTransformations(t *testing.T) {
	t.Run("Abs", func(t *testing.T) {
		negP := Point[Scalar]{X: -3.0, Y: -4.0}
		result := negP.Abs()
		if result.X != 3.0 || result.Y != 4.0 {
			t.Errorf("Abs() = (%v, %v), want (3.0, 4.0)", result.X, result.Y)
		}
	})

	t.Run("Scale", func(t *testing.T) {
		p := Point[Scalar]{X: 3.0, Y: 4.0} // Define p for this test
		result := p.Scale(2.0)
		if result.X != 6.0 || result.Y != 8.0 {
			t.Errorf("Scale(2.0) = (%v, %v), want (6.0, 8.0)", result.X, result.Y)
		}
	})
}

func TestPoint_PointMath(t *testing.T) {
	tests := []struct {
		name     string
		point    Point[Scalar]
		expected Point[Scalar]
		testFunc func(Point[Scalar]) Point[Scalar]
	}{
		{
			"Floor",
			Point[Scalar]{X: 3.7, Y: 4.2},
			Point[Scalar]{X: 3.0, Y: 4.0},
			func(p Point[Scalar]) Point[Scalar] { return p.Floor() },
		},
		{
			"Ceil",
			Point[Scalar]{X: 3.2, Y: 4.7},
			Point[Scalar]{X: 4.0, Y: 5.0},
			func(p Point[Scalar]) Point[Scalar] { return p.Ceil() },
		},
		{
			"Round",
			Point[Scalar]{X: 3.4, Y: 4.6},
			Point[Scalar]{X: 3.0, Y: 5.0},
			func(p Point[Scalar]) Point[Scalar] { return p.Round() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc(tt.point)
			if !NearlyEqual(result.X, tt.expected.X) || !NearlyEqual(result.Y, tt.expected.Y) {
				t.Errorf("%s() = (%v, %v), want (%v, %v)", tt.name, result.X, result.Y, tt.expected.X, tt.expected.Y)
			}
		})
	}
}

func TestPoint_PointComparison(t *testing.T) {
	p1 := Point[Scalar]{X: 1.0, Y: 2.0}
	p2 := Point[Scalar]{X: 1.0, Y: 2.0}
	p3 := Point[Scalar]{X: 3.0, Y: 4.0}

	t.Run("Eq", func(t *testing.T) {
		if !p1.Equal(p2) {
			t.Error("Eq points should be Eq")
		}
		if p1.Equal(p3) {
			t.Error("Different points should not be Eq")
		}
	})
}

func TestPoint_PointProperties(t *testing.T) {
	t.Run("IsZero", func(t *testing.T) {
		zero := Point[Scalar]{X: 0.0, Y: 0.0}
		nonZero := Point[Scalar]{X: 1.0, Y: 0.0}

		if !zero.IsZero() {
			t.Error("Zero point should return true for IsZero()")
		}
		if nonZero.IsZero() {
			t.Error("Non-zero point should return false for IsZero()")
		}
	})

	t.Run("IsFinite", func(t *testing.T) {
		finite := Point[Scalar]{X: 1.0, Y: 2.0}
		infinite := Point[Scalar]{X: Scalar(math.Inf(1)), Y: 2.0}

		if !finite.IsFinite() {
			t.Error("Finite point should return true for IsFinite()")
		}
		if infinite.IsFinite() {
			t.Error("Infinite point should return false for IsFinite()")
		}
	})
}

func TestPoint_PointVectorOperations(t *testing.T) {
	p1 := Point[Scalar]{X: 3.0, Y: 4.0}
	p2 := Point[Scalar]{X: 1.0, Y: 2.0}

	t.Run("Length", func(t *testing.T) {
		length := p1.Length()
		expected := Scalar(5.0) // sqrt(3^2 + 4^2) = 5
		if !NearlyEqual(length, expected) {
			t.Errorf("Length() = %v, want %v", length, expected)
		}
	})

	t.Run("LengthSquared", func(t *testing.T) {
		lengthSq := p1.LengthSquared()
		expected := Scalar(25.0) // 3^2 + 4^2 = 25
		if !NearlyEqual(lengthSq, expected) {
			t.Errorf("LengthSquared() = %v, want %v", lengthSq, expected)
		}
	})

	t.Run("Distance", func(t *testing.T) {
		distance := p1.Distance(p2)
		expected := Scalar(math.Sqrt(8.0)) // sqrt((3-1)^2 + (4-2)^2) = sqrt(8)
		if !NearlyEqual(distance, expected) {
			t.Errorf("Distance() = %v, want %v", distance, expected)
		}
	})

	t.Run("DistanceSquared", func(t *testing.T) {
		distanceSq := p1.DistanceSquared(p2)
		expected := Scalar(8.0) // (3-1)^2 + (4-2)^2 = 8
		if !NearlyEqual(distanceSq, expected) {
			t.Errorf("DistanceSquared() = %v, want %v", distanceSq, expected)
		}
	})

	t.Run("Dot", func(t *testing.T) {
		dot := p1.Dot(p2)
		expected := Scalar(11.0) // 3*1 + 4*2 = 11
		if !NearlyEqual(dot, expected) {
			t.Errorf("Dot() = %v, want %v", dot, expected)
		}
	})

	t.Run("Cross", func(t *testing.T) {
		cross := p1.Cross(p2)
		expected := Scalar(2.0) // 3*2 - 4*1 = 2
		if !NearlyEqual(cross, expected) {
			t.Errorf("Cross() = %v, want %v", cross, expected)
		}
	})
}

func TestPoint_PointNormalize(t *testing.T) {
	t.Run("Normal vector", func(t *testing.T) {
		p := Point[Scalar]{X: 3.0, Y: 4.0}
		normalized := p.Normalize()

		// Should have length 1
		length := normalized.Length()
		if !NearlyEqual(length, Scalar(1.0)) {
			t.Errorf("Normalized point length = %v, want 1.0", length)
		}
	})

	t.Run("Zero vector", func(t *testing.T) {
		zero := Point[Scalar]{X: 0.0, Y: 0.0}
		normalized := zero.Normalize()

		// Should return default unit vector (1, 0)
		if !NearlyEqual(normalized.X, Scalar(1.0)) || !NearlyEqual(normalized.Y, Scalar(0.0)) {
			t.Errorf("Normalized zero point = (%v, %v), want (1.0, 0.0)", normalized.X, normalized.Y)
		}
	})
}

func TestPoint_PointRotate(t *testing.T) {
	p := Point[Scalar]{X: 1.0, Y: 0.0}

	t.Run("90 degrees", func(t *testing.T) {
		rotated := p.Rotate(Radians(math.Pi / 2))
		// Should be approximately (0, 1)
		if !NearlyEqual(rotated.X, Scalar(0.0)) || !NearlyEqual(rotated.Y, Scalar(1.0)) {
			t.Errorf("Rotate(π/2) = (%v, %v), want (0.0, 1.0)", rotated.X, rotated.Y)
		}
	})

	t.Run("180 degrees", func(t *testing.T) {
		rotated := p.Rotate(Radians(math.Pi))
		// Should be approximately (-1, 0)
		if !NearlyEqual(rotated.X, Scalar(-1.0)) || !NearlyEqual(rotated.Y, Scalar(0.0)) {
			t.Errorf("Rotate(π) = (%v, %v), want (-1.0, 0.0)", rotated.X, rotated.Y)
		}
	})
}

func TestPoint_PointReflect(t *testing.T) {
	tests := []struct {
		axis     Vector2[Scalar]
		point    Point[Scalar]
		expected Point[Scalar]
	}{
		{
			axis:     Vector2[Scalar]{X: 0, Y: 1},
			point:    Point[Scalar]{X: 2, Y: 3},
			expected: Point[Scalar]{X: 2, Y: -3},
		},
		{
			axis:     Vector2[Scalar]{X: 1, Y: 1},
			point:    Point[Scalar]{X: 1, Y: 0},
			expected: Point[Scalar]{X: 0, Y: -1},
		},
		{
			axis:     Vector2[Scalar]{X: 1, Y: 1},
			point:    Point[Scalar]{X: -1, Y: -1},
			expected: Point[Scalar]{X: -1, Y: -1}.Neg(),
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			reflected := tt.point.Reflect(tt.axis)
			if !reflected.Equal(tt.expected) {
				t.Errorf("Reflect(%v) = %v, want %v", tt.axis, reflected, tt.expected)
			}
		})
	}
}

func TestPoint_PointLerp(t *testing.T) {
	p1 := Point[Scalar]{X: 0.0, Y: 0.0}
	p2 := Point[Scalar]{X: 10.0, Y: 20.0}

	tests := []struct {
		t        Scalar
		expected Point[Scalar]
	}{
		{Scalar(0.0), p1},
		{Scalar(1.0), p2},
		{Scalar(0.5), Point[Scalar]{X: 5.0, Y: 10.0}},
		{Scalar(0.25), Point[Scalar]{X: 2.5, Y: 5.0}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := p1.Lerp(p2, tt.t)
			if !NearlyEqual(result.X, tt.expected.X) || !NearlyEqual(result.Y, tt.expected.Y) {
				t.Errorf("Lerp(%v) = (%v, %v), want (%v, %v)", tt.t, result.X, result.Y, tt.expected.X, tt.expected.Y)
			}
		})
	}
}

func TestPoint_PointTranslate(t *testing.T) {
	p := Point[Scalar]{X: 1.0, Y: 2.0}
	vector := Vector2[Scalar]{X: 3.0, Y: 4.0}

	result := p.Translate(vector)
	if result.X != 4.0 || result.Y != 6.0 {
		t.Errorf("Translate() = (%v, %v), want (4.0, 6.0)", result.X, result.Y)
	}
}
