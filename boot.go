package nimgobus

import (
	"runtime"
)

// Boot simulates the RM Nimbus "Welcome" boot screen
func (n *Nimbus) Boot() {
	drawBackground(n)
	plotOpts := PlotOptions{
		Font:  1,
		Brush: 3,
	}
	n.Plot(plotOpts, "Please supply an operating system", 188, 100)
}

// drawBackground draws the background of the Welcome screen
func drawBackground(n *Nimbus) {
	// Collect system info
	firmwareVersion := runtime.Version()
	firmwareVersion = firmwareVersion[2:]
	if len(firmwareVersion) > 8 {
		firmwareVersion = firmwareVersion[:8]
	}
	//serialNumber := "21/06809"
	n.SetMode(80)
	n.SetColour(0, 0)
	n.SetColour(1, 9)
	n.SetPaper(1)
	n.SetBorder(1)
	n.SetCursor(-1)
	n.Cls()
	areaOpts := AreaOptions{
		Brush: 2,
	}
	n.Area(areaOpts, 0, 0, 639, 0, 639, 249, 0, 249, 0, 0)
	areaOpts.Brush = 1
	n.Area(areaOpts, 3, 2, 636, 2, 636, 247, 3, 247, 3, 2)
	xl := 10
	yl := 212
	n.PlonkLogo(xl, yl)
	lineOpts := LineOptions{
		Brush: 2,
	}
	n.Line(lineOpts, xl, yl, xl+304, yl, xl+304, yl+32, xl, yl+32, xl, yl)
	plotOpts := PlotOptions{
		SizeX: 3, SizeY: 3, Font: 1,
	}
	n.Plot(plotOpts, "Welcome", 238, 145)
	plotOpts.Brush = 2
	n.Plot(plotOpts, "Welcome", 236, 147)
	// test system info
	n.Print(firmwareVersion)
}
