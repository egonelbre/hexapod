package g

type Orient struct {
	Yaw   Radians // rotation around Y
	Pitch Radians // rotation around Z
	Roll  Radians // rotation around X
}

func (orient Orient) Matrix() Mat {
	sy, cy := Sincos(orient.Yaw)
	sz, cz := Sincos(orient.Pitch)
	sx, cx := Sincos(orient.Roll)

	return Mat{
		cy * cz, (cx * sz) + ((sx * cz) * sy), (sx * sz) - ((cx * cz) * sy), 0,
		-cy * sz, (cx * cz) - ((sx * sz) * sy), (sx * cz) + ((cx * sz) * sy), 0,
		sy, -sx * cy, cx * cy, 0,
		0, 0, 0, 1,
	}
}
