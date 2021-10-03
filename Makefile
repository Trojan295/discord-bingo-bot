IMAGE_NAME=rg.fr-par.scw.cloud/discordbots/bingo-bot:latest

.PHONY: build push

build:
	docker build -t ${IMAGE_NAME} .

push: build
	docker push ${IMAGE_NAME}

