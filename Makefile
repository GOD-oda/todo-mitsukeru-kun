.PHONY: build
build:
	docker compose build --no-cache

.PHONY: deploy
deploy:
	git tag -a -f -m "Description of this release" v0.0.1
	git push -f origin v0.0.1