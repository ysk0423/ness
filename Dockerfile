# マルチステージビルドを使用
FROM golang:1.24.4-alpine AS builder

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係ファイルをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# 静的リンクバイナリをビルド（CGO無効化）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# distrolessイメージを使用（最小サイズ）
FROM gcr.io/distroless/static:nonroot

# 作業ディレクトリを設定
WORKDIR /

# ビルド済みバイナリをコピー
COPY --from=builder /app/main .

# nonrootユーザーで実行
USER nonroot:nonroot

# ポートを公開
EXPOSE 8080

# アプリケーションを実行
ENTRYPOINT ["./main"]