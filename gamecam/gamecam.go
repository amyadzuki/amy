package gamecam

import (
	"github.com/amyadzuki/amystuff/maths"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"math"
)

type Control struct {
	camera  *camera.Camera
	icamera camera.ICamera
	Window  window.IWindow

	position0 math32.Vector3
	target0   math32.Vector3

	rotateStart math32.Vector2
	rotateEnd   math32.Vector2

	MaxAzimuthAngle float32
	MaxDistance     float32
	MaxPolarAngle   float32
	MinAzimuthAngle float32
	MinDistance     float32
	MinPolarAngle   float32
	RotateSpeed     float32
	ZoomSpeed       float32

	mode CamMode

	Zoom       int8
	ZoomStep1P int8 // negate to invert zoom direction
	ZoomStep3P int8 // negate to invert zoom direction

	EnableKeys bool
	EnableZoom bool

	enabled    bool
	rotating   bool
	subsCursor bool
	subsEvents bool
}

const ModeWorld uint8 = 0
const (
	// Controlling the screen for another reason, such as in a menu
	ModeScreenOther uint8 = iota << 0

	// Controlling the screen because a key is held
	ModeScreenHeld

	// Controlling the screen because a key was toggled on
	ModeScreenToggled
)

func New(icamera camera.ICamera, iwindow window.IWindow) (c *Control) {
	c = new(Control)
	c.Init(icamera, iwindow)
	return
}

func (c *Control) Dispose() {
	c.Window.UnsubscribeID(window.OnMouseUp, &c.subsEvents)
	c.Window.UnsubscribeID(window.OnMouseDown, &c.subsEvents)
	c.Window.UnsubscribeID(window.OnScroll, &c.subsEvents)
	c.Window.UnsubscribeID(window.OnKeyDown, &c.subsEvents)
	c.Window.UnsubscribeID(window.OnCursor, &c.subsCursor)
}

func (c *Control) Enabled() bool {
	return c.enabled
}

func (c *Control) Init(icamera camera.ICamera, iwindow window.IWindow) {
	c.camera = icamera.GetCamera()
	c.icamera = icamera
	c.Window = iwindow

	c.position0 = c.camera.Position()
	c.target0 = c.camera.Target()

	c.rotateStart = math32.Vector2{0, 0}
	c.rotateEnd = math32.Vector2{0, 0}

	c.MaxAzimuthAngle = float32(math.Inf(+1))
	c.MaxDistance = float32(math.Inf(+1))
	c.MaxPolarAngle = float32(math.Pi)
	c.MinAzimuthAngle = float32(math.Inf(-1))
	c.MinDistance = 0.01
	c.MinPolarAngle = 0.0
	c.RotateSpeed = 1.0
	c.ZoomSpeed = 1.0

	c.mode = CamMode(DefaultToScreen)

	c.Zoom = -0x21
	c.ZoomStep1P = 0x08
	c.ZoomStep3P = 0x08

	c.EnableKeys = true
	c.EnableZoom = true

	c.enabled = true
	c.rotating = false
	c.subsCursor = false
	c.subsEvents = false

	c.Window.SubscribeID(window.OnMouseUp, &c.subsEvents, c.onMouseButton)
	c.Window.SubscribeID(window.OnMouseDown, &c.subsEvents, c.onMouseButton)
	c.Window.SubscribeID(window.OnScroll, &c.subsEvents, c.onMouseScroll)
	c.Window.SubscribeID(window.OnKeyDown, &c.subsEvents, c.onKeyboardKey)
	return
}

func (c *Control) Mode() CamMode {
	return c.mode
}

func (c *Control) Reset() {
	c.SetMode(c.Mode().ClrCopy(WorldOverrides))
	c.SetMode(c.Mode().SetCopy(DefaultToScreen))
	c.camera.SetPositionVec(&c.position0)
	c.camera.LookAt(&c.target0)
}

func (c *Control) RotateDown(amount float64) {
	c.updateRotate(0, amount)
}

func (c *Control) RotateLeft(amount float64) {
	c.RotateRight(-amount)
}

func (c *Control) RotateRight(amount float64) {
	c.updateRotate(amount, 0)
}

func (c *Control) RotateUp(amount float64) {
	c.RotateDown(-amount)
}

func (c *Control) SetEnabled(enabled bool) (was bool) {
	was = c.enabled
	c.enabled = enabled
	if !enabled {
		c.SetMode(c.Mode().ClrCopy(WorldOverrides))
		c.SetMode(c.Mode().SetCopy(DefaultToScreen))
	}
	return
}

func (c *Control) SetMode(cm CamMode) (was CamMode) {
	was = c.mode
	c.mode = cm
	switch {
	case was.World() && cm.Screen():
		c.rotating = false
		c.Window.UnsubscribeID(window.OnCursor, &c.subsCursor)
	case was.Screen() && cm.World():
		c.rotating = true
		c.rotateStart.Set(float32(ev.Xpos), float32(ev.Ypos))
		c.Window.SubscribeID(window.OnCursor, &c.subsCursor, c.onMouseCursor)
	}
	return
}

func (c *Control) ZoomBySteps(step1P, step3P int) {
	old := int(c.Zoom)
	if old >= 0 {
		new := old + step1P
		switch {
		case new < 0:
			new = -1
		case new > 0x70:
			new = 0x70
		}
	} else {
		new := old + step3P
		switch {
		case new >= 0:
			new = 0
		case new < -0x71:
			new = -0x71
		}
	}
	c.Zoom = int8(new)
}

func (c *Control) ZoomIn(amount float64) {
	step1P = int(amount * float64(c.ZoomStep1P))
	step3P = int(amount * float64(c.ZoomStep3P))
	c.ZoomBySteps(step1P, step3P)
}

func (c *Control) ZoomOut(amount float64) {
	c.ZoomIn(-amount)
}

func (c *Control) onMouseButton(evname string, event interface{}) {
	if !c.Enabled() {
		return
	}
	ev := event.(*window.MouseEvent)
	if ev.Button != window.MouseButtonMiddle {
		return
	}
	switch ev.Action {
	case window.Press:
		c.SetMode(c.Mode().SetCopy(MiddleMouseHeld))
	case window.Release:
		c.SetMode(c.Mode().ClrCopy(MiddleMouseHeld))
	}
}

func (c *Control) onMouseCursor(evname string, event interface{}) {
	if !c.Enabled() || c.Mode().Screen() {
		return
	}
	ev := event.(*window.CursorEvent)
	c.rotateEnd.Set(float32(ev.Xpos), float32(ev.Ypos))
	c.rotateDelta.SubVectors(&c.rotateEnd, &c.rotateStart)
	c.rotateStart = c.rotateEnd
	width, height := c.Window.Size()
	by := 2.0 * math.Pi * c.RotateSpeed
	c.RotateLeft(by / float64(width) * float64(c.rotateDelta.X))
	c.RotateUp(by / float64(height) * float64(c.rotateDelta.Y))
}

func (c *Control) onMouseScroll(evname string, event interface{}) {
	if !c.Enabled() || !c.EnableZoom || c.Mode().Screen() {
		return
	}
	ev := event.(*window.ScrollEvent)
	c.ZoomOut(float64(ev.Yoffset))
}

func (c *Control) onKeyboardKey(evname string, event interface{}) {
	if !c.Enabled() || !c.EnableKeys {
		return
	}
	ev := event.(*window.KeyEvent)
	switch ev.KeyCode {
	case window.KeyLeftAlt:
		switch ev.Action {
		case window.Press:
			c.SetMode(c.Mode().SetCopy(ScreenButtonHeld))
		case window.Release:
			c.SetMode(c.Mode().ClrCopy(ScreenButtonHeld))
		}
	case window.KeyEscape:
		switch ev.Action {
		case window.Press:
			c.SetMode(c.Mode().XorCopy(ScreenToggleOn))
		case window.Release:
		}
	case window.KeyHome:
		switch ev.Action {
		case window.Press:
			//          c.Snap() // TODO:
		case window.Release:
		}
	}
}

const updateRotateEpsilon float64 = 0.01
const updateRotatePiMinusEpsilon float64 = math.Pi - updateRotateEpsilon

func (c *Control) updateRotate(thetaDelta, phiDelta float64) {
	var max, min float64
	if float64(c.MaxPolarAngle) < updateRotatePiMinusEpsilon {
		max = float64(c.MaxPolarAngle)
	} else {
		max = updateRotatePiMinusEpsilon
	}
	if float64(c.MinPolarAngle) > updateRotateEpsilon {
		min = float64(c.MinPolarAngle)
	} else {
		min = updateRotateEpsilon
	}
	position := c.camera.Position()
	target := c.camera.Target()
	up := c.camera.Up()
	vdir := position
	vdir.Sub(&target)
	var quat math32.Quaternion
	quat.SetFromUnitVectors(&up, &math32.Vector3{0, 1, 0})
	quatInverse := quat
	quatInverse.Inverse()
	vdir.ApplyQuaternion(&quat)
	radius := vdir.Length()
	theta := math32.Atan2(vdir.X, vdir.Z)
	phi := math32.Acos(vdir.Y / radius)
	theta += thetaDelta
	phi += phiDelta
	theta = maths.ClampFloat64(theta, c.MinAzimuthAngle, c.MaxAzimuthAngle)
	phi = maths.ClampFloat64(phi, min, max)
	vdir.X = radius * math.Sin(phi) * math.Sin(theta)
	vdir.Y = radius * math.Cos(phi)
	vdir.Z = radius * math.Sin(phi) * math.Cos(theta)
	vdir.ApplyQuaternion(&quatInverse)
	position = target
	position.Add(&vdir)
	c.camera.SetPositionVec(&position)
	c.camera.LookAt(&target)
}

const updateZoomEpsilon float64 = 0.01

func (c *Control) updateZoom(zoomDelta float64) {
	if ortho, ok := c.icamera.(*camera.Orthographic); ok {
		zoom := ortho.Zoom() - updateZoomEpsilon*zoomDelta
		ortho.SetZoom(zoom)
	} else {
		position := c.camera.Position()
		target := c.camera.Target()
		vdir := position
		vdir.Sub(&target)
		dist := float64(vdir.Length()) * (1.0 + zoomDelta*c.ZoomSpeed/10.0)
		dist = maths.ClampFloat64(dist, c.MinDistance, c.MaxDistance)
		vdir.SetLength(float32(dist))
		target.Add(&vdir)
		c.camera.SetPositionVec(&target)
	}
}
