# sekret

[![Deployment](https://github.com/JayJamieson/sekret/actions/workflows/deployment.yml/badge.svg)](https://github.com/JayJamieson/sekret/actions/workflows/deployment.yml)

**Very Simple** secret store. Store a secret to share via self-destructing link on viewing of secret.

## Build locally

`docker build -t sekret .`

## Run locally

`docker-compose up`

## CLI

Create a secret `sekret -create <secret_value>`

Fetch a secret `sekret -fetch sweet_elgamal`

## Test requests

```sh
curl --request POST \
  --url http://localhost:8080/api/secret \
  --header 'Content-type: application/json' \
  --data '{
  "secret": "mysecret"
}'
```

```sh
curl --request GET \
  --url http://localhost:8080/api/secret/sweet_elgamal \
  --header 'Content-type: application/json'
```

## Hosting on Fly.io

1. Install [flyctl](https://fly.io/docs/hands-on/install-flyctl/)
2. `flyctl launch`

## Build container image

```sh
docker build -tag <tag name here> .
```
