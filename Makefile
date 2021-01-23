GIT_SUMMARY := $(shell git describe --tags --dirty --always)
REPO=jlongman/hacker-slides
DOCKER_IMAGE_APP := $(REPO)-app:$(GIT_SUMMARY)
DOCKER_IMAGE := $(REPO):$(GIT_SUMMARY)
default: repo

repo:
	@echo $(DOCKER_IMAGE_APP)
	@echo $(DOCKER_IMAGE)

build-app:
	@docker build -t $(DOCKER_IMAGE_APP) ./app
	@docker tag $(DOCKER_IMAGE_APP) $(REPO)


build: build-app
	@docker build -t $(DOCKER_IMAGE) .
	@docker tag $(DOCKER_IMAGE) $(REPO)

push:
	@docker push $(DOCKER_IMAGE)

r:
	@docker run -it -p 8080:8080 $(DOCKER_IMAGE)

release:
	@make build
	@make push
