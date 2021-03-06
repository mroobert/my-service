SHELL := /bin/bash

# ==============================================================================
# Testing running system

# For testing load on the service.
# hey -m GET -c 100 -n 10000 http://localhost:3000/v1/test

# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem

# Testing Auth
# curl -il http://localhost:3000/v1/testauth
# curl -il -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/testauth

run: 
	go run app/services/sales-api/main.go | go run app/tools/logfmt/main.go

# ==============================================================================
# Administration

admin-genkey: 
	go run app/tools/sales-admin/main.go genkey

admin-gentoken: 
	go run app/tools/sales-admin/main.go gentoken 54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

admin-migrate: 
	go run app/tools/sales-admin/main.go migrate

admin-seed:
	go run app/tools/sales-admin/main.go seed

# ==============================================================================
# Running tests within the local computer

test:
	go test ./... -count=1
	staticcheck -checks=all ./...

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
# Building containers
VERSION := 1.0

all: sales-api

sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within k8s/kind
KIND_CLUSTER := practice-starter-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=sales-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	cd zarf/k8s/kind/sales-pod; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply: 
	kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build zarf/k8s/kind/sales-pod | kubectl apply -f -


kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-sales:
	kubectl get pods -o wide --watch

kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system

kind-logs:
	kubectl logs -l app=sales --all-containers=true -f --tail=100 | go run app/tools/logfmt/main.go

kind-restart:
	kubectl rollout restart deployment sales-depl

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply