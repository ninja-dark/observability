FROM golang:1.18.4 as builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o metrics

ENTRYPOINT ["/app/metrics"]