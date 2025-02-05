.PHONY: build

build:
	docker buildx build --platform linux/amd64 -t '' -f docker/Dockerfile . --push

