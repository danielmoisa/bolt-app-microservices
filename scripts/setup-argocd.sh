#!/bin/bash

# ArgoCD Quick Setup Script for Local Testing
set -e

echo "🚀 Setting up ArgoCD for GitOps testing..."

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is required but not installed."
    exit 1
fi

# Check if minikube is running
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Kubernetes cluster not accessible. Start minikube first."
    echo "Run: minikube start"
    exit 1
fi

echo "✅ Kubernetes cluster is accessible"

# Create ArgoCD namespace
echo "📦 Creating ArgoCD namespace..."
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -

# Install ArgoCD
echo "🔄 Installing ArgoCD..."
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for ArgoCD to be ready
echo "⏳ Waiting for ArgoCD to be ready..."
kubectl wait --for=condition=available --timeout=600s deployment/argocd-server -n argocd

# Get admin password
echo "🔑 Getting ArgoCD admin password..."
ARGOCD_PASSWORD=$(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)

echo ""
echo "🎉 ArgoCD setup complete!"
echo ""
echo "📋 Access Details:"
echo "   URL: https://localhost:8080"
echo "   Username: admin"
echo "   Password: $ARGOCD_PASSWORD"
echo ""
echo "🚀 To access ArgoCD:"
echo "   kubectl port-forward svc/argocd-server -n argocd 8080:443"
echo ""
echo "📖 Next steps:"
echo "   1. Run the port-forward command above"
echo "   2. Open https://localhost:8080 in your browser"
echo "   3. Login with admin/$ARGOCD_PASSWORD"
echo "   4. Add your GitHub repository"
echo "   5. Create applications for development and production"
echo ""
echo "📚 Full guide: See ARGOCD_TESTING.md"
