DOCKER_REGISTRY=docker.io/sarikapt9
IMAGE_TAG=latest

# Paths to Dockerfiles
BACKEND_DOCKERFILE=deploy/docker/backend.Dockerfile
FRONTEND_DOCKERFILE=deploy/docker/frontend.Dockerfile
DEMOCTL_DOCKERFILE=deploy/docker/democtl.Dockerfile

# Build Backend Image
build-backend:
	docker build -t $(DOCKER_REGISTRY)/backend:$(IMAGE_TAG) -f $(BACKEND_DOCKERFILE) .

# Build Frontend Image
build-frontend:
	docker build -t $(DOCKER_REGISTRY)/frontend:$(IMAGE_TAG) -f $(FRONTEND_DOCKERFILE) .

# Build Democtl CLI Image
build-democtl:
	docker build -t $(DOCKER_REGISTRY)/democtl:$(IMAGE_TAG) -f $(DEMOCTL_DOCKERFILE) .

# Push Backend Image
push-backend:
	docker push $(DOCKER_REGISTRY)/backend:$(IMAGE_TAG)

# Push Frontend Image
push-frontend:
	docker push $(DOCKER_REGISTRY)/frontend:$(IMAGE_TAG)

# Push Democtl CLI Image
push-democtl:
	docker push $(DOCKER_REGISTRY)/democtl:$(IMAGE_TAG)

# Deploy Kubernetes Resources
apply-k8s:
	kubectl apply -f deploy/kubernetes/

# Delete Kubernetes Resources
delete-k8s:
	kubectl delete -f deploy/kubernetes/

# Full Build & Deploy
all: build-backend build-frontend build-democtl push-backend push-frontend push-democtl apply-k8s
