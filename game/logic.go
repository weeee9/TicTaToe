package game

func blockIsFill(board [3][3]player, x, y int) bool {
	return board[x][y] != 0
}

func cursorInScreen(x, y int) bool {
	return (x > 0 && y > 0) && (x < 320 && y < 320)
}

func mousePosTotXY(posX, posY int) (int, int) {
	return posY / blockWidth, posX / blockHeight
}

func checkWin(p player, board [3][3]player) bool {
	if checkRow(p, board) {
		return true
	}
	if checkCol(p, board) {
		return true
	}
	return checkDia(p, board)
}

func checkRow(p player, board [3][3]player) bool {
	for _, row := range board {
		if row[0] == p && row[1] == p && row[2] == p {
			return true
		}
	}
	return false
}

func checkCol(p player, board [3][3]player) bool {
	for i := 0; i < 3; i++ {
		if board[0][i] == p && board[1][i] == p && board[2][i] == p {
			return true
		}
	}
	return false
}

func checkDia(p player, board [3][3]player) bool {
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
