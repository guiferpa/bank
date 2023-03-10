FROM golang:1.19 AS builder

WORKDIR /opt/app

COPY . .

RUN go build -o ./dist/api ./cmd/api/main.go

FROM debian

COPY --from=builder /opt/app/dist/api /bin/app

CMD ["/bin/app"]