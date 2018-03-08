package adeept

import (
	"github.com/egonelbre/hexapod/g"
	"github.com/egonelbre/hexapod/pose"
)

const (
	legRot = g.Tau / 8           //TODO: fix
	legY   = 28*g.MM - 45*g.MM/2 // relative to body center
)

func ZeroPose() *pose.Body {
	return &pose.Body{
		Size:   g.Vec{105 * g.MM, 45 * g.MM, 105 * g.MM},
		Origin: g.Vec{0, 22 * g.MM, 0},
		Head: pose.Head{
			Offset: g.Vec{63 * g.MM, 20 * g.MM, 0},
		},
		Leg: pose.Legs{
			RF: ZeroLeg("RF", g.Vec{63 * g.MM, legY, +57 * g.MM}, +legRot, 1),
			LF: ZeroLeg("LF", g.Vec{63 * g.MM, legY, -57 * g.MM}, -legRot, -1),
			RM: ZeroLeg("RM", g.Vec{0, legY, +77 * g.MM}, +g.Tau/4, 1),
			LM: ZeroLeg("LM", g.Vec{0, legY, -77 * g.MM}, -g.Tau/4, -1),
			RB: ZeroLeg("RB", g.Vec{-63 * g.MM, legY, +57 * g.MM}, +g.Tau/2-legRot, 1),
			LB: ZeroLeg("LB", g.Vec{-63 * g.MM, legY, -57 * g.MM}, -g.Tau/2+legRot, -1),
		},
	}
}

func ZeroLeg(name string, offset g.Vec, zero g.Radians, side g.Radians) pose.Leg {
	target := offset
	target.Y = 0
	return pose.Leg{
		Name:   name,
		Offset: offset,
		Coxa: pose.Hinge{
			Axis:   pose.Y,
			Zero:   zero,
			Length: 12 * g.MM,
			Range:  pose.HingeRange{side * -g.Tau / 4, side * g.Tau / 4},
		},
		Femur: pose.Hinge{
			Axis:   pose.Z,
			Length: 38 * g.MM,
			Range:  pose.HingeRange{g.Tau/4 + g.Tau/8, -g.Tau/4 + g.Tau/8},
		},
		Tibia: pose.Hinge{
			Axis:   pose.Z,
			Length: 50 * g.MM,
			Range:  pose.HingeRange{g.Tau / 4, -g.Tau / 4},
		},
		IK: pose.LegIK{
			Origin: offset,
			Target: target,
		},
	}
}
