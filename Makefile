.PHONY: build

build:
	docker buildx build --platform linux/amd64 -t '您的仓库地址' -f docker/Dockerfile . --push

