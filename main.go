// Copyleft (L) 5778 Asaf Ohayon.

package main

import (
	"flag"
	"runtime"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/camera/control"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
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

	// Adds a perspective camera to the scene
	// The camera aspect ratio should be updated if the window is resized.
	aspect := float32(width) / float32(height)
	camera := camera.NewPerspective(65, aspect, 0.01, 1000)
	camera.SetPosition(0, 4, 5)
	camera.LookAt(&math32.Vector3{X: 0, Y: 0, Z: 0})

	// Create orbit control and set limits
	orbitControl := control.NewOrbitControl(camera, win)
	orbitControl.Enabled = false
	orbitControl.EnablePan = false
	orbitControl.MaxPolarAngle = 2 * math32.Pi / 3
	orbitControl.MinDistance = 5
	orbitControl.MaxDistance = 15

	// Create main scene and child levelScene
	scene := core.NewNode()
	levelScene := core.NewNode()
	scene.Add(camera)
	scene.Add(levelScene)
	//stepDelta := math32.NewVector2(0, 0)
	renderer.SetScene(scene)

	// Add white ambient light to the scene
	ambLight := light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.4)
	scene.Add(ambLight)

	for {
		RenderFrame(root, renderer, camera, wmgr, win)
	}
}

// RenderFrame renders a frame of the scene with the GUI overlaid
func RenderFrame(root *gui.Root, renderer *renderer.Renderer, camera *camera.Perspective, wmgr window.IWindowManager, win window.IWindow) {

	// Process GUI timers
	root.TimerManager.ProcessTimers()

	// Render the scene/gui using the specified camera
	rendered, err := renderer.Render(camera)
	if err != nil {
		panic(err)
	}

	// Check I/O events
	wmgr.PollEvents()

	// Update window if necessary
	if rendered {
		win.SwapBuffers()
	}
}
