//go:build js && wasm

package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"syscall/js"
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

var Zoom float64 = 2
var ZoomSpeed float64 = 0.2

// TODO init centered
var originX = 200
var originY = 200

// old mouse pos
var OldOrignX int = -1
var OldOrignY int = -1

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
	Init     bool
	// workaround for web
	// window.size does not retrun web size
	height int
	width  int
}

func (g *Game) Update() error {
	// TODO activate buttons based on the mouse
	// right click should auto pan
	// left auto select
	// scroll zoom in/out

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
		// operate out of button bar
		if y > int(barHeight) {
			// pan screen
			if buttons[0].toggle {

				// move
				if OldOrignX == -1 {
					OldOrignX = originX - x
					OldOrignY = originY - y
				}
				originX = x + OldOrignX
				originY = y + OldOrignY
			}
		}
	} else {
		// mouse up
		OldOrignX = -1
		OldOrignY = -1
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	ebitenutil.DebugPrintAt(screen, strconv.Itoa(g.width), 200, 200)
	ebitenutil.DebugPrintAt(screen, strconv.Itoa(g.height), 200, 220)
	doc := js.Global().Get("document")
	ebitenutil.DebugPrintAt(screen, doc.Get("URL").String(), 600, 600) // this is url to server
	// web being weird
	barHeight = float64(g.height) / float64(barScale)
	// draw map/settings
	if buttons[3].toggle {
		// TODO make it proper size
		var MapOpt = &ebiten.DrawImageOptions{}
		var Map = ebiten.NewImage(100, 100)
		// TODO Gen map
		Map.Fill(color.White) //temp
		// TODO draw map

		MapOpt.GeoM.Scale(Zoom, Zoom)
		MapOpt.GeoM.Translate(float64(originX), float64(originY))
		// center map
		MapOpt.GeoM.Translate(-1*float64(Map.Bounds().Dx())*Zoom/2, -1*float64(Map.Bounds().Dy())*Zoom/2)

		screen.DrawImage(Map, MapOpt)
		ebitenutil.DebugPrintAt(screen, "MAP", 100, 100)
	} else {
		// TODO draw settings
		ebitenutil.DebugPrintAt(screen, "Settings", 100, 100)
	}

	// other draws
	vector.DrawFilledRect(screen, 0, 0, float32(g.width), float32(barHeight), color.White, true) //top banner/button tray
	// Debug click info
	var b, x, y = getClick()
	ebitenutil.DebugPrintAt(screen, strconv.FormatBool(b)+" : "+strconv.Itoa(x)+" : "+strconv.Itoa(y), 20, g.height-80)
	// draw buttons
	// distance from left edge, buttons get added to this
	var dist = 0.0
	var op = &ebiten.DrawImageOptions{}
	for b := 0; b < len(buttons); b++ {
		// icon sort
		var iconInfo = ebiten.Image{}
		if !buttons[b].toggle {
			iconInfo = *buttons[b].icon1
		} else {
			iconInfo = *buttons[b].icon2
		}
		// scale image
		var scale = barHeight / float64(iconInfo.Bounds().Dy())
		op.GeoM.Scale(float64(scale), float64(scale))
		op.GeoM.Translate(float64(dist), 0)
		screen.DrawImage(&iconInfo, op)
		buttons[b].x = dist
		dist += float64(iconInfo.Bounds().Dx()) * scale
		op.GeoM.Reset()

	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.width = outsideWidth
	g.height = outsideHeight
	return outsideWidth, outsideHeight
}

func main() {
	// ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("CheckerWars")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// temp colors/button icons
	// toggle pan/select
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
	// zoom in
	buttons[1].icon1 = ebiten.NewImage(40, 40)
	buttons[1].icon1.Fill(r)
	buttons[1].run = zoomIn
	// zoom out
	buttons[2].icon1 = ebiten.NewImage(20, 20)
	buttons[2].icon1.Fill(g)
	buttons[2].run = zoomOut
	// toggle map/settings
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
func quit(toggle bool) {}

// finished
func zoomIn(toggle bool) {
	if Zoom > 2*ZoomSpeed {
		Zoom -= ZoomSpeed
	}
}
func zoomOut(toggle bool) { Zoom += ZoomSpeed }

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
