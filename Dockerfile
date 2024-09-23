
FROM golang:1.22.2-alpine3.18 as builder

WORKDIR /
COPY . sunrise-alt-da
RUN apk add --no-cache make
WORKDIR /sunrise-alt-da
RUN make da-server

FROM alpine:3.18

COPY --from=builder /sunrise-alt-da/bin/da-server /usr/local/bin/da-server

CMD ["da-server"]