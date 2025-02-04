package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// | 0 - toggle pan/select |
// 1 - zoom in |
// 2 - zoom out |
// 3 - toggle settings |
var buttons = [4]Button{}
var barScale = 12
var panning = false
var barHeight float64
var gray = color.RGBA{R: 20, G: 20, B: 20, A: 255}
var prevClick bool

// debug
var (
	r = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	g = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	b = color.RGBA{R: 0, G: 0, B: 255, A: 255}
)

type Piece struct {
	piece   int
	timeout int
}

type Button struct {
	run       func(toggle bool)
	toggle    bool
	togglable bool
	icon1     *ebiten.Image
	icon2     *ebiten.Image
	x         float64
}

func (b *Button) Toggle() {
	if b.toggle {
		b.toggle = false
	} else {
		b.toggle = true
	}
}

type Game struct {
	Map      []Piece
	MapWidth int
	MapType  int
}

func (g *Game) Update() error {
	// click/activate buttons
	var clicked, x, y = getClick()
	// click buttons only once
	if !clicked {
		prevClick = false
	}
	// main button loop loop
	if !prevClick && clicked && y > 0 && y < int(barHeight) {
		prevClick = true
		// var offset = 0
		for b := 0; b < len(buttons); b++ {
			if !buttons[b].toggle {
				// icon 1
				var scale = barHeight / float64(buttons[b].icon1.Bounds().Dy())
				// in x bounds
				if x > int(buttons[b].x) && x < int(buttons[b].x)+buttons[b].icon1.Bounds().Dx()*int(scale) {
					if buttons[b].togglable {
						buttons[b].Toggle()
					}
					buttons[b].run(buttons[b].toggle)
				}

			} else {
				// icon 2
				var scale = barHeight / float64(buttons[b].icon2.Bounds().Dy())
				// in x bounds
				if x > int(buttons[b].x) && x < int(buttons[b].x)+buttons[b].icon2.Bounds().Dx()*int(scale) {
					if buttons[b].togglable {
						buttons[b].Toggle()
					}
					buttons[b].run(buttons[b].toggle)
				}

			}

		}
	} else if clicked {
		prevClick = true
		// do something else with mouse input
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var width, height = ebiten.WindowSize()
	// web being weird
	barHeight = float64(height) / float64(barScale)
	if height == 0 {
		height = 1080
		barHeight = 50
	}
	if width == 0 {
		width = 1080
	}
	vector.DrawFilledRect(screen, 0, 0, float32(width), float32(barHeight), color.White, true)
	// message for web with too big screen
	vector.DrawFilledRect(screen, float32(width), 0, float32(width), 1080, gray, true)
	ebitenutil.DebugPrintAt(screen, "get the app for fullscreen", width+10, 0)
	ebitenutil.DebugPrintAt(screen, "github.com/DynomiteWaffle/CheckerWarsClient", width+10, 30) //github link
	ebitenutil.DebugPrintAt(screen, "dynomitewaffle.itch.io/checker-wars", width+10, 50)         //itch.io link
	// Debug click info
	var b, x, y = getClick()
	ebitenutil.DebugPrintAt(screen, strconv.FormatBool(b)+" : "+strconv.Itoa(x)+" : "+strconv.Itoa(y), width+10, 80)
	// draw buttons
	// distance from left edge, buttons get added to this
	var dist = 0.0
	var op = &ebiten.DrawImageOptions{}
	for b := 0; b < len(buttons); b++ {
		if !buttons[b].toggle {
			// icon 1
			// scale image
			var scale = barHeight / float64(buttons[b].icon1.Bounds().Dy())
			op.GeoM.Scale(float64(scale), float64(scale))
			op.GeoM.Translate(float64(dist), 0)
			screen.DrawImage(buttons[b].icon1, op)
			buttons[b].x = dist
			dist += float64(buttons[b].icon1.Bounds().Dx()) * scale
			op.GeoM.Reset()
		} else {
			// icon 2
			var scale = barHeight / float64(buttons[b].icon2.Bounds().Dy())
			op.GeoM.Scale(float64(scale), float64(scale))
			op.GeoM.Translate(float64(dist), 0)
			screen.DrawImage(buttons[b].icon2, op)
			buttons[b].x = dist
			dist += float64(buttons[b].icon2.Bounds().Dx()) * scale
			op.GeoM.Reset()
		}
	}
	// draw map/settings
	if buttons[3].toggle {
		// TODO draw map
		ebitenutil.DebugPrintAt(screen, "MAP", 100, 100)
	} else {
		// TODO draw settings
		ebitenutil.DebugPrintAt(screen, "Settings", 100, 100)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	// ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("CheckerWars")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// temp colors/button icons
	buttons[0].icon1 = ebiten.NewImage(30, 20)
	buttons[0].icon1.Fill(color.Black)
	buttons[0].icon2 = ebiten.NewImage(30, 20)
	buttons[0].icon2.Fill(color.White)
	buttons[0].run = func(toggle bool) {
		if toggle {
			fmt.Println("Button1 toggle 1")
		} else {
			fmt.Println("Button1 toggle 2")

		}
	}
	buttons[0].togglable = true
	buttons[1].icon1 = ebiten.NewImage(40, 40)
	buttons[1].icon1.Fill(r)
	buttons[1].run = func(toggle bool) {
		fmt.Println("Button2")
	}
	buttons[2].icon1 = ebiten.NewImage(20, 20)
	buttons[2].icon1.Fill(g)
	buttons[2].run = func(toggle bool) {
		fmt.Println("Button3")
	}
	buttons[3].togglable = true
	buttons[3].icon1 = ebiten.NewImage(20, 20)
	buttons[3].icon1.Fill(b)
	buttons[3].icon2 = ebiten.NewImage(20, 20)
	buttons[3].icon2.Fill(color.White)
	buttons[3].run = func(toggle bool) {
		fmt.Println("Button4")
	}
	// setup buttons
	//
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

// button functions
// TODO these functions
func quit(toggle bool)      {}
func zoomIn(toggle bool)    {}
func zoomOut(toggle bool)   {}
func togglePan(toggle bool) {}

// TODO
func readMap(Map string) []Piece {
	// get map format
	// prosses map format
	switch 1 {
	case 1:
		break
	case 2:
		break
	}
	// return prossesed map format
	return []Piece{}
}

// Returns if clicked,x,y
func getClick() (bool, int, int) {
	// vars
	var (
		succeded bool = false
		x        int
		y        int
		touches  []ebiten.TouchID
	)
	// update touches
	touches = ebiten.AppendTouchIDs(touches)
	// if empty look at mouse
	if len(touches) == 0 {
		// mouse is pressed
		if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			succeded = true
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButton1) {
			succeded = true
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButton2) {
			succeded = true
		}
		// mouse pos
		x, y = ebiten.CursorPosition()

	} else {
		succeded = true
		// touches
		x, y = ebiten.TouchPosition(0)
	}

	return succeded, x, y
}
