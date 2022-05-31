APP := factorialsucks
VERSION := latest

.PHONY: all

all: build run

build:
	@docker build \
		--build-arg APP=$(APP) \
		. \
		-t $(APP):$(VERSION)

run:
	@docker run \
		-it \
		--rm $(APP) \
		--help
