package nimgobus

import (
	"image"

	"github.com/StephaneBunel/bresenham"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/vector"
)

// PlonkLogo draws the RM Nimbus logo
func (n *Nimbus) PlonkLogo(x, y int) {
	// Convert position
	_, height := n.logoImage.Size()
	ex, ey := n.convertPos(x, y, height)

	// Draw the logo at the location
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ex, ey)
	n.paper.DrawImage(n.logoImage, op)
}

// PlotOptions describes optional parameters for the Plot command.  Plot will
// interpret zero values for SizeX and SizeY as 1.
type PlotOptions struct {
	Brush     int
	Font      int
	Direction int
	SizeX     int
	SizeY     int
}

// Plot draws a string of characters on the paper at a given location
// with the colour and size of your choice.
func (n *Nimbus) Plot(opt PlotOptions, text string, x, y int) {
	// Handle default size values
	if opt.SizeX == 0 {
		opt.SizeX = 1
	}
	if opt.SizeY == 0 {
		opt.SizeY = 1
	}
	// Validate brush
	n.validateColour(opt.Brush)
	// Create a new image big enough to contain the plotted chars
	// (without scaling)
	img, _ := ebiten.NewImage(len(text)*10, 10, ebiten.FilterDefault)
	// draw chars on the image
	xpos := 0
	for _, c := range text {
		n.drawChar(img, int(c), xpos, 0, opt.Brush, opt.Font)
		xpos += 8
	}
	// Scale img and draw on paper
	ex, ey := n.convertPos(x, y, 10)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(opt.SizeX), float64(opt.SizeY))
	op.GeoM.Translate(ex, ey)
	n.paper.DrawImage(img, op)
}

// LineOptions describes optional parameters for the Line command.
type LineOptions struct {
	Brush int
}

// Line draws connected lines on the screen.  x, y values are passed in the variadic
// p parameter.
func (n *Nimbus) Line(opt LineOptions, p ...int) {
	// Validate colour
	n.validateColour(opt.Brush)
	// Use drawLine to draw connected lines
	for i := 0; i < len(p)-2; i += 2 {
		n.drawLine(p[i], p[i+1], p[i+2], p[i+3], opt.Brush)
	}
}

// AreaOptions describes optional parameters for the Area command
type AreaOptions struct {
	Brush int
}

// Area draws a filled polygon on the screen.  x, y values are passed in the variadic
// p parameter.
func (n *Nimbus) Area(opt AreaOptions, p ...int) {
	// Validate colour
	n.validateColour(opt.Brush)
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
		Color: n.convertColour(opt.Brush),
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
