FROM golang:1.15-alpine as builder
RUN apk add --no-cache git
WORKDIR /go/ssh-manager
COPY . .
RUN go build -o ssh-manager

FROM alpine:latest
RUN mkdir ssh-manager
COPY --from=builder /go/ssh-manager/ssh-manager /go/ssh-manager/static /go/ssh-manager/ui /ssh-manager/
CMD ["/ssh-manager/ssh-manager"]