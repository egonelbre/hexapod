package g

import "math"

// Fixed point lengths
const (
	M  Length = MM * 1000
	MM Length = 100
)

type Length float32

func (length Length) Meters() float32      { return length.Float32() / M.Float32() }
func (length Length) Millimeters() float32 { return length.Float32() / MM.Float32() }

func (length Length) Float64() float64 { return float64(length) }
func (length Length) Float32() float32 { return float32(length) }

func (length Length) Scale(v float32) Length { return Length(float32(length) * v) }
func (length Length) Sqrt() Length           { return Length(math.Sqrt(float64(length))) }

func Abs(v Length) Length {
	if v < 0 {
		return -v
	}
	return v
}
