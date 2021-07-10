# sekret
sekret store

docker pull jayjamieson/sekret:version

docker tag jayjamieson/sekret:53fec6f.3 registry.heroku.com/jay-sekret/web
docker push registry.heroku.com/jay-sekret/web
heroku container:release web -a jay-sekret
