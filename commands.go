package nimgobus

import (
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
	n.paper.Fill(n.convertColour(n.paperColour)) // Apply paper colour
	// Redefine textboxes
	for i := 0; i < 10; i++ {
		n.textBoxes[i] = textBox{1, 1, 25, columns}
	}
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
	for i := 1; i < 10; i++ {
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
