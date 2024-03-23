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

```bash
curl --request POST \
  --url http://localhost:8080/api/secret \
  --header 'Content-type: application/json' \
  --data '{
  "secret": "mysecret"
}'
```

```bash
curl --request GET \
  --url http://localhost:8080/api/secret/sweet_elgamal \
  --header 'Content-type: application/json'
```

## Hosting on Fly.io

1. Install [flyctl](https://fly.io/docs/hands-on/install-flyctl/)
2. `flyctl launch`

## Hosting on heroku

```bash
docker pull jayjamieson/sekret:<some version>
// or
docker build -tag sekret .

docker tag <docker hub name>/sekret:<some version> registry.heroku.com/<heroku project>/web
docker push registry.heroku.com/<heroku project>/web
heroku container:release web -a <heroku project>
```

## TODO

- [ ] Add cli and server auto release creation on push to main

CLI

- [ ] Parse secret key if URL is provided
- [ ] Secret sync to environment variable

Server

- [x] Add simple UI to create/view secrets
- [ ] Add secret key encryption
- [ ] Add storage backend support e.g. Redis, SQLite
- [ ] Add configuration support for backend storage, encryption secret
