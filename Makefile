IMAGE_NAME=146986152083.dkr.ecr.eu-central-1.amazonaws.com/bingo-bot

.PHONY: build push

build:
	docker build -t ${IMAGE_NAME} .

push: build
	docker push ${IMAGE_NAME}

