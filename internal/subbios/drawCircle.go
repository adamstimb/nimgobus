package subbios

import "math"

// drawCircle implements the midpoint circle drawing algorithm to draw a
// circle in a 2D array. Looking at circles drawn in BASIC on the Nimbus emulator
// this algorithm gives very similar results, so it was likely a midpoint jobbie as well.
func (v *video) drawCircle(img [][]int, xCentre, yCentre, radius, theta1, theta2, colour int) {
	x := radius
	y := 0
	d := 0
	// generate points of the whole cirtle using the midpoint algorithm
	points := []coord{}
	for x >= y {
		points = append(points, coord{X: xCentre + x, Y: yCentre + y})
		points = append(points, coord{X: xCentre + y, Y: yCentre + x})
		points = append(points, coord{X: xCentre - y, Y: yCentre + x})
		points = append(points, coord{X: xCentre - x, Y: yCentre + y})
		points = append(points, coord{X: xCentre - x, Y: yCentre - y})
		points = append(points, coord{X: xCentre - y, Y: yCentre - x})
		points = append(points, coord{X: xCentre + y, Y: yCentre - x})
		points = append(points, coord{X: xCentre + x, Y: yCentre - y})
		if d <= 0 {
			y = y + 1
			d = d + 2*y + 1
		} else {
			x = x - 1
			d = d - 2*x + 1
		}
	}
	// draw the filled cirtle within the range of theta1 and theta2
	for _, point := range points {
		theta := int((math.Atan2(float64(-point.X+radius), float64(-point.Y+radius)) + math.Pi) * 1000.0)
		if (theta >= theta1 && theta <= theta2) || // Draw within range of theta1 and theta1.
			(theta1 == theta2) || // Ensure a full cirtle is drawn if they're the same value.
			(theta1 == 0 && point.X == xCentre && point.Y > yCentre) || // This is some fudging to make sure 0, 90, 180, 270
			(theta1 == 1571 && point.X > xCentre && point.Y == yCentre) || // start angles are straight and don't have an
			(theta1 == 3142 && point.X == xCentre && point.Y < yCentre) || // annoying wonkyness to them.
			(theta1 == 4713 && point.X < xCentre && point.Y == yCentre) {
			img = v.drawLine(img, xCentre, yCentre, point.X, point.Y, 1, []int{0}, colour, 0, 0)
		}
	}
}
