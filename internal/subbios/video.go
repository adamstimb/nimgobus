package subbios

import (
	"image"
	"image/color"
	"math"
	"sync"
	"time"

	"github.com/adamstimb/nimgobusdev/internal/make2darray"
	"github.com/adamstimb/nimgobusdev/internal/subbios/colour"
	"github.com/hajimehoshi/ebiten/v2"
)

// saveTableRow describes a row in the SaveTable type
type saveTableRow struct {
	X int // x position
	Y int // y position
	C int // colour
}

// SaveTable describes a save table used in sprites
type SaveTable struct {
	Rows []saveTableRow
}

// FillStyle describes the fill settings for AREA, FLOOD, CIRCLE, and SLICE
type fillStyle struct {
	Style    int // 1 for solid/dithered, 2 for hatched, 3 for hollow (edge)
	Hatching int // Hatching type if Style==2
	Colour2  int // 2nd hatching colour if Style==2
}

// Feature defines a feature to be drawn that contains a 2d image array, a screen co-ordinate,
// colour and overwrite (XOR) mode
type feature struct {
	pixels                      [][]int
	x, y                        int
	colour                      int
	xor                         bool
	fillStyle                   fillStyle
	isSprite                    bool
	isDrawSprite                bool
	isMoveSprite                bool
	isDeleteSprite              bool
	isConsoleText               bool
	saveTable                   *SaveTable
	overrideCurrentClippingArea bool
	clippingArea                int
}

// Coord defines an x, y coordinate
type coord struct {
	X int
	Y int
}

// ClippingArea defines a clipping area.
type clippingArea struct {
	MinX int
	MinY int
	MaxX int
	MaxY int
}

// video holds all the video processing malarky.
type video struct {
	monitor              *ebiten.Image       // The display image
	screenImage          *ebiten.Image       // The video memory converted into an Ebiten image
	borderImage          *ebiten.Image       // The background/border image on which the screenImage is drawn
	borderSize           int                 // The border size/thickness
	borderColour         int                 // The current border colour
	clippingAreaTable    [10]clippingArea    // A table of clipping areas
	clippingArea         int                 // The index of the currently selected clipping area
	screenWidth          int                 // The current screen width in chars (either 80 or 40)
	colourLookupTable    []colour.CltElement // The current colour lookup table
	LineStyles           [][]int             // A table of preset line styles
	lineStyleCounter     int                 // A counter for rendering line styles
	ditherPatterns       [16][4][4]int       // The current set of dither patterns
	hatchingPatterns     [6][16][16]int      // The current set of hatching patterns
	ditherLookupTables   [][][]int           // A bunch of lookup tables for rendering dither patterns
	hatchingLookupTables [][][]int           // A bunch of lookup tables for rendering hatching patterns
	polymarkers          [][][2]int          // Polymarkers are defined here
	charSet0             [256][][]int        // The default Nimbus charset
	charSet1             [256][][]int        // The alternative Nimbus charset
	muMemory             sync.Mutex          //
	memory               [250][640]int       // The video memory, a 640x250 array of integers represents the logical colours of each pixel
	muVideoMemoryOverlay sync.Mutex          //
	videoMemoryOverlay   [250][640]int       // A copy of the video memory where temporal things like cursors can be drawn
	colourFlashCounter   int                 // This counter is used to time fast and slow flashing colours
	muHoldDrawQueue      sync.Mutex          //
	holdDrawQueue        bool                // This flag will pause updateVideoMemory() if set to true
	muDrawQueue          sync.Mutex          //
	drawQueue            []feature           // A queue of features to be written to video memory
	con                  *console            // Connect the console here
	logo                 [][]int             // RM Nimbus branding
}

// advanceLineStyleCounter increments the line style counter and resets it to 0 if it exceeds 15
func (v *video) advanceLineStyleCounter() {
	v.lineStyleCounter++
	if v.lineStyleCounter > 15 {
		v.lineStyleCounter = 0
	}
}

// colourFlashTicker increments the colourFlash counter every 500 ms
func (v *video) colourFlashTicker() {
	for {
		v.colourFlashCounter++
		if v.colourFlashCounter > 4 {
			v.colourFlashCounter = 0
		}
		time.Sleep(250 * time.Millisecond)
	}
}

// purgeDrawQueue empties the draw queue
func (v *video) purgeDrawQueue() {
	v.muDrawQueue.Lock()
	v.drawQueue = []feature{}
	v.muDrawQueue.Unlock()
}

// waitForEmptyDrawQueue waits until the draw queue is empty
func (v *video) waitForEmptyDrawQueue() {
	for {
		v.muDrawQueue.Lock()
		if len(v.drawQueue) == 0 {
			v.muDrawQueue.Unlock()
			return
		}
		v.muDrawQueue.Unlock()
	}
}

// resetVideoMemory wipes the video memory (set all pixels to 0)
func (v *video) resetVideoMemory() {
	v.muMemory.Lock()
	for x := 0; x < 640; x++ {
		for y := 0; y < 250; y++ {
			v.memory[y][x] = 0
		}
	}
	v.muMemory.Unlock()
}

// resetColourLookupTable restores the default colour lookup table for the
// current screen mode.
func (v *video) resetColourLookupTable() {
	v.colourLookupTable = []colour.CltElement{}
	if v.screenWidth == 80 {
		v.colourLookupTable = append(v.colourLookupTable, colour.DefaultHighResColours...)
	} else {
		v.colourLookupTable = append(v.colourLookupTable, colour.DefaultLowResColours...)
	}
}

// resetClippingAreas (re-)initializes the clipping areas
func (v *video) resetClippingAreas() {
	maxX := 319
	if v.screenWidth == 80 {
		maxX = 639
	}
	for i := 0; i < 10; i++ {
		v.clippingAreaTable[i] = clippingArea{
			MinX: 0,
			MinY: 0,
			MaxX: maxX,
			MaxY: 249,
		}
	}
}

// initLineStyles initializes the line styles table.
func (v *video) initLineStyles() {
	v.LineStyles = [][]int{}
	v.LineStyles = append(v.LineStyles, colour.LineStyles...)
}

// initDitherPatterns initializes the dither pattern for the current screen mode.
func (v *video) initDitherPatterns() {
	for patternId := 0; patternId < 16; patternId++ {
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				if v.screenWidth == 40 {
					v.ditherPatterns[patternId][y][x] = colour.DefaultLowResDithers[patternId][y][x]
				} else {
					v.ditherPatterns[patternId][y][x] = colour.DefaultHighResDithers[patternId][y][x]
				}
			}
		}
	}
}

// initHatchingPatterns initializes the hatching patterns.
func (v *video) initHatchingPatterns() {
	for patternId := 0; patternId < 6; patternId++ {
		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				v.hatchingPatterns[patternId][y][x] = colour.DefaultHatchings[patternId][y][x]
			}
		}
	}
}

// initDitherLookupTables initializes the dither lookup tables.
// Must be called initDitherPatterns is called or after a dither pattern is changed.
func (v *video) initDitherLookupTables() {
	v.ditherLookupTables = [][][]int{}
	for d := 0; d < 16; d++ {
		// Paint dither pattern over the array
		img := make2darray.Make2dArray(640, 250, 0)
		xd := 0
		for x := 0; x < 640; x++ {
			yd := 0
			for y := 0; y < 250; y++ {
				img[y][x] = v.ditherPatterns[d][yd][xd]
				yd++
				if yd > 3 {
					yd = 0
				}
			}
			xd++
			if xd > 3 {
				xd = 0
			}
		}
		// Add array to the lookup tables
		v.ditherLookupTables = append(v.ditherLookupTables, img)
	}
}

// initHatchingLookupTables initializes the hatching lookup tables.
// Must be called whenever screen width changes.
func (v *video) initHatchingLookupTables() {
	v.hatchingLookupTables = [][][]int{}
	for d := 0; d < 6; d++ {
		// Paint dither pattern over the array
		img := make2darray.Make2dArray(640, 250, 0)
		xd := 0
		for x := 0; x < 640; x++ {
			yd := 0
			for y := 0; y < 250; y++ {
				img[y][x] = v.hatchingPatterns[d][yd][xd]
				yd++
				if yd > 15 {
					yd = 0
				}
			}
			xd++
			if xd > 15 {
				xd = 0
			}
		}
		// Add array to the lookup tables
		v.hatchingLookupTables = append(v.hatchingLookupTables, img)
	}
}

// initPolymarkers initializes the polymarkers
func (v *video) initPolymarkers() {
	v.polymarkers = [][][2]int{}
	v.polymarkers = append(v.polymarkers, colour.DefaultPolymarkers...)
}

// handleFlash is a helper function for renderScreenImage to handle flashing colours
func (v *video) handleFlash(x, y int) color.RGBA {
	// On start-up the colour lookup table may not be initialized in time so return black
	if len(v.colourLookupTable) == 0 {
		return colour.PhysicalColours[0]
	}
	cltElement := v.colourLookupTable[v.videoMemoryOverlay[y][x]] // Assumes not flashing at first
	col := colour.PhysicalColours[cltElement.FirstPhysicalColour]
	switch cltElement.FlashSpeed {
	case 1:
		// Slow flash
		if v.colourFlashCounter == 0 || v.colourFlashCounter == 1 {
			col = colour.PhysicalColours[cltElement.FirstPhysicalColour]
		} else {
			col = colour.PhysicalColours[cltElement.SecondPhysicalColour]
		}
	case 2:
		// Fast flash
		if v.colourFlashCounter == 0 || v.colourFlashCounter == 2 {
			col = colour.PhysicalColours[cltElement.FirstPhysicalColour]
		} else {
			col = colour.PhysicalColours[cltElement.SecondPhysicalColour]
		}
	}
	return col
}

// renderScreenImage converts video memory into an image.
// This function assumes the drawQueue is locked before being called!
func (v *video) renderScreenImage() {
	// Define a new image and colourises it according to video memory overlay
	img := image.NewRGBA(image.Rect(0, 0, 640, 500))
	imgX := 0
	imgY := 0
	if v.screenWidth == 40 {
		// low-res render
		for memX := 0; memX < 320; memX++ {
			for memY := 0; memY < 250; memY++ {
				col := v.handleFlash(memX, memY)
				img.Set(imgX, imgY, col)
				img.Set(imgX+1, imgY, col)
				img.Set(imgX+1, imgY+1, col)
				img.Set(imgX, imgY+1, col)
				imgY = imgY + 2
			}
			imgX = imgX + 2
			imgY = 0
		}
	} else {
		// high-res render
		for memX := 0; memX < 640; memX++ {
			for memY := 0; memY < 250; memY++ {
				col := v.handleFlash(memX, memY)
				img.Set(imgX, imgY, col)
				img.Set(imgX, imgY+1, col)
				imgY = imgY + 2
			}
			imgX = imgX + 1
			imgY = 0
		}
	}
	v.screenImage = ebiten.NewImageFromImage(img)
}

// renderMonitor draws the final monitor image to be displayed by Ebiten
func (v *video) renderMonitor() {
	// Render border (Todo: only fill when border colour changes)
	v.borderImage.Fill(colour.PhysicalColours[colour.DefaultLowResColours[v.borderColour].FirstPhysicalColour]) // Border colour does not flash and cannot be alterned by CLT
	// Draw screenImage on border image
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(v.borderSize), float64(v.borderSize))
	v.borderImage.DrawImage(v.screenImage, op)
	// Draw border image on monitor
	op = &ebiten.DrawImageOptions{}
	v.monitor.DrawImage(v.borderImage, op)
}

// Update should be called on each Ebiten Update call
func (v *video) update() {
	v.updateVideoMemory()
	v.renderMonitor()
}

// determineFeatureSize receives list of x, y coordinates and determines the size of img array
// needed to draw them as a feature.  It returns the width and height of the array and the x-
// and y-offsets (minima)
func determineFeatureSize(geometricData []coord) (width, height, offsetX, offsetY int) {
	// Seed values
	offsetX = 1000
	maxX := 0
	offsetY = 1000
	maxY := 0
	// Search
	for _, coord := range geometricData {
		if coord.X < offsetX {
			offsetX = coord.X
		}
		if coord.X > maxX {
			maxX = coord.X
		}
		if coord.Y < offsetY {
			offsetY = coord.Y
		}
		if coord.Y > maxY {
			maxY = coord.Y
		}
	}
	width = (maxX - offsetX) + 10 // deliberately padded for fill algorithm
	height = (maxY - offsetY) + 1
	return
}

// drawFeature waits until the drawQueue is unlocked then adds a feature for drawing to the drawQueue.
// It also populates the save table if a sprite is to be drawn.
func (v *video) drawFeature(f feature) {
	v.muDrawQueue.Lock()
	// handle sprite
	if f.isSprite && f.isDrawSprite {
		// Update saveTable
		clippingAreaId := v.clippingArea
		if f.overrideCurrentClippingArea {
			clippingAreaId = f.clippingArea
		}
		clip := v.clippingAreaTable[clippingAreaId]
		// Lock video memory and scan
		v.muMemory.Lock()
		f.saveTable.Rows = []saveTableRow{}
		for x := f.x; x < (f.x + len(f.pixels[0])); x++ {
			for y := 250 - f.y - len(f.pixels); y < 250-f.y; y++ {
				// Skip any coordinates outside the clipping area
				//if (x < clip.MinX || x > clip.MaxX) || (y < clip.MinY || y > clip.MaxY) {
				if (x < clip.MinX || x > clip.MaxX) || (y > 249-clip.MinY || y < 249-clip.MaxY) {
					continue
				}
				// Populate save table
				f.saveTable.Rows = append(f.saveTable.Rows, saveTableRow{X: x, Y: y, C: v.memory[y][x]})
			}
		}
		v.muMemory.Unlock()
	}
	// add to queue and unlock
	v.drawQueue = append(v.drawQueue, f)
	v.muDrawQueue.Unlock()
}

// writeFeature writes a feature directly to video memory
func (v *video) writeFeature(f feature) {
	// assumes drawQueue is locked!
	v.muMemory.Lock()
	defer v.muMemory.Unlock()
	// Get clipping area and translate y values to avoid confusion and delay
	clippingAreaId := v.clippingArea
	if f.overrideCurrentClippingArea {
		clippingAreaId = f.clippingArea
	}
	clip := v.clippingAreaTable[clippingAreaId]
	// But override clipping area if console text and use full screen
	if f.isConsoleText {
		maxX := 640
		if v.screenWidth == 40 {
			maxX = 320
		}
		clip = clippingArea{0, 0, maxX, 250}
	}
	// redraw saveTable data and then update it if it's a sprite being moved
	if f.isSprite && f.isMoveSprite {
		// draw saveTable
		for _, r := range f.saveTable.Rows {
			v.memory[r.Y][r.X] = r.C
		}
		// update saveTable
		f.saveTable.Rows = []saveTableRow{}
		for x := f.x; x < (f.x + len(f.pixels[0])); x++ {
			for y := 250 - f.y - len(f.pixels); y < 250-f.y; y++ {
				// Skip any coordinates outside the clipping area
				//if (x < clip.MinX || x > clip.MaxX) || (y < clip.MinY || y > clip.MaxY) {
				if (x < clip.MinX || x > clip.MaxX) || (y > 249-clip.MinY || y < 249-clip.MaxY) {
					continue
				}
				// Populate save table
				f.saveTable.Rows = append(f.saveTable.Rows, saveTableRow{X: x, Y: y, C: v.memory[y][x]})
			}
		}
	}
	// Write the feature
	fX := 0
	for x := f.x; x < (f.x + len(f.pixels[0])); x++ {
		fY := 0
		for y := 250 - f.y - len(f.pixels); y < 250-f.y; y++ {
			// Skip any coordinates outside the clipping area
			//if (x < clip.MinX || x > clip.MaxX) || (y < clip.MinY || y > clip.MaxY) {
			if (x < clip.MinX || x > clip.MaxX) || (y > 249-clip.MinY || y < 249-clip.MaxY) {
				fY++
				continue
			}
			// Draw sprite
			// colour >= 0 represents a b+w feature so colourise it with specified colour
			// otherwise use the colour given by the pixel (but if feature pixel has negative
			// colour leave it as transparent)
			if f.colour >= 0 {
				if f.pixels[fY][fX] == 1 {
					if f.xor {
						// XOR mode
						v.memory[y][x] = v.memory[y][x] ^ f.colour
					} else {
						v.memory[y][x] = f.colour
					}
				}
			} else {
				if f.pixels[fY][fX] >= 0 {
					if f.xor {
						// XOR mode
						v.memory[y][x] = v.memory[y][x] ^ f.pixels[fY][fX]
					} else {
						v.memory[y][x] = f.pixels[fY][fX]
					}
				}
			}
			fY++
		}
		fX++
	}
}

// writeFeatureToOverlay writes a feature directly to video memory overlay
func (v *video) writeFeatureToOverlay(f feature) {
	v.muVideoMemoryOverlay.Lock()
	defer v.muVideoMemoryOverlay.Unlock()

	// Write the feature
	fX := 0
	for x := f.x; x < (f.x + len(f.pixels[0])); x++ {
		fY := 0
		for y := 250 - f.y - len(f.pixels); y < 250-f.y; y++ {
			// Skip any coordinates outside the screen area
			if (x < 0 || x > 640) || (y < 0 || y > 250) {
				fY++
				continue
			}
			// Draw sprite
			// colour >= 0 represents a b+w feature so colourise it with specified colour
			// otherwise use the colour given by the pixel (but if feature pixel has negative
			// colour leave it as transparent)
			if f.colour >= 0 {
				if f.pixels[fY][fX] == 1 {
					if f.xor {
						// XOR mode
						v.videoMemoryOverlay[y][x] = v.videoMemoryOverlay[y][x] ^ f.colour
					} else {
						v.videoMemoryOverlay[y][x] = f.colour
					}
				}
			} else {
				if f.pixels[fY][fX] >= 0 {
					if f.xor {
						// XOR mode
						v.videoMemoryOverlay[y][x] = v.videoMemoryOverlay[y][x] ^ f.pixels[fY][fX]
					} else {
						v.videoMemoryOverlay[y][x] = f.pixels[fY][fX]
					}
				}
			}
			fY++
		}
		fX++
	}
}

// rotateFeature rotates a feature 90 degrees counterclockwise r times
func (v *video) rotateFeature(f feature, r int) feature {
	for i := 0; i < r; i++ {
		f = v.rotateFeature90(f)
	}
	return f
}

// rotateFeature90 rotates a feature 90 degrees counterclockwise
func (v *video) rotateFeature90(f feature) feature {
	img := f.pixels
	imgWidth := len(img[0])
	imgHeight := len(img)
	newWidth := imgHeight
	newHeight := imgWidth
	newImg := make2darray.Make2dArray(newWidth, newHeight, -1)
	for x1 := 0; x1 < imgWidth; x1++ {
		for y1 := 0; y1 < imgHeight; y1++ {
			x2 := y1
			y2 := (newHeight - 1) - x1
			newImg[y2][x2] = img[y1][x1]
		}
	}
	return feature{pixels: newImg, x: f.x, y: f.y, colour: f.colour, xor: f.xor}
}

// resizeFeature does a nearest-neighbour resize of a feature
func (v *video) resizeFeature(f feature, newWidth, newHeight int) feature {
	img := f.pixels
	newImg := make2darray.Make2dArray(newWidth, newHeight, -1)
	imgWidth := len(img[0])
	imgHeight := len(img)
	xScale := float64(imgWidth) / float64(newWidth)
	yScale := float64(imgHeight) / float64(newHeight)
	for y2 := 0; y2 < newHeight; y2++ {
		for x2 := 0; x2 < newWidth; x2++ {
			x1 := int(math.Floor((float64(x2) + 0.5) * xScale))
			y1 := int(math.Floor((float64(y2) + 0.5) * yScale))
			newImg[y2][x2] = img[y1][x1]
		}
	}
	return feature{pixels: newImg, x: f.x, y: f.y, colour: f.colour, xor: f.xor}
}

// updateVideoMemory writes all sprites in the drawQueue to video memory
func (v *video) updateVideoMemory() {
	// Skip if holdDrawQueue is true
	v.muHoldDrawQueue.Lock()
	if v.holdDrawQueue {
		v.muHoldDrawQueue.Unlock()
		return
	}
	v.muHoldDrawQueue.Unlock()
	// Write all the features in the draw queue
	v.muDrawQueue.Lock()
	for _, f := range v.drawQueue {
		v.writeFeature(f)
	}
	// update video overlay
	for y := 0; y < 250; y++ {
		v.videoMemoryOverlay[y] = v.memory[y]
	}
	// draw cursor on overlay if enabled
	if v.con.cursorDisplayed {
		v.drawCursor()
	}
	// flush drawQueue and render screen image
	v.drawQueue = []feature{}
	v.muDrawQueue.Unlock()
	v.renderScreenImage()
}

// GetXY returns the value in the video memory at x, y
func (v *video) GetXY(x, y int) int {
	v.muMemory.Lock()
	c := v.memory[249-y][x]
	v.muMemory.Unlock()
	return c
}

// SetXY sets the pixel value in video memory at x, y
func (v *video) SetXY(x, y, c int) {
	v.muMemory.Lock()
	v.memory[249-y][x] = c
	v.muMemory.Unlock()
}

// makeConsoleCharImg returns the image of a console char
func (v *video) makeConsoleCharImg(c, fg, bg, charset int, underline bool) [][]int {

	var charPixels [][]int
	switch charset {
	case 0:
		charPixels = v.charSet0[c]
	case 1:
		charPixels = v.charSet1[c]
	}
	// add paper and pen colour to charPixels
	newCharPixels := make2darray.Make2dArray(8, 10, bg)
	for x := 0; x < 8; x++ {
		for y := 0; y < 10; y++ {
			if charPixels[y][x] == 1 {
				newCharPixels[y][x] = fg
			}
			// underline?
			if underline && y == 8 {
				newCharPixels[y][x] = fg
			}
		}
	}
	return newCharPixels
}

// plonkChar is called by the console to render an ASCII char on the screen
func (v *video) plonkChar(c, x, y, fg, bg, charset int, xor, underline bool) {
	img := v.makeConsoleCharImg(c, fg, bg, charset, underline)
	v.drawFeature(feature{pixels: img, x: x, y: y, colour: -1, xor: xor, isConsoleText: true})
}

// drawCursor draws the cursor at the current curpos
func (v *video) drawCursor() {
	x, y := v.con.convertCurposToXY()
	// don't draw if outside video memory bounds
	if x > 639 || y < 0 {
		return
	}

	// otherwise draw cursor
	img := v.makeConsoleCharImg(v.con.cursorChar, v.con.penColour, v.con.paperColour, v.con.cursorCharSet, v.con.cursorUnderlined)

	// Steady cursor
	if !v.con.cursorFlashing {
		v.writeFeatureToOverlay(feature{pixels: img, x: x, y: y, colour: -1, xor: true})
	}

	// Flashing cursor
	if v.con.cursorFlashing && (v.colourFlashCounter < 2) {
		v.writeFeatureToOverlay(feature{pixels: img, x: x, y: y, colour: -1, xor: true})
	}
}
