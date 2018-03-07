package robot

import "github.com/egonelbre/hexapod/g"

type Body struct {
	Leg struct {
		LF, RF Leg
		LM, RM Leg
		LB, RB Leg
	}
}

type Leg struct {
	Origin g.Vec
	Coxa   g.Hinge
	Femur  g.Hinge
	Tibia  g.Hinge
}

type RotationAxis byte

const (
	RotationX = RotationAxis(iota)
	RotationY
	RotationZ
)

type HingeX = Hinge
type HingeY = Hinge
type HingeZ = Hinge

type Hinge struct {
	// const
	Axis   RotationAxis
	Zero   g.Radians
	Length g.Length

	Range struct {
		Min g.Radians
		Max g.Radians
	}

	// runtime
	Angle g.Radians
}
