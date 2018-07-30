package gamecam

import (
	"math"

	"github.com/amyadzuki/amygolib/maths"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"

	"github.com/LK4D4/trylock"
)

type Control struct {
	Followee Followee

	IWindow window.IWindow

	mutexMouseCursor trylock.Mutex // sizes: sync.Mutex: 8, trylock.Mutex: 8

	Gui *gui.Panel

	camera *camera.Camera
	persp  *camera.Perspective

	position0 math32.Vector3
	target0   math32.Vector3

	rotateEnd   math32.Vector2
	rotateStart math32.Vector2

	AzimuthMax   float32
	AzimuthMin   float32
	DistanceMax  float32
	DistanceMin  float32
	PolarMax     float32
	PolarMin     float32
	RotateSpeedX float32
	RotateSpeedY float32
	Xoffset      float32
	Yoffset      float32
	ZoomSpeed    float32

	mode CamMode

	Zoom       int8
	ZoomStep1P int8 // negate to invert zoom direction
	ZoomStep3P int8 // negate to invert zoom direction

	EnableKeys bool
	EnableZoom bool

	enabled    bool
	rotating   bool
	subsEvents bool
}

func New(followee Followee, persp *camera.Perspective, iWindow window.IWindow) (c *Control) {
	c = new(Control)
	c.Init(followee, persp, iWindow)
	return
}

func (c *Control) DefaultToScreen() bool {
	return c.Mode().All(DefaultToScreen)
}

func (c *Control) Dispose() {
	c.IWindow.UnsubscribeID(window.OnCursor, &c.subsEvents)
	c.IWindow.UnsubscribeID(window.OnMouseUp, &c.subsEvents)
	c.IWindow.UnsubscribeID(window.OnMouseDown, &c.subsEvents)
	c.IWindow.UnsubscribeID(window.OnScroll, &c.subsEvents)
	c.IWindow.UnsubscribeID(window.OnKeyUp, &c.subsEvents)
	c.IWindow.UnsubscribeID(window.OnKeyDown, &c.subsEvents)
}

func (c *Control) Enabled() bool {
	return c.enabled
}

func (c *Control) Init(followee Followee, persp *camera.Perspective, iWindow window.IWindow) {
	c.Followee = followee

	c.IWindow = iWindow
	width, height := c.IWindow.Size()
	w64, h64 := float64(width), float64(height)
	x, y := w64*0.5, h64*0.5
	c.IWindow.SetCursorPos(x, y)

	c.Gui = nil

	c.camera = persp.GetCamera()
	c.persp = persp
	c.initPositionAndTarget3P()

	c.position0 = c.camera.Position()
	c.target0 = c.camera.Target()

	c.rotateEnd = math32.Vector2{float32(x), float32(y)}
	c.rotateStart = math32.Vector2{float32(x), float32(y)}

	c.AzimuthMax = float32(math.Inf(+1))
	c.AzimuthMin = float32(math.Inf(-1))
	c.DistanceMax = float32(math.Inf(+1))
	c.DistanceMin = 0.01
	c.PolarMax = float32(math.Pi)
	c.PolarMin = 0.0
	c.RotateSpeedX = 0.25
	c.RotateSpeedY = 0.25
	c.ZoomSpeed = 0.1

	c.mode.Init(DefaultToScreen)
	c.updateRotate(0, 0)

	c.Zoom = -0x21
	c.ZoomStep1P = 0x04
	c.ZoomStep3P = 0x04
	c.updateZoomAbsolute()

	c.EnableKeys = true
	c.EnableZoom = true

	c.enabled = true
	c.rotating = false
	c.subsEvents = false

	c.IWindow.SubscribeID(window.OnCursor, &c.subsEvents, c.onMouseCursor)
	c.IWindow.SubscribeID(window.OnMouseUp, &c.subsEvents, c.onMouseButton)
	c.IWindow.SubscribeID(window.OnMouseDown, &c.subsEvents, c.onMouseButton)
	c.IWindow.SubscribeID(window.OnScroll, &c.subsEvents, c.onMouseScroll)
	c.IWindow.SubscribeID(window.OnKeyUp, &c.subsEvents, c.onKeyboardKey)
	c.IWindow.SubscribeID(window.OnKeyDown, &c.subsEvents, c.onKeyboardKey)
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

func (c *Control) SetDefaultToScreen(defaultToScreen bool) (was bool) {
	wasMode := c.Mode()
	if defaultToScreen {
		c.SetMode(CamMode{wasMode.SetCopy(DefaultToScreen)})
	} else {
		c.SetMode(CamMode{wasMode.ClrCopy(DefaultToScreen)})
	}
	was = wasMode.All(DefaultToScreen)
	return
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
		c.IWindow.SetInputMode(window.CursorMode, window.CursorNormal)
		if c.Gui != nil {
			c.Gui.SetEnabled(true)
		}
	case was.Screen() && cm.World():
		if c.Gui != nil {
			c.Gui.SetEnabled(false)
		}
		c.IWindow.SetInputMode(window.CursorMode, window.CursorDisabled)
		w, h := c.IWindow.Size()
		x, y := 0.5*float64(w), 0.5*float64(h)
		c.IWindow.SetCursorPos(x, y)
		c.rotating = true
	}
	return
}

func (c *Control) ZoomBySteps(step1P, step3P int) {
	old := int(c.Zoom)
	var new, init int
	if old >= 0 {
		new = old + step1P
		switch {
		case new < 0:
			new = -1
			init = 3
		case new > 0x70:
			new = 0x70
		}
	} else {
		new = old + step3P
		switch {
		case new >= 0:
			new = 0
			init = 1
		case new < -0x71:
			new = -0x71
		}
	}
	zoom := int8(new)
	c.Zoom = zoom
	switch init {
	case 0:
		c.updateZoomAbsolute()
	case 1:
		c.initPositionAndTarget1P()
	case 3:
		c.initPositionAndTarget3P()
	}
}

func (c *Control) ZoomIn(amount float64) {
	step1P := int(amount * float64(c.ZoomStep1P))
	step3P := int(amount * float64(c.ZoomStep3P))
	c.ZoomBySteps(step1P, step3P)
}

func (c *Control) ZoomOut(amount float64) {
	c.ZoomIn(-amount)
}

func (c *Control) initPositionAndTarget1P() {
	vec := c.Followee.Position()
	x, y, z := float64(vec.X), float64(vec.Y), float64(vec.Z)
	dx, dy := c.Followee.FacingNormalized()
	y -= dy * c.Followee.FrontOfEye()
	z += c.Followee.HeightToEye()
	vec.Z = float32(z)
	c.camera.SetPositionVec(&vec)
	vec.X, vec.Y = float32(x+dx), float32(y+dy)
	c.camera.LookAt(&vec)
	c.updateZoomAbsolute()
}

func (c *Control) initPositionAndTarget3P() {
	target := c.Followee.Position()
	x, y, z := float64(target.X), float64(target.Y), float64(target.Z)
	z += 0.70260 * c.Followee.HeightToEye()
	target.Z = float32(z)
	var vec math32.Vector3
	dx, dy := c.Followee.FacingNormalized()
	vec.X, vec.Y = float32(x-dx), float32(y-dy)
	vec.Z = float32(z * math.Phi)
	c.camera.SetPositionVec(&vec)
	c.camera.LookAt(&target)
	c.persp.SetFov(65)
	c.updateZoomAbsolute()
}

func (c *Control) onKeyboardKey(evname string, event interface{}) {
	if !c.Enabled() || !c.EnableKeys {
		return
	}
	ev := event.(*window.KeyEvent)
	switch ev.Keycode {
	case window.KeyLeftAlt, window.KeyLeftSuper:
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
			// c.Snap() // TODO:
		case window.Release:
		}
	}
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
		c.SetMode(CamMode{c.Mode().SetCopy(MiddleMouseHeld)})
	case window.Release:
		c.SetMode(CamMode{c.Mode().ClrCopy(MiddleMouseHeld)})
	}
}

func (c *Control) onMouseCursor(evname string, event interface{}) {
	if !c.mutexMouseCursor.TryLock() {
		return
	}
	defer c.mutexMouseCursor.Unlock()

	ev := event.(*window.CursorEvent)
	xOffset, yOffset := ev.Xpos, ev.Ypos
	c.Xoffset, c.Yoffset = xOffset, yOffset

	if !c.rotating || !c.Enabled() || c.Mode().Screen() {
		return
	}

	c.rotateEnd.Set(xOffset, yOffset)
	var rotateDelta math32.Vector2 // TODO: don't use vectors for this
	rotateDelta.SubVectors(&c.rotateEnd, &c.rotateStart)
	c.rotateStart = c.rotateEnd
	width, height := c.IWindow.Size()
	w64, h64 := float64(width), float64(height)
	by := 2.0 * math.Pi
	c.RotateLeft(by * float64(c.RotateSpeedX) / float64(w64) * float64(rotateDelta.X))
	c.RotateUp(by * float64(c.RotateSpeedY) / float64(h64) * float64(rotateDelta.Y))
}

func (c *Control) onMouseScroll(evname string, event interface{}) {
	if !c.Enabled() || !c.EnableZoom || c.Mode().Screen() {
		return
	}
	ev := event.(*window.ScrollEvent)
	c.ZoomIn(float64(ev.Yoffset))
}

const updateRotateEpsilon float64 = 0.01
const updateRotatePiMinusEpsilon float64 = math.Pi - updateRotateEpsilon

func (c *Control) updateRotate(thetaDelta, phiDelta float64) {
	var max, min float64
	if float64(c.PolarMax) < updateRotatePiMinusEpsilon {
		max = float64(c.PolarMax)
	} else {
		max = updateRotatePiMinusEpsilon
	}
	if float64(c.PolarMin) > updateRotateEpsilon {
		min = float64(c.PolarMin)
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
	theta = maths.ClampFloat64(theta, float64(c.AzimuthMin), float64(c.AzimuthMax))
	phi = maths.ClampFloat64(phi, float64(min), float64(max))
	vdir.X = float32(radius * math.Sin(phi) * math.Sin(theta))
	vdir.Y = float32(radius * math.Cos(phi))
	vdir.Z = float32(radius * math.Sin(phi) * math.Cos(theta))
	vdir.ApplyQuaternion(&quatInverse)
	if c.Zoom < 0 {
		position = target
		position.Add(&vdir)
		c.camera.SetPositionVec(&position)
		c.camera.LookAt(&target)
	} else {
		target = position
		target.Sub(&vdir)
		c.camera.SetPositionVec(&position)
		c.camera.LookAt(&target)
	}
}

const updateZoomEpsilon float64 = 0.01
const updateZoomEpsilonNegated float64 = -updateZoomEpsilon
const updateZoomAbsoluteScalar1P float64 = 1.0 / 16.0
const updateZoomAbsoluteScalar3P float64 = 1.0 / 16.0

func (c *Control) updateZoomAbsolute() {
	zoom := c.Zoom
	if zoom < 0 {
		// Lock the target and change the position
		power := float64(-(zoom + 1)) * updateZoomAbsoluteScalar3P
		distance := math.Pow(math.Phi, power)
		position := c.camera.Position()
		target := c.camera.Target()
		position.Sub(&target)
		distance = maths.ClampFloat64(distance,
			float64(c.DistanceMin), float64(c.DistanceMax))
		position.SetLength(float32(distance))
		target.Add(&position)
		c.camera.SetPositionVec(&target)
	} else {
		// Just change the field of view
		power := float64(zoom) * updateZoomAbsoluteScalar3P
		scalar := 1.0 / math.Pow(math.Phi, power)
		c.persp.SetFov(float32(65.0 * scalar))
	}
}
