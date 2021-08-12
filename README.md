# sekret
sekret store

docker pull jayjamieson/sekret:version

docker tag jayjamieson/sekret:7d290f4.4 registry.heroku.com/jay-sekret/web
docker push registry.heroku.com/jay-sekret/web
heroku container:release web -a jay-sekret

## Build locally

`docker build -t sekret .`

## Run locally

`docker-compose up`

## Test requests
curl --request POST \
  --url http://localhost:8080/secret \
  --header 'content-type: application/json' \
  --data '{
  "data": "mysecret",
  "owner": "owner"
}'

curl --request GET \
  --url https://localhost:8080/secret/sweet_elgamal \
  --header 'content-type: application/json'