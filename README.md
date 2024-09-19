# Deployment with `minikube`

## What we will deploy

1. DB Dashboard with Adminer
2. PGSQL
3. The application

## Pre-requisites

1. `minikube`
2. `helm`
3. `nginx-ingress-controller`

# Get started

Start a cluster

```bash
minikube start
```

## Deploying `PGSQL` in Kubernetes

### DB dashboard with Adminer

1. Enable NGINX Ingress controller

```
minikube addons enable ingress
```

2. Tunnel

```
minikube tunnel
```

3. Create resources in `/kubernetes/postgres/adminer`

```bash
kubectl create -f adminer.yaml
kubectl create -f adminer-svc.yaml
kubectl create -f ingress.yaml
```

4. To visit Adminer from browser:
   Add a line to `/etc/hosts`

```
127.0.0.1 arkarsg.split
```

---

### PGSQL StatefulSet

In `kubernetes/postgres`,

1. Create resources

```
k create -f postgres-configmap.yaml
k create -f postgres-pv.yaml
k create -f postgres-pvc.yaml
k create -f postgres-statefulset.yaml
k create -f postgres-service.yaml
```

2. Login to Adminer with username and password

---

## Deploying the application

In `/kubernetes

```
k create -f deployment.yaml
```
