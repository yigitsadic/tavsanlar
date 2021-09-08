FROM golang:1.15.15-alpine3.14 as builder

WORKDIR /src

COPY . .

RUN go build -o output

FROM alpine:3.14

COPY --from=builder /src/output /src/output

CMD ["/src/output"]
