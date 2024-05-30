FROM golang:1.22.3

WORKDIR /app
RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN go test -v ./internal/tests/
EXPOSE 8000
CMD [ "air" ]
