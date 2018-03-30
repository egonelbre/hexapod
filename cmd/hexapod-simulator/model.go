package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
	"github.com/gen2brain/raylib-go/raymath"

	"github.com/egonelbre/hexapod/g"
	"github.com/egonelbre/hexapod/pose"
)

type Model struct {
	Pose *pose.Body

	Shader raylib.Shader

	Head raylib.Model
	Body raylib.Model

	LegPlateSize g.Vec
	LegPlate     raylib.Model

	Bone     raylib.Model
	Hinge    raylib.Model
	Effector raylib.Model
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

func matmul(m raylib.Matrix, xs ...raylib.Matrix) raylib.Matrix {
	for i := range xs {
		m = raymath.MatrixMultiply(m, xs[i])
	}
	return m
}

func (model *Model) Draw() {
	DrawLabel3D(".", model.Pose.Origin.Meters(), raylib.Black)

	zero := raylib.Vector3{}

	bodyTransform := matmul(
		raymath.MatrixTranslate(model.Pose.Origin.XYZ()),
		raymath.MatrixRotateY(model.Pose.Orient.Yaw),
		raymath.MatrixRotateZ(model.Pose.Orient.Pitch),
		raymath.MatrixRotateX(model.Pose.Orient.Roll),
	)

	model.Head.Transform = raymath.MatrixMultiply(bodyTransform, raymath.MatrixTranslate(model.Pose.Head.Offset.XYZ()))
	model.Body.Transform = bodyTransform
	raylib.DrawModel(model.Head, zero, 1, HeadColor)
	raylib.DrawModel(model.Body, zero, 1, BodyColor)

	zmm := (model.Pose.Size.Z/2 + model.LegPlateSize.Z/2).Meters()
	ymm := model.Pose.Leg.RF.Offset.Y.Meters()
	model.LegPlate.Transform = raymath.MatrixMultiply(bodyTransform, raymath.MatrixTranslate(0, ymm, +zmm))
	raylib.DrawModel(model.LegPlate, zero, 1, LegPlateColor)
	model.LegPlate.Transform = raymath.MatrixMultiply(bodyTransform, raymath.MatrixTranslate(0, ymm, -zmm))
	raylib.DrawModel(model.LegPlate, zero, 1, LegPlateColor)

	mm := g.MM.Meters()

	for _, leg := range model.Pose.Legs() {
		transform := raymath.MatrixMultiply(bodyTransform, raymath.MatrixTranslate(leg.Offset.XYZ()))

		var labelPosition raylib.Vector3
		raymath.Vector3Transform(&labelPosition, transform)
		labelPosition.Y += 30 * mm
		DrawLabel3D(leg.Name+" "+leg.IK.Debug, labelPosition, raylib.Black)

		for _, hinge := range leg.Hinges() {
			var newRotation func(v float32) raylib.Matrix
			var hingeScale raylib.Matrix
			switch hinge.Axis {
			case pose.X:
				newRotation = raymath.MatrixRotateX
				hingeScale = raymath.MatrixScale(25*mm, 5*mm, 5*mm)
			case pose.Y:
				newRotation = raymath.MatrixRotateY
				hingeScale = raymath.MatrixScale(5*mm, 25*mm, 5*mm)
			case pose.Z:
				newRotation = raymath.MatrixRotateZ
				hingeScale = raymath.MatrixScale(5*mm, 5*mm, 25*mm)
			}

			var hingeCenter raylib.Vector3
			raymath.Vector3Transform(&hingeCenter, transform)
			DrawLabel3D(fmt.Sprintf("%.0f", hinge.Angle*g.RadToDeg), hingeCenter, raylib.Black)

			hingePointMin := raylib.Vector3{30 * mm, 0, 0}
			hingePointZero := raylib.Vector3{30 * mm, 0, 0}
			hingePointMax := raylib.Vector3{30 * mm, 0, 0}

			raymath.Vector3Transform(&hingePointMin, raymath.MatrixMultiply(transform, newRotation(hinge.Zero+hinge.Range.Min)))
			raymath.Vector3Transform(&hingePointZero, raymath.MatrixMultiply(transform, newRotation(hinge.Zero)))
			raymath.Vector3Transform(&hingePointMax, raymath.MatrixMultiply(transform, newRotation(hinge.Zero+hinge.Range.Max)))

			raylib.DrawLine3D(hingeCenter, hingePointMin, HingeMinColor)
			raylib.DrawLine3D(hingeCenter, hingePointZero, HingeZeroColor)
			raylib.DrawLine3D(hingeCenter, hingePointMax, HingeMaxColor)

			transform = raymath.MatrixMultiply(transform, newRotation(hinge.Zero+hinge.Angle))
			model.Hinge.Transform = raymath.MatrixMultiply(transform, hingeScale)

			hingeLength := hinge.Length.Meters()

			boneTransform := raymath.MatrixMultiply(transform,
				raymath.MatrixTranslate(hingeLength/2, 0, 0))

			if hinge == &leg.Tibia {
				model.Bone.Transform = raymath.MatrixMultiply(boneTransform,
					raymath.MatrixScale(hingeLength, 8*mm, 8*mm))
			} else {
				model.Bone.Transform = raymath.MatrixMultiply(boneTransform,
					raymath.MatrixScale(hingeLength, 12*mm, 20*mm))
			}

			if hinge.InBounds() {
				raylib.DrawModel(model.Hinge, zero, 1, HingeColor)
			} else {
				raylib.DrawModel(model.Hinge, zero, 1, HingeErrorColor)
			}

			raylib.DrawModel(model.Bone, zero, 1, BoneColor)

			transform = raymath.MatrixMultiply(transform, raymath.MatrixTranslate(hingeLength, 0, 0))
		}

		effectorColor := EffectorColor
		if !leg.IK.Solved {
			effectorColor = EffectorInvalidColor
		}
		model.Effector.Transform = raymath.MatrixMultiply(transform,
			raymath.MatrixScale(3*mm, 10*mm, 10*mm))
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

		raylib.DrawCubeV(leg.IK.Target.Meters(), raylib.Vector3{8 * mm, 1 * mm, 8 * mm}, EffectorTargetColor)
	}
}
