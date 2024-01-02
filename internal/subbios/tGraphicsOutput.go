package subbios

import (
	"math"

	"github.com/adamstimb/nimgobusdev/internal/make2darray"
	"github.com/adamstimb/nimgobusdev/internal/subbios/colour"
	"github.com/adamstimb/nimgobusdev/internal/subbios/errorcode"
)

// TGraphicsOutput has all the t_graphics_output functions attached to it.
type TGraphicsOutput struct {
	s  *Subbios
	v  *video
	On bool
}

// FGraphicsOutputColdStart initializes the graphics system.  If the graphics system
// is not initialized with either cold start or warm start then the other graphics
// functions will return the error ENotInitialized (except the colour lookup table and
// border colour functions).
// If the error EAlreadyOn is returned it may be useful to call function 4 (reinitialize)
// to restore the screen to a known starting point.
// WARNING - all of the internals of the graphics system are reinitialized whenever the
// the screen width is changed, even when graphics are switched off.
// On the Nimbus it was very important to switch off the graphics system before exiting
// (FGraphicsOutputOff) before exiting - not so important here, but it's implemented
// nonetheless.
func (t *TGraphicsOutput) FGraphicsOutputColdStart() {
	t.s.FunctionError = errorcode.EOk
	// Handle already on
	if t.On {
		t.s.FunctionError = errorcode.EAlreadyOn
		return
	}
	t.v.initLineStyles()
	t.v.purgeDrawQueue()
	t.v.resetColourLookupTable()
	t.v.initDitherPatterns()
	t.v.initHatchingPatterns()
	t.v.initDitherLookupTables()
	t.v.initHatchingLookupTables()
	t.v.initPolymarkers()
	t.v.resetClippingAreas()
	t.v.resetVideoMemory()
	// Set the flag and we're done
	t.On = true
}

// FGraphicsOutputWarmStart is the same as FGraphicsOutputWarmStart except that it
// does not re-initialize anything, e.g. user defined dither patterns.
func (t *TGraphicsOutput) FGraphicsOutputWarmStart() {
	t.s.FunctionError = errorcode.EOk
	// Handle already on
	if t.On {
		t.s.FunctionError = errorcode.EAlreadyOn
		return
	}
	t.v.purgeDrawQueue()
	t.v.resetVideoMemory()
	// Set the flag and we're done
	t.On = true
}

// FGraphicsOutputOff closes down the graphics system, which means that any future
// access to TGraphicsOutput results in the error ENotInitialized (exceptions are
// cold start, warm start, colour lookup table functions and border colour functions).
func (t *TGraphicsOutput) FGraphicsOutputOff() {
	t.s.FunctionError = errorcode.EOk
	// Set the flag and we're done
	t.On = false
}

// FReinitGraphicsOutput re-initializes the graphics system if it's currently on:
// - the colour lookup table is initialized.
// - hatching and dither patterns are set to their default values.
// - clipping areas 1-9 are set to the full screen.
// This function is executed during a screen width change and at start-up(?). If the
// graphics system is currently switch off, ENotInitialized is returned, and the
// function has not effect.
func (t *TGraphicsOutput) FReinitGraphicsOutput() {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	t.v.initLineStyles()
	t.v.initDitherPatterns()
	t.v.initHatchingPatterns()
	t.v.initDitherLookupTables()
	t.v.initHatchingLookupTables()
	t.v.initPolymarkers()
	t.v.purgeDrawQueue()
	t.v.resetColourLookupTable()
	t.v.resetVideoMemory()
}

// FSetBorderColour sets the colour of the screen border.  The border cannot flash
// and its colour is independent of screen resolution. For example, if the entry
// parameter is 6, the border colour will be set to brown.  The error EInvalidParameter
// is given if the border colour is not in the range 0-15.
func (t *TGraphicsOutput) FSetBorderColour(c int) {
	t.s.FunctionError = errorcode.EOk
	// Validate c
	if c < 0 || c > 15 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	t.v.borderColour = c
}

// FGetBorder returns the current border colour.
func (t *TGraphicsOutput) FGetBorderColour() int {
	t.s.FunctionError = errorcode.EOk
	return t.v.borderColour
}

// FSetCltElement sets an element in the colour lookup table.
// elementNumber can be in the range 0-15 for low-resolution mode, or 0-3 in high-resolution mode.
// firstPhysicalColour and secondPhysicalColour can be in the range 0-15.
// flashSpeed can be in the range 0-2: 0 - no flash, 1 - slow flash, 2 - fast flash.
// If any parameters are out of range, the error EInvalidParameter is given.
func (t *TGraphicsOutput) FSetCltElement(elementNumber, firstPhysicalColour, flashSpeed, secondPhysicalColour int) {
	t.s.FunctionError = errorcode.EOk
	// Validate params
	maxElement := 15
	if t.v.screenWidth == 80 {
		maxElement = 3
	}
	if elementNumber < 0 || elementNumber > maxElement {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if firstPhysicalColour < 0 || firstPhysicalColour > 15 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if flashSpeed < 0 || flashSpeed > 2 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if secondPhysicalColour < 0 || secondPhysicalColour > 15 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	// Set element
	t.v.colourLookupTable[elementNumber] = colour.CltElement{
		FirstPhysicalColour:  firstPhysicalColour,
		FlashSpeed:           flashSpeed,
		SecondPhysicalColour: secondPhysicalColour,
	}
}

// FGetCltElement returns an element from the colour lookup table.
// elementNumber can be in the range 0-15 for low-resolution mode, or 0-3 in high-resolution mode.
// If any parameters are out of range, the error EInvalidParameter is given.
// See FSetCltElement for a description of the return values.
func (t *TGraphicsOutput) FGetCltElement(elementNumber int) (firstPhysicalColour, flashSpeed, secondPhysicalColour int) {
	t.s.FunctionError = errorcode.EOk
	// Validate params
	maxElement := 15
	if t.v.screenWidth == 80 {
		maxElement = 3
	}
	if elementNumber < 0 || elementNumber > maxElement {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	// Get element, extract values and return
	element := t.v.colourLookupTable[elementNumber]
	firstPhysicalColour = element.FirstPhysicalColour
	flashSpeed = element.FlashSpeed
	secondPhysicalColour = element.SecondPhysicalColour
	return
}

// FSetOutputClippingAreaLimits defines a clipping area.
// id is the clipping area id (range: 1 - 9 (0 is not user-definable))
// minX, minY, maxX, maxY define the rectangular shape of the clipping area.
// If any coordinates are outside the screen we get EInvalidParameter.
// If the minima are greater than the maxima we get EInvalidParameter.
// If id is out-of-range we get EInvalidParameter.
func (t *TGraphicsOutput) FSetOutputClippingAreaLimits(id, minX, minY, maxX, maxY int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Validate params
	limitX := 319
	if t.v.screenWidth == 80 {
		limitX = 639
	}
	if minX >= maxX || minY >= maxY {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if minY < 0 || minX < 0 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if maxX > limitX || maxY > 249 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if id < 1 || id > 9 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	// Set
	t.v.waitForEmptyDrawQueue()
	t.v.clippingAreaTable[id] = clippingArea{
		MinX: minX,
		MinY: minY,
		MaxX: maxX,
		MaxY: maxY,
	}
}

// FGetOutputClippingAreaLimits returns the limits of a clipping area.
// id is the clipping area id (range: 1 - 9 (0 is not user-definable))
// If id is out-of-range we get EInvalidParameter.
// Returns the min x, min y, max x and max y co-ordinates of the clipping area boundary.
func (t *TGraphicsOutput) FGetOutputClippingAreaLimits(id int) (minX, minY, maxX, maxY int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Validate params
	if id < 1 || id > 9 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	// Get
	return t.v.clippingAreaTable[id].MinX, t.v.clippingAreaTable[id].MinY, t.v.clippingAreaTable[id].MaxX, t.v.clippingAreaTable[id].MaxY
}

// FSetCurrentOutputClippingArea sets the current clipping area.
// id can be in the range 0-9.  0 is always the full screen.
// If id is outside 0-9 the error EInvalidParameter is raised.
func (t *TGraphicsOutput) FSetCurrentOutputClippingArea(id int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Validate
	if id < 0 || id > 9 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	// Set
	t.v.waitForEmptyDrawQueue()
	t.v.clippingArea = id
}

// FGetCurrentOutputClippingArea gets the current clipping area.
func (t *TGraphicsOutput) FGetCurrentOutputClippingArea() (id int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Get
	return t.v.clippingArea
}

// FGetCltContents returns the entire colour lookup table.
func FGetCltContents(t *TGraphicsOutput) (clt []int) {
	t.s.FunctionError = errorcode.EOk
	for i := 0; i < len(t.v.colourLookupTable); i = i + 3 {
		clt = append(clt, t.v.colourLookupTable[i].FirstPhysicalColour)
		clt = append(clt, t.v.colourLookupTable[i].FlashSpeed)
		clt = append(clt, t.v.colourLookupTable[i].SecondPhysicalColour)
	}
	return clt
}

// FSetNewClt replaces the entire colour lookup table with new values.
// newClt is a list representing the new CLT.  If in low-res mode there must be 3*16 values in
// the list, otherwise 3*4.  Each value must be in the range 0-15.
func (t *TGraphicsOutput) FSetNewClt(newClt []int) {
	t.s.FunctionError = errorcode.EOk
	// Validate
	if (t.v.screenWidth == 80 && len(newClt) != 3*4) || (t.v.screenWidth == 40 && len(newClt) != 3*16) {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	for _, c := range newClt {
		if c < 0 || c > 15 {
			t.s.FunctionError = errorcode.EInvalidParameter
			return
		}
	}
	// Set
	c := 0
	for i := 0; i < len(newClt); i = i + 3 {
		t.v.colourLookupTable[c].FirstPhysicalColour = newClt[i]
		t.v.colourLookupTable[c].FlashSpeed = newClt[i+1]
		t.v.colourLookupTable[c].SecondPhysicalColour = newClt[i+3]
	}
}

// FPolyLine draws a series of conmnected lines.
// lineStyle selects the line style to draw with and is in the range 0-6.  0 is a solid line using a dither colour,
// 1 is a solid line using the first logical colour, 2 is dashed, 3 is dotted, 4 is dash-dotted, 5 is irregular dashed
// and 6 is user-defined.
// lineStyleIndex[0] selects the dither pattern to use if lineStyle=0, or defines a custom line style if lineStyle=6
// using the entire array, e.g. []int{0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1}.
// transparency is zero if the second logical colour is to be visible, or 1 if it should be transparent.
// The shape is drawn in XOR mode if 256 is passed as the first logical colour for dither styles or adding 256 to first
// logical colour for all other styles.
// TODO: use [][2]int{} for geometricData
func (t *TGraphicsOutput) FPolyLine(lineStyle int, lineStyleIndex []int, firstLogicalColour, secondLogicalColour, transparency int, geometricData []int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Validate
	if lineStyle < 0 || lineStyle > 6 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if len(geometricData) < 2 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if len(geometricData)%2 != 0 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if lineStyle == 0 {
		// Must have dith pattern selected in lineStyleIndex[0]
		if len(lineStyleIndex) != 1 {
			t.s.FunctionError = errorcode.EInvalidParameter
			return
		}
		if lineStyleIndex[0] < 0 || lineStyleIndex[0] > 15 {
			t.s.FunctionError = errorcode.EInvalidParameter
			return
		}
	}
	if lineStyle == 6 {
		// Must have dith pattern defined in lineStyleIndex[0:15]
		if len(lineStyleIndex) != 16 {
			t.s.FunctionError = errorcode.EInvalidParameter
			return
		}
		for _, i := range lineStyleIndex {
			maxC := 15
			if t.v.screenWidth == 80 {
				maxC = 3
			}
			if i < 0 || i > maxC {
				t.s.FunctionError = errorcode.EInvalidParameter
				return
			}
		}
	}
	if transparency < 0 || transparency > 1 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	maxCol := 15
	if t.v.screenWidth == 80 {
		maxCol = 3
	}
	if firstLogicalColour < 0 || (firstLogicalColour > maxCol && firstLogicalColour < 256) || (firstLogicalColour > 256+maxCol) {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if secondLogicalColour < 0 || secondLogicalColour > maxCol {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	// XOR mode?
	xor := false
	if lineStyle == 0 && firstLogicalColour == 256 {
		xor = true
	}
	if lineStyle > 0 && firstLogicalColour > 255 {
		xor = true
		firstLogicalColour = firstLogicalColour - 256
	}
	// Convert into slice of Coord
	coords := []coord{}
	for i := 0; i < len(geometricData)-1; i += 2 {
		coords = append(coords, coord{X: geometricData[i], Y: geometricData[i+1]})
	}
	// Draw lines on array
	width, height, offsetX, offsetY := determineFeatureSize(coords)
	img := make2darray.Make2dArray(width, height, -1)
	for i := 0; i < len(coords)-1; i++ {
		fromXY := coords[i]
		toXY := coords[i]
		if i != len(coords)-1 {
			toXY = coords[i+1]
		}
		img = t.v.drawLine(img, fromXY.X-offsetX, fromXY.Y-offsetY, toXY.X-offsetX, toXY.Y-offsetY, lineStyle, lineStyleIndex, firstLogicalColour, secondLogicalColour, transparency)
	}
	// Load array into a feature and draw it
	t.v.drawFeature(feature{pixels: img, x: offsetX, y: offsetY, colour: -1, xor: xor}) // colour=-1 because we have a colour feature with transparent (-1) background
}

// FFillArea fills the area described by a set of vertices given in the geometricData parameter.
// fillStyle: 0 - hollow, 1 - solid using fillColour1, 2 - solid using dither, 3 - solid using hatched pattern.
// fillStyleIndex (fillStyle=0, 1) - ignored.
// fillStyleIndex (fillStyle=2) - dither pattern to use.
// fillStyleIndex (fillStyle=3) - hatching pattern to use.
// fillColour1 - the first fill colour
// fillColour2 - the second fill colour (ignored unless hatching selected).
// transparency - fillColour2 is transparent if set to 1 in hatching mode.
// XOR mode can be selected by passing 256 in fillColour1 when using a dither pattern,
// or by adding 256 to fillColour1's value if otherwise.
// geometricData... TODO: use [][2]int{} for geometricData
func (t *TGraphicsOutput) FFillArea(fillStyle, fillStyleIndex, fillColour1, fillColour2, transparency int, geometricData []int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Validate
	if len(geometricData) < 2 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if len(geometricData)%2 != 0 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if transparency < 0 || transparency > 1 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	maxCol := 15
	if t.v.screenWidth == 80 {
		maxCol = 3
	}
	if fillColour1 < 0 || (fillColour1 > maxCol && fillColour1 < 256) || (fillColour1 > 256+maxCol) {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if fillColour2 < 0 || fillColour2 > maxCol {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	// XOR mode ?
	xor := false
	if fillColour1 > 255 {
		xor = true
		fillColour1 = fillColour1 - 256
	}
	if fillStyle == 2 && fillColour1 == 256 {
		xor = true
	}
	// Convert into slice of Coord
	coords := []coord{}
	for i := 0; i < len(geometricData)-1; i += 2 {
		coords = append(coords, coord{X: geometricData[i], Y: geometricData[i+1]})
	}
	// Prepare an array to draw on
	width, height, offsetX, offsetY := determineFeatureSize(coords)
	img := make2darray.Make2dArray(width, height, -1)
	// Close the shape if necessary
	if coords[0].X != coords[len(coords)-1].X || coords[0].Y != coords[len(coords)-1].Y {
		coords = append(coords, coord{coords[0].X, coords[0].Y})
	}
	// Draw outline
	for i := 0; i < len(coords)-1; i++ {
		img = t.v.drawLine(img, coords[i].X-offsetX, coords[i].Y-offsetY, coords[i+1].X-offsetX, coords[i+1].Y-offsetY, 1, []int{0}, fillColour1, fillColour2, transparency)
	}
	// If hollow shape then we're already done
	if fillStyle == 0 {
		t.v.drawFeature(feature{pixels: img, x: offsetX, y: offsetY, colour: -1, xor: xor})
		return
	}
	// Otherwise let's cheat and use draw2d to draw a filled polygon
	img = t.v.d2dFilledPolygon(geometricData, width, height, offsetX, offsetY, fillStyle, fillStyleIndex, fillColour1, fillColour2, transparency)

	t.v.drawFeature(feature{pixels: img, x: offsetX, y: offsetY, colour: -1, xor: xor})
}

// FFloodFillArea fills the screen out from the point (x, y) to a boundary.
// fillStyle: 0 or 1 - solid using fillColour1, 2 - solid using dither, 3 - solid using hatched pattern.
// fillStyleIndex (fillStyle=0, 1) - ignored.
// fillStyleIndex (fillStyle=2) - dither pattern to use.
// fillStyleIndex (fillStyle=3) - hatching pattern to use.
// fillColour1 - the first fill colour.
// fillColour2 - the second fill colour (ignored unless hatching selected).
// transparency - of fillColour2.
// boundarySpecification: 0 - any colour different from that of the seed position is a boundary, 1 - the boundary is any pixel of colour colourOfBoundary.
// colourOfBoundary - the colour of the boundary (ignored if boundarySpecification == 0).
// x - x-coordinate of the seed position.
// y - y-coordinate of the seed position.
func (t *TGraphicsOutput) FFloodFillArea(fillStyle, fillStyleIndex, fillColour1, fillColour2, transparency, boundarySpecification, colourOfBoundary, x, y int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Validate
	if transparency < 0 || transparency > 1 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	maxCol := 15
	if t.v.screenWidth == 80 {
		maxCol = 3
	}
	if fillColour1 < 0 || fillColour1 > maxCol {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if fillColour2 < 0 || fillColour2 > maxCol {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if fillStyleIndex == 0 {
		fillStyleIndex = 1
	}
	if boundarySpecification < 0 || boundarySpecification > 1 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if colourOfBoundary < 0 || colourOfBoundary > maxCol {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	t.v.waitForEmptyDrawQueue()
	t.v.floodFill(fillStyle, fillStyleIndex, fillColour1, fillColour2, transparency, boundarySpecification, colourOfBoundary, x, y)
}

// FPolymarker draws a series of unconnected markers on the screen.
// markerStyle: (1-6) the style of marker, 6 being a custom marker defined in markerShape.
// markerSizeX: (1-50) the x magnification of the marker.
// markerSizeY: (1-50) the y magnification of the marker.
// logicalColour: the logical colour of the marker.
// markerShape: the shape of a custom marker if markerStyle==6, otherwise ignored.
func (t *TGraphicsOutput) FPolymarker(markerStyle, markerSizeX, markerSizeY, logicalColour int, markerShape [][2]int, geometricData [][2]int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Validate
	if markerStyle < 1 || markerStyle > 6 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if (markerSizeX < 1 || markerSizeX > 50) || (markerSizeY < 1 || markerSizeY > 50) {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	maxCol := 15
	if t.v.screenWidth == 80 {
		maxCol = 3
	}
	if logicalColour < 0 || (logicalColour > maxCol && logicalColour < 256) || (logicalColour > 256+maxCol) {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	// XOR mode ?
	xor := false
	if logicalColour > 255 {
		xor = true
		logicalColour = logicalColour - 256
	}
	// Handle custom marker
	if markerStyle != 6 {
		markerShape = t.v.polymarkers[markerStyle-1]
	}
	// Create an array to draw the marker on (find out how big it needs to be first)
	maxX := 0
	maxY := 0
	for _, p := range markerShape {
		if (p[0] == -128 && p[1] == 0) || (p[0] == -128 && p[1] == -128) {
			// Ignore pen-up / end
			continue
		}
		absX := int(math.Abs(float64(p[0])))
		absY := int(math.Abs(float64(p[1])))
		if absX > maxX {
			maxX = absX
		}
		if absY > maxY {
			maxY = absY
		}
	}
	sizeX := (markerSizeX * 2 * maxX) + 1
	sizeY := (markerSizeY * 2 * maxY) + 1
	img := make2darray.Make2dArray(sizeX, sizeY, -1)
	// Draw polymarker in img
	penUp := true
	x, y := 0, 0
	offsetX, offsetY := sizeX/2, sizeY/2
	for _, p := range markerShape {
		// Set penUp if requested
		if p[0] == -128 && p[1] == 0 {
			penUp = true
			continue
		}
		// Break if end of marker
		if p[0] == -128 && p[1] == -128 {
			break
		}
		// Move-to, if pen is up
		if penUp {
			x, y = markerSizeX*p[0], markerSizeY*p[1]
			penUp = false
			continue
		}
		// Otherwise draw-to
		img = t.v.drawLine(img, x+offsetX, y+offsetY, (markerSizeX*p[0])+offsetX, (markerSizeY*p[1])+offsetY, 1, []int{0}, logicalColour, 0, 0)
		x, y = markerSizeX*p[0], markerSizeY*p[1]
	}
	// Render imgs on screen
	for _, p := range geometricData {
		t.v.drawFeature(feature{pixels: img, x: p[0] - offsetX, y: p[1] - offsetY, colour: -1, xor: xor})
	}
}

// FPlotCharacterString plots a character string on the screen.
// orientation: (0-3) the direction to plot, rotated 90 degress anticlockwise each time.
// yMagnification: (1-50) the amount to enlarge vertically.
// xMagnification: (1-50) the amount to enlarge horizontally.
// logicalColour: the colour to plot in, add 256 to plot in XOR mode.
// font: 0 - use the standard character set, 1 - use the alternative character set (custom charset not yet implemented)
// chars: the string of chars to plot.
// x, y: The co-ordinates to plot at.
func (t *TGraphicsOutput) FPlotCharacterString(orientation, yMagnification, xMagnification, logicalColour, font int, chars string, x, y int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Validate
	if (yMagnification < 1 || yMagnification > 50) || (xMagnification < 1 || xMagnification > 50) {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if font < 0 || font > 1 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if orientation < 0 || orientation > 3 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	maxCol := 15
	if t.v.screenWidth == 80 {
		maxCol = 3
	}
	if logicalColour < 0 || (logicalColour > maxCol && logicalColour < 256) || (logicalColour > 256+maxCol) {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	// XOR mode ?
	xor := false
	if logicalColour > 255 {
		xor = true
		logicalColour = logicalColour - 256
	}
	// Plot chars and applying scaling/direction
	imgWidth := len(chars) * 8
	imgHeight := 10
	img := make2darray.Make2dArray(imgWidth, imgHeight, -1)
	// Select charset and draw chars on image
	xOffset := 0
	for _, c := range chars {
		// draw char on image
		var charPixels [][]int
		switch font {
		case 0:
			charPixels = t.v.charSet0[c]
		case 1:
			charPixels = t.v.charSet1[c]
		}
		for x := 0; x < 8; x++ {
			for y := 0; y < 10; y++ {
				img[y][x+xOffset] = charPixels[y][x]
			}
		}
		xOffset += 8
	}
	resizedFeature := t.v.resizeFeature(feature{pixels: img, x: x, y: y, colour: logicalColour, xor: xor}, imgWidth*xMagnification, imgHeight*yMagnification)
	rotatedFeature := t.v.rotateFeature(feature{pixels: resizedFeature.pixels, x: x, y: y, colour: logicalColour, xor: xor}, orientation)
	// Correct x, y for orientation
	if orientation == 2 {
		rotatedFeature.x = rotatedFeature.x - ((len(chars) - 1) * 8 * xMagnification)
	}
	if orientation == 3 {
		rotatedFeature.y = rotatedFeature.y - ((len(chars) - 1) * 8 * yMagnification)
	}
	t.v.drawFeature(rotatedFeature)
}

// Sprite represents (as much as is practical) Nimbus sprite data.
type Sprite struct {
	HighResolution bool      // Set to true if the sprite is intended for high-resolution mode.
	Hotspot        [2]int    // The x, y vector from the bottom-left of the sprite to the hotspot.
	Poses          [][][]int // The sprite images stored as 2D arrays in individual poses.  Up to 2 poses allowed in low-resolution mode, up to 4 in high-resolution mode.
}

// FDrawSprite draws a sprite on the screen and stores the overwritten data in an saveTable array.
func (t *TGraphicsOutput) FDrawSprite(s Sprite, saveTable *SaveTable, x, y, pose int, xor bool, clippingAreaId int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Validate
	if pose >= len(s.Poses) {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	maxPoses := 2
	if t.v.screenWidth == 40 {
		maxPoses = 4
	}
	if len(s.Poses) > maxPoses {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if clippingAreaId < 0 || clippingAreaId >= 10 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	// Do it
	t.v.drawFeature(feature{
		pixels:                      s.Poses[pose],
		x:                           x - s.Hotspot[0],
		y:                           y - s.Hotspot[1],
		colour:                      -1,
		xor:                         xor,
		isSprite:                    true,
		isDrawSprite:                true,
		saveTable:                   saveTable,
		overrideCurrentClippingArea: true,
		clippingArea:                clippingAreaId,
	})
}

// FMoveSprite moves an existing sprite on the screen and stores the overwritten data in an saveTable array.
// Because of the way saveTable has been implemented, this command does not require the "old x, y" as originally
// documented.
func (t *TGraphicsOutput) FMoveSprite(s Sprite, saveTable *SaveTable, x, y, pose int, xor bool, clippingAreaId int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Validate
	if pose >= len(s.Poses) {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	maxPoses := 2
	if t.v.screenWidth == 40 {
		maxPoses = 4
	}
	if len(s.Poses) > maxPoses {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if clippingAreaId < 0 || clippingAreaId >= 10 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	t.v.drawFeature(feature{
		pixels:                      s.Poses[pose],
		x:                           x - s.Hotspot[0],
		y:                           y - s.Hotspot[1],
		colour:                      -1,
		xor:                         xor,
		isSprite:                    true,
		isMoveSprite:                true,
		saveTable:                   saveTable,
		overrideCurrentClippingArea: true,
		clippingArea:                clippingAreaId,
	})
}

// FEraseSprite erases a sprite associated with a saveTable.  Unlike in the original implementation, it is only necessary
// to pass the sprite and saveTable as arguments.
func (t *TGraphicsOutput) FEraseSprite(saveTable *SaveTable) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	pixels := [][]int{{}} // only need an empty pose because we're not going to draw anything
	t.v.drawFeature(feature{
		pixels:       pixels,
		x:            0,
		y:            0,
		colour:       -1,
		isSprite:     true,
		isMoveSprite: true,
		saveTable:    saveTable,
	})
}

// FPlonkLogo draws the RM Nimbus logo on the screen, starting with the
// bottom-left at x, y.
func (t *TGraphicsOutput) FPlonkLogo(x, y int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	f := feature{pixels: t.v.logo, x: x, y: y, colour: -1, xor: false}
	t.v.drawFeature(f)
}

// FPieSlice draws pie slices and filled circles.  xCentre and yCentre represent the position of the
// centre of the circle or pie slice.  Radius is the radius along the vertical axis (this is because
// in high-resolution mode the circle must be stretched across the horizontal axis so it is still a
// circle).  theta1 and theta2 are the starting and stopping angles of the slice, respectively.  If
// theta1 == theta2 then a complete circle will be drawn.  theta1 and theta2 are measured in thousandths
// of a radian, with 0 or 6283 being vertically up (don't ask me, this is how it was originally) and
// 3142 being vertically down.  colour is the colour of the circle (outline and fill colour).
func (t *TGraphicsOutput) FPieSlice(xCentre, yCentre, radius, theta1, theta2, colour int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	img := make2darray.Make2dArray((radius*2)+1, (radius*2)+1, -1)
	t.v.drawCircle(img, radius, radius, radius, theta1, theta2, colour)
	f := feature{pixels: img, x: xCentre - radius, y: yCentre - radius, colour: -1, xor: false}
	if t.v.screenWidth == 80 {
		rescaledF := t.v.resizeFeature(f, len(img[0])*2, len(img))
		rescaledF.x = xCentre - (2 * radius)
		t.v.drawFeature(rescaledF) // Draw a rescaled feature if in high-resolution mode
	} else {
		t.v.drawFeature(f) // Otherwise draw as-is
	}
}

// FArcOfEllipse is not implemented
func (t *TGraphicsOutput) FArcOfEllipse() {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionError = errorcode.EFuncNotImplemented
}

// FSetDitherPattern sets one of the user-definable dither patterns.  ditherId is the number
// of the dither pattern to be defined (8-15). ditherPattern is the new dither pattern defined
// in an array, e.g. [4][4]int{{1,2,1,2},{2,1,2,1},{1,2,1,2},{2,1,2,1}}.  Note that the maximum
// colour value in the dither is set according to screen mode, i.e. 3 in high-resolution mode
// and 15 in low-resolution mode.
func (t *TGraphicsOutput) FSetDitherPattern(ditherId int, ditherPattern [4][4]int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	maxColour := 3
	if t.v.screenWidth == 40 {
		maxColour = 15
	}
	if ditherId < 8 || ditherId > 15 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			colour := ditherPattern[y][x]
			if colour > maxColour || colour < 0 {
				t.s.FunctionError = errorcode.EInvalidParameter
				return
			} else {
				t.v.ditherPatterns[ditherId][y][x] = colour
			}
		}
	}
	t.v.initDitherLookupTables()
}

// FGetDitherPattern returns a dither pattern.
func (t *TGraphicsOutput) FGetDitherPattern(ditherId int) (ditherPattern [4][4]int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	if ditherId < 8 || ditherId > 15 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	ditherPattern = t.v.ditherPatterns[ditherId] // This *should* be copy by reference
	return ditherPattern
}

// FSetHatchingPattern sets one of the user-definable dither patterns.  hatching is the number
// of the hacching pattern to be defined (0-5). hatchingPattern is the new hatching pattern defined
// in an array, e.g. [16][16]int{{1,1,..,0,0},..,{1,1,..,0,0}}.  Note that hatching patterns do
// not store implicit colour information and so the allowed values are either 0 or 1.
func (t *TGraphicsOutput) FSetHatchingPattern(hatchingId int, hatchingPattern [16][16]int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	if hatchingId < 0 || hatchingId > 5 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			colour := hatchingPattern[y][x]
			if colour > 1 || colour < 0 {
				t.s.FunctionError = errorcode.EInvalidParameter
				return
			} else {
				t.v.hatchingPatterns[hatchingId][y][x] = colour
			}
		}
	}
	t.v.initHatchingLookupTables()
}

// FGetHatchingPattern returns a hatching pattern.
func (t *TGraphicsOutput) FGetHatchingPattern(hatchingId int) (hatchingPattern [16][16]int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	if hatchingId < 0 || hatchingId > 5 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	hatchingPattern = t.v.hatchingPatterns[hatchingId] // This *should* be copy by reference
	return hatchingPattern
}

// FGetDisplayLine is not implemented as it's redundant in nimgobus.
func (t *TGraphicsOutput) FGetDisplayLine() {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionError = errorcode.EFuncNotImplemented
}

// FReadPixel returns the logical colour of a pixel on the screen.
func (t *TGraphicsOutput) FReadPixel(x, y int) (colour int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	maxX := 320
	if t.v.screenWidth == 80 {
		maxX = 640
	}
	if x < 0 || x > maxX || y < 0 || y > 250 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	t.v.waitForEmptyDrawQueue()
	t.v.muDrawQueue.Lock()
	t.v.muMemory.Lock()
	colour = t.v.memory[250-y][x]
	t.v.muMemory.Unlock()
	t.v.muDrawQueue.Unlock()
	return colour
}

// FReadToLimit is not implemented because it looks like an almighty faff to me.
func (t *TGraphicsOutput) FReadToLimit() {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionError = errorcode.EFuncNotImplemented
}

// FReadAreaWord is not implemented because it's redundant in nimgobus.
func (t *TGraphicsOutput) FReadAreaWord() {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionError = errorcode.EFuncNotImplemented
}

// FWriteAreaWord is not implemented because it's redundant in nimgobus.
func (t *TGraphicsOutput) FWriteAreaWord() {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionError = errorcode.EFuncNotImplemented
}

// FCopyAreaWord is not implemented because it's redundant in nimgobus.
func (t *TGraphicsOutput) FCopyAreaWord() {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionError = errorcode.EFuncNotImplemented
}

// FSwapAreaWord is not implemented because it's redundant in nimgobus.
func (t *TGraphicsOutput) FSwapAreaWord() {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionError = errorcode.EFuncNotImplemented
}

// FIsine is not implemented because it's redundant in nimgobus.
func (t *TGraphicsOutput) FIsine() {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionError = errorcode.EFuncNotImplemented
}

// FIcos is not implemented because it's redundant in nimgobus.
func (t *TGraphicsOutput) FIcos() {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionError = errorcode.EFuncNotImplemented
}

// FReadAreaPixel returns the screen memory from a given area of the screen.
func (t *TGraphicsOutput) FReadAreaPixel(xMin, yMin, xMax, yMax int) (img [][]int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Validate params
	xLimit := 320
	if t.v.screenWidth == 80 {
		xLimit = 640
	}
	if xMin < 0 || xMin > xLimit || yMin < 0 || yMin > 249 || xMax < 0 || xMax > xLimit || yMax < 0 || yMax > 249 {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	if xMin > xMax || yMin > yMax {
		t.s.FunctionError = errorcode.EInvalidParameter
		return
	}
	// Capture screen memory
	t.v.waitForEmptyDrawQueue()
	t.v.muDrawQueue.Lock()
	t.v.muMemory.Lock()
	for y := yMax; y >= yMin; y-- {
		column := []int{}
		for x := xMin; x <= xMax; x++ {
			column = append(column, t.v.memory[249-y][x])
		}
		img = append(img, column)
	}
	t.v.muMemory.Unlock()
	t.v.muDrawQueue.Unlock()
	return img
}

// FWriteAreaPixel writes an image array captured by FReadAreaPixel to the
// screen.
func (t *TGraphicsOutput) FWriteAreaPixel(img [][]int, xMin, yMin int, xor bool, ignoreLogicalColour int) {
	t.s.FunctionError = errorcode.EOk
	// Handle not on
	if !t.On {
		t.s.FunctionError = errorcode.ENotInitialized
		return
	}
	// Handle ignoreLogicalColour
	maxColour := 3
	if t.v.screenWidth == 80 {
		maxColour = 15
	}
	if ignoreLogicalColour >= 0 && ignoreLogicalColour <= maxColour {
		for y := 0; y < len(img); y++ {
			for x := 0; x < len(img[0]); x++ {
				if img[y][x] == ignoreLogicalColour {
					img[y][x] = -1
				}
			}
		}
	}
	// Draw it
	t.v.waitForEmptyDrawQueue()
	f := feature{pixels: img, x: xMin, y: yMin, xor: xor, colour: -1}
	t.v.drawFeature(f)
}

// FCopyAreaPixel was not implemented because it's redundant in nimgobus - just use FReadAreaPixel followed
// by FWriteAreaPixel instead.
func (t *TGraphicsOutput) FCopyAreaPixel() {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionError = errorcode.EFuncNotImplemented
}
