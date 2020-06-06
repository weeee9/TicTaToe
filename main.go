package main

import (
	"errors"
	"image/color"
	"log"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

const (
	ScreenWidth  int = 320
	ScreenHeight int = 320
	blockWidth   int = (ScreenWidth - 4) / 3
	blockHeight  int = (ScreenHeight - 4) / 3

	fontSize float64 = 36

	xCenter int = blockHeight/2 - int(fontSize/4)
	yCenter int = blockWidth/2 + int(fontSize/4)

	gameoverStr = "      GAME OVER\n\nPRESE SPACE TO RESTART\n\n      ESC TO END"
)

// Player type
type Player int

// TicTacToe game type
type TicTacToe struct {
	Blocks     [3][3]Player
	IsGameOver bool
	BlocksFill int
}

var (
	white = color.White
	black = color.Black

	player1 Player = 1
	player2 Player = -1
	non     Player = 0

	x string = "X"
	o string = "O"

	ttt        = &TicTacToe{}
	noramlFont font.Face
	// currentPlayer started from player 1 X
	currentPlayer = player1

	imageGameover *ebiten.Image
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

	// Gameover
	imageGameover, _ = ebiten.NewImage(ScreenWidth, ScreenHeight, ebiten.FilterDefault)
	imageGameover.Fill(color.NRGBA{0x00, 0x00, 0x00, 0x80})
	y := (ScreenHeight - blockHeight) / 2
	drawTextWithShadowCenter(imageGameover, gameoverStr, 0, y, 1, color.White, ScreenWidth)
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.RunGame(ttt)
}

// Update updates a game by one tick. The given argument represents a screen image.
func (ttt *TicTacToe) Update(screen *ebiten.Image) error {
	return ttt.Draw(screen)
}

// Draw draw game screen
func (ttt *TicTacToe) Draw(screen *ebiten.Image) error {
	screen.Clear()
	screen.Fill(white)
	ttt.putBlock(screen)

	if ttt.IsGameOver {
		screen.DrawImage(imageGameover, &ebiten.DrawImageOptions{})
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			ttt.reset()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return errors.New("END")
		}
		return nil
	}

	posX, posY := ebiten.CursorPosition()
	x, y := mousePosTotXY(posX, posY)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && cursorInScreen(posX, posY) {
		ttt.setBlock(currentPlayer, x, y)
	}

	return nil
}

// Layout accepts a native outside size in device-independent pixels and returns the game's logical screen
// size.
func (ttt *TicTacToe) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// putBlock put block on game screen
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

// setBlock set block to "X" or "O"
func (ttt *TicTacToe) setBlock(player Player, x, y int) {
	if ttt.isFill(x, y) {
		return
	}
	ttt.Blocks[x][y] = player
	ttt.BlocksFill++
	if ttt.BlocksFill == 9 {
		ttt.IsGameOver = true
	}
	if checkWin(currentPlayer, ttt.Blocks) {
		ttt.IsGameOver = true
	}
	currentPlayer *= -1
}

func (ttt *TicTacToe) isFill(x, y int) bool {
	return ttt.Blocks[x][y] != 0
}

func (ttt *TicTacToe) reset() {
	for i := range ttt.Blocks {
		for j := range ttt.Blocks[i] {
			ttt.Blocks[i][j] = 0
		}
	}
	ttt.IsGameOver = false
	ttt.BlocksFill = 0
}

// newBlock return a new block
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

// cursorInScreen check if curcor is in the game screen
func cursorInScreen(x, y int) bool {
	return (x > 0 && y > 0) && (x < 320 && y < 320)
}

func mousePosTotXY(posX, posY int) (int, int) {
	return posY / blockWidth, posX / blockHeight
}

func checkWin(p Player, board [3][3]Player) bool {
	if checkRow(p, board) {
		return true
	}
	if checkCol(p, board) {
		return true
	}
	return checkDia(p, board)
}

func checkRow(p Player, board [3][3]Player) bool {
	for _, row := range board {
		if row[0] == p && row[1] == p && row[2] == p {
			return true
		}
	}
	return false
}

func checkCol(p Player, board [3][3]Player) bool {
	for i := 0; i < 3; i++ {
		if board[0][i] == p && board[1][i] == p && board[2][i] == p {
			return true
		}
	}
	return false
}

func checkDia(p Player, board [3][3]Player) bool {
	// check main diag
	if board[0][0] == p && board[1][1] == p && board[2][2] == p {
		return true
	}

	// check anit-diag
	if board[0][2] == p && board[1][1] == p && board[2][0] == p {
		return true
	}
	return false
}

// below's function is from ebiten/examples/blocks
// used to set game over iamge in this package
const (
	arcadeFontBaseSize = 8
)

var (
	arcadeFonts map[int]font.Face
)

func getArcadeFonts(scale int) font.Face {
	if arcadeFonts == nil {
		tt, err := truetype.Parse(fonts.ArcadeN_ttf)
		if err != nil {
			log.Fatal(err)
		}

		arcadeFonts = map[int]font.Face{}
		for i := 1; i <= 4; i++ {
			const dpi = 72
			arcadeFonts[i] = truetype.NewFace(tt, &truetype.Options{
				Size:    float64(arcadeFontBaseSize * i),
				DPI:     dpi,
				Hinting: font.HintingFull,
			})
		}
	}
	return arcadeFonts[scale]
}

func textWidth(str string) int {
	maxW := 0
	for _, line := range strings.Split(str, "\n") {
		b, _ := font.BoundString(getArcadeFonts(1), line)
		w := (b.Max.X - b.Min.X).Ceil()
		if maxW < w {
			maxW = w
		}
	}
	return maxW
}

var (
	shadowColor = color.NRGBA{0, 0, 0, 0x80}
)

func drawTextWithShadow(rt *ebiten.Image, str string, x, y, scale int, clr color.Color) {
	offsetY := arcadeFontBaseSize * scale
	for _, line := range strings.Split(str, "\n") {
		y += offsetY
		text.Draw(rt, line, getArcadeFonts(scale), x+1, y+1, shadowColor)
		text.Draw(rt, line, getArcadeFonts(scale), x, y, clr)
	}
}

func drawTextWithShadowCenter(rt *ebiten.Image, str string, x, y, scale int, clr color.Color, width int) {
	w := textWidth(str) * scale
	x += (width - w) / 2
	drawTextWithShadow(rt, str, x, y, scale, clr)
}

func drawTextWithShadowRight(rt *ebiten.Image, str string, x, y, scale int, clr color.Color, width int) {
	w := textWidth(str) * scale
	x += width - w
	drawTextWithShadow(rt, str, x, y, scale, clr)
}
