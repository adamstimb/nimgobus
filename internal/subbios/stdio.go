package subbios

import (
	"fmt"
	"time"

	"github.com/adamstimb/nimgobusdev/internal/queue"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Stdio implements all the console commands in kind-of old-skool C stylee.
type Stdio struct {
	s                                    *Subbios
	c                                    *console
	SuppressCtrlCInterrupt               bool // Set to true to prevent CTRL-C keyboard interrupts breaking input/output (you can still check the status though).
	SuppressCtrlShiftScrollLockInterrupt bool // Set to true to prevent CTRL-SHIFT-SCROLL LOCK keyboard interrupts  breaking input/output (you can still check the status though).
	ctrlCInterrupt                       bool
	ctrlShiftScrollLock                  bool
}

// checkKeyboardInterrupts will set the relevant interrupt flag if
// a keyboard interrupt is sent by the user.
func (sio *Stdio) checkKeyboardInterrupts() {

	if (inpututil.KeyPressDuration(ebiten.KeyControlLeft) > 1 || inpututil.KeyPressDuration(ebiten.KeyControlRight) > 1) &&
		inpututil.KeyPressDuration(ebiten.KeyC) > 1 {
		sio.ctrlCInterrupt = true
	}

	if (inpututil.KeyPressDuration(ebiten.KeyControlLeft) > 1 || inpututil.KeyPressDuration(ebiten.KeyControlRight) > 1) &&
		inpututil.KeyPressDuration(ebiten.KeyShiftLeft) > 1 || inpututil.KeyPressDuration(ebiten.KeyShiftRight) > 1 &&
		inpututil.KeyPressDuration(ebiten.KeyScrollLock) > 1 {
		sio.ctrlShiftScrollLock = true
	}

}

// GetCtrlCInterrupt checks if a CTRL+C interrupt flag is set.  Passing true will unset
// the flag after checking, otherwise the flag remains unchanged.
func (sio *Stdio) GetCtrlCInterrupt(unsetAfterChecking bool) bool {
	val := sio.ctrlCInterrupt
	if val {
		time.Sleep(250 * time.Millisecond)
	}
	if unsetAfterChecking {
		sio.ctrlCInterrupt = false
	}
	return val
}

// GetCtrlShiftScrollLockInterrupt checks if the CTRL+SHIFT+SCROLL LOCK interrupt flag is
// set. Passing true will unset the flag after checking, otherwise the flag remains unchanged.
func (sio *Stdio) GetCtrlShiftScrollLockInterrupt(unsetAfterChecking bool) bool {
	val := sio.ctrlShiftScrollLock
	if val {
		time.Sleep(250 * time.Millisecond)
	}
	if unsetAfterChecking {
		sio.ctrlShiftScrollLock = false
	}
	return val
}

// gotAnyInterrupts returns true if either CTRL+C or CTRL+SHIFT+SCROLL LOCK interrupt flags
// are set.
func (sio *Stdio) gotAnyInterrupts() bool {
	if (sio.GetCtrlCInterrupt(false) && !sio.SuppressCtrlCInterrupt) ||
		(sio.GetCtrlShiftScrollLockInterrupt(false) && !sio.SuppressCtrlShiftScrollLockInterrupt) {
		return true
	}
	return false
}

// GetCurpos returns the current cursor position and is a bit of cheat really but nvm.
func (sio *Stdio) GetCurpos() (row, col int) {
	return sio.c.cursorPositionReport()
}

// Printf imitates C's printf command but does not support formatting strings - this should be
// done by Go itself using fmt.sprintf() or similar.  The escape chars for newline `\n` and `\t` tab are supported.
// ANSI escape sequences can also be sent with this function - [a complete description of supported escape sequences and their effects is given below](#escape-sequences).
func (sio *Stdio) Printf(st string) {
	_ = sio.GetCtrlCInterrupt(true)
	_ = sio.GetCtrlShiftScrollLockInterrupt(true)
	// Load the buffer and flush it
	for _, r := range st {
		if sio.gotAnyInterrupts() {
			return
		}
		sio.c.stdoutBuffer.Enqueue(r)
	}
	sio.c.flushStdoutBuffer()
}

// Putchar prints a single rune to the screen.
func (sio *Stdio) Putchar(r rune) {
	sio.c.stdoutBuffer.Enqueue(r)
	sio.c.flushStdoutBuffer()
}

// Putchars prints a slice of runes to the screen.
func (sio *Stdio) Putchars(runes []rune) {
	for _, r := range runes {
		sio.c.stdoutBuffer.Enqueue(r)
	}
	sio.c.flushStdoutBuffer()
}

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

// Getchar gets a rune from keyboard input buffer without waiting for ENTER to be pressed.
// Stdio.KeyboardBufferFlush() should be called first though!
func (sio *Stdio) Getchar() (r rune) {
	_ = sio.GetCtrlCInterrupt(true)
	_ = sio.GetCtrlShiftScrollLockInterrupt(true)
	for {
		if sio.gotAnyInterrupts() {
			return 0
		}
		newr, ok := sio.c.stdinBuffer.Dequeue()
		if !ok {
			time.Sleep(5 * time.Millisecond)
			continue
		}
		if newr == 0 {
			time.Sleep(5 * time.Millisecond)
			continue
		} else {
			r = newr
			break
		}
	}
	return r
}

// Getch gets a rune from keyboard input buffer or returns 0 if the buffer is empty.
func (sio *Stdio) Getch() (r rune) {
	_ = sio.GetCtrlCInterrupt(true)
	_ = sio.GetCtrlShiftScrollLockInterrupt(true)

	if sio.gotAnyInterrupts() {
		time.Sleep(5 * time.Millisecond)
		return 0
	}

	if r, ok := sio.c.stdinBuffer.Dequeue(); ok {
		time.Sleep(5 * time.Millisecond)
		return r
	}

	time.Sleep(5 * time.Millisecond)
	return 0
}

// KeyboardBufferFlush removes all chars from the keyboard input buffer.  It should be called before
// Scanf() or Getchar().
func (sio *Stdio) KeyboardBufferFlush() {
	sio.c.stdinBuffer.Reset()
	sio.c.stdinBufferIndex = 0
}

// Scanf is a little bit different to the usual Scanf.  A buffer of runes (which can be empty but not nil)
// is passed to the function and echoed into the scrolling area from the current cursor position.  The user
// can then edit the buffer as they please and the changes returned via the buffer's pointer when they hit
// ENTER.
func (sio *Stdio) Scanf(textBuffer *[]rune) {
	// Setup
	_ = sio.GetCtrlCInterrupt(true)
	_ = sio.GetCtrlShiftScrollLockInterrupt(true)

	// load the buffer into a queue
	q := queue.New[rune]()
	for _, r := range *textBuffer {
		q.Enqueue(r)
	}
	bufferIndex := 0 // position of cursor in the buffer
	// Echo unedited buffer contents to screen
	for i := 0; i <= bufferIndex; i++ {
		if sio.gotAnyInterrupts() {
			return
		}
		if r, ok := q.PeekAt(i); ok {
			sio.c.stdoutBuffer.Enqueue(r)
			bufferIndex++
		}
	}
	sio.c.flushStdoutBuffer()

	// Loop
	h, w := sio.c.getScrollingAreaSize()
	lastSize := sio.c.stdinBuffer.Size()
	for {
		if sio.gotAnyInterrupts() {
			return
		}

		// Sleep if no change in buffer size
		if lastSize == sio.c.stdinBuffer.Size() {
			time.Sleep(5 * time.Millisecond)
			continue
		}
		// Get latest char in the buffer
		newR, ok := sio.c.stdinBuffer.Dequeue()
		// Ignore if something went wrong, or "empty" char r==0, or any FKEY ('\x05')
		if !ok {
			lastSize = sio.c.stdinBuffer.Size()
			time.Sleep(5 * time.Millisecond)
			continue
		}
		if newR == 0 || newR == '\x07' {
			lastSize = sio.c.stdinBuffer.Size()
			time.Sleep(5 * time.Millisecond)
			continue
		}

		// Handle key
		switch newR {
		case '\n':
			// ENTER pressed so we're finished.  Update buffer and return.
			newBuffer := []rune{}
			for i := 0; i < q.Size(); i++ {
				if r, ok := q.PeekAt(i); ok {
					newBuffer = append(newBuffer, r)
				}
			}
			*textBuffer = newBuffer
			return
		case '\x01':
			// Left arrow
			if bufferIndex == 0 {
				// no left
				// TODO: Bell?
				lastSize = sio.c.stdinBuffer.Size()
				continue
			}
			oldCursorDisplayed := sio.c.cursorDisplayed
			sio.c.cursorDisplayed = false
			oldRow := sio.c.curpos[0]
			sio.c.cursorBackward(1)
			bufferIndex--
			// Scroll down?
			if sio.c.curpos[0] == 1 && sio.c.curpos[1] == w && oldRow == 1 {
				// Scroll down and echo buffer in front of cursor
				sio.c.scrollDown(1)
				for i := bufferIndex; i >= 0; i-- {
					if r, ok := q.PeekAt(i); ok {
						x, y := sio.c.convertCurposToXY()
						// Control char
						if r < 32 && !sio.c.printControlChars {
							r = 0
						}
						sio.c.v.plonkChar(int(r), x, y, sio.c.penColour, sio.c.paperColour, sio.c.charSet, sio.c.xorWriting, sio.c.underlined)
						if sio.c.curpos[0] == 1 && sio.c.curpos[1] == 1 {
							sio.c.curpos = [2]int{1, w}
							break
						} else {
							sio.c.cursorBackward(1)
						}
					}
					if sio.gotAnyInterrupts() {
						sio.c.curpos = [2]int{1, w} // Return cursor to top-right
						return
					}
				}
				sio.c.curpos = [2]int{1, w} // Return cursor to top-right
			}
			lastSize = sio.c.stdinBuffer.Size()
			sio.c.cursorDisplayed = oldCursorDisplayed
		case '\x02':
			// Right arrow
			if bufferIndex >= q.Size() {
				// no right
				// TODO: Bell?
				lastSize = sio.c.stdinBuffer.Size()
				continue
			}
			oldCursorDisplayed := sio.c.cursorDisplayed
			sio.c.cursorDisplayed = false
			oldRow := sio.c.curpos[0]
			sio.c.cursorForward(1)
			bufferIndex++
			returnCurpos := [2]int{sio.c.curpos[0], sio.c.curpos[1]}
			// Detect carriage return.  If that happened then scroll up and we might need to echo more buffer.
			if sio.c.curpos[0] == h && sio.c.curpos[1] == 1 && oldRow == h {
				fmt.Printf("  scroll up\n")
				sio.c.scrollUp(1)
				for i := bufferIndex; i < q.Size(); i++ {
					if sio.gotAnyInterrupts() {
						sio.c.curpos = [2]int{returnCurpos[0], returnCurpos[1]} // Return cursor to original position
						return
					}
					if r, ok := q.PeekAt(i); ok {
						x, y := sio.c.convertCurposToXY()
						// Control char
						if r < 32 && !sio.c.printControlChars {
							r = 0
						}
						sio.c.v.plonkChar(int(r), x, y, sio.c.penColour, sio.c.paperColour, sio.c.charSet, sio.c.xorWriting, sio.c.underlined)
						oldRow := sio.c.curpos[0]
						sio.c.cursorForward(1)
						// Detect carriage return.  Stop echoing if that happened.
						if sio.c.curpos[0] == h && sio.c.curpos[1] == 1 && oldRow == h {
							break
						}
					}
				}
				sio.c.curpos = [2]int{returnCurpos[0], returnCurpos[1]} // Return cursor to original position
				lastSize = sio.c.stdinBuffer.Size()
			}
			sio.c.cursorDisplayed = oldCursorDisplayed
		case '\x08':
			// Backspace
			if bufferIndex == 0 {
				// nothing to delete
				// TODO: Bell?
				lastSize = sio.c.stdinBuffer.Size()
				continue
			}
			// Move cursor back and delete char at that position in the buffer
			oldCursorDisplayed := sio.c.cursorDisplayed
			sio.c.cursorDisplayed = false
			sio.c.cursorBackward(1)
			bufferIndex--
			q.RemoveAt(bufferIndex)
			// Echo remaining section of buffer
			// Echo the new char and subsequent chars but stop if it becomes necessary to scroll
			oldCurpos := [2]int{sio.c.curpos[0], sio.c.curpos[1]}
			for i := bufferIndex; i < q.Size(); i++ {
				if sio.gotAnyInterrupts() {
					sio.c.curpos = [2]int{oldCurpos[0], oldCurpos[1]} // Return cursor to original position
					return
				}
				if r, ok := q.PeekAt(i); ok {
					x, y := sio.c.convertCurposToXY()
					// Control char
					if r < 32 && !sio.c.printControlChars {
						r = 0
					}
					sio.c.v.plonkChar(int(r), x, y, sio.c.penColour, sio.c.paperColour, sio.c.charSet, sio.c.xorWriting, sio.c.underlined)
					if sio.c.curpos[0] == h && sio.c.curpos[1] == w {
						break
					} else {
						sio.c.cursorForward(1)
					}
				}
			}
			sio.c.eraseInLine(0)                              // Clear remainder of line
			sio.c.curpos = [2]int{oldCurpos[0], oldCurpos[1]} // Return cursor to original position
			lastSize = sio.c.stdinBuffer.Size()
			sio.c.cursorDisplayed = oldCursorDisplayed
		default:
			// Insert new char into buffer
			q.InsertAt(bufferIndex, newR)
			// Echo the new char and subsequent chars
			oldCursorDisplayed := sio.c.cursorDisplayed
			sio.c.cursorDisplayed = false
			oldCurpos := [2]int{sio.c.curpos[0], sio.c.curpos[1]}
			for i := bufferIndex; i < q.Size(); i++ {
				if sio.gotAnyInterrupts() {
					sio.c.curpos = [2]int{oldCurpos[0], oldCurpos[1]} // Return cursor to original position
					return
				}
				if r, ok := q.PeekAt(i); ok {
					x, y := sio.c.convertCurposToXY()
					oldRow := sio.c.curpos[0]
					// Control char
					if r < 32 && !sio.c.printControlChars {
						r = 0
					}
					sio.c.v.plonkChar(int(r), x, y, sio.c.penColour, sio.c.paperColour, sio.c.charSet, sio.c.xorWriting, sio.c.underlined)
					// If we're just echoing stuff after the cursor then break if we hit the bottom of the scrolling area
					if i > bufferIndex && sio.c.curpos[0] == h && sio.c.curpos[1] == w {
						break
					}
					// Don't cursorForward/linefeed if cursor is not at the end of the buffer but *is* in the bottom-right
					if !(sio.c.curpos[0] == h && sio.c.curpos[1] == w && i != bufferIndex) {
						sio.c.cursorForward(1)
						if sio.c.curpos[0] == oldRow && sio.c.curpos[1] == 1 {
							// scroll up required
							sio.c.scrollUp(1)
						}
					}
				}
			}
			lastSize = sio.c.stdinBuffer.Size()
			// Return cursor to original position then move it right one column
			sio.c.curpos = [2]int{oldCurpos[0], oldCurpos[1]}
			sio.c.cursorForward(1)
			sio.c.cursorDisplayed = oldCursorDisplayed
			// Advanced index
			bufferIndex++
		}
	}
}
