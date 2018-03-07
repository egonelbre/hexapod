package pose

import "github.com/egonelbre/hexapod/g"

type Body struct {
	Size g.Vec
	Leg  Legs
}

type Legs struct {
	LF, RF Leg
	LM, RM Leg
	LB, RB Leg
}

type Leg struct {
	Origin g.Vec
	Coxa   Hinge
	Femur  Hinge
	Tibia  Hinge
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
	Range  HingeRange

	// runtime
	Angle g.Radians
}

type HingeRange struct {
	Min, Max g.Radians
}
