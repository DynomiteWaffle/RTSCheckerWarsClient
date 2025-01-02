package main

import (
	"fmt"
	"image/color"
	"log"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type displayMap struct {
	m string
	w int
	s int
	x int
	y int
}

var pallette []color.Color = slices.Repeat([]color.Color{color.White}, 255)

var Map displayMap = displayMap{
	m: "012111310",
	w: 3,
	s: 20,
	x: 1,
	y: 1,
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	// draw map
	for i := 0; i < len(Map.m); i++ {
		fmt.Println(int(Map.m[i]) - 48)
		var Colum int = i % Map.w
		var Row int = i / Map.w
		vector.DrawFilledRect(screen, float32((Map.x+Colum)*Map.s), float32((Map.y+Row)*Map.s), float32(Map.s), float32(Map.s), pallette[int(Map.m[i])-48], true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	// TODO Build Color Pallette
	pallette[0] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	pallette[1] = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	pallette[2] = color.RGBA{R: 255, A: 255}
	pallette[3] = color.RGBA{B: 255, A: 255}

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
