package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/gen2brain/raylib-go/raymath"

	"github.com/egonelbre/hexapod/pose"
)

type Minimap struct {
	Body *pose.Body

	Min  raylib.Vector2
	Size raylib.Vector2
}

func NewMinimap(body *pose.Body) *Minimap {
	minimap := &Minimap{}
	minimap.Body = body
	return minimap
}

func (minimap *Minimap) Draw() {
	body := minimap.Body
	hudScale := float32(minimap.Size.X) / 0.6

	center := minimap.Min
	center.X += minimap.Size.X / 2
	center.Y += minimap.Size.Y / 2

	raylib.DrawRectangleV(minimap.Min, minimap.Size, raylib.Fade(raylib.SkyBlue, 0.5))

	var bodyOrigin raylib.Vector3 = body.Origin.Scale(hudScale).Meters()
	var bodySize raylib.Vector3 = body.Size.Scale(hudScale).Meters()
	bodyMin := raymath.Vector2Add(center, raylib.Vector2{bodyOrigin.Z, -bodyOrigin.X})
	bodyMin.X -= bodySize.Z / 2
	bodyMin.Y -= bodySize.X / 2

	raylib.DrawRectangleV(bodyMin, raylib.Vector2{bodySize.Z, bodySize.X}, raylib.DarkGray)

	plantedPoints := []raylib.Vector2{}

	for _, leg := range body.Legs() {
		effectorColor := EffectorColor
		if !leg.IK.Solved {
			effectorColor = EffectorInvalidColor
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
