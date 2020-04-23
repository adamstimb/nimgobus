package nimgobus

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/adamstimb/nimgobus/images"
	"github.com/hajimehoshi/ebiten"
)

// Nimbus acts as a container for all the components of the monitor
type Nimbus struct {
	Monitor               *ebiten.Image
	paper                 *ebiten.Image
	basicColours          []color.RGBA
	borderSize            int
	borderColour          int
	paperColour           int
	defaultHighResPalette []int
	defaultLowResPalette  []int
	palette               []int
	logoImage             *ebiten.Image
	charsetZeroImage      *ebiten.Image
}

// Init initializes a new Nimbus
func (n *Nimbus) Init() {
	n.loadLogoImage()
	n.loadCharsetImages()
	n.borderSize = 50
	n.Monitor, _ = ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2), ebiten.FilterDefault)
	n.paper, _ = ebiten.NewImage(640, 250, ebiten.FilterDefault)
	n.basicColours = basicColours
	n.defaultHighResPalette = defaultHighResPalette
	n.defaultLowResPalette = defaultLowResPalette
	n.palette = defaultHighResPalette
	n.borderColour = 0
	n.paperColour = 0
}

// Update draws the monitor image
func (n *Nimbus) Update() {

	// calculate y scale for paper and apply scaling
	_, paperY := n.paper.Size()
	scaleY := 500.0 / float64(paperY)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, scaleY)
	scaledPaper, _ := ebiten.NewImage(640, 500, ebiten.FilterDefault)
	scaledPaper.DrawImage(n.paper, op)

	// Add the border around the paper
	withBorder, _ := ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2), ebiten.FilterDefault)
	withBorder.Fill(n.convertColour(n.borderColour)) // Apply border colour
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(n.borderSize), float64(n.borderSize))
	withBorder.DrawImage(scaledPaper, op)

	// Draw paper with border on monitor
	op = &ebiten.DrawImageOptions{}
	n.Monitor.DrawImage(withBorder, op)
}

// loadLogoImage loads the Nimbus logo image
func (n *Nimbus) loadLogoImage() {
	img, _, err := image.Decode(bytes.NewReader(images.NimbusLogoImage))
	if err != nil {
		log.Fatal(err)
	}
	n.logoImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

// loadCharsetImages loads the charset images
func (n *Nimbus) loadCharsetImages() {
	img, _, err := image.Decode(bytes.NewReader(images.CharSetZeroImage))
	if err != nil {
		log.Fatal(err)
	}
	n.charsetZeroImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

// convertPos receives Nimbus-style screen coords and returns then as ebiten-style
func (n *Nimbus) convertPos(x, y, imageHeight int) (ex, ey float64) {
	_, paperHeight := n.paper.Size()
	return float64(x), float64(paperHeight) - float64(y) - float64(imageHeight)
}

// validateColour checks that a Nimbus colour index is valid for the current
// screen mode.  If validation fails then a panic is issued.
func (n *Nimbus) validateColour(c int) {
	// Negative values and anything beyond the pallete range is not allowed
	if c < 0 {
		panic("Negative values are not allowed for colours")
	}
	if c > len(n.palette)-1 {
		panic("Colour is out of range for this screen mode")
	}
}

// convertColour receives a Nimbus colour index and returns the RGBA
func (n *Nimbus) convertColour(c int) color.RGBA {
	return n.basicColours[n.palette[c]]
}

// charImageSelecta returns the subimage pointer of a char from the charset
// image for any given ASCII code.  If control char is received, a blank char
// is returned instead.
func (n *Nimbus) charImageSelecta(c int) *ebiten.Image {
	// select blank char 127 if control char
	if c < 33 {
		c = 127
	}

	// Calculate row and column position of the char on the charset image
	mapNumber := c - 32 // codes < 33 are not on the map
	row := int(math.Ceil(float64(mapNumber) / float64(30)))
	column := mapNumber - (30 * (row - 1))

	// Calculate corners of rectangle around the char
	x1 := (column - 1) * 10
	x2 := x1 + 10
	y1 := (row - 1) * 10
	y2 := y1 + 10

	// Return pointer to sub image
	return n.charsetZeroImage.SubImage(image.Rect(x1, y1, x2, y2)).(*ebiten.Image)
}

// drawChar draws a character at a specific location on an image
func (n *Nimbus) drawChar(image *ebiten.Image, c, x, y, colour int) {
	// Convert position
	ex, ey := n.convertPos(x, y, 10)
	// Draw char on paper
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ex, ey)
	// Apply colour
	rgba := n.convertColour(colour)
	r := float64(rgba.R)
	g := float64(rgba.G)
	b := float64(rgba.B)
	op.ColorM.Translate(r, g, b, 0)
	image.DrawImage(n.charImageSelecta(c), op)
}
