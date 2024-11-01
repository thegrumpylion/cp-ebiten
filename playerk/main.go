package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Player")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
