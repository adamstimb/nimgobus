package nimgobus

import (
	"fmt"
	"runtime"

	"github.com/elastic/go-sysinfo"
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

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// drawBackground draws the background of the Welcome screen
func drawBackground(n *Nimbus) {

	// Collect system info
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	firmwareVersion := runtime.Version() // Use Go version instead
	firmwareVersion = firmwareVersion[2:]
	if len(firmwareVersion) > 8 {
		firmwareVersion = firmwareVersion[:8]
	}
	host, err := sysinfo.Host()
	if err != nil {
		panic("Could detect system information")
	}
	firmwareVersion = fmt.Sprintf("Firmware version: %s", firmwareVersion)
	serialNumber := "Serial number:  21/06809" // In honour of whichever physical machine donate its ROM to MAME
	memInfo, err := host.Memory()
	mainMemSize := fmt.Sprintf("main    memory size %7d Mbytes", bToMb(memInfo.Available))
	virtualMemSize := fmt.Sprintf("virtual memory size %7d Mbytes", bToMb(memInfo.VirtualTotal))
	totalMemSize := fmt.Sprintf("total   memory size %7d Mbytes", bToMb(memInfo.Available+memInfo.VirtualTotal))

	// Red frame, light blue paper, Nimbus logo in a red frame
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

	// Welcome
	n.Plot(plotOpts, "Welcome", 238, 145)
	plotOpts.Brush = 2
	n.Plot(plotOpts, "Welcome", 236, 147)

	// Firmware version and serial number
	areaOpts.Brush = 2
	n.Area(areaOpts, 393, 4, 632, 4, 632, 30, 393, 30, 393, 4)
	areaOpts.Brush = 3
	n.Area(areaOpts, 395, 5, 629, 5, 629, 29, 395, 29, 395, 5)
	plotOpts.Brush = 0
	plotOpts.SizeX = 1
	plotOpts.SizeY = 1
	n.Plot(plotOpts, firmwareVersion, 400, 17)
	n.Plot(plotOpts, serialNumber, 400, 7)

	// Memory
	plotOpts.Brush = 0
	n.Plot(plotOpts, mainMemSize, 15, 25)
	n.Plot(plotOpts, virtualMemSize, 15, 15)
	n.Plot(plotOpts, totalMemSize, 15, 5)
}
