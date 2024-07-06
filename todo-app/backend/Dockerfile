# ベースイメージを指定
FROM golang:1.18-alpine

# 作業ディレクトリを設定
WORKDIR /app

# モジュールファイルをコピーして依存関係をダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションのビルド
RUN go build -o main .

# アプリケーションがリッスンするポートを指定
EXPOSE 8080

# コンテナ起動時に実行されるコマンドを指定
CMD ["./main"]

