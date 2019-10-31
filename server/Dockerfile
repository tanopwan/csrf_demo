FROM golang:1.13 as builder
WORKDIR /go/src/csrf_demo/server
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -extldflags "-static"' \
    -o main .

FROM scratch
COPY --from=builder /go/src/csrf_demo/server /

CMD ["/main"]