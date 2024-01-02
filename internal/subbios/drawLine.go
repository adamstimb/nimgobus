package subbios

// lineColour determines the colour - if any - of a particular pixel on a line
func (v *video) lineColour(x, y, lineStyle int, lineStyleIndex []int, firstLogicalColour, secondLogicalColour, transparency int) (colour int) {
	if transparency == 1 {
		secondLogicalColour = -1 // second logical colour is transparent
	}
	if lineStyle == 6 {
		v.LineStyles[6] = lineStyleIndex // setup user-defined line style
	}
	// Dither?
	if lineStyle == 0 {
		return v.ditherLookupTables[lineStyleIndex[0]][y][x]
	}
	// Or determine if it's first or second logical colour
	if v.LineStyles[lineStyle][v.lineStyleCounter] == 1 {
		return firstLogicalColour
	} else {
		return secondLogicalColour
	}
}

// drawLine implements Bresenham's line algorithm to draw a line on a 2d array
// adapted from https://github.com/StephaneBunel/bresenham/blob/master/drawline.go
// Todo: fix stepping effect.
func (v *video) drawLine(img [][]int, x1, y1, x2, y2, lineStyle int, lineStyleIndex []int, firstLogicalColour, secondLogicalColour, transparency int) [][]int {
	// Todo: fix cool trippy effect.

	imgHeight := len(img) - 1
	var dx, dy, e, slope int
	v.lineStyleCounter = 0

	// Because drawing p1 -> p2 is equivalent to draw p2 -> p1,
	// I sort points in x-axis order to handle only half of possible cases.
	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dx, dy = x2-x1, y2-y1
	// Because point is x-axis ordered, dx cannot be negative
	if dy < 0 {
		dy = -dy
	}

	switch {

	// Is line a point ?
	case x1 == x2 && y1 == y2:
		img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)

	// Is line an horizontal ?
	case y1 == y2:
		for ; dx != 0; dx-- {
			img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)
			v.advanceLineStyleCounter()
			x1++
		}
		img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)

	// Is line a vertical ?
	case x1 == x2:
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for ; dy != 0; dy-- {
			img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)
			v.advanceLineStyleCounter()
			y1++
		}
		img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)

	// Is line a diagonal ?
	case dx == dy:
		if y1 < y2 {
			for ; dx != 0; dx-- {
				img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)
				v.advanceLineStyleCounter()
				x1++
				y1++
			}
		} else {
			for ; dx != 0; dx-- {
				img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)
				v.advanceLineStyleCounter()
				x1++
				y1--
			}
		}
		img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)

	// wider than high ?
	case dx > dy:
		if y1 < y2 {
			// BresenhamDxXRYD(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)
				v.advanceLineStyleCounter()
				x1++
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}
			}
		} else {
			// BresenhamDxXRYU(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)
				v.advanceLineStyleCounter()
				x1++
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}
			}
		}
		img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)

	// higher than wide.
	default:
		if y1 < y2 {
			// BresenhamDyXRYD(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)
				v.advanceLineStyleCounter()
				y1++
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		} else {
			// BresenhamDyXRYU(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)
				v.advanceLineStyleCounter()
				y1--
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		}
		img[imgHeight-y1][x1] = v.lineColour(x1, imgHeight-y1, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)
	}
	return img
}
