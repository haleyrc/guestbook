.PHONY: build push release tag

build:
	docker build -t ghcr.io/haleyrc/guestbook .

push:
	docker push ghcr.io/haleyrc/guestbook:latest

tag:
	docker tag ghcr.io/haleyrc/guestbook ghcr.io/haleyrc/guestbook:$(git rev-parse --abbrev-ref HEAD)-$(git rev-parse --short HEAD)

release: build tag push
