.PHONY: start-cluster stop-cluster status deploy-dashboard restart-cluster postgres-url tunnel swagger-docs

# Start minikube cluster
start-cluster:
	minikube start

# Stop minikube cluster
stop-cluster:
	minikube stop

# Check cluster status
status:
	minikube status
	kubectl cluster-info

# Deploy Kubernetes dashboard
deploy-dashboard:
	kubectl apply -f /etc/kubernetes/addons/dashboard-svc.yaml

# Clean restart cluster
restart-cluster:
	minikube delete
	minikube start

# Get PostgreSQL service URL
postgres-url:
	minikube service postgres --url

# Generate Swagger documentation
swagger-docs:
	go run github.com/swaggo/swag/cmd/swag init -g services/api-gateway/main.go -o services/api-gateway/docs

# Start minikube tunnel (for LoadBalancer services)
tunnel:
	minikube tunnel
