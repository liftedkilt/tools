FROM golang:latest as builder

WORKDIR /tools
COPY . /tools/
RUN /tools/build_all

FROM abiosoft/caddy:latest

COPY --from=builder /tools/bin .

EXPOSE 443