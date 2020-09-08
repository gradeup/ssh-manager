# SSH-Manager
**Manage Employee public keys centrally onto all servers**

[![license](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/PytorchLightning/pytorch-lightning/blob/master/LICENSE)

## Running in Production
- Project not yet there yet

## Prerequisites
- Postgres
- Go (1.13+)

## Steps to setup
- $ `cp .env.example .env`
- modify .env variables
- $ `go run main.go`

## TODOs
- Use POST requests for modifying data
- Add Role based Users
- Implement CSRF tokens
- Add option to use single user instances or create seperate users on instances
- Pick instance username from servers (if single user machine), else use user's name
- Add Listen Address from env (0.0.0.0 or 127.0.0.1 or some private IP)

---

## Licence

Please observe the Apache 2.0 license that is listed in this repository.