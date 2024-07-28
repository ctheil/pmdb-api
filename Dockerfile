FROM golang:1.22

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading deps and only redownloading in build on change
COPY go.mod ./

RUN go mod download && go mod verify

COPY . .
RUN go build -v -o main ./cmd/pmdb/main.go

CMD ["/app/main"]
