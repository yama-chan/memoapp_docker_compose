# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.13 as builder

# 命令[RUN, CMD, ENYRYPOINT, COPY, ADD]を実行するための作業用ディレクトリ
# もし指定したディレクトリがなければ、新たに作成します
WORKDIR /app

# mysql のドライバ
RUN go get github.com/go-sql-driver/mysql;

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# ホストのカレントディレクトリ「 go.* 」に当てはまるファイルをDockerイメージの作業用のディレクトリにコピー
COPY go.* ./
# 明示的な依存パッケージのダウンロードは go mod download で可能
RUN go mod download

# Copy local code to the container image.
# ビルドするためのローカルコードをイメージにコピー
COPY . ./

# Build the binary.
# https://godoc.org/github.com/gophersjp/go/src/cmd/go#hdr-Go________________
# 
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates
# bashが使えるようにする
RUN apk add bash
RUN apk add --no-cache mysql-client
# アプリ起動用のスクリプトをイメージにコピー
COPY start_app.sh ./start_app.sh
# 静的ファイル（htmlファイルなど）をコピーしたいがとりあえず全てコピーしとく
COPY . ./
# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# 今回はdocker-compose時にアプリを起動するようにする。
# そのため、イメージ内でアプリは起動しないようにコメント
# CMD ["/server"]

EXPOSE 8080
# Run the web service on container startup.
# CMD ["/server"]
