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
	p int
}

var pallette []color.Color = slices.Repeat([]color.Color{color.White}, 255)

var Map displayMap = displayMap{
	m: ``,
	// m:  "�",
	w: 3,  //width
	s: 20, // scale
	x: 1,  // origin x
	y: 1,  //origin y
	p: 2,  //padding
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	// draw map
	for i := 0; i < len(Map.m); i++ {
		var Colum int = i % Map.w
		var Row int = i / Map.w
		vector.DrawFilledRect(screen, float32((Map.x+Colum)*Map.s), float32((Map.y+Row)*Map.s), float32(Map.s-Map.p), float32(Map.s-Map.p), pallette[int(Map.m[i])], true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	Map.m += string(0) + string(1) + string(2) + string(1) + string(1) + string(1) + string(3) + string(1) + string(0) // test map
	fmt.Sprint(Map.m)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("CheckerWars")
	// TODO Build Color Pallette
	pallette[0] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	pallette[1] = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	pallette[2] = color.RGBA{R: 255, A: 255}
	pallette[3] = color.RGBA{B: 255, A: 255}

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
