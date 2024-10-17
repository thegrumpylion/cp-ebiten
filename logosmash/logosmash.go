package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp/v2"
	"github.com/jakecoffman/cpebiten"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 600
	screenHeight = 480
)

type Game struct {
	*cpebiten.Game
	dot *ebiten.Image
}

func NewGame() *Game {
	const (
		imageWidth  = 188
		imageHeight = 35
	)

	dot := ebiten.NewImage(1, 1)
	dot.Fill(color.White)

	space := cp.NewSpace()
	space.Iterations = 1

	// The space will contain a very large number of similarly sized objects.
	// This is the perfect candidate for using the spatial hash.
	// Generally you will never need to do this.
	space.UseSpatialHash(2.0, 10000)

	var body *cp.Body
	var shape *cp.Shape

	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			if getPixel(uint(x), uint(y)) == 0 {
				continue
			}

			xJitter := 0.05 * rand.Float64()
			yJitter := 0.05 * rand.Float64()

			shape = makeBall(2.0*(float64(x)+imageWidth/2+xJitter)-75, 2*(imageHeight/2.0+float64(y)+yJitter)+150)
			space.AddBody(shape.Body())
			space.AddShape(shape)
		}
	}

	body = space.AddBody(cp.NewBody(1e9, cp.INFINITY))
	body.SetPosition(cp.Vector{X: -1000, Y: 225})
	body.SetVelocity(400, 0)

	shape = space.AddShape(cp.NewCircle(body, 8, cp.Vector{}))
	shape.SetElasticity(0)
	shape.SetFriction(0)

	return &Game{
		Game: cpebiten.NewGame(space, 60),
		dot:  dot,
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebiten can't handle drawing this many triangles in one go, so drawing with image instead
	op := &ebiten.DrawImageOptions{}
	g.Space.EachBody(func(body *cp.Body) {
		op.GeoM.Translate(body.Position().X, body.Position().Y)
		screen.DrawImage(g.dot, op)
		op.GeoM.Reset()
	})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
}

func (g *Game) Layout(int, int) (int, int) {
	return screenWidth, screenHeight
}

func getPixel(x, y uint) int {
	const imageRowLength = 24
	return (imageBitmap[(x>>3)+y*imageRowLength] >> (^x & 0x7)) & 1
}

func makeBall(x, y float64) *cp.Shape {
	body := cp.NewBody(1.0, cp.INFINITY)
	body.SetPosition(cp.Vector{X: x, Y: y})

	shape := cp.NewCircle(body, 0.95, cp.Vector{})
	shape.SetElasticity(0)
	shape.SetFriction(0)

	return shape
}

var imageBitmap = []int{
	15, -16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, -64, 15, 63, -32, -2, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 31, -64, 15, 127, -125, -1, -128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 127, -64, 15, 127, 15, -1, -64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, -1, -64, 15, -2,
	31, -1, -64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, -1, -64, 0, -4, 63, -1, -32, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, -1, -64, 15, -8, 127, -1, -32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, -1, -64, 0, -8, -15, -1, -32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, -31, -1, -64, 15, -8, -32,
	-1, -32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, -15, -1, -64, 9, -15, -32, -1, -32, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 31, -15, -1, -64, 0, -15, -32, -1, -32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 63, -7, -1, -64, 9, -29, -32, 127, -61, -16, 63, 15, -61, -1, -8, 31, -16, 15, -8, 126, 7, -31,
	-8, 31, -65, -7, -1, -64, 9, -29, -32, 0, 7, -8, 127, -97, -25, -1, -2, 63, -8, 31, -4, -1, 15, -13,
	-4, 63, -1, -3, -1, -64, 9, -29, -32, 0, 7, -8, 127, -97, -25, -1, -2, 63, -8, 31, -4, -1, 15, -13,
	-2, 63, -1, -3, -1, -64, 9, -29, -32, 0, 7, -8, 127, -97, -25, -1, -1, 63, -4, 63, -4, -1, 15, -13,
	-2, 63, -33, -1, -1, -32, 9, -25, -32, 0, 7, -8, 127, -97, -25, -1, -1, 63, -4, 63, -4, -1, 15, -13,
	-1, 63, -33, -1, -1, -16, 9, -25, -32, 0, 7, -8, 127, -97, -25, -1, -1, 63, -4, 63, -4, -1, 15, -13,
	-1, 63, -49, -1, -1, -8, 9, -57, -32, 0, 7, -8, 127, -97, -25, -8, -1, 63, -2, 127, -4, -1, 15, -13,
	-1, -65, -49, -1, -1, -4, 9, -57, -32, 0, 7, -8, 127, -97, -25, -8, -1, 63, -2, 127, -4, -1, 15, -13,
	-1, -65, -57, -1, -1, -2, 9, -57, -32, 0, 7, -8, 127, -97, -25, -8, -1, 63, -2, 127, -4, -1, 15, -13,
	-1, -1, -57, -1, -1, -1, 9, -57, -32, 0, 7, -1, -1, -97, -25, -8, -1, 63, -1, -1, -4, -1, 15, -13, -1,
	-1, -61, -1, -1, -1, -119, -57, -32, 0, 7, -1, -1, -97, -25, -8, -1, 63, -1, -1, -4, -1, 15, -13, -1,
	-1, -61, -1, -1, -1, -55, -49, -32, 0, 7, -1, -1, -97, -25, -8, -1, 63, -1, -1, -4, -1, 15, -13, -1,
	-1, -63, -1, -1, -1, -23, -49, -32, 127, -57, -1, -1, -97, -25, -1, -1, 63, -1, -1, -4, -1, 15, -13,
	-1, -1, -63, -1, -1, -1, -16, -49, -32, -1, -25, -1, -1, -97, -25, -1, -1, 63, -33, -5, -4, -1, 15,
	-13, -1, -1, -64, -1, -9, -1, -7, -49, -32, -1, -25, -8, 127, -97, -25, -1, -1, 63, -33, -5, -4, -1,
	15, -13, -1, -1, -64, -1, -13, -1, -32, -49, -32, -1, -25, -8, 127, -97, -25, -1, -2, 63, -49, -13,
	-4, -1, 15, -13, -1, -1, -64, 127, -7, -1, -119, -17, -15, -1, -25, -8, 127, -97, -25, -1, -2, 63,
	-49, -13, -4, -1, 15, -13, -3, -1, -64, 127, -8, -2, 15, -17, -1, -1, -25, -8, 127, -97, -25, -1,
	-8, 63, -49, -13, -4, -1, 15, -13, -3, -1, -64, 63, -4, 120, 0, -17, -1, -1, -25, -8, 127, -97, -25,
	-8, 0, 63, -57, -29, -4, -1, 15, -13, -4, -1, -64, 63, -4, 0, 15, -17, -1, -1, -25, -8, 127, -97,
	-25, -8, 0, 63, -57, -29, -4, -1, -1, -13, -4, -1, -64, 31, -2, 0, 0, 103, -1, -1, -57, -8, 127, -97,
	-25, -8, 0, 63, -57, -29, -4, -1, -1, -13, -4, 127, -64, 31, -2, 0, 15, 103, -1, -1, -57, -8, 127,
	-97, -25, -8, 0, 63, -61, -61, -4, 127, -1, -29, -4, 127, -64, 15, -8, 0, 0, 55, -1, -1, -121, -8,
	127, -97, -25, -8, 0, 63, -61, -61, -4, 127, -1, -29, -4, 63, -64, 15, -32, 0, 0, 23, -1, -2, 3, -16,
	63, 15, -61, -16, 0, 31, -127, -127, -8, 31, -1, -127, -8, 31, -128, 7, -128, 0, 0,
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Logosmash")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
