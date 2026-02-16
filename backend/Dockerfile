# ビルド用イメージ
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/server/main.go

# 実行用イメージ
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
# Firebaseの秘密鍵を環境変数から読み込む場合は不要だが、
# もしファイルが必要な場合はここで調整
EXPOSE 8080
CMD ["./main"]