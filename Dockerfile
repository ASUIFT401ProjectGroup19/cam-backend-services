FROM golang:1.17.2 as builder

RUN mkdir /src
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o authenticationServer ./cmd/authentication/

FROM alpine as cert

RUN apk add --update --no-cache ca-certificates

FROM scratch
COPY --from=cert /etc/ssl/certs/* /etc/ssl/certs
COPY --from=builder /src/authenticationServer /

CMD ["/authenticationServer"]
