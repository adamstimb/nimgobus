package subbios

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// TGraphicsInput implements the one surviving function of t_graphics_input.
type TGraphicsInput struct {
	s        *Subbios
	v        *video
	muStatus sync.Mutex
	status   [3]int //x,y,b
}

// FEnquirePositionAndButtonStatus returns the x, y position of the mouse and the button status.  Unlike
// in the original Nimbus the mouse position is always being monitored and does not need to be initialized.
// There is also no support for resetting the origin or scale, since nimgobus applications run in Windowed
// environments and this would be really weird to implement and use.  Instead, the x, y positions are simply
// scaled down and reported relative to the nimgobus screen co-ordinates.
func (t *TGraphicsInput) FEnquirePositionAndButtonStatus() (x, y, b int) {
	t.muStatus.Lock()
	defer t.muStatus.Unlock()
	return t.status[0], t.status[1], t.status[2]
}

// getScale returns the scale and x, y offset of the application screen
func (t *TGraphicsInput) getScale() (scale, offsetX, offsetY float64) {
	// Get Nimbus monitor screen size
	monitorWidth, monitorHeight := t.s.Monitor.Size()

	// Get ebiten window size so we can scale the Nimbus screen up or down
	// but if (0, 0) is returned we're not running on a desktop so don't do any scaling
	windowWidth, windowHeight := ebiten.WindowSize()

	// Calculate aspect ratios of Nimbus monitor and ebiten screen
	monitorRatio := float64(monitorWidth) / float64(monitorHeight)
	windowRatio := float64(windowWidth) / float64(windowHeight)

	// If windowRatio > monitorRatio then clamp monitorHeight to windowHeight otherwise
	// clamp monitorWidth to screenWidth
	switch {
	case windowRatio > monitorRatio:
		scale = float64(windowHeight) / float64(monitorHeight)
		offsetX = (float64(windowWidth) - float64(monitorWidth)*scale) / 2
		offsetY = 0
	case windowRatio <= monitorRatio:
		scale = float64(windowWidth) / float64(monitorWidth)
		offsetX = 0
		offsetY = (float64(windowHeight) - float64(monitorHeight)*scale) / 2
	}

	return scale, offsetX, offsetY
}

// Retrieves the current mouse position and translates it to a position on the Nimbus screen
// if it's overlapping.
func (t *TGraphicsInput) update() {
	// Get absolute mouse position (we'll translate it later)
	x, y := ebiten.CursorPosition()
	// Get button status
	var b int
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		b = 1
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		b = 2
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		b = 3
	}
	// Translate x, y values
	// Scale x, y to Nimbus screen
	scale, offsetX, offsetY := t.getScale()
	x = int(float64(x) / scale)
	y = int(float64(y) / scale)
	x -= t.v.borderSize
	y -= t.v.borderSize
	x -= int(offsetX)
	y -= int(offsetY)
	videoWidth := 640
	if t.v.screenWidth == 40 {
		videoWidth = 320
	}
	if videoWidth == 640 {
		y = y / 2
	} else {
		x = x / 2
		y = y / 2
	}
	y = 250 - y // Flip vertical
	// Clamp values
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x > videoWidth {
		x = videoWidth
	}
	if y > 250 {
		y = 250
	}
	// Update tMouse status
	t.muStatus.Lock()
	defer t.muStatus.Unlock()
	t.status = [3]int{x, y, b}
}
