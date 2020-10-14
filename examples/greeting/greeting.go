package main

import (
	_ "image/png" // import only for side-effects
	"log"

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
	nim.Boot()       // Boot the Nimbus! (this is optional)
	nim.SetMode(40)  // Low-res, high-colour mode
	nim.SetBorder(9) // Light blue border
	nim.SetPaper(1)  // Dark blue paper
	nim.Cls()        // Clear screen
	// Plot some text with a shadow effect
	op := nimgobus.PlotOptions{
		SizeX: 3, SizeY: 6, Brush: 0,
	}
	nim.Plot(op, "Nimgobus", 65, 150)
	op.Brush = 14
	nim.Plot(op, "Nimgobus", 67, 152)
	op.SizeX = 1
	op.SizeY = 1
	op.Brush = 13
	nim.Plot(op, "it ain't no real thing", 70, 70)
	nim.PlonkLogo(8, 8) // Draw the Nimbus BIOS logo
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
