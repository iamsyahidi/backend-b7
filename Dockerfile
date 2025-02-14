FROM golang:alpine as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 
COPY go.mod go.sum /go/src/backend-b7/
WORKDIR /go/src/backend-b7
RUN go mod download
COPY . /go/src/backend-b7
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o build/backend-b7 backend-b7

FROM alpine:latest
# RUN apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /go/src/backend-b7/build/backend-b7 .
COPY --from=builder /go/src/backend-b7/.env .
EXPOSE 3001
ENTRYPOINT ["/app/backend-b7"]



