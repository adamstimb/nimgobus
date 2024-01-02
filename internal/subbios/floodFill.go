package subbios

// Flood fill algorithm adapted from https://stackoverflow.com/questions/2783204/flood-fill-using-a-stack
func (v *video) floodFillDo(maxX int, hits [250][640]bool, srcColor, fillStyle, fillStyleIndex, fillColour1, fillColour2, transparency, boundarySpecification, colourOfBoundary, x, y int) bool {
	if (y < 0) || (x < 0) || (y > 249) || (x > maxX) {
		return false
	}
	if hits[y][x] {
		return false
	}
	newCol := v.GetXY(x, y)
	if boundarySpecification == 1 {
		if newCol == colourOfBoundary {
			return false
		}
	} else {
		if newCol != srcColor {
			return false
		}
	}
	// valid, paint it
	v.SetXY(x, y, v.fillColour(x, y, fillStyle, fillStyleIndex, fillColour1, fillColour2, transparency))
	return true
}

func (v *video) floodFill(fillStyle, fillStyleIndex, fillColour1, fillColour2, transparency, boundarySpecification, colourOfBoundary, x, y int) {
	maxX := 639
	if v.screenWidth == 40 {
		maxX = 319
	}
	srcColor := v.GetXY(x, y)
	hits := [250][640]bool{}
	stack := []coord{{x, y}}
	for {
		p := stack[len(stack)-1] // pop the stack
		stack = stack[:len(stack)-1]
		result := v.floodFillDo(maxX, hits, srcColor, fillStyle, fillStyleIndex, fillColour1, fillColour2, transparency, boundarySpecification, colourOfBoundary, p.X, p.Y)
		if result {
			hits[p.Y][p.X] = true
			stack = append(stack, coord{p.X, p.Y + 1})
			stack = append(stack, coord{p.X, p.Y - 1})
			stack = append(stack, coord{p.X + 1, p.Y})
			stack = append(stack, coord{p.X - 1, p.Y})
		}
		if len(stack) == 0 {
			// Nothing left to fill
			break
		}
	}
}
