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

	game := visual.NewGame(6)

	ebiten.SetWindowTitle("Christofides Visualizer")
	ebiten.SetWindowSize(800, 600)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
