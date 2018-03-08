package g

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

func (a Vec) Scale(s float32) Vec {
	return Vec{
		X: Length(float32(a.X) * s),
		Y: Length(float32(a.Y) * s),
		Z: Length(float32(a.Z) * s),
	}
}

func (a Vec) Dot(b Vec) Length      { return a.X*b.X + a.Y*b.Y + a.Z*b.Z }
func (a Vec) Length() Length        { return a.Dot(a).Sqrt() }
func (a Vec) Length2() Length       { return a.Dot(a) }
func (a Vec) Distance(b Vec) Length { return a.Sub(b).Length() }

func (a Vec) Meters() struct{ X, Y, Z float32 } {
	return struct{ X, Y, Z float32 }{a.X.Meters(), a.Y.Meters(), a.Z.Meters()}
}
func (a Vec) Millimeters() struct{ X, Y, Z float32 } {
	return struct{ X, Y, Z float32 }{a.X.Millimeters(), a.Y.Millimeters(), a.Z.Millimeters()}
}
func (a Vec) XYZ() (x, y, z float32) {
	return a.X.Meters(), a.Y.Meters(), a.Z.Meters()
}

func (a Vec) NormalizedTo(targetLength Length) Vec {
	length := a.Length()
	return Vec{
		X: a.X * targetLength / length,
		Y: a.Y * targetLength / length,
		Z: a.Z * targetLength / length,
	}
}
