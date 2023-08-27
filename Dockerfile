FROM golang:1.21.0-bullseye as builder

WORKDIR /go/src/github.com/ddung1203/go
RUN git clone https://github.com/ddung1203/go.git .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s' -o main main.go

FROM scratch
WORKDIR /usr/src/app
COPY --from=builder /go/src/github.com/ddung1203/go/ .

CMD ["./main"]