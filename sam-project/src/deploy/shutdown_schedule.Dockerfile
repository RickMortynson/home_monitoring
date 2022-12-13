FROM golang:1.19.3-alpine

WORKDIR /app

COPY . /app

RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux go build -o shutdown-schedule cmd/shutdown-schedule/main.go

CMD [ "./shutdown-schedule" ]