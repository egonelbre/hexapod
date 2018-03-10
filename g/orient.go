package g

type Orient struct {
	Yaw   Radians // rotation around Y
	Pitch Radians // rotation around Z
	Roll  Radians // rotation around X
}

func (orient Orient) Mat() Mat {
	return RotateY(orient.Yaw).
		Mul(RotateZ(orient.Pitch)).
		Mul(RotateX(orient.Roll))
}

func (orient Orient) InvMat() Mat {
	return RotateX(-orient.Roll).
		Mul(RotateZ(-orient.Pitch)).
		Mul(RotateY(-orient.Yaw))
}
