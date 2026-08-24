# GameMate — Nền Tảng Chat & Voice/Video Call Tìm Đồng Đội Chơi Game

## 1. Tổng quan

**GameMate** là nền tảng giúp game thủ tìm đồng đội chơi cùng, giao tiếp qua chat và voice/video call theo thời gian thực, tổ chức theo **party/room** giống mô hình Discord nhưng tập trung vào tính năng **matchmaking tìm đồng đội theo tiêu chí** (game, rank, khung giờ chơi, ngôn ngữ, mục tiêu chơi — giải trí hay tryhard rank).

Dự án được xây dựng theo kiến trúc **modular monolith** bằng Go (backend) và Next.js (frontend), tập trung thể hiện năng lực xử lý hệ thống real-time, matchmaking logic, và giao tiếp WebRTC.

**Tech Stack:**

| Thành phần | Công nghệ |
|---|---|
| Backend | Go (Gin/Echo/Fiber) |
| Frontend | Next.js (App Router, TypeScript) |
| Real-time | WebSocket (gorilla/websocket) |
| Voice/Video | WebRTC, Pion (Go WebRTC library) |
| Message Broker | Redis Pub/Sub |
| Database | PostgreSQL |
| Object Storage | S3 / MinIO |
| STUN/TURN | coturn |
| Containerization | Docker, Docker Compose |
| CI/CD | GitHub Actions |
| Monitoring | Prometheus + Grafana |
| Deployment | VPS / Railway / Fly.io (backend), Vercel (frontend) |

---

## 2. Ý tưởng sản phẩm

### Vấn đề
Game thủ (đặc biệt chơi game team-based: LMHT, Valorant, CS2, PUBG...) thường gặp khó khăn khi muốn tìm đồng đội "hợp gu" — cùng rank, cùng mục tiêu (giải trí/rank), cùng khung giờ rảnh — thay vì random với người lạ trong queue của game. Discord server thì quá rộng, không có cơ chế match tự động; group chat thường thì không biết ai đang rảnh, đang chơi game gì.

### Giải pháp
GameMate cho phép:
- Tạo **hồ sơ gamer**: game đang chơi, rank/level, khung giờ rảnh, mục tiêu (chill/tryhard), ngôn ngữ giao tiếp
- **Matchmaking tìm đồng đội**: tìm người theo tiêu chí, hoặc "Quick Match" ghép ngẫu nhiên trong tiêu chí phù hợp
- **Party Room**: tạo phòng chat + voice call cho nhóm nhỏ (giống Discord voice channel nhưng nhẹ hơn), có thể public (mở cho ai muốn join) hoặc private (mời link)
- **Push-to-talk & voice activity** trong lúc chơi
- **Group video call** khi cần xem màn hình nhau bàn chiến thuật (screen share optional)
- Lịch sử các buổi chơi cùng nhau, đánh giá đồng đội (rate teammate) để xây dựng uy tín trong cộng đồng

---

## 3. Tính năng chính

### Hồ sơ & Quản lý Profile
- Đăng ký/đăng nhập, xác thực JWT (access + refresh token)
- Người dùng chỉ cần thiết lập hồ sơ (`game_profiles`) **lần đầu tiên**; các lần đăng nhập sau hệ thống tự động tải Profile từ PostgreSQL
- Hồ sơ gamer: game yêu thích, rank, khung giờ rảnh, mục tiêu chơi (chill/tryhard), ngôn ngữ, mô tả ngắn
- **Đa dạng Game & Rank**: Mỗi game có hệ thống rank riêng (LMHT: Sắt → Thách Đấu, Valorant: Iron → Radiant). Rank được lưu trữ linh hoạt dạng **JSONB** trong DB, sử dụng file Config ở Frontend để map và hiển thị đúng rank theo từng game
- Gửi lời mời kết bạn / mời vào party

### Matchmaking & Tìm Đồng Đội
- Tìm đồng đội theo bộ lọc (game, rank range, giờ rảnh, mục tiêu)
- **Quick Match**: ghép ngẫu nhiên người phù hợp tiêu chí đang chờ trong hàng đợi
- **Group Matchmaking**: Hỗ trợ ghép nhóm linh hoạt 2, 3, 4, 5 người — không chỉ ghép cặp 1-1
- **Slot Filling**: Thuật toán Redis Worker quét hàng đợi để ghép linh hoạt các nhóm (ví dụ: nhóm 3 + nhóm 2, hoặc nhóm 3 + 1 + 1 để đủ team 5 người)
- **Bạn bè chơi chung**: Tạo Party Room trước → gửi link/mã mời hoặc thông báo real-time qua WebSocket → Trưởng phòng bấm tìm thêm số slot còn thiếu qua matchmaking

### Chat
- Chat 1-1 và group chat (trong party room) real-time qua WebSocket
- Typing indicator, read receipt, presence (online/offline/đang chơi game gì)
- Gửi ảnh/screenshot, emoji reaction
- Lịch sử chat với phân trang (cursor-based)

### Party Room (điểm nhấn khác biệt)
- Tạo room theo game cụ thể, giới hạn số lượng thành viên
- Room public (browse & join) hoặc private (link mời/mã mời)
- Voice channel trong room: bật/tắt mic, push-to-talk, hiển thị ai đang nói
- Video call khi cần (screen share để bàn chiến thuật)
- **Vòng đời phòng (Lifecycle)**: `waiting` (đang chờ/tìm người) → `in_game` (đang chơi) → `closed` (kết thúc)
- **Giải phóng tài nguyên**: Khi phòng đóng (tất cả thoát hoặc Host giải tán), hệ thống tự động ngắt stream WebRTC và kết nối WebSocket để giải phóng băng thông
- **Lưu vết dữ liệu**: Record phòng vẫn được giữ lại trong DB để phục vụ xem lịch sử chat, danh sách buổi chơi, và đánh giá đồng đội sau trận

### Voice/Video Call
- Signaling server tự xây dựng bằng Go (qua WebSocket)
- NAT traversal qua STUN/TURN (coturn)
- Voice call nhóm nhỏ và video call qua kiến trúc SFU (Pion)
- Voice activity detection (hiển thị ai đang nói bằng viền sáng quanh avatar)

### Cộng đồng & Uy tín (Teammate Rating & Anti-Abuse)
- Đánh giá đồng đội sau mỗi buổi chơi — **tự nguyện**, không bắt buộc
- Màn hình kết thúc luôn có nút "Bỏ qua" / tự động đóng sau 10-15s để không gây phiền
- **Giao diện 1-Click**: Dạng thẻ/icon chọn nhanh (`#friendly`, `#skilled`, `#toxic`, `#communicative`) cho 4 đồng đội, hoàn tất trong 2-3 giây
- Lịch sử các buổi chơi cùng nhau
- Badge/huy hiệu đơn giản dựa trên rating (optional, nice-to-have)
- **Chống đánh giá xấu đơn phương (Anti-Abuse)**:
  - *So sánh lịch sử (Anomaly Detection)*: Nếu B có lịch sử 5-star, 1 đánh giá 1-star bất ngờ từ A sẽ bị đánh nhãn `ANOMALY` và tạm thời không trừ điểm B
  - *Xác thực chéo (Cross-Check)*: Nếu 3 người khác đánh giá B tốt, chỉ có A đánh giá tệ → hệ thống tự động hủy bỏ đánh giá của A
  - *Phạt người đánh giá xấu (Shadow Penalty)*: Nếu A thường xuyên cố tình đánh giá 1-star cho người khác, trọng số đánh giá của A sẽ bị hạ về 0 (vô hiệu hóa)

### Hạ tầng & Vận hành
- Scale ngang nhiều instance backend nhờ Redis Pub/Sub
- Container hóa toàn bộ hệ thống, chạy bằng 1 lệnh `docker-compose up`
- CI/CD tự động build/test/deploy
- Metrics và dashboard giám sát hệ thống real-time (số room đang hoạt động, số user online, số cuộc gọi đồng thời)

---

## 4. Kiến trúc hệ thống

### 4.1 Sơ đồ tổng quan

```
                         ┌──────────────────┐
                         │   Next.js Web    │
                         │ (Client/Browser) │
                         └───────┬──────────┘
                                 │ HTTPS / WSS
                                 ▼
                       ┌─────────────────────┐
                       │   API Gateway /     │
                       │   Load Balancer     │
                       │   (Nginx / Traefik) │
                       └──────────┬──────────┘
     ┌────────────────┬───────────┼───────────────┬────────────────┐
     ▼                ▼           ▼               ▼                ▼
┌───────────┐  ┌───────────────┐ ┌───────────┐ ┌─────────────┐ ┌────────────────┐
│  Auth     │  │ Matchmaking   │ │  Chat     │ │ Party Room  │ │ Call/Signaling │
│  Module   │  │  Module       │ │  Module   │ │   Module    │ │     Module     │
└─────┬─────┘  └───────┬───────┘ └──────┬────┘ └───────┬─────┘ └────────┬───────┘
      │                │                │              │                │
      │                └───────┬────────┴──────┬───────┘                │
      │                        ▼               ▼                        │
      │              ┌───────────────────────────────┐                  │
      │              │         Redis                 │                  │
      │              │ (Pub/Sub fanout, matchmaking  │                  │
      │              │  queue, presence, room state) │                  │
      │              └───────────────┬───────────────┘                  │
      │                              │                                  │
      ▼                              ▼                                  ▼
┌─────────────────────────────────────────────┐                 ┌──────────────┐
│                PostgreSQL                   │                 │   coturn     │
│(users, profiles, rooms, messages, ratings..)│                 │ (STUN/TURN)  │
└─────────────────────────────────────────────┘                 └──────────────┘
      │
      ▼
┌─────────────┐
│ S3 / MinIO  │
│(avatar, ảnh)│
└─────────────┘
```

### 4.2 Luồng Matchmaking (điểm đặc trưng của project)

**Kịch bản 1: Quick Match (solo hoặc nhóm nhỏ tìm thêm người)**
1. User (hoặc nhóm bạn đã trong Party Room) bấm "Quick Match" với tiêu chí (game, rank range, khung giờ) → request được đẩy vào **matchmaking queue** lưu trong Redis (Sorted Set hoặc List theo từng "bucket" tiêu chí)
2. Một worker (goroutine chạy định kỳ) quét queue, sử dụng thuật toán **Slot Filling** để ghép linh hoạt: nhóm 3 + nhóm 2, nhóm 3 + 1 + 1, hoặc 5 solo players — miễn sao lắp đầy đủ team size theo yêu cầu game
3. Khi tìm được nhóm phù hợp → tạo Party Room tự động (hoặc merge vào room có sẵn nếu nhóm đã có phòng), gửi thông báo real-time qua WebSocket ("đã tìm thấy đồng đội, chấp nhận trong 30s")
4. Nếu tất cả chấp nhận → đưa vào room, khởi tạo chat + voice channel
5. Nếu timeout hoặc có người từ chối → trả các user/nhóm còn lại về queue

**Kịch bản 2: Bạn bè chơi chung**
1. Trưởng nhóm tạo Party Room trước, gửi link/mã mời hoặc thông báo real-time qua WebSocket cho bạn bè
2. Bạn bè join vào room qua link/mã mời
3. Trưởng phòng bấm "Tìm thêm" → đẩy request vào matchmaking queue với số slot còn thiếu (ví dụ: room 3/5 → cần thêm 2)
4. Worker ghép slot thiếu theo cùng cơ chế Slot Filling ở trên

> Đây là phần thể hiện tư duy thuật toán + xử lý concurrent tốt nhất trong dự án — nên đầu tư kỹ và chuẩn bị giải thích rõ trong phỏng vấn.

### 4.3 Vòng đời Party Room (Lifecycle)

1. **`waiting`**: Phòng mới tạo, đang chờ đủ người hoặc đang tìm thêm qua matchmaking
2. **`in_game`**: Đủ người, trưởng phòng bấm bắt đầu hoặc hệ thống tự chuyển khi đủ team
3. **`closed`**: Tất cả thoát, Host giải tán, hoặc buổi chơi kết thúc
4. Khi chuyển sang `closed` → hệ thống tự động **ngắt stream WebRTC** và **đóng kết nối WebSocket** để giải phóng băng thông server
5. Record phòng (messages, call logs, members) **vẫn giữ lại** trong PostgreSQL để phục vụ xem lịch sử và đánh giá đồng đội

### 4.4 Luồng xử lý Chat & Presence (WebSocket + Redis Pub/Sub)

1. Client mở kết nối WebSocket, xác thực bằng JWT trong handshake
2. Khi user gửi message trong room/DM → server lưu message vào PostgreSQL, publish lên Redis channel theo room/conversation ID
3. Mọi instance backend subscribe channel đó forward tới client đang kết nối với mình → đảm bảo hoạt động đúng khi scale ngang nhiều instance
4. Presence (online/offline/đang chơi game gì) cũng lưu trong Redis với TTL heartbeat, publish thay đổi trạng thái real-time cho bạn bè/room liên quan

### 4.5 Luồng xử lý Voice/Video Call (WebRTC Signaling)

1. Khi vào Party Room, client mở kết nối signaling qua WebSocket
2. Trao đổi SDP offer/answer và ICE candidate qua signaling server để thiết lập kết nối WebRTC
3. Với voice channel nhóm nhỏ (≤4 người): có thể dùng mesh network (mỗi client kết nối trực tiếp với nhau) để đơn giản hóa
4. Với video call/group lớn hơn: dùng **SFU** (Pion) — mỗi client chỉ upload 1 stream, SFU forward tới các client còn lại, tránh quá tải băng thông
5. STUN/TURN (coturn) xử lý NAT traversal khi client nằm sau firewall/NAT không thể kết nối P2P trực tiếp

### 4.6 Modular Monolith — ranh giới module

- **`auth`**: đăng ký, đăng nhập, JWT, quản lý session
- **`profile`**: hồ sơ gamer (game, rank lưu JSONB, giờ rảnh, mục tiêu chơi); auto-load profile khi đăng nhập lại
- **`matchmaking`**: quản lý queue, thuật toán ghép nhóm linh hoạt (slot filling), worker xử lý định kỳ
- **`chat`**: conversation, message, xử lý trong cả DM và party room
- **`party`**: quản lý room lifecycle (tạo, join, rời, trạng thái waiting → in_game → closed, giải phóng tài nguyên khi đóng)
- **`realtime`**: quản lý WebSocket connection, Redis Pub/Sub fanout, presence
- **`call`**: signaling, quản lý voice/video channel, tích hợp SFU, cleanup stream khi room closed
- **`rating`**: đánh giá đồng đội tự nguyện, anti-abuse (anomaly detection, cross-check, shadow penalty)
- **`notification`**: thông báo match tìm thấy, lời mời vào room, v.v.
- **`storage`**: upload/download avatar, ảnh qua S3/MinIO

---

## 5. Cấu trúc thư mục dự kiến

### Backend (Go)

```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   ├── profile/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go               # GameProfile: game, rank, playtime, goal
│   ├── matchmaking/
│   │   ├── handler.go
│   │   ├── service.go             # thuật toán ghép cặp/nhóm
│   │   ├── queue.go               # quản lý queue trong Redis
│   │   └── worker.go              # goroutine xử lý matching định kỳ
│   ├── chat/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   ├── party/
│   │   ├── handler.go             # tạo/join/leave room
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   ├── realtime/
│   │   ├── hub.go
│   │   ├── client.go
│   │   ├── pubsub.go
│   │   └── protocol.go
│   ├── call/
│   │   ├── signaling.go
│   │   ├── sfu.go
│   │   └── room.go
│   ├── rating/
│   │   ├── handler.go
│   │   └── service.go
│   ├── notification/
│   ├── storage/
│   └── middleware/
├── pkg/
│   ├── config/
│   ├── logger/
│   ├── database/
│   └── redis/
├── migrations/
├── deployments/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── k8s/
├── test/
├── go.mod
└── README.md
```

### Frontend (Next.js)

```
frontend/
├── app/
│   ├── (auth)/
│   │   ├── login/page.tsx
│   │   └── register/page.tsx
│   ├── (main)/
│   │   ├── profile/
│   │   │   └── page.tsx           # thiết lập hồ sơ gamer
│   │   ├── find-teammate/
│   │   │   └── page.tsx           # trang matchmaking / browse room
│   │   ├── party/
│   │   │   └── [roomId]/page.tsx  # party room: chat + voice + video
│   │   ├── chat/
│   │   │   ├── page.tsx
│   │   │   └── [id]/page.tsx
│   │   └── layout.tsx
│   └── layout.tsx
├── components/
│   ├── matchmaking/
│   │   ├── QuickMatchButton.tsx
│   │   ├── FilterPanel.tsx
│   │   └── MatchFoundModal.tsx
│   ├── party/
│   │   ├── PartyMemberList.tsx
│   │   ├── VoiceChannel.tsx
│   │   ├── VideoGrid.tsx
│   │   └── RoomBrowser.tsx
│   ├── chat/
│   │   ├── MessageList.tsx
│   │   ├── MessageInput.tsx
│   │   └── ConversationList.tsx
│   └── ui/
├── lib/
│   ├── websocket.ts
│   ├── webrtc.ts
│   ├── api.ts
│   └── auth.ts
├── hooks/
│   ├── useWebSocket.ts
│   ├── useWebRTC.ts
│   ├── useMatchmaking.ts
│   └── usePresence.ts
├── store/
├── types/
├── public/
├── next.config.js
└── package.json
```

---

## 6. Database Schema (rút gọn)

```
users
├── id (PK)
├── email
├── password_hash
├── display_name
├── avatar_url
├── status (online/offline/in-game)
├── last_seen_at
└── created_at

game_profiles
├── id (PK)
├── user_id (FK)
├── game_name
├── rank_data (JSONB: {"tier": "Diamond", "division": "II", "lp": 45} — cấu trúc linh hoạt theo từng game)
├── play_goal (chill/tryhard/competitive)
├── available_hours (JSONB: khung giờ rảnh trong tuần)
├── languages
└── updated_at

matchmaking_requests
├── id (PK)
├── user_id (FK)
├── party_id (FK, nullable — nếu nhóm bạn bè đã có phòng và tìm thêm slot)
├── game_name
├── group_size (số người hiện tại trong nhóm, mặc định 1 nếu solo)
├── target_team_size (số người cần đủ team, ví dụ: 5 cho LMHT)
├── criteria (JSONB: rank range, giờ, mục tiêu...)
├── status (waiting/matched/cancelled/expired)
└── created_at

parties
├── id (PK)
├── name
├── game_name
├── visibility (public/private)
├── max_members
├── status (waiting/in_game/closed)
├── created_by
├── created_at
└── closed_at (nullable — thời điểm phòng đóng, phục vụ lưu vết lịch sử)

party_members
├── party_id (FK)
├── user_id (FK)
├── role (host/member)
└── joined_at

conversations
├── id (PK)
├── type (direct / party)
├── party_id (FK, nullable)
└── created_at

messages
├── id (PK)
├── conversation_id (FK)
├── sender_id (FK)
├── content
├── type (text/image/system)
└── created_at

calls
├── id (PK)
├── party_id (FK)
├── type (voice/video)
├── started_at
└── ended_at

call_participants
├── call_id (FK)
├── user_id (FK)
├── joined_at
└── left_at

teammate_ratings
├── id (PK)
├── party_id (FK)
├── rater_id (FK)
├── rated_user_id (FK)
├── score (1-5)
├── tags (JSONB: ["friendly", "skilled", "communicative"] hoặc ["toxic"])
├── status (valid/anomaly/rejected — kết quả anti-abuse check)
├── weight (float, mặc định 1.0 — trọng số bị hạ nếu rater bị shadow penalty)
└── created_at

rater_reputation
├── user_id (PK, FK)
├── total_ratings_given (tổng số lần đánh giá)
├── anomaly_count (số lần bị phát hiện anomaly)
├── rating_weight (float, mặc định 1.0 — bị hạ dần nếu lạm dụng, về 0 = vô hiệu hóa)
└── updated_at
```

---

## 7. Các quyết định thiết kế đáng chú ý (chuẩn bị cho phỏng vấn)

- **Vì sao dùng Redis cho matchmaking queue thay vì chỉ query trực tiếp PostgreSQL?** → Matchmaking cần xử lý nhanh, tần suất cao, dữ liệu có tính tạm thời (chỉ tồn tại lúc đang chờ ghép) — Redis Sorted Set/List phù hợp hơn nhiều so với việc query liên tục vào DB quan hệ.
- **Vì sao mesh network cho voice channel nhỏ nhưng SFU cho video/group lớn?** → Mesh đơn giản, độ trễ thấp, phù hợp nhóm nhỏ (party game thường 2-5 người); nhưng khi số người/luồng video tăng, mesh làm client phải upload N-1 lần gây quá tải — SFU giải quyết vấn đề này bằng cách chỉ upload 1 lần.
- **Vì sao tách `matchmaking` thành module riêng thay vì gộp vào `party`?** → Matchmaking có domain logic riêng biệt (thuật toán ghép, quản lý queue, timeout xử lý) khác hẳn với quản lý vòng đời room — tách riêng giúp code dễ test và maintain hơn, đồng thời dễ tách thành service riêng sau này nếu cần scale matchmaking độc lập.
- **Vì sao cần Redis Pub/Sub cho cả chat lẫn presence?** → Để hệ thống scale ngang nhiều backend instance mà vẫn đảm bảo message/trạng thái đến đúng người, bất kể client đang kết nối tới instance nào.
- **Vì sao dùng JSONB cho rank thay vì cột text đơn giản?** → Mỗi game có cấu trúc rank khác nhau (LMHT: tier + division + LP, Valorant: tier + RR, CS2: rating number). JSONB cho phép lưu linh hoạt mà không cần thay đổi schema khi thêm game mới, đồng thời vẫn hỗ trợ query/indexing trong PostgreSQL.
- **Vì sao cần cơ chế Slot Filling thay vì chỉ ghép cặp 1-1?** → Game team-based cần đủ team size (3, 4, 5 người). Slot Filling cho phép ghép linh hoạt nhóm bạn bè đã có sẵn phòng + solo players, tránh bắt tất cả phải queue solo, tăng tốc độ tìm trận và trải nghiệm thực tế hơn.
- **Vì sao cần hệ thống anti-abuse cho rating thay vì tin tưởng trực tiếp?** → Trong môi trường gaming, đánh giá tiêu cực do cảm xúc (tilt) rất phổ biến. Anomaly detection + cross-check + shadow penalty bảo vệ người chơi tốt khỏi bị dìm điểm oan, đồng thời không tiết lộ cho người lạm dụng biết rằng đánh giá của họ đã bị vô hiệu hóa (tránh kích động thêm hành vi tiêu cực).

---

## 8. Định hướng mở rộng trong tương lai (optional)

- Tích hợp API của game thật (Riot API cho LMHT/Valorant) để tự động lấy rank thay vì user tự khai
- Thuật toán matchmaking nâng cao hơn (ELO-based compatibility scoring)
- Voice activity detection để hiển thị "đang nói" chính xác hơn (thay vì chỉ dựa vào mic on/off)
- AI gợi ý đồng đội dựa trên lịch sử chơi & rating
- Mobile app (React Native) dùng chung backend
