package g_test

import (
	"testing"

	"github.com/egonelbre/hexapod/g"
)

func approxEqual(a, b g.Vec) bool {
	return a.Sub(b).Length() < 0.1*g.MM
}

func TestRotation(t *testing.T) {
	const one = 1 * g.MM

	v := g.Vec{one, one, one}

	type Case struct {
		M g.Mat
		R g.Vec
	}
	cases := []Case{
		// Zero transforms
		{g.RotateX(0), v},
		{g.RotateY(0), v},
		{g.RotateZ(0), v},

		{g.RotateX(g.Tau), v},
		{g.RotateY(g.Tau * 2), v},
		{g.RotateZ(g.Tau * 3), v},

		{g.RotateX(-g.Tau), v},
		{g.RotateY(-g.Tau * 2), v},
		{g.RotateZ(-g.Tau * 3), v},

		// half-turns
		{g.RotateX(g.Tau / 2), g.Vec{+one, -one, -one}},
		{g.RotateY(g.Tau / 2), g.Vec{-one, +one, -one}},
		{g.RotateZ(g.Tau / 2), g.Vec{-one, -one, +one}},

		{g.RotateX(-g.Tau / 2), g.Vec{+one, -one, -one}},
		{g.RotateY(-g.Tau / 2), g.Vec{-one, +one, -one}},
		{g.RotateZ(-g.Tau / 2), g.Vec{-one, -one, +one}},

		// quarter-turns
		{g.RotateX(g.Tau / 4), g.Vec{+one, +one, -one}},
		{g.RotateY(g.Tau / 4), g.Vec{-one, +one, +one}},
		{g.RotateZ(g.Tau / 4), g.Vec{+one, -one, +one}},

		{g.RotateX(-g.Tau / 4), g.Vec{+one, -one, +one}},
		{g.RotateY(-g.Tau / 4), g.Vec{+one, +one, -one}},
		{g.RotateZ(-g.Tau / 4), g.Vec{-one, +one, +one}},
	}

	for i, c := range cases {
		r := v.Transform(c.M)
		if !approxEqual(r, c.R) {
			t.Errorf("%d: got %v; exp %v", i, r, c.R)
		}
	}
}
