// Game state
let currentPlayer = "";
let myPlayer = "";
let gameBoard = [];
let gameActive = true;
let gameMode = ""; // "two_player" or "vs_computer"
const size = 15;
const cells = [];
let lastMove = null; // Track the last move

// DOM elements
const modeSelection = document.getElementById("mode-selection");
const gameInterface = document.getElementById("game-interface");
const twoPlayerBtn = document.getElementById("two-player-btn");
const vsComputerBtn = document.getElementById("vs-computer-btn");
const board = document.getElementById("board");
const status = document.getElementById("status");
const playerSymbol = document.getElementById("player-symbol");
const resetBtn = document.getElementById("reset-btn");
const backToMenuBtn = document.getElementById("back-to-menu-btn");
const gameModeInfo = document.getElementById("game-mode-info");

// WebSocket connection
let socket = null;

// Initialize game board array
function initializeBoard() {
  for (let i = 0; i < size; i++) {
    gameBoard[i] = [];
    for (let j = 0; j < size; j++) {
      gameBoard[i][j] = "";
    }
  }
}

// Create visual board
function createBoard() {
  board.innerHTML = "";
  for (let row = 0; row < size; row++) {
    cells[row] = [];
    for (let col = 0; col < size; col++) {
      const cell = document.createElement("div");
      cell.classList.add("cell");
      cell.dataset.row = row;
      cell.dataset.col = col;
      cell.textContent = "";
      cell.addEventListener("click", () => handleCellClick(row, col, cell));
      board.appendChild(cell);
      cells[row][col] = cell;
    }
  }
}

// Handle cell click
function handleCellClick(row, col, cell) {
  if (!gameActive || cell.textContent !== "" || currentPlayer !== myPlayer) {
    return;
  }

  const move = {
    row: row,
    col: col,
  };

  socket.send(JSON.stringify(move));
}

// Update cell display
function updateCell(row, col, player) {
  const cell = cells[row][col];
  cell.textContent = player;
  cell.classList.add(player.toLowerCase());
  gameBoard[row][col] = player;
}

// Start game with selected mode
function startGame(mode) {
  gameMode = mode;
  modeSelection.style.display = "none";
  gameInterface.style.display = "block";

  // Initialize game
  initializeBoard();
  createBoard();

  // Update game mode info
  if (mode === "two_player") {
    gameModeInfo.textContent = "Chế độ: Chơi 2 người";
  } else {
    gameModeInfo.textContent = "Chế độ: Chơi với máy";
  }

  // Connect to WebSocket
  socket = new WebSocket("ws://localhost:8080/ws");
  setupWebSocketHandlers();
}

// Setup WebSocket event handlers
function setupWebSocketHandlers() {
  socket.onopen = function () {
    status.textContent = "Đã kết nối với server";
    status.style.color = "#28a745";

    // Send game mode to server
    socket.send(JSON.stringify({ type: "game_mode", mode: gameMode }));
  };

  socket.onclose = function () {
    status.textContent = "Mất kết nối với server";
    status.style.color = "#dc3545";
  };

  socket.onerror = function () {
    status.textContent = "Lỗi kết nối";
    status.style.color = "#dc3545";
  };

  socket.onmessage = function (event) {
    const data = JSON.parse(event.data);

    if (data.type === "win") {
      // Handle win message
      handleWin(data);
    } else if (data.type === "move") {
      // Handle move message
      handleMove(data);
    } else if (data.type === "player_assignment") {
      // Handle player assignment
      handlePlayerAssignment(data);
    } else if (data.type === "turn_update") {
      // Handle turn update
      handleTurnUpdate(data);
    } else if (data.type === "reset") {
      // Handle reset message
      handleReset();
    } else {
      // Handle regular move (backward compatibility)
      handleMove(data);
    }
  };
}

// Handle player assignment
function handlePlayerAssignment(data) {
  myPlayer = data.player;
  status.textContent = `Bạn là người chơi: ${myPlayer}`;
}

// Handle turn update
function handleTurnUpdate(data) {
  currentPlayer = data.current_player;
  playerSymbol.textContent = currentPlayer;
}

// Handle move
function handleMove(data) {
  if (data.row !== undefined && data.col !== undefined && data.player) {
    // Clear previous last move highlight
    if (lastMove) {
      const prevCell = cells[lastMove.row][lastMove.col];
      prevCell.classList.remove("last-move");
      // Reset color if not winning
      if (!prevCell.classList.contains("win")) {
        prevCell.style.color = "";
        prevCell.style.textShadow = "";
      }
    }

    updateCell(data.row, data.col, data.player);

    // Highlight the new last move
    if (data.last_move) {
      const cell = cells[data.last_move.row][data.last_move.col];
      cell.classList.add("last-move");
      // Ensure text is visible
      cell.style.color = "white";
      cell.style.textShadow = "1px 1px 2px rgba(0, 0, 0, 0.5)";
      lastMove = { row: data.last_move.row, col: data.last_move.col };
    }
  }
}

// Handle win
function handleWin(data) {
  gameActive = false;

  // First update the cell with the player symbol
  if (data.row !== undefined && data.col !== undefined && data.player) {
    updateCell(data.row, data.col, data.player);
  }

  // Then highlight winning cells
  if (data.winning_cells) {
    data.winning_cells.forEach((cell) => {
      const cellElement = cells[cell.row][cell.col];
      cellElement.classList.add("win");
      // Ensure text is visible
      cellElement.style.color = "white";
      cellElement.style.textShadow = "1px 1px 2px rgba(0, 0, 0, 0.5)";
    });
  }

  // Show win message
  showGameOver(`Người chơi ${data.player} đã thắng!`);
}

// Show game over dialog
function showGameOver(message) {
  const overlay = document.createElement("div");
  overlay.style.position = "fixed";
  overlay.style.top = "0";
  overlay.style.left = "0";
  overlay.style.width = "100%";
  overlay.style.height = "100%";
  overlay.style.backgroundColor = "rgba(0,0,0,0.5)";
  overlay.style.zIndex = "999";

  const dialog = document.createElement("div");
  dialog.className = "game-over";
  dialog.innerHTML = `
    <h2>${message}</h2>
    <button onclick="resetGame()">Chơi lại</button>
    <button onclick="closeDialog()">Đóng</button>
  `;

  overlay.appendChild(dialog);
  document.body.appendChild(overlay);
}

// Close dialog
function closeDialog() {
  const overlay = document.querySelector('div[style*="position: fixed"]');
  if (overlay) {
    overlay.remove();
  }
}

// Reset game
function resetGame() {
  // Clear board
  for (let row = 0; row < size; row++) {
    for (let col = 0; col < size; col++) {
      cells[row][col].textContent = "";
      cells[row][col].className = "cell";
      cells[row][col].style.color = "";
      cells[row][col].style.textShadow = "";
      gameBoard[row][col] = "";
    }
  }

  gameActive = true;
  lastMove = null; // Reset last move
  closeDialog();

  // Send reset request to server
  socket.send(JSON.stringify({ type: "reset" }));
}

// Handle reset from server
function handleReset() {
  // Clear board
  for (let row = 0; row < size; row++) {
    for (let col = 0; col < size; col++) {
      cells[row][col].textContent = "";
      cells[row][col].className = "cell";
      cells[row][col].style.color = "";
      cells[row][col].style.textShadow = "";
      gameBoard[row][col] = "";
    }
  }

  gameActive = true;
  lastMove = null; // Reset last move
  closeDialog(); // Close any open dialogs
}

// Back to menu
function backToMenu() {
  if (socket) {
    socket.close();
  }
  gameInterface.style.display = "none";
  modeSelection.style.display = "block";
  gameMode = "";
  gameActive = true;
  currentPlayer = "";
  myPlayer = "";
}

// Event listeners
twoPlayerBtn.addEventListener("click", () => startGame("two_player"));
vsComputerBtn.addEventListener("click", () => startGame("vs_computer"));
resetBtn.addEventListener("click", resetGame);
backToMenuBtn.addEventListener("click", backToMenu);
