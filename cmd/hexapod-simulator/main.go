package main

import (
	"github.com/egonelbre/hexapod/adeept"
	_ "github.com/egonelbre/hexapod/cmd/hexapod-simulator/internal"

	"github.com/egonelbre/hexapod/g"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	const screenWidth = 1024
	const screenHeight = 768

	rl.SetConfigFlags(rl.FlagMsaa4xHint)

	rl.InitWindow(screenWidth, screenHeight, "Hexapod Simulator")
	// rl.SetWindowPosition(-1920/2-screenWidth/2, 1200/2-screenHeight/2)

	camera := rl.Camera{}
	camera.Position = rl.NewVector3(0.5, 0.5, 0.5)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 30.0

	pose := adeept.ZeroPose()
	robot := NewRobot(pose)
	model := NewModel(pose)
	minimap := NewMinimap(pose)
	minimap.Min = rl.Vector2{10, 100}
	minimap.Size = rl.Vector2{200, 200}

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		//rl.UpdateCamera(&camera, rl.CameraFree)
		rl.UpdateCamera(&camera, rl.CameraOrbital)
		robot.Update(rl.GetFrameTime())

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)
		{
			rl.DrawGrid(40, 0.01*g.M.Meters())

			model.Draw()
		}
		rl.EndMode3D()

		DrawLabels3D(camera)
		minimap.Draw()

		rl.DrawFPS(10, 10)
		if raygui.Button(rl.Rectangle{10, 40, 20, 20}, "<") {
			robot.Toggle(-1)
		}
		raygui.Label(rl.Rectangle{40, 40, 140, 20}, robot.ModeName())
		if raygui.Button(rl.Rectangle{190, 40, 20, 20}, ">") {
			robot.Toggle(1)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

type Label3D struct {
	Position rl.Vector3
	Text     string
	Color    rl.Color
}

var GlobalLabels []Label3D

func DrawLabel3D(text string, pos rl.Vector3, col rl.Color) {
	GlobalLabels = append(GlobalLabels, Label3D{pos, text, col})
}

func DrawLabels3D(camera rl.Camera) {
	for i := range GlobalLabels {
		label := &GlobalLabels[i]
		screen := rl.GetWorldToScreen(label.Position, camera)
		rl.DrawText(label.Text, int32(screen.X), int32(screen.Y), 18, rl.Fade(label.Color, 0.7))
	}
	GlobalLabels = GlobalLabels[:0]
}
