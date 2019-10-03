# use some sensible default shell settings
SHELL := /bin/bash -o pipefail
.SILENT:
.DEFAULT_GOAL := help

REPOSITORY = gcr.io/$(PROJECT_ID)
IMAGE_TAG = $(REPOSITORY)/xonotic-example:0.6

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

.PHONY: destroy
destroy: 
	gcloud container clusters delete agones-cluster --zone=australia-southeast1-a

.PHONY: custom-pools
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

.PHONY: deploy-agones
deploy-agones:
	kubectl create namespace agones-system
	kubectl apply -f https://raw.githubusercontent.com/googleforgames/agones/release-1.0.0/install/yaml/install.yaml

.PHONY: destroy-agones
destroy-agones: 
	kubectl delete -f https://raw.githubusercontent.com/googleforgames/agones/release-1.0.0/install/yaml/install.yaml
	kubectl delete namespace agones-system

.PHONY: deploy-all
deploy-all: cluster deploy-agones deploy-open-match firewall


.PHONY: deploy-open-match
deploy-open-match:
	kubectl create namespace open-match
	kubectl apply -f open-match-example/openmatch-core.yaml --namespace open-match


.PHONY: build-xonotic
build-xonotic:
	cd agones-game-servers/xonotic; docker build --tag=$(IMAGE_TAG) .

.PHONY: push-xonotic
push-xonotic:
	@echo 'publish $(IMAGE_TAG)'
	docker push $(IMAGE_TAG)

.PHONY: deploy-xonotic
deploy-xonotic: #build-xonotic push-xonotic
	kubectl apply -f agones-game-servers/xonotic/gameserver.yaml

.PHONY: deploy-xonotic-fleet
deploy-xonotic-fleet: 
	cd agones-game-servers/xonotic ; kubectl apply -f fleet.yaml	

.PHONY: deploy-xonotic-fleetautoscaler
deploy-xonotic-fleetautoscaler:
	cd agones-game-servers/xonotic ; kubectl apply -f fleetautoscaler.yaml	

.PHONY: deploy-xonotic-allocation
deploy-xonotic-allocation: 
	cd agones-game-servers/xonotic ; kubectl create -f gameserverallocation.yaml

.PHONY: deploy-openmatch-example
deploy-openmatch-example:
	kubectl apply -f open-match-example/openmatchexample.yaml -n open-match

.PHONY: proxy-openmatch-example
proxy-openmatch-example:
	@echo "View Demo: http://localhost:51507"
	kubectl port-forward --namespace open-match $(shell kubectl get pod --namespace open-match --selector="app=openmatchclient,component=frontend,release=open-match" --output jsonpath='{.items[0].metadata.name}') 51507:51507

.PHONY: deploy-allocator-service
deploy-allocator-service:
	kubectl create -f agones-game-servers/allocator-service/allocator-service-accounts.yaml
	kubectl create -f agones-game-servers/allocator-service/allocator-service.yaml

.PHONY: deploy-openmatch-agones-example
deploy-openmatch-agones-example:
	kubectl apply -f open-match-example/openmatchexample-agones.yaml -n open-match

.PHONY: proxy-openmatch-agones-example
proxy-openmatch-agones-example:
	@echo "View Demo: http://localhost:51507"
	kubectl port-forward --namespace open-match $(shell kubectl get pod --namespace open-match --selector="app=openmatchagonesclient,component=frontend,release=open-match" --output jsonpath='{.items[0].metadata.name}') 51507:51507

.PHONY: proxy-grafana-openmatch
proxy-grafana-openmatch:
	@echo "User: admin"
	@echo "Password: openmatch"
	kubectl port-forward --namespace open-match $(shell kubectl get pod --namespace open-match --selector="app=grafana,release=open-match" --output jsonpath='{.items[0].metadata.name}') 3000:3000

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
