FROM golang:alpine as builder

RUN mkdir /build
ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -o main .

FROM scratch

COPY --from=builder /build/main /app/
WORKDIR /app

CMD ["./main"]