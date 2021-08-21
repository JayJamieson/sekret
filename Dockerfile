FROM golang:1.16-alpine as builder

ARG PKG_NAME=github.com/JayJamieson/sekret

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

COPY . /go/src/${PKG_NAME}

RUN cd /go/src/${PKG_BASE}/${PKG_NAME} && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /sekret

FROM scratch

COPY --from=builder /sekret .

ENV PORT=8080

ARG VERSION="local"
ENV ENV_VERSION=$VERSION

EXPOSE 8080

ENTRYPOINT [ "./sekret" ]
