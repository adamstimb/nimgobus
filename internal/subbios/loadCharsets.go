package subbios

import (
	"bytes"
	"image"
	"log"

	"github.com/adamstimb/nimgobus/internal/make2darray"
	"github.com/adamstimb/nimgobus/internal/resources/font"
)

// charImageSelecta returns the subimage of a char from the charset image for any given ASCII code or charset.
func (v *video) charImageSelecta(img [][]int, c int, charSet int) [][]int {

	if charSet == 1 {
		c += 256 // because char 0 of charset 1 is the 256th char plotted in the chardump image
	}

	// Scan along image until we find the right char.  This is slow but is only done once on startup.
	x1 := -8
	y1 := 39
	for i := 0; i <= c; i++ {
		x1 += 8
		if x1 >= 639 {
			x1 = 0
			y1 += 10
		}
	}

	// Get char and return
	x2 := x1 + 7
	y2 := y1 + 9
	charImgArray := make2darray.Make2dArray(8, 10, -1)
	charImgArrayX := 0
	for x := x1; x <= x2; x++ {
		charImgArrayY := 0
		for y := y1; y <= y2; y++ {
			charImgArray[charImgArrayY][charImgArrayX] = img[y][x]
			charImgArrayY++
		}
		charImgArrayX++
	}

	return charImgArray
}

// loadCharsetImages loads the charset images
func (v *video) loadCharsetImages(charset int) {

	// convertToArray receives an char image and returns it as black-and-white 2d array
	convertToArray := func(img image.Image) [][]int {
		b := img.Bounds()
		width := b.Max.X
		height := b.Max.Y
		newArray := make2darray.Make2dArray(width, height, -1)
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				// Get colour at x, y and if black set it to 1 in the 2d array
				c := img.At(x, y)
				r, g, b, _ := c.RGBA()
				if r == 0 && g == 0 && b == 0 {
					newArray[y][x] = 1
				}
			}
		}
		return newArray
	}

	var imgArray [][]int
	img, _, err := image.Decode(bytes.NewReader(font.Charsets_png))
	if err != nil {
		log.Fatal(err)
	}
	imgArray = convertToArray(img)
	for i := 0; i <= 255; i++ {
		if charset == 0 {
			v.charSet0[i] = v.charImageSelecta(imgArray, i, charset)
		} else {
			v.charSet1[i] = v.charImageSelecta(imgArray, i, charset)
		}
	}
}
