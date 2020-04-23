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
