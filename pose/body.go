package pose

import "github.com/egonelbre/hexapod/g"

type Body struct {
	Size   g.Vec
	Origin g.Vec
	Head   Head
	Leg    Legs
}

func (body *Body) Legs() []*Leg {
	return []*Leg{
		&body.Leg.LF, &body.Leg.RF,
		&body.Leg.LM, &body.Leg.RM,
		&body.Leg.LB, &body.Leg.RB,
	}
}

type Head struct {
	Offset g.Vec
}

type Legs struct {
	LF, RF Leg
	LM, RM Leg
	LB, RB Leg
}

type Leg struct {
	Offset g.Vec // relative to body Origin
	Coxa   Hinge
	Femur  Hinge
	Tibia  Hinge

	// Mainly for debug purposes
	IK LegIK
}

type LegIK struct {
	Origin g.Vec
	Target g.Vec
}

func (leg *Leg) Hinges() []*Hinge {
	return []*Hinge{
		&leg.Coxa,
		&leg.Femur,
		&leg.Tibia,
	}
}

type Axis byte

const (
	X = Axis(iota)
	Y
	Z
)

type Hinge struct {
	// const
	Axis   Axis
	Zero   g.Radians
	Length g.Length
	Range  HingeRange

	// runtime
	Angle g.Radians
}

type HingeRange struct {
	Min, Max g.Radians
}
