# sekret
**Very Simple** secret store. Store a secret to share via self destructing link.

## Build locally

`docker build -t sekret .`

## Run locally

`docker-compose up`

## Test requests

```bash
curl --request POST \
  --url http://localhost:8080/secret \
  --header 'content-type: application/json' \
  --data '{
  "data": "mysecret",
  "owner": "owner"
}'
```

```bash
curl --request GET \
  --url https://localhost:8080/secret/sweet_elgamal \
  --header 'content-type: application/json'
```

## Hosting on heroku

```bash
docker pull jayjamieson/sekret:<some version>

docker tag <docker hub name>/sekret:<some version> registry.heroku.com/<heroku project>/web
docker push registry.heroku.com/<heroku project>/web
heroku container:release web -a <heroku project>
```