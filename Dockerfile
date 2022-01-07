FROM golang:latest
ENV GO111MODULE=on
ENV TZ=Asia/Jakarta
WORKDIR /app/ELKExample
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
CMD ["go","run","/app/ELKExample/main.go"]