.PHONY: postgres-url swagger-docs setup-all run-argo apply-argo-apps run-tilt clean-all check-context setup-minikube status help build-images lint security test lint-fix

# Show available commands
help:
	@echo "🚀 Bolt App Microservices"
	@echo "========================"
	@echo ""
	@echo "Setup:     setup-all, setup-minikube"
	@echo "Dev:       run-tilt, run-argo, build-images"
	@echo "Quality:   lint, security, test, lint-fix"
	@echo "Utils:     status, check-context, clean-all"
	@echo ""

# Setup minikube
setup-minikube: check-context
	minikube start
	@eval $$(minikube docker-env)

# Complete setup
setup-all: setup-minikube
	@echo "🚀 Complete setup..."
	./scripts/setup-argocd.sh
	make swagger-docs
	make build-images
	kubectl apply -f deploy/development/k8s/argocd-applications.yaml
	@echo "✅ Ready! Next: make run-argo or make run-tilt"

# Get PostgreSQL URL
postgres-url: check-context
	minikube service postgres --url

# Apply ArgoCD applications
apply-argo-apps: check-context build-images
	kubectl apply -f deploy/development/k8s/argocd-applications.yaml
	
# Generate Swagger docs
swagger-docs:
	go run github.com/swaggo/swag/cmd/swag init -g services/api-gateway/main.go -o services/api-gateway/docs

# Build Go binaries
build-binaries: check-context
	@echo "Building Go binaries..."
	@eval $$(minikube docker-env) && \
	mkdir -p build && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o build/api-gateway ./services/api-gateway && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o build/trip-service ./services/trip-service/cmd && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o build/driver-service ./services/driver-service && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o build/payment-service ./services/payment-service/cmd

# Build Docker images
build-images: check-context build-binaries
	@echo "Building images..."
	@eval $$(minikube docker-env) && \
	docker build -t bolt-app/api-gateway -f deploy/development/docker/api-gateway.Dockerfile . && \
	docker build -t bolt-app/trip-service -f deploy/development/docker/trip-service.Dockerfile . && \
	docker build -t bolt-app/driver-service -f deploy/development/docker/driver-service.Dockerfile . && \
	docker build -t bolt-app/payment-service -f deploy/development/docker/payment-service.Dockerfile .

# Run Tilt development
run-tilt: check-context build-images
	@echo "Starting Tilt..."
	@echo "Services: API Gateway:8081, Swagger:8082, RabbitMQ:15672, Jaeger:16686"
	@eval $$(minikube docker-env) && tilt up

# Access ArgoCD UI
run-argo: check-context
	@echo "ArgoCD: https://localhost:8080"
	@echo "User: admin"
	@echo "Pass: $$(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d 2>/dev/null || echo "ArgoCD not installed")"
	kubectl port-forward svc/argocd-server -n argocd 8080:443

# Code quality and security checks
lint:
	@echo "🔍 Running golangci-lint..."
	golangci-lint run ./...

lint-fix:
	@echo "🔧 Running golangci-lint with auto-fix..."
	golangci-lint run --fix ./...

security:
	@echo "🔒 Running security analysis with gosec..."
	gosec -fmt=json -out=gosec-report.json ./...

test:
	@echo "🧪 Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Combined quality check
quality: lint security test
	@echo "✅ All quality checks completed!"

# Check missing targets
check-context:
	@echo "Checking kubectl context..."
	@kubectl config current-context

status: check-context
	@echo "📊 Cluster Status:"
	@kubectl get pods -n bolt-app
	@echo ""
	@echo "🌐 Services:"
	@kubectl get svc -n bolt-app

clean-all: check-context
	@echo "🧹 Cleaning up..."
	kubectl delete namespace bolt-app --ignore-not-found=true
	minikube stop