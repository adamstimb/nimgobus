package nimgobus

import (
	"image"

	"github.com/StephaneBunel/bresenham"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/vector"
)

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
		n.drawChar(img, int(c), xpos, 0, colour, n.charset)
		xpos += 8
	}
	// Scale img and draw on paper
	ex, ey := n.convertPos(x, y, 10)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(xsize), float64(ysize))
	op.GeoM.Translate(ex, ey)
	n.paper.DrawImage(img, op)
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

// Area draws a filled polygon on the screen.  The first n pairs of parameters are
// co-ordinates, and the final parameter is the brush colour.
func (n *Nimbus) Area(p ...int) {
	// Extract colour
	colour := p[len(p)-1]
	// Remove colour parameter
	p = p[:len(p)-1]
	// Use vector to draw the polygon
	var path vector.Path
	ex, ey := n.convertPos(p[0], p[1], 1)
	path.MoveTo(float32(ex), float32(ey)) // Go to start position
	for i := 2; i < len(p)-1; i += 2 {
		ex, ey = n.convertPos(p[i], p[i+1], 1)
		path.LineTo(float32(ex), float32(ey))
	}
	// Is the shape closed?  If not, draw a line back to start position
	if p[len(p)-2] != p[0] || p[len(p)-1] != p[1] {
		// Shape is open so close it
		ex, ey = n.convertPos(p[0], p[1], 1)
		path.MoveTo(float32(ex), float32(ey))
	}
	// Fill the shape on paper
	op := &vector.FillOptions{
		Color: n.convertColour(colour),
	}
	path.Fill(n.paper, op)
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
