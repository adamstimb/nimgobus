// This is an RM Basic demo program called ROPE.BAS rewritten 40 years later in Go.
package main

import (
	_ "image/png" // import only for side-effects
	"log"
	"time"

	nimgobus "github.com/adamstimb/nimgobus/pkg"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	launch          int
	nimgobus.Nimbus // Embed the Nimbus in the Game struct
}

func NewGame() *Game {
	game := &Game{}
	game.Init() // Initialize the Nimbus
	return game
}

func (g *Game) Update() error {
	if g.launch == 0 {
		go App(g) // Launch the Nimbus app on first iteration
	}
	g.launch = 1
	g.Nimbus.Update() // Update the app on all subsequent iterations
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// The main loop
func mainLoop(g *Game) {
	g.Subbios.TGraphicsOutput.FGraphicsOutputColdStart()
	g.Subbios.Stdio.Printf("\x1b0h") // Mode 40
	xOld := 0
	yOld := 0
	paint := 1
	g.Subbios.TGraphicsOutput.FPolymarker(2, 1, 1, 259, [][2]int{}, [][2]int{{xOld, yOld}})
	for {
		xPos, yPos, button := g.Subbios.TGraphicsInput.FEnquirePositionAndButtonStatus()
		g.Subbios.TGraphicsOutput.FPolymarker(2, 1, 1, 259, [][2]int{}, [][2]int{{xOld, yOld}})
		switch button {
		case 1:
			g.Subbios.TGraphicsOutput.FPieSlice(xPos, yPos, 10, 0, 0, paint)
		case 2:
			changeCol(g)
		case 3:
			g.Subbios.Stdio.Printf("\x1b0h")
		}
		g.Subbios.TGraphicsOutput.FPolymarker(2, 1, 1, 259, [][2]int{}, [][2]int{{xPos, yPos}})
		xOld = xPos
		yOld = yPos
		paint = (paint + 1) & 15
	}
}

func changeCol(g *Game) {
	var colr int
	for {
		for ink := 1; ink <= 15; ink++ {
			colr = (colr + 1) & 15
			g.Subbios.TGraphicsOutput.FSetCltElement(ink, colr, 0, colr)
			time.Sleep(20 * time.Millisecond) // Because computers are slightly faster than 40 years ago...
		}
		_, _, button := g.Subbios.TGraphicsInput.FEnquirePositionAndButtonStatus()
		if button > 0 {
			return
		}
	}
}

func App(g *Game) {
	// This is the Nimbus app itself
	mainLoop(g)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the Nimbus monitor on the screen and scale to current window size.
	monitorWidth, monitorHeight := g.Monitor.Size()

	// Get ebiten window size so we can scale the Nimbus screen up or down
	// but if (0, 0) is returned we're not running on a desktop so don't do any scaling
	windowWidth, windowHeight := ebiten.WindowSize()

	// Calculate aspect ratios of Nimbus monitor and ebiten screen
	monitorRatio := float64(monitorWidth) / float64(monitorHeight)
	windowRatio := float64(windowWidth) / float64(windowHeight)

	// If windowRatio > monitorRatio then clamp monitorHeight to windowHeight otherwise
	// clamp monitorWidth to screenWidth
	var scale, offsetX, offsetY float64
	switch {
	case windowRatio > monitorRatio:
		scale = float64(windowHeight) / float64(monitorHeight)
		offsetX = (float64(windowWidth) - float64(monitorWidth)*scale) / 2
		offsetY = 0
	case windowRatio <= monitorRatio:
		scale = float64(windowWidth) / float64(monitorWidth)
		offsetX = 0
		offsetY = (float64(windowHeight) - float64(monitorHeight)*scale) / 2
	}

	// Apply scale and centre monitor on screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(offsetX, offsetY)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(g.Monitor, op)
}

func main() {
	// Set up resizeable window
	ebiten.SetWindowSize(1000, 800)
	ebiten.SetWindowTitle("Nimgobus")
	ebiten.SetWindowResizable(true)

	// Create a new game and pass it to RunGame method
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
