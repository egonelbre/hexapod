package g

import "math"

// Fixed point lengths
const (
	M  = 100000
	MM = 100

	TAU      = 2 * math.Pi
	RadToDeg = 360 / TAU
)

type Radians = float32
type Length int32

func (length Length) Meters() float32      { return length.Float32() / M }
func (length Length) Millimeters() float32 { return length.Float32() / MM }

func (length Length) Float64() float64 { return float64(length) }
func (length Length) Float32() float32 { return float32(length) }

func (length Length) Sqrt() Length { return Length(math.Sqrt(float64(length))) }

type Vec struct{ X, Y, Z Length }

var (
	Forward = Vec{+1, 0, 0}
	Back    = Vec{-1, 0, 0}
	Left    = Vec{0, 0, -1}
	Right   = Vec{0, 0, +1}
	Up      = Vec{0, +1, 0}
	Down    = Vec{0, -1, 0}
)

func (a Vec) Add(b Vec) Vec {
	return Vec{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

func (a Vec) Sub(b Vec) Vec {
	return Vec{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func (a Vec) Mul(b Vec) Vec {
	return Vec{
		X: a.X * b.X,
		Y: a.Y * b.Y,
		Z: a.Z * b.Z,
	}
}

func (a Vec) Dot(b Vec) Length      { return a.X*b.X + a.Y*b.Y + a.Z*b.Z }
func (a Vec) Length() Length        { return a.Dot(a).Sqrt() }
func (a Vec) Distance(b Vec) Length { return a.Sub(b).Length() }

func (a Vec) Meters() struct{ X, Y, Z float32 } {
	return struct{ X, Y, Z float32 }{a.X.Meters(), a.Y.Meters(), a.Z.Meters()}
}
func (a Vec) Millimeters() struct{ X, Y, Z float32 } {
	return struct{ X, Y, Z float32 }{a.X.Millimeters(), a.Y.Millimeters(), a.Z.Millimeters()}
}
