FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./cmd/main.go

# ---

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root

COPY --from=build /app/app .

RUN adduser -D -g '' app && chown -R app:app /root
USER app

EXPOSE 50052
CMD ["./app"]
