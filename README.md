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

# Running the service in the cluster
`make dev-up-local` to start the cluster
`make dev-load` to load the service image into KIND
`make dev-apply` to build kustomize and apply changes to the cluster
`make dev-logs` to see the logs of the service running in the cluster
`make dev-down-local` to stop the local cluster

after updating k8s configuration `make dev-update-apply`
after updating source code `make dev-update`

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

# K8 Quotas
Setting cpu limits on a container, will not automatically change the maximum processors the go code uses to run. GOMAXPROCS will not change if only the k8s configuration is changed. This is because Docker is not Kubernetes aware. A Go program with 4 threads will be context switching if Kubernetes is limited to 2 cores.

# Configuration
The only place configuration is allowed to be read is in `main.go`. Configuration should be read in at the start of the app and passed to where it needs to be. All configuration should have a default value which works in the dev environment. Ideally those defaults work in staging and production environments. The service should have allow for `--help` which shows everything that can be configured, default values and how to override them. The service should allow config to be overridden by an env var or a cli flag, with cli flag taking precedence. Arden labs `conf` package does this.

# Debugging
A mux has been brought up on port 4000 that has the ability to get stack traces, cpu and memory profiles, locking profiles and traces of the code while it is running. A mux is a piece of code that registers handler functions to a specific url. Then when a request comes in, if the url matches the route, execute those handlers. The standard library serve mux is being used as a router. 
The routes are defined in the mux and associated with a handler function. Never use the default serve mux as any package imports can expose endpoints on it. Always create a new mux, specifying exactly what should be exposed. When running locally debug information will be exposed at `http://localhost:4000/debug/pprof` and metrics information at `http://localhost:4000/debug/vars`. When running locally, expvarmon can be used to view a metrics dashboard locally with `make metrics-local`
