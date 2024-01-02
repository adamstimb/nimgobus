package subbios

import (
	"image"
	"image/color"

	"github.com/adamstimb/nimgobus/internal/make2darray"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

// fillColour determines the colour of a particular pixel in a fill area
func (v *video) fillColour(x, y, fillStyle int, fillStyleIndex int, firstLogicalColour, secondLogicalColour, transparency int) (colour int) {
	if transparency == 1 {
		secondLogicalColour = -1 // second logical colour is transparent
	}
	if fillStyle == 1 {
		// Solid filled
		return firstLogicalColour
	}
	// Dither?
	if fillStyle == 2 {
		return v.ditherLookupTables[fillStyleIndex][y][x]
	}
	// Hatching?
	if fillStyle == 3 {
		if v.hatchingLookupTables[fillStyleIndex][y][x] == 1 {
			return firstLogicalColour
		} else {
			if transparency == 0 {
				return secondLogicalColour
			} else {
				return -1
			}
		}
	}
	return firstLogicalColour
}

// d2dFilledPolygon uses the draw2d package to draw a filled polygon
func (v *video) d2dFilledPolygon(geometricData []int, width, height, offsetX, offsetY, fillStyle, fillStyleIndex, fillColour1, fillColour2, transparency int) (filledImg [][]int) {
	// draw the polygon in an image using gg - we'll convert it back to a simple array later
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	dc := draw2dimg.NewGraphicContext(img)
	p := color.RGBA{1, 1, 1, 255}
	dc.SetFillColor(p)
	dc.SetFillRule(draw2d.FillRuleWinding)
	dc.SetStrokeColor(p)
	dc.SetLineWidth(1)
	dc.MoveTo(float64(geometricData[0]-offsetX), float64(geometricData[1]-offsetY))
	for i := 2; i < len(geometricData); i += 2 {
		dc.LineTo(float64(geometricData[i]-offsetX), float64(geometricData[i+1]-offsetY))
	}
	dc.Close()
	dc.FillStroke()

	// convert image to array
	filledImg = make2darray.Make2dArray(width, height, -1)
	filledCol := color.RGBA{1, 1, 1, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			col := img.At(x, y)
			if rgb, ok := col.(color.RGBA); ok {
				if rgb == filledCol {
					filledImg[(height-1)-y][x] = v.fillColour(x, y, fillStyle, fillStyleIndex, fillColour1, fillColour2, transparency)
				}
			}
		}
	}

	// draw true outline and correct for overfilling and underfilling
	// Convert into slice of Coord
	coords := []coord{}
	for i := 0; i < len(geometricData)-1; i += 2 {
		coords = append(coords, coord{X: geometricData[i], Y: geometricData[i+1]})
	}
	for i := 0; i < len(coords)-1; i++ {
		filledImg = v.outline(filledImg, coords[i].X-offsetX, coords[i].Y-offsetY, coords[i+1].X-offsetX, coords[i+1].Y-offsetY, 1, fillStyleIndex, -255, fillColour2, transparency)
	}
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			f := v.fillColour(x, y, fillStyle, fillStyleIndex, fillColour1, fillColour2, transparency)
			// left underfill?
			if filledImg[y][x-1] == -255 && filledImg[y][x] == -1 && filledImg[y][x+1] == f {
				filledImg[y][x] = f
			}
			// right underfill?
			if filledImg[y][x-1] == f && filledImg[y][x] == -1 && filledImg[y][x+1] == -255 {
				filledImg[y][x] = f
			}
			// up underfill?
			if filledImg[y-1][x] == -255 && filledImg[y][x] == -1 && filledImg[y+1][x] == f {
				filledImg[y][x] = f
			}
			// down underfill?
			if filledImg[y-1][x] == f && filledImg[y][x] == -1 && filledImg[y+1][x] == -255 {
				filledImg[y][x] = f
			}
			// left overfill?
			if filledImg[y][x-1] == -1 && filledImg[y][x] == f && filledImg[y][x+1] == -255 {
				filledImg[y][x] = -1
			}
			// right overfill?
			if filledImg[y][x-1] == -255 && filledImg[y][x] == f && filledImg[y][x+1] == -1 {
				filledImg[y][x] = -1
			}
			// up overfill?
			if filledImg[y-1][x] == -1 && filledImg[y][x] == f && filledImg[y+1][x] == -255 {
				filledImg[y][x] = -1
			}
			// down overfill?
			if filledImg[y-1][x] == -255 && filledImg[y][x] == f && filledImg[y+1][x] == -1 {
				filledImg[y][x] = -1
			}
		}
	}
	for i := 0; i < len(coords)-1; i++ {
		filledImg = v.outline(filledImg, coords[i].X-offsetX, coords[i].Y-offsetY, coords[i+1].X-offsetX, coords[i+1].Y-offsetY, fillStyle, fillStyleIndex, fillColour1, fillColour2, transparency)
	}
	return filledImg
}
