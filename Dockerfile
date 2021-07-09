FROM golang:1.16-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /go/bin/sekret

FROM scratch

COPY --from=builder /go/bin/sekret /go/bin/sekret

EXPOSE 8080

ENTRYPOINT [ "/go/bin/sekret" ]
