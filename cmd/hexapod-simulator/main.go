package main

import (
	_ "github.com/egonelbre/hexapod/cmd/hexapod-simulator/internal"
	"github.com/egonelbre/hexapod/pose"

	"github.com/egonelbre/hexapod/adeept"
	"github.com/egonelbre/hexapod/g"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1024)
	screenHeight := int32(768)

	raylib.SetConfigFlags(raylib.FlagMsaa4xHint)

	raylib.InitWindow(screenWidth, screenHeight, "ik3d")
	raylib.SetWindowPosition(-1920/2-screenWidth/2, 1200/2-screenHeight/2)

	camera := raylib.Camera{}
	camera.Position = raylib.NewVector3(0.5, 0.5, 0.5)
	camera.Target = raylib.NewVector3(0.0, 0.0, 0.0)
	camera.Up = raylib.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 30.0

	raylib.SetCameraMode(camera, raylib.CameraFree)
	//raylib.SetCameraMode(camera, raylib.CameraOrbital)

	model := NewModel(adeept.ZeroPose())

	raylib.SetTargetFPS(60)
	for !raylib.WindowShouldClose() {
		raylib.UpdateCamera(&camera)

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RayWhite)

		raylib.Begin3dMode(camera)
		{
			raylib.DrawGrid(40, 0.01*g.M.Meters())
			raylib.DrawGizmo(raylib.Vector3{})

			model.Draw()
		}
		raylib.End3dMode()

		{
			raylib.DrawFPS(10, 10)

			raylib.DrawRectangle(10, 30, 320, 133, raylib.Fade(raylib.SkyBlue, 0.5))
			raylib.DrawRectangleLines(10, 30, 320, 133, raylib.Blue)

			raylib.DrawText("Free camera default controls:", 20, 40, 10, raylib.Black)
			raylib.DrawText("- Mouse Wheel to Zoom in-out", 40, 60, 10, raylib.DarkGray)
			raylib.DrawText("- Mouse Wheel Pressed to Pan", 40, 80, 10, raylib.DarkGray)
			raylib.DrawText("- Alt + Mouse Wheel Pressed to Rotate", 40, 100, 10, raylib.DarkGray)
			raylib.DrawText("- Alt + Ctrl + Mouse Wheel Pressed for Smooth Zoom", 40, 120, 10, raylib.DarkGray)
			raylib.DrawText("- Z to zoom to (0, 0, 0)", 40, 140, 10, raylib.DarkGray)
		}

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}

type Model struct {
	Pose *pose.Body
}

func NewModel(pose *pose.Body) *Model {
	model := &Model{}
	model.Pose = pose
	return model
}

func (model *Model) Draw() {
	//
}
