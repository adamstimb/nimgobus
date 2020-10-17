package nimgobus

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/adamstimb/nimgobus/images"
	"github.com/hajimehoshi/ebiten"
)

// colRow defines a column, row position
type colRow struct {
	col int
	row int
}

// textBox defines the bounding box of a scrollable text box
type textBox struct {
	col1 int
	row1 int
	col2 int
	row2 int
}

// Nimbus acts as a container for all the components of the Nimbus monitor.  You
// only need to call the Init() method after declaring a new Nimbus.
type Nimbus struct {
	Monitor               *ebiten.Image
	paper                 *ebiten.Image
	basicColours          []color.RGBA
	borderSize            int
	borderColour          int
	paperColour           int
	penColour             int
	charset               int
	cursorChar            int
	defaultHighResPalette []int
	defaultLowResPalette  []int
	palette               []int
	logoImage             *ebiten.Image
	textBoxes             [10]textBox
	imageBlocks           [16]*ebiten.Image
	selectedTextBox       int
	cursorPosition        colRow
	cursorMode            int
	cursorCharset         int
	cursorFlash           bool
	charImages0           [256]*ebiten.Image
	charImages1           [256]*ebiten.Image
}

// Init initializes a new Nimbus.  You must call this method after declaring a
// new Nimbus variable.
func (n *Nimbus) Init() {
	// in case any randomonia is required we can run a seed on startup
	rand.Seed(time.Now().UnixNano())

	// Load Nimbus logo image and both charsets
	n.loadLogoImage()
	n.loadCharsetImages(0)
	n.loadCharsetImages(1)
	// Set init values of everything else
	n.borderSize = 50
	n.Monitor = ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2))
	n.paper = ebiten.NewImage(640, 250)
	n.basicColours = basicColours
	n.defaultHighResPalette = defaultHighResPalette
	n.defaultLowResPalette = defaultLowResPalette
	n.palette = defaultHighResPalette
	n.borderColour = 0
	n.paperColour = 0
	n.penColour = 3
	n.charset = 0
	n.cursorMode = -1
	n.cursorChar = 95
	n.cursorCharset = 0
	n.cursorPosition = colRow{1, 1}
	n.cursorFlash = false
	n.selectedTextBox = 0
	// Initialize with mode 80 textboxes
	for i := 0; i < 10; i++ {
		n.textBoxes[i] = textBox{1, 1, 25, 80}
	}
	// Start flashCursor
	go n.flashCursor()
}

// flashCursor flips the cursorFlash flag every half second
func (n *Nimbus) flashCursor() {
	for {
		time.Sleep(500 * time.Millisecond)
		if n.cursorMode == 0 {
			// Flashing cursor
			n.cursorFlash = !n.cursorFlash
		}
		if n.cursorMode < 0 {
			// Invisible cursor
			n.cursorFlash = false
		}
		if n.cursorMode > 1 {
			// Visible cursor but not flashing
			n.cursorFlash = true
		}
	}
}

// Update redraws the Nimbus monitor image
func (n *Nimbus) Update() {

	// Copy paper so we can apply overlays (e.g. cursor)
	paperCopy := ebiten.NewImageFromImage(n.paper)

	// Apply overlays
	// Cursor
	if n.cursorFlash {
		curX, curY := n.convertColRow(n.cursorPosition)
		n.drawChar(paperCopy, n.cursorChar, int(curX), int(curY), n.penColour, n.cursorCharset)
	}

	// calculate y scale for paper and apply scaling
	paperX, paperY := paperCopy.Size()
	scaleX := 640.0 / float64(paperX)
	scaleY := 500.0 / float64(paperY)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	scaledPaper := ebiten.NewImage(640, 500)
	scaledPaper.DrawImage(paperCopy, op)

	// Add the border around the paper
	withBorder := ebiten.NewImage(640+(n.borderSize*2), 500+(n.borderSize*2))
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
	n.logoImage = ebiten.NewImageFromImage(img)
}

// loadCharsetImages loads the charset images
func (n *Nimbus) loadCharsetImages(charset int) {
	var img image.Image
	var err error
	if charset == 0 {
		img, _, err = image.Decode(bytes.NewReader(images.CharSetZeroImage))
	} else {
		img, _, err = image.Decode(bytes.NewReader(images.CharSetOneImage))
	}
	if err != nil {
		log.Fatal(err)
	}
	img2 := ebiten.NewImageFromImage(img)
	for i := 0; i <= 255; i++ {
		if charset == 0 {
			n.charImages0[i] = n.charImageSelecta(img2, i)
		} else {
			n.charImages1[i] = n.charImageSelecta(img2, i)
		}
	}
}

// convertPos receives Nimbus-style screen coords and returns then as ebiten-style
func (n *Nimbus) convertPos(x, y, imageHeight int) (ex, ey float64) {
	_, paperHeight := n.paper.Size()
	return float64(x), float64(paperHeight) - float64(y) - float64(imageHeight)
}

// convertColRow receives a Nimbus-style column, row position and returns an
// ebiten-style graphical coordinate
func (n *Nimbus) convertColRow(cr colRow) (ex, ey float64) {
	ex = (float64(cr.col) * 8) - 8
	ey = (float64(cr.row) * 10) - 10
	return ex, ey
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
func (n *Nimbus) charImageSelecta(img *ebiten.Image, c int) *ebiten.Image {

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
	return img.SubImage(image.Rect(x1, y1, x2, y2)).(*ebiten.Image)
}

// drawChar draws a character at a specific location on an image
func (n *Nimbus) drawChar(image *ebiten.Image, c, x, y, colour, charset int) {
	// Draw char on image and apply colour
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	rgba := n.convertColour(colour)
	r := float64(rgba.R) / 0xff
	g := float64(rgba.G) / 0xff
	b := float64(rgba.B) / 0xff
	op.ColorM.Translate(r, g, b, 0)
	if charset == 0 {
		image.DrawImage(n.charImages0[c], op)
	} else {
		image.DrawImage(n.charImages1[c], op)
	}
}
