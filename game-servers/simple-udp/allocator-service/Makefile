.PHONY: secrets
secrets:
	kubectl create secret tls allocatorsecret --cert=./sec/allocsvc.crt --key=./sec/allocsvc.key
	kubectl create secret tls allocatorw3secret --cert=./sec/tls.crt --key=./sec/tls.key

.PHONY: service
service:
	kubectl create -f service-account.yaml
	kubectl create -f allocator-service.yaml

.PHONY: ingress
ingress:
	kubectl apply -f allocator-ingress.yaml