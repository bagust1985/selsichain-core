FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o selsichain ./cmd/selsichain/main.go

EXPOSE 7690-7699

CMD ["./selsichain", "--p2p-port=7690", "--testnet"]