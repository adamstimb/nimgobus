// Demonstrating RM Nimbus graphics with nimgobus
package main

import (
	"fmt"
	_ "image/png" // import only for side-effects
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/adamstimb/nimgobus/internal/subbios"
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

// set colour lookup table so that the physical colours are indexed in the order of logical colours
func sortClt(g *Game) {
	for c := 0; c <= 15; c++ {
		msg := fmt.Sprintf("\x1b40;%d;%d;0;0~C", c, c)
		g.Subbios.Stdio.Printf(msg)
	}
}

// outlinePlot is a plot character string function that outlines the font
func outlinePlot(g *Game, orientation, yMagnification, xMagnification, outlineColour, logicalColour, font int, chars string, x, y int) {
	xo := x
	yo := y
	for i := 0; i <= 7; i++ {
		x := xo
		y := yo
		switch i {
		case 0:
			x--
			y--
		case 1:
			y--
		case 2:
			x++
			y--
		case 3:
			x++
		case 4:
			x++
			y++
		case 5:
			y++
		case 6:
			x--
			y++
		case 7:
			x--
		}
		g.Subbios.TGraphicsOutput.FPlotCharacterString(orientation, yMagnification, xMagnification, outlineColour, font, chars, x, y)
	}
	g.Subbios.TGraphicsOutput.FPlotCharacterString(orientation, yMagnification, xMagnification, logicalColour, font, chars, xo, yo)
}

func introScreen(g *Game) {
	g.Subbios.Stdio.Printf("\x1b0h") // Mode 40
	g.Subbios.TGraphicsOutput.FGraphicsOutputColdStart()
	sortClt(g)
	g.Subbios.Stdio.Printf("\x1b~F")              // Hide cursor
	g.Subbios.TGraphicsOutput.FSetBorderColour(1) // Blue border
	g.Subbios.Stdio.Printf("\x1b45;51m")          // Blue background, white text
	g.Subbios.Stdio.Printf("\x1b2J")              // CLS
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 3, 2, 0, 0, "RM Nimbus Graphics", 20, 165)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 3, 2, 14, 0, "RM Nimbus Graphics", 22, 167)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 0, 0, "with", 145, 145)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 9, 0, "with", 146, 146)
	outlinePlot(g, 0, 4, 5, 0, 13, 0, "nimgobus", 5, 100)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue", 60, 0)
}

func screenModes(g *Game) {
	g.Subbios.Stdio.Printf("\x1b2J") // CLS
	//                      1234567890123456789012345678901234567890
	g.Subbios.Stdio.Printf(" Two screen modes are available: In Mode")
	g.Subbios.Stdio.Printf(" 40 the screen is 40 characters wide by\n")
	g.Subbios.Stdio.Printf(" 25 characters high, or 320x250 pixels.\n\n")
	g.Subbios.Stdio.Printf(" In Mode 80 the screen is 80 characters\n")
	g.Subbios.Stdio.Printf(" wide by 25 high, or 640x250 pixels. The")
	g.Subbios.Stdio.Printf(" screen is also scaled horizontally to\n")
	g.Subbios.Stdio.Printf(" maintain the correct aspect ratio.\n")
	g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 15, 0, 0, []int{20, 40, 150, 40, 150, 150, 20, 150, 20, 40})
	g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 15, 0, 0, []int{170, 40, 300, 40, 300, 150, 170, 150, 170, 40})
	g.Subbios.TGraphicsOutput.FPlotCharacterString(1, 1, 1, 15, 0, "250px", 10, 75)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(1, 1, 1, 15, 0, "250px", 160, 75)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "320px", 70, 30)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "640px", 220, 30)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "(0,0)", 25, 45)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "(0,0)", 175, 45)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "(319,249)", 75, 135)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "(0,0)", 175, 45)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "(639,249)", 225, 135)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Mode 40", 60, 155)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Mode 80", 210, 155)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue...", 0, 0)
}

func colours(g *Game) {
	g.Subbios.Stdio.Printf("\x1b2J") // CLS
	//                      1234567890123456789012345678901234567890
	g.Subbios.Stdio.Printf("16 physical colours are available. Mode\n")
	g.Subbios.Stdio.Printf("40 has 16 logical colour slots, meaning\n")
	g.Subbios.Stdio.Printf("all colours can be used simultaneously.\n")
	g.Subbios.Stdio.Printf("Mode 80 only has 4 slots, so you must \n")
	g.Subbios.Stdio.Printf("choose wisely!")
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Physical Colours", 100, 184)
	y := 18
	for c := 0; c <= 15; c++ {
		x := 90
		colourStr := fmt.Sprintf("%d", c)
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, colourStr, x-25, y)
		g.Subbios.TGraphicsOutput.FFillArea(1, 0, c, 0, 0, []int{x, y, x + 60, y, x + 60, y + 10, x, y + 10, x, y})
		g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 15, 0, 0, []int{x, y, x + 60, y, x + 60, y + 10, x, y + 10, x, y})
		switch c {
		case 0:
			colourStr = "Black"
		case 1:
			colourStr = "Dark Blue"
		case 2:
			colourStr = "Dark Red"
		case 3:
			colourStr = "Dark Purple"
		case 4:
			colourStr = "Dark Green"
		case 5:
			colourStr = "Dark Cyan"
		case 6:
			colourStr = "Brown"
		case 7:
			colourStr = "Light Grey"
		case 8:
			colourStr = "Dark Grey"
		case 9:
			colourStr = "Light Blue"
		case 10:
			colourStr = "Light Red"
		case 11:
			colourStr = "Light Purple"
		case 12:
			colourStr = "Light Green"
		case 13:
			colourStr = "Light Cyan"
		case 14:
			colourStr = "Yellow"
		case 15:
			colourStr = "White"
		}
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, colourStr, x+75, y)
		y += 10
		delay()
	}
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue...", 0, 0)
	// TODO: color reassignment and flashing colours
}

func dithers(g *Game) {
	g.Subbios.Stdio.Printf("\x1b2J") // CLS
	//                      1234567890123456789012345678901234567890
	g.Subbios.Stdio.Printf("Both screen modes have 8 dither colours.")
	g.Subbios.Stdio.Printf("This allows for more subtle shading and\n")
	g.Subbios.Stdio.Printf("colouring.")
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Mode 40 Preset Dithers", 80, 195)
	y := 30
	for c := 0; c <= 7; c++ {
		x := 115
		colourStr := fmt.Sprintf("%d", c)
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, colourStr, x-25, y+5)
		g.Subbios.TGraphicsOutput.FFillArea(2, c, 0, 0, 0, []int{x, y, x + 120, y, x + 120, y + 20, x, y + 20, x, y})
		g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 15, 0, 0, []int{x, y, x + 120, y, x + 120, y + 20, x, y + 20, x, y})
		y += 20
		delay()
	}
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue...", 0, 0)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	g.Subbios.Stdio.Printf("\x1b2h")           // Mode 80
	g.Subbios.Stdio.Printf("\x1b80;0;1;0;0~C") // RM Basic Mode 80 colour scheme
	g.Subbios.Stdio.Printf("\x1b80;1;12;0;0~C")
	g.Subbios.Stdio.Printf("\x1b80;2;10;0;0~C")
	g.Subbios.Stdio.Printf("\x1b80;3;15;0;0~C")
	g.Subbios.Stdio.Printf("\x1b33;50m") // Blue background, white text
	//                      12345678901234567890123456789012345678901234567890123456789012345678901234567890
	g.Subbios.Stdio.Printf("You can even define your own dithers. Both screen modes support 8 user-defined\n")
	g.Subbios.Stdio.Printf("dither patterns.  A dither is a 4x4 pattern of colours. For example:\n\n")
	g.Subbios.Stdio.Printf("    2,0,0,3  3,0,0,3  0,0,0,0  2,2,2,2  1,0,0,3  3,2,1,0  0,1,0,1  1,1,1,1\n")
	g.Subbios.Stdio.Printf("    3,2,0,0  0,3,3,0  0,1,1,0  2,3,3,2  0,1,0,0  0,1,2,3  1,1,0,1  1,2,1,1\n")
	g.Subbios.Stdio.Printf("    0,3,2,0  0,3,3,0  0,1,1,0  2,3,3,2  0,0,1,0  3,2,1,0  0,0,0,1  1,1,2,1\n")
	g.Subbios.Stdio.Printf("    0,0,3,2  3,0,0,3  0,0,0,0  2,2,2,2  3,0,0,1  0,1,2,3  1,1,1,1  1,1,1,1\n")
	g.Subbios.TGraphicsOutput.FPlotCharacterString(1, 2, 1, 3, 0, "Preset Dithers", 40, 30)
	y = 13
	for c := 0; c <= 7; c++ {
		x := 95
		colourStr := fmt.Sprintf("%d", c)
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 3, 0, colourStr, x-25, y+5)
		g.Subbios.TGraphicsOutput.FFillArea(2, c, 0, 0, 0, []int{x, y, x + 120, y, x + 120, y + 20, x, y + 20, x, y})
		g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 3, 0, 0, []int{x, y, x + 120, y, x + 120, y + 20, x, y + 20, x, y})
		y += 20
		delay()
	}
	// Set up 8 user-defined dithers
	userDithers := [8][4][4]int{}
	userDithers[0] = [4][4]int{
		{2, 0, 0, 3},
		{3, 2, 0, 0},
		{0, 3, 2, 0},
		{0, 0, 3, 2},
	}
	userDithers[1] = [4][4]int{
		{3, 0, 0, 3},
		{0, 3, 3, 0},
		{0, 3, 3, 0},
		{3, 0, 0, 3},
	}
	userDithers[2] = [4][4]int{
		{0, 0, 0, 0},
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
	}
	userDithers[3] = [4][4]int{
		{2, 2, 2, 2},
		{2, 3, 3, 2},
		{2, 3, 3, 2},
		{2, 2, 2, 2},
	}
	userDithers[4] = [4][4]int{
		{1, 0, 0, 3},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{3, 0, 0, 1},
	}
	userDithers[5] = [4][4]int{
		{3, 2, 1, 0},
		{0, 1, 2, 3},
		{3, 2, 1, 0},
		{0, 1, 2, 3},
	}
	userDithers[6] = [4][4]int{
		{0, 1, 0, 1},
		{1, 1, 0, 1},
		{0, 0, 0, 1},
		{1, 1, 1, 1},
	}
	userDithers[7] = [4][4]int{
		{1, 1, 1, 1},
		{1, 2, 1, 1},
		{1, 1, 2, 1},
		{1, 1, 1, 1},
	}
	for d := 0; d <= 7; d++ {
		g.Subbios.TGraphicsOutput.FSetDitherPattern(d+8, userDithers[d])
	}
	g.Subbios.TGraphicsOutput.FPlotCharacterString(1, 2, 1, 3, 0, "User-Defined", 370, 40)
	y = 13
	for c := 8; c <= 15; c++ {
		x := 425
		colourStr := fmt.Sprintf("%d", c)
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 3, 0, colourStr, x-25, y+5)
		g.Subbios.TGraphicsOutput.FFillArea(2, c, 0, 0, 0, []int{x, y, x + 120, y, x + 120, y + 20, x, y + 20, x, y})
		g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 3, 0, 0, []int{x, y, x + 120, y, x + 120, y + 20, x, y + 20, x, y})
		y += 20
		delay()
	}
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 3, 0, "Press any key to continue...", 0, 0)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	g.Subbios.Stdio.Printf("\x1b0h")     // Mode 40
	g.Subbios.Stdio.Printf("\x1b45;51m") // Blue background, white text
	g.Subbios.Stdio.Printf("\x1b2J")     // CLS
}

func hatchings(g *Game) {
	g.Subbios.Stdio.Printf("\x1b2J") // CLS
	//                      1234567890123456789012345678901234567890
	g.Subbios.Stdio.Printf("Hatchings are another colouring trick.\n")
	g.Subbios.Stdio.Printf("These are 16x16 patterns of zeros and\n")
	g.Subbios.Stdio.Printf("ones, representing a primary and second-")
	g.Subbios.Stdio.Printf("ary colour, respectively. 6 hatchings\n")
	g.Subbios.Stdio.Printf("are available to both screen modes and\n")
	g.Subbios.Stdio.Printf("can be redefined by the user.\n")
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Hatching     Primary Col  Secondary Col", 0, 165)
	y := 30
	col1 := 0
	col2 := 0
	for c := 0; c <= 5; c++ {
		x := 15
		switch c {
		case 0:
			col1 = 13
			col2 = 0
		case 1:
			col1 = 14
			col2 = 8
		case 2:
			col1 = 9
			col2 = 5
		case 3:
			col1 = 6
			col2 = 1
		case 4:
			col1 = 7
			col2 = 4
		case 5:
			col1 = 14
			col2 = 10
		}
		colourStr := fmt.Sprintf("%d", c)
		col1Str := fmt.Sprintf("%d", col1)
		col2Str := fmt.Sprintf("%d", col2)
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, colourStr, x-15, y+5)
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, col1Str, x+130, y+5)
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, col2Str, x+240, y+5)
		g.Subbios.TGraphicsOutput.FFillArea(3, c, col1, col2, 0, []int{x, y, x + 70, y, x + 70, y + 20, x, y + 20, x, y})
		g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 15, 0, 0, []int{x, y, x + 70, y, x + 70, y + 20, x, y + 20, x, y})
		y += 20
		delay()
	}
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue...", 0, 0)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	g.Subbios.TGraphicsOutput.FFillArea(1, 0, 1, 0, 0, []int{0, 0, 319, 0, 319, 164, 0, 164, 0, 0})
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "(user redefined)", 0, 155)
	// Set up 6 user-defined hatchings
	userHatchings := [8][16][16]int{}
	userHatchings[0] = [16][16]int{
		//  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16
		{1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 1
		{0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0}, // 2
		{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0}, // 3
		{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0}, // 4
		{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0}, // 5
		{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0}, // 6
		{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0}, // 7
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0}, // 8
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1}, // 9
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0}, // 10
		{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0}, // 11
		{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0}, // 12
		{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0}, // 13
		{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0}, // 14
		{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0}, // 15
		{0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0}, // 16
	}
	userHatchings[1] = [16][16]int{
		//  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0}, // 1
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0}, // 2
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0}, // 3
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0}, // 4
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0}, // 5
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1}, // 6
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 7
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, // 8
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 9
		{1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 10
		{0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 11
		{0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 12
		{0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 13
		{0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 14
		{0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 15
		{0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 16
	}
	userHatchings[2] = [16][16]int{
		//  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1}, // 1
		{0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0}, // 2
		{0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0}, // 3
		{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0}, // 4
		{0, 1, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0}, // 5
		{0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}, // 6
		{0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0}, // 7
		{0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0}, // 8
		{0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0}, // 9
		{0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0}, // 10
		{0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}, // 11
		{0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0}, // 12
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0}, // 13
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}, // 14
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}, // 15
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, // 16
	}
	userHatchings[3] = [16][16]int{
		//  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // 1
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // 2
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // 3
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // 4
		{0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // 5
		{0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}, // 6
		{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0}, // 7
		{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0}, // 8
		{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}, // 9
		{1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0}, // 10
		{1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0}, // 11
		{1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0}, // 12
		{1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0}, // 13
		{1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0}, // 14
		{1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0}, // 15
		{1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0}, // 16
	}
	userHatchings[4] = [16][16]int{
		{1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1}, // 1
		{0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0}, // 2
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // 3
		{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}, // 4
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0}, // 5
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0}, // 6
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // 7
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // 8
		{1, 1, 1, 0, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1}, // 9
		{0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0}, // 10
		{0, 0, 0, 1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0}, // 11
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}, // 12
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}, // 13
		{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0}, // 14
		{0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // 15
		{0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0}, // 16
	}
	userHatchings[5] = [16][16]int{
		//  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, // 1
		{1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0}, // 2
		{1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0}, // 3
		{1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0}, // 4
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, // 5
		{0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0}, // 6
		{0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0}, // 7
		{0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0}, // 8
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, // 9
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1}, // 10
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1}, // 11
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1}, // 12
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, // 13
		{0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0}, // 14
		{0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0}, // 15
		{0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0}, // 16
	}
	for d := 0; d <= 7; d++ {
		g.Subbios.TGraphicsOutput.FSetHatchingPattern(d, userHatchings[d])
	}
	y = 30
	col1 = 0
	col2 = 0
	for c := 0; c <= 5; c++ {
		x := 15
		switch c {
		case 0:
			col1 = 13
			col2 = 0
		case 1:
			col1 = 14
			col2 = 8
		case 2:
			col1 = 9
			col2 = 5
		case 3:
			col1 = 14
			col2 = 0
		case 4:
			col1 = 15
			col2 = 9
		case 5:
			col1 = 8
			col2 = 7
		}
		colourStr := fmt.Sprintf("%d", c)
		col1Str := fmt.Sprintf("%d", col1)
		col2Str := fmt.Sprintf("%d", col2)
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, colourStr, x-15, y+5)
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, col1Str, x+130, y+5)
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, col2Str, x+240, y+5)
		g.Subbios.TGraphicsOutput.FFillArea(3, c, col1, col2, 0, []int{x, y, x + 70, y, x + 70, y + 20, x, y + 20, x, y})
		g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 15, 0, 0, []int{x, y, x + 70, y, x + 70, y + 20, x, y + 20, x, y})
		y += 20
		delay()
	}
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue...", 0, 0)
}

func circles(g *Game) {
	g.Subbios.Stdio.Printf("\x1b2J") // CLS
	//                      1234567890123456789012345678901234567890
	g.Subbios.Stdio.Printf("Circles or pie slices are easy to draw.\n")
	g.Subbios.Stdio.Printf("Note that the 'stretching' artefact in\n")
	g.Subbios.Stdio.Printf("Mode 80 is automatically corrected.")
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Circles", 50, 185)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Pie Slices", 200, 185)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue...", 0, 0)
	c := 0
	r := 75
	incr := 1
	for i := 0; i <= 10; i++ {
		g.Subbios.TGraphicsOutput.FPieSlice(80, 100, r, 0, 0, c)
		r -= incr
		incr += 1
		c++
		delay()
	}
	c = 0
	r1 := 0
	incr = 30
	for d := 0; d <= 15; d++ {
		g.Subbios.TGraphicsOutput.FPieSlice(240, 100, 75, r1, 6283, c)
		r1 += incr
		incr += 40
		c++
		delay()
	}
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue...", 0, 0)
}

func plotChar(g *Game) {
	g.Subbios.Stdio.Printf("\x1b2J") // CLS
	//                      1234567890123456789012345678901234567890
	g.Subbios.Stdio.Printf("Character strings can be drawn in any\n")
	g.Subbios.Stdio.Printf("colour, size, orientation and font. They")
	g.Subbios.Stdio.Printf("even be")
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 2, 15, 0, "stretched!", 64, 220)
	g.Subbios.TGraphicsOutput.FSetOutputClippingAreaLimits(1, 0, 30, 319, 200)
	g.Subbios.TGraphicsOutput.FSetCurrentOutputClippingArea(1)
	for i := 0; i <= 40; i++ {
		x := rand.Intn(319)
		y := rand.Intn(249)
		r := rand.Intn(3)
		f := rand.Intn(1)
		magX := rand.Intn(4-1) + 1
		magY := rand.Intn(4-1) + 1
		c := rand.Intn(15)
		g.Subbios.TGraphicsOutput.FPlotCharacterString(r, magY, magX, c, f, "Hello!", x, y)
		delay()
	}
	g.Subbios.TGraphicsOutput.FSetCurrentOutputClippingArea(0)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue...", 0, 0)
}

func lines(g *Game) {
	g.Subbios.Stdio.Printf("\x1b2J") // CLS
	//                      1234567890123456789012345678901234567890
	g.Subbios.Stdio.Printf("Line drawing can draw straight lines,\n")
	g.Subbios.Stdio.Printf("outline simple geometric shapes, or 2D\n")
	g.Subbios.Stdio.Printf("shapes of any kind.")
	xo := 80
	yo := 100
	d := 5
	for c := 0; c <= 15; c++ {
		if c == 1 {
			continue // because blue background
		}
		g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, c, 0, 0, []int{xo - d, yo, xo + d, yo, xo, yo + d, xo - d, yo})
		xo = xo + 160
		g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, c, 0, 0, []int{xo, yo, xo, yo + d, xo + d, yo + d, xo + d, yo - d, xo - d, yo - d, xo - d, yo + d})
		xo = 80
		d += 5
		delay()
	}
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue...", 0, 0)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	g.Subbios.Stdio.Printf("\x1b2J") // CLS
	//                      1234567890123456789012345678901234567890
	g.Subbios.Stdio.Printf("Lines can be solid, dithered or stylised")
	g.Subbios.Stdio.Printf("with 4 preset styles. You can even draw\n")
	g.Subbios.Stdio.Printf("with a user-defined style.")
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Line style    Primary col  Secondary Col", 0, 180)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue...", 0, 0)
	y := 40
	col1 := 0
	col2 := 15
	style := []int{}
	var name string
	for l := 2; l <= 6; l++ {
		msg := fmt.Sprintf("%d           %d", col1, col2)
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, msg, 160, y-5)
		switch l {
		case 2:
			name = "dashed"
		case 3:
			name = "dotted"
		case 4:
			name = "dash-dotted"
		case 5:
			name = "irregular"
		case 6:
			name = "user-defined"
		}
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, name, 0, y-15)

		if l == 6 {
			style = []int{1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0}
		}
		g.Subbios.TGraphicsOutput.FPolyLine(l, style, col2, col1, 0, []int{0, y, 90, y})
		col1++
		col2--
		y += 30
		delay()
	}
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue...", 0, 0)
}

func areas(g *Game) {

	drawMountain := func(x, y, size int) {
		colour := rand.Intn(9-7) + 7
		g.Subbios.TGraphicsOutput.FFillArea(1, 0, colour, 0, 0, []int{x, y, x + 2*size, y, x + size, y + 2*size, x, y})
	}

	g.Subbios.Stdio.Printf("\x1b2J") // CLS
	//                      1234567890123456789012345678901234567890
	g.Subbios.Stdio.Printf("Similarly, we can draw filled shapes of\n")
	g.Subbios.Stdio.Printf("arbitrary complexity but in extreme\n")
	g.Subbios.Stdio.Printf("cases it starts to get a bit jittery!\n")
	// outline shape of Poland
	polska := []int{7, 163, 11, 140, 5, 133, 5, 128, 15, 120, 12, 115, 12, 112, 16, 108, 16, 99, 13, 96, 17, 91, 17, 85, 21, 84, 23, 76, 21, 77, 27, 68, 29, 63, 43, 59, 50, 59, 50, 54, 47, 51, 55, 41, 62, 44, 59, 51, 70, 46, 76, 46, 77, 43, 80, 37, 83, 39, 93, 35, 95, 29, 103, 19, 112, 26, 118, 16, 125, 16, 133, 20, 140, 17, 146, 21, 155, 21, 155, 27, 161, 17, 178, 8, 175, 24, 196, 50, 201, 51, 204, 64, 191, 86, 195, 103, 185, 109, 199, 124, 200, 138, 195, 165, 188, 175, 177, 181, 143, 177, 117, 182, 108, 176, 94, 181, 90, 192, 98, 189, 88, 196, 80, 196, 53, 186, 45, 176, 7, 163}
	// reposition
	for i := 0; i < len(polska); i += 2 {
		polska[i] += 50
		polska[i+1] += 10
	}
	g.Subbios.TGraphicsOutput.FFillArea(1, 0, 4, 0, 0, polska)
	delay()
	// mountains
	drawMountain(123+50, 16+10, 7)
	delay()
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Giwont", 133+50, 16+10)
	delay()
	drawMountain(117+50, 17+10, 6)
	delay()
	drawMountain(59+50, 44+10, 4)
	delay()
	drawMountain(36+50, 61+10, 5)
	delay()
	drawMountain(55+50, 65+10, 3)
	delay()
	drawMountain(144+50, 67+10, 2)
	delay()
	// 62, 72
	g.Subbios.TGraphicsOutput.FPieSlice(62+50, 72+10, 2, 0, 0, 2)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Wroclaw", 72+50, 72+10)
	delay()
	g.Subbios.TGraphicsOutput.FPieSlice(120+50, 40+10, 2, 0, 0, 2)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Krakow", 130+50, 40+10)
	delay()
	g.Subbios.TGraphicsOutput.FPieSlice(141+50, 109+10, 4, 0, 0, 2)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "WARSZAWA", 151+50, 109+10)
	delay()
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 15, 0, "Press any key to continue...", 0, 0)
}

func flood(g *Game) {
	g.Subbios.Stdio.Printf("\x1b2J") // CLS
	//                      1234567890123456789012345678901234567890
	g.Subbios.Stdio.Printf("Flood fill is another handy painting\n")
	g.Subbios.Stdio.Printf("tool.\n")
	delay()
	g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 5, 0, 0, []int{0, 208, 319, 208})
	delay()
	g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 13, 0, 0, []int{0, 180, 319, 180})
	delay()
	g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 15, 0, 0, []int{0, 110, 319, 110})
	delay()
	g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 4, 0, 0, []int{200, 110, 320, 110, 280, 175, 250, 135, 230, 160, 200, 110})
	delay()
	g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 14, 0, 0, []int{0, 50, 319, 50})
	delay()
	g.Subbios.TGraphicsOutput.FFloodFillArea(0, 0, 5, 0, 0, 0, 0, 50, 190)
	delay()
	g.Subbios.TGraphicsOutput.FFloodFillArea(0, 0, 13, 0, 0, 0, 0, 50, 120)
	delay()
	g.Subbios.TGraphicsOutput.FFloodFillArea(0, 0, 14, 0, 0, 0, 0, 50, 10)
	delay()
	g.Subbios.TGraphicsOutput.FFloodFillArea(0, 0, 4, 0, 0, 0, 0, 210, 120)
	delay()
	g.Subbios.TGraphicsOutput.FFloodFillArea(0, 0, 9, 0, 0, 0, 0, 50, 60)
	delay()
	g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 0, 0, 0, []int{80, 111, 75, 120, 140, 115, 140, 111, 80, 111})
	delay()
	g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 15, 0, 0, []int{85, 120, 87, 130, 130, 125, 137, 116})
	delay()
	g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 1, 0, 0, []int{110, 128, 111, 134, 117, 132, 120, 127})
	delay()
	g.Subbios.TGraphicsOutput.FFloodFillArea(0, 0, 0, 0, 0, 0, 0, 85, 115)
	delay()
	g.Subbios.TGraphicsOutput.FFloodFillArea(0, 0, 15, 0, 0, 0, 0, 90, 125)
	delay()
	g.Subbios.TGraphicsOutput.FFloodFillArea(0, 0, 1, 0, 0, 0, 0, 115, 130)
	delay()
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 0, 0, "Press any key to continue...", 0, 0)
}

func sprites(g *Game) {
	//                      1234567890123456789012345678901234567890
	g.Subbios.Stdio.Printf("\nAnd let's not forget animated sprites!")
	poses := [][][]int{
		{
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, 15, 15, -1, 15, 15, -1},
			{15, -1, -1, 15, -1, -1, 15},
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
		},
		{
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
			{15, 15, 15, 15, 15, 15, 15},
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
		},
		{
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
			{15, -1, -1, 15, -1, -1, 15},
			{-1, 15, 15, -1, 15, 15, -1},
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
		},
		{
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
			{15, 15, 15, 15, 15, 15, 15},
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
			{-1, -1, -1, -1, -1, -1, -1},
		},
	}
	sprite := subbios.Sprite{HighResolution: false, Hotspot: [2]int{0, 0}, Poses: poses}
	st := &subbios.SaveTable{}
	x := 125
	y := 160
	p := 0
	goingRight := true
	goingUp := false
	g.Subbios.TGraphicsOutput.FDrawSprite(sprite, st, x, y, p, false, 0)
	for {
		time.Sleep(100 * time.Millisecond)
		if goingRight {
			x += 2
		} else {
			x -= 2
		}
		if goingUp {
			y++
		} else {
			y--
		}
		g.Subbios.TGraphicsOutput.FMoveSprite(sprite, st, x, y, p, false, 0)
		p++
		if p > 4 {
			p = 0
		}
		if x > 300 {
			goingRight = false
		}
		if x < 20 {
			goingRight = true
		}
		if y > 200 {
			goingUp = false
		}
		if y < 100 {
			goingUp = true
		}
		if g.Subbios.Stdio.Getch() != 0 {
			break
		}
	}
	g.Subbios.TGraphicsOutput.FEraseSprite(st)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 14, 0, "Press any key to continue...", 0, 0)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 0, 0, "End of demo! Press any key to quit.", 0, 0)
}

func delay() {
	time.Sleep(100 * time.Millisecond)
}

// The main loop
func mainLoop(g *Game) {
	introScreen(g)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	screenModes(g)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	colours(g)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	dithers(g)
	hatchings(g)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	circles(g)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	plotChar(g)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	lines(g)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	areas(g)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	flood(g)
	g.Subbios.Stdio.KeyboardBufferFlush()
	g.Subbios.Stdio.Getchar()
	sprites(g)
	g.Subbios.Stdio.KeyboardBufferFlush()
	for {
		if g.Subbios.Stdio.Getch() != 0 {
			os.Exit(0)
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
