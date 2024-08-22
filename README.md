# ports service assignment

This creates a simple go service that reads a json file containing port data and stores it in an in-memory database, and makes that available over HTTP. The service can be run in a docker container in Kubernetes. The service aims to emit telemetry in the form of logs, traces and metrics.

### Scope

I have decided to try to meet most of the criteria from the original assignment, but I have aimed to sway it a bit more towards SRE and DevOps practices, deviating from the original assignment specification and expanding the scope. I think that this might be a bit more representative of what I would see my self doing day to day, below are some things I have moved out of scope, and into scope.

Hopefully this is ok, and if not, I am happy to do the assignment again with the original scope. I have tried to keep the service simple, and not over engineer it, and it definitely cuts some corners or does things in less than optimal or ideal ways, which i will hopefully explain in comments throughout.

#### Out of scope
- I am not following hexagonal architecture, i probably wont do it justice in the time i have.
- The file is not streamed into memory, but read into memory in one go. This is because the file is small, and the service is simple. In future there are ways i would do this and i put comments in the code about it.
- Tests are rather lacking. There isn't really any business logic or behavior in here, and its mostly implementation. I did not want to spend time creating extensive mocks, but if i were to do this properly I would spend alot more time on testing.
- The directory layout and structure is not ideal.

#### In scope
- Simple /ready endpoint
- Exposed a Otel metrics and prometheus metrics endpoint.
- HTTP server metrics so that we can use them as high level SLIs (total requests and 99th percentile) to track high level SLOs.
- Emits Otel Tracing to a jaeger instance, this could be tempo, datadog etc if we wish.
- Otel Traces down the stack to the database layer.
- Json logging, and a log level that can be set via an environment variable.
- Multi stage docker build.
- Helm chart to deploy the service.
- CI/CD to build, test and publish the service to a docker registry.
- Dummy deployment pipeline that pretends to deploy the service to a k8s cluster.

#### Future work
I am not really happy with the main application code, i think if i spent a bit longer and had some conversations i could clean it up, write more tests, adhere to ddd and make it a bit more idiomatic. I would probably move it to a hexagonal architecture. However, i wanted to showcase some more skills around telemetry, deployment, ci/cd etc, so i tried to timebox myself on this, therefore it feels a bit half finished.

- I would look at creating an elegant logging abstraction that could be used throughout the service, but also other services.
- Add more context to the traces and logs (like a request id), so that we can trace a request through the system easily.
- Fully fledged gitops pipeline that would deploy the service to k8s clusters.
- Plug this into a developer platform like backstage, and create a service catalog entry for it.
- Create a shared package for otel metrics and tracing.
- Use gorilla mux or similar to handle routing, this would make it easier to introduce handlers and middleware for things like tracing and logs.
- Use a real database and auto instrument it with otel, using otel contrib libraries.
- Decide upon some sort of acceptable testing strategy.


### Testing, Building and Running the service

Requirements:
- Docker
- Kubernetes (kind, k3s, docker desktop etc)
- Helm

You can use visual studio code dev containers to run the tests and the service, or you can run it locally.

Steps to run locally:
1. Clone the repository
2. Run `make test` to run the tests
3. Run `make run` to run the service
4. The service will be available on `localhost:8080`
   1. /ports, /metrics, /ready endpoints are available
5. Run `make build` to build the docker image
6. Run `make deploy` to deploy the docker image to your current k8s context
7. Observe the helm message to test/explore the service

### Bonus
I have added a run-jaeger target to the makefile, this will run a jaeger instance in a docker container, and you can view the traces that are emitted from the service at http://localhost:16686/.

### Debugging in k8s
1. find the pods with `kubectl get pods`
2. kubectl debug -it $pod-name-here --image=ubuntu --target=port-service -- /bin/bash 
