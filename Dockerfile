FROM golang:1.22-alpine as builder

COPY . /go/src/github.com/JayJamieson/sekret
WORKDIR /go/src/github.com/JayJamieson/sekret

RUN apk update && apk add build-base && go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags='-s -w -linkmode external -extldflags "-static"' -o /sekret

FROM gcr.io/distroless/static

COPY --from=builder --chown=nonroot:nonroot /sekret /sekret

ARG VERSION="local"
ENV ENV_VERSION=$VERSION

EXPOSE 8080

ENTRYPOINT [ "/sekret" ]
