## Local Devnet setup - MacOS

### Steps

*start docker*

1. `docker system prune --all --volumes`
2. `minikube delete`
3. `minikube start --memory=8192`
4. `minikube tunnel`

*open new window in mev-commit dir*

1. `eval $(minikube docker-env)`
2. `make build-minikube-mac`

*cd infrastructure/charts*

1. `make check-images`
2. `make deploy-devnet`
