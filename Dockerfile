FROM golang:latest AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main/main.go

FROM alpine:3.7

WORKDIR /app

COPY --from=build-stage /app/app /app

CMD ["./app"]
