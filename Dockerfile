# goのイメージをDockerHubから流用する(Alpine Linux)
FROM golang:1.22.0-alpine3.19
# Linuxパッケージ情報の最新化+gitがないのでgitを入れる
RUN apk update && apk add git
# ログのタイムゾーンを指定
ENV TZ /usr/share/zoneinfo/Asia/Tokyo
# コンテナ内の作業ディレクトリを指定
WORKDIR /app
# ソースコードをコンテナ内にコピー
COPY /app/* ./
# /app/go.modに記載された依存関係の解決＋必要なパッケージのダウンロードを実行
RUN go mod download
# Airのバイナリをインストール
RUN go install github.com/cosmtrek/air@latest
# コンテナの公開するポートを指定
EXPOSE 8989
# 起動時のコマンド(airを使用するため)
CMD ["air", "-c", ".air.toml"]
