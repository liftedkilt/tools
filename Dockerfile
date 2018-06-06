FROM golang:latest as builder

WORKDIR /srv
COPY . /srv/

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/vX.X.X/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN /srv/build_all

FROM abiosoft/caddy:latest

COPY --from=builder /srv/bin /srv/bin

COPY Caddyfile /etc/Caddyfile

EXPOSE 80