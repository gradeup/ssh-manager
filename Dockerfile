FROM golang:1.17-alpine as builder
RUN apk add --no-cache git
WORKDIR /go/ssh-manager
COPY . .
RUN go build -o ssh-manager

FROM alpine:latest
RUN mkdir ssh-manager
WORKDIR /ssh-manager
COPY --from=builder /go/ssh-manager/ssh-manager /ssh-manager/
COPY --from=builder /go/ssh-manager/static /ssh-manager/static/
COPY --from=builder /go/ssh-manager/ui /ssh-manager/ui/
CMD ["./ssh-manager"]