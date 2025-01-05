package main

import (
	"image/color"
	"log"
	"slices"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type displayMap struct {
	m []int
	w int
	s int
	o point
	p int
}
type point struct {
	x int
	y int
}

var pallette []color.Color = slices.Repeat([]color.Color{color.White}, 255)

// TODO convert json hex map to array of ints
func Hex2Map(input string) []int {
	var nums []int = []int{}
	for i := 0; i < len(input); i += 2 {
		// fmt.Println(input[i : i+2])
		num, _ := strconv.ParseInt(input[i:i+2], 16, 16)
		nums = append(nums, int(num))

	}
	return nums
}

var Map displayMap = displayMap{
	m: []int{},
	w: 3,  //width
	s: 20, // scale
	p: 2,  //padding
}
var mouseDist point
var oldMousePos point
var scrollSpeed int = 10

type Game struct{}

func (g *Game) Update() error {
	// move map
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton2) {
		mouseDist.x, mouseDist.y = ebiten.CursorPosition()
		mouseDist.x -= oldMousePos.x
		mouseDist.y -= oldMousePos.y

	} else {
		oldMousePos.x, oldMousePos.y = ebiten.CursorPosition()
		Map.o.x += mouseDist.x
		Map.o.y += mouseDist.y
		mouseDist.x, mouseDist.y = 0, 0

	}
	// zoom map
	// TODO zoom from center
	var _, wheel = ebiten.Wheel()
	Map.s += int(wheel) * scrollSpeed
	if Map.s < 1 {
		Map.s = 1
	}
	// TODO mobile touches/pans

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello World")
	// draw map
	for i := 0; i < len(Map.m); i++ {
		var Colum int = i % Map.w
		var Row int = i / Map.w
		vector.DrawFilledRect(screen, float32(mouseDist.x+Map.o.x+(Colum)*Map.s), float32(mouseDist.y+Map.o.y+(Row)*Map.s), float32(Map.s-Map.p), float32(Map.s-Map.p), pallette[Map.m[i]], true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	Map.m = Hex2Map("000102010101030100")
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
