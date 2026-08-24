<#
.SYNOPSIS
    GameMate Project Scaffolding Script
.DESCRIPTION
    Tạo toàn bộ cấu trúc thư mục và file rỗng cho dự án GameMate
    theo đúng chuẩn phần 5 của PRD (PROJECT_DESCRIPTION.md).
    Bao gồm: Backend (Go) + Frontend (Next.js App Router, TypeScript)
.NOTES
    Author: Senior Software Architect
    Date:   2026-08-21
#>

$ErrorActionPreference = "Stop"

$ROOT = Split-Path -Parent $MyInvocation.MyCommand.Path

Write-Host ""
Write-Host "=============================================" -ForegroundColor Cyan
Write-Host "  GameMate - Project Scaffolding Script"      -ForegroundColor Cyan
Write-Host "=============================================" -ForegroundColor Cyan
Write-Host ""

# ──────────────────────────────────────────────────────
# Helper: tạo file rỗng (bao gồm tạo thư mục cha)
# ──────────────────────────────────────────────────────
function New-EmptyFile {
    param([string]$Path)
    $fullPath = Join-Path $ROOT $Path
    $dir = Split-Path -Parent $fullPath
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
    }
    if (-not (Test-Path $fullPath)) {
        New-Item -ItemType File -Path $fullPath -Force | Out-Null
    }
}

# ──────────────────────────────────────────────────────
# Helper: tạo thư mục rỗng
# ──────────────────────────────────────────────────────
function New-EmptyDir {
    param([string]$Path)
    $fullPath = Join-Path $ROOT $Path
    if (-not (Test-Path $fullPath)) {
        New-Item -ItemType Directory -Path $fullPath -Force | Out-Null
    }
}

# ══════════════════════════════════════════════════════
#  BACKEND (Go)
# ══════════════════════════════════════════════════════
Write-Host "[1/2] Scaffolding Backend (Go)..." -ForegroundColor Yellow

# --- cmd/server ---
New-EmptyFile "backend/cmd/server/main.go"

# --- internal/auth ---
New-EmptyFile "backend/internal/auth/handler.go"
New-EmptyFile "backend/internal/auth/service.go"
New-EmptyFile "backend/internal/auth/repository.go"
New-EmptyFile "backend/internal/auth/model.go"

# --- internal/profile ---
New-EmptyFile "backend/internal/profile/handler.go"
New-EmptyFile "backend/internal/profile/service.go"
New-EmptyFile "backend/internal/profile/repository.go"
New-EmptyFile "backend/internal/profile/model.go"

# --- internal/matchmaking ---
New-EmptyFile "backend/internal/matchmaking/handler.go"
New-EmptyFile "backend/internal/matchmaking/service.go"
New-EmptyFile "backend/internal/matchmaking/queue.go"
New-EmptyFile "backend/internal/matchmaking/worker.go"

# --- internal/chat ---
New-EmptyFile "backend/internal/chat/handler.go"
New-EmptyFile "backend/internal/chat/service.go"
New-EmptyFile "backend/internal/chat/repository.go"
New-EmptyFile "backend/internal/chat/model.go"

# --- internal/party ---
New-EmptyFile "backend/internal/party/handler.go"
New-EmptyFile "backend/internal/party/service.go"
New-EmptyFile "backend/internal/party/repository.go"
New-EmptyFile "backend/internal/party/model.go"

# --- internal/realtime ---
New-EmptyFile "backend/internal/realtime/hub.go"
New-EmptyFile "backend/internal/realtime/client.go"
New-EmptyFile "backend/internal/realtime/pubsub.go"
New-EmptyFile "backend/internal/realtime/protocol.go"

# --- internal/call ---
New-EmptyFile "backend/internal/call/signaling.go"
New-EmptyFile "backend/internal/call/sfu.go"
New-EmptyFile "backend/internal/call/room.go"

# --- internal/rating ---
New-EmptyFile "backend/internal/rating/handler.go"
New-EmptyFile "backend/internal/rating/service.go"

# --- internal/notification (module rỗng, sẽ bổ sung sau) ---
New-EmptyDir  "backend/internal/notification"

# --- internal/storage (module rỗng, sẽ bổ sung sau) ---
New-EmptyDir  "backend/internal/storage"

# --- internal/middleware ---
New-EmptyDir  "backend/internal/middleware"

# --- pkg ---
New-EmptyDir  "backend/pkg/config"
New-EmptyDir  "backend/pkg/logger"
New-EmptyDir  "backend/pkg/database"
New-EmptyDir  "backend/pkg/redis"

# --- migrations ---
New-EmptyDir  "backend/migrations"

# --- deployments ---
New-EmptyFile "backend/deployments/Dockerfile"
New-EmptyFile "backend/deployments/docker-compose.yml"
New-EmptyDir  "backend/deployments/k8s"

# --- test ---
New-EmptyDir  "backend/test"

# --- root files ---
New-EmptyFile "backend/go.mod"
New-EmptyFile "backend/README.md"

Write-Host "  -> Backend structure created." -ForegroundColor Green

# ══════════════════════════════════════════════════════
#  FRONTEND (Next.js / App Router / TypeScript)
# ══════════════════════════════════════════════════════
Write-Host "[2/2] Scaffolding Frontend (Next.js)..." -ForegroundColor Yellow

# --- app/(auth) ---
New-EmptyFile "frontend/app/(auth)/login/page.tsx"
New-EmptyFile "frontend/app/(auth)/register/page.tsx"

# --- app/(main)/profile ---
New-EmptyFile "frontend/app/(main)/profile/page.tsx"

# --- app/(main)/find-teammate ---
New-EmptyFile "frontend/app/(main)/find-teammate/page.tsx"

# --- app/(main)/party/[roomId] ---
New-EmptyFile "frontend/app/(main)/party/[roomId]/page.tsx"

# --- app/(main)/chat ---
New-EmptyFile "frontend/app/(main)/chat/page.tsx"
New-EmptyFile "frontend/app/(main)/chat/[id]/page.tsx"

# --- app/(main)/layout.tsx ---
New-EmptyFile "frontend/app/(main)/layout.tsx"

# --- app/layout.tsx (root layout) ---
New-EmptyFile "frontend/app/layout.tsx"

# --- components/matchmaking ---
New-EmptyFile "frontend/components/matchmaking/QuickMatchButton.tsx"
New-EmptyFile "frontend/components/matchmaking/FilterPanel.tsx"
New-EmptyFile "frontend/components/matchmaking/MatchFoundModal.tsx"

# --- components/party ---
New-EmptyFile "frontend/components/party/PartyMemberList.tsx"
New-EmptyFile "frontend/components/party/VoiceChannel.tsx"
New-EmptyFile "frontend/components/party/VideoGrid.tsx"
New-EmptyFile "frontend/components/party/RoomBrowser.tsx"

# --- components/chat ---
New-EmptyFile "frontend/components/chat/MessageList.tsx"
New-EmptyFile "frontend/components/chat/MessageInput.tsx"
New-EmptyFile "frontend/components/chat/ConversationList.tsx"

# --- components/ui (rỗng, sẽ thêm shared components sau) ---
New-EmptyDir  "frontend/components/ui"

# --- lib ---
New-EmptyFile "frontend/lib/websocket.ts"
New-EmptyFile "frontend/lib/webrtc.ts"
New-EmptyFile "frontend/lib/api.ts"
New-EmptyFile "frontend/lib/auth.ts"

# --- hooks ---
New-EmptyFile "frontend/hooks/useWebSocket.ts"
New-EmptyFile "frontend/hooks/useWebRTC.ts"
New-EmptyFile "frontend/hooks/useMatchmaking.ts"
New-EmptyFile "frontend/hooks/usePresence.ts"

# --- store ---
New-EmptyDir  "frontend/store"

# --- types ---
New-EmptyDir  "frontend/types"

# --- public ---
New-EmptyDir  "frontend/public"

# --- root files ---
New-EmptyFile "frontend/next.config.js"
New-EmptyFile "frontend/package.json"

Write-Host "  -> Frontend structure created." -ForegroundColor Green

# ══════════════════════════════════════════════════════
#  SUMMARY
# ══════════════════════════════════════════════════════
Write-Host ""
Write-Host "=============================================" -ForegroundColor Cyan
Write-Host "  Scaffolding complete!" -ForegroundColor Green
Write-Host "=============================================" -ForegroundColor Cyan
Write-Host ""

# Đếm số file và thư mục đã tạo
$backendFiles = (Get-ChildItem -Path (Join-Path $ROOT "backend") -Recurse -File).Count
$backendDirs  = (Get-ChildItem -Path (Join-Path $ROOT "backend") -Recurse -Directory).Count
$frontendFiles = (Get-ChildItem -Path (Join-Path $ROOT "frontend") -Recurse -File).Count
$frontendDirs  = (Get-ChildItem -Path (Join-Path $ROOT "frontend") -Recurse -Directory).Count

Write-Host "  Backend:  $backendDirs directories, $backendFiles files" -ForegroundColor White
Write-Host "  Frontend: $frontendDirs directories, $frontendFiles files" -ForegroundColor White
Write-Host "  Total:    $($backendDirs + $frontendDirs) directories, $($backendFiles + $frontendFiles) files" -ForegroundColor White
Write-Host ""
Write-Host "  Project root: $ROOT" -ForegroundColor DarkGray
Write-Host ""
