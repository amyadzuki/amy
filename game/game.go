package game

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/amyadzuki/amystuff/logs"
	"github.com/amyadzuki/amystuff/str"
	"github.com/amyadzuki/amystuff/styles"

	//	"github.com/g3n/engine/audio"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	//	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

type Game struct {
	Camera *camera.Perspective // camera.ICamera if my PR gets accepted
	Win    window.IWindow
	Wm     window.IWindowManager

	DockBotLeft  *gui.Panel
	DockBotRight *gui.Panel
	DockTopLeft  *gui.Panel
	DockTopRight *gui.Panel

	LabelCharaChangerBlank string
	LabelFullScreen        string
	LabelWindow            string
	Title                  string

	WidgetCharaChanger *gui.Button
	WidgetClose        *gui.Button
	WidgetFullScreen   *gui.Button
	WidgetHelp         *gui.Button
	WidgetIconify      *gui.Button
	WidgetFps          [2]*gui.Label
	WidgetPing         [2]*gui.Label

	MaxWidthFps  float64
	MaxWidthPing float64

	Frame    int64
	SecFrame int64 // frame at start of this second

	Gs    *gls.GLS
	Logs  *logs.Logs
	Rend  *renderer.Renderer
	Root  *gui.Root
	Scene *core.Node
	w, h  int

	AskQuit   int8
	WantHelp  bool
	HaveAudio bool
	InfoDebug bool
	InfoTrace bool
	MusicHush bool
	MusicMute bool
}

func New(title string) (game *Game) {
	game = new(Game)
	game.Title = title
	game.Scene = core.NewNode()
	return
}

func (game *Game) AddDockTopRight() {
	game.DockTopRight = gui.NewPanel(0, 0)
	game.DockTopRight.SetLayout(gui.NewDockLayout())
	game.Root.Add(game.DockTopRight)
}

func (game *Game) AddWidgetCharaChanger(label string) {
	if game.DockTopLeft == nil {
		game.AddDockTopLeft()
	}
	game.LabelCharaChangerBlank = label
	game.WidgetCharaChanger = gui.NewButton(label)
	game.WidgetCharaChanger.SetLayoutParams(&gui.DockLayoutParams{gui.DockLeft})
	game.WidgetCharaChanger.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		//
	})
	game.DockTopLeft.Add(game.WidgetCharaChanger)
}

func (game *Game) AddWidgetClose(label string) {
	if game.DockTopRight == nil {
		game.AddDockTopRight()
	}
	game.WidgetClose = gui.NewButton(label)
	game.WidgetClose.SetLayoutParams(&gui.DockLayoutParams{gui.DockRight})
	game.WidgetClose.SetStyles(&styles.CloseButton)
	game.WidgetClose.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		game.WidgetClose.SetStyles(&styles.ClosingButton)
		if game.SoftQuit() > 0 {
			game.Quit()
		}
		go func() {
			time.Sleep(time.Second)
			game.AskQuit--
			if game.AskQuit < 0 {
				game.AskQuit = 0
			}
			game.WidgetClose.SetStyles(&styles.CloseButton)
		}()
	})
	game.dockSize(game.DockTopRight, game.WidgetClose)
	game.DockTopRight.Add(game.WidgetClose)
}

func (game *Game) AddWidgetFps() {
	game.WidgetFps[0] = gui.NewLabel("9999")
	game.MaxWidthFps = float64(game.WidgetFps[0].TotalWidth())
	game.WidgetFps[0].SetText("")
	game.WidgetFps[1] = gui.NewLabel(" fps")
	game.Root.Add(game.WidgetFps[0])
	game.Root.Add(game.WidgetFps[1])
}

func (game *Game) AddWidgetFullScreen(labelFullScreen, labelWindow string) {
	if game.DockTopRight == nil {
		game.AddDockTopRight()
	}
	game.LabelFullScreen = labelFullScreen
	game.LabelWindow = labelWindow
	label := labelFullScreen
	if game.FullScreen() {
		label = labelWindow
	}
	game.WidgetFullScreen = gui.NewButton(label)
	game.WidgetFullScreen.SetLayoutParams(&gui.DockLayoutParams{gui.DockRight})
	game.WidgetFullScreen.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		game.ToggleFullScreen()
	})
	game.dockSize(game.DockTopRight, game.WidgetFullScreen)
	game.DockTopRight.Add(game.WidgetFullScreen)
}

func (game *Game) AddWidgetHelp(label string) {
	if game.DockTopRight == nil {
		game.AddDockTopRight()
	}
	game.WidgetHelp = gui.NewButton(label)
	game.WidgetHelp.SetLayoutParams(&gui.DockLayoutParams{gui.DockRight})
	game.WidgetHelp.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		game.WantHelp = !game.WantHelp
	})
	game.dockSize(game.DockTopRight, game.WidgetHelp)
	game.DockTopRight.Add(game.WidgetHelp)
}

func (game *Game) AddWidgetIconify(label string) {
	if game.DockTopRight == nil {
		game.AddDockTopRight()
	}
	game.WidgetIconify = gui.NewButton(label)
	game.WidgetIconify.SetLayoutParams(&gui.DockLayoutParams{gui.DockRight})
	game.WidgetIconify.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		// TODO
	})
	game.dockSize(game.DockTopRight, game.WidgetIconify)
	game.DockTopRight.Add(game.WidgetIconify)
}

func (game *Game) AddWidgetPing() {
	game.WidgetPing[0] = gui.NewLabel("9999")
	game.MaxWidthPing = float64(game.WidgetPing[0].TotalWidth())
	game.WidgetPing[0].SetText("")
	game.WidgetPing[1] = gui.NewLabel(" ms")
	game.Root.Add(game.WidgetPing[0])
	game.Root.Add(game.WidgetPing[1])
}

func (game *Game) FullScreen() bool {
	return game.Win.FullScreen()
}

func (game *Game) Quit() {
	game.Win.SetShouldClose(true)
}

func (game *Game) RecalcDocks() {
	w, h := game.Size()
	w64, h64 := float64(w), float64(h)
	if game.DockTopLeft != nil {
		game.DockTopLeft.SetPosition(0, 0)
	}
	if game.DockTopRight != nil {
		x := float32(w64 - float64(game.DockTopRight.TotalWidth()))
		game.DockTopRight.SetPosition(x, 0)
	}
	if game.DockBotLeft != nil {
		y := float32(h64 - float64(game.DockBotLeft.TotalHeight()))
		game.DockBotLeft.SetPosition(0, y)
	}
	if game.DockBotRight != nil {
		x := float32(w64 - float64(game.DockBotRight.TotalWidth()))
		y := float32(h64 - float64(game.DockBotRight.TotalHeight()))
		game.DockBotRight.SetPosition(x, y)
	}
}

func (game *Game) RecalcPerformanceWidgets() {
	w, _ := game.Size()
	x := float64(w)
	if game.DockTopRight != nil {
		x -= float64(game.DockTopRight.TotalWidth())
	}
	if game.WidgetPing[0] != nil && game.WidgetPing[1] != nil {
		a := game.MaxWidthPing
		b := float64(game.WidgetPing[1].TotalWidth())
		x -= b
		game.WidgetPing[1].SetPosition(float32(x), 0)
		x -= a
		game.WidgetPing[0].SetPosition(float32(x), 0)
	}
	if game.WidgetFps[0] != nil && game.WidgetFps[1] != nil {
		a := game.MaxWidthFps
		b := float64(game.WidgetFps[1].TotalWidth())
		x -= b
		game.WidgetFps[1].SetPosition(float32(x), 0)
		x -= a
		game.WidgetFps[0].SetPosition(float32(x), 0)
	}
}

func (game *Game) SetFullScreen(fullScreen bool) {
	game.Win.SetFullScreen(fullScreen)
	if game.WidgetFullScreen != nil {
		label := game.LabelFullScreen
		if fullScreen {
			label = game.LabelWindow
		}
		game.WidgetFullScreen.Label.SetText(label)
	}
	game.onWinCh("", nil)
}

func (game *Game) Size() (w, h int) {
	w, h = game.w, game.h
	return
}

func (game *Game) SizeRecalc() (w, h int) {
	w, h = game.Win.Size()
	game.w, game.h = w, h
	return
}

func (game *Game) SoftQuit() int8 {
	was := game.AskQuit
	now := int(was) + 1
	if now > 127 {
		now = 127
	}
	game.AskQuit = int8(now)
	return was
}

func (game *Game) StartUp(logPath string) (err error) {
	flag_debug := CommandLine.Bool("debug", false,
		"Log debug info (may slightly slow the game)")
	flag_trace := CommandLine.Bool("debugextra", false,
		"Log trace info (may drastically slow the game)")
	flag_quiet := CommandLine.Bool("quiet", false,
		"Silence -info- messages from the console")
	flag_fullscreen := CommandLine.Bool("fullscreen", false,
		"Launch fullscreen")
	flag_geometry := CommandLine.String("geometry", "960x720",
		"Window geometry (H, WxH, or WxH+X+Y)")
	flag_wm := CommandLine.String("wm", "glfw",
		"Window manager (one of: \"glfw\")")

	func() {
		defer func() {
			if r := recover(); r != nil {
				CommandLine.Bool("x", false, "This help text")
				switch r {
				case flag.ErrHelp:
					os.Exit(0)
				default:
					os.Exit(2)
				}
			}
		}()
		CommandLine.Parse(os.Args[1:])
	}()

	info, debug, trace := !*flag_quiet, *flag_debug, *flag_trace
	game.InfoDebug = debug || trace
	game.InfoTrace = trace
	if game.Logs, err = logs.New(logPath, info, debug, trace); err != nil {
		return
	}

	game.Major("Created process and initialized logging for \"" + game.Title + "\"")

	w, h, x, y, n := 960, 720, 0, 0, 0
	n, err = fmt.Sscanf(*flag_geometry, "%dx%d+%d+%d", &w, &h, &x, &y)
	if n < 1 || n > 4 || (n == 4 && err != nil) {
		game.Warn("could not parse window geometry \"" + *flag_geometry + "\"")
	}
	if n == 1 {
		h = w
		w = h * 4 / 3
	}

	if h < 120 {
		game.Warn("height was capped at the minimum height of 120")
		h = 120
	}

	if w < 160 {
		game.Warn("width was capped at the minimum width of 160")
		w = 160
	}

	fs := *flag_fullscreen

	wm := str.Simp(*flag_wm)
	switch wm {
	case "glfw":
	default:
		game.Warn("unsupported window manager \"" + wm + "\" changed to \"glfw\"")
		wm = "glfw"
	}

	if game.Wm, err = window.Manager(wm); err != nil {
		game.Error("window.Manager")
		return
	}

	startupMessage :=
		"Launching \"" + game.Title + "\" " + strconv.Itoa(w) + "x" +
			strconv.Itoa(h) + " at (" + strconv.Itoa(x) + ", " +
			strconv.Itoa(y) + ") "
	if fs {
		startupMessage += "fullscreen"
	} else {
		startupMessage += "windowed"
	}
	startupMessage += " using \""
	startupMessage += wm
	startupMessage += "\""
	game.Info(startupMessage)

	if game.Win, err = game.Wm.CreateWindow(w, h, game.Title, fs); err != nil {
		game.Error("game.Wm.CreateWindow")
		return
	}

	// OpenGL functions must be executed in the same thread
	// where the context was created (by `CreateWindow`)
	runtime.LockOSThread()

	// Create the OpenGL state
	if game.Gs, err = gls.New(); err != nil {
		game.Error("gls.New")
		return
	}

	width, height := game.SizeRecalc()
	game.ViewportFull()
	aspect := float32(float64(width) / float64(height))
	game.Camera = camera.NewPerspective(65, aspect, 1.0/128.0, 1024.0)

	game.Root = gui.NewRoot(game.Gs, game.Win)
	game.Root.SetSize(float32(width), float32(height))

	game.Win.Subscribe(window.OnWindowSize, game.onWinCh)
	game.Win.Subscribe(window.OnKeyDown, game.onKeyboardKey)
	game.Win.Subscribe(window.OnKeyUp, game.onKeyboardKey)
	game.Win.Subscribe(window.OnMouseDown, game.onMouseButton)
	game.Win.Subscribe(window.OnMouseUp, game.onMouseButton)
	game.Win.Subscribe(window.OnCursor, game.onMouseCursor)

	return
}

func (game *Game) ToggleFullScreen() {
	game.SetFullScreen(!game.FullScreen())
}

func (game *Game) ViewportFull() {
	w, h := game.Size()
	game.Gs.Viewport(0, 0, int32(w), int32(h))
	return
}

/*
func (game *Game) VolumeChanged() {
	if game.HaveAudio {
		loud := game.Settings.MusVolume.Value()
		if game.MusicHush {
			quiet := float32(float64(loud) * 0.5)
			game.PlayerMusic.SetGain(quiet)
		} else {
			game.PlayerMusic.SetGain(loud)
		}
	}
}
*/

// Logging functions, in order of ascending importance

func (game *Game) Minor(v ...interface{}) {
	if game.Logs != nil {
		game.Logs.Minor.Println(v...)
	}
}

func (game *Game) Major(v ...interface{}) {
	if game.Logs != nil {
		game.Logs.Major.Println(v...)
	}
}

func (game *Game) Debug(v ...interface{}) {
	if game.Logs != nil {
		game.Logs.Debug.Println(v...)
	}
}

func (game *Game) Info(v ...interface{}) {
	if game.Logs != nil {
		game.Logs.Info.Println(v...)
	}
}

func (game *Game) Warn(v ...interface{}) {
	if game.Logs != nil {
		game.Logs.Warn.Println(v...)
	}
}

func (game *Game) Error(v ...interface{}) {
	if game.Logs != nil {
		game.Logs.Error.Println(v...)
	}
}

func (game *Game) Fatal(v ...interface{}) {
	if game.Logs != nil {
		game.Logs.Fatal.Fatalln(v...)
	}
}

// Internal functions

func (game *Game) dockSize(p *gui.Panel, w gui.IPanel) {
	oldW, oldH := p.TotalWidth(), p.TotalHeight()
	newW, newH := oldW + w.TotalWidth(), w.TotalHeight()
	if oldH > newH {
		newH = oldH
	}
	p.SetWidth(newW)
	p.SetHeight(newH)
}

func (game *Game) onKeyboardKey(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)
	switch kev.Keycode {
	case window.KeyF11:
		game.ToggleFullScreen()
	}
}

func (game *Game) onMouseButton(evname string, ev interface{}) {
}

func (game *Game) onMouseCursor(evname string, ev interface{}) {
}

func (game *Game) onWinCh(evname string, ev interface{}) {
	if game.Win == nil {
		game.Warn("onWinCh but game.Win was nil")
		return
	}
	w, h := game.SizeRecalc()
	if game.Gs != nil {
		game.Gs.Viewport(0, 0, int32(w), int32(h))
	} else {
		game.Warn("onWinCh but game.GS was nil")
	}
	if game.Root != nil {
		game.Root.SetSize(float32(w), float32(h))
	} else {
		game.Warn("onWinCh but game.Root was nil")
	}
	game.RecalcDocks()
	game.RecalcPerformanceWidgets()
	if game.Camera != nil {
		game.Camera.SetAspect(float32(float64(w) / float64(h)))
	} else {
		game.Warn("onWinCh but game.Camera was nil")
	}
}

func printDefault(f *flag.Flag) {
	s := "  " + FlgHelpBeforeFlag + "-" + f.Name + FlgHelpAfterFlag
	name, usage := flag.UnquoteUsage(f)
	if len(name) > 0 {
		s += " " + FlgHelpBeforeType
		s += name // type name
		s += FlgHelpAfterType + " "
	}
	if len(s) <= 4 { // space, space, hyphen, character
		s += "\t"
	} else {
		s += ConsoleNewLineAndIndent
	}
	s += strings.Replace(usage, "\n", ConsoleNewLineAndIndent, -1)
	if !flag_isZeroValue(f, f.DefValue) {
		if reflect.TypeOf(f.Value).Elem().Kind() == reflect.String {
			s += fmt.Sprintf(" (default: %q)", f.DefValue)
		} else {
			s += fmt.Sprintf(" (default: %v)", f.DefValue)
		}
	}
	fmt.Fprint(CommandLine.Output(), s, "\n")
}

func flag_isZeroValue(flg *flag.Flag, value string) bool {
	typ := reflect.TypeOf(flg.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	if value == z.Interface().(flag.Value).String() {
		return true
	}
	switch value {
	case "false", "", "0":
		return true
	}
	return false
}

var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	CommandLine.VisitAll(printDefault)
	fmt.Fprintln(CommandLine.Output(),
		"  "+FlgHelpBeforeFlag+"-h"+FlgHelpAfterFlag+" | "+
			FlgHelpBeforeFlag+"-help"+FlgHelpAfterFlag+" | "+
			FlgHelpBeforeFlag+"-?"+FlgHelpAfterFlag+
			ConsoleNewLineAndIndent+"Show this help message\n  "+
			"Note that you can use either one or two hyphens wherever one is shown")
}

var CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)

func init() {
	CommandLine.Usage = Usage
}

const (
	ConsoleNewLineAndIndent = "\n      \t"
	VT100Bold               = "\x1b[1m"
	VT100Italic             = "\x1b[3m"
	VT100Underline          = "\x1b[4m"
	VT100Strike             = "\x1b[9m"
	VT100Reset              = "\x1b[0m\x1b[m"
	FlgHelpBeforeFlag       = VT100Bold
	FlgHelpAfterFlag        = VT100Reset
	FlgHelpBeforeType       = VT100Underline
	FlgHelpAfterType        = VT100Reset
)
