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
	x int
	y int
	p int
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
		vector.DrawFilledRect(screen, float32((Map.x+Colum)*Map.s), float32((Map.y+Row)*Map.s), float32(Map.s-Map.p), float32(Map.s-Map.p), pallette[Map.m[i]], true)
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
