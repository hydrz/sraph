package geom

import (
	"math"
	"testing"
)

func TestPoint_Point(t *testing.T) {
	p := NewPoint[Scalar](3.0, 4.0)
	if p.X() != 3.0 || p.Y() != 4.0 {
		t.Errorf("Point: expected (3,4), got (%v,%v)", p.X(), p.Y())
	}
}

func TestPoint_PointArithmetic(t *testing.T) {
	p1 := NewPoint[Scalar](1.0, 2.0)
	p2 := NewPoint[Scalar](3.0, 4.0)

	t.Run("Add", func(t *testing.T) {
		res := p1.Add(p2)
		if res.X() != 4.0 || res.Y() != 6.0 {
			t.Errorf("Add: expected (4,6), got (%v,%v)", res.X(), res.Y())
		}
	})

	t.Run("Sub", func(t *testing.T) {
		res := p2.Sub(p1)
		if res.X() != 2.0 || res.Y() != 2.0 {
			t.Errorf("Sub: expected (2,2), got (%v,%v)", res.X(), res.Y())
		}
	})

	t.Run("Mul", func(t *testing.T) {
		res := p1.Mul(p2)
		if res.X() != 3.0 || res.Y() != 8.0 {
			t.Errorf("Mul: expected (3,8), got (%v,%v)", res.X(), res.Y())
		}
	})

	t.Run("Div", func(t *testing.T) {
		res := p2.Div(p1)
		if res.X() != 3.0 || res.Y() != 2.0 {
			t.Errorf("Div: expected (3,2), got (%v,%v)", res.X(), res.Y())
		}
	})
}

func TestPoint_PointTransformations(t *testing.T) {
	t.Run("Abs", func(t *testing.T) {
		negP := NewPoint[Scalar](-3.0, -4.0)
		res := negP.Abs()
		if res.X() != 3.0 || res.Y() != 4.0 {
			t.Errorf("Abs: expected (3,4), got (%v,%v)", res.X(), res.Y())
		}
	})

	t.Run("Scale", func(t *testing.T) {
		p := NewPoint[Scalar](3.0, 4.0)
		res := p.Scale(2.0)
		if res.X() != 6.0 || res.Y() != 8.0 {
			t.Errorf("Scale(2.0): expected (6,8), got (%v,%v)", res.X(), res.Y())
		}
	})
}

func TestPoint_PointMath(t *testing.T) {
	tests := []struct {
		name     string
		point    Point
		expected Point
		testFunc func(Point) Point
	}{
		{
			"Floor",
			NewPoint[Scalar](3.7, 4.2),
			NewPoint[Scalar](3.0, 4.0),
			func(p Point) Point { return p.Floor() },
		},
		{
			"Ceil",
			NewPoint[Scalar](3.2, 4.7),
			NewPoint[Scalar](4.0, 5.0),
			func(p Point) Point { return p.Ceil() },
		},
		{
			"Round",
			NewPoint[Scalar](3.4, 4.6),
			NewPoint[Scalar](3.0, 5.0),
			func(p Point) Point { return p.Round() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.testFunc(tt.point)
			if !NearlyEqual(res.X(), tt.expected.X()) || !NearlyEqual(res.Y(), tt.expected.Y()) {
				t.Errorf("%s: expected (%v,%v), got (%v,%v)", tt.name, tt.expected.X(), tt.expected.Y(), res.X(), res.Y())
			}
		})
	}
}

func TestPoint_PointComparison(t *testing.T) {
	p1 := NewPoint(1.0, 2.0)
	p2 := NewPoint(1.0, 2.0)
	p3 := NewPoint(3.0, 4.0)

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
		zero := NewPoint(0.0, 0.0)
		nonZero := NewPoint(1.0, 0.0)

		if !zero.IsZero() {
			t.Error("Zero point should return true for IsZero()")
		}
		if nonZero.IsZero() {
			t.Error("Non-zero point should return false for IsZero()")
		}
	})

	t.Run("IsFinite", func(t *testing.T) {
		finite := NewPoint(1.0, 2.0)
		infinite := NewPoint(Scalar(math.Inf(1)), 2.0)

		if !finite.IsFinite() {
			t.Error("Finite point should return true for IsFinite()")
		}
		if infinite.IsFinite() {
			t.Error("Infinite point should return false for IsFinite()")
		}
	})
}

func TestPoint_PointVectorOperations(t *testing.T) {
	p1 := NewPoint(3.0, 4.0)
	p2 := NewPoint(1.0, 2.0)

	t.Run("Length", func(t *testing.T) {
		length := p1.Length()
		expected := Scalar(5) // sqrt(3^2 + 4^2) = 5
		if !NearlyEqual(length, expected) {
			t.Errorf("Length: expected %v, got %v", expected, length)
		}
	})

	t.Run("LengthSquared", func(t *testing.T) {
		lengthSq := p1.LengthSquared()
		expected := Scalar(25.0) // 3^2 + 4^2 = 25
		if !NearlyEqual(lengthSq, expected) {
			t.Errorf("LengthSquared: expected %v, got %v", expected, lengthSq)
		}
	})

	t.Run("Distance", func(t *testing.T) {
		distance := p1.Distance(p2)
		expected := Scalar(math.Sqrt(8.0)) // sqrt((3-1)^2 + (4-2)^2) = sqrt(8)
		if !NearlyEqual(distance, expected) {
			t.Errorf("Distance: expected %v, got %v", expected, distance)
		}
	})

	t.Run("DistanceSquared", func(t *testing.T) {
		distanceSq := p1.DistanceSquared(p2)
		expected := 8.0 // (3-1)^2 + (4-2)^2 = 8
		if !NearlyEqual(distanceSq, expected) {
			t.Errorf("DistanceSquared: expected %v, got %v", expected, distanceSq)
		}
	})

	t.Run("Dot", func(t *testing.T) {
		dot := p1.Dot(p2)
		expected := 11.0 // 3*1 + 4*2 = 11
		if !NearlyEqual(dot, expected) {
			t.Errorf("Dot: expected %v, got %v", expected, dot)
		}
	})

	t.Run("Cross", func(t *testing.T) {
		cross := p1.Cross(p2)
		expected := 2.0 // 3*2 - 4*1 = 2
		if !NearlyEqual(cross, expected) {
			t.Errorf("Cross: expected %v, got %v", expected, cross)
		}
	})
}

func TestPoint_PointNormalize(t *testing.T) {
	t.Run("Normal vector", func(t *testing.T) {
		p := NewPoint(3.0, 4.0)
		normalized := p.Normalize()

		// Should have length 1
		length := normalized.Length()
		if !NearlyEqual(length, 1.0) {
			t.Errorf("Normalized point length = %v, want 1.0", length)
		}
	})

	t.Run("Zero vector", func(t *testing.T) {
		zero := NewPoint(0.0, 0.0)
		normalized := zero.Normalize()

		// Should return default unit vector (1, 0)
		if !NearlyEqual(normalized.X(), 1.0) || !NearlyEqual(normalized.Y(), 0.0) {
			t.Errorf("Normalized zero point = (%v, %v), want (1.0, 0.0)", normalized.X(), normalized.Y())
		}
	})
}

func TestPoint_PointRotate(t *testing.T) {
	p := NewPoint(1.0, 0.0)

	t.Run("90 degrees", func(t *testing.T) {
		rotated := p.Rotate(Radians(math.Pi / 2))
		// Should be approximately (0, 1)
		if !NearlyEqual(rotated.X(), 0.0) || !NearlyEqual(rotated.Y(), 1.0) {
			t.Errorf("Rotate(π/2): expected (%v,%v), got (%v,%v)", 0.0, 1.0, rotated.X(), rotated.Y())
		}
	})

	t.Run("180 degrees", func(t *testing.T) {
		rotated := p.Rotate(Radians(math.Pi))
		// Should be approximately (-1, 0)
		if !NearlyEqual(rotated.X(), Scalar(-1.0)) || !NearlyEqual(rotated.Y(), 0.0) {
			t.Errorf("Rotate(π): expected (%v,%v), got (%v,%v)", Scalar(-1.0), 0.0, rotated.X(), rotated.Y())
		}
	})
}

func TestPoint_PointReflect(t *testing.T) {
	tests := []struct {
		axis     Vector2
		point    Point
		expected Point
	}{
		{
			axis:     NewVector2(0, 1),
			point:    NewPoint(2, 3),
			expected: NewPoint(2, -3),
		},
		{
			axis:     NewVector2(1, 1).Normalize(),
			point:    NewPoint(1, 0),
			expected: NewPoint(0, -1),
		},
		{
			axis:     NewVector2(1, 1).Normalize(),
			point:    NewPoint(-1, -1),
			expected: NewPoint(1, 1), // Reflecting across the normalized axis should yield the negative of the original point
		},
	}
	for _, tt := range tests {
		res := tt.point.Reflect(tt.axis)
		if !res.Equal(tt.expected) {
			t.Errorf("Reflect: expected %v, got %v", tt.expected, res)
		}
	}
}

func TestPoint_PointLerp(t *testing.T) {
	p1 := NewPoint(0.0, 0.0)
	p2 := NewPoint(10.0, 20.0)

	tests := []struct {
		t        Scalar
		expected Point
	}{
		{0.0, p1},
		{1.0, p2},
		{0.5, NewPoint(5.0, 10.0)},
		{0.25, NewPoint(2.5, 5.0)},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := p1.Lerp(p2, tt.t)
			if !NearlyEqual(result.X(), tt.expected.X()) || !NearlyEqual(result.Y(), tt.expected.Y()) {
				t.Errorf("Lerp(%v): expected (%v,%v), got (%v,%v)", tt.t, tt.expected.X(), tt.expected.Y(), result.X(), result.Y())
			}
		})
	}
}

func TestPoint_PointTranslate(t *testing.T) {
	p := NewPoint(1.0, 2.0)
	vector := NewVector2(3.0, 4.0)

	result := p.Translate(vector)
	if result.X() != 4.0 || result.Y() != 6.0 {
		t.Errorf("Translate: expected (4,6), got (%v,%v)", result.X(), result.Y())
	}
}
