FROM golang:latest AS builder

WORKDIR /tools
ADD ./* ./
RUN ./build_all

FROM abiosoft/caddy:latest

COPY --from=builder /tools/bin .

EXPOSE 443