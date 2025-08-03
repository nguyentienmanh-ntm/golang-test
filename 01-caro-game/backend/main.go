package main

import (
	"caro-game/backend/game"
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Cho phép tất cả origin
	},
}

type GameState struct {
	Board        *game.Board
	CurrentPlayer string
	Players      map[*websocket.Conn]string
	GameActive   bool
	GameMode     string // "two_player" or "vs_computer"
}

type Move struct {
	Row    int    `json:"row"`
	Col    int    `json:"col"`
	Player string `json:"player"`
}

type Message struct {
	Type         string      `json:"type"`
	Row          int         `json:"row,omitempty"`
	Col          int         `json:"col,omitempty"`
	Player       string      `json:"player,omitempty"`
	CurrentPlayer string     `json:"current_player,omitempty"`
	WinningCells []Cell      `json:"winning_cells,omitempty"`
	Mode         string      `json:"mode,omitempty"`
	LastMove     *Cell       `json:"last_move,omitempty"`
}

type Cell struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

var (
	clients []*websocket.Conn
	mu      sync.Mutex
	gameState = &GameState{
		Board:        game.NewBoard(),
		CurrentPlayer: "X",
		Players:      make(map[*websocket.Conn]string),
		GameActive:   true,
		GameMode:     "two_player",
	}
)

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
	
	app := iris.New()

	// Serve static files
	app.HandleDir("/", iris.Dir("./frontend"))

	app.Get("/ws", handleWebSocket)

	app.Listen(":8080")
}

func handleWebSocket(ctx iris.Context) {
	conn, err := upgrader.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)
	if err != nil {
		return
	}

	mu.Lock()
	clients = append(clients, conn)
	mu.Unlock()

	// Handle disconnection
	defer func() {
		mu.Lock()
		removeClient(conn)
		mu.Unlock()
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var data map[string]interface{}
		json.Unmarshal(msg, &data)

		if data["type"] == "reset" {
			handleReset()
		} else if data["type"] == "game_mode" {
			handleGameMode(conn, data)
		} else {
			handleMove(conn, msg)
		}
	}
}

func handleGameMode(conn *websocket.Conn, data map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()

	mode, ok := data["mode"].(string)
	if !ok {
		return
	}

	if gameState.GameMode != mode {
		gameState.Board = game.NewBoard()
		gameState.CurrentPlayer = "X"
		gameState.GameActive = true
		gameState.Players = make(map[*websocket.Conn]string)
		gameState.GameMode = mode
	}

	if mode == "vs_computer" {
		// Single player mode - assign player as X
		gameState.Players[conn] = "X"
		sendMessage(conn, Message{
			Type:   "player_assignment",
			Player: "X",
		})
		
		// Send current turn info
		sendMessage(conn, Message{
			Type:         "turn_update",
			CurrentPlayer: gameState.CurrentPlayer,
		})
	} else {
		// Two player mode - assign player symbol
		if len(gameState.Players) == 0 {
			gameState.Players[conn] = "X"
			sendMessage(conn, Message{
				Type:   "player_assignment",
				Player: "X",
			})
		} else if len(gameState.Players) == 1 {
			gameState.Players[conn] = "O"
			sendMessage(conn, Message{
				Type:   "player_assignment",
				Player: "O",
			})
		} else {
			// Too many players
			conn.Close()
			return
		}
		
		// Send current turn info
		sendMessage(conn, Message{
			Type:         "turn_update",
			CurrentPlayer: gameState.CurrentPlayer,
		})
	}
}

func handleMove(conn *websocket.Conn, msg []byte) {
	mu.Lock()
	defer mu.Unlock()

	var move Move
	json.Unmarshal(msg, &move)

	// Check if it's the player's turn
	player, exists := gameState.Players[conn]
	if !exists || player != gameState.CurrentPlayer || !gameState.GameActive {
		return
	}

	// Check if cell is empty
	if !game.IsValidMove(gameState.Board, move.Row, move.Col) {
		return
	}

	// Set the player for this move based on current player
	move.Player = gameState.CurrentPlayer

	// Update board with the current player's symbol
	gameState.Board.Cells[move.Row][move.Col] = gameState.CurrentPlayer

	// Check for win using the current player's symbol
	if game.CheckWin(gameState.Board, move.Row, move.Col, gameState.CurrentPlayer) {
		gameState.GameActive = false
		
		// Find winning cells
		winningCells := findWinningCells(move.Row, move.Col, gameState.CurrentPlayer)
		
		// Broadcast win message
		winMsg := Message{
			Type:         "win",
			Row:          move.Row,    
			Col:          move.Col,  
			Player:       gameState.CurrentPlayer,
			WinningCells: winningCells,
		}
		broadcastMessage(winMsg)
	} else {
		// Switch player
		gameState.CurrentPlayer = togglePlayer(gameState.CurrentPlayer)
		
		// Broadcast move
		moveMsg := Message{
			Type:   "move",
			Row:    move.Row,
			Col:    move.Col,
			Player: move.Player,
			LastMove: &Cell{Row: move.Row, Col: move.Col},
		}
		broadcastMessage(moveMsg)
		
		// Broadcast turn update
		turnMsg := Message{
			Type:         "turn_update",
			CurrentPlayer: gameState.CurrentPlayer,
		}
		broadcastMessage(turnMsg)
		
		// If playing against computer and it's computer's turn
		if gameState.GameMode == "vs_computer" && gameState.CurrentPlayer == "O" && gameState.GameActive {
			// Make computer move after a short delay
			go func() {
				time.Sleep(500 * time.Millisecond)
				makeComputerMove()
			}()
		}
	}
}

func makeComputerMove() {
	mu.Lock()
	defer mu.Unlock()

	if !gameState.GameActive || gameState.CurrentPlayer != "O" {
		return
	}

	// Find all empty cells
	var emptyCells []Cell
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if gameState.Board.Cells[i][j] == "" {
				emptyCells = append(emptyCells, Cell{Row: i, Col: j})
			}
		}
	}

	if len(emptyCells) == 0 {
		return
	}

	// Choose random empty cell
	randomIndex := rand.Intn(len(emptyCells))
	computerMove := emptyCells[randomIndex]

	// Make computer move
	gameState.Board.Cells[computerMove.Row][computerMove.Col] = "O"

	// Check for win
	if game.CheckWin(gameState.Board, computerMove.Row, computerMove.Col, "O") {
		gameState.GameActive = false
		
		// Find winning cells
		winningCells := findWinningCells(computerMove.Row, computerMove.Col, "O")
		
		// Broadcast win message
		winMsg := Message{
			Type:         "win",
			Player:       "O",
			Row:          computerMove.Row,
			Col:          computerMove.Col,
			WinningCells: winningCells,
		}
		broadcastMessage(winMsg)
	} else {
		// Switch player
		gameState.CurrentPlayer = "X"
		
		// Broadcast computer move
		moveMsg := Message{
			Type:   "move",
			Row:    computerMove.Row,
			Col:    computerMove.Col,
			Player: "O",
			LastMove: &Cell{Row: computerMove.Row, Col: computerMove.Col},
		}
		broadcastMessage(moveMsg)
		
		// Broadcast turn update
		turnMsg := Message{
			Type:         "turn_update",
			CurrentPlayer: gameState.CurrentPlayer,
		}
		broadcastMessage(turnMsg)
	}
}

func handleReset() {
	mu.Lock()
	defer mu.Unlock()

	// Reset game state
	gameState.Board = game.NewBoard()
	gameState.CurrentPlayer = "X"
	gameState.GameActive = true

	// Broadcast reset message
	resetMsg := Message{
		Type: "reset",
	}
	broadcastMessage(resetMsg)
	
	// Broadcast turn update
	turnMsg := Message{
		Type:         "turn_update",
		CurrentPlayer: gameState.CurrentPlayer,
	}
	broadcastMessage(turnMsg)
}

func findWinningCells(row, col int, player string) []Cell {
	dirs := [][]int{
		{1, 0}, {0, 1}, {1, 1}, {1, -1},
	}

	for _, dir := range dirs {
		count := 1
		cells := []Cell{{Row: row, Col: col}}

		// Check forward direction
		for i := 1; i < 5; i++ {
			nx, ny := row+i*dir[0], col+i*dir[1]
			if inBounds(nx, ny) && gameState.Board.Cells[nx][ny] == player {
				count++
				cells = append(cells, Cell{Row: nx, Col: ny})
			} else {
				break
			}
		}

		// Check backward direction
		for i := 1; i < 5; i++ {
			nx, ny := row-i*dir[0], col-i*dir[1]
			if inBounds(nx, ny) && gameState.Board.Cells[nx][ny] == player {
				count++
				cells = append(cells, Cell{Row: nx, Col: ny})
			} else {
				break
			}
		}

		if count >= 5 {
			return cells
		}
	}
	return nil
}

func inBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < 15 && y < 15
}

func broadcastMessage(message Message) {
	msg, _ := json.Marshal(message)
	for _, conn := range clients {
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func sendMessage(conn *websocket.Conn, message Message) {
	msg, _ := json.Marshal(message)
	conn.WriteMessage(websocket.TextMessage, msg)
}

func removeClient(conn *websocket.Conn) {
	for i, client := range clients {
		if client == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	delete(gameState.Players, conn)
}

func togglePlayer(p string) string {
	if p == "X" {
		return "O"
	}
	return "X"
}
