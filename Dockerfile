FROM golang:1.26.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . ./
RUN CGO_ENABLED=0 go build -o /build/scheduler ./cmd/main.go

FROM alpine:3.24 AS runtime
WORKDIR /app
COPY --from=builder /build/scheduler ./scheduler
COPY --from=builder /app/web ./web
RUN mkdir -p /app/data
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db
ENV TODO_PASSWORD=""
ENV TODO_SECRET_KEY=""
ENTRYPOINT ["./scheduler"]