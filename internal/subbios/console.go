package subbios

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/adamstimb/nimgobusdev/internal/make2darray"
	"github.com/adamstimb/nimgobusdev/internal/queue"
	"github.com/adamstimb/nimgobusdev/internal/subbios/colour"
	"github.com/hajimehoshi/ebiten/v2"
)

// console holds all the console io malarky.
type console struct {
	curpos                  [2]int             // Cursor position {row, col}
	savedCurpos             [2]int             // SCP stores the current curpos here, RCP and CPR pull the value from here
	penColour               int                // Pen colour
	paperColour             int                // Paper colour
	charSet                 int                // Selected charset
	wordWrap                bool               // Set to true is word wrap is on (Default)
	underlined              bool               // Set to true for underlined chars
	xorWriting              bool               // Set to true for XOR writing
	printControlChars       bool               // Set to true to enable printing of control chars
	cursorUnderlined        bool               // Set to true for underlined cursor
	cursorCharSet           int                // Charset for the cursor
	cursorFlashing          bool               // Set to true for flashing cursor
	cursorDisplayed         bool               // Set to true for displayed cursor, otherwise invisible
	cursorChar              int                // ASCII code of the cursor char
	scrollingArea           [4]int             // Current scrolling area (r1, c1, r2, c2 where r1,c1 == top-left and r2,c2 == bottom-right)
	stdoutBuffer            *queue.Queue[rune] // The buffer of chars (runes) sent to the console
	stdinBuffer             *queue.Queue[rune] // The buffer of ASCII chars received from the keyboard
	stdinBufferIndex        int                // The position of the cursor within the stdin buffer
	lowResColourLookupTable [16][3]int         // The console's Mode 40 colour lookup table {{physicalColour, physicalFlashColour, flashRate}}
	hiResColourLookupTable  [4][3]int          // The console's Mode 80 colour lookup table {{physicalColour, physicalFlashColour, flashRate}}
	v                       *video
}

// Update should be called on each Ebiten Update call
func (c *console) update() {
	// Transfer runs from Ebiten buffer to stdinBuffer - this is mainly to support the Stdio.Scanf() feature.
	newRunes := make([]rune, 1)
	ebiten.AppendInputChars(newRunes[:0])
	// Detect any other keys that the TextBox needs to know about
	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		c.stdinBuffer.Enqueue('\n')
		return
	}
	if repeatingKeyPressed(ebiten.KeyBackspace) {
		c.stdinBuffer.Enqueue('\x08')
		return
	}
	// Left, Right, Up, Down are encoded as x01, x02, x03, x04 respectively.
	if repeatingKeyPressed(ebiten.KeyLeft) {
		c.stdinBuffer.Enqueue('\x01')
		return
	}
	if repeatingKeyPressed(ebiten.KeyRight) {
		c.stdinBuffer.Enqueue('\x02')
		return
	}
	// Keys that were considered a "hit" by Nimbus but don't have any effect on Scanf.  They are queued as \x05.  This
	// is useful for Getchar if applied to a "Press any key..." scenario.
	if repeatingKeyPressed(ebiten.KeyF1) || repeatingKeyPressed(ebiten.KeyF2) || repeatingKeyPressed(ebiten.KeyF3) || repeatingKeyPressed(ebiten.KeyF4) ||
		repeatingKeyPressed(ebiten.KeyF5) || repeatingKeyPressed(ebiten.KeyF6) || repeatingKeyPressed(ebiten.KeyF7) || repeatingKeyPressed(ebiten.KeyF8) ||
		repeatingKeyPressed(ebiten.KeyF9) || repeatingKeyPressed(ebiten.KeyF10) || repeatingKeyPressed(ebiten.KeyF11) || repeatingKeyPressed(ebiten.KeyF12) ||
		repeatingKeyPressed(ebiten.KeyTab) || repeatingKeyPressed(ebiten.KeyUp) || repeatingKeyPressed(ebiten.KeyDown) || repeatingKeyPressed(ebiten.KeyHome) ||
		repeatingKeyPressed(ebiten.KeyEnd) {
		c.stdinBuffer.Enqueue('\x07')
		return
	}
	// Detect printable chars
	for _, r := range newRunes {
		if r != 0 {
			c.stdinBuffer.Enqueue(r)
		}
	}
}

// getScrollingAreaSize returns the height and width of the scrolling area
func (c *console) getScrollingAreaSize() (height, width int) {
	height = c.scrollingArea[2] - (c.scrollingArea[0] - 1)
	width = c.scrollingArea[3] - (c.scrollingArea[1] - 1)
	return height, width
}

// convertCurposToXY converts the current cursor position (relative to the scrolling area) and
// returns the absolute x, y position of the cursor.
func (c *console) convertCurposToXY() (x, y int) {
	absR := c.curpos[0] + (c.scrollingArea[0] - 1)
	absC := c.curpos[1] + (c.scrollingArea[1] - 1)
	x = (absC - 1) * 8
	y = (25 - absR) * 10
	//y = 249 - (absR * 10)
	return x, y
}

// convertAnyCurposToXY converts any given cursor position (relative to the scrolling area) and
// returns the absolute x, y position of the cursor.
func (c *console) convertAnyCurposToXY(row, col int) (x, y int) {
	absR := row + (c.scrollingArea[0] - 1)
	absC := col + (c.scrollingArea[1] - 1)
	x = (absC - 1) * 8
	y = (25 - absR) * 10
	//y = 250 - (absR * 10)
	return x, y
}

// syncVideoColourTable updates the video colour table with settings from the console colour lookup table
func (c *console) syncVideoColourTable() {
	if c.v.screenWidth == 40 {
		// sync from low-res table
		for i := 0; i < 16; i++ {
			c.v.colourLookupTable[i] = colour.CltElement{
				FirstPhysicalColour:  c.lowResColourLookupTable[i][0],
				FlashSpeed:           c.lowResColourLookupTable[i][2],
				SecondPhysicalColour: c.lowResColourLookupTable[i][1],
			}
		}
	} else {
		// sync from high-res table
		for i := 0; i < 4; i++ {
			c.v.colourLookupTable[i] = colour.CltElement{
				FirstPhysicalColour:  c.hiResColourLookupTable[i][0],
				FlashSpeed:           c.hiResColourLookupTable[i][2],
				SecondPhysicalColour: c.hiResColourLookupTable[i][1],
			}
		}
	}
}

func (c *console) lineFeed() {
	c.curpos[0]++ // Move down 1 row
	h, _ := c.getScrollingAreaSize()
	if c.curpos[0] > h {
		c.scrollUp(1)
		c.curpos[0]--
	}
}

func (c *console) carriageReturn() {
	// Move to left-side of scrolling area
	c.curpos[1] = 1
}

// tab moves the cursor to the right by one tab
func (c *console) tab() {

	// nextMultipleOfFour returns the next highest integer that is a multiple of 4.
	nextMultipleOfFour := func(n int) int {
		remainder := n % 4
		if remainder == 0 {
			return n + 4 // n is already a multiple of 4, i.e. already on a tab so return next tab
		}
		return n + 4 - remainder
	}

	delImg := make2darray.Make2dArray(8, 10, c.paperColour)
	nextCol := nextMultipleOfFour(c.curpos[1])
	_, w := c.getScrollingAreaSize()

	if nextCol > w {
		nextCol = 1
	}

	for c.curpos[1] != nextCol {
		x, y := c.convertCurposToXY()
		c.v.drawFeature(feature{pixels: delImg, x: x, y: y, colour: -1, xor: c.xorWriting})
		c.cursorForward(1)
	}

}

// executeEscapeSequence(params, esType) executes an escape sequence
func (c *console) executeEscapeSequence(params []int, esType string) {

	switch esType {
	case "_H":
		c.cursorPosition(params[0], params[1])
	case "_s":
		c.saveCursorPosition()
	case "_u":
		c.restoreCursorPosition()
	case "~B":
		c.defineScrollingArea(params[0], params[1], params[2], params[3])
	case "~E":
		c.setCharacterAttribute(params[0], params[1], params[2])
	case "~A":
		c.setCursorMode(params[0], params[1], params[2], params[3], params[4])
	case "~G":
		c.cursorVisible()
	case "~F":
		c.cursorNotVisible()
	case "_h":
		c.setMode(params[0])
	case "_m":
		c.setGraphicsRendition(params)
	case "~C":
		c.setColourLookupTable(params[0], params[1], params[2], params[3], params[4])
	case "_l":
		c.resetMode(params[0])
	case "_c":
		c.resetToInitialState()
	case "_C":
		c.cursorForward(params[0])
	case "_D":
		c.cursorBackward(params[0])
	case "_S":
		c.scrollUp(params[0])
	case "_T":
		c.scrollDown(params[0])
	case "_A":
		c.cursorUp(params[0])
	case "_B":
		c.cursorDown(params[0])
	case "_J":
		c.eraseInDisplay(params[0])
	case "_K":
		c.eraseInLine(params[0])
	case "~D":
		c.enablePrintControlChars(params[0])
	}
}

// parseEscapeSequence continues to flush the buffer but interprets it
// as an escape sequence
func (c *console) parseEscapeSequence() bool {
	value := make([]rune, 3) // parameters are first collected as rune values
	value = []rune{' ', ' ', ' '}
	valueIndex := 0
	params := make([]int, 16) // if parsed successfully each parameter can then be stored in an array of ints
	params = []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
	paramIndex := 0
	esType := make([]rune, 2) // the escape code type is then stored her (the 1st rune is reserved for ~)
	// Valid, so go ahead
	for c.stdoutBuffer.Size() > 0 {
		if r, ok := c.stdoutBuffer.Dequeue(); ok {
			if r == '[' {
				// new sequence started
				value = []rune{' ', ' ', ' '}
				valueIndex = 0
				params = make([]int, 16)
				params = []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
				paramIndex = 0
				esType = make([]rune, 2)
				continue
			}
			if unicode.IsDigit(r) {
				// collecting value
				if valueIndex > 3 {
					// value is too big therefore invalid sequence
					return false
				}
				value[valueIndex] = r
				valueIndex++
				continue
			}
			if unicode.IsLetter(r) || r == ';' {
				// convert and store current value and move to next
				// catch empty/unset value
				if string(value) == "   " {
					value = []rune{'-', '1'}
				}
				p, err := strconv.Atoi(strings.TrimSpace(string(value)))
				if err != nil {
					panic(err)
				}
				params[paramIndex] = p
				valueIndex = 0
				paramIndex++
				value = []rune{' ', ' ', ' '}
				if paramIndex > 31 {
					// too many params therefore invalid
					return false
				}
				// skip if ;, otherwise continue to parse the rune
				if r == ';' {
					continue
				}
			}
			if r == '~' {
				// collect ~ prefix for esType
				esType[0] = r
				continue
			}
			if unicode.IsLetter(r) {
				// get esType and execute
				esType[1] = r
				// tidy-up esType
				if esType[0] == 0 {
					esType[0] = '_' // otherwise evaluator gets into trouble doing the string matching
				}
				c.executeEscapeSequence(params, string(esType))
				// Is there another sequence following? Return if not.
				if val, ok := c.stdoutBuffer.Peek(); ok {
					if val != '[' {
						return true
					} else {
						continue
					}
				} else {
					return true
				}
			}
			// must be invalid if anything else crops up
			return false
		}
	}
	// end of the world
	return false
}

// flushStdoutBuffer flushes all runes from the consoleBuffer
func (c *console) flushStdoutBuffer() {
	for c.stdoutBuffer.Size() > 0 {
		// TODO: move flushStdOutBuffer to stdio?
		//if c.sio.gotAnyInterrupts() {
		//	return
		//}
		if r, ok := c.stdoutBuffer.Dequeue(); ok {
			// detect control chars
			switch r {
			case '\x1b':
				// ESC
				c.parseEscapeSequence()
				continue
			case '\t':
				c.tab()
				continue
			case '\n':
				c.lineFeed()
				c.carriageReturn()
				continue
			}
			// Control char
			if r < 32 && !c.printControlChars {
				r = 32
			}
			// Otherwise plonk the char
			x, y := c.convertCurposToXY()
			oldRow := c.curpos[0]
			c.v.plonkChar(int(r), x, y, c.penColour, c.paperColour, c.charSet, c.xorWriting, c.underlined)
			c.cursorForward(1)
			if c.curpos[0] == oldRow && c.curpos[1] == 1 {
				// scroll up required
				c.scrollUp(1)
			}
		}
	}
}
