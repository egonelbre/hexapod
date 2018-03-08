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
		SolveLeg(body, leg, leg.IK.Target)
	}
}

func SolveLeg(body *pose.Body, leg *pose.Leg, worldTarget g.Vec) {
	// TODO: handle body rotation

	// Calculate Coxa.Angle by looking in Body-Space top-down
	coxaOrigin := body.Origin.Add(leg.Offset)
	target := worldTarget.Sub(coxaOrigin)
	coxaAngleRelativeToForward := -g.Atan2(target.X, target.Z)
	leg.Coxa.Angle = g.Radians(coxaAngleRelativeToForward) - leg.Coxa.Zero + g.Tau/4

	if leg.Coxa.Angle > g.Tau/2 {
		leg.Coxa.Angle = leg.Coxa.Angle - g.Tau
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
	target = target.Transform(g.RotateY(-coxaAngleRelativeToForward))

	// move 0,0,0 to Femur root
	target.X = 0
	target.Z -= leg.Coxa.Length

	femurLength2 := leg.Femur.Length * leg.Femur.Length
	tibiaLength2 := leg.Tibia.Length * leg.Tibia.Length
	targetDistance2 := target.Length2()
	targetDistance := targetDistance2.Sqrt()

	femurToTargetAngle := g.Atan2(-target.Y, target.Z)

	if leg.Femur.Length+leg.Tibia.Length < targetDistance {
		// too far away, straighten leg and point at
		leg.Femur.Angle = femurToTargetAngle
		leg.Tibia.Angle = 0

		if pose.VectorPlanted(worldTarget) {
			// if the target should be planted, try to plant foot
			footLength := leg.Femur.Length + leg.Tibia.Length
			if footLength > g.Abs(coxaOrigin.Y) {
				leg.Femur.Angle = g.Asin(coxaOrigin.Y.Float32() / footLength.Float32())
			} else {
				leg.Femur.Angle = g.Tau / 4
			}
		}

		leg.Femur.Clamp()
		return
	}

	femurInternalAngle := g.Acos((femurLength2 + targetDistance2 - tibiaLength2).Float32() / (2 * leg.Femur.Length * targetDistance).Float32())
	leg.Femur.Angle = femurToTargetAngle - femurInternalAngle

	tibiaInternalAngle := g.Acos((femurLength2 + tibiaLength2 - targetDistance2).Float32() / (2 * leg.Femur.Length * leg.Tibia.Length).Float32())
	leg.Tibia.Angle = g.Tau/2 - tibiaInternalAngle

	if leg.Tibia.Clamp() {
		if pose.VectorPlanted(worldTarget) {
			footLength2 := femurLength2 + tibiaLength2 - g.Length((2*leg.Femur.Length*leg.Tibia.Length).Float32()*g.Cos(leg.Tibia.Angle))
			footLength := footLength2.Sqrt()

			footInternalAngle := g.Asin(leg.Tibia.Length.Float32() * g.Sin(leg.Tibia.Angle) / footLength.Float32())
			footAngle := g.Asin(coxaOrigin.Y.Float32() / footLength.Float32())

			leg.Femur.Angle = footAngle - footInternalAngle
		}
	}
	leg.Femur.Clamp()
}
