// Copyleft (L) 5778 Asaf Ohayon.

package main

import (
	"flag"
	"runtime"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/logger"
	"github.com/g3n/engine/window"
)

var log *logger.Logger

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

func main() {
	// Parse command line flags
	showLog := flag.Bool("debug", false, "display the debug log")
	flag.Parse()

	// Create logger
	log = logger.New("Cube", nil)

	log.AddWriter(logger.NewConsole(false))
	log.SetFormat(logger.FTIME | logger.FMICROS)

	if *showLog == true {
		log.SetLevel(logger.DEBUG)
	} else {
		log.SetLevel(logger.INFO)
	}

	log.Info("Initializing Cube")

	// Get the window manager
	wmgr, err := window.Manager("glfw")
	if err != nil {
		panic(err)
	}

	win, err := wmgr.CreateWindow(1200, 900, "Cube", false)
	if err != nil {
		panic(err)
	}

	// Create OpenGL state
	gs, err := gls.New()
	if err != nil {
		panic(err)
	}

	// Speed up a bit by not checking OpenGL errors
	gs.SetCheckErrors(false)

	// Sets window background color
	gs.ClearColor(0.1, 0.1, 0.1, 1.0)

	// Sets the OpenGL viewport size the same as the window size
	// This normally should be updated if the window is resized.
	width, height := win.Size()
	gs.Viewport(0, 0, int32(width), int32(height))

	// Creates GUI root panel
	root := gui.NewRoot(gs, win)
	root.SetSize(float32(width), float32(height))

	// Creates a renderer and adds default shaders
	renderer := renderer.NewRenderer(gs)

	err = renderer.AddDefaultShaders()
	if err != nil {
		panic(err)
	}
	renderer.SetGui(root)
}
