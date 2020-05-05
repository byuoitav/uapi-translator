# vars
NAME=av-uapi
ORG=byuoitav
BRANCH:= $(shell git rev-parse --abbrev-ref HEAD)

# go
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# docker
DOCKER=docker
DOCKER_BUILD=$(DOCKER) build
DOCKER_LOGIN=$(DOCKER) login -u $(UNAME) -p $(PASS)
DOCKER_PUSH=$(DOCKER) push
DOCKER_FILE=dockerfile

UNAME=$(shell echo $(DOCKER_USERNAME))
EMAIL=$(shell echo $(DOCKER_EMAIL))
PASS=$(shell echo $(DOCKER_PASSWORD))

build: build-x86

build-x86:
	env GOOS=linux CGO_ENABLED=0 $(GOBUILD) -o $(NAME) -v

test:
	$(GOTEST) -v -race $(go list ./...)

clean:
	$(GOCLEAN)
	rm -f $(NAME)

run: $(NAME)
	./$(NAME)

deps:
	$(GOGET) -d -v

docker: docker-x86

docker-x86: $(NAME)
ifeq "$(BRANCH)" "master"
	$(eval BRANCH=production)
endif
	$(DOCKER_BUILD) --build-arg NAME=$(NAME) -f $(DOCKER_FILE) -t $(ORG)/$(NAME):$(BRANCH) .
	@echo logging in to dockerhub...
	@$(DOCKER_LOGIN)
	$(DOCKER_PUSH) $(ORG)/$(NAME):$(BRANCH)
ifeq "$(BRANCH)" "production"
	$(eval BRANCH=master)
endif

### deps
$(NAME):
	$(MAKE) build-x86
