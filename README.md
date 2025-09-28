# Football Web App

## Descrição

bla bla bla

## MENU

- [Getting Started](#getting-started)
- [Setup](#setup)
- [Server](#server)
- [Client](#client)

---
Veja a [Documentação da API](docs/api/guia.md) para mais detalhes.
---

## Getting Started

### Prerequisites
- Docker Desktop (https://www.docker.com/products/docker-desktop/)
- Goose (https://github.com/pressly/goose)
- Air (https://github.com/cosmtrek/air)
- pnpm (https://pnpm.io/)
- Go 1.25.1
- Node 23.10.0

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

goose status // check if database is connected
goose up

curl http://localhost:4000/health

// go to http://localhost:3000

```

## Server
```bash
cd server

docker compose up

goose status // check if database is connected
goose up

curl http://localhost:4000/health

```

## Client
```bash
cd client

pnpm install

pnpm run dev
```

