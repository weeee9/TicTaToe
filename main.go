package main

import (
	"log"

	"github.com/weeee9/tic-tac-toe/game"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  int = 320
	screenHeight int = 320
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tic Tac Toe")

	ttt := &game.TicTacToe{}

	if err := ebiten.RunGame(ttt); err != nil {
		log.Fatal(err)
	}
}
