package main

import (
	"fmt"

	"github.com/egonelbre/hexapod/g"
	"github.com/egonelbre/hexapod/ik/legik"
	"github.com/egonelbre/hexapod/pose"
	"github.com/gen2brain/raylib-go/raylib"
)

type Robot struct {
	Body     *pose.Body
	Time     float32
	Active   int
	Controls []Control
}

func NewRobot(body *pose.Body) *Robot {
	robot := &Robot{}
	robot.Body = body
	robot.Controls = []Control{
		&Controller{},
		&Stand{},
		&Stand{},
		&TippyTaps1{},
		&TippyTaps2{},
		&TippyTaps3{},
		&Tapping{},
		&Impatient{},
		&Yay{},
	}
	return robot
}

func (robot *Robot) ModeName() string {
	return ControlName(robot.Controls[robot.Active])
}
func (robot *Robot) Toggle(inc int) {
	robot.Active += inc
	for robot.Active < 0 {
		robot.Active += len(robot.Controls)
	}
	for robot.Active >= len(robot.Controls) {
		robot.Active -= len(robot.Controls)
	}
}

func (robot *Robot) Update(dt float32) {
	robot.Time += dt

	control := robot.Controls[robot.Active]
	control.Update(robot.Body, robot.Time, dt)

	legik.Solve(robot.Body)
}

type Control interface {
	Update(body *pose.Body, time, dt float32)
}

func ControlName(control Control) string { return fmt.Sprintf("%T", control) }

type Controller struct {
	Move g.Vec

	Yaw   float32
	Pitch float32
	Roll  float32

	Elevation g.Length
}

func (ctrl *Controller) readInput() {
	if raylib.IsGamepadAvailable(0) {
		leftX := raylib.GetGamepadAxisMovement(0, raylib.GamepadXboxAxisLeftX)
		leftY := raylib.GetGamepadAxisMovement(0, raylib.GamepadXboxAxisLeftY)

		rightX := raylib.GetGamepadAxisMovement(0, raylib.GamepadXboxAxisRightX)
		rightY := raylib.GetGamepadAxisMovement(0, raylib.GamepadXboxAxisRightY)

		rightTrigger := raylib.GetGamepadAxisMovement(0, raylib.GamepadXboxAxisRt)

		ctrl.Move.X = g.Length(leftX * float32(g.M))
		ctrl.Move.Y = g.Length(leftY * float32(g.M))
		ctrl.Move = ctrl.Move.NormalizedTo(20 * g.MM)

		ctrl.Roll = rightX * g.Tau / 16
		ctrl.Pitch = rightY * g.Tau / 16

		ctrl.Elevation = (20 * g.MM).Scale(rightTrigger*0.5 + 0.5)
	}
}

func (ctrl *Controller) Update(body *pose.Body, time, dt float32) {
	ctrl.readInput()

	bodyOscX := g.Sin(time * 0.5)
	bodyOscZ := g.Cos(time + g.Tau/16)

	body.Origin.X = g.Length(3*bodyOscX)*g.MM - 5*g.MM
	body.Origin.Z = g.Length(7*bodyOscZ) * g.MM
	body.Origin.Y = ctrl.Elevation + body.Size.Y*0.5 + 30*g.MM

	body.Orient.Yaw = ctrl.Yaw
	body.Orient.Pitch = ctrl.Pitch
	body.Orient.Roll = ctrl.Roll

	for _, leg := range body.Legs() {
		sn, cs := g.Sincos(time + leg.Phase)

		leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(70 * g.MM))
		leg.IK.Target.X = leg.Offset.X

		side := g.Length(1)
		if leg.IK.Target.Z < 0 {
			side = -1
		}

		if leg.Name[1] != 'M' {
			leg.IK.Target.Z += side * 10 * g.MM
		} else {
			leg.IK.Target.Z -= side * 10 * g.MM
		}

		leg.IK.Target.Y = g.Length(20 * cs * g.MM.Float32())
		leg.IK.Target.X += g.Length(20 * sn * g.MM.Float32())

		leg.IK.Planted = false
		if leg.IK.Target.Y < 0 {
			leg.IK.Target.Y = 0
			leg.IK.Planted = true
		}
	}
}

type Stand struct{}

func (*Stand) Update(body *pose.Body, time, dt float32) {
	bodyOscY := g.Sin(time)
	if bodyOscY < 0 {
		bodyOscY = -bodyOscY
	}
	body.Origin.Y = g.Length(bodyOscY*float32(5*g.MM)) + body.Size.Y

	body.Orient.Pitch = g.Sin(time*0.1) * g.Tau / 32
	body.Orient.Yaw = g.Sin(time*0.14) * g.Tau / 32
	body.Orient.Roll = g.Sin(time*0.20) * g.Tau / 32

	bodyOscX := g.Sin(time * 0.5)
	bodyOscZ := g.Cos(time + g.Tau/16)
	body.Origin.X = g.Length(3*bodyOscX)*g.MM - 5*g.MM
	body.Origin.Z = g.Length(7*bodyOscZ) * g.MM

	for _, leg := range body.Legs() {
		leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(70 * g.MM))
		leg.IK.Target.Y = 0
		leg.IK.Planted = true
	}
}

type TippyTaps1 struct{}

func (*TippyTaps1) Update(body *pose.Body, time, dt float32) {
	body.Origin.Y = body.Size.Y/2 + body.Size.Y/4

	bodyOscY := g.Sin(time * 2)
	body.Origin.Y = g.Length(bodyOscY*float32(5*g.MM)) + body.Size.Y

	bodyOscX := g.Sin(time * 0.5)
	bodyOscZ := g.Cos(time + g.Tau/16)
	body.Origin.X = g.Length(3*bodyOscX)*g.MM - 5*g.MM
	body.Origin.Z = g.Length(5*bodyOscZ) * g.MM

	for _, leg := range body.Legs() {
		_, cs := g.Sincos(time + leg.Phase)

		leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(70 * g.MM))
		leg.IK.Target.Y = g.Length(20 * cs * g.MM.Float32())

		leg.IK.Planted = false
		if leg.IK.Target.Y < 0 {
			leg.IK.Target.Y = 0
			leg.IK.Planted = true
		}
	}
}

type TippyTaps2 struct{}

func (*TippyTaps2) Update(body *pose.Body, time, dt float32) {
	body.Origin.Y = body.Size.Y/2 + body.Size.Y/4

	bodyOscY := g.Sin(time * 2)
	body.Origin.Y = g.Length(bodyOscY*float32(5*g.MM)) + body.Size.Y

	bodyOscX := g.Sin(time * 0.5)
	bodyOscZ := g.Cos(time + g.Tau/16)
	body.Origin.X = g.Length(3*bodyOscX)*g.MM - 5*g.MM
	body.Origin.Z = g.Length(10*bodyOscZ) * g.MM

	for _, leg := range body.Legs() {
		if leg.Name == "RB" || leg.Name == "LB" {
			leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(70 * g.MM))
			leg.IK.Target.Y = 0
			leg.IK.Planted = true
			continue
		}

		_, cs := g.Sincos(time + leg.Phase)

		leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(70 * g.MM))
		leg.IK.Target.Y = g.Length(20 * cs * g.MM.Float32())

		leg.IK.Planted = false
		if leg.IK.Target.Y < 0 {
			leg.IK.Target.Y = 0
			leg.IK.Planted = true
		}
	}
}

type TippyTaps3 struct{}

func (*TippyTaps3) Update(body *pose.Body, time, dt float32) {
	body.Origin.Y = body.Size.Y/2 + body.Size.Y/4

	bodyOscY := g.Sin(time)
	if bodyOscY < 0 {
		bodyOscY = -bodyOscY
	}
	body.Origin.Y = g.Length(bodyOscY*float32(10*g.MM)) + body.Size.Y

	bodyOscX := g.Sin(time * 0.5)
	bodyOscZ := g.Cos(time + g.Tau/16)
	body.Origin.X = g.Length(3*bodyOscX)*g.MM - 5*g.MM
	body.Origin.Z = g.Length(10*bodyOscZ) * g.MM

	for _, leg := range body.Legs() {
		if leg.Name == "RB" || leg.Name == "LB" {
			leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(70 * g.MM))
			leg.IK.Target.Y = 0
			leg.IK.Planted = true
			continue
		}

		phase := leg.Phase
		if leg.Name[0] == 'R' {
			phase = g.Tau / 2
		} else {
			phase = 0
		}

		sn, cs := g.Sincos(time + phase)

		leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(70 * g.MM))
		leg.IK.Target.Y = g.Length(20 * cs * g.MM.Float32())
		_ = sn
		//leg.IK.Target.X += g.Length(20 * sn * g.MM.Float32())
		leg.IK.Planted = false
		if leg.IK.Target.Y < 0 {
			leg.IK.Target.Y = 0
			leg.IK.Planted = true
		}
	}
}

type Tapping struct{}

func (*Tapping) Update(body *pose.Body, time, dt float32) {
	time *= 0.2
	body.Origin.Y = body.Size.Y/2 + body.Size.Y/4

	bodyOscZ := g.Cos(time)
	body.Origin.Z = g.Length(10*bodyOscZ) * g.MM

	for _, leg := range body.Legs() {
		if leg.Name != "RF" {
			leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(70 * g.MM))
			leg.IK.Target.Y = 0
			leg.IK.Planted = true
			continue
		}

		_, cs := g.Sincos(time*2 + g.Tau/2)

		leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(70 * g.MM))
		leg.IK.Target.Y = g.Length(20 * cs * g.MM.Float32())

		//leg.IK.Target.X += g.Length(20 * sn * g.MM.Float32())
		leg.IK.Planted = false
		if leg.IK.Target.Y < 0 {
			leg.IK.Target.Y = 0
			leg.IK.Planted = true
		}
	}
}

type Impatient struct{}

func (*Impatient) Update(body *pose.Body, time, dt float32) {
	time *= 0.2
	body.Origin.Y = body.Size.Y/2 + body.Size.Y/4

	bodyOscY := g.Sin(time*0.5 + g.Tau/8)
	if bodyOscY < 0 {
		bodyOscY = -bodyOscY
	}
	body.Origin.Y = g.Length(bodyOscY*float32(5*g.MM)) + body.Size.Y

	for _, leg := range body.Legs() {
		if leg.Name != "RF" {
			leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(70 * g.MM))
			leg.IK.Target.Y = 0
			leg.IK.Planted = true
			continue
		}

		_, cs := g.Sincos(time*6 + g.Tau/2)

		leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(70 * g.MM))
		leg.IK.Target.Y = g.Length(20 * cs * g.MM.Float32())

		//leg.IK.Target.X += g.Length(20 * sn * g.MM.Float32())
		leg.IK.Planted = false
		if leg.IK.Target.Y < 0 {
			leg.IK.Target.Y = 0
			leg.IK.Planted = true
		}
	}
}

type Yay struct{}

func (*Yay) Update(body *pose.Body, time, dt float32) {
	time *= 0.2
	body.Origin.Y = body.Size.Y/2 + body.Size.Y/4

	bodyOscY := g.Sin(time + g.Tau/8)
	if bodyOscY < 0 {
		bodyOscY = -bodyOscY
	}
	body.Origin.Y = g.Length(bodyOscY*float32(5*g.MM)) + body.Size.Y

	for _, leg := range body.Legs() {
		if leg.Name != "RF" && leg.Name != "LF" {
			leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(70 * g.MM))
			leg.IK.Target.Y = 0
			leg.IK.Planted = true
			continue
		}

		legOscY := g.Sin(time - g.Tau/4)
		leg.IK.Target = leg.Offset.Add(leg.Offset.NormalizedTo(60 * g.MM))
		//leg.IK.Target.Y = 115*g.MM + g.Length(legOscY*float32(5*g.MM))
		leg.IK.Target.Y = 60*g.MM + g.Length(legOscY*float32(60*g.MM))

		//leg.IK.Target.X += g.Length(20 * sn * g.MM.Float32())
		leg.IK.Planted = false
		if leg.IK.Target.Y < 0 {
			leg.IK.Target.Y = 0
			leg.IK.Planted = true
		}
	}
}
