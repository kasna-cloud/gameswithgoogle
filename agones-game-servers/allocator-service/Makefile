REPOSITORY = gcr.io/$(PROJECT_ID)


.PHONY: build
build: 
	docker build -t allocator .
.PHONY: tag
tag:
	docker tag allocator:latest $(REPOSITORY)/allocator:latest

.PHONY: push
push:
	docker push $(REPOSITORY)/allocator:latest	

.PHONY: create-image
create-image: build tag push


.PHONY: secrets
secrets:
	kubectl create secret tls allocatorsecret --cert=./sec/allocsvc.crt --key=./sec/allocsvc.key
	kubectl create secret tls allocatorw3secret --cert=./sec/tls.crt --key=./sec/tls.key

.PHONY: service
service:
	kubectl create -f service-account.yaml
	kubectl create -f allocator-service.yaml

.PHONY: create-allocator
create-allocator: create-image secrets service