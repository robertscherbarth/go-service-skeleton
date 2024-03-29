FROM golang:alpine as builder

# create build folder
RUN mkdir /build

# add build files to image and make to work dir
ADD . /build
WORKDIR /build

# Generate a binary
RUN CGO_ENABLED=0 go build -o main ./cmd/service/.

FROM scratch
# the tls certificates:
# this pulls directly from the upstream image, which already has ca-certificates:
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# copy program
COPY --from=builder /build/main /app/
WORKDIR /app

# add config file
ADD ./configs /app/configs

CMD ["./main"]
