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

# Building Docker images locally
Use `make all`


# Running the cluster
Use `make dev-up-local` to start the cluster
Use `make dev-down-local` to stop the local cluster

# Interacting with the cluster
To see a list of clusters names in kind `kind get clusters`
To see cluster info, use the cluster name e.g. `kubectl cluster-info --context kind-starter-cluster`

# Project layers
## App
The `app` layer is where binaries live. It has two sub-layers `services` for services that we are building and `tooling` for cli tools we are building. Within these sub-layers the name of the folder is the name of the binary and the folder will contain a `main.go` file. The `app` layer is concerned with startup, shutdown, receiving external input, sending external output. This is the presentation layer. The `app` layer should have it's own models for what is coming in and what is going out so that it is independent of the layers below. This layer is allowed to log.

## Business
The `business` layer contains packages which are specific to solving the business problem. This layer can not import from the `app` layer that sits above it, but it can import from layers below it like `foundation` this layer is allowed to log. 

## Foundation
The `foundation` layer is packages which are not necessarily specific to the business layer, but we need them to build out services and tooling. They could potentially become third party dependencies one day and may one day move to their own repos. So we need to make sure that packages at this layer do not import each other. Packages at this layer should be treated like the standard library for our projects.

## Zarf
A zarf is a cardboard sleeve that goes around a coffee cup to protect you. The zarf layer protects us from being burnt by containers. This is where Docker and Kubernetes code lives. In the `docker` sub folder, every service which is being built has its own dockerfile. Only build version and build date are passed into dockerfiles as args. The more things that are hardcoded in the dockerfile, the easier they are to manage. 

In the `k8s` sub folder, the `dev` folder contains configuration for a local cluster (KIND). The `base` folder contains all the kubernetes manifests for services which are being deployed, grouped by service name. Base config is the same regardless of the environment the service is being deployed into. Each service will have a `service-name.yaml` and a `kustomization.yaml`. Changes will be made to what's in `base` to deploy to different environments. 

The `---` in the manifest represents where a new configuration definition begins. Each section separate by `---` could be a separate file. First thing to do is define a namespace. Second thing to do is to define a deployment.

Kustomize takes multiple yaml files and stitches them together to create one yaml file that can be applied. Kustomize does this with the command `kustomize build zarf/k8s/dev/sales`

# KIND
Kind requires that local images need to be loaded into a staging area.
