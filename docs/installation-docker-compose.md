# Installation using Docker Compose

- Setup a bare linux system<br>
`sudo apt-get update`<br>
`sudo apt-get upgrade`<br>
- Generate new ssh key pair
- Install Docker<br>
`sudo apt-get install docker`<br>
`sudo apt-get install docker-compose`<br>
- Create sshmanager directory<br>
`mkdir ~/sshmanager`
`mkdir ~/sshmanager/pg`
`touch ~/sshmanager/.env`
- Update environment variables in ~/sshmanager/.env
- Update Postgres credentials in docker-compose.yaml
- Start services with docker compose<br>
`docker-compose up -d`