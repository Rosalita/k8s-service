# kind

# Setup
1. Install Docker Desktop
https://docs.docker.com/desktop/install/mac-install/

Check installed ok with command to list all containers
```
docker ps
```
2. Get brew, see https://brew.sh/
3. Install Make `brew install make`
4. Use make to install cli tools with brew `make dev-brew`
5. Use make to fetch docker images `make dev-docker`

# Project layers
## Zarf
A zarf is a cardboard sleeve that goes around a coffee cup to protect you. The zarf layer protects us from being burnt by containers. This is where Docker and Kubernetes code lives.