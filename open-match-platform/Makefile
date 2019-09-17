# use some sensible default shell settings
SHELL := /bin/bash -o pipefail
.SILENT:
.DEFAULT_GOAL := help

.PHONY: apply
apply:
	kubectl create namespace open-match
	kubectl apply -f https://open-match.dev/install/v0.6.0/yaml/install.yaml --namespace open-match


.PHONY: port
port:
	kubectl port-forward --namespace open-match service/om-demo 51507:51507