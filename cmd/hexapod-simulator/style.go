package main

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	HeadColor     = rl.DarkGray
	BodyColor     = rl.Gray
	LegPlateColor = rl.LightGray

	HingeColor      = rl.Blue
	HingeErrorColor = rl.NewColor(255, 102, 191, 255)

	HingeMinColor  = rl.Red
	HingeZeroColor = rl.Green
	HingeMaxColor  = rl.DarkPurple

	BoneColor = rl.SkyBlue

	EffectorTargetColor  = rl.Green
	EffectorColor        = rl.Blue
	EffectorInvalidColor = rl.Red
)
