package particles

import (
	"math"
	"math/rand"
)

// ScreenShake manages trauma decay and camera offsets.
type ScreenShake struct {
	Trauma float64 // 0.0 to 1.0
	MaxOff float64 // Maximum pixel offset at trauma = 1.0
}

func NewScreenShake() *ScreenShake {
	return &ScreenShake{
		Trauma: 0.0,
		MaxOff: 12.0,
	}
}

func (s *ScreenShake) AddTrauma(amount float64) {
	s.Trauma = math.Min(1.0, s.Trauma+amount)
}

func (s *ScreenShake) Update(dt float64) {
	if s.Trauma > 0 {
		s.Trauma = math.Max(0, s.Trauma-dt*1.6)
	}
}

func (s *ScreenShake) GetOffset() (float64, float64) {
	if s.Trauma <= 0.01 {
		return 0, 0
	}
	shake := s.Trauma * s.Trauma // Non-linear response curve
	offX := (rand.Float64()*2.0 - 1.0) * s.MaxOff * shake
	offY := (rand.Float64()*2.0 - 1.0) * s.MaxOff * shake
	return offX, offY
}
