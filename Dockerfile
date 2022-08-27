FROM golang:alpine
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o goapp ./cmd/main.go
EXPOSE 8080
ENTRYPOINT ./goapp
