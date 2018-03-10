// package legik solves IK for a 3DOF leg
package legik

import (
	"github.com/egonelbre/hexapod/g"
	"github.com/egonelbre/hexapod/pose"
)

// Assumes:
//   1. Coxa is parallel to the body
//   2. There is no offset between Coxa direction and End Effector

func Solve(body *pose.Body) {
	for _, leg := range body.Legs() {
		leg.IK.Solved = SolveLeg(body, leg, leg.IK.Target)
	}
}

func SolveLeg(body *pose.Body, leg *pose.Leg, worldTarget g.Vec) bool {
	// TODO: handle body rotation

	// Calculate Coxa.Angle by looking in Body-Space top-down

	legTarget := g.Identity().
		Mul(g.Translate(-leg.Offset.X, -leg.Offset.Y, -leg.Offset.Z)).
		Mul(body.Orient.InvMat()).
		Mul(g.Translate(-body.Origin.X, -body.Origin.Y, -body.Origin.Z)).
		Transform(worldTarget)

	// bodyTarget := body.Orient.InvMat().Transform(worldTarget.Sub(body.Origin))
	// legTarget := bodyTarget.Sub(leg.Offset)

	coxaAngleRelativeToForward := -g.Atan2(legTarget.X, legTarget.Z)
	leg.Coxa.Angle = g.Radians(coxaAngleRelativeToForward) - leg.Coxa.Zero + g.Tau/4

	if leg.Coxa.Angle > g.Tau/2 {
		leg.Coxa.Angle = leg.Coxa.Angle - g.Tau
	}
	if leg.Coxa.Clamp() {
		coxaAngleRelativeToForward = coxaAngleRelativeToForward + leg.Coxa.Zero - g.Tau/4
	}

	//  +---------------------------------------------------------+
	//  | in Coxa coordinate space (leg side-view)                |
	//  |                                  .                      |
	//  |                Tibia.Angle ---> / . Tibia.Length        |
	//  |                                O---------------# Goal   |
	//  |                               /                         |
	//  |                              /                          |
	//  |                             /                           |
	//  | ^ y                        / Femur.length               |
	//  | |                         /.                            |
	//  | |                        /  .                           |
	//  | |  #--------------------O-------                        |
	//  | |       Coxa.Length        ^--- Femur.Angle             |
	//  | +-----> x                                               |
	//  +---------------------------------------------------------+

	// move into Coxa space
	target := g.RotateY(-coxaAngleRelativeToForward).Transform(legTarget)

	// move 0,0,0 to Femur root
	target.X = 0
	target.Z -= leg.Coxa.Length

	femurLength2 := leg.Femur.Length * leg.Femur.Length
	tibiaLength2 := leg.Tibia.Length * leg.Tibia.Length
	targetDistance2 := target.Length2()
	targetDistance := targetDistance2.Sqrt()

	footToTargetAngle := g.Atan2(-target.Y, target.Z)

	if leg.Femur.Length+leg.Tibia.Length < targetDistance {
		// too far away, straighten leg and point at
		leg.Femur.Angle = footToTargetAngle
		leg.Tibia.Angle = 0
		if pose.VectorPlanted(worldTarget) {
			/*
				// if the target should be planted, try to plant foot
				footLength := leg.Femur.Length + leg.Tibia.Length
				if footLength > g.Abs(coxaOrigin.Y) {
					leg.Femur.Angle = g.Asin(coxaOrigin.Y.Float32() / footLength.Float32())
				} else {
					leg.Femur.Angle = g.Tau / 4
				}
			*/
		}
		leg.Femur.Clamp()
		return false
	}

	femurInternalAngle := g.Acos((femurLength2 + targetDistance2 - tibiaLength2).Float32() / (2 * leg.Femur.Length * targetDistance).Float32())
	tibiaInternalAngle := g.Acos((femurLength2 + tibiaLength2 - targetDistance2).Float32() / (2 * leg.Femur.Length * leg.Tibia.Length).Float32())

	leg.Femur.Angle = footToTargetAngle - femurInternalAngle
	leg.Tibia.Angle = g.Tau/2 - tibiaInternalAngle

	if false /* leg flipped */ {
		leg.Femur.Angle = footToTargetAngle + femurInternalAngle
		leg.Tibia.Angle = -(g.Tau/2 - tibiaInternalAngle)
	}

	if leg.Tibia.Clamp() {
		if pose.VectorPlanted(worldTarget) {
			/*
				footLength2 := femurLength2 + tibiaLength2 - g.Length((2*leg.Femur.Length*leg.Tibia.Length).Float32()*g.Cos(leg.Tibia.Angle))
				footLength := footLength2.Sqrt()

				footInternalAngle := g.Asin(leg.Tibia.Length.Float32() * g.Sin(leg.Tibia.Angle) / footLength.Float32())
				footAngle := g.Asin(coxaOrigin.Y.Float32() / footLength.Float32())

				leg.Femur.Angle = footAngle - footInternalAngle
			*/
		}
		leg.Femur.Clamp()
		return false
	}
	if leg.Femur.Clamp() {
		return false
	}
	return true
}
