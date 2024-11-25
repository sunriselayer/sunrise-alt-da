
FROM golang:1.22.2-alpine3.18 as builder

WORKDIR /
COPY . sunrise-op-da-server
RUN apk add --no-cache make
WORKDIR /sunrise-op-da-server
RUN make da-server

FROM alpine:3.18

COPY --from=builder /sunrise-op-da-server/bin/da-server /usr/local/bin/da-server

CMD ["da-server"]