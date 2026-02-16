# ビルド用イメージ
FROM golang:1.25-alpine AS builder

# ビルドに必要なパッケージをインストール
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# バイナリサイズを削減してビルド
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/server/main.go

# 実行用イメージ
FROM alpine:3.19

# CA証明書をインストール（HTTPS通信に必要）
RUN apk --no-cache add ca-certificates tzdata

# 非rootユーザーを作成
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

WORKDIR /home/appuser

# ビルドしたバイナリをコピー
COPY --from=builder --chown=appuser:appuser /app/main .

# 非rootユーザーに切り替え
USER appuser

EXPOSE 8080

# ヘルスチェック（Railway用）
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]