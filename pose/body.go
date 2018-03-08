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
	Name   string
	Offset g.Vec // relative to body Origin
	Coxa   Hinge
	Femur  Hinge
	Tibia  Hinge

	// Mainly for debug purposes
	IK LegIK
}

type LegIK struct {
	// Robot Coordinate Space
	Origin g.Vec
	Target g.Vec
	Solved bool
	Debug  string
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

func (hinge *Hinge) lowhigh() (low, high g.Radians) {
	if hinge.Range.Min < hinge.Range.Max {
		return hinge.Range.Min, hinge.Range.Max
	}
	return hinge.Range.Max, hinge.Range.Min
}

func (hinge *Hinge) InBounds() bool {
	low, high := hinge.lowhigh()
	return low <= hinge.Angle && hinge.Angle <= high
}

func (hinge *Hinge) Clamp() bool {
	low, high := hinge.lowhigh()
	if hinge.Angle < low {
		hinge.Angle = low
		return true
	} else if hinge.Angle > high {
		hinge.Angle = high
		return true
	}
	return false
}

func VectorPlanted(v g.Vec) bool {
	return g.Abs(v.Y) < 1*g.MM
}
