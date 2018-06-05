FROM golang:latest as builder

WORKDIR /srv
COPY . /srv/
RUN /srv/build_all

FROM abiosoft/caddy:latest

COPY --from=builder /srv/bin /srv/bin