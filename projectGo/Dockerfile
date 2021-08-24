FROM golang:latest

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o /docker-run

EXPOSE 8080

CMD [ "/docker-run" ]