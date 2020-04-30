package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png" // import only for side-effects
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adamstimb/nimgobus"
	"github.com/adamstimb/nimgobus/examples/isstracker/issImages"
	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 1400
	screenHeight = 1000
)

var (
	nim nimgobus.Nimbus
)

type Game struct {
	count int
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.count == 0 {
		go Test()
	}
	g.count++
	nim.Update()
	return nil
}

func SplashScreen() {
	// Load images from string vars
	img, _, err := image.Decode(bytes.NewReader(issImages.Iss))
	if err != nil {
		log.Fatal(err)
	}
	issImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	nim.SetMode(40)
	nim.Fetch(issImg, 1)
	nim.Writeblock(1, 0, 0)
	plotOpts := nimgobus.PlotOptions{
		Brush: 14, Font: 1, SizeX: 1, SizeY: 2,
	}
	nim.Plot(plotOpts, "ISS TRACKER", 230, 30)
	plotOpts.SizeY = 1
	nim.Plot(plotOpts, "copyright (c) P.P. Bottoms-Farts 1986", 22, 20)
	nim.Plot(plotOpts, "Fegg-Heyes Primary School, North Staffs", 8, 10)
	time.Sleep(3 * time.Second)
}

func getPosition() (float64, float64) {
	r, _ := http.Get("http://api.open-notify.org/iss-now.json")
	type Position struct {
		Latitude  string
		Longitude string
	}
	type ApiBody struct {
		Timestamp    int
		Message      string
		Iss_position Position
	}
	var body ApiBody
	rawBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(string(rawBody)), &body)
	long, _ := strconv.ParseFloat(body.Iss_position.Longitude, 64)
	lat, _ := strconv.ParseFloat(body.Iss_position.Latitude, 64)
	return long, lat
}

func Track() {
	nim.SetMode(80)
	nim.SetColour(0, 9)
	nim.SetColour(1, 1)
	nim.SetColour(2, 2)
	nim.SetBorder(0)
	nim.SetPaper(0)
	nim.SetCharset(1)
	nim.SetPen(3)
	nim.Cls()
	img, _, err := image.Decode(bytes.NewReader(issImages.World))
	if err != nil {
		log.Fatal(err)
	}
	worldImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	nim.Fetch(worldImg, 1)
	nim.Writeblock(1, 0, 0)
	longScale := 640.0 / 360.0
	latScale := 250.0 / 180.0
	for {
		long, lat := getPosition()
		circleOpts := nimgobus.CircleOptions{
			Brush: 2,
		}
		x := int(long*longScale) + 320
		if x > 640 {
			x -= 640
		}
		y := 125 + int(lat*latScale)
		nim.Circle(circleOpts, 4, x, y)
		nim.SetCurpos(1, 1)
		status := fmt.Sprintf("Latitude: %f  Longitude: %f   ", long, lat)
		nim.Print(status)
		time.Sleep(1 * time.Second)
	}
}

func Test() {
	SplashScreen()
	Track()
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Draw the Nimbus monitor on the screen and scale to current window size.
	monitorWidth, monitorHeight := nim.Monitor.Size()
	// Calculate aspect ratios of Nimbus monitor and ebiten screen
	monitorRatio := float64(monitorWidth) / float64(monitorHeight)
	screenRatio := float64(screenWidth) / float64(screenHeight)

	// If screenRatio > monitorRatio then clamp monitorHeight to screenHeight otherwise
	// clamp monitorWidth to screenWidth
	var scale, offsetX, offsetY float64
	if screenRatio > monitorRatio {
		scale = float64(screenHeight) / float64(monitorHeight)
		offsetX = (float64(screenWidth) - float64(monitorWidth)*scale) / 2
		offsetY = 0
	}
	if screenRatio <= monitorRatio {
		scale = float64(screenWidth) / float64(monitorWidth)
		offsetX = 0
		offsetY = (float64(screenHeight) - float64(monitorHeight)*scale) / 2
	}

	// Apply scale and centre monitor on screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(offsetX, offsetY)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(nim.Monitor, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {

	// Initialize the Nimbus
	nim.Init()

	// set up window
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Nimgobus Test")

	// Call RunGame method, passing the address of the pointer to an empty Game struct
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
