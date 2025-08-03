package game

type Board struct {
	Cells [15][15]string
}

func NewBoard() *Board {
	return &Board{}
}

// CheckWin kiểm tra xem người chơi có thắng không
func CheckWin(b *Board, x, y int, symbol string) bool {
	dirs := [][]int{
		{1, 0},  // Hàng ngang
		{0, 1},  // Hàng dọc
		{1, 1},  // Đường chéo xuống phải
		{1, -1}, // Đường chéo xuống trái
	}

	for _, dir := range dirs {
		count := 1

		// Đếm về phía trước
		for i := 1; i < 5; i++ {
			nx, ny := x+i*dir[0], y+i*dir[1]
			if inBounds(nx, ny) && b.Cells[nx][ny] == symbol {
				count++
			} else {
				break
			}
		}

		// Đếm về phía sau
		for i := 1; i < 5; i++ {
			nx, ny := x-i*dir[0], y-i*dir[1]
			if inBounds(nx, ny) && b.Cells[nx][ny] == symbol {
				count++
			} else {
				break
			}
		}

		if count >= 5 {
			return true
		}
	}
	return false
}

// inBounds kiểm tra xem tọa độ có trong bàn cờ không
func inBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < 15 && y < 15
}

// IsValidMove kiểm tra xem nước đi có hợp lệ không
func IsValidMove(b *Board, x, y int) bool {
	return inBounds(x, y) && b.Cells[x][y] == ""
}

// GetBoardState trả về trạng thái hiện tại của bàn cờ
func GetBoardState(b *Board) [15][15]string {
	return b.Cells
}

// ResetBoard làm mới bàn cờ
func ResetBoard(b *Board) {
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			b.Cells[i][j] = ""
		}
	}
}

// CountEmptyCells đếm số ô trống còn lại
func CountEmptyCells(b *Board) int {
	count := 0
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if b.Cells[i][j] == "" {
				count++
			}
		}
	}
	return count
}
