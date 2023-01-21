
IMAGE_NAME := sample-mock-service
TAG := latest

.PHONY:
make-build:
	docker build -f Dockerfile -t ghcr.io/malike/mock-operator/${IMAGE_NAME}:${TAG} .
	docker push ghcr.io/malike/mock-operator/${IMAGE_NAME}:${TAG}

