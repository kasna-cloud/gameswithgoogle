REPOSITORY = gcr.io/$(PROJECT_ID)


.PHONY: build
build: 
	cd openmatchclient; docker build -t openmatchclient .
	cd openmatchfunction; docker build -t openmatchfunction .

.PHONY: tag
tag:
	docker tag openmatchclient:latest $(REPOSITORY)/openmatchclient:latest
	docker tag openmatchfunction:latest $(REPOSITORY)/openmatchfunction:latest

.PHONY: push
push:
	docker push $(REPOSITORY)/openmatchclient:latest
	docker push $(REPOSITORY)/openmatchfunction:latest

.PHONY: build-new-director
build-new-director: 
	cd openmatchclient; docker build -t openmatchclientagones -f DockerfileDirector .
	docker tag openmatchclientagones:latest $(REPOSITORY)/openmatchclientagones:latest
	docker push $(REPOSITORY)/openmatchclientagones:latest

.PHONY: create-images
create-images: build tag push
