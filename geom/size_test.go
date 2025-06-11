package geom

import (
	"math"
	"testing"
)

func TestSize_NewSize(t *testing.T) {
	s := NewSize(Float32(10.0), Float32(20.0))
	if s.Width() != 10.0 || s.Height() != 20.0 {
		t.Errorf("NewSize(10.0, 20.0) = (%v, %v), want (10.0, 20.0)", s.Width(), s.Height())
	}
}

func TestSize_NewSizeInfinite(t *testing.T) {
	s := NewSizeInfinite[Float32]()
	maxVal := Max[Float32]()
	if s.Width() != maxVal || s.Height() != maxVal {
		t.Errorf("NewSizeInfinite() should return max values")
	}
}

func TestSize_SizeArithmetic(t *testing.T) {
	s1 := NewSize(Float32(10.0), Float32(20.0))
	s2 := NewSize(Float32(5.0), Float32(10.0))

	t.Run("Add", func(t *testing.T) {
		result := s1.Add(s2)
		if result.Width() != 15.0 || result.Height() != 30.0 {
			t.Errorf("Add() = (%v, %v), want (15.0, 30.0)", result.Width(), result.Height())
		}
	})

	t.Run("Sub", func(t *testing.T) {
		result := s1.Sub(s2)
		if result.Width() != 5.0 || result.Height() != 10.0 {
			t.Errorf("Sub() = (%v, %v), want (5.0, 10.0)", result.Width(), result.Height())
		}
	})

	t.Run("Mul", func(t *testing.T) {
		result := s1.Mul(s2)
		if result.Width() != 50.0 || result.Height() != 200.0 {
			t.Errorf("Mul() = (%v, %v), want (50.0, 200.0)", result.Width(), result.Height())
		}
	})

	t.Run("Div", func(t *testing.T) {
		result := s1.Div(s2)
		if result.Width() != 2.0 || result.Height() != 2.0 {
			t.Errorf("Div() = (%v, %v), want (2.0, 2.0)", result.Width(), result.Height())
		}
	})

	t.Run("Neg", func(t *testing.T) {
		result := s1.Neg()
		if result.Width() != -10.0 || result.Height() != -20.0 {
			t.Errorf("Neg() = (%v, %v), want (-10.0, -20.0)", result.Width(), result.Height())
		}
	})
}

func TestSize_SizeScale(t *testing.T) {
	s := NewSize(Float32(10.0), Float32(20.0))

	t.Run("Scale", func(t *testing.T) {
		result := s.ScaleWH(Float32(2.0), Float32(1.5))
		if result.Width() != 20.0 || result.Height() != 30.0 {
			t.Errorf("Scale(2.0, 1.5) = (%v, %v), want (20.0, 30.0)", result.Width(), result.Height())
		}
	})

	t.Run("ScaleDim", func(t *testing.T) {
		result := s.Scale(Float32(2.0))
		if result.Width() != 20.0 || result.Height() != 40.0 {
			t.Errorf("ScaleDim(2.0) = (%v, %v), want (20.0, 40.0)", result.Width(), result.Height())
		}
	})
}

func TestSize_SizeComparison(t *testing.T) {
	s1 := NewSize(Float32(10.0), Float32(20.0))
	s2 := NewSize(Float32(10.0), Float32(20.0))
	s3 := NewSize(Float32(5.0), Float32(15.0))

	t.Run("Equal", func(t *testing.T) {
		if !s1.Equal(s2) {
			t.Error("Equal sizes should be equal")
		}
		if s1.Equal(s3) {
			t.Error("Different sizes should not be equal")
		}
	})

	t.Run("Min", func(t *testing.T) {
		result := s1.Min(s3)
		if result.Width() != 5.0 || result.Height() != 15.0 {
			t.Errorf("Min() = (%v, %v), want (5.0, 15.0)", result.Width(), result.Height())
		}
	})

	t.Run("Max", func(t *testing.T) {
		result := s1.Max(s3)
		if result.Width() != 10.0 || result.Height() != 20.0 {
			t.Errorf("Max() = (%v, %v), want (10.0, 20.0)", result.Width(), result.Height())
		}
	})
}

func TestSize_SizeDimensions(t *testing.T) {
	s := NewSize(Float32(15.0), Float32(10.0))

	t.Run("MinDimension", func(t *testing.T) {
		min := s.MinDimension()
		if min != 10.0 {
			t.Errorf("MinDimension() = %v, want 10.0", min)
		}
	})

	t.Run("MaxDimension", func(t *testing.T) {
		max := s.MaxDimension()
		if max != 15.0 {
			t.Errorf("MaxDimension() = %v, want 15.0", max)
		}
	})

	t.Run("Area", func(t *testing.T) {
		area := s.Area()
		if area != 150.0 {
			t.Errorf("Area() = %v, want 150.0", area)
		}
	})
}

func TestSize_SizeAbs(t *testing.T) {
	s := NewSize(Float32(-10.0), Float32(-20.0))
	result := s.Abs()

	if result.Width() != 10.0 || result.Height() != 20.0 {
		t.Errorf("Abs() = (%v, %v), want (10.0, 20.0)", result.Width(), result.Height())
	}
}

func TestSize_SizeMathFunctions(t *testing.T) {
	s := NewSize(Float32(10.7), Float32(20.3))

	t.Run("Floor", func(t *testing.T) {
		result := s.Floor()
		if result.Width() != 10.0 || result.Height() != 20.0 {
			t.Errorf("Floor() = (%v, %v), want (10.0, 20.0)", result.Width(), result.Height())
		}
	})

	t.Run("Ceil", func(t *testing.T) {
		result := s.Ceil()
		if result.Width() != 11.0 || result.Height() != 21.0 {
			t.Errorf("Ceil() = (%v, %v), want (11.0, 21.0)", result.Width(), result.Height())
		}
	})

	t.Run("Round", func(t *testing.T) {
		result := s.Round()
		if result.Width() != 11.0 || result.Height() != 20.0 {
			t.Errorf("Round() = (%v, %v), want (11.0, 20.0)", result.Width(), result.Height())
		}
	})
}

func TestSize_SizeProperties(t *testing.T) {
	t.Run("IsZero", func(t *testing.T) {
		zero := NewSize(Float32(0.0), Float32(0.0))
		nonZero := NewSize(Float32(1.0), Float32(0.0))

		if !zero.IsZero() {
			t.Error("Zero size should return true for IsZero()")
		}
		if nonZero.IsZero() {
			t.Error("Non-zero size should return false for IsZero()")
		}
	})

	t.Run("IsFinite", func(t *testing.T) {
		finite := NewSize(Float32(10.0), Float32(20.0))
		infinite := NewSize(Float32(math.Inf(1)), Float32(20.0))

		if !finite.IsFinite() {
			t.Error("Finite size should return true for IsFinite()")
		}
		if infinite.IsFinite() {
			t.Error("Infinite size should return false for IsFinite()")
		}
	})

	t.Run("IsInfinite", func(t *testing.T) {
		finite := NewSize(Float32(10.0), Float32(20.0))
		infinite := NewSize(Float32(math.Inf(1)), Float32(20.0))

		if finite.IsInfinite() {
			t.Error("Finite size should return false for IsInfinite()")
		}
		if !infinite.IsInfinite() {
			t.Error("Infinite size should return true for IsInfinite()")
		}
	})

	t.Run("IsSquare", func(t *testing.T) {
		square := NewSize(Float32(10.0), Float32(10.0))
		rectangle := NewSize(Float32(10.0), Float32(20.0))

		if !square.IsSquare() {
			t.Error("Square size should return true for IsSquare()")
		}
		if rectangle.IsSquare() {
			t.Error("Rectangle size should return false for IsSquare()")
		}
	})
}

func TestSize_SizeMipCount(t *testing.T) {
	tests := []struct {
		name     string
		size     Size[Int32]
		expected int
	}{
		{"1x1", NewSize(Int32(1), Int32(1)), 1},
		{"2x2", NewSize(Int32(2), Int32(2)), 2},
		{"4x4", NewSize(Int32(4), Int32(4)), 3},
		{"8x8", NewSize(Int32(8), Int32(8)), 4},
		{"16x16", NewSize(Int32(16), Int32(16)), 5},
		{"256x256", NewSize(Int32(256), Int32(256)), 9},
		{"4x2", NewSize(Int32(4), Int32(2)), 3},
		{"8x4", NewSize(Int32(8), Int32(4)), 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.size.MipCount()
			if result != tt.expected {
				t.Errorf("MipCount() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSize_SizeString(t *testing.T) {
	s := NewSize(Float32(10.5), Float32(20.5))
	str := s.String()
	expected := "Size(10.5, 20.5)"
	if str != expected {
		t.Errorf("String() = %v, want %v", str, expected)
	}
}

func BenchmarkSizeOperations(b *testing.B) {
	s1 := NewSize(Float32(10.0), Float32(20.0))
	s2 := NewSize(Float32(5.0), Float32(10.0))

	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s1.Add(s2)
		}
	})

	b.Run("Area", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s1.Area()
		}
	})

	b.Run("IsSquare", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s1.IsSquare()
		}
	})
}
