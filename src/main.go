package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TODO web wont render

// | 0 - toggle pan/select |
// 1 - zoom in |
// 2 - zoom out |
// 3 - toggle settings |
var buttons = [4]Button{}
var barScale = 12
var panning = false

var (
	r = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	g = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	b = color.RGBA{R: 0, G: 0, B: 255, A: 255}
)

type Game struct {
}

type Button struct {
	run       func(toggle bool)
	toggle    bool
	togglable bool
	icon1     *ebiten.Image
	icon2     ebiten.Image
}

func (b *Button) Toggle() {
	if b.toggle {
		b.toggle = false
	} else {
		b.toggle = true
	}
}

func (g *Game) Update() error {
	// TODO click/activate buttons
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var width, height = ebiten.WindowSize()
	vector.DrawFilledRect(screen, 0, 0, float32(width), float32(height)/float32(barScale), color.White, true)
	// draw buttons
	var barHeight float64
	barHeight = float64(height) / float64(barScale)
	// distance from left edge, buttons get added to this
	var dist = 0.0
	var op = &ebiten.DrawImageOptions{}
	for b := 0; b < len(buttons); b++ {
		if !buttons[b].toggle {
			// scale image
			var scale = barHeight / float64(buttons[b].icon1.Bounds().Dy())
			op.GeoM.Scale(float64(scale), float64(scale))
			op.GeoM.Translate(float64(dist), 0)
			screen.DrawImage(buttons[b].icon1, op)
			dist += float64(buttons[b].icon1.Bounds().Dx()) * scale
			op.GeoM.Reset()
		} else {
			// TODO draw icon 2
			var scale = barHeight / float64(buttons[b].icon2.Bounds().Dy())
			op.GeoM.Scale(float64(scale), float64(scale))
			op.GeoM.Translate(float64(dist), 0)
			screen.DrawImage(&buttons[b].icon2, op)
			dist += float64(buttons[b].icon2.Bounds().Dx()) * scale
			op.GeoM.Reset()
		}
	}
	// screen.DrawImage(buttons[0].icon1, op)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("CheckerWars")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	buttons[0].icon1 = ebiten.NewImage(20, 20)
	buttons[0].icon1.Fill(color.Black)
	buttons[1].icon1 = ebiten.NewImage(40, 40)
	buttons[1].icon1.Fill(r)
	buttons[2].icon1 = ebiten.NewImage(20, 20)
	buttons[2].icon1.Fill(g)
	buttons[3].icon1 = ebiten.NewImage(20, 20)
	buttons[3].icon1.Fill(b)
	// setup buttons
	//
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
