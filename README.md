# Football Web App

## Descrição

I choose to unity both projects into one, so here you can register, login, logout, get championships, get and filter matches, broadcast matches and more.

## MENU

- [Getting Started](#getting-started)
- [Setup](#setup)
- [Server](#server)
- [Client](#client)
- [Choices](#choices)

---
Look at [API Documentation](docs/api/guia.md) for more details.
For API Only you may use Postman collection [here](https://web.postman.co/54934cc3-4386-4d24-ad9c-76441e3e236d).
---

## Getting Started

### Prerequisites
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Goose](https://github.com/pressly/goose)
- [Air](https://github.com/cosmtrek/air)
- [pnpm](https://pnpm.io/)
- [Go](https://go.dev/dl/)
- [Node](https://nodejs.org/en/download/)

## Setup
### Install Docker Desktop
**Linux**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```
**MacOS**
```bash
brew install docker
```
**Windows**
https://docs.docker.com/desktop/windows/install/

### Install Nodejs
**Linux**
```bash
curl -fsSL https://deb.nodesource.com/setup_23.x | sudo -E bash -
sudo apt-get install -y nodejs
```
**MacOS**
```bash
brew install node
```
**Windows**
```bash
scoop install main/nodejs
```
### Install pnpm
**Linux**
```bash
curl -fsSL https://get.pnpm.io/install.sh | sh -
```
**MacOS**
```bash
brew install pnpm
```
**Windows**
```bash
scoop install main/pnpm
```
### Install Go
**Linux**
```bash
wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz
rm go1.25.1.linux-amd64.tar.gz
```
**MacOS**
```bash
brew install go
```
**Windows**
```bash
scoop install main/go
```
### Install Goose
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```
### Install Air
```bash
go install github.com/cosmtrek/air@latest
```

## Run both
```bash
docker compose up
# or 
docker compose up -d

goose status # check if database is connected
goose up # make migrations

// go to http://localhost:3000

# For admin = username: admin password: admin123
# For regular user just register one

```

## Server
```bash
cd server

docker compose up

goose status // check if database is connected
goose up

curl http://localhost:4000/health

# Alternativily you can run postgres and use air to run a hot-reload server

air init

air

```

## Client
```bash
cd client

pnpm install

pnpm run dev
```

## Choices
- I use next.js on frontend because its the company's choice
- On the backend I choose to apply Clean Architecture with some more reasonable naming conventions, clean architecture allows modules to be replaced without affecting the rest of the application, and it's easier to test and maintain.
- Since its a demo project I did not made 100% test coverage, but the coverage shows what needs to be shown.
- docker images are optimized for production, but they can be better using distroless images.
- I made auth token as requested on PDF, but note that it's better to use token in secure cookies with httpOnly.
- For a production environment would be nice to use a caching solution too.
