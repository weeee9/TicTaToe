package main

import (
	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

const (
	screenWidth  int = 320
	screenHeight int = 320
	blockWidth   int = (screenWidth - 4) / 3
	blockHeight  int = (screenHeight - 4) / 3

	fontSize float64 = 36

	xCenter int = blockHeight/2 - int(fontSize/4)
	yCenter int = blockWidth/2 + int(fontSize/4)
)

type Player int

type TicTacToe struct {
	Blocks [3][3]Player
}

var (
	white = color.White
	black = color.Black

	player1 Player = 1
	player2 Player = 2
	non     Player = 0

	x string = "X"
	o string = "O"

	ttt           = &TicTacToe{}
	noramlFont    font.Face
	currentPlayer = player1
)

func init() {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	noramlFont = truetype.NewFace(tt, &truetype.Options{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

}

func main() {
	ebiten.Run(update, screenWidth, screenHeight, 1, "Hello World!")

}

func update(screen *ebiten.Image) error {
	screen.Fill(white)
	ttt.putBlock(screen)
	x, y := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && cursorInScreen(x, y) {
		ttt.setBlock(currentPlayer, x, y)
	}
	return nil
}

func newBlock(player Player) *ebiten.Image {
	block, err := ebiten.NewImage(blockWidth, blockHeight, ebiten.FilterDefault)
	if err != nil {
		return nil
	}
	block.Fill(black)
	if player == player1 {
		text.Draw(block, x, noramlFont, xCenter, yCenter, white)
	} else if player == player2 {
		text.Draw(block, o, noramlFont, xCenter, yCenter, white)
	}

	return block
}

func (ttt *TicTacToe) putBlock(screen *ebiten.Image) {
	for i := range ttt.Blocks {
		for j, p := range ttt.Blocks[i] {
			opts := &ebiten.DrawImageOptions{}
			// 修改選項，新增 Translate 變形效果
			opts.GeoM.Translate(float64(1*(j+1)+j*blockWidth), float64(1*(i+1)+i*blockHeight))
			block := newBlock(p)
			screen.DrawImage(block, opts)
		}
	}
}

func (ttt *TicTacToe) setBlock(player Player, x, y int) {
	xLine := y / blockWidth
	yLine := x / blockHeight
	ttt.Blocks[xLine][yLine] = player
}

func (ttt *TicTacToe) checkWin() {}

func cursorInScreen(x, y int) bool {
	return (x > 0 && y > 0) && (x < 320 && y < 320)
}
