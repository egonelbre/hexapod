package main

import (
	_ "github.com/egonelbre/hexapod/cmd/hexapod-simulator/internal"

	"github.com/egonelbre/hexapod/adeept"
	"github.com/egonelbre/hexapod/g"

	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1024)
	screenHeight := int32(768)

	raylib.SetConfigFlags(raylib.FlagMsaa4xHint)

	raylib.InitWindow(screenWidth, screenHeight, "Hexapod Simulator")
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

			model.Update()
			model.Draw()
		}
		raylib.End3dMode()

		model.DrawUI(camera)
		raylib.DrawFPS(10, 10)

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}
