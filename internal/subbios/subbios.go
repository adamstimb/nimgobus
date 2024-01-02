// Subbios provides all commands to control the Nimbus.  It's similar (but different) to the original
// RM Nimbus SUBBIOS but has the advantage that you don't need to write assembler!
package subbios

import (
	"github.com/adamstimb/nimgobusdev/internal/queue"
	"github.com/adamstimb/nimgobusdev/internal/subbios/errorcode"
	"github.com/hajimehoshi/ebiten/v2"
)

// The subbios commands are grouped by Type as represented by this struct.
type Subbios struct {
	Monitor         *ebiten.Image
	borderSize      int
	FunctionStatus  int
	FunctionError   int
	THardSums       THardSums
	TGraphicsOutput TGraphicsOutput
	TRawConsole     tRawConsole
	TGraphicsInput  TGraphicsInput
	Stdio           Stdio
}

// Initializes the subbios commands.
func (s *Subbios) Init() {
	s.borderSize = 50
	// Initialize video
	s.Monitor = ebiten.NewImage(640+(s.borderSize*2), 500+(s.borderSize*2))
	borderImage := ebiten.NewImage(640+(s.borderSize*2), 500+(s.borderSize*2))
	screenImage := ebiten.NewImage(640, 500)
	// Initialize Subbios functions/devices
	s.FunctionStatus = 0
	s.FunctionError = errorcode.EOk
	s.THardSums = THardSums{s: s}
	s.TGraphicsOutput = TGraphicsOutput{
		s: s, v: &video{
			monitor:      s.Monitor,
			borderImage:  borderImage,
			borderSize:   s.borderSize,
			screenImage:  screenImage,
			borderColour: 0,
		},
		On: false}
	s.TGraphicsOutput.v.loadCharsetImages(0)
	s.TGraphicsOutput.v.loadCharsetImages(1)
	s.TGraphicsInput = TGraphicsInput{
		s: s,
		v: s.TGraphicsOutput.v,
	}
	s.Stdio = Stdio{
		s: s,
		c: &console{
			curpos:       [2]int{1, 1},
			stdoutBuffer: queue.New[rune](),
			stdinBuffer:  queue.New[rune](),
			v:            s.TGraphicsOutput.v,
			lowResColourLookupTable: [16][3]int{
				{0, 0, 0},
				{2, 0, 0},
				{4, 0, 0},
				{6, 0, 0},
				{1, 0, 0},
				{3, 0, 0},
				{5, 0, 0},
				{7, 0, 0},
				{8, 0, 0},
				{10, 0, 0},
				{12, 0, 0},
				{14, 0, 0},
				{9, 0, 0},
				{11, 0, 0},
				{13, 0, 0},
				{15, 0, 0},
			},
			hiResColourLookupTable: [4][3]int{
				{0, 0, 0},
				{7, 0, 0},
				{15, 0, 1},
				{15, 0, 0},
			},
			scrollingArea: [4]int{1, 1, 25, 80},
		},
	}
	s.TGraphicsOutput.v.con = s.Stdio.c
	s.TGraphicsOutput.FGraphicsOutputColdStart()
	s.Stdio.c.resetToInitialState()
	s.loadLogoImage()
	// Start background processes
	go s.TGraphicsOutput.v.colourFlashTicker()
}

// Update needs to be called on each Ebiten update, ideally by Nimbus.Update()
func (s *Subbios) Update() {
	s.TGraphicsOutput.v.update()
	s.TGraphicsInput.update()
	s.Stdio.c.update()
	s.Stdio.checkKeyboardInterrupts()
}
