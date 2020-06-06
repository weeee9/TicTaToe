package game

import (
	"errors"
	"image/color"
	"log"

	"github.com/weeee9/tic-tac-toe/textimage"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"

	"golang.org/x/image/font"
)

type player int

// TicTacToe tic tac toe game type
type TicTacToe struct {
	board      [3][3]player
	isGameOver bool
	blocksFill int
}

var (
	playerX player = 1
	playerO player = -1
	non     player = 0

	currentPlayer player

	white = color.White
	black = color.Black

	symbolX string = "X"
	symbolO string = "O"

	normalFont font.Face

	imageGameover *ebiten.Image
)

var _ ebiten.Game = &TicTacToe{}

func init() {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	normalFont = truetype.NewFace(tt, &truetype.Options{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	// Gameover image
	imageGameover, _ = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	imageGameover.Fill(color.NRGBA{0x00, 0x00, 0x00, 0x80})
	y := (screenHeight - blockHeight) / 2
	textimage.DrawTextWithShadowCenter(imageGameover, gameoverStr, 0, y, 1, color.White, screenWidth)

	// set currentPlayer
	currentPlayer = playerX
}

// Update updates a game by one tick. The given argument represents a screen image.
func (ttt *TicTacToe) Update(screen *ebiten.Image) error {
	posX, posY := ebiten.CursorPosition()
	x, y := mousePosTotXY(posX, posY)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && cursorInScreen(posX, posY) {
		ttt.setBlock(currentPlayer, x, y)
	}
	return ttt.Draw(screen)
}

// Layout accepts a native outside size in device-independent pixels and returns the game's logical screen
// size.
func (ttt *TicTacToe) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Draw draw game screen
func (ttt *TicTacToe) Draw(screen *ebiten.Image) error {
	screen.Clear()
	screen.Fill(white)
	ttt.putBlock(screen)

	if ttt.isGameOver {
		screen.DrawImage(imageGameover, &ebiten.DrawImageOptions{})
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			ttt.reset()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return errors.New("END")
		}
		return nil
	}

	return nil
}

// putBlock put block on game screen
func (ttt *TicTacToe) putBlock(screen *ebiten.Image) {
	for i := range ttt.board {
		for j, p := range ttt.board[i] {
			opts := &ebiten.DrawImageOptions{}
			// 修改選項，新增 Translate 變形效果
			opts.GeoM.Translate(float64(1*(j+1)+j*blockWidth), float64(1*(i+1)+i*blockHeight))
			block := newBlock(p)
			screen.DrawImage(block, opts)
		}
	}
}

// setBlock set block to "X" or "O"
func (ttt *TicTacToe) setBlock(player player, x, y int) {
	if blockIsFill(ttt.board, x, y) {
		return
	}
	ttt.board[x][y] = player
	ttt.blocksFill++
	if ttt.blocksFill == 9 {
		ttt.isGameOver = true
	}
	if checkWin(currentPlayer, ttt.board) {
		ttt.isGameOver = true
	}
	currentPlayer *= -1
}

func (ttt *TicTacToe) reset() {
	for i := range ttt.board {
		for j := range ttt.board[i] {
			ttt.board[i][j] = 0
		}
	}
	ttt.isGameOver = false
	ttt.blocksFill = 0
}

// newBlock return a new block image
func newBlock(player player) *ebiten.Image {
	block, err := ebiten.NewImage(blockWidth, blockHeight, ebiten.FilterDefault)
	if err != nil {
		return nil
	}
	block.Fill(black)
	if player == playerX {
		text.Draw(block, symbolX, normalFont, xCenter, yCenter, white)
	} else if player == playerO {
		text.Draw(block, symbolO, normalFont, xCenter, yCenter, white)
	}
	return block
}
