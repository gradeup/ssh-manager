# Installation using Docker Compose

- Setup a bare linux system
- Generate new ssh key pair
- Install Docker
- Create sshmanager directory<br>
`mkdir ~/sshmanager`
`mkdir ~/sshmanager/pg`
`touch ~/sshmanager/.env`
- Update environment variables in ~/sshmanager/.env
- Update Postgres credentials in docker-compose.yaml
- Start services with docker compose<br>
`docker-compose up -d`