package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/egonelbre/hexapod/g"
	"github.com/egonelbre/hexapod/pose"
)

type Model struct {
	Pose *pose.Body

	Shader rl.Shader

	Head rl.Model
	Body rl.Model

	LegPlateSize g.Vec
	LegPlate     rl.Model

	Bone     rl.Model
	Hinge    rl.Model
	Effector rl.Model
}

func NewModel(pose *pose.Body) *Model {
	model := &Model{}
	model.Pose = pose

	model.Shader = rl.LoadShader("shader/base.vs", "shader/lighting.fs")
	model.Shader.UpdateLocation(rl.ShaderLocMatrixView, rl.GetShaderLocation(model.Shader, "view"))
	model.Shader.UpdateLocation(rl.ShaderLocVectorView, rl.GetShaderLocation(model.Shader, "viewPos"))
	model.Shader.UpdateLocation(rl.ShaderLocMatrixModel, rl.GetShaderLocationAttrib(model.Shader, "mMatrix"))

	size := pose.Size.Meters()
	mm := g.MM.Meters()

	model.Head = rl.LoadModelFromMesh(rl.GenMeshCube(30*mm, 30*mm, 30*mm))
	model.Body = rl.LoadModelFromMesh(rl.GenMeshCube(size.X, size.Y, size.Z))

	model.LegPlateSize = g.Vec{
		X: pose.Leg.RF.Offset.X - pose.Leg.RB.Offset.X - 5*g.MM,
		Y: 3 * g.MM,
		Z: pose.Leg.RM.Offset.Z - pose.Size.Z/2 - 2*g.MM,
	}
	model.LegPlate = rl.LoadModelFromMesh(rl.GenMeshCube(model.LegPlateSize.XYZ()))

	// size is assigned dynamically with scaling
	model.Bone = rl.LoadModelFromMesh(rl.GenMeshCube(1, 1, 1))
	model.Hinge = rl.LoadModelFromMesh(rl.GenMeshCube(1, 1, 1))
	model.Effector = rl.LoadModelFromMesh(rl.GenMeshCube(1, 1, 1))

	model.Head.Materials.Shader = model.Shader
	model.Body.Materials.Shader = model.Shader

	model.LegPlate.Materials.Shader = model.Shader

	model.Bone.Materials.Shader = model.Shader
	model.Hinge.Materials.Shader = model.Shader
	model.Effector.Materials.Shader = model.Shader

	return model
}

func matmul(m rl.Matrix, xs ...rl.Matrix) rl.Matrix {
	for i := range xs {
		m = rl.MatrixMultiply(m, xs[i])
	}
	return m
}

func (model *Model) Draw() {
	DrawLabel3D(".", model.Pose.Origin.Meters(), rl.Black)

	zero := rl.Vector3{}

	bodyTransform := matmul(
		rl.MatrixTranslate(model.Pose.Origin.XYZ()),
		rl.MatrixRotateY(model.Pose.Orient.Yaw),
		rl.MatrixRotateZ(model.Pose.Orient.Pitch),
		rl.MatrixRotateX(model.Pose.Orient.Roll),
	)

	model.Head.Transform = rl.MatrixMultiply(bodyTransform, rl.MatrixTranslate(model.Pose.Head.Offset.XYZ()))
	model.Body.Transform = bodyTransform
	rl.DrawModel(model.Head, zero, 1, HeadColor)
	rl.DrawModel(model.Body, zero, 1, BodyColor)

	zmm := (model.Pose.Size.Z/2 + model.LegPlateSize.Z/2).Meters()
	ymm := model.Pose.Leg.RF.Offset.Y.Meters()
	model.LegPlate.Transform = rl.MatrixMultiply(bodyTransform, rl.MatrixTranslate(0, ymm, +zmm))
	rl.DrawModel(model.LegPlate, zero, 1, LegPlateColor)
	model.LegPlate.Transform = rl.MatrixMultiply(bodyTransform, rl.MatrixTranslate(0, ymm, -zmm))
	rl.DrawModel(model.LegPlate, zero, 1, LegPlateColor)

	mm := g.MM.Meters()

	for _, leg := range model.Pose.Legs() {
		transform := rl.MatrixMultiply(bodyTransform, rl.MatrixTranslate(leg.Offset.XYZ()))

		labelPosition := rl.Vector3Transform(rl.Vector3Zero(), transform)
		labelPosition.Y += 30 * mm
		DrawLabel3D(leg.Name+" "+leg.IK.Debug, labelPosition, rl.Black)

		for _, hinge := range leg.Hinges() {
			var newRotation func(v float32) rl.Matrix
			var hingeScale rl.Matrix
			switch hinge.Axis {
			case pose.X:
				newRotation = rl.MatrixRotateX
				hingeScale = rl.MatrixScale(25*mm, 5*mm, 5*mm)
			case pose.Y:
				newRotation = rl.MatrixRotateY
				hingeScale = rl.MatrixScale(5*mm, 25*mm, 5*mm)
			case pose.Z:
				newRotation = rl.MatrixRotateZ
				hingeScale = rl.MatrixScale(5*mm, 5*mm, 25*mm)
			}

			hingeCenter := rl.Vector3Transform(rl.Vector3Zero(), transform)
			DrawLabel3D(fmt.Sprintf("%.0f", hinge.Angle*g.RadToDeg), hingeCenter, rl.Black)

			hingePointMin := rl.Vector3Transform(rl.Vector3{30 * mm, 0, 0}, rl.MatrixMultiply(transform, newRotation(hinge.Zero+hinge.Range.Min)))
			hingePointZero := rl.Vector3Transform(rl.Vector3{30 * mm, 0, 0}, rl.MatrixMultiply(transform, newRotation(hinge.Zero)))
			hingePointMax := rl.Vector3Transform(rl.Vector3{30 * mm, 0, 0}, rl.MatrixMultiply(transform, newRotation(hinge.Zero+hinge.Range.Max)))

			rl.DrawLine3D(hingeCenter, hingePointMin, HingeMinColor)
			rl.DrawLine3D(hingeCenter, hingePointZero, HingeZeroColor)
			rl.DrawLine3D(hingeCenter, hingePointMax, HingeMaxColor)

			transform = rl.MatrixMultiply(transform, newRotation(hinge.Zero+hinge.Angle))
			model.Hinge.Transform = rl.MatrixMultiply(transform, hingeScale)

			hingeLength := hinge.Length.Meters()

			boneTransform := rl.MatrixMultiply(transform,
				rl.MatrixTranslate(hingeLength/2, 0, 0))

			if hinge == &leg.Tibia {
				model.Bone.Transform = rl.MatrixMultiply(boneTransform,
					rl.MatrixScale(hingeLength, 8*mm, 8*mm))
			} else {
				model.Bone.Transform = rl.MatrixMultiply(boneTransform,
					rl.MatrixScale(hingeLength, 12*mm, 20*mm))
			}

			if hinge.InBounds() {
				rl.DrawModel(model.Hinge, zero, 1, HingeColor)
			} else {
				rl.DrawModel(model.Hinge, zero, 1, HingeErrorColor)
			}

			rl.DrawModel(model.Bone, zero, 1, BoneColor)

			transform = rl.MatrixMultiply(transform, rl.MatrixTranslate(hingeLength, 0, 0))
		}

		effectorColor := EffectorColor
		if !leg.IK.Solved {
			effectorColor = EffectorInvalidColor
		}
		model.Effector.Transform = rl.MatrixMultiply(transform,
			rl.MatrixScale(3*mm, 10*mm, 10*mm))
		rl.DrawModel(model.Effector, zero, 1, effectorColor)

		effectorWorldSpace := rl.Vector3Transform(rl.Vector3Zero(), transform)
		effectorWorldGround := effectorWorldSpace
		effectorWorldGround.Y = 0

		rl.DrawLine3D(effectorWorldSpace, effectorWorldGround, effectorColor)

		plantSize := 5 * mm
		if leg.IK.Planted {
			plantSize = 20 * mm
		}
		rl.DrawCubeV(effectorWorldGround, rl.Vector3{plantSize, 1 * mm, plantSize}, effectorColor)

		rl.DrawCubeV(leg.IK.Target.Meters(), rl.Vector3{8 * mm, 1 * mm, 8 * mm}, EffectorTargetColor)
	}
}
