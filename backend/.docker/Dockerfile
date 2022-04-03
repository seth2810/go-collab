FROM golang:1.18-alpine3.14 as builder
WORKDIR /go/src
COPY . .
RUN go build -o /opt/service/backend main.go

FROM alpine:3.9
WORKDIR /usr/service
COPY --from=builder /opt/service/backend /usr/service/backend
CMD ["/usr/service/backend", "serve"]
