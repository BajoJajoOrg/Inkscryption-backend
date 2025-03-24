FROM golang:1.21.0-alpine AS build-stage

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/images/main.go

EXPOSE 8087

CMD [ "./main" ]