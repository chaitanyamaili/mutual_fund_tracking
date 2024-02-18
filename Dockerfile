ARG GOLANG_VER=1.21
FROM golang:${GOLANG_VER}-alpine

WORKDIR /app

ENV APP_ENV "docker"
ENV GOCACHE /tmp/

CMD ["go", "run", "main.go"]
