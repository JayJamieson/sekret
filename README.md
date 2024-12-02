# sekret

[![Deployment](https://github.com/JayJamieson/sekret/actions/workflows/deployment.yml/badge.svg)](https://github.com/JayJamieson/sekret/actions/workflows/deployment.yml)

**Very Simple** secret store.

Store a secret, share via self-destructing link on revealing of secret.

## Build locally

`docker build -t sekret .`

## Hosting on Fly.io

1. Install [flyctl](https://fly.io/docs/hands-on/install-flyctl/)
2. `flyctl launch`

## Build container image

```sh
docker build -tag <tag name here> .
```
