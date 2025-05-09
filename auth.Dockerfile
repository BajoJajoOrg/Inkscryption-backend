FROM golang:1.24.2-alpine AS build-stage

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/sso/main.go

EXPOSE 8082

CMD [ "./main" ]