package main

import "github.com/gen2brain/raylib-go/raylib"

var (
	HeadColor     = raylib.DarkGray
	BodyColor     = raylib.Gray
	LegPlateColor = raylib.LightGray

	HingeColor      = raylib.Blue
	HingeErrorColor = raylib.NewColor(255, 102, 191, 255)

	HingeMinColor  = raylib.Red
	HingeZeroColor = raylib.Green
	HingeMaxColor  = raylib.DarkPurple

	BoneColor = raylib.SkyBlue

	EffectorTargetColor  = raylib.Green
	EffectorColor        = raylib.Blue
	EffectorInvalidColor = raylib.Red
)
