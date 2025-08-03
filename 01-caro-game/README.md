# Game Caro - 15x15

Một game caro hoàn chỉnh được xây dựng bằng Go (backend) và HTML/CSS/JavaScript (frontend) với WebSocket để giao tiếp real-time.

## Tính năng

### 1. Bàn cờ 15x15

- Bàn cờ caro kích thước 15x15 ô
- Giao diện đẹp mắt với CSS Grid
- Hiệu ứng hover và transition mượt mà

### 2. Chế độ chơi

- **Chơi 2 người**: Hai người chơi thay phiên nhau
- **Chơi với máy**: Người chơi đánh với AI (đánh ngẫu nhiên)
- Giao diện chọn chế độ chơi khi bắt đầu

### 3. Quản lý người chơi

- Hỗ trợ 2 người chơi (X và O) trong chế độ 2 người
- Tự động phân bổ ký hiệu người chơi khi kết nối
- Hiển thị lượt chơi hiện tại
- AI đánh ngẫu nhiên trong chế độ chơi với máy

### 4. Logic game

- Kiểm tra thắng thua theo 3 hướng: ngang, dọc, chéo
- Yêu cầu 5 ký tự liên tiếp để thắng
- Ngăn chặn nước đi không hợp lệ
- Highlight các ô thắng cuộc
- **Highlight nước đi vừa thực hiện**: Ô vừa đánh sẽ được highlight màu xanh để dễ nhận biết

### 5. WebSocket Communication

- Kết nối real-time giữa các người chơi
- Đồng bộ trạng thái game
- Thông báo thắng thua tức thì

### 6. Giao diện người dùng

- Dialog thông báo kết quả
- Nút "Chơi lại" để reset game
- Nút "Về menu" để quay lại chọn chế độ

## Cách chạy

### Yêu cầu

- Go 1.24.5 hoặc cao hơn
- Trình duyệt web hiện đại

### Bước 1: Chạy server

```bash
cd 01-caro-game
go run backend/main.go (phải chạy đúng như này)
```
Server sẽ chạy tại `http://localhost:8080`

### Bước 2: Mở game

1. Mở trình duyệt và truy cập `http://localhost:8080`
2. Chọn chế độ chơi:
   - **Chơi 2 người**: Mở thêm một tab/cửa sổ trình duyệt khác với cùng URL
   - **Chơi với máy**: Chỉ cần một tab trình duyệt

## Cấu trúc project

```
01-caro-game/
├── backend/
│   ├── game/
│   │   └── game.go          # Logic game caro
│   └── main.go              # Server WebSocket + AI
├── frontend/
│   ├── index.html           # Giao diện chính
│   ├── style.css            # CSS styling
│   └── app.js               # JavaScript logic
├── go.mod                   # Go modules
└── README.md               # Hướng dẫn này
```

## Luật chơi

1. **Bắt đầu**: Người chơi X đi trước
2. **Lượt chơi**: Hai người chơi thay phiên nhau (hoặc người chơi với máy)
3. **Mục tiêu**: Tạo được 5 ký tự liên tiếp theo hàng ngang, dọc hoặc chéo
4. **Thắng cuộc**: Người đầu tiên tạo được 5 ký tự liên tiếp sẽ thắng
5. **Chơi lại**: Nhấn nút "Chơi lại" để bắt đầu game mới
6. **Về menu**: Nhấn nút "Về menu" để chọn lại chế độ chơi

## Chế độ chơi

### Chơi 2 người

- Cần 2 người chơi kết nối vào server
- Người chơi đầu tiên sẽ là X, người thứ hai là O
- Hai người thay phiên nhau đánh

### Chơi với máy

- Chỉ cần 1 người chơi
- Người chơi luôn là X, máy là O
- Máy sẽ đánh ngẫu nhiên sau khi người chơi đánh
- Có độ trễ 500ms để tạo cảm giác tự nhiên

## Công nghệ sử dụng

### Backend

- **Go**: Ngôn ngữ lập trình chính
- **Iris Framework**: Web framework
- **Gorilla WebSocket**: WebSocket implementation
- **Goroutines**: Xử lý đồng thời
- **Random AI**: AI đánh ngẫu nhiên

### Frontend

- **HTML5**: Cấu trúc trang web
- **CSS3**: Styling và layout
- **JavaScript ES6+**: Logic client-side
- **WebSocket API**: Giao tiếp real-time

## Tính năng kỹ thuật

### Backend

- **Concurrent handling**: Xử lý nhiều kết nối đồng thời
- **State management**: Quản lý trạng thái game tập trung
- **Validation**: Kiểm tra tính hợp lệ của nước đi
- **Win detection**: Thuật toán kiểm tra thắng thua hiệu quả
- **AI implementation**: AI đánh ngẫu nhiên với goroutines

### Frontend

- **Real-time updates**: Cập nhật giao diện tức thì
- **Event handling**: Xử lý sự kiện người dùng
- **State synchronization**: Đồng bộ trạng thái với server
- **Error handling**: Xử lý lỗi kết nối và game
- **Mode selection**: Giao diện chọn chế độ chơi

## Mở rộng

Game có thể được mở rộng với các tính năng:

- Chat giữa người chơi
- Lưu lịch sử game
- AI thông minh hơn (minimax, alpha-beta pruning)
- Tournament mode
- Custom board sizes
- Spectator mode
- Multiple AI difficulty levels

## Troubleshooting

### Lỗi kết nối

- Kiểm tra server có đang chạy không
- Đảm bảo port 8080 không bị chiếm
- Kiểm tra firewall settings

### Game không hoạt động

- Refresh trang web
- Kiểm tra console browser để xem lỗi
- Đảm bảo WebSocket được hỗ trợ

### Chế độ 2 người

- Chỉ hỗ trợ tối đa 2 người chơi
- Người chơi thứ 3 sẽ bị từ chối kết nối

### Chế độ chơi với máy

- Chỉ cần 1 người chơi
- Máy sẽ tự động đánh sau lượt của người chơi
- Có thể chơi nhiều lần mà không cần reset server
