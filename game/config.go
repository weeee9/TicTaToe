package game

const (
	screenWidth  int = 320
	screenHeight int = 320

	blockWidth  int = (screenWidth - 4) / 3
	blockHeight int = (screenHeight - 4) / 3

	fontSize float64 = 36

	xCenter int = blockHeight/2 - int(fontSize/4)
	yCenter int = blockWidth/2 + int(fontSize/4)

	gameoverStr = "      GAME OVER\n\nPRESE SPACE TO RESTART\n\n      ESC TO END"
)
