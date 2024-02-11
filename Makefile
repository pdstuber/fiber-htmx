.DEFAULT_GOAL := build

IMAGE ?= fiber-htmx:latest
PORT := 4000

.PHONY: build
build:
	@docker buildx create --use --name=crossplat --node=crossplat && \
	docker buildx build \
		--progress plain \
		--output "type=docker,push=false" \
		--tag $(IMAGE) \
		--file build/Dockerfile \
		.

.PHONY: run
run:
	docker run -p "${PORT}:${PORT}" "${IMAGE}" run