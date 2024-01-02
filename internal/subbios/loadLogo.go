package subbios

import (
	"bytes"
	"image"
	"log"

	"github.com/adamstimb/nimgobus/internal/make2darray"
	"github.com/adamstimb/nimgobus/internal/resources/logo"
)

// loadLogoImage loads the Nimbus logo image
func (s *Subbios) loadLogoImage() {

	// convertToArray receives the logo image and returns it as 3-colour 2d array
	convertToArray := func(img image.Image) [][]int {
		b := img.Bounds()
		width := b.Max.X
		height := b.Max.Y
		newArray := make2darray.Make2dArray(width, height, -1)
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				// Get colour at x, y and downsample to 3 colours
				c := img.At(x, y)
				r, g, b, a := c.RGBA()
				if r > 60000 && g > 60000 && b > 60000 && a > 60000 {
					// white
					newArray[y][x] = 3
					continue
				}
				if r > 60000 && (g > 20000 && g < 30000) && (b > 2000 && b < 30000) && a > 60000 {
					// red
					newArray[y][x] = 2
					continue
				}
				if r < 10000 && g < 10000 && b < 10000 && a > 60000 {
					// green
					newArray[y][x] = 1
					continue
				}
			}
		}
		return newArray
	}

	img, _, err := image.Decode(bytes.NewReader(logo.NimbusLogoFinal_png))
	if err != nil {
		log.Fatal(err)
	}
	s.TGraphicsOutput.v.logo = convertToArray(img)
}
