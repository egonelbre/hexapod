package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/egonelbre/hexapod/pose"
)

type Minimap struct {
	Body *pose.Body

	Min  rl.Vector2
	Size rl.Vector2
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

	rl.DrawRectangleV(minimap.Min, minimap.Size, rl.Fade(rl.SkyBlue, 0.5))

	var bodyOrigin rl.Vector3 = body.Origin.Scale(hudScale).Meters()
	var bodySize rl.Vector3 = body.Size.Scale(hudScale).Meters()
	bodyMin := rl.Vector2Add(center, rl.Vector2{bodyOrigin.Z, -bodyOrigin.X})
	bodyMin.X -= bodySize.Z / 2
	bodyMin.Y -= bodySize.X / 2

	rl.DrawRectangleV(bodyMin, rl.Vector2{bodySize.Z, bodySize.X}, rl.DarkGray)

	plantedPoints := []rl.Vector2{}

	for _, leg := range body.Legs() {
		effectorColor := EffectorColor
		if !leg.IK.Solved {
			effectorColor = EffectorInvalidColor
		}
		if !leg.IK.Planted {
			effectorColor = rl.Fade(effectorColor, 0.5)
		}

		footSize := rl.Vector2{10, 10}

		var p rl.Vector3
		p = leg.IK.Target.Scale(hudScale).Meters()

		t := rl.Vector2Add(center, rl.Vector2{p.Z, -p.X})
		if leg.IK.Planted {
			plantedPoints = append(plantedPoints, t)
		}

		t.X -= footSize.X / 2
		t.Y -= footSize.Y / 2

		rl.DrawRectangleV(t, footSize, effectorColor)
	}

	if len(plantedPoints) >= 2 {
		p := plantedPoints[len(plantedPoints)-1]
		for _, n := range plantedPoints {
			rl.DrawLineV(p, n, rl.Blue)
			p = n
		}
	}
}
