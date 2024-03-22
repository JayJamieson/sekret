FROM golang:1.22-alpine as builder

ARG PKG_NAME=github.com/JayJamieson/sekret

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN apk add build-base

COPY . /go/src/${PKG_NAME}

RUN cd /go/src/${PKG_BASE}/${PKG_NAME} && \
    CGO_ENABLED=1 GOOS=linux go build -ldflags='-s -w -linkmode external -extldflags "-static"' -o /sekret

FROM scratch

COPY --from=builder /sekret /sekret

ENV PORT=8080

ARG VERSION="local"
ENV ENV_VERSION=$VERSION

EXPOSE 8080

ENTRYPOINT [ "/sekret" ]
