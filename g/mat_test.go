package g_test

import (
	"testing"
	"testing/quick"

	"github.com/egonelbre/hexapod/g"
)

func approxEqualVec(a, b g.Vec) bool {
	return a.Sub(b).Length() < 0.1*g.MM
}

func approxEqualMat(a, b g.Mat) bool {
	return a.Sub(b).Abs().GrandSum() < 1e-6
}

func TestRotate(t *testing.T) {
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
		r := c.M.Transform(v)
		if !approxEqualVec(r, c.R) {
			t.Errorf("%d: got %v; exp %v", i, r, c.R)
		}
	}
}

func TestRotateX(t *testing.T) {
	f := func(x int) bool {
		r := float32(x) / 1e6
		w := g.RotateX(r).Mul(g.RotateX(-r))
		return approxEqualMat(w, g.Identity())
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestRotateY(t *testing.T) {
	f := func(x int) bool {
		r := float32(x) / 1e6
		w := g.RotateY(r).Mul(g.RotateY(-r))
		return approxEqualMat(w, g.Identity())
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestRotateZ(t *testing.T) {
	f := func(x int) bool {
		r := float32(x) / 1e6
		w := g.RotateZ(r).Mul(g.RotateZ(-r))
		return approxEqualMat(w, g.Identity())
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
