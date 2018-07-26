package gamecam

import (
	"fmt"

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
	Xoffset         float32
	Yoffset         float32
	ZoomSpeed       float32

	mode CamMode

	Zoom       int8
	ZoomStep1P int8 // negate to invert zoom direction
	ZoomStep3P int8 // negate to invert zoom direction

	EnableKeys bool
	EnableZoom bool

	enabled    bool
	rotating   bool
	subsEvents int
}

func New(icamera camera.ICamera, iwindow window.IWindow) (c *Control) {
	c = new(Control)
	c.Init(icamera, iwindow)
	return
}

func (c *Control) Dispose() {
	c.Window.UnsubscribeID(window.OnCursor, &c.subsEvents)
	c.Window.UnsubscribeID(window.OnMouseUp, &c.subsEvents)
	c.Window.UnsubscribeID(window.OnMouseDown, &c.subsEvents)
	c.Window.UnsubscribeID(window.OnScroll, &c.subsEvents)
	c.Window.UnsubscribeID(window.OnKeyDown, &c.subsEvents)
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
	c.ZoomSpeed = 0.1

	c.mode.Init(DefaultToScreen)

	c.Zoom = -0x21
	c.ZoomStep1P = 0x08
	c.ZoomStep3P = 0x08

	c.EnableKeys = true
	c.EnableZoom = true

	c.enabled = true
	c.rotating = false
	c.subsEvents = false

	c.Window.SubscribeID(window.OnCursor, &c.subsEvents, c.onMouseCursor)
	c.Window.SubscribeID(window.OnMouseUp, &c.subsEvents, c.onMouseButton)
	c.Window.SubscribeID(window.OnMouseDown, &c.subsEvents, c.onMouseButton)
	c.Window.SubscribeID(window.OnScroll, &c.subsEvents, c.onMouseScroll)
	c.Window.SubscribeID(window.OnKeyDown, &c.subsEvents, c.onKeyboardKey)
	fmt.Println("Subscribed!")
	return
}

func (c *Control) Mode() CamMode {
	return c.mode
}

func (c *Control) Reset() {
	c.SetMode(CamMode{c.Mode().ClrCopy(WorldOverrides)})
	c.SetMode(CamMode{c.Mode().SetCopy(DefaultToScreen)})
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
		c.SetMode(CamMode{c.Mode().ClrCopy(WorldOverrides)})
		c.SetMode(CamMode{c.Mode().SetCopy(DefaultToScreen)})
	}
	return
}

func (c *Control) SetMode(cm CamMode) (was CamMode) {
	was = c.mode
	c.mode = cm
	switch {
	case was.World() && cm.Screen():
		c.rotating = false
	case was.Screen() && cm.World():
		c.rotating = true
		c.rotateStart.Set(float32(c.Xoffset), float32(c.Yoffset))
	}
	return
}

func (c *Control) ZoomBySteps(step1P, step3P int) {
	old := int(c.Zoom)
	var new int
	if old >= 0 {
		new = old + step1P
		switch {
		case new < 0:
			new = -1
		case new > 0x70:
			new = 0x70
		}
	} else {
		new = old + step3P
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
	step1P := int(amount * float64(c.ZoomStep1P))
	step3P := int(amount * float64(c.ZoomStep3P))
	c.ZoomBySteps(step1P, step3P)
}

func (c *Control) ZoomOut(amount float64) {
	c.ZoomIn(-amount)
}

func (c *Control) onMouseButton(evname string, event interface{}) {
	if !c.Enabled() {
		fmt.Println("    >>> Quick return")
		return
	}
	ev := event.(*window.MouseEvent)
	if ev.Button != window.MouseButtonMiddle {
		return
	}
	switch ev.Action {
	case window.Press:
		c.SetMode(CamMode{c.Mode().SetCopy(MiddleMouseHeld)})
	case window.Release:
		c.SetMode(CamMode{c.Mode().ClrCopy(MiddleMouseHeld)})
	}
}

func (c *Control) onMouseCursor(evname string, event interface{}) {
	fmt.Println("onMouseCursor")
	ev := event.(*window.CursorEvent)
	xOffset, yOffset := ev.Xpos, ev.Ypos
	c.Xoffset, c.Yoffset = xOffset, yOffset
	if !c.rotating || !c.Enabled() || c.Mode().Screen() {
		fmt.Println("    >>> Quick return")
		return
	}
	c.rotateEnd.Set(xOffset, yOffset)
	var rotateDelta math32.Vector2 // TODO: don't use vectors for this
	rotateDelta.SubVectors(&c.rotateEnd, &c.rotateStart)
	c.rotateStart = c.rotateEnd
	width, height := c.Window.Size()
	by := 2.0 * math.Pi * float64(c.RotateSpeed)
	c.RotateLeft(by / float64(width) * float64(rotateDelta.X))
	c.RotateUp(by / float64(height) * float64(rotateDelta.Y))
}

func (c *Control) onMouseScroll(evname string, event interface{}) {
	fmt.Println("onMouseScroll")
	if !c.Enabled() || !c.EnableZoom || c.Mode().Screen() {
		fmt.Println("    >>> Quick return")
		return
	}
	ev := event.(*window.ScrollEvent)
	c.ZoomOut(float64(ev.Yoffset))
}

func (c *Control) onKeyboardKey(evname string, event interface{}) {
	fmt.Println("onKeyboardKey")
	if !c.Enabled() || !c.EnableKeys {
		fmt.Println("    >>> Quick return")
		return
	}
	ev := event.(*window.KeyEvent)
	switch ev.Keycode {
	case window.KeyLeftAlt:
		switch ev.Action {
		case window.Press:
			c.SetMode(CamMode{c.Mode().SetCopy(ScreenButtonHeld)})
		case window.Release:
			c.SetMode(CamMode{c.Mode().ClrCopy(ScreenButtonHeld)})
		}
	case window.KeyEscape:
		switch ev.Action {
		case window.Press:
			c.SetMode(CamMode{c.Mode().XorCopy(ScreenToggleOn)})
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
	fmt.Printf("updateRotate: %f, %f\n", thetaDelta, phiDelta)
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
	radius := float64(vdir.Length())
	theta := float64(math32.Atan2(vdir.X, vdir.Z)) // TODO: 64-bit
	phi := math.Acos(float64(vdir.Y) / radius)
	theta += thetaDelta
	phi += phiDelta
	theta = maths.ClampFloat64(theta, float64(c.MinAzimuthAngle), float64(c.MaxAzimuthAngle))
	phi = maths.ClampFloat64(phi, float64(min), float64(max))
	vdir.X = float32(radius * math.Sin(phi) * math.Sin(theta))
	vdir.Y = float32(radius * math.Cos(phi))
	vdir.Z = float32(radius * math.Sin(phi) * math.Cos(theta))
	vdir.ApplyQuaternion(&quatInverse)
	position = target
	position.Add(&vdir)
	c.camera.SetPositionVec(&position)
	c.camera.LookAt(&target)
}

const updateZoomEpsilon float64 = 0.01

func (c *Control) updateZoom(zoomDelta float64) {
	if ortho, ok := c.icamera.(*camera.Orthographic); ok {
		zoom := float64(ortho.Zoom()) - updateZoomEpsilon*zoomDelta
		ortho.SetZoom(float32(zoom))
	} else {
		fmt.Printf("updateZoom:else: %f\n", zoomDelta)
		position := c.camera.Position()
		target := c.camera.Target()
		vdir := position
		vdir.Sub(&target)
		dist := float64(vdir.Length()) * (1.0 + zoomDelta*float64(c.ZoomSpeed))
		dist = maths.ClampFloat64(dist, float64(c.MinDistance), float64(c.MaxDistance))
		vdir.SetLength(float32(dist))
		target.Add(&vdir)
		c.camera.SetPositionVec(&target)
	}
}
