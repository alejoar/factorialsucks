APP := factorialsucks
VERSION := latest
EMAIL ?= your@email.com

.PHONY: all

all: build help

build: go_build

go_build:
	go build -o $(APP) .

today_continuous: go_today_continuous

go_today_continuous: go_build
	./$(APP) \
		--email $(EMAIL) \
		--today \
		--clock-in 10:00 \
		--clock-out 18:00

docker_build:
	@docker build \
		--build-arg APP=$(APP) \
		. \
		-t $(APP):$(VERSION)

docker_today_continuous:
	@docker run \
		-it \
		--rm \
		$(APP) \
			--email $(EMAIL) \
			--today \
			--clock-in 7:00 \
			--clock-out 15:00

docker_today_splitshift:
	@docker run \
		-it \
		--rm \
		$(APP) \
			--email $(EMAIL) \
			--today \
			--clock-in 7:00 \
			--clock-out 13:00
	@docker run \
		-it \
		--rm \
		$(APP) \
			--email $(EMAIL) \
			--today \
			--clock-in 14:00 \
			--clock-out 16:00

help:
	@docker run \
		-it \
		--rm \
		$(APP) \
			--help
