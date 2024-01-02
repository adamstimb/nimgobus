package subbios

import (
	"github.com/adamstimb/nimgobusdev/internal/make2darray"
)

// All the ANSI-style escape sequences for the console are implemented here.
// Note that as with real escape sequences there is no error reporting - it either
// works or it doesn't lol.

// Where parameters are optional, pass -1 if no value was given by the caller.

func (c *console) defineScrollingArea(r1, c1, r2, c2 int) {
	// Validate scroll area is fully within the screen and r1c1 is top-left and r2c2 is bottom-right
	if (r1 < 1 || r1 > r2) || (r2 > 25) || (c1 < 1 || c1 > c2) || (c2 > c.v.screenWidth) {
		return
	}
	// Set scroll area and send cursor home
	c.scrollingArea = [4]int{r1, c1, r2, c2}
	c.curpos = [2]int{1, 1}
}

func (c *console) setCharacterAttribute(m, n, p int) {
	// Unset params have no effect in this sequence
	if m == 0 {
		c.underlined = false
	}
	if m == 1 {
		c.underlined = true
	}
	if n == 0 {
		c.charSet = 0
	}
	if n == 1 {
		c.charSet = 1
	}
	if p == 0 {
		c.xorWriting = false
	}
	if p == 1 {
		c.xorWriting = true
	}
}

func (c *console) setCursorMode(n, q, r, m, p int) {
	// Unset params have no effect in this sequence
	if n == 0 {
		c.cursorUnderlined = false
	}
	if n == 1 {
		c.cursorUnderlined = true
	}
	if q == 0 {
		c.cursorCharSet = 0
	}
	if q == 1 {
		c.cursorCharSet = 1
	}
	if r == 0 {
		c.cursorFlashing = false
	}
	if r == 1 {
		c.cursorFlashing = true
	}
	if m >= 0 && m <= 255 {
		c.cursorChar = m
	}
	if p == 0 {
		c.cursorDisplayed = true
	}
	if p == 1 {
		c.cursorDisplayed = false
	}
}

func (c *console) cursorVisible() {
	c.cursorDisplayed = true
}

func (c *console) cursorNotVisible() {
	c.cursorDisplayed = false
}

func (c *console) setGraphicsRendition(params []int) {
	for _, p := range params {
		switch p {
		case 0:
			// all attributes off
			c.underlined = false
			c.charSet = 0
			c.paperColour = 0
			if c.v.screenWidth == 80 {
				c.penColour = 1
			} else {
				c.penColour = 7
			}
			continue
		case 4:
			// underline on
			c.underlined = true
			continue
		case 10:
			// standard charset
			c.charSet = 0
			continue
		case 11:
			// alternative charset
			c.charSet = 1
			continue
		case 24:
			// underline off
			c.underlined = false
			continue
		}
		// handle foreground/background colours
		if c.v.screenWidth == 80 && p >= 30 && p <= 33 {
			c.penColour = p - 30
			continue
		}
		if c.v.screenWidth == 40 && p >= 30 && p <= 45 {
			c.penColour = p - 30
			continue
		}
		if c.v.screenWidth == 80 && p >= 50 && p <= 53 {
			c.paperColour = p - 50
			continue
		}
		if c.v.screenWidth == 40 && p >= 50 && p <= 65 {
			c.paperColour = p - 50
			continue
		}
	}
}

func (c *console) setColourLookupTable(q, n, m, f, p int) {
	// validate
	if q != 40 && q != 80 {
		return
	}
	if f < 0 || f > 2 {
		return
	}
	highestColour := 15
	if q == 80 {
		highestColour = 3
	}
	if n < 0 || m < 0 || p < 0 || n > highestColour || m > 15 || p > highestColour {
		return
	}

	// set
	if q == 40 {
		// low-res table
		c.lowResColourLookupTable[n] = [3]int{m, p, f}
	} else {
		// high-res table
		c.hiResColourLookupTable[n] = [3]int{m, p, f}
	}
	c.syncVideoColourTable()
}

func (c *console) setMode(n int) {
	switch n {
	case 0:
		// 40 column mode, clear screen, home cursor, use low-res CLT
		c.v.waitForEmptyDrawQueue()
		c.v.resetVideoMemory()
		c.v.screenWidth = 40
		c.syncVideoColourTable()
		c.scrollingArea = [4]int{1, 1, 25, 40}
		c.v.resetClippingAreas()
		c.curpos = [2]int{1, 1}
		return
	case 2:
		// 80 column mode, clear screen, home cursor, use high-res CLT
		c.v.waitForEmptyDrawQueue()
		c.v.resetVideoMemory()
		c.v.screenWidth = 80
		c.syncVideoColourTable()
		c.scrollingArea = [4]int{1, 1, 25, 80}
		c.v.resetClippingAreas()
		c.curpos = [2]int{1, 1}
		return
	case 7:
		// set word wrap on
		c.wordWrap = true
		return
	}
}

func (c *console) resetMode(n int) {
	switch n {
	case 0:
		// 40 column mode, clear screen, home cursor, use low-res CLT
		c.v.waitForEmptyDrawQueue()
		c.v.resetVideoMemory()
		c.syncVideoColourTable()
		c.v.screenWidth = 40
		c.scrollingArea = [4]int{1, 1, 25, 40}
		c.curpos = [2]int{1, 1}
		return
	case 2:
		// 80 column mode, clear screen, home cursor, use high-res CLT
		c.v.waitForEmptyDrawQueue()
		c.v.resetVideoMemory()
		c.syncVideoColourTable()
		c.v.screenWidth = 80
		c.scrollingArea = [4]int{1, 1, 25, 80}
		c.curpos = [2]int{1, 1}
		return
	case 7:
		// set word wrap off
		c.wordWrap = false
		return
	}
}

func (c *console) resetToInitialState() {
	c.v.waitForEmptyDrawQueue()
	c.v.resetVideoMemory()
	c.lowResColourLookupTable = [16][3]int{
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
	}
	c.hiResColourLookupTable = [4][3]int{
		{0, 0, 0},
		{7, 0, 0},
		{15, 0, 1},
		{15, 0, 0},
	}
	c.syncVideoColourTable()
	c.v.screenWidth = 80
	c.penColour = 1
	c.paperColour = 0
	c.charSet = 0
	c.wordWrap = true
	c.cursorChar = 95
	c.cursorCharSet = 0
	c.cursorUnderlined = true
	c.cursorFlashing = true
	c.cursorDisplayed = true
	c.scrollingArea = [4]int{1, 1, 25, 80}
	c.curpos = [2]int{1, 1}
	c.savedCurpos = [2]int{1, 1}
	c.stdoutBuffer.Reset()
	c.stdinBuffer.Reset()
	c.stdinBufferIndex = 0
}

func (c *console) scrollUp(n int) {
	// Define bounding rectangle for the textbox
	h, w := c.getScrollingAreaSize()
	x1, y1 := c.convertAnyCurposToXY(1, 1)
	x2, y2 := c.convertAnyCurposToXY(h, w)
	y1 += 10
	x2 += 8

	// We have to manipulate videoMemory itself next, so wait for drawQueue to empty and get the drawQueue lock
	c.v.waitForEmptyDrawQueue()

	// Draw queue is empty and locked, lock video memory
	c.v.muDrawQueue.Lock()
	c.v.muMemory.Lock()

	// Copy the textbox segment of videoMemory
	textBoxImg := make2darray.Make2dArray((x2-x1)+1, y1-y2, -1)
	for y := y2; y < y1; y++ {
		textBoxImg[(len(textBoxImg)-1)-(y-y2)] = c.v.memory[249-y][x1:x2]
	}
	// Empty paper on bottom row of textbox
	paperImg := make2darray.Make2dArray((x2 - x1), 10, c.paperColour)

	// Unlock everything and send images
	c.v.muDrawQueue.Unlock()
	c.v.muMemory.Unlock()
	c.v.drawFeature(feature{pixels: textBoxImg[10:], x: x1, y: y2 + 10, colour: -1, xor: false})
	c.v.drawFeature(feature{pixels: paperImg, x: x1, y: y2, colour: -1, xor: false})
}

func (c *console) scrollDown(n int) {
	// Define bounding rectangle for the textbox
	h, w := c.getScrollingAreaSize()
	x1, y1 := c.convertAnyCurposToXY(1, 1)
	x2, y2 := c.convertAnyCurposToXY(h, w)
	y1 += 10
	x2 += 8

	// We have to manipulate videoMemory itself next, so wait for drawQueue to empty and get the drawQueue lock
	c.v.waitForEmptyDrawQueue()

	// Draw queue is empty and locked, lock video memory
	c.v.muDrawQueue.Lock()
	c.v.muMemory.Lock()

	// Copy the textbox segment of videoMemory
	textBoxImg := make2darray.Make2dArray((x2-x1)+1, y1-y2, -1)
	for y := y2; y < y1; y++ {
		textBoxImg[(len(textBoxImg)-1)-(y-y2)] = c.v.memory[249-y][x1:x2]
	}

	// Truncate bottom row of textbox
	truncated := make2darray.Make2dArray(len(textBoxImg[0]), len(textBoxImg)-10, -1)
	for y := 0; y < len(truncated); y++ {
		for x := 0; x < len(truncated[0]); x++ {
			truncated[y][x] = textBoxImg[y][x]
		}
	}

	// Empty paper on *top* row of textbox
	paperImg := make2darray.Make2dArray((x2 - x1), 10, c.paperColour)

	// Unlock everything and send images
	c.v.muDrawQueue.Unlock()
	c.v.muMemory.Unlock()
	c.v.drawFeature(feature{pixels: truncated, x: x1, y: y2, colour: -1, xor: false})
	c.v.drawFeature(feature{pixels: paperImg, x: x1, y: y1 - 10, colour: -1, xor: false})
}

func (c *console) cursorForward(n int) {
	// Handle unset params
	if n < 0 {
		n = 1
	}
	// Move cursor forward n steps
	h, w := c.getScrollingAreaSize()
	for i := 0; i < n; i++ {
		c.curpos[1]++
		if c.curpos[1] > w {
			if c.wordWrap {
				// move the cursor to the next line down (if there is one)
				c.curpos[1] = 1
				c.curpos[0]++
				if c.curpos[0] > h {
					c.curpos[0] = h
				}
			} else {
				// but if word wrap is off keep cursor in original position
				c.curpos[1]--
			}
		}
	}
}

func (c *console) cursorBackward(n int) {
	// Handle unset params
	if n < 0 {
		n = 1
	}
	// Move cursor back n steps
	_, w := c.getScrollingAreaSize()
	for i := 0; i < n; i++ {
		c.curpos[1]--
		if c.curpos[1] < 1 {
			c.curpos[0]--
			c.curpos[1] = w
		}
		if c.curpos[0] < 1 {
			c.curpos[0] = 1
		}
	}
	// TODO: Special case when editing input buffer and scrollup required
}

func (c *console) cursorUp(n int) {
	// Handle unset params
	if n < 0 {
		n = 1
	}
	// TODO: Special case when editing input buffer and scrollup required (*and* handle tab char)
	if c.curpos[0] == 1 {
		return
	}
	startCol := c.curpos[1]
	for i := 0; i < n; i++ {
		c.cursorBackward(1)
		for startCol != c.curpos[1] {
			c.cursorBackward(1)
		}
		if c.curpos[0] == 1 {
			return
		}
	}
}

func (c *console) cursorDown(n int) {
	// Handle unset params
	if n < 0 {
		n = 1
	}
	// TODO: Special case when editing input buffer and scrolldown required (*and* handle tab char)
	if c.curpos[0] == c.scrollingArea[2] {
		return
	}
	h, _ := c.getScrollingAreaSize()
	startCol := c.curpos[1]
	for i := 0; i < n; i++ {
		c.cursorForward(1)
		for startCol != c.curpos[1] {
			c.cursorForward(1)
		}
		if c.curpos[0] == h {
			return
		}
	}
}

func (c *console) cursorPosition(r1, c1 int) {
	// r1 and c1 are relative position to the scrolling area so we need the absolute position:
	row := r1 - 1 + c.scrollingArea[0]
	col := c1 - 1 + c.scrollingArea[1]
	// If outside scrolling area, cursor goes home
	if (row < c.scrollingArea[0] || row > c.scrollingArea[2]) && (col < c.scrollingArea[1] || col > c.scrollingArea[3]) {
		c.curpos = [2]int{1, 1}
		return
	}
	// If row is within but col is outside go to r1 and left-most column
	if (row >= c.scrollingArea[0] && row <= c.scrollingArea[2]) && (col < c.scrollingArea[1] || col > c.scrollingArea[3]) {
		c.curpos = [2]int{r1, 1}
		return
	}
	// If col is within but row is outside go to c1 and top row:
	if (col >= c.scrollingArea[1] && col <= c.scrollingArea[3]) && (row < c.scrollingArea[0] || row > c.scrollingArea[3]) {
		c.curpos = [2]int{1, c1}
		return
	}
	// Otherwise set curpos with given params
	c.curpos = [2]int{r1, c1}
}

func (c *console) saveCursorPosition() {
	c.savedCurpos = [2]int{c.curpos[0], c.curpos[1]}
}

func (c *console) restoreCursorPosition() {
	c.curpos = [2]int{c.savedCurpos[0], c.savedCurpos[1]}
}

func (c *console) cursorPositionReport() (row, col int) {
	return c.curpos[0], c.curpos[1]
}

func (c *console) eraseInDisplay(n int) {
	c.v.waitForEmptyDrawQueue()
	h, w := c.getScrollingAreaSize()
	switch n {
	case 0:
		// Erase chars from the cursor position to the end of the row (including the
		// char at the cursor position), and all rows below. Cursor does not move.

		// First call eraseInLine(0) to erase current row
		c.eraseInLine(0)

		// Then wipe all rows below if there are any:
		if c.curpos[0] == h {
			return
		}

		// Define bounding rectangle for the section to wipe:
		x1, y1 := c.convertAnyCurposToXY(c.curpos[0], 1)
		x2, y2 := c.convertAnyCurposToXY(h, w)
		x2 += 8

		// Draw empty paper
		paperImg := make2darray.Make2dArray((x2 - x1), (y1 - y2), c.paperColour)
		c.v.drawFeature(feature{pixels: paperImg, x: x1, y: y2, colour: -1, xor: false})
	case 1:
		// Erase chars from the beginning of the row to the cursor (including the one
		// at the cursor position), and all rows above. Cursor does not move.

		// First call eraseInLine(1) to erase current row
		c.eraseInLine(1)

		// Then wipe all rows above if there are any:
		if c.curpos[0] == 1 {
			return
		}

		// Define bounding rectangle for the section to wipe:
		x1, y1 := c.convertAnyCurposToXY(1, 1)
		x2, y2 := c.convertAnyCurposToXY(h, w)
		y1 += 10
		x2 += 8

		// Draw empty paper
		paperImg := make2darray.Make2dArray((x2 - x1), (y1 - y2), c.paperColour)
		c.v.drawFeature(feature{pixels: paperImg, x: x1, y: y2, colour: -1, xor: false})
	case 2:
		// Erase entire scrolling area and send cursor home.

		// Define bounding rectangle for the textbox:
		x1, y1 := c.convertAnyCurposToXY(1, 1)
		x2, y2 := c.convertAnyCurposToXY(h, w)
		x2 += 8

		// Draw empty paper over textbox and send cursor home
		paperImg := make2darray.Make2dArray((x2 - x1), (y1-y2)+10, c.paperColour)
		c.v.drawFeature(feature{pixels: paperImg, x: x1, y: y2, colour: -1, xor: false})
		c.cursorPosition(1, 1)
	}
}

func (c *console) eraseInLine(n int) {
	_, w := c.getScrollingAreaSize()
	switch n {
	case 0:
		// Erase chars from the cursor position to the end of the row (including the
		// char at the cursor position). Cursor does not move.

		// Define bounding rectangle for the section to wipe:
		x1, y1 := c.convertAnyCurposToXY(c.curpos[0], c.curpos[1])
		x2, y2 := c.convertAnyCurposToXY(c.curpos[0], w)
		y2 += 10
		x2 += 8

		// Draw empty paper
		paperImg := make2darray.Make2dArray((x2 - x1), 10, c.paperColour)
		c.v.drawFeature(feature{pixels: paperImg, x: x1, y: y1, colour: -1, xor: false})
	case 1:
		// Erase chars from the beginning of the row to the cursor (including the one
		// at the cursor position). Cursor does not move.

		// Define bounding rectangle for the section to wipe:
		x1, y1 := c.convertAnyCurposToXY(c.curpos[0], 1)
		x2, y2 := c.convertAnyCurposToXY(c.curpos[0], c.curpos[1])
		y2 += 10
		x2 += 8

		// Draw empty paper
		paperImg := make2darray.Make2dArray((x2 - x1), 10, c.paperColour)
		c.v.drawFeature(feature{pixels: paperImg, x: x1, y: y1, colour: -1, xor: false})
	case 2:
		// Erase whole row. Cursor does not move.

		// Define bounding rectangle for the textbox:
		x1, y1 := c.convertAnyCurposToXY(c.curpos[0], 1)
		x2, y2 := c.convertAnyCurposToXY(c.curpos[0], w)
		y1 += 10
		x2 += 8

		// Draw empty paper
		paperImg := make2darray.Make2dArray((x2 - x1), 10, c.paperColour)
		c.v.drawFeature(feature{pixels: paperImg, x: x1, y: y2, colour: -1, xor: false})
	}
}

func (c *console) enablePrintControlChars(n int) {
	if n == 0 {
		c.printControlChars = true
	}
	if n == 1 {
		c.printControlChars = false
	}
}
