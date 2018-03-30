package main

import (
	"github.com/egonelbre/hexapod/adeept"
	_ "github.com/egonelbre/hexapod/cmd/hexapod-simulator/internal"

	"github.com/egonelbre/hexapod/g"

	"github.com/gen2brain/raylib-go/raygui"
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

	pose := adeept.ZeroPose()
	robot := NewRobot(pose)
	model := NewModel(pose)
	minimap := NewMinimap(pose)
	minimap.Min = raylib.Vector2{10, 100}
	minimap.Size = raylib.Vector2{200, 200}

	raylib.SetTargetFPS(60)
	for !raylib.WindowShouldClose() {
		raylib.UpdateCamera(&camera)
		robot.Update(raylib.GetFrameTime())

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RayWhite)

		raylib.Begin3dMode(camera)
		{
			raylib.DrawGrid(40, 0.01*g.M.Meters())
			raylib.DrawGizmo(raylib.Vector3{})

			model.Draw()
		}
		raylib.End3dMode()

		DrawLabels3D(camera)
		minimap.Draw()

		raylib.DrawFPS(10, 10)
		if raygui.Button(raylib.Rectangle{10, 40, 20, 20}, "<") {
			robot.Toggle(-1)
		}
		raygui.Label(raylib.Rectangle{40, 40, 140, 20}, robot.ModeName())
		if raygui.Button(raylib.Rectangle{190, 40, 20, 20}, ">") {
			robot.Toggle(1)
		}

		raylib.EndDrawing()
	}

	raylib.CloseWindow()
}

type Label3D struct {
	Position raylib.Vector3
	Text     string
	Color    raylib.Color
}

var GlobalLabels []Label3D

func DrawLabel3D(text string, pos raylib.Vector3, col raylib.Color) {
	GlobalLabels = append(GlobalLabels, Label3D{pos, text, col})
}

func DrawLabels3D(camera raylib.Camera) {
	for i := range GlobalLabels {
		label := &GlobalLabels[i]
		screen := raylib.GetWorldToScreen(label.Position, camera)
		raylib.DrawText(label.Text, int32(screen.X), int32(screen.Y), 18, raylib.Fade(label.Color, 0.7))
	}
	GlobalLabels = GlobalLabels[:0]
}
