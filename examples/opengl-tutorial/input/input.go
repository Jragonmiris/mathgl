package input

import (
	"github.com/Jragonmiris/mathgl"
	"github.com/go-gl/glfw"
	"math"
)

type Camera struct {
	pos            mathgl.Vec3f
	hAngle, vAngle float64
	time           float64
}

const (
	speed      float64 = 3.0
	mouseSpeed         = .005
	width              = 1024.0
	height             = 768.0
	initialFOV         = 45.0
)

func NewCamera() *Camera {
	return &Camera{pos: mathgl.Vec3f{0., 0., 5.}, hAngle: math.Pi, vAngle: 0.0, time: -1.0} // Make time -1 since it will never naturally be, this acts as a "first time?" flag
}

// Since go has multiple return values, I just went ahead and made it return the view and perspective matrices (in that order) rather than messing with getter methods
func (c *Camera) ComputeViewPerspective() (mathgl.Mat4f, mathgl.Mat4f) {
	if mathgl.FloatEqual(-1.0, c.time) {
		c.time = glfw.Time()
	}

	currTime := glfw.Time()
	deltaT := currTime - c.time

	xPos, yPos := glfw.MousePos()
	glfw.SetMousePos(width/2.0, height/2.0)

	c.hAngle += mouseSpeed * ((width / 2.0) - float64(xPos))
	c.vAngle += mouseSpeed * ((height / 2.0) - float64(yPos))

	dir := mathgl.Vec3f{
		float32(math.Cos(c.vAngle) * math.Sin(c.hAngle)),
		float32(math.Sin(c.vAngle)),
		float32(math.Cos(c.vAngle) * math.Cos(c.hAngle))}

	right := mathgl.Vec3f{
		float32(math.Sin(c.hAngle - math.Pi/2.0)),
		0.0,
		float32(math.Cos(c.hAngle - math.Pi/2.0))}

	up := right.Cross(dir)

	if glfw.Key(glfw.KeyUp) == glfw.KeyPress || glfw.Key('W') == glfw.KeyPress {
		c.pos = c.pos.Add(dir.Mul(float32(deltaT * speed)))
	}

	if glfw.Key(glfw.KeyDown) == glfw.KeyPress || glfw.Key('S') == glfw.KeyPress {
		c.pos = c.pos.Sub(dir.Mul(float32(deltaT * speed)))
	}

	if glfw.Key(glfw.KeyRight) == glfw.KeyPress || glfw.Key('D') == glfw.KeyPress {
		c.pos = c.pos.Add(right.Mul(float32(deltaT * speed)))
	}

	if glfw.Key(glfw.KeyLeft) == glfw.KeyPress || glfw.Key('A') == glfw.KeyPress {
		c.pos = c.pos.Sub(right.Mul(float32(deltaT * speed)))
	}

	// Adding to the original tutorial, Space goes up
	if glfw.Key(glfw.KeySpace) == glfw.KeyPress {
		c.pos = c.pos.Add(up.Mul(float32(deltaT * speed)))
	}

	// Adding to the original tutorial, left control goes down
	if glfw.Key(glfw.KeyLctrl) == glfw.KeyPress {
		c.pos = c.pos.Sub(up.Mul(float32(deltaT * speed)))
	}

	fov := initialFOV - 5.0*float64(glfw.MouseWheel())

	proj := mathgl.Perspective(fov, 4.0/3.0, 0.1, 100.0)
	view := mathgl.LookAtV(c.pos, c.pos.Add(dir), up)

	c.time = currTime

	return view, proj
}
