FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o url-shortener ./cmd/url-shortener

EXPOSE 8083

CMD ["/app/url-shortener", "-config=/app/config/local.yaml"]