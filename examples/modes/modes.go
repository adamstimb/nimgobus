package main

import (
	_ "image/png" // import only for side-effects
	"log"
	"time"

	"github.com/adamstimb/nimgobus"
	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 1400
	screenHeight = 1000
)

var (
	nim nimgobus.Nimbus
)

type Game struct {
	count int
}

func (g *Game) Update() error {
	if g.count == 0 {
		go Test()
	}
	g.count++
	nim.Update()
	return nil
}

func Test() {
	Mode80()
	time.Sleep(4 * time.Second)
	Mode40()
}

func Mode80() {
	nim.SetMode(80)
	plotOpts := nimgobus.PlotOptions{
		Brush: 2,
		SizeX: 2,
		SizeY: 2,
	}
	nim.Plot(plotOpts, "Mode 80", 30, 220)
	plotOpts.Brush = 3
	nim.Plot(plotOpts, "Mode 80", 32, 221)
	nim.SetWriting(1, 5, 7, 75, 11)
	nim.SetWriting(1)
	nim.Print("This is Mode 80.  The screen is 80 character columns wide and 25 columns tall, or 640 pixels wide and 250 pixels tall.  Pixels are doubled in length along the vertical so everything has this wacky stretched-out look!  4 colours are available.")
	// Draw colour swatch
	areaOpts := nimgobus.AreaOptions{}
	var x, y int
	width := 143
	for i := 0; i < 4; i++ {
		areaOpts.Brush = i
		y = 50
		x = 30 + (i * width)
		areaOpts.Brush = 3
		nim.Area(areaOpts, x-1, y-1, x+width+1, y-1, x+width+1, y+81, x-1, y+81, x-1, y-1)
		areaOpts.Brush = i
		nim.Area(areaOpts, x, y, x+width, y, x+width, y+80, x, y+80, x, y)
	}
}

func Mode40() {
	nim.SetMode(40)
	nim.SetPaper(1)
	nim.SetBorder(1)
	nim.Cls()
	plotOpts := nimgobus.PlotOptions{
		Brush: 0,
		SizeX: 2,
		SizeY: 2,
	}
	nim.Plot(plotOpts, "Mode 40", 15, 220)
	plotOpts.Brush = 14
	nim.Plot(plotOpts, "Mode 40", 16, 221)
	nim.SetWriting(1, 3, 6, 38, 13)
	nim.SetWriting(1)
	nim.Print("This is Mode 40.  The screen is 40 character columns wide and 25 columns tall, or 320 pixels wide and 250 pixels tall.  16 sumptious colours are available.")
	// Draw colour swatch
	areaOpts := nimgobus.AreaOptions{}
	var x, y int
	width := 18
	for i := 0; i < 16; i++ {
		areaOpts.Brush = i
		y = 50
		x = 15 + (i * width)
		areaOpts.Brush = 15
		nim.Area(areaOpts, x-1, y-1, x+width+1, y-1, x+width+1, y+81, x-1, y+81, x-1, y-1)
		areaOpts.Brush = i
		nim.Area(areaOpts, x, y, x+width, y, x+width, y+80, x, y+80, x, y)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Draw the Nimbus monitor on the screen and scale to current window size.
	monitorWidth, monitorHeight := nim.Monitor.Size()
	// Calculate aspect ratios of Nimbus monitor and ebiten screen
	monitorRatio := float64(monitorWidth) / float64(monitorHeight)
	screenRatio := float64(screenWidth) / float64(screenHeight)

	// If screenRatio > monitorRatio then clamp monitorHeight to screenHeight otherwise
	// clamp monitorWidth to screenWidth
	var scale, offsetX, offsetY float64
	if screenRatio > monitorRatio {
		scale = float64(screenHeight) / float64(monitorHeight)
		offsetX = (float64(screenWidth) - float64(monitorWidth)*scale) / 2
		offsetY = 0
	}
	if screenRatio <= monitorRatio {
		scale = float64(screenWidth) / float64(monitorWidth)
		offsetX = 0
		offsetY = (float64(screenHeight) - float64(monitorHeight)*scale) / 2
	}

	// Apply scale and centre monitor on screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(offsetX, offsetY)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(nim.Monitor, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {

	// Initialize the Nimbus
	nim.Init()

	// set up window
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Nimgobus Test")

	// Call RunGame method, passing the address of the pointer to an empty Game struct
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
