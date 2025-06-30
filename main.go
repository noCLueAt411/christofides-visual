package main

import (
	"christofides-algo/visual"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	//Make it possible to generate new graph on press 1
	//Make it possible to generate show points
	//Make it possible to solve it with 2
	//Make it posible to show solution with 3

	var windowWidth float32 = 1200
	var windowHeight float32 = 600

	game := visual.NewGame(40, windowWidth, windowHeight)

	ebiten.SetWindowTitle("Christofides Visualizer")
	ebiten.SetWindowSize(int(windowWidth), int(windowHeight))

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
