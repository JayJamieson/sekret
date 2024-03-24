## TODO

- [ ] Add cli auto release creation on push to main
- [x] Add server auto release creation on push to main
  - [x] push to docker hub
  - [x] deploy to fly server

CLI

- [ ] Parse secret key if URL is provided
- [ ] Secret sync to environment variable

Server

- [x] Add simple UI to create/view secrets
- [ ] Add secret key encryption using passphrase
- [ ] Add storage backend support e.g. Redis, SQLite
- [ ] Add configuration support for backend storage, encryption secret
