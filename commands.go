package nimgobus

import (
	"image"

	"github.com/StephaneBunel/bresenham"
	"github.com/hajimehoshi/ebiten"
)

// SetMode sets the screen mode. 40 is low-resolution, high-colour mode (320x250) and
// 80 is high-resolutions, low-colour mode (640x250)
func (n *Nimbus) SetMode(columns int) {
	if columns != 40 && columns != 80 {
		// RM Basic would just pass if an invalid column number was given so
		// we'll do the same
		return
	}
	if columns == 40 {
		// low-resolution, high-colour mode (320x250)
		n.paper, _ = ebiten.NewImage(320, 250, ebiten.FilterDefault)
		n.palette = n.defaultLowResPalette
	}
	if columns == 80 {
		// high-resolutions, low-colour mode (640x250)
		n.paper, _ = ebiten.NewImage(640, 250, ebiten.FilterDefault)
		n.palette = n.defaultHighResPalette
	}
	n.cursorPosition = colRow{1, 1}              // Relocate cursor
	n.paper.Fill(n.convertColour(n.paperColour)) // Apply paper colour
	// Redefine textboxes and clear screen
	for i := 0; i < 10; i++ {
		n.textBoxes[i] = textBox{1, 1, columns, 25}
	}
	n.Cls()
}

// PlonkLogo draws the RM Nimbus logo with bottom left-hand corner at (x, y)
func (n *Nimbus) PlonkLogo(x, y int) {
	// Convert position
	_, height := n.logoImage.Size()
	ex, ey := n.convertPos(x, y, height)

	// Draw the logo at the location
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ex, ey)
	n.paper.DrawImage(n.logoImage, op)
}

// SetBorder sets the border colour
func (n *Nimbus) SetBorder(c int) {
	n.validateColour(c)
	n.borderColour = c
}

// PlonkChar plots a character on the paper at a given location
func (n *Nimbus) PlonkChar(c, x, y, colour int) {
	n.drawChar(n.paper, c, x, y, colour)
}

// Plot draws a string of characters on the paper at a given location
// with the colour and size of your choice
func (n *Nimbus) Plot(text string, x, y, xsize, ysize, colour int) {
	// Validate colour
	n.validateColour(colour)
	// Create a new image big enough to contain the plotted chars
	// (without scaling)
	img, _ := ebiten.NewImage(len(text)*10, 10, ebiten.FilterDefault)
	// draw chars on the image
	xpos := 0
	for _, c := range text {
		n.drawChar(img, int(c), xpos, 0, colour)
		xpos += 8
	}
	// Scale img and draw on paper
	ex, ey := n.convertPos(x, y, 10)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(xsize), float64(ysize))
	op.GeoM.Translate(ex, ey)
	n.paper.DrawImage(img, op)
}

// Mode returns the current screen mode (40 column or 80 column)
func (n *Nimbus) Mode() int {
	width, _ := n.paper.Size()
	if width == 320 {
		return 40 // low-res mode 40
	}
	if width == 640 {
		return 80 // high-res mode 80
	}
	return 0 // this never happens
}

// SetWriting selects a textbox if only 1 parameter is passed (index), or
// defines a textbox if 5 parameters are passed (index, col1, row1, col2,
// row2)
func (n *Nimbus) SetWriting(p ...int) {
	// Validate number of parameters
	if len(p) != 1 && len(p) != 5 {
		// invalid
		panic("SetWriting accepts either 1 or 5 parameters")
	}
	if len(p) == 1 {
		// Select textbox - validate choice first then set it
		// and return
		if p[0] < 0 || p[0] > 9 {
			panic("SetWriting index out of range")
		}
		n.selectedTextBox = p[0]
		return
	}
	// Otherwise define textbox if index is not 0
	if p[0] == 0 {
		panic("SetWriting cannot define index zero")
	}
	// Validate column and row values
	for i := 1; i < 5; i++ {
		if p[i] < 0 {
			panic("Negative row or column values are not allowed")
		}
	}
	if p[2] > 25 || p[4] > 25 {
		panic("Row values above 25 are not allowed")
	}
	maxColumns := n.Mode()
	if p[1] > maxColumns || p[3] > maxColumns {
		panic("Column value out of range for this screen mode")
	}
	// Validate passed - set the textbox
	n.textBoxes[p[0]] = textBox{p[1], p[2], p[3], p[4]}
	return
}

// SetPaper sets the paper colour
func (n *Nimbus) SetPaper(c int) {
	n.validateColour(c)
	n.paperColour = c
}

// Cls clears the selected textbox if no parameters are passed, or clears another
// textbox if one parameter is passed
func (n *Nimbus) Cls(p ...int) {
	// Validate number of parameters
	if len(p) != 0 && len(p) != 1 {
		// invalid
		panic("Cls accepts either 0 or 1 parameters")
	}
	// Pick the textbox
	var box textBox
	if len(p) == 0 {
		// No parameters passed so clear currently selected textbox
		box = n.textBoxes[n.selectedTextBox]
	} else {
		// One parameter passed so chose another textbox
		box = n.textBoxes[p[0]]
	}
	// Define bounding rectangle for the textbox
	x1, y1 := n.convertColRow(colRow{box.col1, box.row1})
	x2, y2 := n.convertColRow(colRow{box.col2, box.row2})
	x2 += 8
	y2 += 10
	// Create temp image and fill it with paper colour, then paste on the
	// paper
	img, _ := ebiten.NewImage(int(x2-x1), int(y2-y1), ebiten.FilterDefault)
	img.Fill(n.convertColour(n.paperColour))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x1, y1)
	n.paper.DrawImage(img, op)
}

// SetCurpos sets the cursor position within the selected text box
func (n *Nimbus) SetCurpos(col, row int) {
	// Pick the textbox
	box := n.textBoxes[n.selectedTextBox]
	// Validate col and row position
	if col < 0 || row < 0 {
		panic("Negative column or row values are not allowed")
	}
	width := box.col2 - box.col1
	height := box.row2 - box.row1
	if col > width {
		panic("Column value is outside selected textbox")
	}
	if row > height {
		panic("Row value is outside selected textbox")
	}
	// Validation passed, set cursor position
	n.cursorPosition = colRow{col, row}
}

// SetPen sets the pen colour
func (n *Nimbus) SetPen(c int) {
	n.validateColour(c)
	n.penColour = c
}

// Put plots an ASCII char at the cursor position
func (n *Nimbus) Put(c int) {
	// todo: validate c
	// Pick the textbox
	box := n.textBoxes[n.selectedTextBox]
	width := box.col2 - box.col1
	height := box.row2 - box.row1
	// Get x, y coordinate of cursor and draw the char
	relCurPos := n.cursorPosition
	var absCurPos colRow // we need the absolute cursor position
	absCurPos.col = relCurPos.col + box.col1
	absCurPos.row = relCurPos.row + box.row1
	ex, ey := n.convertColRow(absCurPos)
	ex -= 8
	ey -= 10

	// Draw paper under the char
	img, _ := ebiten.NewImage(8, 10, ebiten.FilterDefault)
	img.Fill(n.convertColour(n.paperColour))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ex, ey)
	n.paper.DrawImage(img, op)
	// Draw the char
	n.drawChar(n.paper, c, int(ex), int(ey), n.penColour)
	// Update relative cursor position
	relCurPos.col++
	// Carriage return?
	if relCurPos.col > width+1 || c == 13 {
		// over the edge so carriage return
		relCurPos.col = 1
		relCurPos.row++
	}
	// Line feed?
	if relCurPos.row > height+1 {
		// move cursor up and scroll textbox
		relCurPos.row--
		// Scroll up:
		// Define bounding rectangle for the textbox
		x1, y1 := n.convertColRow(colRow{box.col1, box.row1})
		x2, y2 := n.convertColRow(colRow{box.col2, box.row2})
		x2 += 8
		y2 += 10
		// Copy actual textbox image
		oldTextBoxImg := n.paper.SubImage(image.Rect(int(x1), int(y1), int(x2), int(y2))).(*ebiten.Image)
		// Create a new textbox image and fill it with paper colour
		newTextBoxImg, _ := ebiten.NewImage(int(x2-x1), int(y2-y1), ebiten.FilterDefault)
		newTextBoxImg.Fill(n.convertColour(n.paperColour))
		// Place old textbox image on new image 10 pixels higher
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, -10)
		newTextBoxImg.DrawImage(oldTextBoxImg, op)
		// Redraw the textbox on the paper
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x1, y1)
		n.paper.DrawImage(newTextBoxImg, op)
	}
	// Set new cursor position
	n.cursorPosition = relCurPos
}

// Print prints a string in the selected textbox
func (n *Nimbus) Print(text string) {
	for _, c := range text {
		n.Put(int(c))
	}
	// Send carriage return
	n.Put(13)
}

// SetColour assigns one of the basic colours to a slot in the current palette
func (n *Nimbus) SetColour(paletteSlot, basicColour int) {
	// Validate basicColour and paletteSlot
	if basicColour < 0 || basicColour > 16 {
		panic("basicColour is out of range")
	}
	n.validateColour(paletteSlot)
	// Validation passed, assign colour
	n.palette[paletteSlot] = basicColour
}

// Line draws connected lines on the screen.  The first n pairs of parameters are
// co-ordinates, and the final parameter is the brush colour.
func (n *Nimbus) Line(p ...int) {
	// Extract colour
	colour := p[len(p)-1]
	// Remove colour parameter
	p = p[:len(p)-1]
	// Use drawLine to draw connected lines
	for i := 0; i < len(p)-2; i += 2 {
		n.drawLine(p[i], p[i+1], p[i+2], p[i+3], colour)
	}
}

// drawLine uses the Bresenham algorithm to draw a straight line on the Nimbus paper
func (n *Nimbus) drawLine(x1, y1, x2, y2, colour int) {
	// convert coordinates
	ex1, ey1 := n.convertPos(x1, y1, 1)
	ex2, ey2 := n.convertPos(x2, y2, 1)
	// create a temp image on which to draw the line
	paperWidth, paperHeight := n.paper.Size()
	dest := image.NewRGBA(image.Rect(0, 0, paperWidth, paperHeight))
	bresenham.Bresenham(dest, int(ex1), int(ey1), int(ex2), int(ey2), n.convertColour(colour))
	// create a copy of the image as an ebiten.image and paste it on to the Nimbus paper
	img, _ := ebiten.NewImageFromImage(dest, ebiten.FilterDefault)
	op := &ebiten.DrawImageOptions{}
	n.paper.DrawImage(img, op)
}
