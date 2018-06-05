FROM golang:latest as builder

WORKDIR /srv
COPY . /srv/
RUN /srv/build_all

FROM hacdias/filebrowser:latest

COPY --from=builder /srv/bin /srv/bin

WORKDIR /srv
