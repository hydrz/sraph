package gio

import (
	"math"
	"time"

	"github.com/opensraph/sraph/geom"
)

// TouchID represents a unique touch point identifier
type TouchID int

// GestureID represents a unique gesture identifier
type GestureID int

// TouchPhase represents the phase of a touch
type TouchPhase int

const (
	TouchPhaseBegan TouchPhase = iota
	TouchPhaseMoved
	TouchPhaseStationary
	TouchPhaseEnded
	TouchPhaseCancelled
)

func (p TouchPhase) String() string {
	switch p {
	case TouchPhaseBegan:
		return "began"
	case TouchPhaseMoved:
		return "moved"
	case TouchPhaseStationary:
		return "stationary"
	case TouchPhaseEnded:
		return "ended"
	case TouchPhaseCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// GestureType represents the type of gesture
type GestureType int

const (
	GestureTypeUnknown GestureType = iota
	GestureTypeTap
	GestureTypeDoubleTap
	GestureTypeLongPress
	GestureTypePan
	GestureTypePinch
	GestureTypeRotation
	GestureTypeSwipe
	GestureTypeFling
)

func (g GestureType) String() string {
	switch g {
	case GestureTypeTap:
		return "tap"
	case GestureTypeDoubleTap:
		return "double_tap"
	case GestureTypeLongPress:
		return "long_press"
	case GestureTypePan:
		return "pan"
	case GestureTypePinch:
		return "pinch"
	case GestureTypeRotation:
		return "rotation"
	case GestureTypeSwipe:
		return "swipe"
	case GestureTypeFling:
		return "fling"
	default:
		return "unknown"
	}
}

// GestureState represents the state of a gesture
type GestureState int

const (
	GestureStatePossible GestureState = iota
	GestureStateBegan
	GestureStateChanged
	GestureStateEnded
	GestureStateCancelled
	GestureStateFailed
)

func (s GestureState) String() string {
	switch s {
	case GestureStatePossible:
		return "possible"
	case GestureStateBegan:
		return "began"
	case GestureStateChanged:
		return "changed"
	case GestureStateEnded:
		return "ended"
	case GestureStateCancelled:
		return "cancelled"
	case GestureStateFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// TouchPoint represents a single touch point
type TouchPoint struct {
	ID        TouchID
	Position  geom.Point[geom.F32]
	Pressure  float32
	Size      geom.Point[geom.F32]
	StartTime time.Time
	LastTime  time.Time
	StartPos  geom.Point[geom.F32]
	LastPos   geom.Point[geom.F32]
	Phase     TouchPhase
}

// Gesture represents a recognized gesture
type Gesture struct {
	ID          GestureID
	Type        GestureType
	State       GestureState
	Position    geom.Point[geom.F32]
	StartPos    geom.Point[geom.F32]
	Translation geom.Point[geom.F32]
	Velocity    geom.Point[geom.F32]
	Scale       float32
	Rotation    float32
	StartTime   time.Time
	Duration    time.Duration
	TouchPoints []TouchID
	WindowID    WindowID
}

// GestureThresholds contains threshold values for gesture recognition
type GestureThresholds struct {
	TapMaxDuration       time.Duration
	TapMaxDistance       float32
	DoubleTapMaxInterval time.Duration
	DoubleTapMaxDistance float32
	LongPressMinDuration time.Duration
	LongPressMaxDistance float32
	PanMinDistance       float32
	PinchMinDistance     float32
	RotationMinAngle     float32
	SwipeMinVelocity     float32
	FlingMinVelocity     float32
}

// DefaultGestureThresholds returns default gesture recognition thresholds
func DefaultGestureThresholds() GestureThresholds {
	return GestureThresholds{
		TapMaxDuration:       300 * time.Millisecond,
		TapMaxDistance:       10.0,
		DoubleTapMaxInterval: 300 * time.Millisecond,
		DoubleTapMaxDistance: 20.0,
		LongPressMinDuration: 500 * time.Millisecond,
		LongPressMaxDistance: 10.0,
		PanMinDistance:       10.0,
		PinchMinDistance:     20.0,
		RotationMinAngle:     math.Pi / 12, // 15 degrees
		SwipeMinVelocity:     100.0,
		FlingMinVelocity:     500.0,
	}
}
