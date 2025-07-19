# 🔄 ArgoCD GitOps Testing Setup

## Overview
Use ArgoCD locally to test GitOps deployments without real cloud infrastructure.

## 🚀 Quick Setup

### 1. Install ArgoCD in Minikube
```bash
# Make sure Minikube is running
minikube start

# Create ArgoCD namespace
kubectl create namespace argocd

# Install ArgoCD
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for ArgoCD to be ready
kubectl wait --for=condition=available --timeout=600s deployment/argocd-server -n argocd
```

### 2. Access ArgoCD UI
```bash
# Port-forward ArgoCD server
kubectl port-forward svc/argocd-server -n argocd 8080:443

# Get admin password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d && echo
```

**Access ArgoCD at**: https://localhost:8080
- **Username**: `admin`
- **Password**: (from command above)

### 3. Add Your Repository
In ArgoCD UI:
1. **Settings** → **Repositories**
2. **Connect Repo**
3. **Repository URL**: `https://github.com/danielmoisa/bolt-app-microservices`
4. **Connection Method**: `HTTPS`
5. **Connect**

### 4. Create Applications

#### Development App (Auto-sync)
```bash
kubectl apply -f deploy/development/k8s/argocd-applications.yaml
```

Or via UI:
1. **Applications** → **New App**
2. **Application Name**: `bolt-app-dev`
3. **Project**: `default`
4. **Repository URL**: `https://github.com/danielmoisa/bolt-app-microservices`
5. **Path**: `deploy/development/k8s`
6. **Cluster URL**: `https://kubernetes.default.svc`
7. **Namespace**: `bolt-app`
8. **Sync Policy**: `Automatic`

#### Production App (Manual sync)
1. **Application Name**: `bolt-app-prod`
2. **Path**: `deploy/production/k8s`
3. **Namespace**: `bolt-app-prod`
4. **Sync Policy**: `Manual`

## 🧪 Testing Scenarios

### Test 1: GitOps Workflow
```bash
# 1. Make a change to a deployment
vim deploy/development/k8s/api-gateway-deployment.yaml

# 2. Commit and push
git add .
git commit -m "Test GitOps deployment"
git push origin main

# 3. Watch ArgoCD auto-sync in UI
# 4. See changes applied automatically
```

### Test 2: Production Deployment Simulation
```bash
# 1. Use ArgoCD UI to manually sync production app
# 2. Compare dev vs prod configurations
# 3. Test rollback functionality
```

### Test 3: Configuration Drift Detection
```bash
# 1. Manually change a deployment
kubectl edit deployment api-gateway -n bolt-app

# 2. Watch ArgoCD detect drift
# 3. See auto-healing in action (if enabled)
```

## 📊 ArgoCD Features to Test

### 1. **Application Health**
- Monitor pod status
- Check sync status
- View deployment history

### 2. **Sync Strategies**
- **Auto-sync**: Changes applied automatically
- **Manual sync**: Require approval
- **Sync waves**: Control deployment order

### 3. **Rollbacks**
- View deployment history
- One-click rollbacks
- Compare configurations

### 4. **Multi-Environment**
- Development (auto-sync)
- Production (manual approval)
- Different configurations per environment

## 🔄 Advanced Testing

### Test Blue-Green Deployments
Update `argocd-applications.yaml` for blue-green:
```yaml
spec:
  syncPolicy:
    automated:
      prune: false
    syncOptions:
    - CreateNamespace=true
    - ApplyOutOfSyncOnly=true
```

### Test Canary Deployments
Add Argo Rollouts:
```bash
kubectl create namespace argo-rollouts
kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml
```

## 🎯 Benefits of This Setup

1. **No Cloud Costs**: Everything runs locally
2. **Real GitOps Experience**: Same workflows as production
3. **Safe Testing**: Can't break anything important
4. **Learning**: Understanding ArgoCD concepts
5. **CI/CD Integration**: Test GitHub Actions + ArgoCD

## 🛠️ Useful Commands

```bash
# ArgoCD CLI (optional)
curl -sSL -o argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
sudo install -m 555 argocd-linux-amd64 /usr/local/bin/argocd

# Login via CLI
argocd login localhost:8080

# Sync application via CLI
argocd app sync bolt-app-dev

# Get application status
argocd app get bolt-app-dev

# View application logs
argocd app logs bolt-app-dev
```

## 🎬 Demo Workflow

1. **Start**: `tilt up` (your existing setup)
2. **Install ArgoCD**: Follow steps above
3. **Create Apps**: Deploy via ArgoCD instead of Tilt
4. **Test Changes**: Push code → Watch ArgoCD sync
5. **Production Sim**: Use manual sync for "production"
6. **Rollback Test**: Revert to previous version

This gives you the full GitOps experience without any cloud dependencies! 🚀
