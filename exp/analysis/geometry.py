# -*- coding: utf-8 -*-

import sympy

sympy.init_printing()

def Identity():
	return sympy.Matrix([
		[1, 0, 0, 0],
		[0, 1, 0, 0],
		[0, 0, 1, 0],
		[0, 0, 0, 1]
	])


def TranslateV(v):
	x, y, z = v[0] / v[3], v[1] / v[3], v[2] / v[3]
	return Translate(x, y, z)

def Translate(x, y, z):
	if x is sympy.Matrix:
		x, y, z = x[0] / x[3], x[1] / x[3], x[2] / x[3]
	return sympy.Matrix([
		[1, 0, 0, x],
		[0, 1, 0, y],
		[0, 0, 1, z],
		[0, 0, 0, 1]
	])

def RotateX(angle):
	cs = sympy.cos(angle)
	sn = sympy.sin(angle)
	return sympy.Matrix([
		[1, 0, 0, 0],
		[0, cs, -sn, 0],
		[0, sn, cs, 0],
		[0, 0, 0, 1]
	])

def RotateY(angle):
	cs = sympy.cos(angle)
	sn = sympy.sin(angle)
	return sympy.Matrix([
		[cs, 0, sn, 0],
		[0, 1, 0, 0],
		[-sn, 0, cs, 0],
		[0, 0, 0, 1]
	])

def RotateZ(angle):
	cs = sympy.cos(angle)
	sn = sympy.sin(angle)
	return sympy.Matrix([
		[cs, -sn, 0, 0],
		[sn, cs, 0, 0],
		[0, 0, 1, 0],
		[0, 0, 0, 1]
	])

def Vector(name):
	x, y, z = sympy.symbols(name + "_x " + name + "_y " + name + "_z")
	return sympy.Matrix([x, y, z, 1])

def NegateV(v):
	vn = -v
	vn[3] = 1
	return vn


def basic():
	body = Vector("B")
	leg = Vector("L")
	coxaAngle, coxaLength = sympy.symbols("C_A C_L")
	femurAngle, femurLength = sympy.symbols("F_A F_L")
	tibiaAngle, tibiaLength = sympy.symbols("T_A T_L")

	target = sympy.simplify(
			TranslateV(body) *
			TranslateV(leg) *
			RotateY(coxaAngle) *
			Translate(coxaLength, 0, 0) *
			RotateX(femurAngle) *
			Translate(femurLength, 0, 0) *
			RotateX(tibiaAngle) *
			Translate(tibiaLength, 0, 0) *
			sympy.Matrix([0,0,0,1])
		)

	sympy.pprint(sympy.cse(sympy.solve(sympy.Eq(target, Vector("T")), coxaAngle)))

def alternate():
	yaw, pitch, roll = sympy.symbols("yaw pitch roll")
	body = Vector("B")
	leg = Vector("L")
	coxaAngle, coxaLength = sympy.symbols("C_A C_L")
	femurAngle, femurLength = sympy.symbols("F_A F_L")
	tibiaAngle, tibiaLength = sympy.symbols("T_A T_L")

	print("START")

	target = sympy.simplify(
			Translate(-tibiaLength, 0, 0) *
			RotateX(-tibiaAngle) *
			Translate(-femurLength, 0, 0) *
			RotateX(-femurAngle) *
			Translate(-coxaLength, 0, 0) *
			RotateY(-coxaAngle) *
			TranslateV(NegateV(leg)) *
			RotateZ(-roll) * RotateX(-pitch) * RotateY(-yaw) *
			TranslateV(NegateV(body)) *
			Vector("T")
		)

	print(sympy.cse(target))

	solution = sympy.cse(sympy.solve(
		sympy.Eq(target, sympy.Matrix([0,0,0,1])), 
		coxaAngle))
	sympy.pprint(solution)


alternate()

#a = sympy.symbols("a")
#sympy.pprint(sympy.simplify(RotateY(a) * RotateY(-a)))

pass
#yaw, pitch, roll = sympy.symbols("y z x")
#
#sympy.pprint(RotateY(yaw) * RotateX(pitch) * RotateZ(roll))
#sympy.pprint(RotateZ(-roll) * RotateX(-pitch) * RotateY(-yaw))
#
#x, y, z = sympy.symbols("x y z")
#sympy.pprint(Rotate())


"""
yaw, pitch, roll = sympy.symbols("yaw pitch roll")
body = Vector("B")
coxa = Vector("C")
target = Vector("T")

# TranslateV(body) * TranslateV(coxa) * Vector("relative") == Vector("target")

coxaToTarget = TranslateV(NegateV(coxa)) * TranslateV(NegateV(body)) * target

coxaAngle = sympy.atan2(coxaToTarget[2], coxaToTarget[0])
localTarget = RotateY(-coxaAngle) * coxaToTarget

sympy.pprint(RotateY(-coxaAngle) * coxaToTarget)
sympy.pprint(RotateY(+coxaAngle) * coxaToTarget)

#bodyMatrix = Translate(body[0], body[1], body[2]) * Translate(coxaX, coxaY, coxaZ)
#sympy.pprint(RotateY(yaw) * TranslateV(body) * Vector("coxa"))
"""
