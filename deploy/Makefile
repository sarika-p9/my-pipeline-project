DOCKER_REGISTRY=docker.io/sarikapt9
IMAGE_TAG=latest

# Paths to Dockerfiles
BACKEND_DOCKERFILE=deploy/docker/backend.Dockerfile
FRONTEND_DOCKERFILE=deploy/docker/frontend.Dockerfile

# Local Build
build-backend:
	docker build -f backend.Dockerfile -t my-backend:$(IMAGE_TAG) ../..

build-frontend:
	docker build -f $(FRONTEND_DOCKERFILE) -t my-frontend:$(IMAGE_TAG) .

# Tag Images
tag-backend:
	docker tag my-backend:$(IMAGE_TAG) $(DOCKER_REGISTRY)/my-backend:$(IMAGE_TAG)

tag-frontend:
	docker tag my-frontend:$(IMAGE_TAG) $(DOCKER_REGISTRY)/my-frontend:$(IMAGE_TAG)

# Docker Login
docker-login:
	docker login

# Push Images
push-backend: tag-backend docker-login
	docker push $(DOCKER_REGISTRY)/my-backend:$(IMAGE_TAG)

push-frontend: tag-frontend docker-login
	docker push $(DOCKER_REGISTRY)/my-frontend:$(IMAGE_TAG)

# Kubernetes Deployment
apply-k8s:
	kubectl apply -f backend-deployment.yaml
	kubectl apply -f frontend-deployment.yaml
	kubectl apply -f rabbitmq-deployment.yaml
	kubectl apply -f nats-deployment.yaml

restart-deployments:
	kubectl get deployments --no-headers -o custom-columns=":metadata.name" | xargs -I{} kubectl rollout restart deployment {}

delete-all:
	kubectl delete all --all --namespace=default

get-all:
	kubectl get all

# Port Forwarding
port-forward-backend:
	kubectl port-forward svc/backend-service 30002:8080 30003:50051

port-forward-frontend:
	kubectl port-forward svc/frontend-service 30001:80

# Full Build & Deploy
all: build-backend build-frontend push-backend push-frontend apply-k8s
