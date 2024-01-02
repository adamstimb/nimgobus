package main

import (
	"fmt"
	_ "image/png" // import only for side-effects
	"log"
	"math/rand"
	"time"

	"github.com/adamstimb/nimgobus/examples/games/worm/queue"
	nimgobus "github.com/adamstimb/nimgobus/pkg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	launch          int
	nimgobus.Nimbus // Embed the Nimbus in the Game struct
	anyKey          bool
	key             ebiten.Key
	direction       int
	speed           int
	score           int
	maxScore        int
	hitEdge         bool
	demoMode        bool
	keyA            bool
	keyS            bool
	keyK            bool
	keyM            bool
	keyEscape       bool
	keyD            bool
}

func NewGame() *Game {
	game := &Game{}
	game.Init() // Initialize the Nimbus
	return game
}

func isAnyKeyPressed() bool {
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}
	return false
}

func (g *Game) Update() error {
	if g.launch == 0 {
		go App(g) // Launch the Nimbus app on first iteration
	}
	g.launch = 1
	g.Nimbus.Update() // Update the app on all subsequent iterations

	// Synchronise keypresses to wormCounter so the worm game loop doesn't miss anything
	//var anyKey bool
	//var key ebiten.Key

	// Handle direction
	if inpututil.KeyPressDuration(ebiten.KeyA) > 0 {
		g.keyA = true
		g.keyS = false
		g.keyK = false
		g.keyM = false
		g.keyD = false
	}
	if inpututil.KeyPressDuration(ebiten.KeyS) > 0 {
		g.keyA = false
		g.keyS = true
		g.keyK = false
		g.keyM = false
		g.keyD = false
	}
	if inpututil.KeyPressDuration(ebiten.KeyK) > 0 {
		g.keyA = false
		g.keyS = false
		g.keyK = true
		g.keyM = false
		g.keyD = false
	}
	if inpututil.KeyPressDuration(ebiten.KeyM) > 0 {
		g.keyA = false
		g.keyS = false
		g.keyK = false
		g.keyM = true
		g.keyD = false
	}
	if inpututil.KeyPressDuration(ebiten.KeyD) > 0 {
		g.keyA = false
		g.keyS = false
		g.keyK = false
		g.keyM = false
		g.keyD = true
	}
	// Handle other control inputs
	if inpututil.KeyPressDuration(ebiten.KeyEscape) > 0 && inpututil.KeyPressDuration(ebiten.KeyEscape) < 2 {
		g.keyEscape = true
	}
	if isAnyKeyPressed() {
		g.anyKey = true
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

type position struct {
	x int
	y int
}

type worm struct {
	tail       *queue.Queue[position]
	foodAmount int
}

type food struct {
	x      int
	y      int
	amount int
}

// drawArena clears the screen and draws the boundary line
func drawArena(g *Game) {
	g.Subbios.Stdio.Printf("\x1b[2J") // CLS
	g.Subbios.TGraphicsOutput.FPolyLine(1, []int{}, 3, 0, 0, []int{0, 20, 639, 20, 639, 249, 0, 249, 0, 20})
}

// resetControlKeys unsets all the control key flags
func resetControlKeys(g *Game) {
	g.keyA = false
	g.keyS = false
	g.keyK = false
	g.keyM = false
	g.keyEscape = false
	g.keyD = false
	g.anyKey = false
}

// playWorm is the game itself.  demoMode runs on autopilot and will return if
// any key is pressed.  x1, y1, x2, y2 define the bottom-left and upper-right bounds
// of the arena, respectively.  speed is the speed setting.
func playWorm(g *Game, showStatus bool, demoQuitsOnAnyKey bool, x1, y1, x2, y2 int) {

	time.Sleep(500 * time.Millisecond)
	resetControlKeys(g)
	diameter := 4 // worm segment diameter

	gotFirstKey := false

	// demoMode runs on autopilot without showing the stats bar
	if g.demoMode {
		gotFirstKey = true
	}

	// spawn the worm in the middle-ish of the arena
	w := worm{
		tail:       queue.New[position](),
		foodAmount: 5,
	}
	xSpawn := x1 + ((x2 - x1) / 2)
	ySpawn := ((y2 - y1) / 2)
	// quantize
	xSpawn = (xSpawn / 20) * 20
	ySpawn = (ySpawn / 20) * 20

	w.tail.Enqueue(position{x: xSpawn, y: ySpawn})
	currentPosition, _ := w.tail.PeekEnd()
	g.Subbios.TGraphicsOutput.FPieSlice(currentPosition.x, currentPosition.y, diameter, 1, 1, 1)

	// dropFood puts food out for the worm but without hitting it
	dropFood := func() food {
		// in a for loop to make sure random positions don't coincide with the worm
		for {
			f := food{
				x:      rand.Intn((x2-20)-(x1+20)) + x1 + 20,
				y:      rand.Intn((y2-20)-(y1+20)) + y1 + 20,
				amount: rand.Intn(8) + 1,
			}
			// quantize food position
			f.x = (f.x / 20) * 20
			f.y = (f.y / 20) * 20

			retry := false
			for i := 0; i < w.tail.Size(); i++ {
				p, _ := w.tail.PeekAt(i)
				if p.x == f.x && p.y == f.y {
					retry = true
					break
				}
			}
			if retry {
				continue
			} else {
				amount := fmt.Sprintf("%d", f.amount)
				g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 3, 0, amount, f.x-4, f.y-5)
				return f
			}
		}
	}

	// spawn food
	f := dropFood()

	// loop
	autopilot := false
	if g.demoMode {
		// pick a random direction (0,1,2,3 == up, down, left, right)
		g.direction = rand.Intn(3)
	}
	dead := false
	g.key = 0
	elapsedFrames := 0
	for {

		// quit demo on anyKey (e.g. for introscreen)
		if g.demoMode && g.anyKey && demoQuitsOnAnyKey {
			return
		}

		// cancel demo and resume in manual mode (e.g. for demo mode)
		if g.demoMode && g.anyKey && !demoQuitsOnAnyKey {
			g.demoMode = false
			// prevent worm veering off-course
			resetControlKeys(g)
			switch g.direction {
			case 0:
				g.keyK = true
			case 1:
				g.keyM = true
			case 2:
				g.keyA = true
			case 3:
				g.keyS = true
			}
		}

		// prevent autopilot toggle racing when ESC is pressed
		minFrames := 5
		if elapsedFrames < minFrames+1 {
			g.keyEscape = false
			g.anyKey = false
			elapsedFrames++
		} else {
			elapsedFrames = minFrames
		}

		if showStatus {
			var msg string
			if g.demoMode || autopilot {
				msg = fmt.Sprintf("    Score : %5d                Press any key              Maximum : %5d", g.score, g.maxScore)
			} else {
				msg = fmt.Sprintf("    Score : %5d                Speed : %3d                Maximum : %5d", g.score, g.speed, g.maxScore)

			}
			g.Subbios.Stdio.Printf("\x1b25;H")
			g.Subbios.Stdio.Printf(msg)
		}

		// worm death sequence
		if dead {
			if g.score > g.maxScore {
				g.maxScore = g.score
			}
			time.Sleep(time.Duration(200/(g.speed+1)) * time.Millisecond)
			if tailpos, ok := w.tail.DequeueFromEnd(); ok {
				// dying
				g.Subbios.TGraphicsOutput.FPieSlice(tailpos.x, tailpos.y, diameter, 1, 1, 0)
				continue
			} else {
				// dead
				return
			}
		} else {
			// delay to playable speed
			time.Sleep(time.Duration(700/(g.speed+1)) * time.Millisecond)
		}

		// prepare to move head
		currentPosition, _ := w.tail.Peek()
		newPosition := position{}
		var newDirection int

		// change direction with autopilot if autopilot on or in demo mode
		if g.demoMode || autopilot {
			if currentPosition.y < f.y && g.direction != 1 {
				newDirection = 0
			}
			if currentPosition.y > f.y && g.direction != 0 {
				newDirection = 1
			}
			if currentPosition.y < f.y && g.direction == 1 {
				newDirection = 2
			}
			if currentPosition.y > f.y && g.direction == 0 {
				newDirection = 3
			}
			if currentPosition.x < f.x && g.direction != 2 {
				newDirection = 3
			}
			if currentPosition.x > f.x && g.direction != 3 {
				newDirection = 2
			}
			g.direction = newDirection
		} else {
			// take direction from key press but don't let it turn back in to itself
			if g.keyA {
				newDirection = 2
				gotFirstKey = true
			}
			if g.keyS {
				newDirection = 3
				gotFirstKey = true
			}
			if g.keyK {
				newDirection = 0
				gotFirstKey = true
			}
			if g.keyM {
				newDirection = 1
				gotFirstKey = true
			}
			if (newDirection == 0 && g.direction == 1) ||
				(newDirection == 1 && g.direction == 0) ||
				(newDirection == 2 && g.direction == 3) ||
				(newDirection == 3 && g.direction == 2) {
				// no dice
			} else {
				g.direction = newDirection
			}
		}

		// switch off autopilot? not available in demo mode
		if !g.demoMode {
			if g.anyKey && autopilot && elapsedFrames == minFrames {
				autopilot = false
				elapsedFrames = 0
				// prevent worm veering off-course
				resetControlKeys(g)
				switch g.direction {
				case 0:
					g.keyK = true
				case 1:
					g.keyM = true
				case 2:
					g.keyA = true
				case 3:
					g.keyS = true
				}
			}

			// switch on autopilot?
			if g.keyEscape && !autopilot && elapsedFrames == minFrames {
				elapsedFrames = 0
				autopilot = true
				gotFirstKey = true
				g.keyEscape = false
				g.anyKey = false
			}
		}

		// pause if we haven't got any key inputs yet
		if !gotFirstKey {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// apply direction
		switch g.direction {
		case 0:
			// up
			newPosition = position{currentPosition.x, currentPosition.y + 10}
		case 1:
			// down
			newPosition = position{currentPosition.x, currentPosition.y - 10}
		case 2:
			// left
			newPosition = position{currentPosition.x - 20, currentPosition.y}
		case 3:
			// right
			newPosition = position{currentPosition.x + 20, currentPosition.y}
		}

		// Add new head position
		w.tail.InsertAt(0, newPosition)
		// paint over old head and paint new head
		g.Subbios.TGraphicsOutput.FPieSlice(currentPosition.x, currentPosition.y, diameter, 1, 1, 2)
		g.Subbios.TGraphicsOutput.FPieSlice(newPosition.x, newPosition.y, diameter, 1, 1, 1)

		// Got food?
		if newPosition.x == f.x && newPosition.y == f.y {
			// score that food
			w.foodAmount += f.amount
			g.score += f.amount
			// drop new food
			f = dropFood()

		}

		// move tail or eat food
		if w.foodAmount == 0 {
			tailPosition, _ := w.tail.PeekEnd()
			// no food left so move tail
			g.Subbios.TGraphicsOutput.FPieSlice(tailPosition.x, tailPosition.y, diameter, 1, 1, 0)
			w.tail.DequeueFromEnd()
		} else {
			// eat food instead
			w.foodAmount--
		}

		// Hit wall?
		if (newPosition.x <= x1 || newPosition.x >= x2) || (newPosition.y <= y1 || newPosition.y >= y2) {
			// dead
			g.Subbios.TGraphicsOutput.FPieSlice(f.x, f.y, 5, 1, 1, 0) // clean up the leftover food!
			dead = true
			g.hitEdge = true
			resetControlKeys(g)
		}

		// Hit tail?
		for i := 1; i < w.tail.Size(); i++ {
			p, _ := w.tail.PeekAt(i)
			if p.x == newPosition.x && p.y == newPosition.y {
				// dead
				g.Subbios.TGraphicsOutput.FPieSlice(f.x, f.y, 5, 1, 1, 0) // clean up the leftover food!
				dead = true
				g.hitEdge = false
				resetControlKeys(g)
			}
		}

	}

}

// introScreen draws the intro screen with roaming worm
func introScreen(g *Game) {
	// Set screen and colours
	g.Subbios.Stdio.Printf("\x1b[2h") // Mode 80
	g.Subbios.Stdio.Printf("\x1b[~F") // Hide cursor
	g.Subbios.TGraphicsOutput.FGraphicsOutputColdStart()
	g.Subbios.TGraphicsOutput.FSetBorderColour(0)
	g.Subbios.TGraphicsOutput.FSetCltElement(0, 0, 0, 0)
	g.Subbios.TGraphicsOutput.FSetCltElement(1, 12, 0, 0)
	g.Subbios.TGraphicsOutput.FSetCltElement(2, 10, 0, 0)
	g.Subbios.TGraphicsOutput.FSetCltElement(3, 14, 0, 0)

	// Title and keys
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 4, 7, 2, 0, "WORM", 100, 200)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 2, 2, 1, 0, "A - left", 100, 150)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 2, 2, 1, 0, "S - right", 100, 130)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 2, 2, 1, 0, "K - up", 100, 110)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 2, 2, 1, 0, "M - down", 100, 90)

	// Print instructions
	g.Subbios.Stdio.Printf("\x1b[33m")   // Yellow text
	g.Subbios.Stdio.Printf("\x1b[21;1H") // Move cursor to 22,1
	g.Subbios.Stdio.Printf("Move the worm to catch the numbers without hitting the side\n\n")
	g.Subbios.Stdio.Printf("or your own tail.  Press <ESC> at any time for auto-pilot\n\n")
	g.Subbios.Stdio.Printf("                Press any key when ready (D - demo)...")
}

func fullGame(g *Game) {
	drawArena(g)
	g.anyKey = false
	g.score = 0
	time.Sleep(500 * time.Millisecond)
	playWorm(g, true, false, 0, 20, 640, 250)
}

func endScreen(g *Game, bestScoreSoFar bool) {
	g.Subbios.Stdio.Printf("\x1b[2J") // CLS
	var msg string
	if g.hitEdge {
		msg = "You hit the edge !"
	} else {
		msg = "You hit your tail !"
	}
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 2, 2, 1, 0, msg, 120, 220)
	msg = fmt.Sprintf("You scored %d", g.score)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 2, 2, 1, 0, msg, 140, 180)
	if bestScoreSoFar {
		g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 2, 2, 1, 0, "This is the best score so far !", 80, 140)
	}
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 2, 2, 1, 0, "Another game (Y, N, D-demo) ? ", 60, 100)
	g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 2, 2, 1, 0, "Or change speed (C) ? ", 60, 60)
}

func App(g *Game) {
	rand.Seed(time.Now().UnixNano())
	g.Subbios.Boot() // yeaahhh
	resetControlKeys(g)

	// Intro screen with worm running in demo mode
	introScreen(g)
	g.speed = 5
	g.maxScore = 0
	g.demoMode = true
	for !g.anyKey {
		playWorm(g, false, true, 320, 60, 640, 250)
	}

	// Handle demo mode
	if g.keyD {
		g.demoMode = true
		resetControlKeys(g)
	} else {
		g.demoMode = false
	}

	time.Sleep(500 * time.Millisecond)

	// Main loop - keep playing and recording max score until user quits
	bestScoreSoFar := false
	for {
		previousMaxScore := g.maxScore
		fullGame(g)

		// New best score?
		if g.maxScore > previousMaxScore {
			bestScoreSoFar = true
		} else {
			bestScoreSoFar = false
		}
		time.Sleep(500 * time.Millisecond)

		// display end message and choices
		endScreen(g, bestScoreSoFar)

		// skip choices and play again if in demo mode
		if g.demoMode {
			continue
		}

		// get choice
		for {
			g.Subbios.Stdio.KeyboardBufferFlush()
			choice := g.Subbios.Stdio.Getchar()
			if choice == 'y' || choice == 'Y' {
				// new game
				g.demoMode = false
				break
			}
			if choice == 'd' || choice == 'D' {
				// new game in demo mode
				g.demoMode = true
				break
			}
			if choice == 'n' || choice == 'N' {
				// Simple goodbye message for now
				g.Subbios.Stdio.Printf("\x1b[c")
				//                      12345678901234567890123456789012345678901234567890123456789012345678901234567890
				g.Subbios.Stdio.Printf("Thank you for trying Worm!  This is a test deployment that can change or \n")
				g.Subbios.Stdio.Printf("disappear without warning.  But fear not: Trains will be next ...\n\n")
				g.Subbios.Stdio.Printf("Signed\n\n")
				g.Subbios.Stdio.Printf("The Management.\n")
				return
			}
			if choice == 'c' || choice == 'C' {
				// Change speed
				g.Subbios.TGraphicsOutput.FPlotCharacterString(0, 2, 2, 3, 0, "Speed (0..9) ?", 0, 0)
				for {
					g.Subbios.Stdio.KeyboardBufferFlush()
					choice = g.Subbios.Stdio.Getchar()
					value := int(choice - '0')
					// valid choice?  Try again if not
					if value >= 0 && value <= 9 {
						// is valid - set and break loop
						g.speed = value
						break
					}
				}
				// new game with new speed setting
				break
			}
		}
	}
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
	ebiten.SetWindowTitle("Nimbus Worms")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Create a new game and pass it to RunGame method
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
