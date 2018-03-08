package main

import (
	"fmt"
	"math"
	"time"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

func ttfFromBytesMust(b []byte, size float64) font.Face {
	ttf, err := truetype.Parse(b)
	if err != nil {
		panic(err)
	}
	return truetype.NewFace(ttf, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
}

var DefaultFont = text.NewAtlas(
	ttfFromBytesMust(goregular.TTF, 14),
	text.ASCII, text.RangeTable(unicode.Latin),
)

type Length int

func (length Length) M() float64  { return length.Float64() / M }
func (length Length) MM() float64 { return length.Float64() / MM }

func (length Length) Screen() float64  { return length.Float64() / 50 }
func (length Length) Float64() float64 { return float64(length) }
func (length Length) Sqrt() Length     { return Length(math.Sqrt(float64(length))) }

const (
	M  = 100000
	MM = 100

	RadToDeg = 360 / TAU
)

type Vector struct{ X, Y Length }

func (a Vector) Add(b Vector) Vector {
	return Vector{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

func (a Vector) Sub(b Vector) Vector {
	return Vector{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}

func (a Vector) Distance(b Vector) Length {
	d := a.Sub(b)
	return (d.X*d.X + d.Y*d.Y).Sqrt()
}

func (a Vector) Pixel() pixel.Vec {
	return pixel.V(a.X.Screen(), a.Y.Screen())
}

func Angle(a, b Vector) float64 {
	d := b.Sub(a)
	return math.Atan2(d.Y.Float64(), d.X.Float64())
}

type Joint struct {
	Pos    Vector
	Angle  float64
	Length Length

	RelativeAngle float64
}

func (joint *Joint) Reach(target Vector, length Length) {
	distance := joint.Pos.Distance(target) + 1
	joint.Pos.X = target.X + (joint.Pos.X-target.X)*length/distance
	joint.Pos.Y = target.Y + (joint.Pos.Y-target.Y)*length/distance
}

func (tail *Joint) RecalculateAngle(head *Joint) {
	tail.Angle = Angle(head.Pos, tail.Pos)
}

func (tail *Joint) String() string {
	return fmt.Sprintf("<%.02f %.02f %.0f°>", tail.Pos.X.MM(), tail.Pos.Y.MM(), tail.Angle*RadToDeg)
}

type Leg struct {
	Base  Joint
	Femur Joint
	Tibia Joint

	Target Vector

	Joints [3]*Joint
}

func NewLeg(femurLength, tibiaLength Length) *Leg {
	leg := &Leg{}
	leg.Femur.Length = femurLength
	leg.Tibia.Length = tibiaLength

	leg.Joints[0] = &leg.Base
	leg.Joints[1] = &leg.Femur
	leg.Joints[2] = &leg.Tibia

	leg.Reset()
	return leg
}

func (leg *Leg) Reset() {
	for i := 0; i < len(leg.Joints)-1; i++ {
		head, tail := leg.Joints[i], leg.Joints[i+1]
		tail.Pos.X = head.Pos.X
		tail.Pos.Y = head.Pos.Y + tail.Length
	}
}

func (leg *Leg) Reach(target Vector) {
	leg.Target = target

	origin := leg.Joints[0].Pos
	n := len(leg.Joints)

	// forward reaching
	leg.Joints[n-1].Pos = target
	for i := n - 2; i >= 0; i-- {
		center, placed := leg.Joints[i], leg.Joints[i+1]
		center.Reach(placed.Pos, placed.Length)
	}

	leg.Joints[0].Pos = origin
	for i := 0; i < n-1; i++ {
		placed, center := leg.Joints[i], leg.Joints[i+1]
		center.Reach(placed.Pos, center.Length)
	}

	for i := 0; i < n-1; i++ {
		placed, center := leg.Joints[i], leg.Joints[i+1]
		center.RecalculateAngle(placed)
		center.RelativeAngle = center.Angle - placed.Angle
	}
}

func (leg *Leg) Render(draw *imdraw.IMDraw) {
	w := float64(len(leg.Joints) * 3)
	for i := 0; i < len(leg.Joints)-1; i++ {
		head := leg.Joints[i]
		tail := leg.Joints[i+1]

		draw.Color = HSL{float32(i) * math.Phi, 0.5, 0.5}
		draw.EndShape = imdraw.SharpEndShape
		draw.Push(
			head.Pos.Pixel(),
			tail.Pos.Pixel())
		draw.Line(w)
		w *= 0.8
	}

	for _, joint := range leg.Joints {
		draw.Color = RGB{0, 0, 255}

		direction := pixel.V(math.Cos(joint.Angle), math.Sin(joint.Angle))
		draw.Push(
			joint.Pos.Pixel(),
			joint.Pos.Pixel().Add(direction.Scaled(30)))
		draw.Line(2)

		t := text.New(joint.Pos.Pixel(), DefaultFont)
		t.Color = RGB{0, 0, 0}
		fmt.Fprintf(t, "%.2f", joint.Angle*RadToDeg)
		t.Draw(draw, pixel.IM)
	}

	draw.Color = HSL{TAU * 3 / 4, 0.8, 0.3}
	draw.Push(leg.Target.Pixel())
	draw.Circle(5, 0)
}

type Robot struct {
	Left  *Leg
	Right *Leg
}

func NewRobot(femurLength, tibiaLength Length) *Robot {
	robot := &Robot{}

	robot.Left = NewLeg(femurLength, tibiaLength)
	robot.Left.Base.Pos.X = -30.00 * MM
	robot.Left.Base.Pos.Y = 48.80 * MM
	robot.Left.Base.Angle = TAU / 2
	robot.Left.Reset()

	robot.Right = NewLeg(femurLength, tibiaLength)
	robot.Right.Base.Pos.Y = 48.80 * MM
	robot.Right.Base.Pos.X = 30.00 * MM
	robot.Right.Reset()

	return robot
}

func (robot *Robot) Update(t, dt float64) {
	tx := 40*MM + Length((math.Sin(t)*0.5+0.5)*50*MM)
	robot.Right.Reach(Vector{tx, 0})
	robot.Left.Reach(Vector{-tx, 0})
}

func (robot *Robot) Render(draw *imdraw.IMDraw) {
	robot.Left.Render(draw)
	robot.Right.Render(draw)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "IK-2D",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	robot := NewRobot(44.43*MM, 50.0*MM)

	start := time.Now()
	for !win.Closed() {
		win.Clear(RGB{255, 255, 255})

		now := time.Since(start).Seconds()
		robot.Update(now, 1/30)

		draw := imdraw.New(DefaultFont.Picture())

		center := win.Bounds().Size().Scaled(0.5)
		draw.SetMatrix(pixel.IM.Moved(center))

		{
			const N = 50
			for t := -N; t <= N; t++ {
				draw.Color = HSL{0, 0, 0.9}
				if t%5 == 0 {
					draw.Color = HSL{0, 0, 0.8}
				}

				draw.Push(
					Vector{Length(t) * 10 * MM, -10 * MM * N}.Pixel(),
					Vector{Length(t) * 10 * MM, 10 * MM * N}.Pixel(),
				)
				draw.Line(1)

				draw.Push(
					Vector{-10 * MM * N, Length(t) * 10 * MM}.Pixel(),
					Vector{10 * MM * N, Length(t) * 10 * MM}.Pixel(),
				)
				draw.Line(1)
			}

			// gizmo
			draw.Color = HSL{0, 0.8, 0.5}
			draw.Push(pixel.ZV, Vector{0, 50 * MM}.Pixel())
			draw.Line(2)

			draw.Color = HSL{TAU / 4, 0.8, 0.5}
			draw.Push(pixel.ZV, Vector{50 * MM, 0}.Pixel())
			draw.Line(2)
		}

		robot.Render(draw)

		draw.Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

const TAU = math.Pi * 2

type RGBA uint32
type RGB struct{ R, G, B uint8 }
type HSL struct{ H, S, L float32 }

func (rgba RGBA) RGBA() (r, g, b, a uint32) {
	r = uint32(rgba >> 24 & 0xFF)
	g = uint32(rgba >> 16 & 0xFF)
	b = uint32(rgba >> 8 & 0xFF)
	a = uint32(rgba >> 0 & 0xFF)
	r |= r << 8
	g |= g << 8
	b |= b << 8
	a |= a << 8
	return
}

func (rgb RGB) RGBA() (r, g, b, a uint32) {
	r, g, b = uint32(rgb.R), uint32(rgb.G), uint32(rgb.B)
	r |= r << 8
	g |= g << 8
	b |= b << 8
	a = 0xFFFF
	return
}

func (hsl HSL) RGBA() (r, g, b, a uint32) {
	r1, g1, b1, _ := hsla(hsl.H, hsl.S, hsl.L, 1)
	return sat16(r1), sat16(g1), sat16(b1), 0xFFFF
}

func hue(v1, v2, h float32) float32 {
	if h < 0 {
		h += 1
	}
	if h > 1 {
		h -= 1
	}
	if 6*h < 1 {
		return v1 + (v2-v1)*6*h
	} else if 2*h < 1 {
		return v2
	} else if 3*h < 2 {
		return v1 + (v2-v1)*(2.0/3.0-h)*6
	}

	return v1
}

func hsla(h, s, l, a float32) (r, g, b, ra float32) {
	if s == 0 {
		return l, l, l, a
	}

	h = float32(math.Mod(float64(h), TAU) / TAU)

	var v2 float32
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - s*l
	}

	v1 := 2*l - v2
	r = hue(v1, v2, h+1.0/3.0)
	g = hue(v1, v2, h)
	b = hue(v1, v2, h-1.0/3.0)
	ra = a

	return
}

// sat16 converts 0..1 float32 to 0..0xFFFF uint32
func sat16(v float32) uint32 {
	v = v * 0xFFFF
	if v >= 0xFFFF {
		return 0xFFFF
	} else if v <= 0 {
		return 0
	}
	return uint32(v) & 0xFFFF
}

// sat8 converts 0..1 float32 to 0..0xFF uint8
func sat8(v float32) uint8 {
	v = v * 0xFF
	if v >= 0xFF {
		return 0xFF
	} else if v <= 0 {
		return 0
	}
	return uint8(v) & 0xFF
}
