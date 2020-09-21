# SSH-Manager
**Manage Employee public keys centrally onto all servers**

[![license](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/PytorchLightning/pytorch-lightning/blob/master/LICENSE)

## Prerequisites
- Postgres (9.6+)
- Go (1.13+)

## Installation
- [Install using docker compose](https://github.com/gradeup/ssh-manager/blob/master/docs/installation-docker-compose.md)
- [Install using Docker Image and self managed Postgres](https://github.com/gradeup/ssh-manager/blob/master/docs/installation-docker-compose.md)
- [Build using Go 1.13+ and self managed Postgres](https://github.com/gradeup/ssh-manager/blob/master/docs/installation-docker-compose.md)

## Configuration
- POSTGRES_HOST<br>
  `Endpoint for postgres server, defaults to 127.0.0.1`
- POSTGRES_PORT=5432<br>
  `Port used for postgres server, defaults to 5432`
- POSTGRES_USER=postgres<br>
  `Username for postgres server, defaults to postgres`
- POSTGRES_PASSWORD=postgres<br>
  `Password for postgres server, defualts to postgres`
- POSTGRES_DATABASE=sshmanager<br>
  `Database used for ssh manager in postgres, defaults to postgres`
- PRIVATE_KEY_PATH=/.ssh/id_rsa<br>
  `Path to private key which has access to all servers that will be managed by this service, defaults to /home/ubuntu/.ssh/id_rsa`
- PUBLIC_KEY_PATH=/.ssh/id_rsa.pub<br>
  `Path to public key which needs to be available on all instances to be managed by this service, defaults to /home/ubuntu/.ssh/id_rsa.pub`
- SERVICE_PORT=8000<br>
  `Port to start ssh-manager web service on, defaults to 8000`

## Steps for local setup
- $ `cp .env.example .env`
- modify .env variables as per configuration
- $ `go run main.go`

## TODOs
- Use POST requests for modifying data
- Add Role based Users (oAuth implementaion as well)
- Implement CSRF tokens
- Add option to use single user instances or create seperate users on instances
- Pick instance username from servers (if single user machine), else use user's name
- Add Listen Address from env (0.0.0.0 or 127.0.0.1 or some private IP)
- Support Private Key as a string from env instead of just file read
- Validate User inputs (unique username/email, trim spaces/special chars)

---

## Licence

Please observe the Apache 2.0 license that is listed in this repository.