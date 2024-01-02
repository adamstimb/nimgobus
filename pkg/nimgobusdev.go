// nimgobus contains the Nimbus struct from which all the pretend 16-bit magic stems!
package nimgobus

import (
	"github.com/adamstimb/nimgobus/internal/subbios"
	"github.com/hajimehoshi/ebiten/v2"
)

// The Subbios commands and Monitor image are accessed here.
type Nimbus struct {
	Subbios subbios.Subbios
	Monitor *ebiten.Image
}

// Initializes nimgobus.
func (n *Nimbus) Init() {
	n.Subbios.Init()
	n.Monitor = n.Subbios.Monitor
}

// Update needs to be called on each ebiten update call.
func (n *Nimbus) Update() {
	if !n.Subbios.TGraphicsOutput.On {
		return
	}
	n.Subbios.Update()
}
