.PHONY:  postgres-url  swagger-docs setup-all forward-argo run-all apply-argo-apps

# Get PostgreSQL service URL
postgres-url:
	minikube service postgres --url

# Apply argo applications
apply-argo-apps:
	kubectl apply -f deploy/development/k8s/argocd-applications.yaml
	
# Generate Swagger documentation
swagger-docs:
	go run github.com/swaggo/swag/cmd/swag init -g services/api-gateway/main.go -o services/api-gateway/docs

# Setup the development environment
setup-all:
	 minikube start
	 ./scripts/setup-argocd.sh
	 
# Run tilt
run-tilt:
	tilt up

# Forward ArgoCD service
run-argo:
	kubectl port-forward svc/argocd-server -n argocd 8080:443