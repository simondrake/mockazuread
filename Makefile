.PHONY: run
run:
	go run .


.PHONY: docker-build
docker-build:
	docker build -f ./Dockerfile -t mockazuread .

.PHONY: docker-tag
docker-tag:
	docker tag mockazuread simondrake/mockazuread:v1alpha1
	docker tag mockazuread simondrake/mockazuread:latest

.PHONY: docker-push
docker-push:
	docker push simondrake/mockazuread:v1alpha1
	docker push simondrake/mockazuread:latest

.PHONY: docker-all
docker-all: docker-build docker-tag docker-push
	@echo "================================================"
	@echo "Docker image has been built, tagged and pushed"
	@echo "================================================"
