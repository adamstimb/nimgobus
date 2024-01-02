package subbios

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/elastic/go-sysinfo"
	"github.com/shirou/gopsutil/host"
)

// randDelay delays for a random number of milliseconds within limits
func randDelay(min, max int) {
	delay := time.Duration(rand.Intn(max-min)+min) * time.Millisecond
	time.Sleep(delay)
}

// convert bytes to Mb
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// welcomeScreen draws the Welcome screen
func welcomeScreen(s *Subbios) {

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
		panic("Could not detect system information")
	}
	firmwareVersion = fmt.Sprintf("Firmware version: %s", firmwareVersion)
	serialNumber := "Serial number:  21/06809" // In honour of whichever physical machine donate its ROM to MAME
	memInfo, err := host.Memory()
	mainMemSize := fmt.Sprintf("main    memory size %7d Mbytes", bToMb(memInfo.Available))
	virtualMemSize := fmt.Sprintf("virtual memory size %7d Mbytes", bToMb(memInfo.VirtualTotal))
	totalMemSize := fmt.Sprintf("total   memory size %7d Mbytes", bToMb(memInfo.Available+memInfo.VirtualTotal))

	// Red frame, light blue paper, Nimbus logo in a red frame
	s.Stdio.Printf("\x1b[2h") // Mode 80
	s.Stdio.Printf("\x1b[~F") // Hide cursor
	s.TGraphicsOutput.FGraphicsOutputColdStart()
	s.TGraphicsOutput.FSetBorderColour(9)
	s.TGraphicsOutput.FSetCltElement(0, 9, 0, 0)
	s.TGraphicsOutput.FSetCltElement(1, 0, 0, 0)
	s.TGraphicsOutput.FSetCltElement(2, 10, 0, 0)
	s.TGraphicsOutput.FSetCltElement(3, 15, 0, 0)
	s.Stdio.Printf("\x1b[50;33m") // Light blue paper and white ink
	s.Stdio.Printf("\x1b[2J")     // CLS
	s.Stdio.Printf("\x1b[~F")     // Hide cursor
	// Frame
	s.TGraphicsOutput.FFillArea(1, 0, 2, 0, 0, []int{0, 0, 639, 0, 639, 249, 0, 249, 0, 0})
	s.TGraphicsOutput.FFillArea(1, 0, 0, 0, 0, []int{3, 2, 636, 2, 636, 247, 3, 247, 3, 2})
	// Logo
	s.TGraphicsOutput.FFillArea(1, 0, 3, 0, 0, []int{7, 218, 312, 218, 312, 244, 7, 244, 7, 218})
	s.TGraphicsOutput.FPolyLine(1, []int{}, 2, 0, 0, []int{7, 218, 312, 218, 312, 244, 7, 244, 7, 218})
	s.TGraphicsOutput.FPlonkLogo(10, 221)
	s.TGraphicsOutput.FPlotCharacterString(0, 3, 3, 1, 1, "Welcome", 244, (249 - 89))
	s.TGraphicsOutput.FPlotCharacterString(0, 3, 3, 2, 1, "Welcome", 242, (249 - 87))
	s.TGraphicsOutput.FFillArea(1, 0, 2, 0, 0, []int{374, 4, 632, 4, 632, 32, 374, 32, 374, 4})
	s.TGraphicsOutput.FFillArea(1, 0, 3, 0, 0, []int{375, 5, 631, 5, 631, 31, 375, 31, 375, 5})
	s.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 1, 1, firmwareVersion, 400, 21)
	s.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 1, 1, serialNumber, 400, 10)
	s.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 1, 1, mainMemSize, 20, 30)
	s.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 1, 1, virtualMemSize, 20, 20)
	s.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 1, 1, totalMemSize, 20, 10)
}

// fakeDosBoot mimicks an old MS-DOS boot
func fakeDosBoot(s *Subbios) {

	// Reset console
	s.TGraphicsOutput.FSetBorderColour(0)
	s.Stdio.Printf("\x1b[c")
	s.Stdio.Printf("\x1b[0~D")
	randDelay(1000, 2000)

	// Print platform information
	platform, _, version, _ := host.PlatformInformation()
	info := fmt.Sprintf("%s - Version %s", platform, version)
	copyright := "\n"
	if strings.Contains(platform, "darwin") {
		copyright = "Copyright (c) Apple Inc. All rights reserved.\n"
	}
	if strings.Contains(platform, "Windows") {
		copyright = "(c) Microsoft Corporation. All rights reserved.\n"
	}
	s.Stdio.Printf("\x1b[30;51m")
	frameWidth := []rune{}
	for _, _ = range info {
		frameWidth = append(frameWidth, 205)
	}
	s.Stdio.Putchars([]rune{201, 205})
	s.Stdio.Putchars(frameWidth)
	s.Stdio.Putchars([]rune{205, 187, '\n', 186, 32})
	s.Stdio.Printf(info)
	s.Stdio.Putchars([]rune{32, 186, '\n'})
	s.Stdio.Putchars([]rune{200, 205})
	s.Stdio.Putchars(frameWidth)
	s.Stdio.Putchars([]rune{205, 188, '\n'})
	s.Stdio.Printf("\x1b[31;50m")
	s.Stdio.Printf(copyright)
	randDelay(1000, 2000)
	s.Stdio.Printf("\nA>")
	randDelay(800, 1500)
	s.Stdio.Printf("autoexec.bat\n") // I know, I know. But the kids are gonna love it
	s.Stdio.KeyboardBufferFlush()
	randDelay(1000, 2000)

}

// Boot simulates the RM Nimbus "Welcome" boot screen and operating system
// loading workflow.  The original Nimbus would also display system info, such
// as firmware version, serial number, memory, etc.  Nimgobus immitates this
// using the Go compiler version as the firmware version, and displays the
// actual physical and virtual memory size.  Serial number is a string constant
// as is the serial number of the Nimbus that provided the ROM dump for the
// emulation on MAME, from which various bits and pieces were reversed
// engineering for nimgobus.
func (s *Subbios) Boot() {
	welcomeScreen(s)
	randDelay(750, 1000)
	s.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 3, 1, "      Please supply an operating system", 137, (249 - 137))
	randDelay(1500, 2000)
	s.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 0, 1, "      Please supply an operating system", 137, (249 - 137))
	s.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 3, 1, "Looking for an operating system - please wait", 137, (249 - 137))
	randDelay(1800, 2900)
	s.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 0, 1, "Looking for an operating system - please wait", 137, (249 - 137))
	s.TGraphicsOutput.FPlotCharacterString(0, 1, 1, 3, 1, "           Loading operating system", 137, (249 - 137))
	randDelay(2100, 2900)
	fakeDosBoot(s)
}
