package nimgobus

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Put draws an ASCII char at the cursor position
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
	n.drawChar(n.paper, c, int(ex), int(ey), n.penColour, n.charset)
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
