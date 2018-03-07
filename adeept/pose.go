package adeept

import (
	"github.com/egonelbre/hexapod/g"
	"github.com/egonelbre/hexapod/pose"
)

const (
	legRot = g.Tau / 8 //TODO: fix
	legY   = 28 * g.MM
)

func ZeroPose() *pose.Body {
	return &pose.Body{
		Size: g.Vec{126 * g.MM, 44 * g.MM, 114 * g.MM},
		Leg: pose.Legs{
			RF: ZeroLeg(g.Vec{63 * g.MM, legY, +57 * g.MM}, legRot),
			LF: ZeroLeg(g.Vec{63 * g.MM, legY, -57 * g.MM}, -legRot),
			RM: ZeroLeg(g.Vec{0, legY, +77 * g.MM}, +g.Tau/4),
			LM: ZeroLeg(g.Vec{0, legY, -77 * g.MM}, -g.Tau/4),
			RB: ZeroLeg(g.Vec{-63 * g.MM, legY, +57 * g.MM}, g.Tau/2-legRot),
			LB: ZeroLeg(g.Vec{-63 * g.MM, legY, -57 * g.MM}, -g.Tau/2-legRot),
		},
	}
}

func ZeroLeg(origin g.Vec, zero g.Radians) pose.Leg {
	return pose.Leg{
		Origin: origin,
		Coxa: pose.Hinge{
			Axis:   pose.RotationY,
			Zero:   zero,
			Length: 12 * g.MM,
			Range:  pose.HingeRange{-g.Tau / 4, g.Tau / 4},
		},
		Femur: pose.Hinge{
			Axis:   pose.RotationZ,
			Length: 38 * g.MM,
			Range:  pose.HingeRange{g.Tau / 4, -g.Tau / 4},
		},
		Tibia: pose.Hinge{
			Axis:   pose.RotationZ,
			Length: 50 * g.MM,
			Range:  pose.HingeRange{g.Tau / 4, -g.Tau / 4},
		},
	}
}
