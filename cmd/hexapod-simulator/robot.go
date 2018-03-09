package main

import (
	"fmt"

	"github.com/egonelbre/hexapod/g"
	"github.com/egonelbre/hexapod/ik/legik"
	"github.com/egonelbre/hexapod/pose"

	"github.com/gen2brain/raylib-go/raylib"
	"github.com/gen2brain/raylib-go/raymath"
)

type Model struct {
	Pose *pose.Body

	Shader raylib.Shader
	Labels []Label

	Head raylib.Model
	Body raylib.Model

	LegPlateSize g.Vec
	LegPlate     raylib.Model

	Bone     raylib.Model
	Hinge    raylib.Model
	Effector raylib.Model
}

type Label struct {
	Position raylib.Vector3
	Text     string
	Color    raylib.Color
}

func NewModel(pose *pose.Body) *Model {
	model := &Model{}
	model.Pose = pose

	model.Shader = raylib.LoadShader("shader/base.vs", "shader/lighting.fs")
	model.Shader.Locs[raylib.LocMatrixModel] = raylib.GetShaderLocation(model.Shader, "mMatrix")
	model.Shader.Locs[raylib.LocMatrixView] = raylib.GetShaderLocation(model.Shader, "view")
	model.Shader.Locs[raylib.LocVectorView] = raylib.GetShaderLocation(model.Shader, "viewPos")

	size := pose.Size.Meters()
	mm := g.MM.Meters()

	model.Head = raylib.LoadModelFromMesh(raylib.GenMeshCube(30*mm, 30*mm, 30*mm))
	model.Body = raylib.LoadModelFromMesh(raylib.GenMeshCube(size.X, size.Y, size.Z))

	model.LegPlateSize = g.Vec{
		X: pose.Leg.RF.Offset.X - pose.Leg.RB.Offset.X - 5*g.MM,
		Y: 3 * g.MM,
		Z: pose.Leg.RM.Offset.Z - pose.Size.Z/2 - 2*g.MM,
	}
	model.LegPlate = raylib.LoadModelFromMesh(raylib.GenMeshCube(model.LegPlateSize.XYZ()))

	// size is assigned dynamically with scaling
	model.Bone = raylib.LoadModelFromMesh(raylib.GenMeshCube(1, 1, 1))
	model.Hinge = raylib.LoadModelFromMesh(raylib.GenMeshCube(1, 1, 1))
	model.Effector = raylib.LoadModelFromMesh(raylib.GenMeshCube(1, 1, 1))

	model.Head.Material.Shader = model.Shader
	model.Body.Material.Shader = model.Shader

	model.LegPlate.Material.Shader = model.Shader

	model.Bone.Material.Shader = model.Shader
	model.Hinge.Material.Shader = model.Shader
	model.Effector.Material.Shader = model.Shader

	return model
}

func (model *Model) Update() {
	// time := 2 * float32(raylib.GetMouseX()) / float32(raylib.GetScreenWidth())
	time := raylib.GetTime() * 8

	model.walk(time)
	//model.tippyTaps1(time)
	//model.tapping(time)

	legik.Solve(model.Pose)
}

func (model *Model) walk(time float32) {
	model.Pose.Origin.Y = model.Pose.Size.Y/2 + model.Pose.Size.Y/4

	bodyOscY := g.Sin(time)
	if bodyOscY < 0 {
		bodyOscY = -bodyOscY
	}
	model.Pose.Origin.Y = g.Length(bodyOscY*float32(5*g.MM)) + model.Pose.Size.Y

	bodyOscX := g.Sin(time * 0.5)
	bodyOscZ := g.Cos(time + g.Tau/16)
	model.Pose.Origin.X = g.Length(3*bodyOscX)*g.MM - 5*g.MM
	model.Pose.Origin.Z = g.Length(7*bodyOscZ) * g.MM

	for _, leg := range model.Pose.Legs() {
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

func (model *Model) tippyTaps1(time float32) {
	model.Pose.Origin.Y = model.Pose.Size.Y/2 + model.Pose.Size.Y/4

	bodyOscY := g.Sin(time * 2)
	model.Pose.Origin.Y = g.Length(bodyOscY*float32(5*g.MM)) + model.Pose.Size.Y

	bodyOscX := g.Sin(time * 0.5)
	bodyOscZ := g.Cos(time + g.Tau/16)
	model.Pose.Origin.X = g.Length(3*bodyOscX)*g.MM - 5*g.MM
	model.Pose.Origin.Z = g.Length(5*bodyOscZ) * g.MM

	for _, leg := range model.Pose.Legs() {
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

func (model *Model) tippyTaps2(time float32) {
	model.Pose.Origin.Y = model.Pose.Size.Y/2 + model.Pose.Size.Y/4

	bodyOscY := g.Sin(time * 2)
	model.Pose.Origin.Y = g.Length(bodyOscY*float32(5*g.MM)) + model.Pose.Size.Y

	bodyOscX := g.Sin(time * 0.5)
	bodyOscZ := g.Cos(time + g.Tau/16)
	model.Pose.Origin.X = g.Length(3*bodyOscX)*g.MM - 5*g.MM
	model.Pose.Origin.Z = g.Length(10*bodyOscZ) * g.MM

	for _, leg := range model.Pose.Legs() {
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

func (model *Model) tippyTaps3(time float32) {
	model.Pose.Origin.Y = model.Pose.Size.Y/2 + model.Pose.Size.Y/4

	bodyOscY := g.Sin(time)
	if bodyOscY < 0 {
		bodyOscY = -bodyOscY
	}
	model.Pose.Origin.Y = g.Length(bodyOscY*float32(10*g.MM)) + model.Pose.Size.Y

	bodyOscX := g.Sin(time * 0.5)
	bodyOscZ := g.Cos(time + g.Tau/16)
	model.Pose.Origin.X = g.Length(3*bodyOscX)*g.MM - 5*g.MM
	model.Pose.Origin.Z = g.Length(10*bodyOscZ) * g.MM

	for _, leg := range model.Pose.Legs() {
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

func (model *Model) tapping(time float32) {
	time *= 0.2
	model.Pose.Origin.Y = model.Pose.Size.Y/2 + model.Pose.Size.Y/4

	//bodyOscY := g.Sin(time)
	//if bodyOscY < 0 {
	//	bodyOscY = -bodyOscY
	//}
	//model.Pose.Origin.Y = g.Length(bodyOscY*float32(10*g.MM)) + model.Pose.Size.Y

	//bodyOscX := g.Sin(time)
	//model.Pose.Origin.X = g.Length(3*bodyOscX)*g.MM - 5*g.MM
	bodyOscZ := g.Cos(time)
	model.Pose.Origin.Z = g.Length(10*bodyOscZ) * g.MM

	for _, leg := range model.Pose.Legs() {
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
func (model *Model) AddLabel(text string, pos raylib.Vector3, col raylib.Color) {
	model.Labels = append(model.Labels, Label{pos, text, col})
}

func (model *Model) Draw() {
	model.Labels = model.Labels[:0]
	model.AddLabel(".", model.Pose.Origin.Meters(), raylib.Black)

	zero := raylib.Vector3{}
	center := raylib.Vector3{}

	bodyTransform := raymath.MatrixTranslate(model.Pose.Origin.XYZ())

	model.Head.Transform = bodyTransform
	model.Body.Transform = bodyTransform
	raylib.DrawModel(model.Head, model.Pose.Head.Offset.Meters(), 1, raylib.DarkGray)
	raylib.DrawModel(model.Body, center, 1, raylib.Gray)

	zmm := (model.Pose.Size.Z/2 + model.LegPlateSize.Z/2).Meters()
	ymm := model.Pose.Leg.RF.Offset.Y.Meters()
	model.LegPlate.Transform = bodyTransform
	raylib.DrawModel(model.LegPlate, raylib.Vector3{0, ymm, +zmm}, 1, raylib.LightGray)
	raylib.DrawModel(model.LegPlate, raylib.Vector3{0, ymm, -zmm}, 1, raylib.LightGray)

	mm := g.MM.Meters()

	for _, leg := range model.Pose.Legs() {
		transform := raymath.MatrixMultiply(bodyTransform, raymath.MatrixTranslate(leg.Offset.XYZ()))

		var labelPosition raylib.Vector3
		raymath.Vector3Transform(&labelPosition, transform)
		labelPosition.Y += 30 * mm
		model.AddLabel(leg.Name+" "+leg.IK.Debug, labelPosition, raylib.Black)

		for _, hinge := range leg.Hinges() {
			var newRotation func(v float32) raylib.Matrix
			var hingeScale raylib.Matrix
			switch hinge.Axis {
			case pose.X:
				newRotation = raymath.MatrixRotateX
				hingeScale = raymath.MatrixScale(20*mm, 5*mm, 5*mm)
			case pose.Y:
				newRotation = raymath.MatrixRotateY
				hingeScale = raymath.MatrixScale(5*mm, 20*mm, 5*mm)
			case pose.Z:
				newRotation = raymath.MatrixRotateZ
				hingeScale = raymath.MatrixScale(5*mm, 5*mm, 20*mm)
			}

			var hingeCenter raylib.Vector3
			raymath.Vector3Transform(&hingeCenter, transform)
			model.AddLabel(fmt.Sprintf("%.0f", hinge.Angle*g.RadToDeg), hingeCenter, raylib.Black)

			hingePointMin := raylib.Vector3{20 * mm, 0, 0}
			hingePointZero := raylib.Vector3{20 * mm, 0, 0}
			hingePointMax := raylib.Vector3{20 * mm, 0, 0}

			raymath.Vector3Transform(&hingePointMin, raymath.MatrixMultiply(transform, newRotation(hinge.Zero+hinge.Range.Min)))
			raymath.Vector3Transform(&hingePointZero, raymath.MatrixMultiply(transform, newRotation(hinge.Zero)))
			raymath.Vector3Transform(&hingePointMax, raymath.MatrixMultiply(transform, newRotation(hinge.Zero+hinge.Range.Max)))

			raylib.DrawLine3D(hingeCenter, hingePointMin, raylib.Red)
			raylib.DrawLine3D(hingeCenter, hingePointZero, raylib.Green)
			raylib.DrawLine3D(hingeCenter, hingePointMax, raylib.Blue)

			transform = raymath.MatrixMultiply(transform, newRotation(hinge.Zero+hinge.Angle))
			model.Hinge.Transform = raymath.MatrixMultiply(transform, hingeScale)

			hingeLength := hinge.Length.Meters()

			boneTransform := raymath.MatrixMultiply(transform,
				raymath.MatrixTranslate(hingeLength/2, 0, 0))

			model.Bone.Transform = raymath.MatrixMultiply(boneTransform,
				raymath.MatrixScale(hingeLength, 4*mm, 4*mm))

			raylib.DrawModel(model.Hinge, zero, 1, raylib.Blue)

			if hinge.InBounds() {
				raylib.DrawModel(model.Bone, zero, 1, raylib.SkyBlue)
			} else {
				raylib.DrawModel(model.Bone, zero, 1, raylib.NewColor(255, 102, 191, 255))
			}

			transform = raymath.MatrixMultiply(transform, raymath.MatrixTranslate(hingeLength, 0, 0))
		}

		effectorColor := raylib.Blue
		if !leg.IK.Solved {
			effectorColor = raylib.Red
		}
		model.Effector.Transform = raymath.MatrixMultiply(transform,
			raymath.MatrixScale(5*mm, 5*mm, 5*mm))
		raylib.DrawModel(model.Effector, zero, 1, effectorColor)

		var effectorWorldSpace raylib.Vector3
		raymath.Vector3Transform(&effectorWorldSpace, transform)
		effectorWorldGround := effectorWorldSpace
		effectorWorldGround.Y = 0

		raylib.DrawLine3D(effectorWorldSpace, effectorWorldGround, effectorColor)

		plantSize := 5 * mm
		if leg.IK.Planted {
			plantSize = 20 * mm
		}
		raylib.DrawCubeV(effectorWorldGround, raylib.Vector3{plantSize, 1 * mm, plantSize}, effectorColor)

		//raylib.DrawCircle3D(leg.IK.Target.Meters(), 5*mm, raylib.Vector3{0, 0, 0}, 0, raylib.DarkGreen)
		raylib.DrawCubeV(leg.IK.Target.Meters(), raylib.Vector3{8 * mm, 1 * mm, 8 * mm}, raylib.Green)
	}
}

func (model *Model) DrawUI(camera raylib.Camera) {
	for i := range model.Labels {
		label := &model.Labels[i]
		screen := raylib.GetWorldToScreen(label.Position, camera)
		raylib.DrawText(label.Text, int32(screen.X), int32(screen.Y), 18, label.Color)
	}

	min := raylib.Vector2{10, 30}
	size := raylib.Vector2{200, 200}
	hudScale := float32(size.X) / 0.6

	center := min
	center.X += size.X / 2
	center.Y += size.Y / 2

	raylib.DrawRectangleV(min, size, raylib.Fade(raylib.SkyBlue, 0.5))

	var bodyOrigin raylib.Vector3 = model.Pose.Origin.Scale(hudScale).Meters()
	var bodySize raylib.Vector3 = model.Pose.Size.Scale(hudScale).Meters()
	bodyMin := raymath.Vector2Add(center, raylib.Vector2{bodyOrigin.Z, -bodyOrigin.X})
	bodyMin.X -= bodySize.Z / 2
	bodyMin.Y -= bodySize.X / 2

	raylib.DrawRectangleV(bodyMin, raylib.Vector2{bodySize.Z, bodySize.X}, raylib.DarkGray)

	plantedPoints := []raylib.Vector2{}

	for _, leg := range model.Pose.Legs() {
		effectorColor := raylib.Blue
		if !leg.IK.Solved {
			effectorColor = raylib.Red
		}
		if !leg.IK.Planted {
			effectorColor = raylib.Fade(effectorColor, 0.5)
		}

		footSize := raylib.Vector2{10, 10}

		var p raylib.Vector3
		p = leg.IK.Target.Scale(hudScale).Meters()

		t := raymath.Vector2Add(center, raylib.Vector2{p.Z, -p.X})
		if leg.IK.Planted {
			plantedPoints = append(plantedPoints, t)
		}

		t.X -= footSize.X / 2
		t.Y -= footSize.Y / 2

		raylib.DrawRectangleV(t, footSize, effectorColor)
	}

	if len(plantedPoints) >= 2 {
		p := plantedPoints[len(plantedPoints)-1]
		for _, n := range plantedPoints {
			raylib.DrawLineV(p, n, raylib.Blue)
			p = n
		}
	}
}
