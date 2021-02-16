FROM golang:1.15-alpine3.12 AS builder

WORKDIR /go/src/github.com/sm43/goa-gorm
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api-server ./cmd/user/...

FROM alpine:3.12

RUN apk --no-cache add ca-certificates && addgroup -S hub && adduser -S hub -G hub
USER hub

WORKDIR /app

COPY --from=builder /go/src/github.com/sm43/goa-gorm/api-server /app/api-server

EXPOSE 8080

CMD [ "/app/api-server" ]