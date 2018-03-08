package g

type Mat struct {
	M00, M01, M02, M03 float32
	M10, M11, M12, M13 float32
	M20, M21, M22, M23 float32
	M30, M31, M32, M33 float32
}

func Translate(x, y, z Length) Mat {
	return Mat{
		1, 0, 0, x.Float32(),
		0, 1, 0, y.Float32(),
		0, 0, 1, z.Float32(),
		0, 0, 0, 1,
	}
}

func Scale(x, y, z float32) Mat {
	return Mat{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	}
}

func RotateX(angle Radians) Mat {
	sn, cs := Sincos(angle)
	return Mat{
		1, 0, 0, 0,
		0, cs, -sn, 0,
		0, sn, cs, 0,
		0, 0, 0, 1,
	}
}

func RotateY(angle Radians) Mat {
	sn, cs := Sincos(angle)
	return Mat{
		cs, 0, sn, 0,
		0, 1, 0, 0,
		-sn, 0, cs, 0,
		0, 0, 0, 1,
	}
}

func RotateZ(angle Radians) Mat {
	sn, cs := Sincos(angle)
	return Mat{
		cs, -sn, 0, 0,
		sn, cs, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func (a Mat) Mul(b Mat) Mat {
	var r Mat

	r.M00 = a.M00*b.M00 + a.M01*b.M10 + a.M02*b.M20 + a.M03*b.M30
	r.M01 = a.M00*b.M01 + a.M01*b.M11 + a.M02*b.M21 + a.M03*b.M31
	r.M02 = a.M00*b.M02 + a.M01*b.M12 + a.M02*b.M22 + a.M03*b.M32
	r.M03 = a.M00*b.M03 + a.M01*b.M13 + a.M02*b.M23 + a.M03*b.M33
	r.M10 = a.M10*b.M00 + a.M11*b.M10 + a.M12*b.M20 + a.M13*b.M30
	r.M11 = a.M10*b.M01 + a.M11*b.M11 + a.M12*b.M21 + a.M13*b.M31
	r.M12 = a.M10*b.M02 + a.M11*b.M12 + a.M12*b.M22 + a.M13*b.M32
	r.M13 = a.M10*b.M03 + a.M11*b.M13 + a.M12*b.M23 + a.M13*b.M33
	r.M20 = a.M20*b.M00 + a.M21*b.M10 + a.M22*b.M20 + a.M23*b.M30
	r.M21 = a.M20*b.M01 + a.M21*b.M11 + a.M22*b.M21 + a.M23*b.M31
	r.M22 = a.M20*b.M02 + a.M21*b.M12 + a.M22*b.M22 + a.M23*b.M32
	r.M23 = a.M20*b.M03 + a.M21*b.M13 + a.M22*b.M23 + a.M23*b.M33
	r.M30 = a.M30*b.M00 + a.M31*b.M10 + a.M32*b.M20 + a.M33*b.M30
	r.M31 = a.M30*b.M01 + a.M31*b.M11 + a.M32*b.M21 + a.M33*b.M31
	r.M32 = a.M30*b.M02 + a.M31*b.M12 + a.M32*b.M22 + a.M33*b.M32
	r.M33 = a.M30*b.M03 + a.M31*b.M13 + a.M32*b.M23 + a.M33*b.M33

	return r
}

func (v Vec) Transform(t Mat) Vec {
	return Vec{
		X: Length(t.M00*v.X.Float32() + t.M10*v.Y.Float32() + t.M20*v.Z.Float32() + t.M30),
		Y: Length(t.M01*v.X.Float32() + t.M11*v.Y.Float32() + t.M21*v.Z.Float32() + t.M31),
		Z: Length(t.M02*v.X.Float32() + t.M12*v.Y.Float32() + t.M22*v.Z.Float32() + t.M32),
	}
}

// TODO: implement inlined Rotate(X|Y|Z)
