# use some sensible default shell settings
SHELL := /bin/bash -o pipefail
.SILENT:
.DEFAULT_GOAL := help

# docker-compose calls
GCLOUD = docker-compose run -w /app gcloud

.PHONY: cluster
cluster: 
	gcloud container clusters create agones-cluster --cluster-version=1.12 \
  		--tags=game-server \
  		--scopes=gke-default \
  		--num-nodes=3 \
		--zone=australia-southeast1-a \
  		--machine-type=n1-standard-4 \
		--enable-autoscaling \
		--max-nodes=10

.PHONY: pools
pools:		
	gcloud container node-pools create agones-system \
  		--cluster=agones-cluster \
  		--node-taints agones.dev/agones-system=true:NoExecute \
  		--node-labels agones.dev/agones-system=true \
  		--num-nodes=1 \
		--zone=australia-southeast1-a 
	gcloud container node-pools create agones-metrics \
		--cluster=agones-cluster \
		--node-taints agones.dev/agones-metrics=true:NoExecute \
		--node-labels agones.dev/agones-metrics=true \
		--num-nodes=1 \
		--zone=australia-southeast1-a
.PHONY: firewall
firewall:
	gcloud compute firewall-rules create game-server-firewall \
		--allow udp:7000-8000 \
		--target-tags game-server \
		--description "Firewall to allow game server udp traffic"

.PHONY: apply
apply:
	kubectl create namespace agones-system
	kubectl apply -f https://raw.githubusercontent.com/googleforgames/agones/master/install/yaml/install.yaml

.PHONY: install-helm
install-helm:
	helm repo add agones https://agones.dev/chart/stable
	#kubectl create namespace xonotic3
	helm install --set "gameservers.namespaces={default,xonotic3}" --namespace agones-system --name agones-install agones/agones

.PHONY: destroy
destroy: 
	kubectl delete -f https://raw.githubusercontent.com/googleforgames/agones/master/install/yaml/install.yaml
	kubectl delete namespace agones-system
##@ Local
.PHONY: shell
shell:
	docker-compose run $(GCLOUD) --entrypoint=sh gcloud

##@ Local
.PHONY: auth
auth: 
	sh "./rbac.sh"

##@ Misc
.PHONY: help
help: ## Display this help
	awk \
	  'BEGIN { \
	    FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n" \
	  } /^[a-zA-Z_-]+:.*?##/ { \
	    printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 \
	  } /^##@/ { \
	    printf "\n\033[1m%s\033[0m\n", substr($$0, 5) \
	  }' $(MAKEFILE_LIST)
