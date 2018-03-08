package g

import "math"

const (
	Tau      Radians = 2 * math.Pi
	RadToDeg         = 360 / Tau
)

type Radians = float32

func Sin(v Radians) Radians { return Radians(math.Sin(float64(v))) }
func Cos(v Radians) Radians { return Radians(math.Cos(float64(v))) }

func Sincos(v Radians) (sin, cos Radians) {
	sn, cs := math.Sincos(float64(v))
	return Radians(sn), Radians(cs)
}

func Atan2(y, x Length) Radians {
	return Radians(math.Atan2(y.Float64(), x.Float64()))
}

func Asin(v Radians) Radians { return Radians(math.Asin(float64(v))) }
func Acos(v Radians) Radians { return Radians(math.Acos(float64(v))) }
