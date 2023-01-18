FROM golang:1.18.3-alpine3.16 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN go build -o main .

FROM alpine:3.16 AS runner

WORKDIR /app

RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta

COPY --from=builder /app/main .

EXPOSE 8002

CMD ["./main"]