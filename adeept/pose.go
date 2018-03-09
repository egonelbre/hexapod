package adeept

import (
	"github.com/egonelbre/hexapod/g"
	"github.com/egonelbre/hexapod/pose"
)

const (
	legRot = g.Tau / 8           //TODO: fix
	legY   = 28*g.MM - 45*g.MM/2 // relative to body center

	femurAngleOffset = 0 // + g.Tau/8
	tibiaAngleOffset = 0 // + g.Tau/8

	servo_sg92r_sec_per60deg = 0.1
	servo_sg92r_speed        = 60.0 * g.DegToRad / servo_sg92r_sec_per60deg
)

func ZeroPose() *pose.Body {
	return &pose.Body{
		Size:   g.Vec{105 * g.MM, 45 * g.MM, 105 * g.MM},
		Origin: g.Vec{0, 22 * g.MM, 0},
		Head: pose.Head{
			Offset: g.Vec{63 * g.MM, 20 * g.MM, 0},
		},
		Leg: pose.Legs{
			RF: ZeroLeg("RF", g.Vec{63 * g.MM, legY, +57 * g.MM}, +legRot, 1, 3*g.Tau/6),
			LF: ZeroLeg("LF", g.Vec{63 * g.MM, legY, -57 * g.MM}, -legRot, -1, 0*g.Tau/6),
			RM: ZeroLeg("RM", g.Vec{0, legY, +77 * g.MM}, +g.Tau/4, 1, 1*g.Tau/6),
			LM: ZeroLeg("LM", g.Vec{0, legY, -77 * g.MM}, -g.Tau/4, -1, 4*g.Tau/6),
			RB: ZeroLeg("RB", g.Vec{-63 * g.MM, legY, +57 * g.MM}, +g.Tau/2-legRot, 1, 5*g.Tau/6),
			LB: ZeroLeg("LB", g.Vec{-63 * g.MM, legY, -57 * g.MM}, -g.Tau/2+legRot, -1, 2*g.Tau/6),
		},
	}
}

func ZeroLeg(name string, offset g.Vec, zero g.Radians, side g.Radians, phase g.Radians) pose.Leg {
	target := offset
	target.Y = 0
	return pose.Leg{
		Name:   name,
		Phase:  -phase,
		Offset: offset,
		Coxa: pose.Hinge{
			Axis:   pose.Y,
			Zero:   zero,
			Length: 12 * g.MM,
			Speed:  servo_sg92r_speed,
			Range:  pose.HingeRange{side * -g.Tau / 4, side * g.Tau / 4},
		},
		Femur: pose.Hinge{
			Axis:   pose.Z,
			Length: 38 * g.MM,
			Speed:  servo_sg92r_speed,
			Range:  pose.HingeRange{g.Tau/4 + femurAngleOffset, -g.Tau/4 + femurAngleOffset},
		},
		Tibia: pose.Hinge{
			Axis:   pose.Z,
			Length: 50 * g.MM,
			Speed:  servo_sg92r_speed,
			Range:  pose.HingeRange{g.Tau/4 + tibiaAngleOffset, -g.Tau/4 + tibiaAngleOffset},
		},
		IK: pose.LegIK{
			Origin: offset,
			Target: target,
		},
	}
}
