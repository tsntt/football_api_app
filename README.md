# Football Web App

## MENU

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
https://nodejs.org/en/download/

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
https://pnpm.io/installation

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
https://go.dev/dl/go1.25.1.windows-amd64.msi

### Install Goose
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```
### Install Air
```bash
go install github.com/cosmtrek/air@latest
```

## Server

## Client

