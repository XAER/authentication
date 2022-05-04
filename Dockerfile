FROM golang:1.18-alpine

ENV PROJECT_DIR=/app \
    GO111MODULE=auto

WORKDIR ${PROJECT_DIR}

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon@latest

COPY . .

CMD CompileDaemon --build="go build main.go" --command=./main
