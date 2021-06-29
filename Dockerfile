FROM golang:1.15.8 AS builder
WORKDIR /src
ADD . /src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o /svc ./cmd

FROM alpine
COPY --from=builder /svc /

ENV APP_HOST_NAME=http://gateway.marvel.com
ENV APP_PUBLIC_KEY=5d04d11a2b5d49b6aa05a49998b5f083
ENV APP_PRIVATE_KEY=481732840d168b3244a9f2a4efb3c21d2b1ce2e0
ENV APP_TS=1

EXPOSE 8080
CMD ["/svc"]